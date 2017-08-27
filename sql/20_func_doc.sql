/*
  Tables for stored proc documenting

*/

-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS func_def(
  code    TEXT PRIMARY KEY
, nspname NAME NOT NULL
, proname NAME NOT NULL
, max_age INTEGER NOT NULL DEFAULT 0    -- 0: forever, >0 secs, -1: no cache
, anno    TEXT NOT NULL
, sample  TEXT NOT NULL
);
COMMENT ON TABLE func_def IS 'Function attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS func_arg_def(
  code  TEXT REFERENCES func_def ON DELETE CASCADE
, is_in BOOL
, arg   TEXT
, anno  TEXT NOT NULL
, CONSTRAINT func_arg_def_pkey PRIMARY KEY (code, is_in, arg)
);
COMMENT ON TABLE func_arg_def IS 'Function in/out argument attributes';

-- -----------------------------------------------------------------------------
