/*
  Tables and functions for stored proc doc registering

*/

-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION method(
  a_code    TEXT
, a_nspname TEXT
, a_proname TEXT
, a_anno    TEXT
, a_args    JSON
, a_result  JSON
, a_sample  TEXT
, a_max_age INT DEFAULT 0
) RETURNS void VOLATILE LANGUAGE 'plpgsql' AS
$_$
DECLARE
  v_version INTEGER;
BEGIN
  SELECT setting INTO v_version FROM pg_settings WHERE name = 'server_version_num';
  IF v_version >= 90500 THEN
    INSERT INTO
      rpc.func_def (  code,   nspname,   proname,   anno,   sample, max_age)
      VALUES   (a_code, a_nspname, a_proname, a_anno, a_sample, a_max_age)
      ON CONFLICT ON CONSTRAINT func_def_pkey DO UPDATE SET -- PG9.3: ERROR:  syntax error at or near "ON"
        nspname = a_nspname
      , proname = a_proname
      , anno = a_anno
      , sample = a_sample
      , max_age = a_max_age
    ;
  ELSE
    UPDATE rpc.func_def SET
        nspname = a_nspname
      , proname = a_proname
      , anno = a_anno
      , sample = a_sample
      , max_age = a_max_age
      WHERE code = a_code
    ;
    IF NOT FOUND THEN
      INSERT INTO rpc.func_def
        (  code,   nspname,   proname,   anno,   sample,   max_age) VALUES
        (a_code, a_nspname, a_proname, a_anno, a_sample, a_max_age)
      ;
    END IF;
  END IF;

  DELETE FROM rpc.func_arg_def WHERE code = a_code;
  INSERT INTO rpc.func_arg_def
    (   code, is_in, arg, anno)
    SELECT
      a_code,  true, key, value
      FROM json_each_text(a_args)
  ;
  INSERT INTO rpc.func_arg_def
    (   code, is_in, arg, anno)
    SELECT
      a_code, false, key, value
      FROM json_each_text(a_result)
  ;
END;
$_$;
COMMENT ON FUNCTION method(TEXT, TEXT, TEXT, TEXT, JSON, JSON, TEXT, INTEGER) IS 'Register RPC method';

-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION alias(
  a_proname TEXT
, a_anno    TEXT
, a_args    JSON
, a_result  JSON
, a_sample  TEXT
, a_max_age INT DEFAULT 0
) RETURNS void VOLATILE LANGUAGE 'sql' AS
$_$
  SELECT rpc.method($1, current_schema(), $1, $2, $3, $4, $5, $6)
$_$;
COMMENT ON FUNCTION alias(TEXT, TEXT, JSON, JSON, TEXT, INTEGER) IS 'Register RPC method with the same name as internal func';
