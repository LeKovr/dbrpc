/*
 Sample of methods which differs by cache lifetime

*/

CREATE SCHEMA IF NOT EXISTS test_cache_tick;

SET SEARCH_PATH = 'test_cache_tick', 'public';

-- seq to check cache reset
CREATE SEQUENCE IF NOT EXISTS cache_tick_seq;


CREATE OR REPLACE FUNCTION cache_tick(a_code TEXT) RETURNS TABLE(code TEXT, seq BIGINT) STABLE LANGUAGE 'sql' AS
$_$
  SELECT $1, nextval('rpc.cache_tick_seq')
$_$;

SELECT rpc.method('cache_tick', 'rpc', 'cache_tick', 'Такт кэша'
, '{"a_code":"Аргумент для кэширования"}'
, '{"code":"Аргумент для кэширования","seq":"Номер из последовательности"}'
,'{"a_code":"test1"}'
, 0
);

SELECT rpc.method('cache_tick5', 'rpc', 'cache_tick', 'Такт кэша 5 сек'
, '{"a_code":"Аргумент для кэширования"}'
, '{"code":"Аргумент для кэширования","seq":"Номер из последовательности"}'
,'{"a_code":"test1"}'
, 5
);

SELECT rpc.method('cache_tick30', 'rpc', 'cache_tick', 'Такт кэша 30 сек'
, '{"a_code":"Аргумент для кэширования"}'
, '{"code":"Аргумент для кэширования","seq":"Номер из последовательности"}'
,'{"a_code":"test1"}'
, 30
);

