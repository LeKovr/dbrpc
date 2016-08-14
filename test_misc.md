# crebas_misc.sql testing results

* [pg_func_args](#pg_func_args)
* [pg_func_result](#pg_func_result)
* [dbsize](#dbsize)
* [echo](#echo)
* [echo_jsonb](#echo_jsonb)
* [echo_single](#echo_single)
* [echo_arr](#echo_arr)
* [echo_not_found](#echo_not_found)
* [test_error](#test_error)

## pg_func_args

### Arguments: nspname=public&proname=dbsize

#### GET

```
curl -gis http://localhost:8081/api/pg_func_args?nspname=public&proname=dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 462.021µs
Content-Length: 93

```
```json
{
  "result": [
    {
      "type": "text",
      "name": "name",
      "id": 1,
      "def_is_null": false,
      "def": ""
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"nspname":"public","proname":"dbsize"} -H Content-type: application/json http://localhost:8081/api/pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 428.72µs
Content-Length: 67

```
```json
[
  {
    "type": "text",
    "name": "name",
    "id": 1,
    "def_is_null": false,
    "def": ""
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_args","params":{"nspname":"public","proname":"dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 395.839µs
Content-Length: 101

```
```json
{
  "result": [
    {
      "type": "text",
      "name": "name",
      "id": 1,
      "def_is_null": false,
      "def": ""
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: nspname=public&proname=pg_func_args

#### GET

```
curl -gis http://localhost:8081/api/pg_func_args?nspname=public&proname=pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 471.319µs
Content-Length: 169

```
```json
{
  "result": [
    {
      "type": "text",
      "name": "nspname",
      "id": 1,
      "def_is_null": false,
      "def": null
    },
    {
      "type": "text",
      "name": "proname",
      "id": 2,
      "def_is_null": false,
      "def": null
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"nspname":"public","proname":"pg_func_args"} -H Content-type: application/json http://localhost:8081/api/pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 497.276µs
Content-Length: 143

```
```json
[
  {
    "type": "text",
    "name": "nspname",
    "id": 1,
    "def_is_null": false,
    "def": null
  },
  {
    "type": "text",
    "name": "proname",
    "id": 2,
    "def_is_null": false,
    "def": null
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_args","params":{"nspname":"public","proname":"pg_func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 554.326µs
Content-Length: 177

```
```json
{
  "result": [
    {
      "type": "text",
      "name": "nspname",
      "id": 1,
      "def_is_null": false,
      "def": null
    },
    {
      "type": "text",
      "name": "proname",
      "id": 2,
      "def_is_null": false,
      "def": null
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## pg_func_result

### Arguments: nspname=public&proname=pg_func_args

#### GET

```
curl -gis http://localhost:8081/api/pg_func_result?nspname=public&proname=pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 295.378µs
Content-Length: 216

```
```json
{
  "result": [
    {
      "type": "TABLE",
      "name": null
    },
    {
      "type": "integer",
      "name": "id"
    },
    {
      "type": "text",
      "name": "name"
    },
    {
      "type": "text",
      "name": "type"
    },
    {
      "type": "text",
      "name": "def"
    },
    {
      "type": "boolean",
      "name": "def_is_null"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"nspname":"public","proname":"pg_func_args"} -H Content-type: application/json http://localhost:8081/api/pg_func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 290.824µs
Content-Length: 190

```
```json
[
  {
    "type": "TABLE",
    "name": null
  },
  {
    "type": "integer",
    "name": "id"
  },
  {
    "type": "text",
    "name": "name"
  },
  {
    "type": "text",
    "name": "type"
  },
  {
    "type": "text",
    "name": "def"
  },
  {
    "type": "boolean",
    "name": "def_is_null"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_result","params":{"nspname":"public","proname":"pg_func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 322.082µs
Content-Length: 224

```
```json
{
  "result": [
    {
      "type": "TABLE",
      "name": null
    },
    {
      "type": "integer",
      "name": "id"
    },
    {
      "type": "text",
      "name": "name"
    },
    {
      "type": "text",
      "name": "type"
    },
    {
      "type": "text",
      "name": "def"
    },
    {
      "type": "boolean",
      "name": "def_is_null"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: nspname=public&proname=dbsize

#### GET

```
curl -gis http://localhost:8081/api/pg_func_result?nspname=public&proname=dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 548.243µs
Content-Length: 147

```
```json
{
  "result": [
    {
      "type": "TABLE",
      "name": null
    },
    {
      "type": "name",
      "name": "name"
    },
    {
      "type": "name",
      "name": "owner"
    },
    {
      "type": "text",
      "name": "size"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"nspname":"public","proname":"dbsize"} -H Content-type: application/json http://localhost:8081/api/pg_func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 270.481µs
Content-Length: 121

```
```json
[
  {
    "type": "TABLE",
    "name": null
  },
  {
    "type": "name",
    "name": "name"
  },
  {
    "type": "name",
    "name": "owner"
  },
  {
    "type": "text",
    "name": "size"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_result","params":{"nspname":"public","proname":"dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 293.16µs
Content-Length: 155

```
```json
{
  "result": [
    {
      "type": "TABLE",
      "name": null
    },
    {
      "type": "name",
      "name": "name"
    },
    {
      "type": "name",
      "name": "owner"
    },
    {
      "type": "text",
      "name": "size"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## dbsize

### Arguments: name=template1

#### GET

```
curl -gis http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 253.23µs
Content-Length: 84

```
```json
{
  "result": [
    {
      "size": "6873 kB",
      "owner": "postgres",
      "name": "template1"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":"template1"} -H Content-type: application/json http://localhost:8081/api/dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 262.568µs
Content-Length: 58

```
```json
[
  {
    "size": "6873 kB",
    "owner": "postgres",
    "name": "template1"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"dbsize","params":{"name":"template1"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 601.997µs
Content-Length: 92

```
```json
{
  "result": [
    {
      "size": "6873 kB",
      "owner": "postgres",
      "name": "template1"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: name=template1

#### GET

```
curl -gis http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 446.059µs
Content-Length: 84

```
```json
{
  "result": [
    {
      "size": "6873 kB",
      "owner": "postgres",
      "name": "template1"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":"template1"} -H Content-type: application/json http://localhost:8081/api/dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 312.007µs
Content-Length: 58

```
```json
[
  {
    "size": "6873 kB",
    "owner": "postgres",
    "name": "template1"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"dbsize","params":{"name":"template1"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 297.964µs
Content-Length: 92

```
```json
{
  "result": [
    {
      "size": "6873 kB",
      "owner": "postgres",
      "name": "template1"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## echo

### Arguments: name=test&id=1

#### GET

```
curl -gis http://localhost:8081/api/echo?name=test&id=1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 204.396µs
Content-Length: 50

```
```json
{
  "result": [
    {
      "name": "test",
      "id": 1
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":"test","id":"1"} -H Content-type: application/json http://localhost:8081/api/echo

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 478.119µs
Content-Length: 24

```
```json
[
  {
    "name": "test",
    "id": 1
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo","params":{"name":"test","id":"1"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 621.614µs
Content-Length: 58

```
```json
{
  "result": [
    {
      "name": "test",
      "id": 1
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: name=test

#### GET

```
curl -gis http://localhost:8081/api/echo?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 267.219µs
Content-Length: 50

```
```json
{
  "result": [
    {
      "name": "test",
      "id": 5
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 613.67µs
Content-Length: 24

```
```json
[
  {
    "name": "test",
    "id": 5
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 337.037µs
Content-Length: 58

```
```json
{
  "result": [
    {
      "name": "test",
      "id": 5
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## echo_jsonb

### Arguments: name=test

#### GET

```
curl -gis http://localhost:8081/api/echo_jsonb?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 397.934µs
Content-Length: 77

```
```json
{
  "result": [
    {
      "name": "test",
      "js": {
        "b": [
          "c",
          "d"
        ],
        "a": 2
      },
      "id": 5
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo_jsonb

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 1.244392ms
Content-Length: 51

```
```json
[
  {
    "name": "test",
    "js": {
      "b": [
        "c",
        "d"
      ],
      "a": 2
    },
    "id": 5
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_jsonb","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 508.233µs
Content-Length: 85

```
```json
{
  "result": [
    {
      "name": "test",
      "js": {
        "b": [
          "c",
          "d"
        ],
        "a": 2
      },
      "id": 5
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## echo_single

### Arguments: name=test

#### GET

```
curl -gis http://localhost:8081/api/echo_single?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 290.893µs
Content-Length: 50

```
```json
{
  "result": [
    {
      "echo_single": "test"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo_single

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 269.007µs
Content-Length: 24

```
```json
[
  {
    "echo_single": "test"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_single","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 486.74µs
Content-Length: 58

```
```json
{
  "result": [
    {
      "echo_single": "test"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## echo_arr

### Arguments: name=test1&name=test2

#### GET

```
curl -gis http://localhost:8081/api/echo_arr?name=test1&name=test2

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 491.109µs
Content-Length: 61

```
```json
{
  "result": [
    {
      "name": [
        "test1",
        "test2"
      ],
      "id": 5
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":["test1","test2"]} -H Content-type: application/json http://localhost:8081/api/echo_arr

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 778.492µs
Content-Length: 35

```
```json
[
  {
    "name": [
      "test1",
      "test2"
    ],
    "id": 5
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_arr","params":{"name":["test1","test2"]}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 343.429µs
Content-Length: 69

```
```json
{
  "result": [
    {
      "name": [
        "test1",
        "test2"
      ],
      "id": 5
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: name=test1

#### GET

```
curl -gis http://localhost:8081/api/echo_arr?name=test1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 241.498µs
Content-Length: 53

```
```json
{
  "result": [
    {
      "name": [
        "test1"
      ],
      "id": 5
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"name":["test1"]} -H Content-type: application/json http://localhost:8081/api/echo_arr

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 295.1µs
Content-Length: 27

```
```json
[
  {
    "name": [
      "test1"
    ],
    "id": 5
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_arr","params":{"name":["test1"]}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 251.5µs
Content-Length: 61

```
```json
{
  "result": [
    {
      "name": [
        "test1"
      ],
      "id": 5
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## echo_not_found

### Arguments: name=test

#### GET

```
curl -gis http://localhost:8081/api/echo_not_found?name=test

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Content-Length: 19

404 page not found
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo_not_found

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Content-Length: 19

404 page not found
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_not_found","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=UTF-8
X-Elapsed: 319.755µs
Content-Length: 145

```
```json
{
  "error": {
    "data": "ERROR: Function not found: echo_not_found (SQLSTATE P0001)",
    "message": "Method not found",
    "code": -32601
  },
  "jsonrpc": "2.0",
  "id": 1
}
```

## test_error

#### GET

```
curl -gis http://localhost:8081/api/test_error

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 356.513µs
Content-Length: 121

```
```json
{
  "error": "ERROR: prepared statement \"SELECT $1 FROM table_not_exists\" does not exist (SQLSTATE 26000)",
  "success": false
}
```

#### Postgrest

```
curl -is -d {} -H Content-type: application/json http://localhost:8081/api/test_error

HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=UTF-8
X-Elapsed: 360.485µs
Content-Length: 137

```
```json
{
  "details": "ERROR: prepared statement \"SELECT $1 FROM table_not_exists\" does not exist (SQLSTATE 26000)",
  "message": "Method call error"
}
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"test_error"}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 852.477µs
Content-Length: 178

```
```json
{
  "error": {
    "data": "ERROR: prepared statement \"SELECT $1 FROM table_not_exists\" does not exist (SQLSTATE 26000)",
    "message": "Internal Error",
    "code": -32603
  },
  "jsonrpc": "2.0",
  "id": 1
}
```

