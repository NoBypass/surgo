package marshal

import (
	"fmt"
	"reflect"
	"time"
)

var units = map[string]time.Duration{
	"y":  365 * 24 * time.Hour,
	"w":  7 * 24 * time.Hour,
	"d":  24 * time.Hour,
	"h":  time.Hour,
	"m":  time.Minute,
	"s":  time.Second,
	"ms": time.Millisecond,
	"µs": time.Microsecond,
	"us": time.Microsecond,
	"ns": time.Nanosecond,
}

func isTime(ts any) bool {
	switch ts.(type) {
	case time.Time, time.Duration:
		return true
	default:
		return false
	}
}

func parseTimes(ts any) any {
	switch ts.(type) {
	case time.Time:
		return ts.(time.Time).Format("d\"2006-01-02T15:04:05Z07:00\"")
	case time.Duration:
		var result string
		d := ts.(time.Duration)
		for unit, duration := range units {
			if d >= duration {
				amount := d / duration
				d -= amount * duration
				result += fmt.Sprintf("%d%s", amount, unit)
			}
		}
		return result
	default:
		return ts
	}
}

func (m *Marshaler) tagOf(field reflect.StructField) string {
	dbTag := field.Tag.Get("db")
	if dbTag == "" {
		dbTag = field.Tag.Get(string(*m))
	}
	if dbTag == "" || dbTag[0] == ',' {
		dbTag = field.Name + dbTag
	}
	return dbTag
}
