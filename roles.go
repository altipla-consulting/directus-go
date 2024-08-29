package directus

import (
	"context"
	"fmt"
	"sync"
)

type Role interface {
	PoliciesKeys() ([]string, error)
	ParentKey() (string, error)
	ChildrenKeys() ([]string, error)
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

func (r *roleV10) PoliciesKeys() ([]string, error) {
	return nil, fmt.Errorf("directus: Policies are not supported in Directus below 11.0.0")
}

func (r *roleV10) ParentKey() (string, error) {
	return "", fmt.Errorf("directus: Parent Role is not supported in Directus below 11.0.0")
}

func (r *roleV10) ChildrenKeys() ([]string, error) {
	return nil, fmt.Errorf("directus: Children Roles are not supported in Directus below 11.0.0")
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

func (r *roleV11) PoliciesKeys() ([]string, error) {
	return r.Policies, nil
}

func (r *roleV11) ParentKey() (string, error) {
	if r.Parent == nil {
		return "", nil
	}
	return *r.Parent, nil
}

func (r *roleV11) ChildrenKeys() ([]string, error) {
	return r.Children, nil
}

type clientRoles struct {
	client *Client
	mu     sync.Mutex
}

func (cr *clientRoles) List(ctx context.Context) ([]Role, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	inf, err := cr.client.Server.Info(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	// Support for Directus 10.0.X - 10.13.X
	case inf.VersionMajor == 10 && (inf.VersionMinor >= 0 && inf.VersionMinor <= 13):
		rolesV10, err := cr.client.rolesV10.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("directus: cannot list roles: %w", err)
		}
		roles := make([]Role, len(rolesV10))
		for i, r := range rolesV10 {
			r.ServerInfo = inf
			roles[i] = r
		}

		return roles, nil

	// Support for Directus 11.0.X
	case inf.VersionMajor == 11 && inf.VersionMinor == 0:
		rolesV11, err := cr.client.rolesV11.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("directus: cannot list roles: %w", err)
		}
		roles := make([]Role, len(rolesV11))
		for i, r := range rolesV11 {
			r.ServerInfo = inf
			roles[i] = r
		}
		return roles, nil

	default:
		return nil, fmt.Errorf("directus: unsupported server version: %v", inf)
	}
}

func (cr *clientRoles) Get(ctx context.Context, id string) (Role, error) {
	cr.mu.Lock()
  defer cr.mu.Unlock()

	inf, err := cr.client.Server.Info(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	// Support for Directus 10.0.X - 10.13.X
	case inf.VersionMajor == 10 && (inf.VersionMinor >= 0 && inf.VersionMinor <= 13):
		roleV10, err := cr.client.rolesV10.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("directus: cannot get role: %w", err)
		}
		roleV10.ServerInfo = inf
		return roleV10, nil

	// Support for Directus 11.0.X
	case inf.VersionMajor == 11 && inf.VersionMinor == 0:
		roleV11, err := cr.client.rolesV11.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("directus: cannot get role: %w", err)
		}
		roleV11.ServerInfo = inf
		return roleV11, nil

	default:
		return nil, fmt.Errorf("directus: unsupported server version: %v", inf)
	}
}

func (cr *clientRoles) Create(ctx context.Context, role Role) (Role, error) {
	cr.mu.Lock()
  defer cr.mu.Unlock()

	inf, err := cr.client.Server.Info(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	// Support for Directus 10.0.X - 10.13.X
	case inf.VersionMajor == 10 && (inf.VersionMinor >= 0 && inf.VersionMinor <= 13):
		roleV10, ok := role.(*roleV10)
		if !ok {
			return nil, fmt.Errorf("directus: unsupported role type: %T", role)
		}
		return cr.client.rolesV10.Create(ctx, roleV10)

	// Support for Directus 11.0.X
	case inf.VersionMajor == 11 && inf.VersionMinor == 0:
		roleV11, ok := role.(*roleV11)
		if !ok {
			return nil, fmt.Errorf("directus: unsupported role type: %T", role)
		}
		return cr.client.rolesV11.Create(ctx, roleV11)

	default:
		return nil, fmt.Errorf("directus: unsupported server version: %v", inf)
	}
}

func (cr *clientRoles) Delete(ctx context.Context, role Role) error {
	cr.mu.Lock()
  defer cr.mu.Unlock()

	inf, err := cr.client.Server.Info(ctx)
	if err != nil {
		return err
	}

	switch {
	// Support for Directus 10.0.X - 10.13.X
	case inf.VersionMajor == 10 && (inf.VersionMinor >= 0 && inf.VersionMinor <= 13):
		roleV10, ok := role.(*roleV10)
		if !ok {
			return fmt.Errorf("directus: unsupported role type: %T", role)
		}
		return cr.client.rolesV10.Delete(ctx, roleV10.ID)

	// Support for Directus 11.0.X
	case inf.VersionMajor == 11 && inf.VersionMinor == 0:
		roleV11, ok := role.(*roleV11)
		if !ok {
			return fmt.Errorf("directus: unsupported role type: %T", role)
		}
		return cr.client.rolesV11.Delete(ctx, roleV11.ID)

	default:
		return fmt.Errorf("directus: unsupported server version: %v", inf)
	}
}

func (cr *clientRoles) Patch(ctx context.Context, role Role) (Role, error) {
	cr.mu.Lock()
  defer cr.mu.Unlock()

	inf, err := cr.client.Server.Info(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	// Support for Directus 10.0.X - 10.13.X
	case inf.VersionMajor == 10 && (inf.VersionMinor >= 0 && inf.VersionMinor <= 13):
		roleV10, ok := role.(*roleV10)
		if !ok {
			return nil, fmt.Errorf("directus: unsupported role type: %T", role)
		}
		return cr.client.rolesV10.Patch(ctx, roleV10.ID, roleV10)

	// Support for Directus 11.0.X
	case inf.VersionMajor == 11 && inf.VersionMinor == 0:
		roleV11, ok := role.(*roleV11)
		if !ok {
			return nil, fmt.Errorf("directus: unsupported role type: %T", role)
		}
		return cr.client.rolesV11.Patch(ctx, roleV11.ID, roleV11)

	default:
		return nil, fmt.Errorf("directus: unsupported server version: %v", inf)
	}
}
