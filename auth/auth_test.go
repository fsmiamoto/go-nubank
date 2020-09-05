package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	t.Run("valid crendentials", func(t *testing.T) {
		server := buildMockAuthServer("1234", "pass")
		a, err := New(&Config{
			CPF:             "1234",
			Password:        "pass",
			LoginServiceURL: server.URL,
		})
		assert.Nil(t, err)
		assert.Nil(t, a.Login())
		assert.Equal(t, accessTokenFixture, a.AccessToken())
	})

	t.Run("invalid crendentials", func(t *testing.T) {
		server := buildMockAuthServer("1234", "pass")
		a, err := New(&Config{
			CPF:             "1234",
			Password:        "potato",
			LoginServiceURL: server.URL,
		})
		assert.Nil(t, err)
		assert.NotNil(t, a.Login())
	})
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
