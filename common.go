package surgo

import (
	"fmt"
	"strings"
	"time"
)

const (
	NONE = "NONE"
)

type (
	ID    []any
	Range [2]ID
)

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
			v = fmt.Sprintf("<datetime>'%s'", v.(time.Time).Format(time.RFC3339))
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
