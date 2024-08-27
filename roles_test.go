package directus

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRolesList(t *testing.T) {

	i, err := initClient(t).Server.Info(context.Background())
	require.NoError(t, err)

	major, err := strconv.Atoi(strings.Split(i.Version, ".")[0])
	require.NoError(t, err)

	if major < 11 {
		t.Skip("Roles are only available in Directus 11 and above")
	}

	roles, err := initClient(t).Roles.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, roles)

	for _, role := range roles {
		r, err := json.Marshal(role)
		require.NoError(t, err)
		slog.Info("Role", slog.Any("role", r))
	}

}
