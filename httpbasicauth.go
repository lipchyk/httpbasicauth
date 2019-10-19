package httpbasicauth

import (
	"fmt"
	"net/http"
)

type Checker interface {
	Check(username, password string) bool
}

// SimpleCredentialMap wraps map[string]string
type SimpleCredentialMap map[string]string

func (c SimpleCredentialMap) Check(u, p string) bool {
	password, ok := c[u]
	if !ok {
		return false
	}

	return password == p
}

// Handle is a middleware that wraps your handler with http basic auth functionality
func Handle(c Checker, realm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				u, p := extractCreds(r)
				if !c.Check(u, p) {
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

func extractCreds(r *http.Request) (username, password string) {
	username, password, _ = r.BasicAuth()
	return
}
