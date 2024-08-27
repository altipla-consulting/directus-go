package directus

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessesList(t *testing.T) {
	i, err := initClient(t).Server.Info(context.Background())
	require.NoError(t, err)

	major, err := strconv.Atoi(strings.Split(i.Version, "")[0])
	require.NoError(t, err)

	if major >= 11 {
		accesses, err := initClient(t).Accesses.List(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, accesses)

		for _, access := range accesses {
			fmt.Printf("%#v\n", access)
		}
	}
}
