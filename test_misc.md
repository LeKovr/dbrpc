# crebas_misc.sql testing results

* [index](#index)
* [func_args](#func_args)
* [func_result](#func_result)
* [echo_not_found](#echo_not_found)
* [test_error](#test_error)

## index

#### GET

```
curl -gis http://localhost:8081/rpc/index

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 13.80067ms
Content-Length: 531

```
```json
{
  "result": [
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_args",
      "nspname": "dbrpc",
      "code": "func_args",
      "anno": "Описание аргументов процедуры"
    },
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_result",
      "nspname": "dbrpc",
      "code": "func_result",
      "anno": "Описание результата процедуры"
    },
    {
      "sample": "{\"a_lang\":\"ru\"}",
      "proname": "index",
      "nspname": "dbrpc",
      "code": "index",
      "anno": "Список описаний процедур"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {} -H Content-type: application/json http://localhost:8081/rpc/index

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 722.404µs
Content-Length: 505

```
```json
[
  {
    "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
    "proname": "func_args",
    "nspname": "dbrpc",
    "code": "func_args",
    "anno": "Описание аргументов процедуры"
  },
  {
    "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
    "proname": "func_result",
    "nspname": "dbrpc",
    "code": "func_result",
    "anno": "Описание результата процедуры"
  },
  {
    "sample": "{\"a_lang\":\"ru\"}",
    "proname": "index",
    "nspname": "dbrpc",
    "code": "index",
    "anno": "Список описаний процедур"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"index"}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 273.229µs
Content-Length: 539

```
```json
{
  "result": [
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_args",
      "nspname": "dbrpc",
      "code": "func_args",
      "anno": "Описание аргументов процедуры"
    },
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_result",
      "nspname": "dbrpc",
      "code": "func_result",
      "anno": "Описание результата процедуры"
    },
    {
      "sample": "{\"a_lang\":\"ru\"}",
      "proname": "index",
      "nspname": "dbrpc",
      "code": "index",
      "anno": "Список описаний процедур"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: lang=ru

#### GET

```
curl -gis http://localhost:8081/rpc/index?lang=ru

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 1.17035ms
Content-Length: 531

```
```json
{
  "result": [
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_args",
      "nspname": "dbrpc",
      "code": "func_args",
      "anno": "Описание аргументов процедуры"
    },
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_result",
      "nspname": "dbrpc",
      "code": "func_result",
      "anno": "Описание результата процедуры"
    },
    {
      "sample": "{\"a_lang\":\"ru\"}",
      "proname": "index",
      "nspname": "dbrpc",
      "code": "index",
      "anno": "Список описаний процедур"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"lang":"ru"} -H Content-type: application/json http://localhost:8081/rpc/index

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 534.091µs
Content-Length: 505

```
```json
[
  {
    "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
    "proname": "func_args",
    "nspname": "dbrpc",
    "code": "func_args",
    "anno": "Описание аргументов процедуры"
  },
  {
    "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
    "proname": "func_result",
    "nspname": "dbrpc",
    "code": "func_result",
    "anno": "Описание результата процедуры"
  },
  {
    "sample": "{\"a_lang\":\"ru\"}",
    "proname": "index",
    "nspname": "dbrpc",
    "code": "index",
    "anno": "Список описаний процедур"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"index","params":{"lang":"ru"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 973.244µs
Content-Length: 539

```
```json
{
  "result": [
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_args",
      "nspname": "dbrpc",
      "code": "func_args",
      "anno": "Описание аргументов процедуры"
    },
    {
      "sample": "{\"code\":\"func_args\", \"lang\":\"ru\"}",
      "proname": "func_result",
      "nspname": "dbrpc",
      "code": "func_result",
      "anno": "Описание результата процедуры"
    },
    {
      "sample": "{\"a_lang\":\"ru\"}",
      "proname": "index",
      "nspname": "dbrpc",
      "code": "index",
      "anno": "Список описаний процедур"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## func_args

### Arguments: code=index

#### GET

```
curl -gis http://localhost:8081/rpc/func_args?code=index

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 6.497184ms
Content-Length: 130

```
```json
{
  "result": [
    {
      "type": "text",
      "def_is_null": false,
      "def": "ru",
      "arg": "lang",
      "anno": "Язык документации"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"index"} -H Content-type: application/json http://localhost:8081/rpc/func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 370.841µs
Content-Length: 104

```
```json
[
  {
    "type": "text",
    "def_is_null": false,
    "def": "ru",
    "arg": "lang",
    "anno": "Язык документации"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"func_args","params":{"code":"index"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 857.852µs
Content-Length: 138

```
```json
{
  "result": [
    {
      "type": "text",
      "def_is_null": false,
      "def": "ru",
      "arg": "lang",
      "anno": "Язык документации"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: code=func_args

#### GET

```
curl -gis http://localhost:8081/rpc/func_args?code=func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 4.908635ms
Content-Length: 225

```
```json
{
  "result": [
    {
      "type": "text",
      "def_is_null": false,
      "def": null,
      "arg": "code",
      "anno": "Имя процедуры"
    },
    {
      "type": "text",
      "def_is_null": false,
      "def": "ru",
      "arg": "lang",
      "anno": "Язык документации"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"func_args"} -H Content-type: application/json http://localhost:8081/rpc/func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 384.206µs
Content-Length: 199

```
```json
[
  {
    "type": "text",
    "def_is_null": false,
    "def": null,
    "arg": "code",
    "anno": "Имя процедуры"
  },
  {
    "type": "text",
    "def_is_null": false,
    "def": "ru",
    "arg": "lang",
    "anno": "Язык документации"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"func_args","params":{"code":"func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 372.615µs
Content-Length: 233

```
```json
{
  "result": [
    {
      "type": "text",
      "def_is_null": false,
      "def": null,
      "arg": "code",
      "anno": "Имя процедуры"
    },
    {
      "type": "text",
      "def_is_null": false,
      "def": "ru",
      "arg": "lang",
      "anno": "Язык документации"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: code=func_result

#### GET

```
curl -gis http://localhost:8081/rpc/func_args?code=func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 2.694983ms
Content-Length: 225

```
```json
{
  "result": [
    {
      "type": "text",
      "def_is_null": false,
      "def": null,
      "arg": "code",
      "anno": "Имя процедуры"
    },
    {
      "type": "text",
      "def_is_null": false,
      "def": "ru",
      "arg": "lang",
      "anno": "Язык документации"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"func_result"} -H Content-type: application/json http://localhost:8081/rpc/func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 586.331µs
Content-Length: 199

```
```json
[
  {
    "type": "text",
    "def_is_null": false,
    "def": null,
    "arg": "code",
    "anno": "Имя процедуры"
  },
  {
    "type": "text",
    "def_is_null": false,
    "def": "ru",
    "arg": "lang",
    "anno": "Язык документации"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"func_args","params":{"code":"func_result"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 199.753µs
Content-Length: 233

```
```json
{
  "result": [
    {
      "type": "text",
      "def_is_null": false,
      "def": null,
      "arg": "code",
      "anno": "Имя процедуры"
    },
    {
      "type": "text",
      "def_is_null": false,
      "def": "ru",
      "arg": "lang",
      "anno": "Язык документации"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

## func_result

### Arguments: code=index

#### GET

```
curl -gis http://localhost:8081/rpc/func_result?code=index

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 10.823916ms
Content-Length: 383

```
```json
{
  "result": [
    {
      "type": "text",
      "arg": "code",
      "anno": "Имя процедуры"
    },
    {
      "type": "name",
      "arg": "nspname",
      "anno": "Имя схемы хранимой функции"
    },
    {
      "type": "name",
      "arg": "proname",
      "anno": "Имя хранимой функции"
    },
    {
      "type": "text",
      "arg": "anno",
      "anno": "Описание"
    },
    {
      "type": "text",
      "arg": "sample",
      "anno": "Пример вызова"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"index"} -H Content-type: application/json http://localhost:8081/rpc/func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 485.613µs
Content-Length: 357

```
```json
[
  {
    "type": "text",
    "arg": "code",
    "anno": "Имя процедуры"
  },
  {
    "type": "name",
    "arg": "nspname",
    "anno": "Имя схемы хранимой функции"
  },
  {
    "type": "name",
    "arg": "proname",
    "anno": "Имя хранимой функции"
  },
  {
    "type": "text",
    "arg": "anno",
    "anno": "Описание"
  },
  {
    "type": "text",
    "arg": "sample",
    "anno": "Пример вызова"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"func_result","params":{"code":"index"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 557.214µs
Content-Length: 391

```
```json
{
  "result": [
    {
      "type": "text",
      "arg": "code",
      "anno": "Имя процедуры"
    },
    {
      "type": "name",
      "arg": "nspname",
      "anno": "Имя схемы хранимой функции"
    },
    {
      "type": "name",
      "arg": "proname",
      "anno": "Имя хранимой функции"
    },
    {
      "type": "text",
      "arg": "anno",
      "anno": "Описание"
    },
    {
      "type": "text",
      "arg": "sample",
      "anno": "Пример вызова"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: code=func_args

#### GET

```
curl -gis http://localhost:8081/rpc/func_result?code=func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 3.271638ms
Content-Length: 401

```
```json
{
  "result": [
    {
      "type": "text",
      "arg": "arg",
      "anno": "Имя аргумента"
    },
    {
      "type": "text",
      "arg": "type",
      "anno": "Тип аргумента"
    },
    {
      "type": "text",
      "arg": "def",
      "anno": "Значение по умолчанию"
    },
    {
      "type": "boolean",
      "arg": "def_is_null",
      "anno": "Значение по умолчанию задано как NULL"
    },
    {
      "type": "text",
      "arg": "anno",
      "anno": "Описание"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"func_args"} -H Content-type: application/json http://localhost:8081/rpc/func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 275.308µs
Content-Length: 375

```
```json
[
  {
    "type": "text",
    "arg": "arg",
    "anno": "Имя аргумента"
  },
  {
    "type": "text",
    "arg": "type",
    "anno": "Тип аргумента"
  },
  {
    "type": "text",
    "arg": "def",
    "anno": "Значение по умолчанию"
  },
  {
    "type": "boolean",
    "arg": "def_is_null",
    "anno": "Значение по умолчанию задано как NULL"
  },
  {
    "type": "text",
    "arg": "anno",
    "anno": "Описание"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"func_result","params":{"code":"func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 935.99µs
Content-Length: 409

```
```json
{
  "result": [
    {
      "type": "text",
      "arg": "arg",
      "anno": "Имя аргумента"
    },
    {
      "type": "text",
      "arg": "type",
      "anno": "Тип аргумента"
    },
    {
      "type": "text",
      "arg": "def",
      "anno": "Значение по умолчанию"
    },
    {
      "type": "boolean",
      "arg": "def_is_null",
      "anno": "Значение по умолчанию задано как NULL"
    },
    {
      "type": "text",
      "arg": "anno",
      "anno": "Описание"
    }
  ],
  "jsonrpc": "2.0",
  "id": 1
}
```

### Arguments: code=func_result

#### GET

```
curl -gis http://localhost:8081/rpc/func_result?code=func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 7.0286ms
Content-Length: 209

```
```json
{
  "result": [
    {
      "type": "text",
      "arg": "arg",
      "anno": "Имя аргумента"
    },
    {
      "type": "text",
      "arg": "type",
      "anno": "Тип аргумента"
    },
    {
      "type": "text",
      "arg": "anno",
      "anno": "Описание"
    }
  ],
  "success": true
}
```

#### Postgrest

```
curl -is -d {"code":"func_result"} -H Content-type: application/json http://localhost:8081/rpc/func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 315.289µs
Content-Length: 183

```
```json
[
  {
    "type": "text",
    "arg": "arg",
    "anno": "Имя аргумента"
  },
  {
    "type": "text",
    "arg": "type",
    "anno": "Тип аргумента"
  },
  {
    "type": "text",
    "arg": "anno",
    "anno": "Описание"
  }
]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"func_result","params":{"code":"func_result"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
X-Elapsed: 450.824µs
Content-Length: 217

```
```json
{
  "result": [
    {
      "type": "text",
      "arg": "arg",
      "anno": "Имя аргумента"
    },
    {
      "type": "text",
      "arg": "type",
      "anno": "Тип аргумента"
    },
    {
      "type": "text",
      "arg": "anno",
      "anno": "Описание"
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
curl -gis http://localhost:8081/rpc/echo_not_found?name=test

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Content-Length: 19

404 page not found
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/rpc/echo_not_found

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Content-Length: 19

404 page not found
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_not_found","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=UTF-8
X-Elapsed: 106.785µs
Content-Length: 87

```
```json
{
  "error": {
    "data": {},
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
curl -gis http://localhost:8081/rpc/test_error

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Content-Length: 19

404 page not found
```

#### Postgrest

```
curl -is -d {} -H Content-type: application/json http://localhost:8081/rpc/test_error

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Content-Length: 19

404 page not found
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"test_error"}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/rpc/

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=UTF-8
X-Elapsed: 77.179µs
Content-Length: 87

```
```json
{
  "error": {
    "data": {},
    "message": "Method not found",
    "code": -32601
  },
  "jsonrpc": "2.0",
  "id": 1
}
```

