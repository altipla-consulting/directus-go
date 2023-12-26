package directus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type testJSON struct {
	Content JSON[*testJSONContent] `json:"content"`
}

type testJSONContent struct {
	Foo string `json:"foo"`
}

func TestJSONMarshal(t *testing.T) {
	value := testJSON{
		Content: NewJSON(&testJSONContent{Foo: "bar"}),
	}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, string(data), `{"content": "{\"foo\":\"bar\"}"}`)
}

func TestJSONMarshalNil(t *testing.T) {
	value := testJSON{}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, string(data), `{"content": "null"}`)
}

func TestJSONUnmarshal(t *testing.T) {
	value := testJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"content": "{\"foo\":\"bar\"}"}`), &value))
	require.NotNil(t, value.Content.Value)
	require.Equal(t, value.Content.Value.Foo, "bar")
}

func TestJSONUnmarshalNil(t *testing.T) {
	value := testJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"content": "null"}`), &value))
	require.Nil(t, value.Content.Value)
}

func TestJSONUnmarshalEmpty(t *testing.T) {
	value := testJSON{}
	require.NoError(t, json.Unmarshal([]byte(`{"content": ""}`), &value))
	require.Nil(t, value.Content.Value)
}
