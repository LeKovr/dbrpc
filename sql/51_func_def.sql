/*
  Functions for stored proc definition fetching

*/

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_args(a_nspname TEXT, a_proname TEXT) 
  RETURNS TABLE(arg TEXT, type TEXT, id INT, required BOOL, def_val TEXT) STABLE LANGUAGE 'plpgsql' AS
$_$
  -- a_code:  название функции
  DECLARE
    v_i          INTEGER;
    v_args       TEXT;
    v_defs       TEXT[];
    v_def        TEXT;
    v_arg        TEXT;
    v_type       TEXT;
    v_default    TEXT;
    v_required   BOOL;
  BEGIN
    SELECT INTO v_args
      pg_get_function_arguments(p.oid)
      FROM pg_catalog.pg_proc p
      JOIN pg_namespace n ON (n.oid = p.pronamespace)
     WHERE n.nspname = a_nspname
       AND p.proname = a_proname
    ;

    IF NOT FOUND THEN
      RAISE EXCEPTION 'Function not found: %', a_proname;
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
        RAISE EXCEPTION 'No required arg name for % arg id %', a_proname, v_i;
      END IF;

      v_required := FALSE;
      IF split_part(v_def, ' ', 4) = 'DEFAULT' THEN
        v_default := substr(v_def, strpos(v_def, ' DEFAULT ') + 9);
        v_default := regexp_replace(v_default, '::[^:]+$', '');
        IF v_default = 'NULL' THEN
          v_default := NULL;
        ELSE
          v_default := btrim(v_default, chr(39)); -- '
        END IF;
      ELSE
        v_default := NULL;
        v_required := TRUE;
      END IF;
      v_arg  := regexp_replace(split_part(v_def, ' ', 2), '^' || rpc.pg_func_arg_prefix(), '');
      v_type := split_part(v_def, ' ', 3);
      RAISE DEBUG '   column %: name=%, type=%, req=%, def=%', v_i, v_arg, v_type, v_required, v_default;
        RETURN QUERY SELECT v_arg, v_type, v_i, v_required, v_default;
    END LOOP;
    RETURN;
  END;
$_$;
COMMENT ON FUNCTION pg_func_args(TEXT, TEXT) IS 'Return function argument definition';

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION pg_func_result(a_nspname TEXT, a_proname TEXT) RETURNS TABLE(arg TEXT, type TEXT) STABLE LANGUAGE 'plpgsql' AS
$_$
  -- a_code:  название функции
  DECLARE
    v_is_set     BOOL;
    v_ret        TEXT;
    v_defs       TEXT[];
    v_i          INTEGER;
  BEGIN
    SELECT INTO v_is_set, v_ret 
      p.proretset, pg_get_function_result(p.oid)
      FROM pg_catalog.pg_proc p
      JOIN pg_namespace n ON (n.oid = p.pronamespace)
     WHERE n.nspname = a_nspname
       AND p.proname = a_proname
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
            RETURN QUERY SELECT split_part(v_defs[v_i], ' ', 1), split_part(v_defs[v_i], ' ', 2);
        END LOOP;
    ELSE
       RETURN QUERY SELECT NULL::TEXT,'SINGLE'::TEXT;
       RETURN QUERY SELECT '-'::TEXT, v_ret; -- function scalar result type
    END IF;
    RETURN;
  END;
$_$;

COMMENT ON FUNCTION pg_func_result(TEXT, TEXT) IS 'Return function result definition';
