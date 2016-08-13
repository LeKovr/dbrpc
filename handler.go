package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/LeKovr/dbrpc/workman"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// ArgDef holds function argument attributes
type ArgDef struct {
	ID        int32
	Name      string
	Type      string
	Default   *string
	AllowNull bool
}

// FuncArgDef holds set of function argument attributes
type FuncArgDef []ArgDef

// JSON-RPC v2.0 structures
type reqParams map[string]interface{}

type serverRequest struct {
	Method  string    `json:"method"`
	Version string    `json:"jsonrpc"`
	ID      uint64    `json:"id"`
	Params  reqParams `json:"params"`
}

type serverResponse struct {
	ID      uint64           `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *json.RawMessage `json:"error,omitempty"`
}

type respRPCError struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}
type respPGTError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// -----------------------------------------------------------------------------

// FunctionDef creates a job for fetching of function argument definition
func FunctionDef(cfg *AplFlags, log *logger.Log, jc chan workman.Job, method string) (FuncArgDef, interface{}) {

	key := []string{cfg.ArgDefFunc, cfg.Schema + "." + method}

	payload, _ := json.Marshal(key)
	respChannel := make(chan workman.Result)

	work := workman.Job{Payload: string(payload), Result: respChannel}

	// Push the work onto the queue.
	jc <- work

	resp := <-respChannel
	log.Debugf("Got def: %s", resp.Result)
	if !resp.Success {
		return nil, resp.Error
	}

	var res FuncArgDef
	err := json.Unmarshal(*resp.Result, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// -----------------------------------------------------------------------------

// FunctionResult creates a job for fetching of function result
func FunctionResult(jc chan workman.Job, payload string) workman.Result {

	respChannel := make(chan workman.Result)
	// let's create a job with the payload
	work := workman.Job{Payload: payload, Result: respChannel}

	// Push the work onto the queue.
	jc <- work

	resp := <-respChannel
	return resp
}

// -----------------------------------------------------------------------------

func httpHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		log.Debugf("Request method: %s", r.Method)
		if r.Method == "GET" {
			getContextHandler(cfg, log, jc, w, r, true)
		} else if r.Method == "HEAD" {
			getContextHandler(cfg, log, jc, w, r, false) // Like get but without data
		} else if r.Method == "POST" && r.URL.Path == cfg.Prefix {
			postContextHandler(cfg, log, jc, w, r)
		} else if r.Method == "POST" {
			postgrestContextHandler(cfg, log, jc, w, r)
		} else if r.Method == "OPTIONS" {
			optionsContextHandler(cfg, log, jc, w, r)
		} else {
			e := fmt.Sprintf("Unsupported request method: %s", r.Method)
			log.Warn(e)
			http.Error(w, e, http.StatusNotImplemented)
		}
	}
}

// -----------------------------------------------------------------------------

func getContextHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job, w http.ResponseWriter, r *http.Request, reply bool) {
	start := time.Now()

	method := strings.TrimPrefix(r.URL.Path, cfg.Prefix)

	argDef, errd := FunctionDef(cfg, log, jc, method)
	if errd != nil {
		// Warning was when fetched from db
		log.Infof("Method %s load def error: %s", method, errd)
		http.NotFound(w, r)
		return
	}

	key := []*string{&method}
	r.ParseForm()

	f404 := []string{}
	for _, a := range argDef {
		v := r.Form[a.Name]

		if len(v) == 0 {
			if !a.AllowNull && a.Default == nil {
				f404 = append(f404, a.Name)
			} else if a.Default != nil {
				log.Debugf("Arg: %s use default", a.Name)
				break // use defaults
			}
			key = append(key, nil) // TODO: nil does not replaced with default
		} else if strings.HasSuffix(a.Type, "[]") {
			// convert array into string
			// TODO: escape ","
			s := "{" + strings.Join(v, ",") + "}"
			key = append(key, &s)
		} else {
			key = append(key, &v[0])
		}

		log.Debugf("Arg: %+v (%d)", v, len(f404))
	}
	var result workman.Result
	if len(f404) > 0 {
		result = workman.Result{Success: false, Error: fmt.Sprintf("Required parameter(s) %+v not found", f404)}
	} else {
		payload, _ := json.Marshal(key)
		log.Debugf("Args: %s", string(payload))
		result = FunctionResult(jc, string(payload))
	}

	if reply {
		out, err := json.Marshal(result)
		if err != nil {
			log.Warnf("Marshall error: %+v", err)
			e := workman.Result{Success: false, Error: err.Error()}
			out, _ = json.Marshal(e)
		}
		setMetric(w, start, http.StatusOK)
		w.Write(out)
		//w.Write([]byte("\n"))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// -----------------------------------------------------------------------------

func setMetric(w http.ResponseWriter, start time.Time, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Elapsed", fmt.Sprint(time.Since(start)))
	w.WriteHeader(http.StatusOK)
}

// -----------------------------------------------------------------------------

func optionsContextHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job, w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")
	var host string
	if origin != "" && len(cfg.Hosts) > 0 { // lookup if host is allowed
		for _, h := range cfg.Hosts {
			if origin == h {
				host = h
				break
			}
		}
	} else {
		host = origin
	}
	if origin != "" && host == "" {
		log.Warningf("Unregistered request source: %s", origin)
		http.Error(w, "Origin not registered", http.StatusForbidden)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", host)
	w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.WriteHeader(http.StatusOK)
	return
}

func getRaw(data interface{}) *json.RawMessage {
	j, _ := json.Marshal(data)
	raw := json.RawMessage(j)
	return &raw
}

// -----------------------------------------------------------------------------

// postContextHandler serve JSON-RPC envelope
func postContextHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job, w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	data, _ := ioutil.ReadAll(r.Body)
	req := serverRequest{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		e := fmt.Sprintf("json parse error: %s", err)
		log.Warn(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	resultRPC := serverResponse{ID: req.ID, Version: req.Version}

	argDef, errd := FunctionDef(cfg, log, jc, req.Method)
	if errd != nil {
		log.Warnf("Method %s load def error: %s", req.Method, errd)
		resultRPC.Error = getRaw(respRPCError{Code: -32601, Message: errd.(string)})
	} else {
		// Load args
		r.ParseForm()
		key, f404 := fetchArgs(argDef, req.Params, req.Method)
		if len(f404) > 0 {
			resultRPC.Error = getRaw(respRPCError{Code: -32602, Message: "Required parameter(s) not found", Data: getRaw(f404)})
		} else {
			payload, _ := json.Marshal(key)
			log.Debugf("Args: %s", string(payload))
			res := FunctionResult(jc, string(payload))
			if res.Success {
				resultRPC.Result = res.Result
			} else {
				resultRPC.Error = getRaw(respRPCError{Code: -32603, Message: "Internal Error", Data: getRaw(res.Error)})
			}
		}
	}

	out, err := json.Marshal(resultRPC)
	if err != nil {
		log.Warnf("Marshall error: %+v", err)
		resultRPC.Result = nil
		resultRPC.Error = getRaw(respRPCError{Code: -32603, Message: "Internal Error", Data: getRaw(err.Error())})

		out, _ = json.Marshal(resultRPC)
	}
	log.Debugf("JSON Resp: %s", string(out))

	setMetric(w, start, http.StatusOK)
	w.Write(out)
	//w.Write([]byte("\n"))
}

// -----------------------------------------------------------------------------

// postgrestContextHandler serve JSON-RPC envelope
// 404 when method not found
func postgrestContextHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job, w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	method := strings.TrimPrefix(r.URL.Path, cfg.Prefix)
	log.Debugf("postgrest call for %s", method)

	argDef, errd := FunctionDef(cfg, log, jc, method)
	if errd != nil {
		log.Warnf("Method %s load def error: %s", method, errd)
		http.NotFound(w, r)
		return
	}
	resultStatus := http.StatusOK

	req := reqParams{}
	var resultRPC interface{}

	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &req)

	if err != nil {
		e := fmt.Sprintf("json parse error: %s", err)
		log.Warn(e)
		resultRPC = respPGTError{Message: "Cannot parse request payload", Details: e}
		resultStatus = http.StatusBadRequest
	} else {
		// Load args
		key, f404 := fetchArgs(argDef, req, method)
		if len(f404) > 0 {
			resultRPC = respPGTError{Code: "42883", Message: "Required parameter(s) not found", Details: strings.Join(f404, ", ")}
			resultStatus = http.StatusBadRequest
		} else {
			payload, _ := json.Marshal(key)
			log.Debugf("Args: %s", string(payload))
			res := FunctionResult(jc, string(payload))
			if res.Success {
				resultRPC = res.Result
			} else {
				resultRPC = respPGTError{Message: "Method call error", Details: res.Error.(string)}
				resultStatus = http.StatusBadRequest // TODO: ?
			}
		}
	}

	out, err := json.Marshal(resultRPC)
	if err != nil {
		log.Warnf("Marshall error: %+v", err)
		resultRPC = respPGTError{Message: "Method result marshall error", Details: err.Error()}
		resultStatus = http.StatusBadRequest // TODO: ?
		out, _ = json.Marshal(resultRPC)
	}
	setMetric(w, start, resultStatus)
	log.Debugf("JSON Resp: %s", string(out))
	w.Write(out)
	//w.Write([]byte("\n"))
}

func fetchArgs(argDef FuncArgDef, req reqParams, method string) ([]*string, []string) {

	key := []*string{&method}
	f404 := []string{}

	for _, a := range argDef {
		v, ok := req[a.Name]
		if !ok {
			if !a.AllowNull && a.Default == nil {
				f404 = append(f404, a.Name)
			} else if a.Default != nil {
				//log.Debugf("Arg: %s use default", a.Name)
				break // use defaults
			}
			key = append(key, nil) // TODO: nil does not replaced with default
		} else if strings.HasSuffix(a.Type, "[]") {
			// wait slice
			s := reflect.ValueOf(v)
			if s.Kind() != reflect.Slice {
				// string or {string}
				vs := v.(string)

				// // convert scalar to postgres array
				// asArray = regexp.MustCompile(`^\{.+\}$`)
				// if !asArray.MatchString(vs) {
				// 	vs = "{" + vs + "}"
				// }
				key = append(key, &vs)
			} else {
				// slice
				ret := make([]string, s.Len())

				for i := 0; i < s.Len(); i++ {
					ret[i] = s.Index(i).Interface().(string)
					//	log.Printf("====== %+v", ret[i])
				}
				// convert array into string
				// TODO: escape ","
				ss := "{" + strings.Join(ret, ",") + "}"
				key = append(key, &ss)
			}
		} else {
			s := v.(string)
			key = append(key, &s)
		}

		//log.Debugf("Arg: %+v (%d)", v, len(f404))
	}
	return key, f404
}
