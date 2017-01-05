
dbrpc
=====

[![GoCard][1]][2]
[![GitHub license][3]][4]

[1]: https://goreportcard.com/badge/LeKovr/dbrpc
[2]: https://goreportcard.com/report/github.com/LeKovr/dbrpc
[3]: https://img.shields.io/badge/license-MIT-blue.svg
[4]: LICENSE

[dbrpc](https://github.com/LeKovr/dbrpc) - Database RPC service written in go language.

This service

* gets http request (GET or POST)
    ```
    http://hostname/api/function?arg1=val1&arg2=val2
    ```
* loads database function signature like
    ```
    CREATE FUNCTION function(arg1 TYPE, arg2 TYPE) RETURNS SETOF ...
    ```
* calls sql
    ```
    select * from function(arg1 := val1, arg2 := val2)
    ```
* and returns query result as json:
    ```
curl "http://localhost:8081/api/public.echo?id=1&name=op" | jq "."
{
  "result": [
    {
      "name": "op",
      "id": 1
    }
  ],
  "success": true
}
    ```

Also, the same functionality may be used via JSON-RPC interface

Features
--------

* [x] configurable limit of simultaneous database connections
* [x] caching with [groupcache](github.com/golang/groupcache)
* [x] gracefull restart
* [x] CORS support
* [x] JSON-RPC over HTTP interface
* [x] required args checking
* [x] method index via /rpc/index[.json]
* [x] [named notation](https://www.postgresql.org/docs/devel/static/sql-syntax-calling-funcs.html)
* [x] cache expiration via max_age function attribute
* [x] JWT result encoding for configured functions
* [x] JWT header validation & func args substitution
* [ ] Authentication
* [ ] Access control
* [ ] RPC interface (gRPC?)
* [ ] Cache warm/bench/test with wget
* [ ] Reset metadata cache on SIGHUP and via LISTEN
* [ ] Metrics for [Prometheus](https://prometheus.io/) via expvar
* [ ] Integrated templates
* [ ] Swagger & human autodoc
* [ ] i18n

### ToDo

* endless uses syscall.Kill which is not portable to Windows yet.
* improve tests
* add `--index` arg - proc name to fetch functions list (and name -> function mapping)
* delay index load (via listen)
* light version - without index table
* ReadOnly function attr (for RO transactions & different db connection)/ Method with RW cached <1sec
* Avoid escaping (\u003cbr\u003e)
* add cron_func and cron_interval for this: `for q := range "select * from cron_func(stamp)" { select * from q(stamp) }`
* check if index query closed correctly
* fatal if no index data

### Arrays

Declaration:
```sql
SELECT ws.register_comment('echo_arr','тест массива','{"a_name":"массив","a_id":"число"}','{"name":"массив","id":"число"}','');

CREATE OR REPLACE FUNCTION echo_arr(
  a_name   TEXT[]
,  a_id     INTEGER DEFAULT 5
) RETURNS TABLE(name TEXT[], id INTEGER) LANGUAGE 'sql' AS
$_$
    SELECT $1, $2;
$_$;
```

Calls:
```
 curl -s 'http://localhost:8081/rpc/echo_arr?a_id=107050&a_name=2&a_name=3'
{
    "success": true,
    "result": [
        {
            "id": 107050,
            "name": [
                "2",
                "3"
            ]
        }
    ]
}

curl gs 'http://localhost:8081/rpc/echo_arr?a_id=107050&a_name=2,3'
{
    "success": true,
    "result": [
        {
            "id": 107050,
            "name": [
                "2",
                "3"
            ]
        }
    ]
}
```

Install
-------

```
go get github.com/LeKovr/dbrpc
```

### Download

See [Latest release](https://github.com/LeKovr/dbrpc/releases/latest)

Acknowledgements
----------------
* [Marcio Castilho](http://marcio.io) for his [blog post](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
* [groupcache](https://github.com/golang/groupcache) authors
* Julio Capote for his [groupcache sample](https://github.com/capotej/groupcache-db-experiment)

License
-------

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2016 Alexey Kovrizhkin ak@elfire.ru
