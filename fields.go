package directus

import (
	"context"
	"net/http"
)

type clientFields struct {
	c *Client
}

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

func (cr clientFields) List(ctx context.Context) ([]Field, error) {
	reply := struct {
		Data []Field `json:"data"`
	}{}
	if err := cr.c.buildSendRequest(ctx, http.MethodGet, cr.c.urlf("/fields"), nil, &reply); err != nil {
		return nil, err
	}

	return reply.Data, nil
}

func (cr clientFields) ListCollection(ctx context.Context, collection string) ([]Field, error) {
	reply := struct {
		Data []Field `json:"data"`
	}{}
	if err := cr.c.buildSendRequest(ctx, http.MethodGet, cr.c.urlf("/fields/%s", collection), nil, &reply); err != nil {
		return nil, err
	}

	return reply.Data, nil
}
