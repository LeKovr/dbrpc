package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"gopkg.in/jackc/pgx.v2"

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
	Schema     string   `long:"db_schema" default:"public" description:"Database functions schema name or comma delimited list"`
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

	mux1, wm := Handlers(&cfg, log, db)
	wm.Run()
	defer wm.Stop()

	/*
	   peers := groupcache.NewHTTPPool("http://localhost:" + *port)
	   http.ListenAndServe("127.0.0.1:"+*port, http.HandlerFunc(peers.ServeHTTP))
	*/

	runServer(cfg, log, mux1)

	log.Println("Server stopped")
	os.Exit(0)
}

// -----------------------------------------------------------------------------

// Handlers used to prepare and http handlers
func Handlers(cfg *Config, log *logger.Log, db *pgx.Conn) (*mux.Router, *workman.WorkMan) {

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
	panicIfError(err)

	r := mux.NewRouter()
	r.PathPrefix(cfg.apl.Prefix).Handler(httpHandler(&cfg.apl, log, wm.JobQueue))

	return r, wm
}

// -----------------------------------------------------------------------------

func makeConfig(cfg *Config) *flags.Parser {
	p := flags.NewParser(nil, flags.Default)
	_, err := p.AddGroup("Application Options", "", cfg)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Applied logic Options", "", &cfg.apl)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Logging Options", "", &cfg.log)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("WorkerManager Options", "", &cfg.wm)
	panicIfError(err) // check Flags parse error
	return p
}

func setUp(cfg *Config) (log *logger.Log, db *pgx.Conn, err error) {

	p := makeConfig(cfg)

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
	c, err := pgx.ParseURI("postgres://" + cfg.Connect)
	panicIfError(err) // check Flags parse error
	db, err = pgx.Connect(c)
	panicIfError(err) // check Flags parse error

	if cfg.apl.Schema != "public" {
		_, err = db.Exec("set search_path = " + cfg.apl.Schema + ", public")
		panicIfError(err)
	}
	return
}

// -----------------------------------------------------------------------------

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
