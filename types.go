package surgo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	NONE           = "NONE"
	dateTimePrefix = "<datetime>"
)

type (
	ID            []any
	Range         [2]ID
	Record[T any] struct {
		ID   ID
		Data *T
	}
)

// Convert converts the ID to the provided types. This is
// useful when an ID was scanned from a database and needs
// to be converted to the correct types.
func (id ID) Convert(types ...reflect.Kind) ID {
	if len(types) < len(id) {
		panic("not enough types provided")
	}
	for i, v := range id {
		switch types[i] {
		case reflect.String:
			id[i] = fmt.Sprintf("%v", v)
		case reflect.Int:
			kind := reflect.ValueOf(v).Kind()
			if kind == reflect.Float64 {
				id[i] = int(v.(float64))
			} else if kind == reflect.String {
				ni, err := strconv.Atoi(v.(string))
				if err != nil {
					panic("cannot convert string to int")
				}
				id[i] = ni
			}
		case reflect.Float64:
			kind := reflect.ValueOf(v).Kind()
			if kind == reflect.Int {
				id[i] = float64(v.(int))
			} else if kind == reflect.String {
				nf, err := strconv.ParseFloat(v.(string), 64)
				if err != nil {
					panic("cannot convert string to float64")
				}
				id[i] = nf
			}
		case reflect.Bool:
			kind := reflect.ValueOf(v).Kind()
			if kind == reflect.String {
				nb, err := strconv.ParseBool(v.(string))
				if err != nil {
					panic("cannot convert string to bool")
				}
				id[i] = nb
			} else if kind == reflect.Int {
				id[i] = v.(int) != 0
			}
		default:
			id[i] = v
		}
	}

	return id
}

func (id ID) string() string {
	if len(id) == 1 {
		v := id[0]
		switch v.(type) {
		case string:
			return fmt.Sprintf("`%s`", strings.Replace(id[0].(string), "`", "", -1))
		case time.Time:
			return fmt.Sprintf("`%s`", v.(time.Time).Format(time.RFC3339))
		default:
			return fmt.Sprintf("%v", id[0])
		}
	}

	str := "["
	for _, v := range id {
		switch v.(type) {
		case string:
			v = fmt.Sprintf("'%s'", strings.Replace(v.(string), "'", "", -1))
		case time.Time:
			v = fmt.Sprintf("%s'%s'", dateTimePrefix, v.(time.Time).Format(time.RFC3339))
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

func parseTimes(ts any) (string, bool) {
	switch ts.(type) {
	case time.Time:
		return ts.(time.Time).Format(time.RFC3339), true
	case time.Duration:
		total := ts.(time.Duration)
		var unit string
		switch {
		case total >= time.Hour:
			total = total / time.Hour
			unit = "h"
		case total >= time.Minute:
			total = total / time.Minute
			unit = "m"
		case total >= time.Second:
			total = total / time.Second
			unit = "s"
		case total >= time.Millisecond:
			total = total / time.Millisecond
			unit = "ms"
		case total >= time.Microsecond:
			total = total / time.Microsecond
			unit = "us"
		default:
			total = total / time.Nanosecond
			unit = "ns"
		}
		return fmt.Sprintf("%d%s", total, unit), true
	}
	return "", false
}

func stringToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func stringToDuration(s string) (time.Duration, error) {
	var (
		total      int64
		unit       string
		multiplier time.Duration
	)
	if _, err := fmt.Sscanf(s, "%d%s", &total, &unit); err != nil {
		return 0, err
	}

	switch unit {
	case "h":
		multiplier = time.Hour
	case "m":
		multiplier = time.Minute
	case "s":
		multiplier = time.Second
	case "ms":
		multiplier = time.Millisecond
	case "us":
		multiplier = time.Microsecond
	case "ns":
		multiplier = time.Nanosecond
	}
	return time.Duration(total) * multiplier, nil
}

func stringToID(s string) ID {
	s = strings.Trim(s, "`")
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	parts := strings.Split(s, ", ")
	if len(parts) == 1 {
		return ID{s}
	}

	var id ID
	for i, part := range parts {
		// TODO support datetime
		if strings.HasPrefix(part, "'") {
			id[i] = strings.Trim(part, "'")
		} else {
			id[i] = part
		}
	}

	return id
}
