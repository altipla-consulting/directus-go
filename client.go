package directus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Client struct {
	instance, token string
	logger          *slog.Logger
	bodyLogger      bool
}

type ClientOption func(client *Client)

func WithLogger(logger *slog.Logger) ClientOption {
	return func(client *Client) {
		client.logger = logger
	}
}

func WithBodyLogger() ClientOption {
	return func(client *Client) {
		client.bodyLogger = true
	}
}

func NewClient(instance string, token string, opts ...ClientOption) *Client {
	client := &Client{
		instance: strings.TrimRight(instance, "/"),
		token:    token,
		logger:   slog.New(slog.Default().Handler()),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (client *Client) urlf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s", client.instance, fmt.Sprintf(format, a...))
}

func (client *Client) do(req *http.Request, dest interface{}) error {
	client.logger.Debug("directus request", "method", req.Method, "url", req.URL.String())
	if client.bodyLogger && req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("directus: cannot read request body: %v", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		client.logger.Debug(string(body))
	}

	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("directus: request failed: %w", err)
	}
	defer resp.Body.Close()

	client.logger.Debug("directus reply", "status", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("directus: cannot read response body: %v", err)
	}
	if client.bodyLogger {
		client.logger.Debug(string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("directus: unexpected status code %v for url %q", resp.StatusCode, req.URL.String())
	}

	if dest != nil {
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(dest); err != nil {
			return fmt.Errorf("directus: cannot decode response: %v", err)
		}
	}

	return nil
}

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
	var buf *bytes.Buffer
	if request != nil {
		buf = bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(request); err != nil {
			return fmt.Errorf("directus: cannot encode request: %v", err)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	for _, opt := range items.opts {
		opt(req)
	}
	return items.c.do(req, &reply)
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

func (items *ItemsClient[T]) Get(ctx context.Context, id string) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodGet, items.c.urlf("/items/%s/%s", items.collection, id), nil, &reply); err != nil {
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
