package user

import (
	"net/http"

	"app"
	"templates"

	log "github.com/sirupsen/logrus"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	email := session.Values["email"].(string)
	user := app.User{
		Email: email,
	}
	//if err := app.DB.View(func(tx *bolt.Tx) error {
	//b := tx.Bucket([]byte("users"))
	//// quick and dirty deserialization of our data into user
	//raw := b.Get(email)
	//return binary.Read(bytes.NewReader(raw), binary.BigEndian, &user)
	//}); err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//log.Error(err)
	//return
	//}

	templates.RenderTemplate(w, "user", user)
}
