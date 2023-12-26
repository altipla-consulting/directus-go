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
	Roles *clientRoles

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

	client.Roles = &clientRoles{c: client}

	return client
}

func (client *Client) urlf(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s", client.instance, fmt.Sprintf(format, a...))
}

func (client *Client) buildSendRequest(ctx context.Context, method, url string, request, reply any) error {
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
	return client.sendRequest(req, &reply)
}

func (client *Client) sendRequest(req *http.Request, dest interface{}) error {
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

	if resp.StatusCode != http.StatusOK && (req.Method != http.MethodDelete && resp.StatusCode != http.StatusNoContent) {
		return &unexpectedStatusError{
			status: resp.StatusCode,
			url:    req.URL,
		}
	}

	if dest != nil && len(body) > 0 {
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(dest); err != nil {
			return fmt.Errorf("directus: cannot decode response: %v", err)
		}
	}

	return nil
}
