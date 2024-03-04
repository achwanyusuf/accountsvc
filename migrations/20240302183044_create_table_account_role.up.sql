CREATE SEQUENCE account_role_id_seq;

CREATE TABLE IF NOT EXISTS account_roles (
  id integer primary key DEFAULT nextval('account_role_id_seq'),
  account_id integer NOT NULL,
  role_id integer NOT NULL,
  created_by integer default 0 NOT NULL,
  created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_by integer default 0 NOT NULL,
  updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_by integer,
  deleted_at timestamp WITH TIME ZONE
);

ALTER SEQUENCE account_role_id_seq OWNED BY account_roles.id;