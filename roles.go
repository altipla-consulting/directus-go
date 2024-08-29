package directus

import (
	"context"
	"fmt"
)

type Role interface{
	ServerVersion() *ServerInfo
}

type roleV10 struct {
	ID          string `json:"id,omitempty"`
	Icon        Icon   `json:"icon,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	AdminAccess *bool `json:"admin_access"`
	AppAccess   *bool `json:"app_access"`

	Users []string `json:"users,omitempty"`

	ServerInfo *ServerInfo
}
func (r *roleV10) ServerVersion() *ServerInfo {
	return r.ServerInfo
}

type roleV11 struct {
	ID          string `json:"id,omitempty"`
	Icon        Icon   `json:"icon,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Policies []string `json:"policies,omitempty"`
	Parent   *string  `json:"parent,omitempty"`
	Children []string `json:"children,omitempty"`

	Users []string `json:"users,omitempty"`

	ServerInfo *ServerInfo
}
func (r *roleV11) ServerVersion() *ServerInfo {
	return r.ServerInfo
}

type clientRoles struct {
	client *Client
}

func (cr *clientRoles) List(ctx context.Context) ([]Role, error) {
	i, err := cr.client.Server.Info(ctx)
	if err != nil {
		return nil, err
	}

	major := i.Major
	switch {

		// Support for Directus 10.0.X - 10.13.X
		case major == 10 && (i.Minor >= 0 && i.Minor <= 13):
			rolesV10, err := cr.client.rolesV10.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("directus: cannot list roles: %w", err)
			}
			roles := make([]Role, len(rolesV10))
        for i, r := range rolesV10 {
            roles[i] = r
        }

        return roles, nil

		// Support for Directus 11.0.X
		case major == 11 && i.Minor == 0:
			rolesV11, err := cr.client.rolesV11.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("directus: cannot list roles: %w", err)
			}
			roles := make([]Role, len(rolesV11))
        for i, r := range rolesV11 {
            roles[i] = r
        }

        return roles, nil
		default:
			return nil, fmt.Errorf("directus: unsupported server version: %v", i)
	}
}
