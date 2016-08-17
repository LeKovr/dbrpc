/*
  Tables and functions for stored proc doc registering

*/

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION register_comment(
  a_lang    TEXT
, a_code    TEXT
, a_nspname TEXT
, a_proname TEXT
, a_anno    TEXT
, a_args    JSON
, a_result  JSON
, a_sample  TEXT
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  INSERT INTO 
    func_def (  code,   nspname,   proname,   anno,   sample)
    VALUES   (a_code, a_nspname, a_proname, a_anno, a_sample)
    ON CONFLICT ON CONSTRAINT func_def_pkey DO UPDATE SET 
      nspname = a_nspname 
    , proname = a_proname
    , anno = a_anno
    , sample = a_sample
  ;

  DELETE FROM func_arg_def WHERE code = a_code;
  INSERT INTO func_arg_def (code, is_in, lang, arg, anno)
    SELECT a_code, true, a_lang, key, value
      FROM json_each_text(a_args) 
  ;
  INSERT INTO func_arg_def (code, is_in, lang, arg, anno)
    SELECT a_code, false, a_lang, key, value
      FROM json_each_text(a_result)
  ;
END;
$_$;

/* ------------------------------------------------------------------------- */


CREATE OR REPLACE FUNCTION register_comment(
  a_proname TEXT
, a_anno    TEXT
, a_args    JSON
, a_result  JSON
, a_sample  TEXT
) RETURNS void VOLATILE LANGUAGE 'sql' AS
$_$
  SELECT register_comment('ru', $1, current_schema(), $1, $2, $3, $4, $5)
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
  ON CONFLICT ON CONSTRAINT func_arg_def_common_pkey DO NOTHING
  ;
END;
$_$;

