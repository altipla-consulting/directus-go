package directus

import (
	"encoding/json"

	"github.com/perimeterx/marshmallow"
)

type Role struct {
	ID          string `json:"id,omitempty"`
	Icon        Icon   `json:"icon,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	AdminAccess bool `json:"admin_access"`
	AppAccess   bool `json:"app_access"`

	Users []string `json:"users,omitempty"`
}

type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`

	Provider           string `json:"provider,omitempty"`
	ExternalIdentifier string `json:"external_identifier,omitempty"`
}

type Icon string

type Collection struct {
	Collection string            `json:"collection"`
	Meta       CollectionMeta    `json:"meta"`
	Schema     *CollectionSchema `json:"schema,omitempty"`
	Fields     []*Field          `json:"fields,omitempty"`
}

type CollectionMeta struct {
	Collection string `json:"collection,omitempty"`

	Color     string `json:"color,omitempty"`
	Icon      Icon   `json:"icon,omitempty"`
	Note      string `json:"note,omitempty"`
	Hidden    bool   `json:"hidden"`
	Singleton bool   `json:"singleton"`

	ArchiveField     string `json:"archive_field,omitempty"`
	ArchiveValue     string `json:"archive_value,omitempty"`
	UnarchiveValue   string `json:"unarchive_value,omitempty"`
	ArchiveAppFilter bool   `json:"archive_app_filter"`

	SortField string `json:"sort_field,omitempty"`

	Group    Nullable[string]   `json:"group,omitempty"`
	Sort     int64              `json:"sort,omitempty"`
	Collapse CollectionCollapse `json:"collapse,omitempty"`

	Versioning bool `json:"versioning,omitempty"`

	Accountability Nullable[Accountability] `json:"accountability"`

	System bool `json:"system,omitempty"`

	PreviewURL string `json:"preview_url,omitempty"`

	DisplayTemplate string `json:"display_template,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (meta *CollectionMeta) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, meta, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	meta.Unknown = values
	return nil
}

func (meta *CollectionMeta) MarshalJSON() ([]byte, error) {
	type alias CollectionMeta
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

type CollectionCollapse string

const (
	CollectionCollapseOpen   CollectionCollapse = "open"
	CollectionCollapseClosed CollectionCollapse = "closed"
	CollectionCollapseLocked CollectionCollapse = "locked"
)

type CollectionSchema struct {
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`

	Collation string `json:"collation,omitempty"`
	Engine    string `json:"engine,omitempty"`
	Schema    string `json:"schema,omitempty"`

	Unknown map[string]any `json:"-"`
}

func (schema *CollectionSchema) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, schema, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	schema.Unknown = values
	return nil
}

func (schema *CollectionSchema) MarshalJSON() ([]byte, error) {
	type alias CollectionSchema
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

type Accountability string

const (
	AccountabilityAll      Accountability = "all"
	AccountabilityActivity Accountability = "activity"
)

type RelationDefinition struct {
	Collection        string `json:"collection"`
	Field             string `json:"field"`
	RelatedCollection string `json:"related_collection"`

	Schema RelationSchema `json:"schema"`
	Meta   RelationMeta   `json:"meta"`
}

type RelationSchema struct {
	Table    string         `json:"table"`
	Column   string         `json:"column"`
	OnUpdate RelationAction `json:"on_update"`
	OnDelete RelationAction `json:"on_delete"`

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

type CustomTranslation struct {
	ID       string `json:"id,omitempty"`
	Key      string `json:"key,omitempty"`
	Language string `json:"language,omitempty"`
	Value    string `json:"value"`
}

type Folder struct {
	ID     string           `json:"id,omitempty"`
	Name   string           `json:"name"`
	Parent Nullable[string] `json:"parent"`
}

type Preset struct {
	ID         int64            `json:"id,omitempty"`
	Bookmark   Nullable[string] `json:"bookmark"`
	User       Nullable[string] `json:"user"`
	Role       Nullable[string] `json:"role"`
	Collection string           `json:"collection"`
	Layout     Nullable[string] `json:"layout"`
	Icon       Nullable[string] `json:"icon"`

	Unknown map[string]any `json:"-"`
}

func (preset *Preset) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, preset, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	preset.Unknown = values
	return nil
}

func (preset *Preset) MarshalJSON() ([]byte, error) {
	type alias Preset
	base, err := json.Marshal((*alias)(preset))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range preset.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type Flow struct {
	ID             string           `json:"id,omitempty"`
	Name           string           `json:"name"`
	Status         string           `json:"status"`
	Description    Nullable[string] `json:"description"`
	Accountability Accountability   `json:"accountability"`
	Color          string           `json:"color,omitempty"`
	Icon           Icon             `json:"icon,omitempty"`
	Operation      Nullable[string] `json:"operation"`

	Unknown map[string]any `json:"-"`
}

func (flow *Flow) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, flow, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	flow.Unknown = values
	return nil
}

func (flow *Flow) MarshalJSON() ([]byte, error) {
	type alias Flow
	base, err := json.Marshal((*alias)(flow))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range flow.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
