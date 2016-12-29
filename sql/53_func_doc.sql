/*
  Functions for stored proc documentation fetching

*/
/* ------------------------------------------------------------------------- */

SELECT register_comment_common('ru','{
  "a_lang":  "Язык документации"
, "a_code":  "Имя процедуры"
, "code":    "Имя процедуры"
, "nspname": "Имя схемы хранимой функции"
, "proname": "Имя хранимой функции"
, "anno":    "Описание"
, "arg":     "Имя аргумента"
, "def_val": "Значение по умолчанию"
}');

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION index(a_lang TEXT DEFAULT 'ru') RETURNS TABLE (
  code TEXT
, nspname NAME
, proname NAME
, permit_code TEXT
, max_age INTEGER
, anno    TEXT
, sample  TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  SELECT *
    FROM rpc.func_def
   ORDER BY code
$_$;

SELECT register_comment('index', 'Список описаний процедур', '{}', '{"permit_code":"Код разрешения", "max_age":"Время хранения в кэше(сек)", "sample":"Пример вызова"}','{"a_lang":"ru"}');

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION func_args(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE (
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
  SELECT f.arg, type, required, def_val, COALESCE(d.anno, c.anno)
   FROM q_def q, rpc.pg_func_args(q.n, q.p) f
   LEFT OUTER JOIN rpc.func_arg_def d ON (d.arg = f.arg AND d.lang = $2 AND d.code = $1 AND d.is_in)
   LEFT OUTER JOIN rpc.func_arg_def_common c ON (c.arg = f.arg AND c.lang = $2)
$_$;

SELECT register_comment('func_args', 'Описание аргументов процедуры','{}','{"type": "Тип аргумента", "required": "Значение обязательно"  }'
  ,'{"a_code":"func_args", "a_lang":"ru"}');

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION func_result(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE (
  arg  TEXT
, type TEXT
, anno TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  WITH q_def (n, p) AS (
    SELECT nspname, proname FROM rpc.func_def where code = $1
  )
  SELECT f.arg, type, COALESCE(d.anno, c.anno)
   FROM q_def q, rpc.pg_func_result(q.n, q.p) f
   LEFT OUTER JOIN rpc.func_arg_def d ON (d.arg = f.arg AND d.lang = $2 AND d.code = $1 AND NOT d.is_in)
   LEFT OUTER JOIN rpc.func_arg_def_common c ON (c.arg = f.arg AND c.lang = $2)
$_$;

SELECT register_comment('func_result', 'Описание результата процедуры', '{}', '{
  "type": "Тип аргумента"
  }', '{"a_code":"func_args", "a_lang":"ru"}');
