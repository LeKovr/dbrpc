//  darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package main

import (
	//	"fmt"
	//	"io/ioutil"
	//	"os"
	//	"path"
	//	"syscall"
	"net/http"
	"time"

	"gopkg.in/tylerb/graceful.v1"
	//	"github.com/LeKovr/endless"
	// "github.com/fvbock/endless"

	"github.com/gorilla/mux"

	"github.com/LeKovr/go-base/logger"
)

func runServer(cfg Config, log *logger.Log, rout *mux.Router) {

	log.Printf("info: Listening in %s", cfg.Addr)
	srv := &graceful.Server{
		Timeout: 5 * time.Second,
		Server:  &http.Server{Addr: cfg.Addr, Handler: rout},
		ShutdownInitiated: func() {
			log.Printf("info: Server is shutting down")
		},
	}
	srv.ListenAndServe()

	/*
		Program := path.Base(os.Args[0])

		server := endless.NewServer(cfg.Addr, rout, log)
		server.BeforeBegin = func(addr string) {
			log.Printf("Listen %s with program pid %d", addr, syscall.Getpid())
			ioutil.WriteFile(Program+".pid", []byte(fmt.Sprintf("%d", syscall.Getpid())), 0644)
		}
		inStop := false
		f := func() { inStop = true }
		server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGINT, f)
		server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGTERM, f)
		err := server.ListenAndServe()
		if err != nil && !inStop {
			log.Debug(err)
		}
	*/
}
