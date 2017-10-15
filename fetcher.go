package main

import (
	"encoding/json"
	"fmt"
	"regexp"
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
		log.Debug("asking groupcache")
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
		log.Info("asking dbserver")

		var isOk byte = 1 // status: success

		data, err := dbQuery(cfg, log, db, key)
		if err != nil {
			log.Warnf("Query for key %s error: %+v", key, err)
			// TODO: send to error logging channel
			data, err = parseError(err)
			if err != nil {
				data, _ = json.Marshal(err.Error())
			}
			isOk = 0
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
	//	Permit  *string // permit code
	MaxAge int  // cache max age
	IsRO   bool // function is read-only (not volatile)
}

// FuncMap holds map of function definitions
type FuncMap map[string]FuncDef

// -----------------------------------------------------------------------------

func indexFetcher(cfg *AplFlags, log *logger.Log, db *pgx.ConnPool) (index *FuncMap, err error) {
	var rows *pgx.Rows

	// q := fmt.Sprintf("select code, nspname, proname, permit_code, max_age, is_ro from %s()", cfg.ArgIndexFunc)
	q := fmt.Sprintf("select code, nspname, proname, max_age, is_ro from %s()", cfg.ArgIndexFunc)
	log.Debugf("Query: %s", q)

	rows, err = db.Query(q)
	if err != nil {
		return
	}

	ind := FuncMap{}
	for rows.Next() {
		var fmr FuncDef
		var code string
		// err = rows.Scan(&code, &fmr.NspName, &fmr.ProName, &fmr.Permit, &fmr.MaxAge, &fmr.IsRO)
		err = rows.Scan(&code, &fmr.NspName, &fmr.ProName, &fmr.MaxAge, &fmr.IsRO)
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
	var rows *pgx.Rows
	var q string
	var vals []interface{}
	var table *TableRows
	inTx := cfg.BeginFunc != "" // call in transaction
	var lang *string
	var tz *string
	if strings.HasPrefix(key, "[") {
		var args []interface{}
		err = json.Unmarshal([]byte(key), &args)
		if err != nil {
			return
		}
		inTx = false // no transaction for internal calls
		q, vals = PrepareFuncSQL(cfg, args)
	} else {
		var args CallDef
		err = json.Unmarshal([]byte(key), &args)
		if err != nil {
			return
		}
		q, vals, lang, tz = PrepareFuncSQLmap(cfg, args)
	}

	log.Debugf("Query: %s args: %+v", q, vals)

	if inTx {
		// call func after begin
		log.Debug("Run in transaction mode")
		var tx *pgx.Tx
		tx, err = db.Begin()
		if err != nil {
			return
		}
		// Rollback is safe to call even if the tx is already closed, so if
		// the tx commits successfully, this is a no-op
		defer tx.Rollback()

		rows, err = tx.Query("select "+cfg.BeginFunc+"($1,$2)", lang, tz)
		if err != nil {
			return
		}
		rows.Next() // just fetch unneded data or will be "conn is busy"
		if rows.Err() != nil {
			err = rows.Err()
			return
		}
		rows.Close()

		log.Debug("Run main call")
		rows, err = tx.Query(q, vals...)
		if err != nil {
			return
		}
		table, err = FetchSQLResult(rows, log)
		if err != nil {
			return
		}

		err = tx.Commit()
		log.Debug("Call committed")

	} else {
		rows, err = db.Query(q, vals...)
		if err != nil {
			return
		}
		table, err = FetchSQLResult(rows, log)
		defer rows.Close()
	}

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

	proc := args[2].(string) // nil, cache_id, method, args..
	argVals := args[3:]

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

// PrepareFuncSQLmap prepares sql query with named args placeholders
func PrepareFuncSQLmap(cfg *AplFlags, args CallDef) (string, []interface{}, *string, *string) {

	var proc string
	ref := args.Proc
	if ref != nil {
		proc = *ref
	} else {
		// fatal - incorrect map structure
	}

	argValPrep := make([]interface{}, len(args.Args))
	argIDs := make([]string, len(args.Args))

	i := 0
	for k, v := range args.Args {
		argIDs[i] = fmt.Sprintf("%s %s $%d", k, cfg.ArgSyntax, i+1)
		argValPrep[i] = v
		i++
	}

	argIDStr := strings.Join(argIDs, ",")

	var q string
	ref = args.Name
	if ref != nil {
		nsp := *ref
		q = fmt.Sprintf("select * from %s.%s(%s)", nsp, proc, argIDStr)
	} else {
		// use search_path
		q = fmt.Sprintf("select * from %s(%s)", proc, argIDStr)
	}

	return q, argValPrep, args.Lang, args.TZ
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
				if val != nil {
					var raw []byte
					raw, err = json.Marshal(val)
					if err != nil {
						log.Warnf("Value marshal error: %s", err.Error())
						return
					}
					ref := json.RawMessage(raw)
					entry[col] = &ref
				}
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

var reError = regexp.MustCompile(`^ERROR: (.+) \(SQLSTATE (\w+)\)$`)
var reHash = regexp.MustCompile(`^\{.+\}$`)

// Vars holds error fields hash
type Vars map[string]interface{}

// ErrorVars holds returned error struct
type ErrorVars struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    Vars   `json:"data,omitempty"`
}

func parseError(e error) (data []byte, err error) {

	// ERROR: {\"code\" : \"YA014\", \"data\" : {\"login\": \"john\"}} (SQLSTATE P0001)"
	// EMFE - Error Message Format Error
	match := reError.FindStringSubmatch(e.Error())
	if len(match) == 3 {
		if match[2] == "P0001" {
			if reHash.MatchString(match[1]) {
				ev := ErrorVars{}
				err = json.Unmarshal([]byte(match[1]), &ev)
				if err != nil {
					err = fmt.Errorf("%s (EMFE: parse json error: %s)", match[1], err.Error())
				} else {
					data, _ = json.Marshal(ev) // Got application error
				}
			} else {
				err = fmt.Errorf("%s (EMFE: not json hash in message)", match[1])
			}
		} else {
			ev := ErrorVars{Code: match[2], Message: match[1]}
			data, _ = json.Marshal(ev) // Got postgresql error
			return
		}
	} else {
		err = fmt.Errorf("%s (EMFE: Unknown error format)", match)
	}
	return
}
