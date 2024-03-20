package surgo

import (
	"fmt"
	"reflect"
	"strings"
)

func (db *DB) parseParams(args any) (map[string]any, error) {
	kind := reflect.ValueOf(args).Kind()
	switch kind {
	case reflect.Slice:
		return sliceToMap(args), nil
	case reflect.Array:
		return sliceToMap(args), nil
	case reflect.Map:
		return args.(map[string]any), nil
	case reflect.Struct:
		return structToMap(args), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", kind.String())
	}
}

func sliceToMap(slice any) map[string]any {
	v := reflect.ValueOf(slice)
	m := make(map[string]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		m[fmt.Sprintf("%d", i+1)] = v.Index(i).Interface()
	}

	return m
}

func structToMap[T any](content T) map[string]any {
	t := reflect.TypeOf(content)
	nFields := t.NumField()
	m := make(map[string]any, nFields)
	for i := range nFields {
		field := t.Field(i)
		name := field.Name
		if tag, ok := field.Tag.Lookup("db"); ok {
			name = tag
		}

		m[name] = reflect.ValueOf(content).Field(i).Interface()
	}

	return m
}

func scanData(scan any, res Result) error {
	v := reflect.ValueOf(scan)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("scan must be a pointer")
	}

	data, ok := res.Data.([]any)
	if !ok {
		data, ok := res.Data.(any)
		if !ok {
			return fmt.Errorf("unexpected type %T", res.Data)
		}
		if v.Elem().Kind() == reflect.Slice {
			slice := reflect.MakeSlice(v.Elem().Type(), 1, 1)
			err := fillData(slice.Index(0).Addr().Interface(), data)
			if err != nil {
				return err
			}
			v.Elem().Set(slice)
			return nil
		} else {
			return fillData(scan, data)
		}
	}

	if v.Elem().Kind() == reflect.Slice {
		slice := reflect.MakeSlice(v.Elem().Type(), len(data), len(data))
		for i, d := range data {
			err := fillData(slice.Index(i).Addr().Interface(), d)
			if err != nil {
				return err
			}
		}
		v.Elem().Set(slice)
	} else {
		return fillData(scan, data[0])
	}

	return nil
}

func fillData(scan any, data any) error {
	m, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map, got %T", data)
	}

	if len(m) == 1 {
		for _, v := range m {
			setVal(reflect.ValueOf(scan).Elem(), reflect.ValueOf(v))
			return nil
		}
	}

	t := reflect.TypeOf(scan).Elem()
	for k, v := range m {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Tag.Get("db") == k || field.Name == k || strings.ToLower(field.Name) == k {
				setVal(reflect.ValueOf(scan).Elem().Field(i), reflect.ValueOf(v))
			}
		}
	}

	return nil
}

func setVal(dest, src reflect.Value) {
	if dest.Kind() == reflect.Int && src.Kind() == reflect.Float64 {
		dest.SetInt(int64(src.Float()))
	} else {
		dest.Set(src)
	}
}
