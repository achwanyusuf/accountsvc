CREATE SEQUENCE account_id_seq;

CREATE TABLE IF NOT EXISTS accounts (
  id integer primary key DEFAULT nextval('account_id_seq'),
  email    varchar(100) unique NOT NULL,
  password   varchar(100) NOT NULL,
  name varchar(255) NOT NULL,
  created_by integer default 0 NOT NULL,
  created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_by integer default 0 NOT NULL,
  updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_by integer,
  deleted_at timestamp WITH TIME ZONE
);

ALTER SEQUENCE account_id_seq OWNED BY accounts.id;