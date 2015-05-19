package nimbusec

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cumulodev/oauth"
)

const (
	// EmptyFilter is a filter that matches all fields.
	EmptyFilter = ""

	// DefaultAPI is the default endpoint of the nimbusec API.
	DefaultAPI = "https://api.nimbusec.com/"
)

// API represents a client to the nimbusec API.
type API struct {
	url    *url.URL
	client *oauth.Consumer
	token  *oauth.AccessToken
}

// NewAPI creates a new nimbusec API client.
func NewAPI(rawurl, key, secret string) (*API, error) {
	client := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
	token := &oauth.AccessToken{}

	parsed, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	return &API{
		url:    parsed,
		client: client,
		token:  token,
	}, nil
}

// geturl builds the fully qualified url to the nimbusec API.
func (a *API) geturl(relpath string, args ...interface{}) string {
	if url, err := a.url.Parse(fmt.Sprintf(relpath, args...)); err == nil {
		return url.String()
	}

	return ""
}

// try is used to encapsulate a HTTP operation and retrieve the optional
// nimbusec error if one happened.
func try(resp *http.Response, err error) (*http.Response, error) {
	if resp == nil {
		return resp, err
	}

	if resp.StatusCode < 300 {
		return resp, err
	}

	msg := resp.Header.Get("x-nimbusec-error")
	if msg != "" {
		return resp, errors.New(msg)
	}

	return resp, err
}

type params map[string]string

func (a *API) post(url string, params params, src interface{}, dst interface{}) error {
	payload, err := json.Marshal(src)
	if err != nil {
		return err
	}

	resp, err := try(a.client.Post(url, "application/json", string(payload), params, a.token))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// no destination, so caller was only interested in the
	// side effects.
	if dst == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dst)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) get(url string, params params, dst interface{}) error {
	resp, err := try(a.client.Get(url, params, a.token))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// no destination, so caller was only interested in the
	// side effects.
	if dst == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dst)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) put(url string, params params, src interface{}, dst interface{}) error {
	payload, err := json.Marshal(src)
	if err != nil {
		return err
	}

	resp, err := try(a.client.Put(url, "application/json", string(payload), params, a.token))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// no destination, so caller was only interested in the
	// side effects.
	if dst == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dst)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) delete(url string, params params) error {
	resp, err := a.client.Delete(url, params, a.token)
	resp.Body.Close()
	return err
}
