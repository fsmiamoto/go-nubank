package auth

import (
	"fmt"
	"net/http"
)

type Auth struct {
	client           HTTPClient
	links            map[string]string
	loginURL         string
	liftURL          string
	loginAccessToken string
	accessToken      string
	login            string
	password         string
	qrCodeID         string
}

type Config struct {
	CPF             string
	Password        string
	LoginServiceURL string
	LiftServiceURL  string
}

func New(cfg *Config) (*Auth, error) {
	return &Auth{
		login:    cfg.CPF,
		password: cfg.Password,
		client:   &http.Client{},
		loginURL: cfg.LoginServiceURL,
		liftURL:  cfg.LiftServiceURL,
	}, nil
}

func (a *Auth) AccessToken() string {
	return a.accessToken
}

func (a *Auth) LoginWithQRCode(qrCodeID string) error {
	err := a.requestAccessToken()
	if err != nil {
		return err
	}
	return a.liftQRCodeID(qrCodeID)
}

func (a *Auth) Links() map[string]string {
	return a.links
}

func (a *Auth) liftQRCodeID(qrCodeID string) error {
	req, err := buildLiftRequest(a.liftURL, qrCodeID, a.loginAccessToken)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	res, err := sendRequestToService(a.client, req)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	token, err := getTokenFromLiftResponse(res)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	links, err := getLinksFromLiftResponse(res)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	a.accessToken = token
	a.links = links

	return err
}

func (a *Auth) requestAccessToken() error {
	req, err := buildLoginRequest(a.loginURL, a.login, a.password)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	res, err := sendRequestToService(a.client, req)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	token, err := getTokenFromLoginResponse(res)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	a.loginAccessToken = token

	return nil
}
