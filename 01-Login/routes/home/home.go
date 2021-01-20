package home

import (
	"net/http"
	"templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	loginUrl := "/login"
	templates.RenderTemplate(w, "home", loginUrl)
}
