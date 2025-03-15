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
	fields   []string
}

type ResourceClientOption[T any, PK string | int64] func(rc *ResourceClient[T, PK])

func WithResourceFields[T any, PK string | int64](fields ...string) ResourceClientOption[T, PK] {
	return func(rc *ResourceClient[T, PK]) {
		rc.fields = fields
	}
}

func NewResourceClient[T any, PK string | int64](client *Client, endpoint string, opts ...ResourceClientOption[T, PK]) *ResourceClient[T, PK] {
	rc := &ResourceClient[T, PK]{
		client:   client,
		endpoint: endpoint,
	}
	for _, opt := range opts {
		opt(rc)
	}
	return rc
}

func (rc *ResourceClient[T, PK]) List(ctx context.Context) ([]*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.client.urlf("/%s", rc.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	if len(rc.fields) > 0 {
		q := req.URL.Query()
		for _, field := range rc.fields {
			q.Add("fields[]", field)
		}
		req.URL.RawQuery = q.Encode()
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.client.urlf("/%s/%v", rc.endpoint, id), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	if len(rc.fields) > 0 {
		q := req.URL.Query()
		for _, field := range rc.fields {
			q.Add("fields[]", field)
		}
		req.URL.RawQuery = q.Encode()
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
	if len(rc.fields) > 0 {
		q := req.URL.Query()
		for _, field := range rc.fields {
			q.Add("fields[]", field)
		}
		req.URL.RawQuery = q.Encode()
	}
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := rc.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (rc *ResourceClient[T, PK]) Delete(ctx context.Context, id PK) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, rc.client.urlf("/%s/%v", rc.endpoint, id), nil)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, rc.client.urlf("/%s/%v", rc.endpoint, id), &buf)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	if len(rc.fields) > 0 {
		q := req.URL.Query()
		for _, field := range rc.fields {
			q.Add("fields[]", field)
		}
		req.URL.RawQuery = q.Encode()
	}
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := rc.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}
