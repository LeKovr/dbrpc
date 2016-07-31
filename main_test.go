package main

import (
	"bytes"
	"fmt"
	json "github.com/gorilla/rpc/v2/json2"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeKovr/go-base/logger"
)

// https://github.com/haisum/rpcexample/blob/master/examples/jrpcclient.go

// Request prepares JSON-RPC v2 request
func Request(url, method string, args interface{}) (req *http.Request, err error) {

	message, err := json.EncodeClientRequest(method, args)
	if err != nil {
		return
	}
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return
}

// Call makes request & decodes response
func Call(req *http.Request, result interface{}) (*http.Response, error) {

	cl := new(http.Client)
	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error in sending request to %s. %s", req.URL, err)
	}
	defer resp.Body.Close()

	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		return nil, fmt.Errorf("Couldn't decode response. %s", err)
	}
	return resp, nil
}

type Params map[string]string

func TestOne(t *testing.T) {

	var cfg Config
	p := makeConfig(&cfg)
	p.Parse()

	log, _ := logger.New(logger.Level("debug"), logger.Dest("err.out"))

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	// columns to be used for result
	columns := []string{"id", "name", "type", "default", "allownull"}

	q := fmt.Sprintf("select (.+) from %s.%s(.+)", cfg.apl.Schema, cfg.apl.ArgDefFunc)
	var sp *string
	mock.ExpectQuery(q).
		WithArgs("public.echo").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(1, "id", "integer", sp, false).
			AddRow(2, "name", "text", sp, false))

	columns = []string{"id", "name"}
	q = fmt.Sprintf("select (.+) from %s.%s(.+)", cfg.apl.Schema, "echo")
	mock.ExpectQuery(q).
		WithArgs("2", "test").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(2, "test"))
	mux1, wm := Handlers(&cfg, log, db)
	wm.Run()
	defer wm.Stop()

	server := httptest.NewServer(mux1)                     //Creating new server with the user handlers
	apiURL := fmt.Sprintf("%s"+cfg.apl.Prefix, server.URL) //Grab the address for the API endpoint

	//	userJson := `{"name": "test"}`

	//	reader = strings.NewReader(userJson) //Convert string to reader

	// {"ID":1,"Name":"id","Type":"integer","Default":null,"AllowNull":false},
	//{"ID":2,"Name":"name","Type":"text","Default":null,"AllowNull":false}]"

	log.Printf("Uri %s", apiURL)
	request, err := Request(apiURL, "echo", Params{"name": "test", "id": "2"}) //Create request with JSON body
	log.Printf("Send %s %+v", apiURL, request)

	var resp []serverResponse
	_, err = Call(request, &resp) // http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Call error: %s", err)
	}
	//	if res.StatusCode != 200 {
	//		t.Errorf("Success expected, but got: %d", res.StatusCode) //Uh-oh this means our test failed
	//	}
	log.Printf("Result: %+v", resp)
	server.CloseClientConnections()

	return
}

/*
   req, err := elclient.Request("http://"+addr+"/api", "App.Multiply", &Service1Request{4, 2})
   if err != nil {
       t.Error("Expected err to be nil, but got:", err)
   }

*/
