package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResourceClient[T any, PK string | int64] struct {
	client   *Client
	endpoint string
}

func NewResourceClient[T any, PK string | int64](client *Client, endpoint string) *ResourceClient[T, PK] {
	return &ResourceClient[T, PK]{client, endpoint}
}

func (rc *ResourceClient[T, PK]) List(ctx context.Context) ([]*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.client.urlf("/%s", rc.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data []*T `json:"data"`
	}{}
	if err := rc.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (rc *ResourceClient[T, PK]) Get(ctx context.Context, id PK) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.client.urlf("/%s/%s", rc.endpoint, id), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply T
	if err := rc.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (rc *ResourceClient[T, PK]) Create(ctx context.Context, item *T) (*T, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(item); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.client.urlf("/%s", rc.endpoint), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply T
	if err := rc.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (rc *ResourceClient[T, PK]) Delete(ctx context.Context, id PK) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, rc.client.urlf("/%s/%s", rc.endpoint, id), nil)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	return rc.client.sendRequest(req, nil)
}

func (rc *ResourceClient[T, PK]) Patch(ctx context.Context, id PK, item *T) (*T, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(item); err != nil {
		return nil, fmt.Errorf("directus: cannot encode request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, rc.client.urlf("/%s/%s", rc.endpoint, id), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	var reply T
	if err := rc.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
