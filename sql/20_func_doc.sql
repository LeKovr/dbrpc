/*
  Tables for stored proc documenting

*/

/* ------------------------------------------------------------------------- */

CREATE TABLE func_permit(
  code    TEXT PRIMARY KEY
, is_managed BOOL NOT NULL DEFAULT FALSE -- may be added by user
-- TODO: lang support
, anno    TEXT NOT NULL
);
COMMENT ON TABLE func_permit IS 'Function permits';

-- -----------------------------------------------------------------------------

CREATE TABLE func_def(
  code    TEXT PRIMARY KEY
, nspname NAME NOT NULL
, proname NAME NOT NULL
, permit_code TEXT REFERENCES func_permit -- NULL for public
, max_age INTEGER NOT NULL DEFAULT 0    -- 0: forever, >0 secs, -1: no cache
-- TODO: lang support
, anno    TEXT NOT NULL
, sample  TEXT NOT NULL
);
COMMENT ON TABLE func_def IS 'Function attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE func_arg_def(
  code  TEXT REFERENCES func_def ON DELETE CASCADE
, is_in BOOL
, lang  TEXT DEFAULT 'ru'
, arg   TEXT
, anno  TEXT NOT NULL
, CONSTRAINT func_arg_def_pkey PRIMARY KEY (code, is_in, lang, arg)
);
COMMENT ON TABLE func_arg_def IS 'Function in/out argument attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE func_arg_def_common(
  lang  TEXT DEFAULT 'ru'
, arg   TEXT
, anno  TEXT NOT NULL
, CONSTRAINT func_arg_def_common_pkey PRIMARY KEY (lang, arg)
);
COMMENT ON TABLE func_arg_def_common IS 'Common function in/out argument attributes';

/* ------------------------------------------------------------------------- */
