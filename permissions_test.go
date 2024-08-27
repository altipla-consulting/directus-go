package directus

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissionst(t *testing.T) {

	permissions, err := initClient(t).Permissions.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, permissions)

	for _, permission := range permissions {
		p, err := json.Marshal(permission)
		require.NoError(t, err)
		slog.Info("Permission", slog.Any("permission", p))
	}

}
