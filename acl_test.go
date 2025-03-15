package directus

import (
	"context"
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
