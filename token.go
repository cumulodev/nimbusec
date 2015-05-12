package nimbusec

import "encoding/json"

// Token represents the credentials of an API or agent for the nimbusec API.
type Token struct {
	Id       int    `json:"id"`       // unique identification of a token
	Name     string `json:"name"`     // given name for a token
	Key      string `json:"key"`      // oauth key
	Secret   string `json:"secret"`   // oauth secret
	LastCall int    `json:"lastCall"` // last timestamp (in ms) an agent used the token
	Version  int    `json:"version"`  // last agent version that was seen for this key
}

// CreateToken issues the nimbusec API to create a new agent token.
func (a *API) CreateToken(token *Token) (*Token, error) {
	payload, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	param := make(map[string]string)
	url := a.geturl("/v2/agent/token")
	resp, err := a.client.Post(url, "application/json", string(payload), param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(Token)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetToken fetches a token by its ID.
func (a *API) GetToken(token int) (*Token, error) {
	param := make(map[string]string)
	url := a.geturl("/v2/agent/token/%d", token)
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(Token)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// FindTOkens searches for tokens that match the given filter criteria.
func (a *API) FindTokens(filter string) ([]Token, error) {
	param := make(map[string]string)
	if filter != EmptyFilter {
		param["q"] = filter
	}

	url := a.geturl("/v2/agent/token")
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := make([]Token, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
