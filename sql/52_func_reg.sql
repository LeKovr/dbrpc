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
, a_max_age INT DEFAULT 0
, a_permit_code TEXT DEFAULT NULL
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
DECLARE
  v_version INTEGER;
BEGIN
  SELECT setting INTO v_version FROM pg_settings WHERE name = 'server_version_num';
  IF v_version >= 90500 THEN
    INSERT INTO
      rpc.func_def (  code,   nspname,   proname,   anno,   sample, max_age, permit_code)
      VALUES   (a_code, a_nspname, a_proname, a_anno, a_sample, a_max_age, a_permit_code)
      ON CONFLICT ON CONSTRAINT func_def_pkey DO UPDATE SET -- PG9.3: ERROR:  syntax error at or near "ON"
        nspname = a_nspname
      , proname = a_proname
      , anno = a_anno
      , sample = a_sample
      , max_age = a_max_age
      , permit_code = a_permit_code
    ;
  ELSE
    UPDATE rpc.func_def SET
        nspname = a_nspname
      , proname = a_proname
      , anno = a_anno
      , sample = a_sample
      , max_age = a_max_age
      , permit_code = a_permit_code
      WHERE code = a_code
    ;
    IF NOT FOUND THEN
      INSERT INTO rpc.func_def
        (  code,   nspname,   proname,   anno,   sample,   max_age,   permit_code) VALUES
        (a_code, a_nspname, a_proname, a_anno, a_sample, a_max_age, a_permit_code)
      ;
    END IF;
  END IF;

  DELETE FROM rpc.func_arg_def WHERE code = a_code;
  INSERT INTO rpc.func_arg_def
    (   code, is_in,   lang, arg,  anno)
    SELECT
      a_code,  true, a_lang, key, value
      FROM json_each_text(a_args)
  ;
  INSERT INTO rpc.func_arg_def
    (   code, is_in,   lang, arg, anno)
    SELECT
      a_code, false, a_lang, key, value
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
, a_max_age INT DEFAULT 0
, a_permit_code TEXT DEFAULT NULL
) RETURNS void VOLATILE LANGUAGE 'sql' AS
$_$
  SELECT rpc.register_comment('ru', $1, current_schema(), $1, $2, $3, $4, $5, $6, $7)
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION register_comment_common(
  a_lang    TEXT
, a_args    JSON
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
DECLARE
  r_rec     RECORD;
  -- v_version INTEGER;
BEGIN
  /*
  SELECT setting INTO v_version FROM pg_settings WHERE name = 'server_version_num';
  IF v_version >= 90500 THEN
    INSERT INTO func_arg_def_common (lang, arg, anno)
      SELECT a_lang, key, value
        FROM json_each_text(a_args)
    ON CONFLICT ON CONSTRAINT func_arg_def_common_pkey DO UPDATE SET
        anno = value -- PG9.5: cannot be referenced from this part of the query
    ;
  ELSE
  */
    FOR r_rec IN SELECT * FROM json_each_text(a_args) LOOP
      UPDATE rpc.func_arg_def_common SET
        anno = r_rec.value
      WHERE lang = a_lang
        AND  arg = r_rec.key
      ;
      IF NOT FOUND THEN
        INSERT INTO rpc.func_arg_def_common
          (  lang,       arg,        anno) VALUES
          (a_lang, r_rec.key, r_rec.value)
        ;
      END IF;
    END LOOP;
  -- END IF;
END;
$_$;

