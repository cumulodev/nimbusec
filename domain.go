package nimbusec

import "fmt"

// Domain represents a nimbusec monitored domain.
type Domain struct {
	Id        int      `json:"id,omitempty"` // unique identification of domain
	Bundle    string   `json:"bundle"`       // id of assigned bundle
	Name      string   `json:"name"`         // name of domain (usually DNS name)
	Scheme    string   `json:"scheme"`       // whether the domain uses http or https
	DeepScan  string   `json:"deepScan"`     // starting point for the domain deep scan
	FastScans []string `json:"fastScans"`    // landing pages of the domain scanned
}

// CreateDomain issues the API to create the given domain.
func (a *API) CreateDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.geturl("/v2/domain")
	err := a.post(url, params{}, domain, dst)
	return dst, err
}

// CreateOrUpdateDomain issues the nimbusec API to create the given domain. Instead
// of failing when attempting to create a duplicate domain, this method will update
// the remote domain instead.
func (a *API) CreateOrUpdateDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.geturl("/v2/domain")
	err := a.post(url, params{"upsert": "true"}, domain, dst)
	return dst, err
}

// CreateOrGetDomain issues the nimbusec API to create the given domain. Instead
// of failing when attempting to create a duplicate domain, this method will fetch
// the remote domain instead.
func (a *API) CreateOrGetDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.geturl("/v2/domain")
	err := a.post(url, params{"upsert": "false"}, domain, dst)
	return dst, err
}

// GetDomain retrieves a domain from the API by its ID.
func (a *API) GetDomain(domain int) (*Domain, error) {
	dst := new(Domain)
	url := a.geturl("/v2/domain/%d", domain)
	err := a.get(url, params{}, dst)
	return dst, err
}

// GetDomainByName fetches an domain by its name.
func (a *API) GetDomainByName(name string) (*Domain, error) {
	domains, err := a.FindDomains(fmt.Sprintf("name eq \"%s\"", name))
	if err != nil {
		return nil, err
	}

	if len(domains) == 0 {
		return nil, fmt.Errorf("name %q did not match any domains", name)
	}

	if len(domains) > 1 {
		return nil, fmt.Errorf("name %q matched too many domains. please contact nimbusec.", name)
	}

	return &domains[0], nil
}

// FindDomains searches for domains that match the given filter criteria.
func (a *API) FindDomains(filter string) ([]Domain, error) {
	params := params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Domain, 0)
	url := a.geturl("/v2/domain")
	err := a.get(url, params, &dst)
	return dst, err
}

// UpdateDOmain issues the nimbusec API to update a domain.
func (a *API) UpdateDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.geturl("/v2/domain/%d", domain.Id)
	err := a.put(url, params{}, domain, dst)
	return dst, err
}

// DeleteDomain issues the API to delete a domain. When clean=false, the domain and
// all assiciated data will only be marked as deleted, whereas with clean=true the data
// will also be removed from the nimbusec system.
func (a *API) DeleteDomain(d *Domain, clean bool) error {
	url := a.geturl("/v2/domain/%d", d.Id)
	return a.delete(url, params{
		"pleaseremovealldata": fmt.Sprintf("%t", clean),
	})
}

// FindInfected searches for domains that have pending Results that match the
// given filter criteria.
func (a *API) FindInfected(filter string) ([]Domain, error) {
	params := make(map[string]string)
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Domain, 0)
	url := a.geturl("/v2/infected")
	err := a.get(url, params, &dst)
	return dst, err
}
