package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ItemsClient[T any] struct {
	c          *Client
	collection string
	opts       []ReadOption
}

func NewItemsClient[T any](client *Client, collection string, opts ...ReadOption) *ItemsClient[T] {
	return &ItemsClient[T]{
		c:          client,
		collection: collection,
		opts:       opts,
	}
}

type ReadOption func(req *http.Request)

func WithFields(fields ...string) ReadOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Set("fields", strings.Join(fields, ","))
		req.URL.RawQuery = q.Encode()
	}
}

func (items *ItemsClient[T]) itemsdo(ctx context.Context, method, url string, request, reply any) error {
	var body io.Reader
	if request != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(request); err != nil {
			return fmt.Errorf("directus: cannot encode request: %v", err)
		}
		body = &buf
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	for _, opt := range items.opts {
		opt(req)
	}
	return items.c.sendRequest(req, &reply)
}

func (items *ItemsClient[T]) List(ctx context.Context) ([]*T, error) {
	reply := struct {
		Data []*T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodGet, items.c.urlf("/items/%s", items.collection), nil, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (items *ItemsClient[T]) Filter(ctx context.Context, filter Filter) ([]*T, error) {
	u, err := url.Parse(items.c.urlf("/items/%s", items.collection))
	if err != nil {
		return nil, err
	}
	qs := u.Query()
	f, err := FilterJSON(filter)
	if err != nil {
		return nil, err
	}
	qs.Set("filter", f)
	u.RawQuery = qs.Encode()

	reply := struct {
		Data []*T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodGet, u.String(), nil, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (items *ItemsClient[T]) Get(ctx context.Context, id string) (*T, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: %v", ErrItemNotFound, id)
	}
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodGet, items.c.urlf("/items/%s/%s", items.collection, id), nil, &reply); err != nil {
		var e *unexpectedStatusError
		if errors.As(err, &e) && e.status == http.StatusForbidden {
			return nil, fmt.Errorf("%w: %v", ErrItemNotFound, id)
		}
		return nil, err
	}
	return reply.Data, nil
}

func (items *ItemsClient[T]) Create(ctx context.Context, item *T) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodPost, items.c.urlf("/items/%s", items.collection), item, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (items *ItemsClient[T]) Update(ctx context.Context, id string, item *T) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodPatch, items.c.urlf("/items/%s/%s", items.collection, id), item, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (items *ItemsClient[T]) Delete(ctx context.Context, id string) error {
	return items.itemsdo(ctx, http.MethodDelete, items.c.urlf("/items/%s/%s", items.collection, id), nil, nil)
}

type SingletonClient[T any] struct {
	items *ItemsClient[T]
}

func NewSingletonClient[T any](client *Client, collection string, opts ...ReadOption) *SingletonClient[T] {
	return &SingletonClient[T]{
		items: NewItemsClient[T](client, collection, opts...),
	}
}

func (s *SingletonClient[T]) Get(ctx context.Context) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := s.items.itemsdo(ctx, http.MethodGet, s.items.c.urlf("/items/%s", s.items.collection), nil, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

func (s *SingletonClient[T]) Update(ctx context.Context, item *T) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := s.items.itemsdo(ctx, http.MethodPatch, s.items.c.urlf("/items/%s", s.items.collection), item, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}
