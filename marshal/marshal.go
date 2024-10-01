package marshal

import (
	"reflect"
	"strings"
)

type Marshaler string

func (m *Marshaler) Marshal(vars map[string]any) map[string]any {
	for k, v := range vars {
		vars[k] = m.marshal(v)
	}

	return vars
}

func (m *Marshaler) marshal(v any) any {
	if isTime(v) {
		return parseTimes(v)
	} else if isStruct(v) {
		return m.handleStruct(v)
	} else if isSlice(v) {
		return m.handleSlice(v)
	} else if isMap(v) {
		return m.handleMap(v)
	} else {
		return v
	}
}

func isMap(x any) bool {
	t := reflect.TypeOf(x)
	return t.Kind() == reflect.Map || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Map)
}

func isSlice(x any) bool {
	t := reflect.TypeOf(x)
	return t.Kind() == reflect.Slice || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice)
}

func isStruct(x any) bool {
	t := reflect.TypeOf(x)
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

func (m *Marshaler) handleMap(x any) map[string]any {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	if t.Kind() == reflect.Ptr {
		return m.Marshal(v.Elem().Interface().(map[string]any))
	}

	return m.Marshal(v.Interface().(map[string]any))
}

func (m *Marshaler) handleSlice(x any) []any {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	resolved := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		resolved[i] = m.marshal(v.Index(i).Interface())
	}

	return resolved
}

func (m *Marshaler) handleStruct(x any) map[string]any {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	resolved := make(map[string]any)
	for i := range t.NumField() {
		field := t.Field(i)
		name := field.Name
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = field.Tag.Get(string(*m))
		}

		if dbTag != "" {
			vals := strings.Split(dbTag, ",")
			if vals[0] == "-" {
				continue
			} else {
				name = vals[0]
			}

			if v.Field(i).IsZero() && len(vals) > 1 && vals[1] == "omitempty" {
				continue
			}
		}

		val := v.Field(i).Interface()
		resolved[name] = val
	}

	return m.Marshal(resolved)
}
