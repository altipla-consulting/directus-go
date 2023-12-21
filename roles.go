package directus

import (
	"context"
	"net/http"
)

type clientRoles struct {
	c *Client
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminAccess bool   `json:"admin_access"`
}

func (cr clientRoles) List(ctx context.Context) ([]Role, error) {
	reply := struct {
		Data []Role `json:"data"`
	}{}
	if err := cr.c.buildSendRequest(ctx, http.MethodGet, cr.c.urlf("/roles"), nil, &reply); err != nil {
		return nil, err
	}

	return reply.Data, nil
}
