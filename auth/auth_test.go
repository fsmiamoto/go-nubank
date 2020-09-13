package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginWithQRCode(t *testing.T) {
	loginServer := buildMockLoginService("1234", "pass")
	liftServer := buildMockLiftService("aaaa-bbbb-cccc", "your_token")

	t.Run("invalid login credentials", func(t *testing.T) {
		a, _ := New(&Config{
			CPF:             "1234",
			Password:        "not-pass",
			LoginServiceURL: loginServer.URL,
			LiftServiceURL:  liftServer.URL,
		})

		assert.Error(t, ErrInvalidCredentials, a.LoginWithQRCode("aaaa-bbbb-cccc"))
	})

	t.Run("registered qr code", func(t *testing.T) {
		a, _ := New(&Config{
			CPF:             "1234",
			Password:        "pass",
			LoginServiceURL: loginServer.URL,
			LiftServiceURL:  liftServer.URL,
		})

		assert.NoError(t, a.LoginWithQRCode("aaaa-bbbb-cccc"))

		assert.Equal(t, "your_lift_token", a.AccessToken())
	})

	t.Run("unregistered qr code", func(t *testing.T) {
		a, _ := New(&Config{
			CPF:             "1234",
			Password:        "pass",
			LoginServiceURL: loginServer.URL,
			LiftServiceURL:  liftServer.URL,
		})

		assert.Error(t, assert.AnError, a.LoginWithQRCode("random-id"))
	})
}

func buildMockLiftService(qrCodeID string, token string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody liftRequest

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(400)
			w.Write([]byte(`{"error"::"(not (some-matching-condition? nil))"}`))
			return
		}

		rawBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(rawBody, &requestBody)

		if r.Header.Get("Authorization") != "Bearer "+token {
			w.WriteHeader(401)
			w.Write([]byte(`
				{
					"error": "Unauthorized"
				}
			`))
			return
		}

		if requestBody.QRCodeID != qrCodeID {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(`
		{
		  "access_token": "your_lift_token",
		  "token_type": "bearer"
		}
		`))
	}))
}

func buildMockLoginService(login, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody loginRequest

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(400)
			w.Write([]byte(`{"error"::"(not (some-matching-condition? nil))"}`))
			return
		}

		rawBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(rawBody, &requestBody)

		if requestBody.Login != login || requestBody.Password != password {
			w.WriteHeader(401)
			w.Write([]byte(`
			{
				"error": "Unauthorized"
			}
			`))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(`
		{
		  "access_token": "your_token",
		  "token_type": "bearer",
		  "refresh_token": "string token",
		  "refresh_before": "2020-08-22T22:38:49Z"
		}
		`))
	}))
}
