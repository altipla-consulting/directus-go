package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResourceClient[T any] struct {
	client   *Client
	endpoint string
}

func NewResourceClient[T any](client *Client, endpoint string) *ResourceClient[T] {
	return &ResourceClient[T]{client, endpoint}
}

func (rc *ResourceClient[T]) List(ctx context.Context) ([]*T, error) {
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

func (rc *ResourceClient[T]) Get(ctx context.Context, id string) (*T, error) {
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

func (rc *ResourceClient[T]) Create(ctx context.Context, item *T) (*T, error) {
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

func (rc *ResourceClient[T]) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, rc.client.urlf("/%s/%s", rc.endpoint, id), nil)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	return rc.client.sendRequest(req, nil)
}

func (rc *ResourceClient[T]) Patch(ctx context.Context, id string, item *T) (*T, error) {
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
