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

### Arguments: code=public.dbsize

#### GET

```
curl -is http://localhost:8081/api/pg_func_args?code=public.dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 1.04222ms
Content-Length: 95

```
```json
{
  "result": [
    {
      "AllowNull": false,
      "Default": "",
      "Type": "text",
      "Name": "name",
      "ID": 1
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"public.dbsize"} -H Content-type: application/json http://localhost:8081/api/pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 545.782µs
Content-Length: 69

```
```json
[
  {
    "AllowNull": false,
    "Default": "",
    "Type": "text",
    "Name": "name",
    "ID": 1
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_args","params":{"code":"public.dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 705.991µs
Content-Length: 103

```
```json
{
  "result": [
    {
      "AllowNull": false,
      "Default": "",
      "Type": "text",
      "Name": "name",
      "ID": 1
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: code=public.pg_func_args

#### GET

```
curl -is http://localhost:8081/api/pg_func_args?code=public.pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 561.158µs
Content-Length: 169

```
```json
{
  "result": [
    {
      "AllowNull": false,
      "Default": null,
      "Type": "text",
      "Name": "code",
      "ID": 1
    },
    {
      "AllowNull": false,
      "Default": "a_",
      "Type": "text",
      "Name": "prefix",
      "ID": 2
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"public.pg_func_args"} -H Content-type: application/json http://localhost:8081/api/pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 535.428µs
Content-Length: 143

```
```json
[
  {
    "AllowNull": false,
    "Default": null,
    "Type": "text",
    "Name": "code",
    "ID": 1
  },
  {
    "AllowNull": false,
    "Default": "a_",
    "Type": "text",
    "Name": "prefix",
    "ID": 2
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_args","params":{"code":"public.pg_func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 295.415µs
Content-Length: 177

```
```json
{
  "result": [
    {
      "AllowNull": false,
      "Default": null,
      "Type": "text",
      "Name": "code",
      "ID": 1
    },
    {
      "AllowNull": false,
      "Default": "a_",
      "Type": "text",
      "Name": "prefix",
      "ID": 2
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## pg_func_result

### Arguments: code=public.pg_func_args

#### GET

```
curl -is http://localhost:8081/api/pg_func_result?code=public.pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 240.075µs
Content-Length: 215

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
      "name": "allow_null"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"public.pg_func_args"} -H Content-type: application/json http://localhost:8081/api/pg_func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 252.72µs
Content-Length: 189

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
    "name": "allow_null"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_result","params":{"code":"public.pg_func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 227.598µs
Content-Length: 223

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
      "name": "allow_null"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: code=public.dbsize

#### GET

```
curl -is http://localhost:8081/api/pg_func_result?code=public.dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 565.746µs
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
curl -is -d {"code":"public.dbsize"} -H Content-type: application/json http://localhost:8081/api/pg_func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 192.087µs
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
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_result","params":{"code":"public.dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 708.861µs
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
curl -is http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 503.265µs
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
X-Elapsed: 193.577µs
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
X-Elapsed: 399.691µs
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
curl -is http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 471.347µs
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
X-Elapsed: 533.098µs
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
X-Elapsed: 319.041µs
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
curl -is http://localhost:8081/api/echo?name=test&id=1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 283.327µs
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
X-Elapsed: 222.146µs
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
X-Elapsed: 309.943µs
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
curl -is http://localhost:8081/api/echo?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 269.413µs
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
X-Elapsed: 266.975µs
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
X-Elapsed: 346.611µs
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
curl -is http://localhost:8081/api/echo_jsonb?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 349.757µs
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
X-Elapsed: 206.698µs
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
X-Elapsed: 814.637µs
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
curl -is http://localhost:8081/api/echo_single?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 353.786µs
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
X-Elapsed: 367.159µs
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
X-Elapsed: 485.678µs
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
curl -is http://localhost:8081/api/echo_arr?name=test1&name=test2

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 448.593µs
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
X-Elapsed: 431.2µs
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
X-Elapsed: 281.599µs
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
curl -is http://localhost:8081/api/echo_arr?name=test1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 259.507µs
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
X-Elapsed: 266.253µs
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
X-Elapsed: 271.469µs
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
curl -is http://localhost:8081/api/echo_not_found?name=test

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
X-Elapsed: 178.601µs
Content-Length: 152

```
```json
{
  "error": {
    "data": "ERROR: Function not found: public.echo_not_found (SQLSTATE P0001)",
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
curl -is http://localhost:8081/api/test_error

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 240.005µs
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
X-Elapsed: 286.673µs
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
X-Elapsed: 448.282µs
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

