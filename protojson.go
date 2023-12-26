package directus

import (
	"encoding/json"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type ProtoJSON[T proto.Message] struct {
	Value T
}

func NewProtoJSON[T proto.Message](value T) ProtoJSON[T] {
	return ProtoJSON[T]{Value: value}
}

func (n ProtoJSON[T]) String() string {
	return prototext.Format(n.Value)
}

func (n ProtoJSON[T]) MarshalJSON() ([]byte, error) {
	content, err := protojson.Marshal(n.Value)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot marshal protobuf json: %w", err)
	}
	return json.Marshal(string(content))
}

func (n *ProtoJSON[T]) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("directus: cannot unmarshal protobuf json: %w", err)
	}
	n.Value = reflect.New(reflect.TypeOf(n.Value).Elem()).Interface().(T)
	if value == "" {
		return nil
	}
	return protojson.Unmarshal([]byte(value), n.Value)
}
