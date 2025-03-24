package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRolesList(t *testing.T) {
	roles, err := initClient(t).Roles.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, roles)

	for _, role := range roles {
		fmt.Printf("%#v\n", role)
	}
}

func TestRolesCreate(t *testing.T) {
	create := &Role{
		Name: "Test",
		Policies: []RolePolicy{
			{ID: "365820c9-c20d-48cf-929e-58638b547a34"},
		},
	}
	role, err := initClient(t).Roles.Create(context.Background(), create)
	require.NoError(t, err)
	require.NotNil(t, role)

	fmt.Printf("%#v\n", role)
}

func TestPoliciesList(t *testing.T) {
	policies, err := initClient(t).Policies.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, policies)

	for _, policy := range policies {
		fmt.Printf("%#v\n", policy)
	}
}

func TestRoleMarshal(t *testing.T) {
	role := &Role{
		ID:    "1234-role",
		Name:  "Gestor",
		Icon:  "supervised_user_circle",
		Users: []string{"1234-user", "3456-user"},
		Policies: []RolePolicy{
			{ID: "1234-policy"},
			{ID: "3456-policy"},
		},
		existingPolicies: map[string]string{
			"1234-policy": "1234-policy-access",
			"3456-policy": "3456-policy-access",
		},
	}
	b, err := json.Marshal(role)
	require.NoError(t, err)
	require.JSONEq(t, `
		{
			"id": "1234-role",
			"name": "Gestor",
			"icon": "supervised_user_circle",
			"admin_access": false,
			"app_access": false,
			"users": ["1234-user", "3456-user"],
			"policies": [{ "id": "1234-policy" }, { "id": "3456-policy" }]
		}
	`, string(b))
}

func TestRoleUnmarshal(t *testing.T) {
	b := `
		{
			"id": "1234-role",
			"name": "Gestor",
			"icon": "supervised_user_circle",
			"admin_access": false,
			"app_access": false,
			"users": ["1234-user", "3456-user"],
			"policies": [{ "id": "1234-policy" }, { "id": "3456-policy" }]
		}
	`

	var role Role
	require.NoError(t, json.Unmarshal([]byte(b), &role))

	require.Equal(t, "1234-role", role.ID)
	require.Equal(t, "Gestor", role.Name)
	require.EqualValues(t, "supervised_user_circle", role.Icon)
	require.False(t, role.AdminAccess)
	require.False(t, role.AppAccess)
	require.Len(t, role.Users, 2)
	require.Equal(t, "1234-user", role.Users[0])
	require.Equal(t, "3456-user", role.Users[1])
	require.Len(t, role.Policies, 2)
	require.Equal(t, "1234-policy", role.Policies[0].ID)
	require.Equal(t, "3456-policy", role.Policies[1].ID)
	require.Empty(t, role.existingPolicies)
}
