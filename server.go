package directus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/perimeterx/marshmallow"
)

type clientServer struct {
	client *Client
}

type Info struct {
	Version string `json:"version"`

	Unknown map[string]any `json:"-"`
}

func (server *Info) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, server, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	server.Unknown = values
	return nil
}

func (cr *clientServer) Info(ctx context.Context) (*Info, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/server/info"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data *Info `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}
