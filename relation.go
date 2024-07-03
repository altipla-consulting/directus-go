package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/perimeterx/marshmallow"
)

type Relation[T any] struct {
	idstr string
	idnum int64
	value *T
}

func NewRelation[T any](data *T) Relation[T] {
	return Relation[T]{value: data}
}

func NewRelationID[T any](id string) Relation[T] {
	return Relation[T]{idstr: id}
}

func NewRelationNumericID[T any](id int64) Relation[T] {
	return Relation[T]{idnum: id, idstr: fmt.Sprintf("%d", id)}
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

func (r Relation[T]) Empty() bool {
	return r.value == nil && r.idstr == "" && r.idnum == 0
}

func (r Relation[T]) MarshalJSON() ([]byte, error) {
	if r.idstr != "" {
		return json.Marshal(r.idstr)
	}
	if r.idnum != 0 {
		return json.Marshal(r.idnum)
	}
	return json.Marshal(r.value)
}

func (r *Relation[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if err := json.Unmarshal(data, &r.idstr); err == nil {
		return nil
	}
	if err := json.Unmarshal(data, &r.idnum); err == nil {
		return nil
	}
	return json.Unmarshal(data, &r.value)
}

type RelationDefinition struct {
	Collection        string `json:"collection"`
	Field             string `json:"field"`
	RelatedCollection string `json:"related_collection"`

	Schema RelationSchema `json:"schema"`
	Meta   RelationMeta   `json:"meta"`
}

type RelationSchema struct {
	Table    string         `json:"table,omitempty"`
	Column   string         `json:"column,omitempty"`
	OnUpdate RelationAction `json:"on_update,omitempty"`
	OnDelete RelationAction `json:"on_delete,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (schema *RelationSchema) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, schema, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	schema.Unknown = values
	return nil
}

func (schema *RelationSchema) MarshalJSON() ([]byte, error) {
	type alias RelationSchema
	base, err := json.Marshal((*alias)(schema))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range schema.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type RelationAction string

const (
	RelationActionCascade  RelationAction = "CASCADE"
	RelationActionSetNull  RelationAction = "SET NULL"
	RelationActionNoAction RelationAction = "NO ACTION"
)

func (action *RelationAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*action = RelationAction(str)
	return nil
}

type RelationMeta struct {
	ID int64 `json:"id,omitempty"`

	System bool `json:"system,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (meta *RelationMeta) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, meta, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	meta.Unknown = values
	return nil
}

func (meta *RelationMeta) MarshalJSON() ([]byte, error) {
	type alias RelationMeta
	base, err := json.Marshal((*alias)(meta))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range meta.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type clientRelations struct {
	client *Client
}

func (cr *clientRelations) List(ctx context.Context) ([]*RelationDefinition, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/relations"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data []*RelationDefinition `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr *clientRelations) ListCollection(ctx context.Context, collection string) ([]*RelationDefinition, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/relations/%s", collection), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data []*RelationDefinition `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr *clientRelations) Get(ctx context.Context, collection, field string) (*RelationDefinition, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/relations/%s/%s", collection, field), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply RelationDefinition
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (cr *clientRelations) Create(ctx context.Context, field *RelationDefinition) (*RelationDefinition, error) {
	if field.Collection == "" {
		return nil, fmt.Errorf("directus: field collection is required")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(field); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cr.client.urlf("/relations"), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply RelationDefinition
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (cr *clientRelations) Delete(ctx context.Context, collection, field string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, cr.client.urlf("/relations/%s/%s", collection, field), nil)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	return cr.client.sendRequest(req, nil)
}

func (cr *clientRelations) Patch(ctx context.Context, f *RelationDefinition) (*RelationDefinition, error) {
	if f.Collection == "" {
		return nil, fmt.Errorf("directus: field collection is required")
	}
	if f.Field == "" {
		return nil, fmt.Errorf("directus: field name is required")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(f); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, cr.client.urlf("/relations/%s/%s", f.Collection, f.Field), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply RelationDefinition
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
