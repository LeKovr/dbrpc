/*
  Tables and functions for stored proc documentation generator

*/

/* ------------------------------------------------------------------------- */

CREATE TABLE func_def(
  code    TEXT PRIMARY KEY
, nspname NAME NOT NULL
, proname NAME NOT NULL
, anno    TEXT
, sample  TEXT
);
COMMENT ON TABLE func_def IS 'Function attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE func_arg_def(
  code  TEXT REFERENCES func_def ON DELETE CASCADE
, is_in BOOL
, lang  TEXT DEFAULT 'ru'
, arg   TEXT
, anno  TEXT
, CONSTRAINT func_arg_def_pk PRIMARY KEY (code, is_in, lang, arg)
);
COMMENT ON TABLE func_arg_def IS 'Function in/out argument attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE func_arg_def_common(
  lang  TEXT NOT NULL DEFAULT 'ru'
, arg   TEXT
, anno  TEXT
, CONSTRAINT func_arg_def_common_pk PRIMARY KEY (lang, arg)
);
COMMENT ON TABLE func_arg_def_common IS 'Common function in/out argument attributes';

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION register_comment(
  a_lang    TEXT
, a_nspname TEXT
, a_proname TEXT
, a_anno    TEXT
, a_args    JSON
, a_result  JSON
, a_sample    TEXT
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  -- code DEFAULTs to proname
  DELETE FROM func_def WHERE nspname = a_nspname AND proname = a_proname;
  INSERT INTO func_def (code, nspname, proname, anno, sample)
    VALUES (a_proname, a_nspname, a_proname, a_anno, a_sample)
  ;
  INSERT INTO func_arg_def (code, is_in, lang, arg, anno)
    SELECT a_proname, true, a_lang, key, value
      FROM json_each_text(a_args) 
  ;
  INSERT INTO func_arg_def (code, is_in, lang, arg, anno)
    SELECT a_proname, false, a_lang, key, value
      FROM json_each_text(a_result)
  ;
END;
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION register_comment_common(
  a_lang    TEXT
, a_args    JSON
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  INSERT INTO func_arg_def_common (lang, arg, anno)
    SELECT a_lang, key, value
      FROM json_each_text(a_args)
  ;
END;
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION index(a_nspname TEXT DEFAULT NULL, a_lang TEXT DEFAULT 'ru') RETURNS TABLE(code TEXT, nspname TEXT, proname TEXT, anno TEXT, sample TEXT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT COALESCE(d.code, p.proname)::TEXT AS code
  , n.nspname::TEXT
  , p.proname::TEXT
  , d.anno::TEXT
  , d.sample
    FROM pg_catalog.pg_proc p
    JOIN pg_namespace n ON (n.oid = p.pronamespace)
    LEFT OUTER JOIN func_def d ON ( d.nspname = n.nspname AND d.proname = p.proname)
   WHERE n.nspname = COALESCE($1, current_schema)
   ORDER BY 1
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_args_ext(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE(name TEXT, type TEXT, def TEXT, def_is_null BOOL, anno TEXT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT name, type, def, def_is_null, d.anno 
   FROM pg_func_args(null,$1) f
   LEFT OUTER JOIN func_arg_def d ON (f.name = d.arg)
  WHERE f.name IS NOT NULL
    AND COALESCE (d.code, $1) = $1
    AND COALESCE (d.lang, $2) = $2
    AND COALESCE (is_in, TRUE)
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_result_ext(a_code TEXT, a_lang TEXT default 'ru') RETURNS TABLE(name TEXT, type TEXT, anno TEXT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT name, type, d.anno 
   FROM pg_func_result(null,$1) f
   LEFT OUTER JOIN func_arg_def d ON (f.name = d.arg)
  WHERE f.name IS NOT NULL
    AND COALESCE (d.code, $1) = $1
    AND COALESCE (d.lang, $2) = $2
    AND NOT COALESCE (is_in, FALSE)
$_$;

