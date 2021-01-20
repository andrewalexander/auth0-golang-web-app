package logout

import (
	"net/http"
	"net/url"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	//domain := os.Getenv("AUTH0_DOMAIN")
	//domain := "https://gotem.auth.us-east-1.amazoncognito.com/logout?client_id=2htjam9t68mkehjg06j0hbb929&logout_uri=http://localhost:4242"
	domain := "https://gotem.auth.us-east-1.amazoncognito.com"

	logoutUrl, err := url.Parse(domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parameters.Add("logout_uri", returnTo.String())
	//parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	parameters.Add("client_id", "2htjam9t68mkehjg06j0hbb929")
	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
