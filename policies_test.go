package directus

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoliciesList(t *testing.T) {

	i, err := initClient(t).Server.Info(context.Background())
	require.NoError(t, err)

	major, err := strconv.Atoi(strings.Split(i.Version, ".")[0])
	require.NoError(t, err)

	if major < 11 {
		t.Skip("Policies are only available in Directus 11 and above")
	}

	policies, err := initClient(t).Policies.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, policies)

	for _, policy := range policies {
		fmt.Printf("%#v\n", policy)
	}

}
