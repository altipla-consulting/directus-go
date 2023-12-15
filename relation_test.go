package directus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type RelationCollection struct {
	ID  int64  `json:"id"`
	Foo string `json:"foo"`
}

func TestRelationUnmarshalListOfNumericIDs(t *testing.T) {
	data := []byte(`[1,2,3]`)
	var rel []Relation[RelationCollection]
	require.NoError(t, json.Unmarshal(data, &rel))

	require.Len(t, rel, 3)
	require.EqualValues(t, rel[0].NumericID(), int64(1))
	require.EqualValues(t, rel[1].NumericID(), int64(2))
	require.EqualValues(t, rel[2].NumericID(), int64(3))
}

func TestRelationUnmarshalListOfStringIDs(t *testing.T) {
	data := []byte(`["foo-1","foo-2","foo-3"]`)
	var rel []Relation[RelationCollection]
	require.NoError(t, json.Unmarshal(data, &rel))

	require.Len(t, rel, 3)
	require.EqualValues(t, rel[0].StringID(), "foo-1")
	require.EqualValues(t, rel[1].StringID(), "foo-2")
	require.EqualValues(t, rel[2].StringID(), "foo-3")
}

func TestRelationUnmarshalListOfEntities(t *testing.T) {
	data := []byte(`[{"id":1,"foo":"bar"},{"id":2,"foo":"baz"}]`)
	var rel []Relation[RelationCollection]
	require.NoError(t, json.Unmarshal(data, &rel))

	require.Len(t, rel, 2)
	require.EqualValues(t, rel[0].Value().ID, 1)
	require.EqualValues(t, rel[0].Value().Foo, "bar")
	require.EqualValues(t, rel[1].Value().ID, 2)
	require.EqualValues(t, rel[1].Value().Foo, "baz")
}

func TestRelationUnmarshalMixed(t *testing.T) {
	data := []byte(`[{"id":1,"foo":"bar"},2,"foo-3"]`)
	var rel []Relation[RelationCollection]
	require.NoError(t, json.Unmarshal(data, &rel))

	require.Len(t, rel, 3)
	require.EqualValues(t, rel[0].Value().ID, 1)
	require.EqualValues(t, rel[0].Value().Foo, "bar")
	require.EqualValues(t, rel[1].NumericID(), 2)
	require.EqualValues(t, rel[2].StringID(), "foo-3")
}
