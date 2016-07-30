package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"syscall"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/fvbock/endless"
	"github.com/golang/groupcache"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"

	"github.com/LeKovr/dbrpc/workman"
	"github.com/LeKovr/go-base/logger"

	_ "expvar"
	_ "net/http/pprof"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Addr       string `long:"http_addr" default:"localhost:8081"  description:"Http listen address"`
	CacheGroup string `long:"cache_group" default:"DBRPC"  description:"Cache group name"`
	CacheSize  int64  `long:"cache_size" default:"67108864"  description:"Cache size in bytes"` // 64<<20
	Version    bool   `long:"version" description:"Show version and exit"`
	Connect    string `long:"db_connect" default:"user:pass@localhost/userdb?sslmode=disable" description:"Database connect string"`
}

// AplFlags defines applied logic flags
type AplFlags struct {
	Prefix     string   `long:"url_prefix" default:"/api/"  description:"Http request prefix"`
	Schema     string   `long:"db_schema" default:"public" description:"Database functions schema"`
	ArgDefFunc string   `long:"db_argdef" default:"pg_func_args" description:"Argument definition function"`
	Hosts      []string `long:"http_origin" description:"Allowed http origin(s)"`
}

// Config defines all of application flags
type Config struct {
	Flags
	apl AplFlags
	log logger.Flags
	wm  workman.Flags
}

// -----------------------------------------------------------------------------

func main() {

	var cfg Config
	log, db, _ := setUp(&cfg)
	defer log.Close()
	defer db.Close()

	Program := path.Base(os.Args[0])
	log.Infof("%s v %s. DataBase RPC service", Program, Version)
	log.Println("Copyright (C) 2016, Alexey Kovrizhkin <ak@elfire.ru>")

	cache := groupcache.NewGroup(
		cfg.CacheGroup,
		cfg.CacheSize,
		groupcache.GetterFunc(dbFetcher(&cfg.apl, log, db)),
	)
	log.Debugf("Cache group %s with size: %d", cfg.CacheGroup, cfg.CacheSize)

	wm, err := workman.New(
		workman.WorkerFunc(cacheFetcher(log, cache)),
		workman.Config(&cfg.wm),
		workman.Logger(log),
	)
	panicIfError(err) // check Flags parse error
	wm.Run()
	defer wm.Stop()

	mux1 := mux.NewRouter()
	mux1.PathPrefix(cfg.apl.Prefix).Handler(httpHandler(&cfg.apl, log, wm.JobQueue))
	//mux1.HandleFunc(cfg.Prefix, httpHandler(&cfg.apl, log, wm.JobQueue))

	/*
	   peers := groupcache.NewHTTPPool("http://localhost:" + *port)
	   http.ListenAndServe("127.0.0.1:"+*port, http.HandlerFunc(peers.ServeHTTP))
	*/

	server := endless.NewServer(cfg.Addr, mux1)
	server.BeforeBegin = func(addr string) {
		log.Printf("Listen %s with program pid %d", addr, syscall.Getpid())
		ioutil.WriteFile(Program+".pid", []byte(fmt.Sprintf("%d", syscall.Getpid())), 0644)
	}
	inStop := false
	server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGTERM, func() { inStop = true })
	err = server.ListenAndServe()
	if err != nil && !inStop {
		log.Debug(err)
	}
	log.Println("Server stopped")
	os.Exit(0)
}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (log *logger.Log, db *sql.DB, err error) {

	p := flags.NewParser(nil, flags.Default)
	_, err = p.AddGroup("Application Options", "", cfg)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Database Options", "", &cfg.apl)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Logging Options", "", &cfg.log)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("WorkerManager Options", "", &cfg.wm)
	panicIfError(err) // check Flags parse error

	_, err = p.Parse()
	if err != nil {
		os.Exit(1) // error message written already
	}

	if cfg.Version {
		// show version & exit
		fmt.Printf("%s\n%s\n%s", Version, Build, Commit)
		os.Exit(0)
	}

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a new instance of the logger
	log, err = logger.New(logger.Dest(cfg.log.Dest), logger.Level(cfg.log.Level))
	panicIfError(err) // check Flags parse error

	// Setup database
	db, err = sql.Open("postgres", "postgres://"+cfg.Connect)
	panicIfError(err) // check Flags parse error

	return
}

// -----------------------------------------------------------------------------

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
