package directus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNullableStringUnmarshalNull(t *testing.T) {
	data := []byte(`null`)
	var n Nullable[string]
	require.NoError(t, json.Unmarshal(data, &n))
	require.False(t, n.Valid)
	require.Empty(t, n.Value)
}

func TestNullableStringUnmarshalFilled(t *testing.T) {
	data := []byte(`"foo"`)
	var n Nullable[string]
	require.NoError(t, json.Unmarshal(data, &n))
	require.True(t, n.Valid)
	require.Equal(t, n.Value, "foo")
}

func TestNullableStringUnmarshalEmpty(t *testing.T) {
	data := []byte(`""`)
	var n Nullable[string]
	require.NoError(t, json.Unmarshal(data, &n))
	require.True(t, n.Valid)
	require.Empty(t, n.Value)
}

func TestNullableStringMarshalNull(t *testing.T) {
	var n Nullable[string]
	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, data, []byte("null"))
}

func TestNullableStringMarshalFilled(t *testing.T) {
	n := Nullable[string]{
		Valid: true,
		Value: "foo",
	}
	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, data, []byte(`"foo"`))
}

func TestNullableStringMarshalEmpty(t *testing.T) {
	n := Nullable[string]{
		Valid: true,
	}
	data, err := json.Marshal(n)
	require.NoError(t, err)
	require.Equal(t, data, []byte(`""`))
}
