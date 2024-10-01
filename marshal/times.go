package marshal

import (
	"fmt"
	"regexp"
	"strconv"
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
		return ts.(time.Time).Format(time.RFC3339)
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

func stringToTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func stringToDuration(s string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)(y|w|d|h|ms|m|s|µs|us|ns)`)
	matches := re.FindAllStringSubmatch(s, -1)

	var duration time.Duration
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
		unit := match[2]

		duration += time.Duration(value) * units[unit]
	}

	return duration, nil
}
