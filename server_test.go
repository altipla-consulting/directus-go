package directus

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerInfoUnmarshall(t *testing.T) {
	data := []byte(`
		{
			"project": {
			"project_name": "Directus",
			"project_descriptor": null,
			"project_logo": null,
			"project_color": "#6644FF",
			"default_appearance": "auto",
			"default_theme_light": null,
			"default_theme_dark": null,
			"theme_light_overrides": null,
			"theme_dark_overrides": null,
			"default_language": "es-ES",
			"public_foreground": null,
			"public_favicon": null,
			"public_note": null,
			"custom_css": null,
			"public_registration": false,
			"public_registration_verify_email": true,
			"public_background": null
			},
			"rateLimit": false,
			"rateLimitGlobal": false,
			"extensions": {
			"limit": null
			},
			"queryLimit": {
			"default": 100,
			"max": -1
			},
			"websocket": false,
			"version": "11.0.2"
		}`,
	)

	var server Info
	require.NoError(t, json.Unmarshal(data, &server))
	require.EqualValues(t, server.Version, "11.0.2")
}

func TestServerInfo(t *testing.T) {
	cli := initClient(t)
	s, err := cli.Server.Info(context.Background())
	require.NoError(t, err)
	require.NotNil(t, s)
}
