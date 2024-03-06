package directus

import (
	"context"
	"net/http"
)

type clientUsers struct {
	c *Client
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
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
