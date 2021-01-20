package main

import (
	"log"
	"login"
	"net/http"

	"github.com/codegangsta/negroni"

	"callback"
	"home"
	"logout"
	"middlewares"
	"user"
)

func StartServer() {
	r := http.NewServeMux()

	r.HandleFunc("/", home.HomeHandler)
	r.HandleFunc("/login", login.LoginHandler)
	r.HandleFunc("/logout", logout.LogoutHandler)
	r.HandleFunc("/callback", callback.CallbackHandler)
	r.Handle("/user", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(user.UserHandler)),
	))
	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	n := negroni.New(negroni.NewLogger())
	n.UseHandler(r)
	//http.Handle("/", r)
	http.Handle("/", n)
	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
