package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ErrInvalidCredentials = errors.New("auth: invalid credentials")

const (
	clientID     = "other.conta"
	clientSecret = "yQPeLzoHuJzlMMSAjC-LgNUJdUecx8XO"
	grantType    = "password"
)

type loginRequest struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	GrantType    string `json:"grant_type"`
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
}

type loginResponse struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	RefreshToken  string `json:"refresh_token"`
	RefreshBefore string `json:"refresh_before"`
}

type Auth struct {
	client      HTTPClientPost
	serviceURL  string
	accessToken string
}

func New(serviceURL string) (*Auth, error) {
	return &Auth{
		client:     &http.Client{},
		serviceURL: serviceURL,
	}, nil
}

func (a *Auth) Login(cpf, password string) error {
	requestBody, err := json.Marshal(loginRequest{
		Login:        cpf,
		Password:     password,
		GrantType:    grantType,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	res, err := a.client.Post(a.serviceURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: got status %v from auth service", res.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	token, err := parseAuthResponse(responseBody)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	a.accessToken = token

	return nil
}

func (a *Auth) AccessToken() string {
	return a.accessToken
}

func parseAuthResponse(rawResponse []byte) (string, error) {
	var response loginResponse

	err := json.Unmarshal(rawResponse, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
