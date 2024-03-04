DELETE FROM account_roles WHERE id in (1, 2, 3);
ALTER SEQUENCE account_role_id_seq RESTART WITH 1;