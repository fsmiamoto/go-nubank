package auth

import (
	"fmt"
	"net/http"
)

type Auth struct {
	client      HTTPClientPost
	serviceURL  string
	accessToken string
	login       string
	password    string
}

type Config struct {
	ServiceURL string
	CPF        string
	Password   string
}

func New(cfg *Config) (*Auth, error) {
	return &Auth{
		client:     &http.Client{},
		serviceURL: cfg.ServiceURL,
		login:      cfg.CPF,
		password:   cfg.Password,
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

	response, err := sendRequestToService(a.client, a.serviceURL, request)
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
