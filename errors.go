package directus

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	// ErrItemNotFound is returned when the item is not found in the collection.
	ErrItemNotFound = errors.New("directus: item not found")
)

type unexpectedStatusError struct {
	status int
	url    *url.URL
}

func (e *unexpectedStatusError) Error() string {
	return fmt.Sprintf("directus: unexpected status code %v for url %q", e.status, e.url.String())
}

type Error struct {
	Message    string          `json:"message"`
	Extensions ErrorExtensions `json:"extensions"`
}

func (e Error) Error() string {
	if e.Extensions.Code != "" {
		return fmt.Sprintf("directus: %s (code: %s)", e.Message, e.Extensions.Code)
	}
	return fmt.Sprintf("directus: %s", e.Message)
}

type ErrorExtensions struct {
	Code ErrorCode `json:"code"`
}

type ErrorCode string

const (
	ErrorCodeRecordNotUnique = "RECORD_NOT_UNIQUE"
)
