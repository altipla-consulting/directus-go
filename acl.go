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

	Users    []string     `json:"users,omitempty"`
	Policies []RolePolicy `json:"policies,omitempty"`

	existingPolicies map[string]string
}

type RolePolicy struct {
	ID string `json:"id"`

	accessID string
}

func (role *Role) UnmarshalJSON(data []byte) error {
	type alias Role
	if err := json.Unmarshal(data, (*alias)(role)); err != nil {
		return err
	}
	role.existingPolicies = make(map[string]string)
	for _, rp := range role.Policies {
		role.existingPolicies[rp.ID] = rp.accessID
	}
	return nil
}

func (role *Role) MarshalJSON() ([]byte, error) {
	type alias Role
	base, err := json.Marshal((*alias)(role))
	if err != nil {
		return nil, err
	}

	m := make(map[string]any)
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}

	if role.existingPolicies == nil {
		role.existingPolicies = make(map[string]string)
	}

	alt := Alterations[rolePolicyInternal, string]{}
	present := make(map[string]bool)
	for _, rp := range role.Policies {
		present[rp.ID] = true

		if _, ok := role.existingPolicies[rp.ID]; !ok {
			alt.Create = append(alt.Create, &rolePolicyInternal{Policy: rp.ID})
		}
	}
	for policy, accessID := range role.existingPolicies {
		if !present[policy] {
			alt.Delete = append(alt.Delete, accessID)
		}
	}
	m["policies"] = alt

	return json.Marshal(m)
}

type rolePolicyInternal struct {
	ID     string `json:"id,omitempty"`
	Policy string `json:"policy"`
}

func (rp *RolePolicy) UnmarshalJSON(data []byte) error {
	var read rolePolicyInternal
	if err := json.Unmarshal(data, &read); err != nil {
		return err
	}
	rp.ID = read.Policy
	rp.accessID = read.ID
	return nil
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
