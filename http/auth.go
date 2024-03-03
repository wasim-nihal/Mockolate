package httpserver

import (
	"flag"
	"net/http"
)

var (
	basicAuthUsername = flag.String("basicauth.username", "", "username for basic authentication")
	basicAuthPassword = flag.String("basicauth.password", "", "password for basic authentication")
)

func CheckBasicAuth(w http.ResponseWriter, r *http.Request) bool {
	if len(*basicAuthUsername) == 0 {
		// HTTP Basic Auth is disabled.
		return true
	}
	username, password, ok := r.BasicAuth()
	if ok {
		if username == *basicAuthUsername && password == *basicAuthPassword {
			return true
		}
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="MockServer"`)
	http.Error(w, "", http.StatusUnauthorized)
	return false
}
