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
)

// ItemsClient access the items API in a type-safe way.
type ItemsClient[T any] struct {
	c          *Client
	collection string
	opts       []ReadOption
}

type readOptionApply struct {
	req  *http.Request
	deep map[string]deepFilter
}

type deepFilter struct {
	Filter Filter `json:"_filter"`
}

// ReadOption configures the returned data from Directus when reading or returning items.
type ReadOption func(apply *readOptionApply)

// WithFields filters the fields of each returned item. It can add relations deep fields to the response to obtain them
// in the same request.
func WithFields(fields ...string) ReadOption {
	return func(apply *readOptionApply) {
		q := apply.req.URL.Query()
		for _, field := range fields {
			q.Add("fields[]", field)
		}
		apply.req.URL.RawQuery = q.Encode()
	}
}

// WithSort sorts the returned items by the given fields. Use a minus sign (-) to sort in descending order.
// It does not order deep relations inside each of the items. To sort deep relations, use WithDeepSort.
func WithSort(sort ...string) ReadOption {
	return func(apply *readOptionApply) {
		q := apply.req.URL.Query()
		for _, s := range sort {
			q.Add("sort[]", s)
		}
		apply.req.URL.RawQuery = q.Encode()
	}
}

// WithLimit limits the number of returned items.
func WithLimit(limit int64) ReadOption {
	return func(apply *readOptionApply) {
		q := apply.req.URL.Query()
		q.Set("limit", fmt.Sprintf("%d", limit))
		apply.req.URL.RawQuery = q.Encode()
	}
}

// WithOffset skips the first n items.
func WithOffset(offset int64) ReadOption {
	return func(apply *readOptionApply) {
		q := apply.req.URL.Query()
		q.Add("offset", fmt.Sprintf("%d", offset))
		apply.req.URL.RawQuery = q.Encode()
	}
}

// WithDeepSort sorts the deep relations of each returned item by the given fields. Use a minus sign (-) to sort in
// descending order. It does not order the items themselves. To sort the items, use WithSort.
func WithDeepSort(field string, sort ...string) ReadOption {
	return func(apply *readOptionApply) {
		q := apply.req.URL.Query()
		for _, s := range sort {
			q.Add(fmt.Sprintf("deep[%s][_sort][]", field), s)
		}
		apply.req.URL.RawQuery = q.Encode()
	}
}

// WithDeepLimit limits the number of returned deep relations of each item.
func WithDeepLimit(field string, limit int64) ReadOption {
	return func(apply *readOptionApply) {
		q := apply.req.URL.Query()
		q.Add(fmt.Sprintf("deep[%s][_limit]", field), fmt.Sprintf("%d", limit))
		apply.req.URL.RawQuery = q.Encode()
	}
}

// WithDeepFilter adds a filter to the deep relations returned from the query. It is applied together with WithFields()
// most of the time to extract the related items.
func WithDeepFilter(field string, filter Filter) ReadOption {
	return func(apply *readOptionApply) {
		apply.deep[field] = deepFilter{Filter: filter}
	}
}

// NewItemsClient creates a new client to access & write items in a type-safe way.
func NewItemsClient[T any](client *Client, collection string, opts ...ReadOption) *ItemsClient[T] {
	return &ItemsClient[T]{
		c:          client,
		collection: collection,
		opts:       opts,
	}
}

func (items *ItemsClient[T]) applyOpts(req *http.Request, opts ...ReadOption) error {
	apply := &readOptionApply{
		req:  req,
		deep: make(map[string]deepFilter),
	}

	for _, opt := range items.opts {
		opt(apply)
	}
	for _, opt := range opts {
		opt(apply)
	}

	q := apply.req.URL.Query()
	if len(apply.deep) > 0 {
		b, err := json.Marshal(apply.deep)
		if err != nil {
			return fmt.Errorf("directus: cannot marshal deep filter: %v", err)
		}
		q.Set("deep", string(b))
	}
	req.URL.RawQuery = q.Encode()

	return nil
}

func (items *ItemsClient[T]) itemsdo(ctx context.Context, method, url string, request, reply any, opts ...ReadOption) error {
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
	if err := items.applyOpts(req, opts...); err != nil {
		return err
	}
	return items.c.sendRequest(req, &reply)
}

// List items of the collection.
func (items *ItemsClient[T]) List(ctx context.Context, opts ...ReadOption) ([]*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, items.c.urlf("/items/%s", items.collection), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}

	// Return all items of the collection by default.
	q := req.URL.Query()
	q.Set("limit", "-1")
	req.URL.RawQuery = q.Encode()

	if err := items.applyOpts(req, opts...); err != nil {
		return nil, err
	}

	reply := struct {
		Data []*T `json:"data"`
	}{}
	if err := items.c.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

// Filter items of the collection.
func (items *ItemsClient[T]) Filter(ctx context.Context, filter Filter, opts ...ReadOption) ([]*T, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}

	if err := items.applyOpts(req, opts...); err != nil {
		return nil, err
	}

	reply := struct {
		Data []*T `json:"data"`
	}{}
	if err := items.c.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

// Get a single item by its primary key. If it cannot be found, it returns ErrItemNotFound.
func (items *ItemsClient[T]) Get(ctx context.Context, id string, opts ...ReadOption) (*T, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: %v", ErrItemNotFound, id)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, items.c.urlf("/items/%s/%s", items.collection, id), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}

	if err := items.applyOpts(req, opts...); err != nil {
		return nil, err
	}

	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.c.sendRequest(req, &reply); err != nil {
		var e *unexpectedStatusError
		if errors.As(err, &e) && e.status == http.StatusForbidden {
			return nil, fmt.Errorf("%w: %v", ErrItemNotFound, id)
		}
		return nil, err
	}
	return reply.Data, nil
}

// Create a new item in the collection.
func (items *ItemsClient[T]) Create(ctx context.Context, item *T) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodPost, items.c.urlf("/items/%s", items.collection), item, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}

// Update an item in the collection by its primary key.
func (items *ItemsClient[T]) Update(ctx context.Context, id string, item *T) (*T, error) {
	reply := struct {
		Data *T `json:"data"`
	}{}
	if err := items.itemsdo(ctx, http.MethodPatch, items.c.urlf("/items/%s/%s", items.collection, id), item, &reply); err != nil {
		if errors.Is(err, ErrItemNotFound) {
			return nil, fmt.Errorf("%w: %v", err, id)
		}
		return nil, err
	}
	return reply.Data, nil
}

// Delete an item from the collection by its primary key.
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
