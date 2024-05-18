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
	ID       []any
	Range    [2]ID
	Datetime Date
	Date     struct {
		time.Time
	}
)

func (id ID) string() string {
	if len(id) == 1 {
		v := id[0]
		switch v.(type) {
		case string:
			return fmt.Sprintf("`%s`", strings.Replace(id[0].(string), "`", "", -1))
		case Datetime:
			return fmt.Sprintf("`%s`", v.(Datetime).string())
		case Date:
			return fmt.Sprintf("`%s`", v.(Date).string())
		default:
			return fmt.Sprintf("%v", id[0])
		}
	}

	str := "["
	for _, v := range id {
		switch v.(type) {
		case string:
			v = fmt.Sprintf("'%s'", strings.Replace(v.(string), "'", "", -1))
		case Datetime:
			v = rangedString(v.(Datetime).string())
		case Date:
			v = rangedString(v.(Date).string())
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

func (dt Datetime) string() string {
	return dt.Format(time.RFC3339)
}

func (dt Date) string() string {
	return dt.Format(time.DateOnly)
}

func rangedString(dt string) string {
	return fmt.Sprintf("<datetime>'%s'", dt)
}
