package directus

import (
	"context"
	"net/http"
)

type clientUsers struct {
	c *Client
}

type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
}

func (cr clientUsers) List(ctx context.Context) ([]User, error) {
	reply := struct {
		Data []User `json:"data"`
	}{}
	if err := cr.c.buildSendRequest(ctx, http.MethodGet, cr.c.urlf("/users"), nil, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr clientUsers) Create(ctx context.Context, user *User) (*User, error) {
	var reply User
	if err := cr.c.buildSendRequest(ctx, http.MethodPost, cr.c.urlf("/users"), user, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
