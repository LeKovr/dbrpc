
# Методы API

* [cache_tick](#cache_tick) - Такт кэша
* [cache_tick30](#cache_tick30) - Такт кэша 30 сек
* [cache_tick5](#cache_tick5) - Такт кэша 5 сек
* [func_args](#func_args) - Описание аргументов процедуры
* [func_result](#func_result) - Описание результата процедуры
* [index](#index) - Список описаний процедур

## cache_tick

Такт кэша

### Аргументы

Имя | Тип | По умолчанию | Обязателен | Описание
----|-----|--------------|------------|---------
a_code | text | | true | Аргумент для кэширования

### Результат

Имя | Тип | Описание
----|-----|---------
code | text | Аргумент для кэширования
seq | bigint | Номер из последовательности

### Пример вызова

```
H=http://localhost:8081/rpc
Q='{"a_code":"test1"}'
curl -gsd "$Q" -H "Content-type: application/json" $H/cache_tick | jq .
```
```json
[
  {
    "seq": 1,
    "code": "test1"
  }
]
```


## cache_tick30

Такт кэша 30 сек

### Аргументы

Имя | Тип | По умолчанию | Обязателен | Описание
----|-----|--------------|------------|---------
a_code | text | | true | Аргумент для кэширования

### Результат

Имя | Тип | Описание
----|-----|---------
code | text | Аргумент для кэширования
seq | bigint | Номер из последовательности

### Пример вызова

```
H=http://localhost:8081/rpc
Q='{"a_code":"test1"}'
curl -gsd "$Q" -H "Content-type: application/json" $H/cache_tick30 | jq .
```
```json
[
  {
    "seq": 1,
    "code": "test1"
  }
]
```


## cache_tick5

Такт кэша 5 сек

### Аргументы

Имя | Тип | По умолчанию | Обязателен | Описание
----|-----|--------------|------------|---------
a_code | text | | true | Аргумент для кэширования

### Результат

Имя | Тип | Описание
----|-----|---------
code | text | Аргумент для кэширования
seq | bigint | Номер из последовательности

### Пример вызова

```
H=http://localhost:8081/rpc
Q='{"a_code":"test1"}'
curl -gsd "$Q" -H "Content-type: application/json" $H/cache_tick5 | jq .
```
```json
[
  {
    "seq": 2,
    "code": "test1"
  }
]
```


## func_args

Описание аргументов процедуры

### Аргументы

Имя | Тип | По умолчанию | Обязателен | Описание
----|-----|--------------|------------|---------
a_code | text | | true | Имя процедуры
a_lang | text | ru | false | Язык документации

### Результат

Имя | Тип | Описание
----|-----|---------
arg | text | Имя аргумента
type | text | Тип аргумента
required | boolean | Значение обязательно
def_val | text | Значение по умолчанию
anno | text | Описание

### Пример вызова

```
H=http://localhost:8081/rpc
Q='{"a_code":"func_args", "a_lang":"ru"}'
curl -gsd "$Q" -H "Content-type: application/json" $H/func_args | jq .
```
```json
[
  {
    "type": "text",
    "required": true,
    "def_val": null,
    "arg": "a_code",
    "anno": "Имя процедуры"
  },
  {
    "type": "text",
    "required": false,
    "def_val": "ru",
    "arg": "a_lang",
    "anno": "Язык документации"
  }
]
```


## func_result

Описание результата процедуры

### Аргументы

Имя | Тип | По умолчанию | Обязателен | Описание
----|-----|--------------|------------|---------
a_code | text | | true | Имя процедуры
a_lang | text | ru | false | Язык документации

### Результат

Имя | Тип | Описание
----|-----|---------
arg | text | Имя аргумента
type | text | Тип аргумента
anno | text | Описание

### Пример вызова

```
H=http://localhost:8081/rpc
Q='{"a_code":"func_args", "a_lang":"ru"}'
curl -gsd "$Q" -H "Content-type: application/json" $H/func_result | jq .
```
```json
[
  {
    "type": "TABLE",
    "arg": null,
    "anno": null
  },
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
    "type": "boolean",
    "arg": "required",
    "anno": "Значение обязательно"
  },
  {
    "type": "text",
    "arg": "def_val",
    "anno": "Значение по умолчанию"
  },
  {
    "type": "text",
    "arg": "anno",
    "anno": "Описание"
  }
]
```


## index

Список описаний процедур

### Аргументы

Имя | Тип | По умолчанию | Обязателен | Описание
----|-----|--------------|------------|---------
a_lang | text | ru | false | Язык документации

### Результат

Имя | Тип | Описание
----|-----|---------
code | text | Имя процедуры
nspname | name | Имя схемы хранимой функции
proname | name | Имя хранимой функции
permit_code | text | Код разрешения
max_age | integer | Время хранения в кэше(сек)
anno | text | Описание
sample | text | Пример вызова

### Пример вызова

```
H=http://localhost:8081/rpc
Q='{"a_lang":"ru"}'
curl -gsd "$Q" -H "Content-type: application/json" $H/index | jq .
```
```json
[
  {
    "sample": "{\"a_code\":\"test1\"}",
    "proname": "cache_tick",
    "permit_code": null,
    "nspname": "rpc",
    "max_age": 0,
    "code": "cache_tick",
    "anno": "Такт кэша"
  },
  {
    "sample": "{\"a_code\":\"test1\"}",
    "proname": "cache_tick",
    "permit_code": null,
    "nspname": "rpc",
    "max_age": 30,
    "code": "cache_tick30",
    "anno": "Такт кэша 30 сек"
  },
  {
    "sample": "{\"a_code\":\"test1\"}",
    "proname": "cache_tick",
    "permit_code": null,
    "nspname": "rpc",
    "max_age": 5,
    "code": "cache_tick5",
    "anno": "Такт кэша 5 сек"
  },
  {
    "sample": "{\"a_code\":\"func_args\", \"a_lang\":\"ru\"}",
    "proname": "func_args",
    "permit_code": null,
    "nspname": "rpc",
    "max_age": 0,
    "code": "func_args",
    "anno": "Описание аргументов процедуры"
  },
  {
    "sample": "{\"a_code\":\"func_args\", \"a_lang\":\"ru\"}",
    "proname": "func_result",
    "permit_code": null,
    "nspname": "rpc",
    "max_age": 0,
    "code": "func_result",
    "anno": "Описание результата процедуры"
  },
  {
    "sample": "{\"a_lang\":\"ru\"}",
    "proname": "index",
    "permit_code": null,
    "nspname": "rpc",
    "max_age": 0,
    "code": "index",
    "anno": "Список описаний процедур"
  }
]
```


---

Generated by doc_gen.sh

Thu, 29 Dec 2016 14:59:07 +0300
