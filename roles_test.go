package directus

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRolesList(t *testing.T) {
	i, err := initClient(t).Server.Info(context.Background())
	require.NoError(t, err)

	major, err := strconv.Atoi(strings.Split(i.Version, "")[0])
	require.NoError(t, err)

	if major >= 11 {
		roles, err := initClient(t).RolesV11.List(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, roles)

		for _, role := range roles {
			fmt.Printf("%#v\n", role)
		}
	}

	if major < 11 {
		roles, err := initClient(t).Roles.List(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, roles)

		for _, role := range roles {
			fmt.Printf("%#v\n", role)
		}
	}
}
