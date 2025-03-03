package linodego

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// User represents a User object
type User struct {
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Restricted bool     `json:"restricted"`
	TFAEnabled bool     `json:"tfa_enabled"`
	SSHKeys    []string `json:"ssh_keys"`
}

// UserCreateOptions fields are those accepted by CreateUser
type UserCreateOptions struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Restricted bool   `json:"restricted"`
}

// UserUpdateOptions fields are those accepted by UpdateUser
type UserUpdateOptions struct {
	Username   string `json:"username,omitempty"`
	Restricted *bool  `json:"restricted,omitempty"`
}

// GetCreateOptions converts a User to UserCreateOptions for use in CreateUser
func (i User) GetCreateOptions() (o UserCreateOptions) {
	o.Username = i.Username
	o.Email = i.Email
	o.Restricted = i.Restricted

	return
}

// GetUpdateOptions converts a User to UserUpdateOptions for use in UpdateUser
func (i User) GetUpdateOptions() (o UserUpdateOptions) {
	o.Username = i.Username
	o.Restricted = copyBool(&i.Restricted)

	return
}

// UsersPagedResponse represents a paginated User API response
type UsersPagedResponse struct {
	*PageOptions
	Data []User `json:"data"`
}

// endpoint gets the endpoint URL for User
func (UsersPagedResponse) endpoint(c *Client, _ ...any) string {
	endpoint, err := c.Users.Endpoint()
	if err != nil {
		panic(err)
	}

	return endpoint
}

func (resp *UsersPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(UsersPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*UsersPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListUsers lists Users on the account
func (c *Client) ListUsers(ctx context.Context, opts *ListOptions) ([]User, error) {
	response := UsersPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetUser gets the user with the provided ID
func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	e, err := c.Users.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%s", e, id)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&User{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*User), nil
}

// CreateUser creates a User.  The email address must be confirmed before the
// User account can be accessed.
func (c *Client) CreateUser(ctx context.Context, createOpts UserCreateOptions) (*User, error) {
	var body string

	e, err := c.Users.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&User{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*User), nil
}

// UpdateUser updates the User with the specified id
func (c *Client) UpdateUser(ctx context.Context, id string, updateOpts UserUpdateOptions) (*User, error) {
	var body string

	e, err := c.Users.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%s", e, id)

	req := c.R(ctx).SetResult(&User{})

	if bodyData, err := json.Marshal(updateOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Put(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*User), nil
}

// DeleteUser deletes the User with the specified id
func (c *Client) DeleteUser(ctx context.Context, id string) error {
	e, err := c.Users.Endpoint()
	if err != nil {
		return err
	}

	e = fmt.Sprintf("%s/%s", e, id)

	_, err = coupleAPIErrors(c.R(ctx).Delete(e))

	return err
}
