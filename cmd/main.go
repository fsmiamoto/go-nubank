package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fsmiamoto/go-nubank/auth"
	"github.com/fsmiamoto/go-nubank/discovery"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

var ErrTokenNotFound = errors.New("token not found")
var tokenFilePath = xdg.ConfigHome + "/go-nubank/token"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Trying to read existing token...")
	token, err := readExistingAccessToken()
	if errors.Is(err, ErrTokenNotFound) {
		fmt.Println("Token not found, fetching new one...")
		token, err = fetchNewAccessToken()
	}

	if err != nil {
		return err
	}

	fmt.Println("Token: ", token)

	return nil
}

func fetchNewAccessToken() (string, error) {
	d, err := discovery.New()
	if err != nil {
		return "", err
	}

	loginServiceURL, err := d.ServiceURL("login")
	if err != nil && !errors.Is(err, discovery.ErrServiceNotFound) {
		return "", err
	}

	liftServiceURL, err := d.ServiceURL("lift")
	if err != nil && !errors.Is(err, discovery.ErrServiceNotFound) {
		return "", err
	}

	id := uuid.New()

	qr, err := qrcode.New(id.String(), qrcode.Medium)
	if err != nil {
		return "", err
	}

	var cpf, password string

	fmt.Printf("Insert your login information: \n")
	fmt.Printf("CPF: ")
	fmt.Scanf("%s", &cpf)
	fmt.Printf("Password: ")
	fmt.Scanf("%s", &password)

	fmt.Printf("Scan the following QR Code on the Nubank app: \n")
	fmt.Println(qr.ToSmallString(false))

	fmt.Println("Press any key to continue after scanning...")
	fmt.Scanln()

	a, err := auth.New(&auth.Config{
		LoginServiceURL: loginServiceURL,
		LiftServiceURL:  liftServiceURL,
		CPF:             cpf,
		Password:        password,
	})
	if err != nil {
		return "", err
	}

	err = a.LoginWithQRCode(id.String())
	if err != nil {
		return "", err
	}

	return a.AccessToken(), nil
}

func saveAccessToken(token string) error {
	trimmed := strings.TrimSpace(token)
	return ioutil.WriteFile(tokenFilePath, []byte(trimmed), 0600)
}

func readExistingAccessToken() (string, error) {
	f, err := os.Open(tokenFilePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrTokenNotFound
		}
		return "", err
	}

	token, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(token), nil
}
