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
