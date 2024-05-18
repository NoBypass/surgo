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
		currIdx++
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
		*idx--
		return arg.(map[string]any), nil
	case reflect.Struct:
		*idx--
		return structToMap(arg), nil
	default:
		m[fmt.Sprintf("%d", *idx+1)] = v.Interface()
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

func (id ID) string() string {
	if len(id) == 1 {
		switch id[0].(type) {
		case string:
			return fmt.Sprintf("`%s`", strings.Replace(id[0].(string), "`", "", -1))
		default:
			return fmt.Sprintf("%v", id[0])
		}
	}

	str := "["
	for _, v := range id {
		switch v.(type) {
		case string:
			v = fmt.Sprintf("'%s'", strings.Replace(v.(string), "'", "", -1))
		default:
			v = fmt.Sprintf("%v", v)
		}

		str += fmt.Sprintf("%v, ", v)
	}
	return str[:len(str)-2] + "]"
}

func (r Range) string() string {
	return fmt.Sprintf("%s..%s", r[0].string(), r[1].string())
}
