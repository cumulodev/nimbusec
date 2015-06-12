package nimbusec

import (
	"fmt"
	"strconv"
	"time"
)

// Domain represents a nimbusec monitored domain.
type Domain struct {
	Id        int      `json:"id,omitempty"` // Unique identification of domain
	Bundle    string   `json:"bundle"`       // ID of assigned bundle
	Name      string   `json:"name"`         // Name of domain (usually DNS name)
	Scheme    string   `json:"scheme"`       // Flag whether the domain uses http or https
	DeepScan  string   `json:"deepScan"`     // Starting point for the domain deep scan
	FastScans []string `json:"fastScans"`    // Landing pages of the domain scanned
}

// DomainBilling represents a billing change event. These happen for example when
// the bundle for a domain changes, or the domain get's disabled or activated.
type DomainBilling struct {
	Time   Timestamp `json:"timestamp"` // Time when the change happend.
	Action string    `json:"action"`    // Type of change, can be either `link` or `unlink`.
	Bundle string    `json:"bundle"`    // ID of bundle this change occured against. Note: The bundle may no longer exist.
	Amount int       `json:"amount"`    // The value of the bundle at the time the change occured.
}

type DomainEvent struct {
	Time    Timestamp `json:"time"`
	Event   string    `json:"event"`
	Human   string    `json:"human"`
	Machine string    `json:"machine"`
}

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(t).Unix()
	stamp := strconv.FormatInt(ts*1000, 10)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(ts/1000, 0))
	return nil
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
		return nil, ErrNotFound
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

// ListDomainConfigs fetches the list of all available configuration keys for the
// given domain.
func (a *API) ListDomainConfigs(domain int) ([]string, error) {
	dst := make([]string, 0)
	url := a.geturl("/v2/domain/%d/config", domain)
	err := a.get(url, params{}, &dst)
	return dst, err
}

// GetDomainConfig fetches the requested domain configuration.
func (a *API) GetDomainConfig(domain int, key string) (string, error) {
	url := a.geturl("/v2/domain/%d/config/%s/", domain, key)
	return a.getTextPlain(url, params{})
}

// SetDomainConfig sets the domain configuration `key` to the requested value.
// This method will create the domain configuration if it does not exist yet.
func (a *API) SetDomainConfig(domain int, key string, value string) (string, error) {
	url := a.geturl("/v2/domain/%d/config/%s/", domain, key)
	return a.putTextPlain(url, params{}, value)
}

// DeleteDomainConfig issues the API to delete the domain configuration with
// the provided key.
func (a *API) DeleteDomainConfig(domain int, key string) error {
	url := a.geturl("/v2/domain/%d/config/%s/", domain, key)
	return a.delete(url, params{})
}

// GetDomainBilling gets the billing change log for the given domain. The returned
// list is sorted by time descending, where up to `limit` items will be returned.
func (a *API) GetDomainBilling(domain int, limit int) ([]DomainBilling, error) {
	dst := make([]DomainBilling, 0)
	url := a.geturl("/v2/domain/%d/billing", domain)
	err := a.get(url, params{"limit": strconv.Itoa(limit)}, &dst)
	return dst, err
}

func (a *API) GetDomainEvent(domain int, filter string, limit int) ([]DomainEvent, error) {
	params := params{
		"limit": strconv.Itoa(limit),
	}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]DomainEvent, 0)
	url := a.geturl("/v2/domain/%d/events", domain)
	err := a.get(url, params, &dst)
	return dst, err
}

func (a *API) CreateDomainEvent(domain int, log *DomainEvent) error {
	url := a.geturl("/v2/domain/%d/events", domain)
	return a.post(url, params{}, log, nil)
}
