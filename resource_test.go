package directus

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourcesCollectionsList(t *testing.T) {
	client := initClient(t)
	collections, err := client.Collections.List(context.Background())
	require.NoError(t, err)

	for _, c := range collections {
		fmt.Printf("%#v\n", c.Schema)
	}
}
