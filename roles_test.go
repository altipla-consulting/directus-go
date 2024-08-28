package directus

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRolesList(t *testing.T) {
	roles, err := initClient(t).Roles.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, roles)

	for _, role := range roles {
		r, err := json.Marshal(role)
		require.NoError(t, err)
		slog.Info("Role", slog.Any("role", r))
	}

}
