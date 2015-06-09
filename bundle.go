package nimbusec

type Bundle struct {
	Id         int       `json:"id,omitempty"`
	Name       string    `json:"name"`
	Start      Timestamp `json:"startDate"`
	End        Timestamp `json:"endDate"`
	Quota      string    `json:"quota"`
	Depth      int       `json:"depth"`
	Fast       int       `json:"fast"`
	Deep       int       `json:"deep"`
	Contingent int       `json:"contingent"`
	Active     int       `json:"active"`
	Engines    []string  `json:"engines"`
	Amount     int       `json:"amount"`
	Currency   string    `json:"currency"`
}

func (a *API) GetBundle(bundle string) (*Bundle, error) {
	dst := new(Bundle)
	url := a.geturl("/v2/bundle/%s", bundle)
	err := a.get(url, params{}, dst)
	return dst, err
}

func (a *API) FindBundles(filter string) ([]Bundle, error) {
	params := params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Bundle, 0)
	url := a.geturl("/v2/bundle")
	err := a.get(url, params, &dst)
	return dst, err
}
