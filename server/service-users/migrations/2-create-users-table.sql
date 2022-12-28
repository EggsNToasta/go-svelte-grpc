-- +migrate Up
CREATE TABLE
  users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    created timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted boolean NOT NULL DEFAULT FALSE,
    email text UNIQUE NOT NULL,
    role text NOT NULL,
    "lastLogin" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

-- +migrate Up
CREATE TRIGGER set_timestamp BEFORE
UPDATE
  ON users FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp ();
