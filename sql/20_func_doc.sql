/*
  Tables for stored proc documenting

*/

/* ------------------------------------------------------------------------- */

CREATE TABLE func_def(
  code    TEXT PRIMARY KEY
, nspname NAME NOT NULL
, proname NAME NOT NULL
, anno    TEXT
, sample  TEXT
);
COMMENT ON TABLE func_def IS 'Function attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE func_arg_def(
  code  TEXT REFERENCES func_def ON DELETE CASCADE
, is_in BOOL
, lang  TEXT DEFAULT 'ru'
, arg   TEXT
, anno  TEXT
, CONSTRAINT func_arg_def_pkey PRIMARY KEY (code, is_in, lang, arg)
);
COMMENT ON TABLE func_arg_def IS 'Function in/out argument attributes';

-- -----------------------------------------------------------------------------

CREATE TABLE func_arg_def_common(
  lang  TEXT NOT NULL DEFAULT 'ru'
, arg   TEXT
, anno  TEXT
, CONSTRAINT func_arg_def_common_pkey PRIMARY KEY (lang, arg)
);
COMMENT ON TABLE func_arg_def_common IS 'Common function in/out argument attributes';

/* ------------------------------------------------------------------------- */
