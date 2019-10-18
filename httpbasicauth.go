package httpbasicauth

import (
	"fmt"
	"net/http"
)

// Handle is a middleware that wraps your handler with http basic auth functionality
func Handle(creds map[string]string, realm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if !checkBasicAuth(creds, r) {
					unauthorized(w, realm)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func checkBasicAuth(creds map[string]string, r *http.Request) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}

	if password, ok := creds[u]; ok {
		return p == password
	}

	return false
}
