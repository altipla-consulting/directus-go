package directus

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JSON[T any] struct {
	Value T
}

func NewJSON[T any](value T) JSON[T] {
	return JSON[T]{Value: value}
}

func (n JSON[T]) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (n JSON[T]) MarshalJSON() ([]byte, error) {
	content, err := json.Marshal(n.Value)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot marshal json: %w", err)
	}
	return json.Marshal(string(content))
}

func (n *JSON[T]) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("directus: cannot unmarshal json: %w", err)
	}
	n.Value = reflect.New(reflect.TypeOf(n.Value)).Elem().Interface().(T)
	if value == "" {
		return nil
	}
	return json.Unmarshal([]byte(value), &n.Value)
}
