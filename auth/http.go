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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type liftRequest struct {
	QRCodeID string `json:"qr_code_id"`
	Type     string `json:"type"`
}

type liftResponse struct {
	AccessToken string                       `json:"access_token"`
	TokenType   string                       `json:"token_type"`
	Links       map[string]map[string]string `json:"_links"`
}

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

func getTokenFromLoginResponse(rawResponse []byte) (string, error) {
	var response loginResponse

	err := json.Unmarshal(rawResponse, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func getTokenFromLiftResponse(rawResponse []byte) (string, error) {
	var response liftResponse

	err := json.Unmarshal(rawResponse, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func getLinksFromLiftResponse(rawResponse []byte) (map[string]string, error) {
	var response liftResponse
	err := json.Unmarshal(rawResponse, &response)
	if err != nil {
		return nil, err
	}

	links := make(map[string]string)

	for service, value := range response.Links {
		links[service] = value["href"]
	}

	return links, nil
}

func buildLoginRequest(serviceURL string, login, password string) (*http.Request, error) {
	body, err := json.Marshal(loginRequest{
		Login:        login,
		Password:     password,
		GrantType:    "password",
		ClientID:     "other.conta",
		ClientSecret: "yQPeLzoHuJzlMMSAjC-LgNUJdUecx8XO",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func sendRequestToService(client HTTPClient, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
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

func buildLiftRequest(serviceURL string, qrCodeID string, token string) (*http.Request, error) {
	body, err := json.Marshal(liftRequest{
		QRCodeID: qrCodeID,
		Type:     "login-webapp",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
