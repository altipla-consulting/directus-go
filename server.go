package directus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/perimeterx/marshmallow"
	"golang.org/x/mod/semver"
)

type clientServer struct {
	client *Client
}

type ServerInfo struct {
	Version      string `json:"version"`

	Unknown map[string]any `json:"-"`
}

func (server *ServerInfo) UnmarshalJSON(data []byte) error {
	values, err := marshmallow.Unmarshal(data, server, marshmallow.WithExcludeKnownFieldsFromMap(true))
	if err != nil {
		return err
	}
	server.Unknown = values
	return nil
}

func (cs *clientServer) Info(ctx context.Context) (*ServerInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cs.client.urlf("/server/info"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data *ServerInfo `json:"data"`
	}{}
	if err := cs.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}

	return reply.Data, nil
}

// Supports Directus 10.0.0 - 10.13.9
func (cs *clientServer) ValidV10(ctx context.Context) bool {
	inf, err := cs.Info(ctx)
	if err != nil {
		return false
	}
	return semver.Compare(inf.Version, "10.0.0") >= 0 && semver.Compare(inf.Version, "10.13.9") <= 0
}

// Supports Directus 11.0.0 and above
func (cs *clientServer) ValidV11(ctx context.Context) bool {
	inf, err := cs.Info(ctx)
	if err != nil {
		return false
	}

	return semver.Compare(inf.Version, "11.0.0") >= 0
}
