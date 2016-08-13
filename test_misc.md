# crebas_misc.sql testing results
## pg_func_args

### code=public.dbsize

#### GET

```
curl -is http://localhost:8081/api/pg_func_args?code=public.dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 96

{"success":true,"result":[{"ID":1,"Name":"name","Type":"text","Default":"","AllowNull":false}]}
```

#### Postgrest

```
curl -is -d {"code":"public.dbsize"} -H Content-type: application/json http://localhost:8081/api/pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 70

[{"ID":1,"Name":"name","Type":"text","Default":"","AllowNull":false}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_args","params":{"code":"public.dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 104

{"id":1,"jsonrpc":"2.0","result":[{"ID":1,"Name":"name","Type":"text","Default":"","AllowNull":false}]}
```

### code=public.pg_func_args

#### GET

```
curl -is http://localhost:8081/api/pg_func_args?code=public.pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 170

{"success":true,"result":[{"ID":1,"Name":"code","Type":"text","Default":null,"AllowNull":false},{"ID":2,"Name":"prefix","Type":"text","Default":"a_","AllowNull":false}]}
```

#### Postgrest

```
curl -is -d {"code":"public.pg_func_args"} -H Content-type: application/json http://localhost:8081/api/pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 144

[{"ID":1,"Name":"code","Type":"text","Default":null,"AllowNull":false},{"ID":2,"Name":"prefix","Type":"text","Default":"a_","AllowNull":false}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_args","params":{"code":"public.pg_func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 178

{"id":1,"jsonrpc":"2.0","result":[{"ID":1,"Name":"code","Type":"text","Default":null,"AllowNull":false},{"ID":2,"Name":"prefix","Type":"text","Default":"a_","AllowNull":false}]}
```

## pg_func_result

### code=public.pg_func_args

#### GET

```
curl -is http://localhost:8081/api/pg_func_result?code=public.pg_func_args

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 216

{"success":true,"result":[{"name":null,"type":"TABLE"},{"name":"id","type":"integer"},{"name":"name","type":"text"},{"name":"type","type":"text"},{"name":"def","type":"text"},{"name":"allow_null","type":"boolean"}]}
```

#### Postgrest

```
curl -is -d {"code":"public.pg_func_args"} -H Content-type: application/json http://localhost:8081/api/pg_func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 190

[{"name":null,"type":"TABLE"},{"name":"id","type":"integer"},{"name":"name","type":"text"},{"name":"type","type":"text"},{"name":"def","type":"text"},{"name":"allow_null","type":"boolean"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_result","params":{"code":"public.pg_func_args"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 224

{"id":1,"jsonrpc":"2.0","result":[{"name":null,"type":"TABLE"},{"name":"id","type":"integer"},{"name":"name","type":"text"},{"name":"type","type":"text"},{"name":"def","type":"text"},{"name":"allow_null","type":"boolean"}]}
```

### code=public.dbsize

#### GET

```
curl -is http://localhost:8081/api/pg_func_result?code=public.dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 148

{"success":true,"result":[{"name":null,"type":"TABLE"},{"name":"name","type":"name"},{"name":"owner","type":"name"},{"name":"size","type":"text"}]}
```

#### Postgrest

```
curl -is -d {"code":"public.dbsize"} -H Content-type: application/json http://localhost:8081/api/pg_func_result

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 122

[{"name":null,"type":"TABLE"},{"name":"name","type":"name"},{"name":"owner","type":"name"},{"name":"size","type":"text"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"pg_func_result","params":{"code":"public.dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 156

{"id":1,"jsonrpc":"2.0","result":[{"name":null,"type":"TABLE"},{"name":"name","type":"name"},{"name":"owner","type":"name"},{"name":"size","type":"text"}]}
```

## dbsize

### dbsize

#### GET

```
curl -is http://localhost:8081/api/dbsize?dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 554

{"success":true,"result":[{"name":"tpro","owner":"tpro","size":"2147 MB"},{"name":"iac","owner":"iac","size":"761 MB"},{"name":"gogs","owner":"gogs","size":"8409 kB"},{"name":"mmost","owner":"mmost","size":"8001 kB"},{"name":"op","owner":"op","size":"7408 kB"},{"name":"tpro-template","owner":"postgres","size":"7121 kB"},{"name":"iac-template","owner":"postgres","size":"7065 kB"},{"name":"postgres","owner":"postgres","size":"7000 kB"},{"name":"template1","owner":"postgres","size":"6873 kB"},{"name":"template0","owner":"postgres","size":"6857 kB"}]}
```

#### Postgrest

```
curl -is -d {"dbsize":"dbsize"} -H Content-type: application/json http://localhost:8081/api/dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 528

[{"name":"tpro","owner":"tpro","size":"2147 MB"},{"name":"iac","owner":"iac","size":"761 MB"},{"name":"gogs","owner":"gogs","size":"8409 kB"},{"name":"mmost","owner":"mmost","size":"8001 kB"},{"name":"op","owner":"op","size":"7408 kB"},{"name":"tpro-template","owner":"postgres","size":"7121 kB"},{"name":"iac-template","owner":"postgres","size":"7065 kB"},{"name":"postgres","owner":"postgres","size":"7000 kB"},{"name":"template1","owner":"postgres","size":"6873 kB"},{"name":"template0","owner":"postgres","size":"6857 kB"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"dbsize","params":{"dbsize":"dbsize"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 562

{"id":1,"jsonrpc":"2.0","result":[{"name":"tpro","owner":"tpro","size":"2147 MB"},{"name":"iac","owner":"iac","size":"761 MB"},{"name":"gogs","owner":"gogs","size":"8409 kB"},{"name":"mmost","owner":"mmost","size":"8001 kB"},{"name":"op","owner":"op","size":"7408 kB"},{"name":"tpro-template","owner":"postgres","size":"7121 kB"},{"name":"iac-template","owner":"postgres","size":"7065 kB"},{"name":"postgres","owner":"postgres","size":"7000 kB"},{"name":"template1","owner":"postgres","size":"6873 kB"},{"name":"template0","owner":"postgres","size":"6857 kB"}]}
```

### name=template1

#### GET

```
curl -is http://localhost:8081/api/dbsize?name=template1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 85

{"success":true,"result":[{"name":"template1","owner":"postgres","size":"6873 kB"}]}
```

#### Postgrest

```
curl -is -d {"name":"template1"} -H Content-type: application/json http://localhost:8081/api/dbsize

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 59

[{"name":"template1","owner":"postgres","size":"6873 kB"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"dbsize","params":{"name":"template1"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 93

{"id":1,"jsonrpc":"2.0","result":[{"name":"template1","owner":"postgres","size":"6873 kB"}]}
```

## echo

### name=test&id=1

#### GET

```
curl -is http://localhost:8081/api/echo?name=test&id=1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 51

{"success":true,"result":[{"id":1,"name":"test"}]}
```

#### Postgrest

```
curl -is -d {"name":"test","id":"1"} -H Content-type: application/json http://localhost:8081/api/echo

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 25

[{"id":1,"name":"test"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo","params":{"name":"test","id":"1"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 59

{"id":1,"jsonrpc":"2.0","result":[{"id":1,"name":"test"}]}
```

### name=test

#### GET

```
curl -is http://localhost:8081/api/echo?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 51

{"success":true,"result":[{"id":5,"name":"test"}]}
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 25

[{"id":5,"name":"test"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 59

{"id":1,"jsonrpc":"2.0","result":[{"id":5,"name":"test"}]}
```

## echo_jsonb

### name=test

#### GET

```
curl -is http://localhost:8081/api/echo_jsonb?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 78

{"success":true,"result":[{"id":5,"js":{"a":2,"b":["c","d"]},"name":"test"}]}
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo_jsonb

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 52

[{"id":5,"js":{"a":2,"b":["c","d"]},"name":"test"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_jsonb","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 86

{"id":1,"jsonrpc":"2.0","result":[{"id":5,"js":{"a":2,"b":["c","d"]},"name":"test"}]}
```

## echo_single

### name=test

#### GET

```
curl -is http://localhost:8081/api/echo_single?name=test

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 51

{"success":true,"result":[{"echo_single":"test"}]}
```

#### Postgrest

```
curl -is -d {"name":"test"} -H Content-type: application/json http://localhost:8081/api/echo_single

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 25

[{"echo_single":"test"}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_single","params":{"name":"test"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 59

{"id":1,"jsonrpc":"2.0","result":[{"echo_single":"test"}]}
```

## echo_arr

### name=test1&name=test2

#### GET

```
curl -is http://localhost:8081/api/echo_arr?name=test1&name=test2

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 62

{"success":true,"result":[{"id":5,"name":["test1","test2"]}]}
```

#### Postgrest

```
curl -is -d {"name":["test1","test2"]} -H Content-type: application/json http://localhost:8081/api/echo_arr

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 36

[{"id":5,"name":["test1","test2"]}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_arr","params":{"name":["test1","test2"]}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 70

{"id":1,"jsonrpc":"2.0","result":[{"id":5,"name":["test1","test2"]}]}
```

### name=test1

#### GET

```
curl -is http://localhost:8081/api/echo_arr?name=test1

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 54

{"success":true,"result":[{"id":5,"name":["test1"]}]}
```

#### Postgrest

```
curl -is -d {"name":["test1"]} -H Content-type: application/json http://localhost:8081/api/echo_arr

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 28

[{"id":5,"name":["test1"]}]
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"echo_arr","params":{"name":["test1"]}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 62

{"id":1,"jsonrpc":"2.0","result":[{"id":5,"name":["test1"]}]}
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

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 131

{"id":1,"jsonrpc":"2.0","error":{"code":-32601,"message":"\"ERROR: Function not found: public.echo_not_found (SQLSTATE P0001)\""}}
```

## test_error

### test_error

#### GET

```
curl -is http://localhost:8081/api/test_error?test_error

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 130

{"success":false,"error":"\"ERROR: prepared statement \\\"SELECT $1 FROM table_not_exists\\\" does not exist (SQLSTATE 26000)\""}
```

#### Postgrest

```
curl -is -d {"test_error":"test_error"} -H Content-type: application/json http://localhost:8081/api/test_error

HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=UTF-8
Content-Length: 146

{"message":"Method call error","details":"\"ERROR: prepared statement \\\"SELECT $1 FROM table_not_exists\\\" does not exist (SQLSTATE 26000)\""}
```

#### JSON-RPC 2.0

```
D='{"jsonrpc":"2.0","id":1,"method":"test_error","params":{"test_error":"test_error"}}'
curl -is -d "$D" -H "Content-type: application/json" http://localhost:8081/api/

HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 187

{"id":1,"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal Error","data":"\"ERROR: prepared statement \\\"SELECT $1 FROM table_not_exists\\\" does not exist (SQLSTATE 26000)\""}}
```

