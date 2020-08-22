package discovery

import "net/http"

type HTTPClientGet interface {
	Get(url string) (*http.Response, error)
}
