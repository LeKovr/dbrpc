package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/jackc/pgx.v2"

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
		} else if len(data) == 0 {
			res = workman.Result{Success: false, Error: "This internal error must be catched earlier. Please contact vendor"}
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

func dbFetcher(cfg *AplFlags, log *logger.Log, db *pgx.Conn) groupcache.GetterFunc {
	return func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
		log.Printf("asking for %s from dbserver", key)

		var isOk byte = 1 // status: success

		data, err := dbQuery(cfg, log, db, key)
		if err != nil {
			log.Warnf("Query for key %s error: %+v", key, err)
			isOk = 0
			data, _ = json.Marshal(err.Error())
		}

		result := []byte{isOk}
		result = append(result, data...)

		dd := result[1:]
		log.Printf("Save data: %s", dd)
		dest.SetBytes([]byte(result))
		return nil
	}
}

func dbQuery(cfg *AplFlags, log *logger.Log, db *pgx.Conn, key string) (data []byte, err error) {
	var args []string
	var rows *pgx.Rows

	err = json.Unmarshal([]byte(key), &args)
	if err != nil {
		return
	}

	if args[0] == cfg.ArgDefFunc {

		q := fmt.Sprintf("select * from %s.%s($1)", cfg.Schema, args[0])

		rows, err = db.Query(q, args[1])
		defer rows.Close()
		if err != nil {
			return
		}

		var res []ArgDef
		for rows.Next() {
			var a ArgDef
			err = rows.Scan(&a.ID, &a.Name, &a.Type, &a.Default, &a.AllowNull)
			if err != nil {
				return
			}
			res = append(res, a)
		}
		if rows.Err() != nil {
			err = rows.Err()
			return
		}
		log.Printf("Func def: %s (%+v)", args[1], res)

		data, err = json.Marshal(res)
		if err != nil {
			return
		}

	} else {
		q, vals := PrepareFuncSQL(cfg, args)
		log.Printf("Query: %s (%+v)", q, vals)
		rows, err = db.Query(q, vals...)
		if err != nil {
			return
		}

		data, err = FetchSQLResult(rows, log)
		defer rows.Close()
		if err != nil {
			return
		}

	}
	return

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
func FetchSQLResult(rows *pgx.Rows, log *logger.Log) (data []byte, err error) {
	// http://stackoverflow.com/a/29164115
	columnDefs := rows.FieldDescriptions()
	//log.Debugf("=========== %+v", columnDefs)
	columns := []string{}
	types := []string{}
	for _, c := range columnDefs {
		columns = append(columns, c.Name)
		types = append(types, c.DataTypeName)
	}

	var tableData []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Warningf("Value fetch error: %s", err.Error())
		}
		log.Debugf("Values: %+v", values)

		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if types[i] == "json" || types[i] == "jsonb" {
				raw := fmt.Sprintf("%s", val)
				ref := json.RawMessage(raw)
				entry[col] = &ref
			} else {
				v = val
				entry[col] = v
			}
		}
		tableData = append(tableData, entry)

	}
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		err = rows.Err()
		return
	}
	data, err = json.Marshal(tableData)
	return
}
