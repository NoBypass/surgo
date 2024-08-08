package surgo

import (
	"reflect"
	"testing"
)

func TestIsStruct(t *testing.T) {
	type TestStruct struct{}
	var testStruct TestStruct
	var testStructPtr *TestStruct

	tests := []struct {
		input any
		want  bool
	}{
		{testStruct, true},
		{testStructPtr, true},
		{123, false},
		{"string", false},
	}

	for _, tt := range tests {
		if got := isStruct(tt.input); got != tt.want {
			t.Errorf("isStruct(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestStructToMap(t *testing.T) {
	type NestedStruct struct {
		Field1 string
		Field2 int
	}
	type TestStruct struct {
		Field1 string `db:"field1"`
		Field2 int    `db:"field2,omitempty"`
		Nested NestedStruct
	}

	tests := []struct {
		input any
		want  map[string]any
	}{
		{
			TestStruct{Field1: "value1", Field2: 0, Nested: NestedStruct{Field1: "nestedValue1", Field2: 2}},
			map[string]any{"field1": "value1", "Nested": map[string]any{"Field1": "nestedValue1", "Field2": 2}},
		},
		{
			TestStruct{Field1: "value1", Field2: 2, Nested: NestedStruct{Field1: "nestedValue1", Field2: 2}},
			map[string]any{"field1": "value1", "field2": 2, "Nested": map[string]any{"Field1": "nestedValue1", "Field2": 2}},
		},
	}

	for _, tt := range tests {
		if got := structToMap(tt.input); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("structToMap(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseVars(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}
	var testStruct = TestStruct{Field1: "value1", Field2: 2}
	var testStructPtr = &TestStruct{Field1: "value1", Field2: 2}

	tests := []struct {
		input map[string]any
		want  map[string]any
	}{
		{
			map[string]any{"key1": testStruct},
			map[string]any{"key1": map[string]any{"Field1": "value1", "Field2": 2}},
		},
		{
			map[string]any{"key1": testStructPtr},
			map[string]any{"key1": map[string]any{"Field1": "value1", "Field2": 2}},
		},
		{
			map[string]any{"key1": 123, "key2": "string"},
			map[string]any{"key1": 123, "key2": "string"},
		},
	}

	for _, tt := range tests {
		if got := parseVars(tt.input); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseVars(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
