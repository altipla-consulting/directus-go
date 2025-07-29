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
	FieldTypeUUID      FieldType = "uuid"
)

type FieldMeta struct {
	ID     int64      `json:"id,omitempty"`
	Hidden bool       `json:"hidden"`
	Width  FieldWidth `json:"width,omitempty"`

	ReadOnly bool `json:"read_only"`
	Required bool `json:"required"`

	Sort    int64          `json:"sort,omitempty"`
	System  bool           `json:"system,omitempty"`
	Special []FieldSpecial `json:"special,omitempty"`

	Translations []*FieldTranslation `json:"translations,omitempty"`

	Options *FieldOptions `json:"options,omitempty"`

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

func (meta *FieldMeta) MarshalJSON() ([]byte, error) {
	type alias FieldMeta
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

func (meta *FieldMeta) HasSpecial(special FieldSpecial) bool {
	for _, s := range meta.Special {
		if s == special {
			return true
		}
	}
	return false
}

func (meta *FieldMeta) Translation(language string) *FieldTranslation {
	for _, t := range meta.Translations {
		if t.Language == language {
			return t
		}
	}
	return nil
}

type FieldWidth string

const (
	FieldWidthFull FieldWidth = "full"
	FieldWidthHalf FieldWidth = "half"
)

func (width *FieldWidth) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*width = FieldWidth(str)
	return nil
}

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
	FieldSpecialGroup       FieldSpecial = "group"
)

func (special *FieldSpecial) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*special = FieldSpecial(value)
	return nil
}

type FieldTranslation struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

type FieldOptions struct {
	Choices FieldChoices `json:"choices,omitempty"`

	unknown map[string]any
}

func (options *FieldOptions) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, options, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	options.unknown = values
	return nil
}

func (options *FieldOptions) MarshalJSON() ([]byte, error) {
	type alias FieldOptions
	base, err := json.Marshal((*alias)(options))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range options.unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type FieldChoices struct {
	Choices []*FieldChoice
	Values  []any
}

func (choices *FieldChoices) MarshalJSON() ([]byte, error) {
	switch {
	case len(choices.Choices) > 0:
		return json.Marshal(choices.Choices)
	case len(choices.Values) > 0:
		return json.Marshal(choices.Values)
	default:
		return []byte("null"), nil
	}
}

func (choices *FieldChoices) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, rawchoice := range raw {
		var c string
		if err := json.Unmarshal(rawchoice, &c); err == nil {
			choices.Values = append(choices.Values, c)
			continue
		}

		var choice FieldChoice
		if err := json.Unmarshal(rawchoice, &choice); err != nil {
			return err
		}
		choices.Choices = append(choices.Choices, &choice)
	}
	return nil
}

type FieldChoice struct {
	Text  string `json:"text"`
	Value any    `json:"value"`
}

type FieldSchema struct {
	Name      string `json:"name,omitempty"`
	Table     string `json:"table,omitempty"`
	DataType  string `json:"data_type"`
	MaxLength int64  `json:"max_length,omitempty"`

	IsNullable       bool `json:"is_nullable"`
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

func (schema *FieldSchema) MarshalJSON() ([]byte, error) {
	type alias FieldSchema
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
	var reply = struct {
		Data *Field `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr *clientFields) Create(ctx context.Context, field *Field) (*Field, error) {
	if field.Collection == "" {
		return nil, fmt.Errorf("directus: field collection is required")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(field); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cr.client.urlf("/fields/%s", field.Collection), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply = struct {
		Data *Field `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (cr *clientFields) Delete(ctx context.Context, collection, field string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, cr.client.urlf("/fields/%s/%s", collection, field), nil)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	return cr.client.sendRequest(req, nil)
}

func (cr *clientFields) Patch(ctx context.Context, f *Field) (*Field, error) {
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, cr.client.urlf("/fields/%s/%s", f.Collection, f.Field), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply = struct {
		Data *Field `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}
