CREATE SEQUENCE role_id_seq;

CREATE TABLE IF NOT EXISTS roles (
  id integer primary key DEFAULT nextval('role_id_seq'),
  scope varchar(3) NOT NULL,
  cid uuid NOT NULL,
  sec text NOT NULL,
  created_by integer default 0 NOT NULL,
  created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_by integer default 0 NOT NULL,
  updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_by integer,
  deleted_at timestamp WITH TIME ZONE
);

ALTER SEQUENCE role_id_seq OWNED BY roles.id;