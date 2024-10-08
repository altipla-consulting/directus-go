package directus

import (
	"encoding/json"

	"github.com/perimeterx/marshmallow"
)

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

	Translations []*CollectionTranslation `json:"translations,omitempty"`

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

func (meta *CollectionMeta) Translation(language string) *CollectionTranslation {
	for _, t := range meta.Translations {
		if t.Language == language {
			return t
		}
	}
	return nil
}

type CollectionCollapse string

const (
	CollectionCollapseOpen   CollectionCollapse = "open"
	CollectionCollapseClosed CollectionCollapse = "closed"
	CollectionCollapseLocked CollectionCollapse = "locked"
)

func (collapse *CollectionCollapse) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*collapse = CollectionCollapse(str)
	return nil
}

func (collapse *CollectionCollapse) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*collapse))
}

type CollectionTranslation struct {
	Language    string `json:"language"`
	Singular    string `json:"singular"`
	Plural      string `json:"plural"`
	Translation string `json:"translation"`
}

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

func (accountability *Accountability) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*accountability = Accountability(str)
	return nil
}

func (accountability *Accountability) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*accountability))
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

type Operation struct {
	ID          string           `json:"id,omitempty"`
	Flow        string           `json:"flow"`
	Key         string           `json:"key"`
	PositionX   int32            `json:"position_x"`
	PositionY   int32            `json:"position_y"`
	Type        string           `json:"type"`
	Name        Nullable[string] `json:"name"`
	Reject      Nullable[string] `json:"reject"`
	Resolve     Nullable[string] `json:"resolve"`
	UserCreated Nullable[string] `json:"user_created"`

	Unknown map[string]any `json:"-"`
}

func (operation *Operation) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, operation, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	operation.Unknown = values
	return nil
}

func (operation *Operation) MarshalJSON() ([]byte, error) {
	type alias Operation
	base, err := json.Marshal((*alias)(operation))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range operation.Unknown {
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
	Operations     []string         `json:"operations,omitempty"`
	UserCreated    Nullable[string] `json:"user_created"`

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

type File struct {
	FileSize        Nullable[int64]  `json:"file_size"`
	ID              string           `json:"id,omitempty"`
	Folder          Nullable[string] `json:"folder"`
	Title           Nullable[string] `json:"title"`
	Type            Nullable[string] `json:"type"`
	Description     Nullable[string] `json:"description"`
	Storage         string           `json:"storage"`
	Charset         Nullable[string] `json:"charset"`
	FilenameDowload string           `json:"filename_download"`
	FocalPointX     Nullable[int32]  `json:"focal_point_x"`
	FocalPointY     Nullable[int32]  `json:"focal_point_y"`
	Width           Nullable[int32]  `json:"width"`
	Height          Nullable[int32]  `json:"height"`
	Duration        Nullable[int32]  `json:"duration"`
	Location        Nullable[string] `json:"location"`
	Tags            Nullable[string] `json:"tags"`
	Embed           Nullable[string] `json:"embed"`
	FilenameDisk    Nullable[string] `json:"filename_disk"`

	Unknown map[string]any `json:"-"`
}

func (file *File) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, file, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	file.Unknown = values
	return nil
}

func (file *File) MarshalJSON() ([]byte, error) {
	type alias File
	base, err := json.Marshal((*alias)(file))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range file.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)

}

type Dashboard struct {
	ID    string           `json:"id,omitempty"`
	Name  string           `json:"name"`
	Icon  Icon             `json:"icon"`
	Color Nullable[string] `json:"color"`
	Note  Nullable[string] `json:"note"`
}

type Panel struct {
	ID         string           `json:"id,omitempty"`
	Dashboard  string           `json:"dashboard"`
	Height     int32            `json:"height"`
	Width      int32            `json:"width"`
	PositionX  int32            `json:"position_x"`
	PositionY  int32            `json:"position_y"`
	ShowHeader bool             `json:"show_header"`
	Type       string           `json:"type"`
	Color      Nullable[string] `json:"color"`
	Icon       Icon             `json:"icon"`
	Name       Nullable[string] `json:"name"`
	Note       Nullable[string] `json:"note"`

	Unknown map[string]any `json:"-"`
}

func (panel *Panel) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, panel, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	panel.Unknown = values
	return nil
}

func (panel *Panel) MarshalJSON() ([]byte, error) {
	type alias Panel
	base, err := json.Marshal((*alias)(panel))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range panel.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
