package surgo

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func (db *DB) Query(query string) (interface{}, error) {
	query = strings.TrimSpace(query) + ";"
	return db.db.Query(query, nil)
}

func scan[T any](scan **T, content T) {
	*scan = &content
}

func fields(fields []string) string {
	if len(fields) == 0 {
		return "*"
	}
	return strings.Join(fields, ", ")
}

func omit(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf(" OMIT %s", strings.Join(fields, ", "))
}

func only(condition bool) string {
	if condition {
		return "ONLY "
	}
	return ""
}

func where(condition string) string {
	if condition == "" {
		return ""
	}
	return fmt.Sprintf("WHERE %s ", condition)
}

func group(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("GROUP BY %s ", strings.Join(fields, ", "))
}

func order(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("ORDER BY %s ", strings.Join(fields, ", "))
}

func limit(limit int) string {
	if limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d ", limit)
}

func start(start int) string {
	if start == 0 {
		return ""
	}
	return fmt.Sprintf("START %d ", start)
}

func fetch(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("FETCH %s ", strings.Join(fields, ", "))
}

func timeout(d time.Duration) string {
	if d == 0 {
		return ""
	}
	// TODO: parse duration to ideal measurement
	return fmt.Sprintf("TIMEOUT %dms ", d.Milliseconds())
}

func id(ID string) string {
	if ID == "" {
		return ""
	}
	return fmt.Sprintf(":%s", ID)
}

func parallel(condition bool) string {
	if condition {
		return "PARALLEL"
	}
	return ""
}

func returns(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("RETURN %s ", strings.Join(fields, ", "))
}

func content[T any](content *T) string {
	contentStr := " CONTENT {"
	v := reflect.ValueOf(content).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if fieldType.Name == "ID" {
			continue
		}
		tag := fieldType.Tag.Get("surreal")
		if tag == "" {
			tag = fieldType.Tag.Get("surreal")
			if tag == "" {
				tag = pascaleToSnake(fieldType.Name)
			}
		}
		switch field.Kind() {
		case reflect.String:
			contentStr += fmt.Sprintf(`%s:"%v",`, tag, field.Interface())
		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				contentStr += fmt.Sprintf(`%s:"%v",`, tag, field.Interface().(time.Time).Format(time.RFC3339))
			} else {
				contentStr += fmt.Sprintf("%s:%v,", tag, field.Interface())
			}
		default:
			contentStr += fmt.Sprintf("%s:%v,", tag, field.Interface())
		}
	}

	contentStr = contentStr[:len(contentStr)-1] + "} "

	return contentStr
}

func pascaleToSnake(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToLower(result)
}

func nameOf[T any]() string {
	return pascaleToSnake(reflect.TypeOf(new(T)).Elem().Name())
}
