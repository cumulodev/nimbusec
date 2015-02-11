package nimbusec

import (
	"fmt"
	"net/url"

	"github.com/cumulodev/oauth"
)

const (
	EMPTY_FILTER = ""
	DEFAULT_API  = "https://api.nimbusec.com/"
)

type API struct {
	url    *url.URL
	client *oauth.Consumer
	token  *oauth.AccessToken
}

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

func (a *API) geturl(relpath string, args ...interface{}) string {
	if url, err := a.url.Parse(fmt.Sprintf(relpath, args...)); err == nil {
		return url.String()
	}

	return ""
}
