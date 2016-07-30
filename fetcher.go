package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/groupcache"

	"github.com/LeKovr/dbrpc/workman"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// Processor gets value from cache and converts it into Result struct
func cacheFetcher(log *logger.Log, cacheGroup *groupcache.Group) workman.WorkerFunc {
	// https://github.com/capotej/groupcache-db-experiment
	return func(payload string) workman.Result {
		var data []byte
		log.Printf("asked for %s from groupcache", payload)
		err := cacheGroup.Get(nil, payload,
			groupcache.AllocatingByteSliceSink(&data))
		var res workman.Result
		if err != nil {
			res = workman.Result{Success: false, Error: err.Error()}
		} else {
			d := data[1:]
			if data[0] == 1 { // First byte stores success state (1: true, 0: false)
				raw := json.RawMessage(d)
				res = workman.Result{Success: true, Result: &raw}
			} else {
				res = workman.Result{Success: false, Error: string(d)}
			}
		}
		return res
	}
}

// -----------------------------------------------------------------------------

func dbFetcher(cfg *AplFlags, log *logger.Log, db *sql.DB) groupcache.GetterFunc {
	return func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
		log.Printf("asking for %s from dbserver", key)

		var args []string
		var data []byte

		//err := json.Unmarshal(key, &args)
		json.Unmarshal([]byte(key), &args)

		if args[0] == cfg.ArgDefFunc {

			q := fmt.Sprintf("select * from %s.%s($1)", cfg.Schema, args[0])

			rows, err := db.Query(q, args[1])
			if err != nil {
				return err
			}
			defer rows.Close()

			var res []ArgDef
			for rows.Next() {
				var a ArgDef
				err = rows.Scan(&a.ID, &a.Name, &a.Type, &a.Default, &a.AllowNull)
				if err != nil {
					return err
				}
				res = append(res, a)
			}

			data, err = json.Marshal(res)
			if err != nil {
				return err
			}

		} else {
			q, vals := PrepareFuncSQL(cfg, args)
			log.Printf("Query: %s (%+v)", q, vals)
			rows, err := db.Query(q, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()

			data, err = FetchSQLResult(rows)
			if err != nil {
				return err
			}

		}
		result := []byte{1} // status: success
		result = append(result, data...)

		dd := result[1:]
		log.Printf("Save data: %s", dd)
		dest.SetBytes([]byte(result))
		return nil
	}
}

// -----------------------------------------------------------------------------

// PrepareFuncSQL prepares sql query with args placeholders
func PrepareFuncSQL(cfg *AplFlags, args []string) (string, []interface{}) {
	mtd := args[0]
	argVals := args[1:]

	argValPrep := make([]interface{}, len(argVals))
	argIDs := make([]string, len(argVals))

	for i, v := range argVals {
		argIDs[i] = fmt.Sprintf("$%d", i+1)
		argValPrep[i] = v
	}

	argIDStr := strings.Join(argIDs, ",")

	q := fmt.Sprintf("select * from %s.%s(%s)", cfg.Schema, mtd, argIDStr)

	return q, argValPrep
}

// -----------------------------------------------------------------------------

// FetchSQLResult fetches sql result and marshalls it into json
func FetchSQLResult(rows *sql.Rows) (data []byte, err error) {
	// http://stackoverflow.com/a/29164115
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	count := len(columns)
	var tableData []map[string]interface{}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	data, err = json.Marshal(tableData)
	return
}
