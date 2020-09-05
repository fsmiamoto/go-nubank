package discovery

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
		if err != nil {
			t.Errorf("expected no error but got %q", err)
		}

		url, err := d.ServiceURL("login")
		expected := "url_login"

		if url != expected {
			t.Errorf("expected %q but got %q", expected, url)
		}
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
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
		if err != nil {
			t.Errorf("expected no error but got %q", err)
		}

		_, err = d.ServiceURL("missing")
		if err == nil {
			t.Errorf("expected an error but got none")
		}
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
		if err != nil {
			t.Errorf("expected no error but got %q", err)
		}

		_, err = d.ServiceURL("ios")
		if err == nil {
			t.Errorf("expected an error but got none")
		}
	})

}

func buildMockDiscoveryServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(response))
	}))
}
