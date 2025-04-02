package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Filter interface {
	content() any
	String() string
}

type filterOperator struct {
	field string
	op    string
	value any
}

func (f filterOperator) content() any {
	return map[string]any{
		f.field: map[string]any{
			f.op: f.value,
		},
	}
}

func (f filterOperator) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.content())
}

func (f filterOperator) String() string {
	return fmt.Sprintf("%s %s %v", f.field, f.op, f.value)
}

func Eq(field string, value any) Filter {
	return filterOperator{field: field, op: "_eq", value: value}
}

func Neq(field string, value any) Filter {
	return filterOperator{field: field, op: "_neq", value: value}
}

func Gt(field string, value any) Filter {
	return filterOperator{field: field, op: "_gt", value: value}
}

func Gte(field string, value any) Filter {
	return filterOperator{field: field, op: "_gte", value: value}
}

func Lt(field string, value any) Filter {
	return filterOperator{field: field, op: "_lt", value: value}
}

func Lte(field string, value any) Filter {
	return filterOperator{field: field, op: "_lte", value: value}
}

func Empty(field string) Filter {
	return filterOperator{field: field, op: "_empty", value: nil}
}

func NotEmpty(field string) Filter {
	return filterOperator{field: field, op: "_nempty", value: nil}
}

func In(field string, values ...any) Filter {
	return filterOperator{field: field, op: "_in", value: values}
}

func Between(field string, from, to any) Filter {
	return filterOperator{field: field, op: "_between", value: []any{from, to}}
}

func StartsWith(field string, value string) Filter {
	return filterOperator{field: field, op: "_starts_with", value: value}
}

type filterLogical struct {
	op     string
	values []Filter
}

func (f filterLogical) content() any {
	var values []any
	for _, v := range f.values {
		values = append(values, v.content())
	}
	return map[string]any{
		f.op: values,
	}
}

func (f filterLogical) String() string {
	vals := make([]string, len(f.values))
	for i, v := range f.values {
		vals[i] = v.String()
	}
	if f.op == "_and" {
		return strings.Join(vals, " && ")
	}
	if f.op == "_or" {
		return strings.Join(vals, " || ")
	}

	// Shouldn't reach here really.
	return strings.Join(vals, " "+f.op+" ")
}

func (f filterLogical) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.content())
}

func And(filters ...Filter) Filter {
	return filterLogical{op: "_and", values: filters}
}

func Or(filters ...Filter) Filter {
	return filterLogical{op: "_or", values: filters}
}

type filterRelated struct {
	field  string
	filter Filter
}

func (f filterRelated) content() any {
	return map[string]any{
		f.field: f.filter.content(),
	}
}

func (f filterRelated) String() string {
	return fmt.Sprintf("%s.%s", f.field, f.filter.String())
}

func (f filterRelated) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.content())
}

func Related(field string, filter Filter) Filter {
	return filterRelated{field, filter}
}

type filterEmpty struct{}

func (f filterEmpty) String() string {
	return ""
}

func (f filterEmpty) content() any {
	return map[string]any{}
}

func (f filterEmpty) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.content())
}

func Noop() Filter {
	return filterEmpty{}
}

func FilterJSON(filter Filter) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(filter.content()); err != nil {
		return "", fmt.Errorf("directus: cannot encode filter: %v", err)
	}
	return buf.String(), nil
}
