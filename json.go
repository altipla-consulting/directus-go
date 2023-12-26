package directus

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JSON[T any] struct {
	Value T
}

func (n JSON[T]) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n JSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

func (n *JSON[T]) UnmarshalJSON(data []byte) error {
	n.Value = reflect.New(reflect.TypeOf(n.Value)).Elem().Interface().(T)
	return json.Unmarshal(data, &n.Value)
}
