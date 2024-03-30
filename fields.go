package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/perimeterx/marshmallow"
)

type Field struct {
	Collection string       `json:"collection"`
	Field      string       `json:"field"`
	Type       FieldType    `json:"type"`
	Meta       FieldMeta    `json:"meta"`
	Schema     *FieldSchema `json:"schema,omitempty"`
}

type FieldType string

const (
	FieldTypeString    FieldType = "string"
	FieldTypeAlias     FieldType = "alias"
	FieldTypeDate      FieldType = "date"
	FieldTypeTimestamp FieldType = "timestamp"
)

type FieldMeta struct {
	ID     int64      `json:"id"`
	Hidden bool       `json:"hidden"`
	Width  FieldWidth `json:"width"`

	ReadOnly bool `json:"read_only"`
	Required bool `json:"required"`

	Sort    int64          `json:"sort,omitempty"`
	System  bool           `json:"system,omitempty"`
	Special []FieldSpecial `json:"special,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (meta *FieldMeta) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, meta, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	meta.Unknown = values
	return nil
}

type FieldWidth string

const (
	FieldWidthFull FieldWidth = "full"
	FieldWidthHalf FieldWidth = "half"
)

type FieldSpecial string

const (
	FieldSpecialManyToOne   FieldSpecial = "m2o"
	FieldSpecialDateCreated FieldSpecial = "date-created"
	FieldSpecialDateUpdated FieldSpecial = "date-updated"
	FieldSpecialUUID        FieldSpecial = "uuid"
	FieldSpecialUserCreated FieldSpecial = "user-created"
	FieldSpecialUserUpdated FieldSpecial = "user-updated"
	FieldSpecialFile        FieldSpecial = "file"
	FieldSpecialAlias       FieldSpecial = "alias"
	FieldSpecialNoData      FieldSpecial = "no-data"
	FieldSpecialCastBoolean FieldSpecial = "cast-boolean"
)

type FieldSchema struct {
	Name      string `json:"name,omitempty"`
	Table     string `json:"table,omitempty"`
	DataType  string `json:"data_type"`
	MaxLength int64  `json:"max_length,omitempty"`

	IsNullable       bool `json:"is_nullable,omitempty"`
	IsUnique         bool `json:"is_unique,omitempty"`
	IsPrimaryKey     bool `json:"is_primary_key,omitempty"`
	HasAutoIncrement bool `json:"has_auto_increment,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (schema *FieldSchema) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, schema, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	schema.Unknown = values
	return nil
}

type clientFields struct {
	client *Client
}

func (cr *clientFields) List(ctx context.Context) ([]*Field, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/fields"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data []*Field `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr *clientFields) ListCollection(ctx context.Context, collection string) ([]*Field, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/fields/%s", collection), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data []*Field `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr *clientFields) Get(ctx context.Context, collection, field string) (*Field, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/fields/%s/%s", collection, field), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply Field
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (cr *clientFields) Create(ctx context.Context, field *Field) (*Field, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(field); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cr.client.urlf("/fields"), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply Field
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (cr *clientFields) Delete(ctx context.Context, collection, field string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, cr.client.urlf("/fields/%s/%s", collection, field), nil)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	return cr.client.sendRequest(req, nil)
}

func (cr *clientFields) Patch(ctx context.Context, collection, field string, f *Field) (*Field, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(f); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, cr.client.urlf("/fields/%s/%s", collection, field), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply Field
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
