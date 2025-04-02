package directus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFiltersMultipleAnd(t *testing.T) {
	json, err := FilterJSON(And(Eq("domain", "foo-domain"), Eq("auth_token", "foo-auth-token"), Eq("status", "PUBLISHED")))
	require.NoError(t, err)
	require.JSONEq(t, string(json), `{
		"_and": [
			{ "domain": { "_eq": "foo-domain" } },
			{ "auth_token": { "_eq": "foo-auth-token" } },
			{ "status": { "_eq": "PUBLISHED" } }
		]
	}`)
}
