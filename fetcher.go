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

// TableRow holds query result row
type TableRow map[string]interface{}

// TableRows holds slice of query result rows
type TableRows []TableRow

// -----------------------------------------------------------------------------

// Processor gets value from cache and converts it into Result struct
func cacheFetcher(log *logger.Log, cacheGroup *groupcache.Group) workman.WorkerFunc {
	// https://github.com/capotej/groupcache-db-experiment
	return func(payload string) workman.Result {
		var data []byte
		log.Debugf("asked for %s from groupcache", payload)
		err := cacheGroup.Get(nil, payload,
			groupcache.AllocatingByteSliceSink(&data))
		var res workman.Result
		if err != nil {
			res = workman.Result{Success: false, Error: err.Error()}
		} else if len(data) == 0 {
			res = workman.Result{Success: false, Error: "This internal error must be catched earlier. Please contact vendor"}
		} else {
			d := data[1:]
			raw := json.RawMessage(d)
			ok := data[0] == 1
			if ok { // First byte stores success state (1: true, 0: false)
				res = workman.Result{Success: ok, Result: &raw}
			} else {
				res = workman.Result{Success: ok, Error: &raw}
			}
		}
		return res
	}
}

// -----------------------------------------------------------------------------

func dbFetcher(cfg *AplFlags, log *logger.Log, db *pgx.ConnPool) groupcache.GetterFunc {
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
		log.Debugf("Save data: %s", dd)
		dest.SetBytes([]byte(result))
		return nil
	}
}

// -----------------------------------------------------------------------------

// FuncDef holds function definition
type FuncDef struct {
	NspName string // function namespace
	ProName string // function name
	// ToDo: cache, permission etc
}

// FuncMap holds map of function definitions
type FuncMap map[string]FuncDef

// -----------------------------------------------------------------------------

func indexFetcher(cfg *AplFlags, log *logger.Log, db *pgx.ConnPool) (index *FuncMap, err error) {
	var rows *pgx.Rows

	q := fmt.Sprintf("select code, nspname, proname from %s()", cfg.ArgIndexFunc)
	log.Debugf("Query: %s", q)

	rows, err = db.Query(q)
	if err != nil {
		return
	}

	ind := FuncMap{}
	for rows.Next() {
		var fmr FuncDef
		var code string
		err = rows.Scan(&code, &fmr.NspName, &fmr.ProName) // ToDo: cache, permission etc
		if err != nil {
			log.Warnf("Value fetch error: %s", err.Error())
			return
		}
		ind[code] = fmr
	}
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		err = rows.Err()
		return
	}
	log.Debugf("Index: %+v", ind)
	index = &ind
	return
}

// -----------------------------------------------------------------------------

func dbQuery(cfg *AplFlags, log *logger.Log, db *pgx.ConnPool, key string) (data []byte, err error) {
	var args []interface{}
	var rows *pgx.Rows

	err = json.Unmarshal([]byte(key), &args)
	if err != nil {
		return
	}

	q, vals := PrepareFuncSQL(cfg, args)
	log.Debugf("Query: %s", q)
	rows, err = db.Query(q, vals...)
	if err != nil {
		return
	}

	table, err := FetchSQLResult(rows, log)
	defer rows.Close()
	if err != nil {
		return
	}

	data, err = json.Marshal(*table)
	if err != nil {
		return
	}

	return
}

// -----------------------------------------------------------------------------

// PrepareFuncSQL prepares sql query with args placeholders
func PrepareFuncSQL(cfg *AplFlags, args []interface{}) (string, []interface{}) {
	proc := args[1].(string)
	argVals := args[2:]

	argValPrep := make([]interface{}, len(argVals))
	argIDs := make([]string, len(argVals))

	for i, v := range argVals {
		argIDs[i] = fmt.Sprintf("$%d", i+1)
		argValPrep[i] = v
	}

	argIDStr := strings.Join(argIDs, ",")

	var q string
	if args[0] != nil {
		nsp := args[0].(string)
		q = fmt.Sprintf("select * from %s.%s(%s)", nsp, proc, argIDStr)
	} else {
		// use search_path
		q = fmt.Sprintf("select * from %s(%s)", proc, argIDStr)
	}

	return q, argValPrep
}

// -----------------------------------------------------------------------------

// FetchSQLResult fetches sql result and marshalls it into json
func FetchSQLResult(rows *pgx.Rows, log *logger.Log) (data *TableRows, err error) {
	// http://stackoverflow.com/a/29164115
	columnDefs := rows.FieldDescriptions()
	//log.Debugf("=========== %+v", columnDefs)
	columns := []string{}
	types := []string{}
	for _, c := range columnDefs {
		columns = append(columns, c.Name)
		types = append(types, c.DataTypeName)
	}

	var tableData TableRows
	for rows.Next() {
		var values []interface{}
		values, err = rows.Values()
		if err != nil {
			log.Warnf("Value fetch error: %s", err.Error())
			return
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
	if tableData == nil {
		log.Warn("Empty result")
		tableData = TableRows{} // empty result is empty array
	}
	data = &tableData
	return
}
