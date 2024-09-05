package surgo

import (
	"reflect"
	"strings"
)

func parseVars(vars map[string]any) map[string]any {
	for k, v := range vars {
		if isTime(v) {
			vars[k] = parseTimes(v)
		} else if isStruct(v) {
			vars[k] = structToMap(v)
		} else {
			vars[k] = v
		}
	}

	return vars
}

// isStruct checks if the given value is a struct or a pointer to a struct.
func isStruct(x any) bool {
	t := reflect.TypeOf(x)
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

func structToMap(x any) map[string]any {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]any)
	for i := range t.NumField() {
		field := t.Field(i)
		name := field.Name
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = field.Tag.Get(fallbackTag)
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
		if isStruct(val) {
			m[name] = structToMap(val)
		} else {
			m[name] = parseTimes(val)
		}
	}

	return m
}
