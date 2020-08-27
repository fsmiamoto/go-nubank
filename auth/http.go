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

func getTokenFromResponse(rawResponse []byte) (string, error) {
	var response loginResponse

	err := json.Unmarshal(rawResponse, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func buildLoginRequestBody(login, password string) ([]byte, error) {
	return json.Marshal(loginRequest{
		Login:        login,
		Password:     password,
		GrantType:    "password",
		ClientID:     "other.conta",
		ClientSecret: "yQPeLzoHuJzlMMSAjC-LgNUJdUecx8XO",
	})
}

func sendRequestToService(client HTTPClientPost, url string, body []byte) ([]byte, error) {
	res, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidCredentials
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %v from service", res.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
