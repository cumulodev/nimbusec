package nimbusec

import "encoding/json"

type Result struct {
	// unique identification of a result
	Id int `json:"id,omitempty`
	// status of the result (pending, acknowledged, falsepositive, removed)
	Status string `json:"status"`
	// event type of result (e.g added file)
	Event string `json:"status"`
	// category of result
	Category string `json:"status"`
	// severity level of result (1 = medium to 3 = severe)
	Severity int `json:"severity"`
	// probability the result is critical
	Probability float64 `json:"probability`
	// flag indicating if the file can be safely deleted without loosing user data
	SafeToDelete bool `json:"safeToDelete"`
	// timestamp (in ms) of the first occurrence
	CreateDate int `json:"createDate"`
	// timestamp (in ms) of the last occurrence
	LastDate int `json:"lastDate"`

	// the following fields contain more details about the result. Not all fields
	// must be filled or present.

	// name identifying the threat of a result. meaning differs per category:
	// malware & webshell: the virus database name of the malicious software
	// blacklist: the name of the blacklist containing the domain
	Threatname string `json:"threatname"`
	// affected resource (e.g. file path or URL)
	Resource string `json:"resource"`
	// MD5 hash sum of the affected file
	MD5 string `json:"md5"`
	// filesize of the affected file
	Filesize int `json:"filesize"`
	// file owner of the affected file
	Owner string `json:"owner"`
	// file group of the affected file
	Group string `json:"group"`
	// permission of the affected file as decimal integer
	Permission int `json:"permission"`
	// diff of a content change between two scans
	Diff string `json:"diff"`
	// reason why a domain/URL is blacklisted
	Reason string `json:"reason"`
}

func (a *API) GetResult(domain, result int) (*Result, error) {
	param := make(map[string]string)
	url := a.geturl("/v2/domain/%d/result/%d", domain, result)
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(Result)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) FindResults(domain int, filter string) ([]Result, error) {
	param := make(map[string]string)
	if filter != EMPTY_FILTER {
		param["q"] = filter
	}

	url := a.geturl("/v2/domain/%d/result", domain)
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := make([]Result, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
