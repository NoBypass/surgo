package surgo

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func (db *DB) query(query string, params map[string]any) ([]Result, error) {
	if !strings.HasSuffix(query, ";") {
		query = query + ";"
	}
	resp, err := db.DB.Query(query, params)
	if err != nil {
		return nil, err
	}

	respSlice := resp.([]any)
	resSlice := make([]Result, len(respSlice))
	for i, s := range respSlice {
		m := s.(map[string]any)
		d, err := time.ParseDuration(m["time"].(string))
		if err != nil {
			return nil, err
		}

		if m["status"] == "ERR" {
			resSlice[i] = Result{
				Data:  nil,
				Error: fmt.Errorf(m["result"].(string)),
			}
		} else {
			resSlice[i] = Result{
				Data:  m["result"],
				Error: nil,
			}
		}

		resSlice[i].Duration = d
		resSlice[i].Query = Query{
			query, params,
		}
	}
	return resSlice, nil
}

func parseParams(args []any) (map[string]any, error) {
	if len(args) == 0 {
		return nil, nil
	} else if len(args) == 1 {
		return parseParam(args[0], 0)
	}

	m := make(map[string]any)
	for i, a := range args {
		nm, err := parseParam(a, i)
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

func parseParam(arg any, idx int) (map[string]any, error) {
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
		m[fmt.Sprintf("%d", idx+1)] = v.Interface()
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
		if tag, ok := field.Tag.Lookup("DB"); ok {
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
			if field.Tag.Get("DB") == k || field.Name == k || strings.ToLower(field.Name) == k {
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
	return str[:len(str)-1] + "]"
}

func (r Range) string() string {
	return fmt.Sprintf("%s..%s", r[0].string(), r[1].string())
}
