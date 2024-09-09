package directus

import (
	"encoding/json"

	"github.com/perimeterx/marshmallow"
)

type Role struct {
	ID          string `json:"id,omitempty"`
	Icon        Icon   `json:"icon,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	AdminAccess bool `json:"admin_access"`
	AppAccess   bool `json:"app_access"`

	Users []string `json:"users,omitempty"`
}

type User struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
	Role      string   `json:"role,omitempty"`
	Policies  []string `json:"policies,omitempty"`

	Provider           string `json:"provider,omitempty"`
	ExternalIdentifier string `json:"external_identifier,omitempty"`
}

type PermissionAction string

const (
	PermissionActionCreate PermissionAction = "create"
	PermissionActionRead   PermissionAction = "read"
	PermissionActionUpdate PermissionAction = "update"
	PermissionActionDelete PermissionAction = "delete"
)

func (action *PermissionAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*action = PermissionAction(str)
	return nil
}

func (action *PermissionAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*action))
}

type Permission struct {
	ID         int64              `json:"id,omitempty"`
	Policy     Nullable[string]   `json:"policy"`
	Collection string             `json:"collection"`
	Action     PermissionAction   `json:"action"`
	Fields     Nullable[[]string] `json:"fields"`
	System     bool               `json:"system,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (permission *Permission) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, permission, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	permission.Unknown = values
	return nil
}

func (permission *Permission) MarshalJSON() ([]byte, error) {
	type alias Permission
	base, err := json.Marshal((*alias)(permission))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range permission.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type Policy struct {
	ID          string `json:"id,omitempty"`
	Icon        Icon   `json:"icon,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	AdminAccess bool `json:"admin_access"`
	AppAccess   bool `json:"app_access"`

	Users       []string `json:"users,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []int64  `json:"permissions,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (policy *Policy) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, policy, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	policy.Unknown = values
	return nil
}

func (policy *Policy) MarshalJSON() ([]byte, error) {
	type alias Policy
	base, err := json.Marshal((*alias)(policy))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range policy.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
