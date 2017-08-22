package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/jackc/pgx.v2"

	"github.com/LeKovr/dbrpc/jwtutil"
	"github.com/LeKovr/dbrpc/workman"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// ArgDef holds function argument attributes
type ArgDef struct {
	ID       int32   `json:"id"`
	Name     string  `json:"arg"`
	Type     string  `json:"type"`
	Default  *string `json:"def_val"`
	Required bool    `json:"required"`
}

// FuncArgDef holds slice of function argument attributes
type FuncArgDef []ArgDef

// RPCServer holds server attributes
type RPCServer struct {
	cacheID uint32
	cfg     *AplFlags
	log     *logger.Log
	jc      chan workman.Job
	funcs   *FuncMap
	started int
	JWT     *jwtutil.App
	lock    *sync.RWMutex
}

// -----------------------------------------------------------------------------

// LoadFuncs stores func attr with locking
func (s *RPCServer) LoadFuncs(db *pgx.ConnPool) {
	fm, err := indexFetcher(s.cfg, s.log, db)
	if err != nil {
		s.log.Fatal("Error loading function index:", err)
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cacheID++
	s.funcs = fm
}

// -----------------------------------------------------------------------------

// Setup initializes RPC server
func (s *RPCServer) Setup(db *pgx.ConnPool) {
	s.cacheID = 0
	s.lock = &sync.RWMutex{}
	s.LoadFuncs(db)
	go s.ListenCounter(db)
}

// -----------------------------------------------------------------------------

// ListenCounter listens db event & resets cache
func (s *RPCServer) ListenCounter(db *pgx.ConnPool) {
	conn, err := db.Acquire()
	if err != nil {
		s.log.Fatal("Error acquiring connection:", err)
	}
	defer db.Release(conn)

	conn.Listen(s.cfg.CacheResetEvent)

	for {
		notification, err := conn.WaitForNotification(60 * time.Second)
		if err != nil && err != pgx.ErrNotificationTimeout {
			// error, no timeout
			s.log.Warn("Error waiting for notification:", err)
			time.Sleep(time.Second * 10) // sleep & repeat
		} else if err == nil {
			// no error, no timeout
			s.log.Warnf("Cache reset event received with payload: %s", notification.Payload)
			//fmt.Println("PID:", notification.Pid, "Channel:", notification.Channel, "Payload:", notification.Payload)
			s.LoadFuncs(db)
		}

	}
}

// -----------------------------------------------------------------------------

// CacheID returns current cache id
func (s RPCServer) CacheID() uint32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.cacheID
}

// -----------------------------------------------------------------------------

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
	Message string           `json:"message"`
	Code    string           `json:"code,omitempty"`
	Details *json.RawMessage `json:"details,omitempty"`
}

// -----------------------------------------------------------------------------

// FunctionArgDef creates a job for fetching of function argument definition
func (s RPCServer) FunctionArgDef(nsp, proc string) (FuncArgDef, interface{}) {

	id := string(s.CacheID())
	key := []*string{nil, &id, &s.cfg.ArgDefFunc, &nsp, &proc}

	payload, _ := json.Marshal(key)
	respChannel := make(chan workman.Result)

	work := workman.Job{Payload: string(payload), Result: respChannel}

	// Push the work onto the queue.
	s.jc <- work

	resp := <-respChannel
	s.log.Debugf("Got def (%v): %s", resp.Success, resp.Result)
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

func (s *RPCServer) httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := s.log
		cfg := s.cfg
		defer r.Body.Close()
		log.Debugf("Request method: %s %s", r.Method, r.RequestURI)

		if origin := r.Header.Get("Origin"); origin != "" {

			log.Debugf("Lookup origin %s in %+v", origin, cfg.Hosts)
			if !stringExists(cfg.Hosts, origin, "*") {
				log.Warningf("Unregistered request source: %s", origin)
				http.Error(w, "Origin not registered", http.StatusForbidden)
				return
			}
			methodsAllowed := "origin, content-type, accept, keep-alive, user-agent, x-requested-with, x-token, authorization"
			if cfg.LangHeader != "" {
				methodsAllowed = methodsAllowed + ", " + cfg.LangHeader
			}
			if cfg.TZHeader != "" {
				methodsAllowed = methodsAllowed + ", " + cfg.TZHeader
			}
			w.Header().Add("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Credentials", "true") // TODO
			w.Header().Add("Access-Control-Allow-Headers", methodsAllowed)
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		}
		var session *jwtutil.Session
		if auth := r.Header.Get("Authorization"); auth != "" {
			// log.Debugf("Lookup auth from %s", auth)
			authData := strings.TrimPrefix(auth, "Bearer ") // todo: error if eq
			se, err := s.JWT.Parse(authData)
			if err != nil {
				log.Warningf("Auth parse error: %s", err)
				http.Error(w, "Auth error", http.StatusForbidden)
				return
			}
			log.Debugf("Lookup auth got %v", se)
			session = se
		} else {
			session = &jwtutil.Session{}
		}

		// setup language
		if len(cfg.Langs) > 0 {
			//log.Debugf("Language supported: %+v", cfg.Langs)
			if langHdr := r.Header.Get(cfg.LangHeader); langHdr != "" {
				if stringExists(cfg.Langs, langHdr, "") {
					log.Debugf("Use lang %s from header", langHdr)
					(*session)["lang"] = langHdr
				} else {
					log.Debugf("Unsupported lang %s in header, using default", langHdr)
					(*session)["lang"] = cfg.Langs[0]
				}
			} else {
				(*session)["lang"] = cfg.Langs[0]
			}
		}
		// setup timezone
		if tzHdr := r.Header.Get(cfg.TZHeader); tzHdr != "" {
			log.Debugf("Use tz %s from header", tzHdr)
			(*session)["tz"] = tzHdr
		} else {
			(*session)["tz"] = ""
		}

		if r.Method == "GET" {
			s.getContextHandler(w, r, session, true, cfg.Compact)
		} else if r.Method == "HEAD" {
			s.getContextHandler(w, r, session, false, false) // Like get but without data
		} else if r.Method == "POST" && r.URL.Path == cfg.Prefix {
			s.postContextHandler(w, r, session)
		} else if r.Method == "POST" {
			s.postgrestContextHandler(w, r, session)
		} else if r.Method == "OPTIONS" {
			w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusNoContent)
		} else {
			e := fmt.Sprintf("Unsupported request method: %s", r.Method)
			log.Warn(e)
			http.Error(w, e, http.StatusNotImplemented)
		}
	}
}

// -----------------------------------------------------------------------------

func setMetric(w http.ResponseWriter, start time.Time, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Elapsed", fmt.Sprint(time.Since(start)))
	w.WriteHeader(status)
}

// -----------------------------------------------------------------------------

func getRaw(data interface{}) *json.RawMessage {
	j, _ := json.Marshal(data)
	raw := json.RawMessage(j)
	return &raw
}

// -----------------------------------------------------------------------------

// Age calculates age tag for caching
func (s RPCServer) Age(max int) int {

	fromStart := int(time.Now().Unix()) - s.started
	var ret int
	if max < 0 {
		ret = fromStart // 1sec cache
	} else if max == 0 {
		ret = 0 // single call per runtime
	} else {
		ret = fromStart / max // will change every max sec
	}
	s.log.Infof("Alive for %d, with max age %d got key %d", fromStart, max, ret)
	return ret
}

// -----------------------------------------------------------------------------

// Check if str or any exists in strings slice
func stringExists(strings []string, str string, any string) bool {
	if len(strings) > 0 { // lookup if host is allowed
		for _, s := range strings {
			if str == s || (any != "" && s == any) {
				return true
			}
		}
	}
	return false
}

// -----------------------------------------------------------------------------

// FunctionDef returns function attributes from index() method
func (s *RPCServer) FunctionDef(method string) (*FuncDef, error) {

	s.lock.RLock()
	defer s.lock.RUnlock()

	fm := *s.funcs

	if def, ok := fm[method]; ok {
		return &def, nil
	}
	return nil, fmt.Errorf("no method %s", method)
}

// -----------------------------------------------------------------------------

func (s *RPCServer) getContextHandler(w http.ResponseWriter, r *http.Request, session *jwtutil.Session, reply bool, compact bool) {
	start := time.Now()
	log := s.log
	method := strings.TrimPrefix(r.URL.Path, s.cfg.Prefix)
	method = strings.TrimSuffix(method, ".json") // Allow use .json in url

	needJWT := strings.HasSuffix(method, s.cfg.JWTSuffix)
	if needJWT {
		log.Infof("Method %s has JWT suffix which does not allowed for GET requests", method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fd, err := s.FunctionDef(method)
	if err != nil {
		// Warning was when fetched from db
		log.Infof("Method %s load def error: %s", method, err)
		http.NotFound(w, r)
		return
	}

	argDef, errd := s.FunctionArgDef(fd.NspName, fd.ProName)
	if errd != nil {
		// Warning was when fetched from db
		log.Infof("Method %s load def error: %s", method, errd)
		http.NotFound(w, r)
		return
	}

	r.ParseForm()
	var lang string
	if langs, ok := (*session)["lang"]; ok {
		lang = langs.(string)
	}
	var tz string
	if tzs, ok := (*session)["tz"]; ok {
		tz = tzs.(string)
	}

	f404 := []string{}
	ret := CallDef{
		Cache: s.CacheID(),
		Age:   s.Age(fd.MaxAge),
		Name:  &fd.NspName,
		Proc:  &fd.ProName,
		Lang:  &lang,
		TZ:    &tz,
		Args:  map[string]interface{}{},
	}

	for _, a := range argDef {
		var v []string
		if strings.HasPrefix(a.Name, s.cfg.JWTArgPrefix) {
			// session arg
			if session != nil {
				if v1, ok := (*session)[strings.TrimPrefix(a.Name, s.cfg.JWTArgPrefix)]; ok {
					v = append(v, v1.(string))
				}
			}
		} else {
			v = r.Form[a.Name]
		}
		if len(v) == 0 {
			if a.Required && a.Default == nil {
				f404 = append(f404, a.Name)
			} else if a.Default != nil {
				log.Debugf("Arg: %s use default", a.Name)
			}
		} else if strings.HasSuffix(a.Type, "[]") {
			// convert array into string
			if v[0] == "{}" {
				// empty array
				ret.Args[a.Name] = &v[0]
			} else {
				str := "{" + strings.Join(v, ",") + "}" // TODO: escape ","
				ret.Args[a.Name] = &str
			}
		} else {
			ret.Args[a.Name] = &v[0]
		}
	}
	var result workman.Result
	if len(f404) > 0 {
		result = workman.Result{Success: false, Error: fmt.Sprintf("Required parameter(s) %+v not found", f404)}
	} else {
		payload, _ := json.Marshal(ret)
		log.Debugf("Args: %s", string(payload))
		result = FunctionResult(s.jc, string(payload))
	}

	if reply {
		var out []byte
		if compact {
			out, err = json.Marshal(result)
		} else {
			out, err = json.MarshalIndent(result, "", "    ")
		}
		if err != nil {
			log.Warnf("Marshall error: %+v", err)
			e := workman.Result{Success: false, Error: err.Error()}
			out, _ = json.Marshal(e)
		}
		setMetric(w, start, http.StatusOK)
		log.Debug("Start writing")
		w.Write(out)
		//w.Write([]byte("\n"))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// -----------------------------------------------------------------------------

// postContextHandler serve JSON-RPC envelope
func (s *RPCServer) postContextHandler(w http.ResponseWriter, r *http.Request, session *jwtutil.Session) {

	start := time.Now()
	log := s.log

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
	resultStatus := http.StatusOK

	var method string
	needJWT := strings.HasSuffix(req.Method, s.cfg.JWTSuffix)
	if needJWT {
		method = strings.TrimSuffix(req.Method, s.cfg.JWTSuffix)
	} else {
		method = req.Method
	}

	fd, err := s.FunctionDef(method)
	if err != nil {
		resultRPC.Error = getRaw(respRPCError{Code: -32601, Message: "Method not found", Data: getRaw(err)})
		resultStatus = http.StatusNotFound
	} else {

		argDef, errd := s.FunctionArgDef(fd.NspName, fd.ProName)
		if errd != nil {
			log.Warnf("Method %s load def error: %s", method, errd)
			resultRPC.Error = getRaw(respRPCError{Code: -32601, Message: "Method not found", Data: getRaw(errd)})
			resultStatus = http.StatusNotFound
		} else {
			// Load args
			r.ParseForm()
			log.Infof("Argument source: %+v", req.Params)
			key, f404 := fetchArgs(log, argDef, req.Params, fd.NspName, fd.ProName, s.Age(fd.MaxAge), s.CacheID(), session, s.cfg.JWTArgPrefix)
			if len(f404) > 0 {
				resultRPC.Error = getRaw(respRPCError{Code: -32602, Message: "Required parameter(s) not found", Data: getRaw(f404)})
			} else {
				payload, _ := json.Marshal(key)
				log.Debugf("Args: %s", string(payload))
				res := FunctionResult(s.jc, string(payload))
				if res.Success {
					if needJWT {
						resultRPC.Result, _ = s.JWT.Create(method, res.Result)
					} else {
						resultRPC.Result = res.Result
					}
				} else {
					resultRPC.Error = getRaw(respRPCError{Code: -32603, Message: "Internal Error", Data: getRaw(res.Error)})
				}
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
	setMetric(w, start, resultStatus)
	log.Debug("Start writing")
	w.Write(out)
	log.Debugf("JSON Resp: %s", string(out))
	//w.Write([]byte("\n"))
}

// -----------------------------------------------------------------------------

// postgrestContextHandler serve JSON-RPC envelope
// 404 when method not found
func (s *RPCServer) postgrestContextHandler(w http.ResponseWriter, r *http.Request, session *jwtutil.Session) {

	start := time.Now()
	log := s.log

	method := strings.TrimPrefix(r.URL.Path, s.cfg.Prefix)
	method = strings.TrimSuffix(method, ".json") // Allow use .json in url
	log.Debugf("postgrest call for %s", method)

	needJWT := strings.HasSuffix(method, s.cfg.JWTSuffix)
	if needJWT {
		method = strings.TrimSuffix(method, s.cfg.JWTSuffix)
	}
	fd, err := s.FunctionDef(method)
	if err != nil {
		log.Warnf("Method %s load def error: %s", method, err)
		http.NotFound(w, r)
		return
	}

	argDef, errd := s.FunctionArgDef(fd.NspName, fd.ProName)
	if errd != nil {
		log.Warnf("Method %s load def error: %s", method, errd)
		http.NotFound(w, r)
		return
	}
	resultStatus := http.StatusOK

	req := reqParams{}
	var resultRPC interface{}

	data, _ := ioutil.ReadAll(r.Body)

	if len(data) == 0 {
		resultRPC = respPGTError{Message: "Cannot parse empty request payload, use '{}'"}
		resultStatus = http.StatusBadRequest
	} else {

		err = json.Unmarshal(data, &req)

		if err != nil {
			e := fmt.Sprintf("json parse error: %s", err)
			log.Warnf("Error parse request(%s): %+v", data, e)
			resultRPC = respPGTError{Message: "Cannot parse request payload", Details: getRaw(e)}
			resultStatus = http.StatusBadRequest
		} else {
			// Load args
			log.Infof("Argument source: %+v (session: %+v)", req, session)
			key, f404 := fetchArgs(log, argDef, req, fd.NspName, fd.ProName, s.Age(fd.MaxAge), s.CacheID(), session, s.cfg.JWTArgPrefix)
			if len(f404) > 0 {
				resultRPC = respPGTError{Code: "42883", Message: "Required parameter(s) not found", Details: getRaw(strings.Join(f404, ", "))}
				resultStatus = http.StatusBadRequest
			} else {
				payload, _ := json.Marshal(key)
				log.Debugf("Args: %s", string(payload))
				res := FunctionResult(s.jc, string(payload))
				if res.Success {
					if needJWT {
						resultRPC, _ = s.JWT.Create(method, res.Result)
					} else {
						resultRPC = res.Result
					}

				} else {
					resultRPC = respPGTError{Message: "Method call error", Details: getRaw(res.Error)}
					resultStatus = http.StatusBadRequest // TODO: ?
				}
			}
		}
	}

	out, err := json.Marshal(resultRPC)
	if err != nil {
		log.Warnf("Marshall error: %+v", err)
		resultRPC = respPGTError{Message: "Method result marshall error", Details: getRaw(err)}
		resultStatus = http.StatusBadRequest // TODO: ?
		out, _ = json.Marshal(resultRPC)
	}
	setMetric(w, start, resultStatus)
	log.Debug("Start writing")
	w.Write(out)
	log.Debugf("JSON Resp: %s", string(out))
	//w.Write([]byte("\n"))
}

func fetchArgs(log *logger.Log, argDef FuncArgDef, req reqParams,
	nsp, proc string, age int, cacheID uint32, session *jwtutil.Session, prefix string) (CallDef, []string) {

	f404 := []string{}
	var lang string
	if langs, ok := (*session)["lang"]; ok {
		lang = langs.(string)
	}
	var tz string
	if tzs, ok := (*session)["tz"]; ok {
		tz = tzs.(string)
	}
	ret := CallDef{
		Cache: cacheID,
		Age:   age,
		Name:  &nsp,
		Proc:  &proc,
		Lang:  &lang,
		TZ:    &tz,
		Args:  map[string]interface{}{},
	}
	for _, a := range argDef {

		var v interface{}
		var ok bool

		if strings.HasPrefix(a.Name, prefix) {
			// get from session
			if session != nil {
				v, ok = (*session)[strings.TrimPrefix(a.Name, prefix)]
			} else {
				ok = false
			}
		} else {
			// get from request
			v, ok = req[a.Name]
		}

		if !ok {
			if a.Required && a.Default == nil {
				f404 = append(f404, a.Name)
			} else if a.Default != nil {
				log.Debugf("Arg: %s use default", a.Name)
			}
		} else if strings.HasSuffix(a.Type, "[]") {
			// wait slice
			s := reflect.ValueOf(v)
			if s.Kind() != reflect.Slice {
				// string or {string}
				vs := v.(string)
				log.Debugf("=Array from no slice: %+v", vs)
				ret.Args[a.Name] = &vs
			} else {
				// slice
				arr := make([]string, s.Len())

				for i := 0; i < s.Len(); i++ {
					arr[i] = s.Index(i).Interface().(string)
					//	log.Printf("====== %+v", ret[i])
				}
				// convert array into string
				// TODO: escape ","
				ss := "{" + strings.Join(arr, ",") + "}"
				log.Debugf("=Array from slice: %+v", ss)
				ret.Args[a.Name] = &ss
			}
		} else if a.Type == "integer" { // may be 2.001380402e+09
			log.Debugf("=Arg (int): %+v", v)
			//	i, err := strconv.ParseInt(v.(string), 10, 64)
			var s string
			f, err := getFloat(v)
			if err != nil {
				// log.Debugf("Cannot convert to int %+v: %+v", v, err)
				s = fmt.Sprintf("%s", v.(string)) //
			} else {
				s = fmt.Sprintf("%.0f", f) //v.(string)
			}
			ret.Args[a.Name] = &s
		} else {
			log.Debugf("=Scalar from iface: %+v", v)
			ret.Args[a.Name] = v
		}

	}
	return ret, f404
}

func getFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	// ...other cases...
	default:
		return math.NaN(), fmt.Errorf("getFloat: unknown value is of incompatible type %+v", i)
	}
}

func getBool(unk interface{}) (bool, error) {
	switch unk.(type) {
	case bool:
		if unk.(bool) {
			return true, nil
		}
		return false, nil

	// ...other cases...
	default:
		b, err := strconv.ParseBool(unk.(string))
		return b, err //fmt.Errorf("getFloat: unknown value is of incompatible type %+v", i)
	}
}
