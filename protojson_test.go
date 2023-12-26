package directus

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/altipla-consulting/directus-go/internal/testproto"
)

type testProtoJSONTimestamp struct {
	Timestamp ProtoJSON[*timestamppb.Timestamp] `json:"ts"`
}

type testProtoJSONSimple struct {
	Simple ProtoJSON[*pb.SimpleMessage] `json:"simple"`
}

func TestProtoJSONMarshalTimestamp(t *testing.T) {
	value := &testProtoJSONTimestamp{
		Timestamp: NewProtoJSON(timestamppb.New(time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC))),
	}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, string(data), `{"ts": "\"2020-01-02T03:04:05.000000006Z\""}`)
}

func TestProtoJSONMarshalSimple(t *testing.T) {
	value := &testProtoJSONSimple{
		Simple: NewProtoJSON(&pb.SimpleMessage{TestValue: "hello"}),
	}
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, string(data), `{"simple": "{\"testValue\":\"hello\"}"}`)
}

func TestProtoJSONMarshalNil(t *testing.T) {
	value := new(testProtoJSONSimple)
	data, err := json.Marshal(value)
	require.NoError(t, err)
	require.JSONEq(t, string(data), `{"simple": "{}"}`)
}

func TestProtoJSONUnmarshal(t *testing.T) {
	value := new(testProtoJSONTimestamp)
	require.NoError(t, json.Unmarshal([]byte(`{"ts": "\"2020-01-02T03:04:05.000000006Z\""}`), &value))
	require.Equal(t, value.Timestamp.Value.AsTime().UTC(), time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC))
}

func TestProtoJSONUnmarshalSimple(t *testing.T) {
	value := new(testProtoJSONSimple)
	require.NoError(t, json.Unmarshal([]byte(`{"simple": "{\"testValue\":\"hello\"}"}`), &value))
	require.NotNil(t, value.Simple.Value)
	require.Equal(t, value.Simple.Value.TestValue, "hello")
}

func TestProtoJSONUnmarshalNil(t *testing.T) {
	value := new(testProtoJSONSimple)
	require.NoError(t, json.Unmarshal([]byte(`{"simple": "{}"}`), &value))
	require.Empty(t, value.Simple.Value.String())
}

func TestProtoJSONUnmarshalEmpty(t *testing.T) {
	value := new(testProtoJSONSimple)
	require.NoError(t, json.Unmarshal([]byte(`{"simple": ""}`), &value))
	require.Empty(t, value.Simple.Value.String())
}
