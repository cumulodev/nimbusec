package nimbusec

import "encoding/json"

const (
	RoleUser          = "user"
	RoleAdministrator = "administrator"
)

type User struct {
	Id           int    `json:"id,omitempty"`           // unique identification of user
	Login        string `json:"login"`                  // login name of user
	Mail         string `json:"mail"`                   // e-mail contact where mail notifications are sent to
	Role         string `json:"role"`                   // role of an user (`administrator` or `user`
	Company      string `json:"company"`                // company name of user
	Surname      string `json:"surname"`                // surname of user
	Forename     string `json:"forename"`               // forename of user
	Title        string `json:"title"`                  // academic title of user
	Mobile       string `json:"mobile"`                 // phone contact where sms notifications are sent to
	Password     string `json:"password,omitempty"`     // password of user (only used when creating or updating a user)
	SignatureKey string `json:"signatureKey,omitempty"` // secret for SSO (only used when creating or updating a user)
}

func (a *API) CreateUser(user *User) (*User, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	param := make(map[string]string)
	url := a.geturl("/v2/user")
	resp, err := a.client.Post(url, "application/json", string(payload), param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(User)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) GetUser(user int) (*User, error) {
	param := make(map[string]string)
	url := a.geturl("/v2/user/%d", user)
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := new(User)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) FindUsers(filter string) ([]User, error) {
	param := make(map[string]string)
	if filter != EMPTY_FILTER {
		param["q"] = filter
	}

	url := a.geturl("/v2/user")
	resp, err := a.client.Get(url, param, a.token)
	if err != nil {
		return nil, err
	}

	body := make([]User, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *API) DeleteUser(user *User) error {
	param := make(map[string]string)
	url := a.geturl("/v2/user/%d", user.Id)
	_, err := a.client.Delete(url, param, a.token)
	return err
}
