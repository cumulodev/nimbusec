package nimbusec

import "fmt"

const (
	// RoleUser is the restricted role for an user
	RoleUser = "user"

	// RoleAdministrator is the unrestricted role for an user
	RoleAdministrator = "administrator"
)

// User represents an human user able to login and receive notifications.
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

// CreateUser issues the nimbusec API to create the given user.
func (a *API) CreateUser(user *User) (*User, error) {
	dst := new(User)
	url := a.geturl("/v2/user")
	err := a.post(url, params{}, user, dst)
	return dst, err
}

// CreateOrUpdateUser issues the nimbusec API to create the given user. Instead of
// failing when attempting to create a duplicate user, this method will update the
// remote user instead.
func (a *API) CreateOrUpdateUser(user *User) (*User, error) {
	dst := new(User)
	url := a.geturl("/v2/user")
	err := a.post(url, params{"upsert": "true"}, user, dst)
	return dst, err
}

// CreateOrGetUser issues the nimbusec API to create the given user. Instead of
// failing when attempting to create a duplicate user, this method will fetch the
// remote user instead.
func (a *API) CreateOrGetUser(user *User) (*User, error) {
	dst := new(User)
	url := a.geturl("/v2/user")
	err := a.post(url, params{"upsert": "false"}, user, dst)
	return dst, err
}

// GetUser fetches an user by its ID.
func (a *API) GetUser(user int) (*User, error) {
	dst := new(User)
	url := a.geturl("/v2/user/%d", user)
	err := a.get(url, params{}, dst)
	return dst, err
}

// GetUserByLogin fetches an user by its login name.
func (a *API) GetUserByLogin(login string) (*User, error) {
	users, err := a.FindUsers(fmt.Sprintf("login eq \"%s\"", login))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("login %q did not match any users", login)
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("login %q matched too many users. please contact nimbusec.", login)
	}

	return &users[0], nil
}

// FindUsers searches for users that match the given filter criteria.
func (a *API) FindUsers(filter string) ([]User, error) {
	params := params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]User, 0)
	url := a.geturl("/v2/user")
	err := a.get(url, params, dst)
	return dst, err
}

// UpdateUser issues the nimbusec API to update an user.
func (a *API) UpdateUser(user *User) (*User, error) {
	dst := new(User)
	url := a.geturl("/v2/user/%d", user.Id)
	err := a.put(url, params{}, user, dst)
	return dst, err
}

// DeleteUser issues the nimbusec API to delete an user. The root user or tennant
// can not be deleted via this method.
func (a *API) DeleteUser(user *User) error {
	url := a.geturl("/v2/user/%d", user.Id)
	return a.delete(url, params{})
}

// UpdateDomainSet updates the set of allowed domains of an restricted user.
func (a *API) UpdateDomainSet(user *User, domains []int) ([]int, error) {
	dst := make([]int, 0)
	url := a.geturl("/v2/user/%d/domains", user.Id)
	err := a.put(url, params{}, domains, dst)
	return dst, err
}

// LinkDomain links the given domain id to the given user and adds the priviledges for
// the user to view the domain.
func (a *API) LinkDomain(user *User, domain int) error {
	url := a.geturl("/v2/user/%d/domains", user.Id)
	return a.post(url, params{}, domain, nil)
}

// UnlinkDomain unlinks the given domain id to the given user and removes the priviledges
// from the user to view the domain.
func (a *API) UnlinkDomain(user *User, domain int) error {
	url := a.geturl("/v2/user/%d/domains/%d", user.Id, domain)
	return a.delete(url, params{})
}
