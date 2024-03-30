package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourcesCollectionsList(t *testing.T) {
	client := initClient(t)
	collections, err := client.Collections.List(context.Background())
	require.NoError(t, err)

	fmt.Printf("%#v\n", collections[0].Schema)

	e, err := json.Marshal(collections[0])
	require.NoError(t, err)
	fmt.Println(string(e))
}
