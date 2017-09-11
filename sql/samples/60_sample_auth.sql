/*
TODO
 * шифрование пароля
 * ф-я проверки доступа

*/
CREATE SCHEMA IF NOT EXISTS test_auth;

SET SEARCH_PATH = 'test_auth', 'public';

CREATE TABLE IF NOT EXISTS role(
  code TEXT PRIMARY KEY
, name TEXT
);

CREATE TABLE IF NOT EXISTS permit(
  role_code TEXT REFERENCES role
, code TEXT
, anno TEXT
, CONSTRAINT permit_pkey PRIMARY KEY(role_code, code)
);

CREATE TABLE IF NOT EXISTS team_type(
  code TEXT PRIMARY KEY
, anno TEXT
);

CREATE TABLE IF NOT EXISTS team(
  id INTEGER PRIMARY KEY
, status_id INTEGER NOT NULL DEFAULT 1
, area LTREE
, inn TEXT NOT NULL UNIQUE -- TODO: pair with area country
, name TEXT
, team_type TEXT REFERENCES team_type
, team_data jsonb
, created_by INTEGER -- REFERENCES account
, created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS account(
  id INTEGER PRIMARY KEY
, role_code TEXT REFERENCES role NOT NULL
, status_id INTEGER NOT NULL DEFAULT 1
, login TEXT UNIQUE
, tele_id TEXT UNIQUE
, phone TEXT UNIQUE
, api_key TEXT
, password TEXT
, name TEXT
, team_id INTEGER REFERENCES team
, created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP
);


/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION login(
  a_login TEXT
, a_password TEXT
) RETURNS TABLE(
  id INTEGER
, role_code TEXT
, status_id INTEGER
, team_id INTEGER
, team_status_id INTEGER
, name TEXT
, team_name TEXT
) STABLE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  RETURN QUERY
    SELECT a.id
    , a.role_code
    , a.status_id
    , a.team_id
    , t.status_id AS team_status_id
    , a.name
    , t.name AS team_name
    FROM account a LEFT OUTER JOIN team t ON (t.id = a.team_id)
    WHERE login = a_login AND password = a_password
  ;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'user % is unknown or password does not match', a_login;
  END IF;
  RETURN;
END
$_$;
SELECT rpc.alias('login', 'Авторизация'
, '{"a_login":"Логин","a_password":"Пароль"}'
, '{"id":"ID пользователя", "name":"Имя пользователя"}'
,'{"a_login":"user","a_password":"see_.config"}'
);

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION api_open(
  a_id INTEGER
, a_key TEXT
) RETURNS TABLE(
  id INTEGER
, role_code TEXT
, status_id INTEGER
, team_id INTEGER
, team_status_id INTEGER
, name TEXT
, team_name TEXT
) STABLE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  -- TODO: проверить permit на открытие АПИ
  RETURN QUERY
    SELECT a.id
    , a.role_code
    , a.status_id
    , a.team_id
    , t.status_id AS team_status_id
    , a.name
    , t.name AS team_name
    FROM account a LEFT OUTER JOIN team t ON (t.id = a.team_id)
    WHERE a.id = a_id AND api_key = a_key AND api_key <> ''
  ;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'account % is unknown or api key does not match', a_id;
  END IF;
  RETURN;
END
$_$;

SELECT rpc.alias('api_open', 'Авторизация API'
, '{"a_id":"ID учетной записи","a_key":"Ключ API"}'
, '{"id":"ID пользователя", "name":"Имя пользователя"}'
,'{"a_id":"1","a_key":"demo"}'
, 5);

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION need_permit(_id INTEGER, a_permit_code TEXT) RETURNS VOID STABLE LANGUAGE 'plpgsql' AS
$_$
 -- internal func
BEGIN
  IF NOT EXISTS(SELECT 1
    FROM account a JOIN permit p USING (role_code)
    WHERE a.id = _id AND p.code = a_permit_code
  ) THEN
    RAISE EXCEPTION 'account % does not have required permission %', _id, a_permit_code;
  END IF;
END;
$_$;

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION api_tele_profile(_id INTEGER, a_id TEXT) RETURNS TABLE(
  id INTEGER
, name TEXT
, status_id INTEGER
, role_code TEXT
, team_id INTEGER
, team_status_id INTEGER
, team_name TEXT
) STABLE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  PERFORM need_permit(_id, 'as_tele_user');
  RETURN QUERY
    SELECT a.id
    , a.name
    , a.status_id
    , a.role_code
    , a.team_id
    , t.status_id AS team_status_id
    , t.name AS team_name
    FROM account a LEFT OUTER JOIN team t ON (t.id = a.team_id)
    WHERE tele_id = a_id
  ;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'account % is unknown', a_id;
  END IF;
END
$_$;

SELECT rpc.alias('api_tele_profile', 'Профиль пользователя Telegram'
, '{"_id":"ID JWT","a_id":"Ключ Telegram"}'
, '{"id":"ID пользователя", "name":"Имя пользователя"}'
,'{"_id":"1","a_id":"telegram_demo_id"}'
, 5
);

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION api_tele_profile_set(_id INTEGER, a_id TEXT, a_phone TEXT) RETURNS TABLE(
  id INTEGER
, name TEXT
, status_id INTEGER
, role_code TEXT
, team_id INTEGER
, team_status_id INTEGER
, team_name TEXT
) VOLATILE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  PERFORM need_permit(_id, 'as_tele_user');
  INSERT INTO account (login, phone, tele_id, role_code) VALUES (a_phone, a_phone, a_id, 'user');
  RETURN QUERY
    SELECT a.id
    , a.name
    , a.status_id
    , a.role_code
    , a.team_id
    , t.status_id AS team_status_id
    , t.name AS team_name
    FROM account a LEFT OUTER JOIN team t ON (t.id = a.team_id)
    WHERE tele_id = a_id
  ;
END
$_$;

SELECT rpc.alias('api_tele_profile_set', 'Создать Профиль пользователя Telegram'
, '{"_id":"ID JWT","a_id":"Ключ Telegram"}'
, '{"id":"ID пользователя", "name":"Имя пользователя"}'
,'{"_id":"1","a_id":"telegram_demo_id"}'
, 5
);

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION api_team_list(
  _id INTEGER
, a_status_id INTEGER DEFAULT 0
) RETURNS TABLE(
  id INTEGER
, status_id INTEGER
, name TEXT
, area LTREE
, inn TEXT
, created_by INTEGER
, created_at TIMESTAMP(0)
) STABLE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  PERFORM need_permit(_id, 'as_tele_user');
  RETURN QUERY SELECT
    t.id, t.status_id, t.name, t.area, t.inn, t.created_by, t.created_at
    FROM team t
    WHERE a_status_id IN (t.status_id,0)
    ORDER BY created_at
  ;
END
$_$;

SELECT rpc.alias('api_team_list', 'Список компаний по статусу'
, '{"_id":"ID JWT","a_status_id":"ID статуса"}'
, '{
    "id":"ID компании"
  , "status_id":"ID статуса"
  , "name":"Название компании"
  , "area":"Регион"
  , "inn":"Налоговый идентификатор"
  , "created_by":"ID создателя"
  , "created_at":"Время создания"
  }'
,'{"_id":"1"}'
, 5
);

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION api_team_confirm(_id INTEGER, a_id INTEGER) RETURNS BOOL VOLATILE LANGUAGE 'plpgsql' AS
$_$
BEGIN
  PERFORM need_permit(_id, 'team_confirm');
  UPDATE team SET status_id = 2 WHERE id = a_id AND status_id = 1;
  RETURN FOUND;
END
$_$;

SELECT rpc.alias('api_team_confirm', 'Подтвердить реквизиты компании'
, '{"_id":"ID JWT","a_id":"ID компании"}'
, '{"api_team_confirm":"Флаг успешного подтверждения"}'
, '{"_id":"1","a_id":"2"}'
);

/* ------------------------------------------------------------------------- */

CREATE OR REPLACE FUNCTION account_team_set(_id INTEGER, a_name TEXT, a_inn TEXT, a_area TEXT) RETURNS INTEGER VOLATILE LANGUAGE 'plpgsql' AS
$_$
DECLARE
  v_id INTEGER;
BEGIN
  PERFORM need_permit(_id, 'team_add');
  INSERT INTO team (name, area, inn, created_by) VALUES (a_name, a_area::ltree, a_inn, _id)
  RETURNING id INTO v_id;

  UPDATE account SET team_id=v_id WHERE id=_id;
  RETURN v_id;
END
$_$;

SELECT rpc.alias('account_team_set', 'Зарегистрировать компанию'
, '{"_id":"ID JWT","a_name":"Название компании", "a_inn": "Налоговый идентификатор", "a_area":"Регион"}'
, '{"account_team_set":"ID компании"}'
,'{"_id":"1","a_name":"team Demo"}'
);

/* ------------------------------------------------------------------------- */

INSERT INTO role (code, name) VALUES
  ('service', 'Сервис')
, ('user', 'Пользователь')
, ('employee', 'Сотрудник компании')
;

INSERT INTO permit (role_code, code, anno) VALUES
  ('service', 'as_tele_user', 'Выполнять действия от имени пользователя tele')
, ('service', 'team_confirm', 'Авторизовать компанию')
, ('user',    'team_add',     'Зарегистрировать компанию')

;

CREATE SEQUENCE account_id_seq;
ALTER TABLE account ALTER COLUMN id SET DEFAULT NEXTVAL('account_id_seq');

INSERT INTO account(id, name, role_code, api_key) VALUES
  (1, 'Bot', 'service', :'SERVICE_KEY')
;
INSERT INTO account(id, name, role_code, login, password) VALUES
  (2, 'User', 'user', 'user', :'USER_PASS')
;
select SETVAL('account_id_seq', 2);

CREATE SEQUENCE team_id_seq;
ALTER TABLE team ALTER COLUMN id SET DEFAULT NEXTVAL('team_id_seq');
