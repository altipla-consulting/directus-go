package directus

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestItemsList(t *testing.T) {
	items := NewItemsClient[Regime](initClient(t), "regimes")
	regimes, err := items.List(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, regimes)

	for _, regime := range regimes {
		fmt.Printf("%#v\n", regime)
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

type PMSRoom struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	DisplayName string `json:"display_name"`
}

func TestItemsFilter(t *testing.T) {
	items := NewItemsClient[PMSRoom](initClient(t), "pms_rooms")
	rooms, err := items.Filter(context.Background(), Gt("max_people", 3))
	require.NoError(t, err)

	for _, room := range rooms {
		fmt.Printf("%#v\n", room)
	}
}

type Foo struct {
	ID           string `json:"id"`
	Translations []FooTranslation
}

type FooTranslation struct {
	LanguagesCode string `json:"languages_code"`
	DisplayName   string `json:"display_name"`
}

func TestItemsDeepFilter(t *testing.T) {
	var u string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u = r.URL.String()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"data": [{"id": "1", "translations": [{"languages_code": "en-GB", "display_name": "Test"}]}]}`)
	}))
	defer s.Close()
	client := NewClient(s.URL, "local-token", WithBodyLogger())
	items := NewItemsClient[Foo](client, "foo")

	results, err := items.List(context.Background(), WithDeepFilter("translations", Eq("languages_code", "en-GB")))
	require.NoError(t, err)

	require.Len(t, results, 1)
	require.Equal(t, "1", results[0].ID)
	require.Len(t, results[0].Translations, 1)
	require.Equal(t, "en-GB", results[0].Translations[0].LanguagesCode)
	require.Equal(t, "Test", results[0].Translations[0].DisplayName)

	got, err := url.QueryUnescape(u)
	require.NoError(t, err)
	require.Equal(t, `/items/foo?deep={"translations":{"_filter":{"languages_code":{"_eq":"en-GB"}}}}&limit=-1`, got)
}

func TestItemsMultipleDeepFilter(t *testing.T) {
	var u string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u = r.URL.String()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"data": []}`)
	}))
	defer s.Close()
	client := NewClient(s.URL, "local-token", WithBodyLogger())
	items := NewItemsClient[Foo](client, "foo")

	results, err := items.List(context.Background(),
		WithDeepFilter("translations", Eq("languages_code", "en-GB")),
		WithDeepFilter("another_field", Eq("display_name", "foo")))
	require.NoError(t, err)
	require.Empty(t, results)

	got, err := url.QueryUnescape(u)
	require.NoError(t, err)
	require.Equal(t, `/items/foo?deep={"another_field":{"_filter":{"display_name":{"_eq":"foo"}}},"translations":{"_filter":{"languages_code":{"_eq":"en-GB"}}}}&limit=-1`, got)
}
