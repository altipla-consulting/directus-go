package directus

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldMetaUnmarshal(t *testing.T) {
	data := []byte(`
		{
			"collection": "contact_data",
			"field": "phone",
			"type": "string",
			"schema": {
				"name": "phone",
				"table": "contact_data",
				"data_type": "varchar",
				"default_value": null,
				"generation_expression": null,
				"max_length": 255,
				"numeric_precision": null,
				"numeric_scale": null,
				"is_generated": false,
				"is_nullable": false,
				"is_unique": false,
				"is_primary_key": false,
				"has_auto_increment": false,
				"foreign_key_column": null,
				"foreign_key_table": null,
				"comment": ""
			},
			"meta": {
				"id": 1406,
				"collection": "contact_data",
				"field": "phone",
				"special": null,
				"interface": "input",
				"options": {
					"trim": true,
					"iconRight": "phone_in_talk"
				},
				"display": null,
				"display_options": null,
				"readonly": false,
				"hidden": false,
				"sort": 4,
				"width": "half",
				"translations": [
					{
						"language": "es-ES",
						"translation": "Teléfono"
					}
				],
				"note": null,
				"conditions": null,
				"required": true,
				"group": null,
				"validation": null,
				"validation_message": null
			}
		}
	`)
	var field Field
	require.NoError(t, json.Unmarshal(data, &field))
	require.EqualValues(t, field.Meta.ID, 1406)
	require.Equal(t, field.Meta.Width, FieldWidthHalf)
}

func TestFieldMetaOptionsChoicesString(t *testing.T) {
	data := []byte(`
		{
			"collection": "contact_data",
			"field": "phone",
			"type": "string",
			"schema": {
				"name": "phone",
				"table": "contact_data",
				"data_type": "varchar",
				"default_value": null,
				"generation_expression": null,
				"max_length": 255,
				"numeric_precision": null,
				"numeric_scale": null,
				"is_generated": false,
				"is_nullable": false,
				"is_unique": false,
				"is_primary_key": false,
				"has_auto_increment": false,
				"foreign_key_column": null,
				"foreign_key_table": null,
				"comment": ""
			},
			"meta": {
				"id": 1406,
				"collection": "contact_data",
				"field": "phone",
				"special": null,
				"interface": "input",
				"options": {
					"choices": ["GET", "POST"]
				},
				"display": null,
				"display_options": null,
				"readonly": false,
				"hidden": false,
				"sort": 4,
				"width": "half",
				"translations": [
					{
						"language": "es-ES",
						"translation": "Teléfono"
					}
				],
				"note": null,
				"conditions": null,
				"required": true,
				"group": null,
				"validation": null,
				"validation_message": null
			}
		}
	`)
	var field Field
	require.NoError(t, json.Unmarshal(data, &field))
	require.EqualValues(t, field.Meta.ID, 1406)
	require.Equal(t, field.Meta.Width, FieldWidthHalf)
	require.Len(t, field.Meta.Options.Choices.Values, 2)
	require.Equal(t, field.Meta.Options.Choices.Values[0], "GET")
	require.Equal(t, field.Meta.Options.Choices.Values[1], "POST")
}

func TestFieldMetaOptionsChoicesNull(t *testing.T) {
	data := []byte(`
		{
			"collection": "contact_data",
			"field": "phone",
			"type": "string",
			"schema": {
				"name": "phone",
				"table": "contact_data",
				"data_type": "varchar",
				"default_value": null,
				"generation_expression": null,
				"max_length": 255,
				"numeric_precision": null,
				"numeric_scale": null,
				"is_generated": false,
				"is_nullable": false,
				"is_unique": false,
				"is_primary_key": false,
				"has_auto_increment": false,
				"foreign_key_column": null,
				"foreign_key_table": null,
				"comment": ""
			},
			"meta": {
				"id": 1406,
				"collection": "contact_data",
				"field": "phone",
				"special": null,
				"interface": "input",
				"options": null,
				"display": null,
				"display_options": null,
				"readonly": false,
				"hidden": false,
				"sort": 4,
				"width": "half",
				"translations": [
					{
						"language": "es-ES",
						"translation": "Teléfono"
					}
				],
				"note": null,
				"conditions": null,
				"required": true,
				"group": null,
				"validation": null,
				"validation_message": null
			}
		}
	`)
	var field Field
	require.NoError(t, json.Unmarshal(data, &field))
	require.EqualValues(t, field.Meta.ID, 1406)
	require.Equal(t, field.Meta.Width, FieldWidthHalf)
}

func TestFieldMetaOptionsChoices(t *testing.T) {
	//Load data from /internal/mocks/fields.json file
	js, err := os.ReadFile("internal/mocks/fields.json")
	require.NoError(t, err)

	reply := struct {
		Data []*Field `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal(js, &reply))

	for _, field := range reply.Data {
		if field.Meta.Options == nil {
			fmt.Println("Field Meta Options: null")
		} else {
			fmt.Printf("Field Meta Options: %+v\n\n", field.Meta.Options)
		}

		// Verificar la serialización y deserialización
		data, err := json.Marshal(field)
		require.NoError(t, err)

		var deserializedField Field
		require.NoError(t, json.Unmarshal(data, &deserializedField))

		if deserializedField.Meta.Options == nil {
			fmt.Println("Deserialized Field Meta Options: null")
		} else {
			fmt.Printf("Deserialized Field Meta Options: %+v\n\n", deserializedField.Meta.Options)
		}
	}
}
