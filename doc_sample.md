
# Схема public. Методы API

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



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text | | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | name | 
owner | name | 
size | text | 

## echo



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text |  | 
id | integer | 5 | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | 
id | integer | 

## echo_arr



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text[] |  | 
id | integer | 5 | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | text[] | 
id | integer | 

## echo_jsonb



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text |  | 
id | integer | 5 | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | 
id | integer | 
js | jsonb | 

## echo_single



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | 

## fsm_path_ok



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
path | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | boolean | 

## fsm_transes_ok



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
transes | text |  | 
final_trans_id | integer |  | 
start_trans_id | integer | 1 | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | boolean | 

## index



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
nspname | text |  | 
lang | text | ru | 

### Результат

Имя | Тип | Описание
----|-----|---------
code | text | 
nspname | text | 
proname | text | 
anno | text | 
sample | text | 

## pg_func_arg_prefix



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | 

## pg_func_args



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
nspname | text |  | 
proname | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
id | integer | 
name | text | 
type | text | 
def | text | 
def_is_ | boolean | 

## pg_func_args_ext



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
code | text |  | 
lang | text | ru | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | 
type | text | 
def | text | 
def_is_ | boolean | 
anno | text | 

## pg_func_result



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
nspname | text |  | 
proname | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | 
type | text | 

## pg_func_result_ext



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
code | text |  | 
lang | text | ru | 

### Результат

Имя | Тип | Описание
----|-----|---------
name | text | 
type | text | 
anno | text | 

## pg_func_search_nsp



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
code | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | 

## pg_schema_oid



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
name | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | oid | 

## register_comment



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
lang | text |  | 
nspname | text |  | 
proname | text |  | 
anno | text |  | 
args | json |  | 
result | json |  | 
sample | text |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | void | 

## register_comment_common



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
lang | text |  | 
args | json |  | 

### Результат

Имя | Тип | Описание
----|-----|---------
- | void | 

## test_error



### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------

### Результат

Имя | Тип | Описание
----|-----|---------
- | text | 
