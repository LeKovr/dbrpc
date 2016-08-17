/*
  Create schema & set search_path

*/

\set SCH dbrpc

CREATE SCHEMA :SCH;

SET SEARCH_PATH = :SCH, public;