package directus

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Regime struct {
	ID     string `json:"id,omitempty"`
	Code   string `json:"code,omitempty"`
	Status string `json:"status"`
}

type RegimeWithTranslations struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`

	Translations []Relation[RegimeTranslation] `json:"translations,omitempty"`
}

type RegimeTranslation struct {
	Lang        string `json:"languages_code"`
	DisplayName string `json:"display_name"`
}

func initClient(t *testing.T) *Client {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})

	if os.Getenv("DIRECTUS_TOKEN") == "" {
		t.Skip("DIRECTUS_TOKEN not set")
	}
	return NewClient("https://compostela.admin.onetbooking.com", os.Getenv("DIRECTUS_TOKEN"), WithLogger(slog.New(handler)), WithBodyLogger())
}

func TestItemsList(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes")
	regimes, err := items.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, regimes)

	for _, regime := range regimes {
		fmt.Printf("%#v\n", regime)
	}
}

func TestListRoles(t *testing.T) {
	roles, err := initClient(t).ListRoles(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, roles)

	for _, role := range roles {
		fmt.Printf("%#v\n", role)
	}
}

func TestItemsListFields(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes", WithFields("id", "code"))
	regimes, err := items.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, regimes)

	for _, regime := range regimes {
		require.Empty(t, regime.Status)
		fmt.Printf("%#v\n", regime)
	}
}

func TestItemsGet(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes")
	regime, err := items.Get(context.Background(), "1158c0ee-bd5d-4e7b-a640-8bceca82b5db")
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%#v\n", regime)
}

func TestItemsGetFields(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes", WithFields("id", "code"))
	regime, err := items.Get(context.Background(), "1158c0ee-bd5d-4e7b-a640-8bceca82b5db")
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	require.Empty(t, regime.Status)
	fmt.Printf("%#v\n", regime)
}

func TestItemsGetNotFound(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes")
	regime, err := items.Get(context.Background(), "foo")
	require.Nil(t, regime)
	require.EqualError(t, err, "directus: item not found: foo")
}

func TestItemsCreate(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes")
	regime, err := items.Create(context.Background(), &Regime{
		Code:   "test",
		Status: "published",
	})
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%#v\n", regime)
}

func TestItemsCreateRelationship(t *testing.T) {
	items := NewItemsClient[RegimeWithTranslations](initClient(t), "regimes")
	regime, err := items.Create(context.Background(), &RegimeWithTranslations{
		Code: "test",
		Translations: []Relation[RegimeTranslation]{
			NewRelation(&RegimeTranslation{Lang: "en-GB", DisplayName: "Test"}),
			NewRelation(&RegimeTranslation{Lang: "es-ES", DisplayName: "Test"}),
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%#v\n", regime)
}

func TestItemsUpdate(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes")
	regime, err := items.Update(context.Background(), "1158c0ee-bd5d-4e7b-a640-8bceca82b5db", &Regime{
		Status: "draft",
	})
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%#v\n", regime)
}

func TestItemsUpdateFields(t *testing.T) {
	items := NewItemsClient[RegimeWithTranslations](initClient(t), "regimes", WithFields("*", "translations.*"))
	req := &RegimeWithTranslations{
		Status: "draft",
	}
	regime, err := items.Update(context.Background(), "1918e55e-8ad0-42bf-bc38-26304f0a6b1b", req)
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%+v\n", regime)
	fmt.Printf("%+v\n", regime.Translations[0].Value())
}

func TestItemsRelationNumeric(t *testing.T) {
	items := NewItemsClient[RegimeWithTranslations](initClient(t), "regimes", WithFields("*"))
	regime, err := items.Get(context.Background(), "1158c0ee-bd5d-4e7b-a640-8bceca82b5db")
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%+v\n", regime)
}

func TestItemsRelationValue(t *testing.T) {
	items := NewItemsClient[RegimeWithTranslations](initClient(t), "regimes", WithFields("*", "translations.*"))
	regime, err := items.Get(context.Background(), "1158c0ee-bd5d-4e7b-a640-8bceca82b5db")
	require.NoError(t, err)
	require.NotEmpty(t, regime)

	fmt.Printf("%+v\n", regime)
	fmt.Printf("%+v\n", regime.Translations[0].Value())
}
