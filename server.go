package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/perimeterx/marshmallow"
)

type Server struct {
	Version string `json:"version"`

	Unknown map[string]any `json:"-"`
}

type clientServer struct {
	client *Client
}

func (server *Server) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, server, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	server.Unknown = values
	return nil
}

func (server *Server) MarshalJSON() ([]byte, error) {
	type alias Server
	base, err := json.Marshal((*alias)(server))
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	for k, v := range server.Unknown {
		m[k] = v
	}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (cr *clientServer) GetInfo(ctx context.Context) (*Server, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/settings"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data *Server `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}
