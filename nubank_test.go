package nubank_test

import (
	"testing"

	"github.com/fsmiamoto/go-nubank"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("all good", func(t *testing.T) {
		_, err := nubank.New("123.456.789-90", "senhaSecreta123")
		assert.NoError(t, err, "expected no error")
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := nubank.New("123.456.789-90", "senhaNaoTaoSecreta123")
		assert.Error(t, err)
	})
}

func TestLoginWithQRCode(t *testing.T) {
	t.Run("registered qr code", func(t *testing.T) {
		nu, _ := nubank.New("123.456.789-90", "senhaSecreta123")
		err := nu.LoginWithQRCode("some-registered-id")
		assert.NoError(t, err)
	})

	t.Run("unregistered qr code", func(t *testing.T) {
		nu, _ := nubank.New("123.456.789-90", "senhaSecreta123")
		err := nu.LoginWithQRCode("an-unregistered-id")
		assert.Error(t, err)
	})
}
