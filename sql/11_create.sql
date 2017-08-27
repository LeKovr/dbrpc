/*

  Create database schema

*/
-- -----------------------------------------------------------------------------

CREATE SCHEMA IF NOT EXISTS rpc;
COMMENT ON SCHEMA rpc IS 'RPC method registry';

SET SEARCH_PATH = 'rpc', 'public';
-- -----------------------------------------------------------------------------
