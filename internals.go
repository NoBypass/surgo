package surgo

import (
	"fmt"
	"strings"
	"time"
)

func (db *DB) Query(query string) (interface{}, error) {
	println(query)
	return nil, nil
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
	return fmt.Sprintf(":%s ", ID)
}

func parallel(condition bool) string {
	if condition {
		return "PARALLEL"
	}
	return ""
}
