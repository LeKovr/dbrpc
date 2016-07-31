
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

* gets http request
    ```
    http://hostname/api/function?arg1=val1&arg2=val2
    ```
* loads database function signature like
    ```
    CREATE FUNCTION function(arg1 TYPE, arg2 TYPE) RETURNS SETOF ...
    ```
* calls sql
    ```
    select * from function(val1, val2)
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
* [ ] RPC interface
* [ ] Cache control
* [ ] Authentication
* [ ] Access control
* [ ] Metrics for [Prometheus](https://prometheus.io/)
* [ ] Integrated templates
* [ ] i18n

### ToDo

* endless uses syscall.Kill which is not portable to Windows yet.
* improve tests
* add `--index` arg - proc name to fetch functions list (and name -> function mapping)

Install
-------

```
go get github.com/LeKovr/dbrpc
```

### Download

See [Latest release](https://github.com/LeKovr/dbrpc/latest)

Acknowledgements
----------------
* [Marcio Castilho](http://marcio.io) for his [blog post](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
* [groupcache](https://github.com/golang/groupcache) authors
* Julio Capote for his [groupcache sample](https://github.com/capotej/groupcache-db-experiment)

License
-------

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2016 Alexey Kovrizhkin ak@elfire.ru
