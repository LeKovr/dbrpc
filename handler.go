package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LeKovr/dbrpc/workman"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// ArgDef holds function argument attributes
type ArgDef struct {
	ID        int
	Name      string
	Type      string
	Default   *string
	AllowNull bool
}

// FuncArgDef holds set of function argument attributes
type FuncArgDef []ArgDef

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
	log.Printf("Got def: %s", resp.Result)
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
		if r.Method == "GET" {
			getContextHandler(cfg, log, jc, w, r)
		} else if r.Method == "POST" {
			postContextHandler(cfg, log, jc, w, r)
		} else if r.Method == "OPTIONS" {
			optionsContextHandler(cfg, log, jc, w, r)
		} else {
			e := fmt.Sprintf("Unsupported request method: %s", r.Method)
			log.Warn(e)
			http.Error(w, e, http.StatusNotImplemented)
			return
		}
	}
}

// -----------------------------------------------------------------------------

func getContextHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job, w http.ResponseWriter, r *http.Request) {

	//	method := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, cfg.Prefix), "/")
	method := strings.TrimPrefix(r.URL.Path, cfg.Prefix)
	log.Printf("GotREquest %s (%s)", method, r.URL.Path)

	argDef, errd := FunctionDef(cfg, log, jc, method)
	if errd != nil {
		log.Printf("mtd def error: %s", errd)
		http.NotFound(w, r)
		return
	}

	key := []string{method}
	r.ParseForm()

	f404 := []string{}
	for _, a := range argDef {
		v := r.Form[a.Name]

		if len(v) == 0 {
			if !a.AllowNull && a.Default == nil {
				f404 = append(f404, a.Name)
			}
			key = append(key, "") // TODO: nil
		} else {
			key = append(key, v[0]) // TODO: array support
		}

		log.Debugf("Arg: %+v (%d)", v, len(f404))
	}
	var result workman.Result
	if len(f404) > 0 {
		result = workman.Result{Success: false, Error: fmt.Sprintf("Required parameter(s) %+v not found", f404)}
	} else {
		payload, _ := json.Marshal(key)
		result = FunctionResult(jc, string(payload))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	out, err := json.Marshal(result)
	if err != nil {
		log.Println("Marshall error: ", err)
		e := workman.Result{Success: false, Error: err.Error()}
		out, _ = json.Marshal(e)
	}
	w.Write(out)
	w.Write([]byte("\n"))
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

func postContextHandler(cfg *AplFlags, log *logger.Log, jc chan workman.Job, w http.ResponseWriter, r *http.Request) {

	data, _ := ioutil.ReadAll(r.Body)
	req := serverRequest{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		e := fmt.Sprintf("json parse error: %s", err)
		log.Warn(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	resultRPC := serverResponse{ID: req.ID, Version: req.Version}

	argDef, errd := FunctionDef(cfg, log, jc, req.Method)
	if errd != nil {
		log.Printf("mtd def error: %s", errd)
		resultRPC.Error = getRaw(respRPCError{Code: -32601, Message: errd.(string)})
	} else {
		// Load args
		key := []string{req.Method}
		r.ParseForm()
		f404 := []string{}
		for _, a := range argDef {
			v, ok := req.Params[a.Name]
			if !ok {
				if !a.AllowNull && a.Default == nil {
					f404 = append(f404, a.Name)
				}
				key = append(key, "") // TODO: nil
			} else {
				key = append(key, v.(string))
			}
		}

		if len(f404) > 0 {
			resultRPC.Error = getRaw(respRPCError{Code: -32602, Message: "Required parameter(s) not found", Data: getRaw(f404)})
		} else {
			payload, _ := json.Marshal(key)
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
		log.Println("Marshall error: ", err)
		resultRPC.Result = nil
		resultRPC.Error = getRaw(respRPCError{Code: -32603, Message: "Internal Error", Data: getRaw(err.Error())})

		out, _ = json.Marshal(resultRPC)
	}
	log.Debugf("JSON Resp: %s", string(out))
	w.Write(out)
	w.Write([]byte("\n"))
}

// JSON-RPC v2.0 structures

type serverRequest struct {
	Method  string                 `json:"method"`
	Version string                 `json:"jsonrpc"`
	ID      uint64                 `json:"id"`
	Params  map[string]interface{} `json:"params"`
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
