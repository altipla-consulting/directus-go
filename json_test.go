package directus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type jsonTest struct {
	Foo string `json:"foo"`
}

func TestJSONMarshal(t *testing.T) {
	value := JSON[jsonTest]{
		Value: jsonTest{Foo: "bar"},
	}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.Equal(t, string(data), `{"foo":"bar"}`)
}

func TestJSONMarshalNil(t *testing.T) {
	var value JSON[*jsonTest]
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.Equal(t, string(data), `null`)
}

func TestJSONUnmarshal(t *testing.T) {
	var value JSON[jsonTest]
	require.NoError(t, json.Unmarshal([]byte(`{"foo":"bar"}`), &value))
	require.Equal(t, value.Value.Foo, "bar")
}

func TestJSONUnmarshalNil(t *testing.T) {
	var value JSON[*jsonTest]
	require.NoError(t, json.Unmarshal([]byte(`null`), &value))
	require.Nil(t, value.Value)
}
