package auth

import (
	"io"
	"net/http"
)

type HTTPClientPost interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}
