/*
  Functions for stored proc documentation fetching

*/
/* ------------------------------------------------------------------------- */

SELECT register_comment_common('ru','{
  "lang":    "Язык документации"
, "code":    "Имя процедуры"
, "nspname": "Имя схемы хранимой функции"
, "proname": "Имя хранимой функции"
, "anno":    "Описание"
, "arg":     "Имя аргумента"
, "def":     "Значение по умолчанию"
}');

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION index(a_lang TEXT DEFAULT 'ru') RETURNS TABLE (
  code TEXT
, nspname NAME
, proname NAME
, anno    TEXT
, sample  TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  SELECT *
    FROM func_def
   ORDER BY code
$_$;

SELECT register_comment('index', 'Список описаний процедур', '{}', '{"sample":"Пример вызова"}','{"a_lang":"ru"}');

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION func_args(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE (
  arg  TEXT
, type TEXT
, def  TEXT
, def_is_null BOOL
, anno TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  WITH q_def (n, p) AS (
    SELECT nspname, proname FROM func_def where code = $1
  )
  SELECT f.arg, type, def, def_is_null, COALESCE(d.anno, c.anno)
   FROM q_def q, pg_func_args(q.n, q.p) f
   LEFT OUTER JOIN func_arg_def d ON (d.arg = f.arg)
   LEFT OUTER JOIN func_arg_def_common c ON (c.arg = f.arg)
  WHERE -- (not needed here) f.arg IS NOT NULL AND
        COALESCE (d.code, $1) = $1
    AND COALESCE (d.lang, $2) = $2
    AND COALESCE (c.lang, $2) = $2
    AND COALESCE (is_in, TRUE)
$_$;

SELECT register_comment('func_args', 'Описание аргументов процедуры','{}','{
  "type": "Тип аргумента"
, "def_is_null": "Значение по умолчанию задано как NULL"
  }','{"code":"func_args", "lang":"ru"}');

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION func_result(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE (
  arg  TEXT
, type TEXT
, anno TEXT
) STABLE LANGUAGE 'sql' AS
$_$
  WITH q_def (n, p) AS (
    SELECT nspname, proname FROM func_def where code = $1
  )
  SELECT f.arg, type, COALESCE(d.anno, c.anno) 
   FROM q_def q, pg_func_result(q.n, q.p) f
   LEFT OUTER JOIN func_arg_def d ON (d.arg = f.arg)
   LEFT OUTER JOIN func_arg_def_common c ON (c.arg = f.arg)
  WHERE f.arg IS NOT NULL -- null in rows with result type (TABLE/SCALAR)
    AND COALESCE (d.code, $1) = $1
    AND COALESCE (d.lang, $2) = $2
    AND COALESCE (c.lang, $2) = $2
    AND NOT COALESCE (is_in, FALSE)
$_$;

SELECT register_comment('func_result', 'Описание результата процедуры', '{}', '{
  "type": "Тип аргумента"
  }', '{"code":"func_args", "lang":"ru"}');
