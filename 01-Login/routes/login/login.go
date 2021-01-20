package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"app"
	"auth"

	log "github.com/sirupsen/logrus"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state; we keep in session store and give it to client.
	// if they don't include it in their response, we reject
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	// should create or open sess...  kept getting "no such
	// file or directory", we don't care about the error - we
	// are overwriting any existing session - gorilla/sessions
	// github explains this behavior in example
	session, _ := app.Store.Get(r, "auth-session")
	// escape it so we guarantee it is url encoded when we send it
	session.Values["state"] = url.QueryEscape(state)
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
