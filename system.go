package directus

import "github.com/perimeterx/marshmallow"

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
}

type CollectionMeta struct {
	Collection string `json:"collection,omitempty"`

	Icon      Icon   `json:"icon,omitempty"`
	Note      string `json:"note,omitempty"`
	Hidden    bool   `json:"hidden"`
	Singleton bool   `json:"singleton"`

	ArchiveField     string `json:"archive_field,omitempty"`
	ArchiveValue     string `json:"archive_value,omitempty"`
	UnarchiveValue   string `json:"unarchive_value,omitempty"`
	ArchiveAppFilter bool   `json:"archive_app_filter"`

	SortField string `json:"sort_field,omitempty"`

	Group    string             `json:"group,omitempty"`
	Sort     int64              `json:"sort,omitempty"`
	Collapse CollectionCollapse `json:"collapse,omitempty"`

	Versioning bool `json:"versioning,omitempty"`

	Accountability *Accountability `json:"accountability"`

	System bool `json:"system,omitempty"`

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

type RelationAction string

const (
	RelationActionCascade  RelationAction = "CASCADE"
	RelationActionSetNull  RelationAction = "SET NULL"
	RelationActionNoAction RelationAction = "NO ACTION"
)

type RelationMeta struct {
	ID int64 `json:"id"`

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

type CustomTranslation struct {
	ID          string `json:"id,omitempty"`
	Key         string `json:"key"`
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

type Folder struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name"`
	Parent string `json:"parent"`
}
