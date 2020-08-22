package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ErrServiceNotFound = errors.New("discovery: service not found")
var ErrAppServiceNotFound = errors.New("discovery: app service not found")

const DISCOVERY_URL = "https://prod-s0-webapp-proxy.nubank.com.br/api/discovery"
const DISCOVERY_APP_URL = "https://prod-s0-webapp-proxy.nubank.com.br/api/app/discovery"

type Discovery struct {
	client   HTTPClientGet
	services map[string]string
}

func New() (*Discovery, error) {
	d, err := create(&http.Client{})
	if err != nil {
		return nil, fmt.Errorf("discovery: %w", err)
	}
	return d, nil
}

func (d *Discovery) ServiceURL(name string) (string, error) {
	url, ok := d.services[name]
	if !ok {
		return "", ErrServiceNotFound
	}
	return url, nil
}

func create(client HTTPClientGet) (*Discovery, error) {
	d := &Discovery{
		client:   client,
		services: make(map[string]string),
	}

	err := d.fetchServicesFrom(DISCOVERY_URL)
	if err != nil {
		return nil, err
	}

	err = d.fetchServicesFrom(DISCOVERY_APP_URL)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Discovery) fetchServicesFrom(url string) error {
	services, err := sendDiscoveryRequest(d.client, url)
	if err != nil {
		return err
	}

	for service, url := range services {
		d.services[service] = url
	}

	return nil
}

func sendDiscoveryRequest(client HTTPClientGet, url string) (map[string]string, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %v from api, expected 200", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	urlsByService, err := parseDiscoveryResponse(body)
	if err != nil {
		return nil, err
	}

	return urlsByService, nil
}

func parseDiscoveryResponse(response []byte) (map[string]string, error) {
	var unmarshalled interface{}

	err := json.Unmarshal(response, &unmarshalled)
	if err != nil {
		return nil, err
	}

	parsedResponse, ok := unmarshalled.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot unmarshal response from api")
	}

	serviceToURL := make(map[string]string)

	for service, obj := range parsedResponse {
		if url, ok := obj.(string); ok {
			serviceToURL[service] = url
		}
	}

	return serviceToURL, nil
}
