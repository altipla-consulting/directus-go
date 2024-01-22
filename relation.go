package directus

import (
	"encoding/json"
	"fmt"
)

type Relation[T any] struct {
	idstr string
	idnum int64
	value *T
	Valid bool
}

func NewRelation[T any](data *T) Relation[T] {
	return Relation[T]{value: data}
}

func (r Relation[T]) Value() *T {
	if r.value == nil {
		panic("directus: do not extract values from a relation without loaded fields")
	}
	return r.value
}

func (r Relation[T]) StringID() string {
	if r.idstr == "" {
		panic("directus: do not extract the string ID of a relation that doesn't have it")
	}
	return r.idstr
}

func (r Relation[T]) NumericID() int64 {
	if r.idnum == 0 {
		panic("directus: do not extract the numeric ID of a relation that doesn't have it")
	}
	return r.idnum
}

func (r Relation[T]) String() string {
	if r.value != nil {
		return fmt.Sprintf("%+v", r.value)
	}
	if r.idstr != "" {
		return r.idstr
	}
	if r.idnum != 0 {
		return fmt.Sprintf("%v", r.idnum)
	}
	return "INVALID_RELATION"
}

func (r Relation[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value)
}

func (r *Relation[T]) UnmarshalJSON(data []byte) error {
	if r.value != nil {
		r.Valid = true
	}
	if err := json.Unmarshal(data, &r.idstr); err == nil {
		return nil
	}
	if err := json.Unmarshal(data, &r.idnum); err == nil {
		return nil
	}
	return json.Unmarshal(data, &r.value)
}
