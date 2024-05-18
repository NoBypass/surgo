package surgo

import (
	"fmt"
	"reflect"
	"strings"
)

func parseParams(args []any) (map[string]any, error) {
	currIdx := 0
	if len(args) == 0 {
		return nil, nil
	} else if len(args) == 1 {
		return parseParam(args[0], &currIdx)
	}

	m := make(map[string]any)
	for _, arg := range args {
		if arg == nil {
			continue
		}
		nm, err := parseParam(arg, &currIdx)
		if err != nil {
			return nil, err
		}
		for k, v := range nm {
			switch v.(type) {
			case []any:
				if _, ok := m[k]; ok {
					m[k] = append(m[k].([]any), v.([]any)...)
					continue
				}
				m[k] = v
			default:
				m[k] = v
			}
		}
	}

	return m, nil
}

func parseParam(arg any, idx *int) (map[string]any, error) {
	if arg == nil {
		return nil, nil
	}
	m := make(map[string]any, 1)

	switch arg.(type) {
	case ID:
		m["$"] = []any{arg}
		return m, nil
	case Range:
		m["$"] = []any{arg}
		return m, nil
	}

	v := reflect.ValueOf(arg)
	switch v.Kind() {
	case reflect.Map:
		return arg.(map[string]any), nil
	case reflect.Struct:
		return structToMap(arg), nil
	default:
		m[fmt.Sprintf("%d", *idx+1)] = v.Interface()
		*idx++
		return m, nil
	}
}

func structToMap[T any](content T) map[string]any {
	t := reflect.TypeOf(content)
	nFields := t.NumField()
	m := make(map[string]any, nFields)
	for i := range nFields {
		field := t.Field(i)
		name := strings.ToLower(field.Name)
		if tag, ok := field.Tag.Lookup("db"); ok {
			name = tag
		}

		m[name] = reflect.ValueOf(content).Field(i).Interface()
	}

	return m
}
