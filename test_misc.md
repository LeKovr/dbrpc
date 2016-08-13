# crebas_misc.sql testing results

## pg_func_args

### code=public.dbsize

#### GET

```
curl -is http://localhost:8081/api/pg_func_args?code=public.dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 433.392µs
Content-Length: 95

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
X-Elapsed: 254.079µs
Content-Length: 69

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
X-Elapsed: 734.051µs
Content-Length: 103

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

### code=public.pg_func_args

#### GET

```
curl -is http://localhost:8081/api/pg_func_args?code=public.pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 332.831µs
Content-Length: 169

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
X-Elapsed: 268.167µs
Content-Length: 143

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
X-Elapsed: 313.701µs
Content-Length: 177

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

### code=public.pg_func_args

#### GET

```
curl -is http://localhost:8081/api/pg_func_result?code=public.pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 286.935µs
Content-Length: 215

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
X-Elapsed: 884.806µs
Content-Length: 189

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
X-Elapsed: 396.522µs
Content-Length: 223

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

### code=public.dbsize

#### GET

```
curl -is http://localhost:8081/api/pg_func_result?code=public.dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 373.247µs
Content-Length: 147

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
X-Elapsed: 252.151µs
Content-Length: 121

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
X-Elapsed: 1.80097ms
Content-Length: 155

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

### name=template1

#### GET

```
curl -is http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 922.327µs
Content-Length: 84

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
X-Elapsed: 747.445µs
Content-Length: 58

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
X-Elapsed: 1.023268ms
Content-Length: 92

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

### name=template1

#### GET

```
curl -is http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 880.968µs
Content-Length: 84

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
X-Elapsed: 333.079µs
Content-Length: 58

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
X-Elapsed: 429.84µs
Content-Length: 92

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

### name=test&id=1

#### GET

```
curl -is http://localhost:8081/api/echo?name=test&id=1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 269.389µs
Content-Length: 50

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
X-Elapsed: 1.036179ms
Content-Length: 24

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
X-Elapsed: 799.25µs
Content-Length: 58

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

### name=test

#### GET

```
curl -is http://localhost:8081/api/echo?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 236.373µs
Content-Length: 50

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
X-Elapsed: 433.086µs
Content-Length: 24

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
X-Elapsed: 269.538µs
Content-Length: 58

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

### name=test

#### GET

```
curl -is http://localhost:8081/api/echo_jsonb?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 248.487µs
Content-Length: 77

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
X-Elapsed: 281.817µs
Content-Length: 51

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
X-Elapsed: 233.221µs
Content-Length: 85

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

### name=test

#### GET

```
curl -is http://localhost:8081/api/echo_single?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 240.478µs
Content-Length: 50

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
X-Elapsed: 976.235µs
Content-Length: 24

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
X-Elapsed: 734.205µs
Content-Length: 58

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

### name=test1&name=test2

#### GET

```
curl -is http://localhost:8081/api/echo_arr?name=test1&name=test2

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 616.845µs
Content-Length: 61

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
X-Elapsed: 251.438µs
Content-Length: 35

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
X-Elapsed: 253.942µs
Content-Length: 69

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

### name=test1

#### GET

```
curl -is http://localhost:8081/api/echo_arr?name=test1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 279.821µs
Content-Length: 53

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
X-Elapsed: 333.144µs
Content-Length: 27

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
X-Elapsed: 357.134µs
Content-Length: 61

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

### name=test

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
X-Elapsed: 215.393µs
Content-Length: 152

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

### test_error

#### GET

```
curl -is http://localhost:8081/api/test_error?test_error

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 214.366µs
Content-Length: 121

{
  "error": "ERROR: prepared statement \"SELECT $1 FROM table_not_exists\" does not exist (SQLSTATE 26000)",
  "success": false
}
```

#### Postgrest

```
curl -is -d {"test_error":"test_error"} -H Content-type: application/json http://localhost:8081/api/test_error

HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=UTF-8
X-Elapsed: 205.857µs
Content-Length: 137

{
  "details": "ERROR: prepared statement \"SELECT $1 FROM table_not_exists\" does not exist (SQLSTATE 26000)",
  "message": "Method call error"
}
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"test_error","params":{"test_error":"test_error"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 312.057µs
Content-Length: 178

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

