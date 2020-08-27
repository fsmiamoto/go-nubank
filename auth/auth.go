package auth

import (
	"fmt"
	"net/http"
)

type Auth struct {
	client      HTTPClientPost
	loginURL    string
	liftURL     string
	accessToken string
	login       string
	password    string
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

func (a *Auth) Login() error {
	return a.requestAccessToken()
}

func (a *Auth) AccessToken() string {
	return a.accessToken
}

func (a *Auth) requestAccessToken() error {
	request, err := buildLoginRequestBody(a.login, a.password)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	response, err := sendRequestToService(a.client, a.loginURL, request)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	token, err := getTokenFromResponse(response)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	a.accessToken = token

	return nil
}
