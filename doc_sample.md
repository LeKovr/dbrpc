
# Проект IAC. Методы API

**Внимание:** это предварительная версия программно сгенеренного документа

* [dbsize](#dbsize)
* [echo](#echo)
* [echo_arr](#echo_arr)
* [echo_jsonb](#echo_jsonb)
* [echo_single](#echo_single)
* [fsm_path_ok](#fsm_path_ok)
* [fsm_transes_ok](#fsm_transes_ok)
* [index](#index)
* [pg_func_arg_prefix](#pg_func_arg_prefix)
* [pg_func_args](#pg_func_args)
* [pg_func_args_ext](#pg_func_args_ext)
* [pg_func_result](#pg_func_result)
* [pg_func_result_ext](#pg_func_result_ext)
* [pg_func_search_nsp](#pg_func_search_nsp)
* [pg_schema_oid](#pg_schema_oid)
* [register_comment](#register_comment)
* [register_comment_common](#register_comment_common)
* [test_error](#test_error)

## dbsize

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text | | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | name | null
owner | name | null
size | text | null

## echo

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text | null | null
id | integer | 5 | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | null
id | integer | null

## echo_arr

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text[] | null | null
id | integer | 5 | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | text[] | null
id | integer | null

## echo_jsonb

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text | null | null
id | integer | 5 | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | null
id | integer | null
js | jsonb | null

## echo_single

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | null

## fsm_path_ok

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
path | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | boolean | null

## fsm_transes_ok

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
transes | text | null | null
final_trans_id | integer | null | null
start_trans_id | integer | 1 | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | boolean | null

## index

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
nspname | text | null | null
lang | text | ru | null

### Результат

Имя | Тип | Описание
----|-----|---------
code | text | null
nspname | text | null
proname | text | null
anno | text | null
sample | text | null

## pg_func_arg_prefix

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | null

## pg_func_args

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
nspname | text | null | null
proname | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
id | integer | null
name | text | null
type | text | null
def | text | null
def_is_null | boolean | null

## pg_func_args_ext

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
code | text | null | null
lang | text | ru | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | null
type | text | null
def | text | null
def_is_null | boolean | null
anno | text | null

## pg_func_result

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
nspname | text | null | null
proname | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | null
type | text | null

## pg_func_result_ext

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
code | text | null | null
lang | text | ru | null

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | null
type | text | null
anno | text | null

## pg_func_search_nsp

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
code | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | null

## pg_schema_oid

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | oid | null

## register_comment

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
lang | text | null | null
nspname | text | null | null
proname | text | null | null
anno | text | null | null
args | json | null | null
result | json | null | null
sample | text | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | void | null

## register_comment_common

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
lang | text | null | null
args | json | null | null

### Результат

Имя | Тип | Описание
----|-----|---------
- | void | null

## test_error

null

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | null
