/*
  Tables and functions for stored proc documentation generator

*/

/* ------------------------------------------------------------------------- */

CREATE TABLE func_def(
  schema NAME
, code NAME
, lang TEXT NOT NULL DEFAULT 'ru'
, anno TEXT
, test TEXT
, CONSTRAINT func_def_pk PRIMARY KEY (schema, code, lang)
);
COMMENT ON TABLE func_def IS 'Function attributes';

/* ------------------------------------------------------------------------- */

CREATE TABLE func_field_def(
  schema NAME
, code NAME
, is_in BOOL NOT NULL
, lang TEXT NOT NULL DEFAULT 'ru'
, sub_code TEXT
, anno TEXT
, CONSTRAINT func_field_def_pk PRIMARY KEY (schema, code, is_in, lang, sub_code)
, CONSTRAINT func_field_def_fk FOREIGN KEY (schema, code, lang) REFERENCES func_def ON DELETE CASCADE 
);
COMMENT ON TABLE func_field_def IS 'Function in/out fields attributes';

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION register_comment(
  a_lang TEXT
, a_schema TEXT
, a_code TEXT
, a_anno TEXT
, a_args JSON
, a_result JSON
, a_test TEXT
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  DELETE FROM func_def WHERE schema = a_schema AND code = a_code;
  INSERT INTO func_def (schema, code, lang, anno, test)
    VALUES (a_schema, a_code, a_lang, a_anno, a_test)
  ;
  INSERT INTO func_field_def (schema, code, is_in, lang, sub_code, anno)
    SELECT a_schema, a_code, true, a_lang, key, value
      FROM json_each_text(a_args) 
  ;
  INSERT INTO func_field_def (schema, code, is_in, lang, sub_code, anno)
    SELECT a_schema, a_code, false, a_lang, key, value
      FROM json_each_text(a_result)
  ;
END;
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION index(a_schema TEXT, a_lang TEXT DEFAULT 'ru') RETURNS TABLE(name TEXT, anno TEXT, example TEXT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT p.proname::TEXT, d.anno::TEXT, d.test
    FROM pg_catalog.pg_proc p
    JOIN pg_namespace n ON (n.oid = p.pronamespace)
    LEFT OUTER JOIN func_def d ON ( d.schema = n.nspname AND d.code = p.proname AND d.lang = $2)
   WHERE n.nspname = $1
   ORDER BY 1
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_args_ext(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE(name TEXT, type TEXT, anno TEXT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT name, type, d.anno 
   FROM pg_func_args($1) f
   LEFT OUTER JOIN func_field_def d ON (f.name = d.sub_code)
  WHERE f.name IS NOT NULL
    AND d.schema = split_part($1, '.', 1)
    AND d.code = split_part($1, '.', 2)
    AND d.lang = $2
    AND is_in
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_result_ext(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE(name TEXT, type TEXT, anno TEXT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT name, type, d.anno 
   FROM pg_func_result($1) f
   LEFT OUTER JOIN func_field_def d ON (f.name = d.sub_code)
  WHERE f.name IS NOT NULL
    AND d.schema = split_part($1, '.', 1)
    AND d.code = split_part($1, '.', 2)
    AND d.lang = $2
    AND NOT is_in
$_$;

