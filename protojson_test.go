package directus

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoJSONMarshal(t *testing.T) {
	value := ProtoJSON[*timestamppb.Timestamp]{
		Value: timestamppb.New(time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)),
	}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.Equal(t, string(data), `"2020-01-02T03:04:05.000000006Z"`)
}

func TestProtoJSONMarshalNil(t *testing.T) {
	var value ProtoJSON[*structpb.Struct]
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.Equal(t, string(data), `{}`)
}

func TestProtoJSONUnmarshal(t *testing.T) {
	var value ProtoJSON[*timestamppb.Timestamp]
	require.NoError(t, json.Unmarshal([]byte(`"2020-01-02T03:04:05.000000006Z"`), &value))
	require.Equal(t, value.Value.AsTime().UTC(), time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC))
}

func TestProtoJSONUnmarshalNil(t *testing.T) {
	var value ProtoJSON[*structpb.Struct]
	require.NoError(t, json.Unmarshal([]byte(`{}`), &value))
	require.Empty(t, value.Value.String())
}

func TestProtoJSONNewInferred(t *testing.T) {
	value := NewProtoJSON(timestamppb.New(time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)))
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.Equal(t, string(data), `"2020-01-02T03:04:05.000000006Z"`)
}
