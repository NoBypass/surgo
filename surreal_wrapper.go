package surgo

import (
	"fmt"
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
				Error: fmt.Errorf(m["result"].(string)),
			}
		} else {
			resSlice[i] = Result{
				Data:  m["result"],
				Error: nil,
			}
		}

		resSlice[i].Duration = d
		resSlice[i].Query = Query{query, params}
	}
	return resSlice, nil
}
