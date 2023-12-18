package directus

import (
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type ProtoJSON[T proto.Message] struct {
	Value T
}

func (n ProtoJSON[T]) String() string {
	return prototext.Format(n.Value)
}

func (n ProtoJSON[T]) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(n.Value)
}

func (n *ProtoJSON[T]) UnmarshalJSON(data []byte) error {
	n.Value = reflect.New(reflect.TypeOf(n.Value).Elem()).Interface().(T)
	return protojson.Unmarshal(data, n.Value)
}
