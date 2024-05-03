package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/perimeterx/marshmallow"
)

type Settings struct {
	ProjectName       string           `json:"project_name"`
	ProjectURL        Nullable[string] `json:"project_url"`
	ProjectDescriptor Nullable[string] `json:"project_descriptor"`
	ProjectColor      string           `json:"project_color"`

	DefaultLanguage string      `json:"default_language"`
	ModuleBar       []ModuleBar `json:"module_bar"`

	AuthPasswordPolicy string `json:"auth_password_policy"`
	AuthLoginAttempts  int32  `json:"auth_login_attempts"`

	CustomCSS Nullable[string] `json:"custom_css"`

	Unknown map[string]any `json:"-"`
}

func (settings *Settings) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, settings, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	settings.Unknown = values
	return nil
}

func (settings *Settings) MarshalJSON() ([]byte, error) {
	type alias Settings
	base, err := json.Marshal((*alias)(settings))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range settings.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type ModuleBar struct {
	ID      string        `json:"id"`
	Locked  bool          `json:"locked"`
	Enabled bool          `json:"enabled"`
	Type    ModuleBarType `json:"type"`
	Name    string        `json:"name,omitempty"`
	Icon    string        `json:"icon,omitempty"`
	URL     string        `json:"url,omitempty"`
}

type ModuleBarType string

const (
	ModuleBarTypeLink   ModuleBarType = "link"
	ModuleBarTypeModule ModuleBarType = "module"
)

type clientSettings struct {
	client *Client
}

func (cr *clientSettings) Get(ctx context.Context) (*Settings, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/settings"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply Settings
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (cr *clientSettings) Update(ctx context.Context, settings *Settings) (*Settings, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(settings); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, cr.client.urlf("/settings"), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply Settings
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
