package discovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("parses the response correctly", func(t *testing.T) {
		server := buildMockDiscoveryServer(`
		{
			"login":          "url_login",
			"reset_password": "url_reset_password",
			"email_verify":   "url_email_verify_token"
		}
		`)

		d, err := fromClient(&http.Client{}, []string{server.URL})
		assert.Nil(t, err)

		got, err := d.ServiceURL("login")
		expected := "url_login"

		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("missing service", func(t *testing.T) {
		server := buildMockDiscoveryServer(`
		{
			"login":          "url_login",
			"reset_password": "url_reset_password",
			"email_verify":   "url_email_verify_token"
		}
		`)

		d, err := fromClient(&http.Client{}, []string{server.URL})
		assert.Nil(t, err)

		_, err = d.ServiceURL("missing")
		assert.NotNil(t, err)
	})

	t.Run("ignore second level urls", func(t *testing.T) {
		server := buildMockDiscoveryServer(`
		{
			"scopes":       "url_scopes",
			"userinfo":     "url_userinfo",
			"revoke_token": "url_revoke_token",
			"faq": {
				"ios": "url_ios",
				"android": "url_android",
				"wp": "url_windows_phone"
			}
		}
		`)

		d, err := fromClient(&http.Client{}, []string{server.URL})
		assert.Nil(t, err)

		_, err = d.ServiceURL("ios")
		assert.NotNil(t, err)
	})

}

func buildMockDiscoveryServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
}
