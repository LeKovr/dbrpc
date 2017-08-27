/*
  Functions for stored proc documentation fetching

  TODO: if search_path contains i18n_?? and exists i18n_??.rpc_func_?? - get anno from there
*/

-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION index() RETURNS TABLE (
  code TEXT
, nspname TEXT
, proname TEXT
, max_age INTEGER
, anno    TEXT
, sample  TEXT
, is_ro   BOOL
) STABLE LANGUAGE 'sql' AS
$_$
  SELECT
  code
, nspname::TEXT
, proname::TEXT
, max_age
, anno
, sample
  , rpc.pg_func_is_ro(nspname, proname) AS is_ro
    FROM rpc.func_def
   ORDER BY code
$_$;

SELECT alias('index'
, 'Список описаний процедур'
, '{}'
, '{
    "code":    "Имя процедуры"
  , "nspname": "Имя схемы хранимой функции"
  , "proname": "Имя хранимой функции"
  , "max_age": "Время хранения в кэше(сек)"
  , "anno":    "Описание"
  , "sample":  "Пример вызова"
  , "is_ro":   "Метод Read-only"
  }'
,'{}'
);

-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION func_args(a_code TEXT) RETURNS TABLE (
  arg  TEXT
, type TEXT
, required BOOL
, def_val  TEXT
, anno TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  WITH q_def (n, p) AS (
    SELECT nspname, proname FROM rpc.func_def where code = $1
  )
  SELECT f.arg, type, required, def_val, d.anno
   FROM q_def q, rpc.pg_func_args(q.n, q.p) f
   LEFT OUTER JOIN rpc.func_arg_def d ON (d.arg = f.arg AND d.code = $1 AND d.is_in)
$_$;

SELECT alias('func_args'
, 'Описание аргументов процедуры'
, '{"a_code":  "Имя процедуры"}'
, '{
    "arg":      "Имя аргумента"
  , "type":     "Тип аргумента"
  , "required": "Значение обязательно"
  , "def_val":  "Значение по умолчанию"
  , "anno":     "Описание"
  }'
,'{"a_code": "func_args"}'
);

-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION func_result(a_code TEXT) RETURNS TABLE (
  arg  TEXT
, type TEXT
, anno TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  WITH q_def (n, p) AS (
    SELECT nspname, proname FROM rpc.func_def where code = $1
  )
  SELECT f.arg, type, d.anno
   FROM q_def q, rpc.pg_func_result(q.n, q.p) f
   LEFT OUTER JOIN rpc.func_arg_def d ON (d.arg = f.arg AND d.code = $1 AND NOT d.is_in)
$_$;

SELECT alias('func_result'
, 'Описание результата процедуры'
, '{"a_code":  "Имя процедуры"}'
, '{
    "arg":  "Имя аргумента"
  , "type": "Тип аргумента"
  , "anno": "Описание"
  }'
, '{"a_code":"func_args"}'
);
