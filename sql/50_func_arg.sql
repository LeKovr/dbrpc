/*
  Utility functions for stored proc definition fetching

*/

-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION pg_func_arg_prefix() RETURNS TEXT IMMUTABLE LANGUAGE 'sql' AS
$_$
  SELECT ''::TEXT
$_$;
COMMENT ON FUNCTION pg_func_arg_prefix() IS 'Prefix removed from argument names in public definitions';
