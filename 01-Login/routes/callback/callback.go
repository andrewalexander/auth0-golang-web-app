package callback

import (
	"auth"
	"context"
	"net/http"
	"net/url"
	"os"

	"app"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	// '+' to ' ' wasn't working correctly.. perform query unescape to
	// normalize state variable for both inputs
	inState, err := url.QueryUnescape(r.URL.Query().Get("state"))
	if err != nil {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		log.Errorf("error unescaping provided state variable. %s", err)
		return
	}
	s := session.Values["state"].(string)
	ss, err := url.QueryUnescape(s)
	if err != nil {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		log.Errorf("error unescaping stored state variable. %s", err)
		return
	}
	storeState, _ := url.QueryUnescape(ss)
	if inState != storeState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		log.Errorf("mismatching state parameters.\n%s\n%s", inState, storeState)
		return
	}

	// use go-oidc wrapped in authenticator to exchange response code for token
	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Errorf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	// verify it
	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("AUTH_CLIENT_ID"),
	}
	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// get our claims (userInfo) out of the token
	profile := map[string]interface{}{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	// store anything we want to persist as part of this session - csrf tokens
	// should be added here. keep it small though.. 4096 limit or we end up
	// with `securecookie: the value is too long` returned when we go to Save()
	email := profile["email"].([]byte)
	session.Values["email"] = email
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("could not save session. %s", err)
		return
	}
	// quick and dirty serialization of our data into bytes buffer
	//u := app.User{
	//Email: string(email),
	//}
	//buf, _ := json.Marshal(u)
	//binary.Write(&buf, binary.BigEndian, u)
	//if err := app.DB.Update(func(tx *bolt.Tx) error {
	//b := tx.Bucket([]byte("users"))
	//return b.Put(email, buf)
	//}); err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//log.Error("could not persist session and user state. %s", err)
	//return
	//}

	// Redirect to logged in page
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}
