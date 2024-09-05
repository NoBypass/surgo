package surgo

import (
	"fmt"
	"time"
)

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
		return ts.(time.Time).Format(time.RFC3339)
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
		return fmt.Sprintf("%d%s", total, unit)
	default:
		return ts
	}
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

func (db *DB) respToResult(resp any) ([]Result, error) {
	respSlice, ok := resp.([]any)
	if !ok {
		respSlice = []any{resp}
	}
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
				Error: fmt.Errorf("%w: %s", newErrQuery(err), m["result"].(string)),
			}
		} else {
			res := Result{
				Data:  m["result"],
				Error: nil,
			}

			if res.Data == nil {
				res.Error = newErrNoResult(fmt.Errorf("no result found"))
			}

			resSlice[i] = res
		}

		resSlice[i].Duration = d
	}

	return resSlice, nil
}
