package directus

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/perimeterx/marshmallow"
)

type clientServer struct {
	client *Client
}

type ServerInfo struct {
	Version      string `json:"version"`
	VersionMajor int
	VersionMinor int
	VersionPatch int

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

func (cr *clientServer) Info(ctx context.Context) (*ServerInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cr.client.urlf("/server/info"), nil)
	if err != nil {
		return nil, fmt.Errorf("directus: cannot prepare request: %v", err)
	}
	reply := struct {
		Data *ServerInfo `json:"data"`
	}{}
	if err := cr.client.sendRequest(req, &reply); err != nil {
		return nil, err
	}

	chunks := strings.Split(reply.Data.Version, ".")
	if chunks == nil || len(chunks) < 3 {
		return nil, fmt.Errorf("directus: cannot parse version: %v", reply.Data.Version)
	}

	var major, minor, patch int
	major, err = strconv.Atoi(chunks[0])
	if err == nil {
		reply.Data.VersionMajor = major
	}
	minor, err = strconv.Atoi(chunks[1])
	if err == nil {
		reply.Data.VersionMinor = minor
	}
	patch, err = strconv.Atoi(chunks[2])
	if err == nil {
		reply.Data.VersionPatch = patch
	}

	return reply.Data, nil
}
