package nimbusec

import (
	"encoding/json"
	"fmt"
)

type Domain struct {
	Id        int      `json:"id,omitempty"` // unique identification of domain
	Bundle    string   `json:"bundle"`       // id of assigned bundle
	Name      string   `json:"name"`         // name of domain (usually DNS name)
	Scheme    string   `json:"scheme"`       // whether the domain uses http or https
	DeepScan  string   `json:"deepScan"`     // starting point for the domain deep scan
	FastScans []string `json:"fastScans"`    // landing pages of the domain scanned
}

func (a *API) CreateDomain(domain *Domain) (*Domain, error) {
	payload, err := json.Marshal(domain)
	if err != nil {
		return nil, err
	}

	param := make(map[string]string)
	url := a.geturl("/v2/domain")
	resp, err := a.client.Post(url, "application/json", string(payload), param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(Domain)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) GetDomain(domain int) (*Domain, error) {
	param := make(map[string]string)
	url := a.geturl("/v2/domain/%d", domain)
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(Domain)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) FindDomains(filter string) ([]Domain, error) {
	param := make(map[string]string)
	if filter != EMPTY_FILTER {
		param["q"] = filter
	}

	url := a.geturl("/v2/domain")
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := make([]Domain, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) DeleteDomain(d *Domain, clean bool) error {
	url := a.geturl("/v2/domain/%d", d.Id)
	_, err := a.client.Delete(url, map[string]string{
		"pleaseremovealldata": fmt.Sprintf("%t", clean),
	}, a.token)
	return err
}

func (a *API) FindInfected(filter string) ([]Domain, error) {
	param := make(map[string]string)
	if filter != EMPTY_FILTER {
		param["q"] = filter
	}

	url := a.geturl("/v2/infected")
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := make([]Domain, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
