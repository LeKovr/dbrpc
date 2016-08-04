/* ------------------------------------------------------------------------- */
CREATE OR REPLACE FUNCTION pg_schema_oid(a_name TEXT) RETURNS oid STABLE LANGUAGE 'sql' AS
$_$
  -- a_name: название пакета
  SELECT oid FROM pg_namespace WHERE nspname = $1
$_$;
/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_arg_anno(
  a_src     TEXT
, a_argname TEXT
) RETURNS TEXT IMMUTABLE LANGUAGE 'sql' AS
$_$
  -- a_src:     путь к функции
  -- a_argname: название аргумента
  SELECT (regexp_matches($1, E'--\\s+' || $2 || E':\\s+(.*)$', 'gm'))[1];
$_$;
/* ------------------------------------------------------------------------- */

DROP FUNCTION IF EXISTS pg_func_args(a_code TEXT, a_prefix TEXT);
CREATE OR REPLACE FUNCTION pg_func_args(a_code TEXT, a_prefix TEXT DEFAULT 'a_') RETURNS TABLE(id INT, name TEXT, type TEXT, def TEXT, allow_null BOOL) STABLE LANGUAGE 'plpgsql' AS
$_$
  -- a_code:  название функции
  DECLARE
    v_i          INTEGER;
    v_args       TEXT;
    v_defs       TEXT[];
    v_def        TEXT;
    v_name       TEXT;
    v_type       TEXT;
    v_default    TEXT;
    v_allow_null BOOL;
  BEGIN
    SELECT INTO v_args
      pg_get_function_arguments(oid)
      FROM pg_catalog.pg_proc p
        WHERE p.pronamespace = pg_schema_oid(split_part(a_code, '.', 1))
        AND p.proname = split_part(a_code, '.', 2)
    ;

    IF NOT FOUND THEN
      RAISE EXCEPTION 'Function not found: %', a_code;
    END IF;
    IF v_args = '' THEN
      -- ф-я не имеет аргументов
      RETURN;
    END IF;

    RAISE DEBUG 'args: %',v_args;

    v_defs := regexp_split_to_array(v_args, E',\\s+');
    FOR v_i IN 1 .. pg_catalog.array_upper(v_defs, 1) LOOP
      v_def := v_defs[v_i];
      RAISE DEBUG 'PARSING ARG DEF (%)', v_def;
      IF v_def !~ E'^(IN)?OUT ' THEN
        v_def := 'IN ' || v_def;
      END IF;
      IF split_part(v_def, ' ', 1) = 'OUT' THEN
        CONTINUE;
      END IF;
      IF split_part(v_def, ' ', 3) IN ('', 'DEFAULT') THEN
        -- аргумент без имени - автогенерация невозможна
        RAISE EXCEPTION 'No required arg name for % arg id %',a_code, v_i;
      END IF;

      v_allow_null := FALSE;
      IF split_part(v_def, ' ', 4) = 'DEFAULT' THEN
        v_default := substr(v_def, strpos(v_def, ' DEFAULT ') + 9);
        v_default := regexp_replace(v_default, '::[^:]+$', '');
        IF v_default = 'NULL' THEN
          v_default := NULL;
          v_allow_null := TRUE;
        ELSE
          v_default := btrim(v_default, chr(39)); -- '
        END IF;
      ELSE
        v_default := NULL;
      END IF;
      v_name := regexp_replace(split_part(v_def, ' ', 2), '^'||a_prefix, '');
      v_type := split_part(v_def, ' ', 3);
      RAISE DEBUG '   column %: name=%, type=%, def=%, null=%', v_i, v_name, v_type, v_default, v_allow_null;
        RETURN QUERY SELECT v_i, v_name, v_type, v_default, v_allow_null;
    END LOOP;
    RETURN;
  END;
$_$;
/* ------------------------------------------------------------------------- */

DROP FUNCTION IF EXISTS pg_func_result(text);
CREATE OR REPLACE FUNCTION pg_func_result(a_code TEXT) RETURNS TABLE(name TEXT, type TEXT) STABLE LANGUAGE 'plpgsql' AS
$_$
  -- a_code:  название функции
  DECLARE
    v_is_set BOOL;
    v_ret TEXT;
    v_defs       TEXT[];
    v_i INTEGER;   
  BEGIN

    SELECT INTO v_is_set, v_ret 
      proretset, pg_get_function_result(oid)
      FROM pg_catalog.pg_proc p
        WHERE p.pronamespace = pg_schema_oid(split_part(a_code, '.', 1))
        AND p.proname = split_part(a_code, '.', 2)
    ;

    IF v_ret = '' THEN
      -- ф-я не имеет результата
      RETURN;
    END IF;
    RAISE DEBUG 'result1: % (%)',v_ret,v_is_set;
    IF v_is_set THEN
            RETURN QUERY SELECT NULL::TEXT,'TABLE'::TEXT;
        v_ret := regexp_replace(v_ret,'(TABLE\()(.+)\)',E'\\2','i');
        v_defs := regexp_split_to_array(v_ret, E',\\s+');
        FOR v_i IN 1 .. pg_catalog.array_upper(v_defs, 1) LOOP
            RETURN QUERY SELECT split_part(v_defs[v_i], ' ', 1),split_part(v_defs[v_i], ' ', 2);
        END LOOP;
    ELSE
       RETURN QUERY SELECT NULL::TEXT,'SINGLE'::TEXT;
    END IF;
    RETURN;
  END;
$_$;
/* ------------------------------------------------------------------------- */



DROP FUNCTION IF EXISTS dbsize(name TEXT);
CREATE OR REPLACE FUNCTION dbsize(a_name TEXT DEFAULT '') RETURNS TABLE(
    name NAME -- f1
    , owner name /* f2 */
    , size TEXT
    )
LANGUAGE 'sql' AS $_$
-- a_name: имя БД
SELECT d.datname
,  pg_catalog.pg_get_userbyid(d.datdba)
,  CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
        THEN pg_catalog.pg_size_pretty(pg_catalog.pg_database_size(d.datname))
        ELSE 'No Access'
    END
FROM pg_catalog.pg_database d
WHERE COALESCE(a_name,'') IN (d.datname, '')
    ORDER BY
    CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
        THEN pg_catalog.pg_database_size(d.datname)
        ELSE NULL
    END DESC -- nulls first
    LIMIT 20
$_$;


CREATE OR REPLACE FUNCTION echo(
  a_name   TEXT
,  a_id     INTEGER DEFAULT 5
) RETURNS TABLE(name TEXT, id INTEGER) LANGUAGE 'sql' AS
$_$
    SELECT $1, $2;
$_$;

CREATE OR REPLACE FUNCTION echo_single(
 a_name   TEXT
) RETURNS TEXT LANGUAGE 'sql' AS
$_$
    SELECT $1;
$_$;
/* ------------------------------------------------------------------------- */
/* ------------------------------------------------------------------------- */

SET client_min_messages to debug;
SELECT * from dbsize();


/*
SELECT -- INTO v_args, v_src
      pg_get_function_arguments(oid), prosrc, pg_get_function_result(oid)
      FROM pg_catalog.pg_proc p
        WHERE --p.pronamespace = ws.pg_schema_oid(split_part('public', '.', 1))
        -- AND 
        p.proname = 'dbsize'; -- split_part(a_code, '.', 2)
*/


select * from pg_func_args('public.dbsize');        
select * from pg_func_args('public.pg_func_args');        
select * from pg_func_result('public.pg_func_args');        
select * from pg_func_result('public.dbsize');        

select * from echo('test',1);
select * from echo('test');
select echo_single('test');
