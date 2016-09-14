

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

CREATE OR REPLACE FUNCTION echo_jsonb(
  a_name   TEXT
,  a_id     INTEGER DEFAULT 5
) RETURNS TABLE(name TEXT, id INTEGER, js jsonb) LANGUAGE 'sql' AS
$_$
    SELECT $1, $2, ' {"a": 2, "b": ["c", "d"]}'::jsonb;
$_$;

SELECT ws.register_comment('echo_arr','тест массива','{"a_name":"массив","a_id":"число"}','{"name:"массив","id":"число"}');

CREATE OR REPLACE FUNCTION echo_arr(
  a_name   TEXT[]
,  a_id     INTEGER DEFAULT 5
) RETURNS TABLE(name TEXT[], id INTEGER) LANGUAGE 'sql' AS
$_$
    SELECT $1, $2;
$_$;

CREATE OR REPLACE FUNCTION echo_single(
 a_name   TEXT
) RETURNS TEXT LANGUAGE 'sql' AS
$_$
    SELECT $1;
$_$;

CREATE OR REPLACE FUNCTION test_error() RETURNS TEXT LANGUAGE 'sql' AS
$_$EXECUTE "SELECT $1 FROM table_not_exists"; SELECT 'xx'::TEXT$_$;

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
select * from echo_jsonb('test');
select echo_single('test');

select * from echo_arr('{test1,test2}');
