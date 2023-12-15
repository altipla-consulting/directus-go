package directus

import (
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	Value T
	Valid bool
}

func (n Nullable[T]) String() string {
	if n.Valid {
		return fmt.Sprintf("%+v", n.Value)
	}
	return "NULL"
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Value)
	}
	return []byte("null"), nil
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Valid = false
	n.Value = *new(T)

	if string(data) == "null" {
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.Value)
}
