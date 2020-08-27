package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var sucessResponseFixture = []byte(`
	{
	  "access_token": "your_token",
	  "token_type": "bearer",
	  "refresh_token": "string token",
	  "refresh_before": "2020-08-22T22:38:49Z"
	}
`)
var accessTokenFixture = "your_token"

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		server   *httptest.Server
		login    string
		password string
		wantErr  bool
	}{
		{
			name:     "can login with valid credentials",
			server:   buildMockAuthServer("1234", "pass"),
			login:    "1234",
			password: "pass",
			wantErr:  false,
		},
		{
			name:     "cannot login with invalid credentials",
			server:   buildMockAuthServer("1234", "pass"),
			login:    "1234",
			password: "not_pass",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				client:     &http.Client{},
				serviceURL: tt.server.URL,
				login:      tt.login,
				password:   tt.password,
			}

			err := a.Login()
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && a.AccessToken() != accessTokenFixture {
				t.Errorf("Auth.Login() want accessToken to be %q but got %q ", accessTokenFixture, a.AccessToken())
			}
		})
	}
}

func buildMockAuthServer(login, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody map[string]string

		rawBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(rawBody, &requestBody)

		if requestBody["login"] != login || requestBody["password"] != password {
			w.WriteHeader(401)
			w.Write([]byte(`
				{
					"error": "Unauthorized"
				}
			`))
			return
		}

		w.WriteHeader(200)
		w.Write(sucessResponseFixture)
		return
	}))
}
