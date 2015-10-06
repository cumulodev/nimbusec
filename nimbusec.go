package nimbusec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

// ErrNotFound is returned by GetXYByName functions if the requested entity can not be found.
var ErrNotFound = errors.New("not found")

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

// get is a helper for all GET request with json payload.
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

// post is a helper for all POST request with json payload.
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

// put is a helper for all PUT request with json payload.
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

// delete is a helper for all DELETE request with json payload.
func (a *API) delete(url string, params params) error {
	resp, err := a.client.Delete(url, params, a.token)
	resp.Body.Close()
	return err
}

// getTextPlain is a helper for all GET request with plain text payload.
func (a *API) getTextPlain(url string, params params) (string, error) {
	data, err := a.getBytes(url, params)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// putTextPlain is a helper for all PUT request with plain text payload.
func (a *API) putTextPlain(url string, params params, payload string) (string, error) {
	resp, err := try(a.client.Put(url, "text/plain", string(payload), params, a.token))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// getBytes is a helper for all GET request with raw byte payload.
func (a *API) getBytes(url string, params params) ([]byte, error) {
	resp, err := try(a.client.Get(url, params, a.token))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
