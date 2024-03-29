package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Field struct {
	Collection string    `json:"collection"`
	Field      string    `json:"field"`
	Type       FieldType `json:"type"`
	Meta       FieldMeta `json:"meta"`
}

type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeAlias  FieldType = "alias"
)

type FieldMeta struct {
	ID     int32      `json:"id"`
	Hidden bool       `json:"hidden"`
	Width  FieldWidth `json:"width"`
}

type FieldWidth string

const (
	FieldWidthFull FieldWidth = "full"
	FieldWidthHalf FieldWidth = "half"
)

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
