package directus

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrItemNotFound = errors.New("directus: item not found")
)

type unexpectedStatusError struct {
	status int
	url    *url.URL
}

func (e *unexpectedStatusError) Error() string {
	return fmt.Sprintf("directus: unexpected status code %v for url %q", e.status, e.url.String())
}
