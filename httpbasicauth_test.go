package httpbasicauth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestHandle(t *testing.T) {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "YOLO!")
		},
	)

	creds := map[string]string{"foo": "bar"}
	authhandler := Handle(creds, "Restricted zone")(handler)

	tests := []struct {
		name                 string
		user                 string
		password             string
		expectedResponseCode int
		expectedResponseBody string
		expectedRealm        string
	}{
		{
			name:                 "auth with wrong user",
			user:                 "f00",
			password:             "bar",
			expectedResponseCode: 401,
			expectedRealm:        "Restricted zone",
			expectedResponseBody: "401 Unauthorized",
		},
		{
			name:                 "auth with wrong password",
			user:                 "foo",
			password:             "b@r",
			expectedResponseCode: 401,
			expectedRealm:        "Restricted zone",
			expectedResponseBody: "401 Unauthorized",
		},
		{
			name:                 "auth without creds",
			expectedResponseCode: 401,
			expectedRealm:        "Restricted zone",
			expectedResponseBody: "401 Unauthorized",
		},
		{
			name:                 "auth with correct creds",
			user:                 "foo",
			password:             "bar",
			expectedResponseCode: 200,
			expectedResponseBody: "YOLO!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.user != "" {
				req.SetBasicAuth(tt.user, tt.password)
			}

			rr := httptest.NewRecorder()
			authhandler.ServeHTTP(rr, req)

			// check realm
			if rr.Code != http.StatusOK {
				wwwauth := getrealm(rr)
				if wwwauth == "" || wwwauth != tt.expectedRealm {
					t.Fatalf("Unexpected realm, got %v, expected = %v", wwwauth, tt.expectedRealm)
				}
			}

			if tt.expectedResponseCode != rr.Code {
				t.Fatalf("Wrong status code, got %v, expected = %v", rr.Code, tt.expectedResponseCode)
			}

			responseBody := strings.TrimSpace(rr.Body.String())

			if responseBody != tt.expectedResponseBody {
				t.Errorf("Unexpected body: got %v, want %v", responseBody, tt.expectedResponseBody)
			}
		})
	}
}

func getrealm(r http.ResponseWriter) string {
	re := regexp.MustCompile(`^Basic realm="(.+?)"$`)
	wwwauth := r.Header().Get("WWW-Authenticate")
	return re.FindStringSubmatch(wwwauth)[1]
}
