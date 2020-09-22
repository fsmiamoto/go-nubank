package nubank

import (
	"github.com/fsmiamoto/go-nubank/auth"
	"github.com/fsmiamoto/go-nubank/discovery"
)

type Nubank struct {
	auth      *auth.Auth
	discovery *discovery.Discovery
	login     string
	password  string
}

func New(cpf, password string) (*Nubank, error) {
	return nil, nil
}

func (nu *Nubank) LoginWithQRCode(qrCodeID string) error {
	return nil
}
