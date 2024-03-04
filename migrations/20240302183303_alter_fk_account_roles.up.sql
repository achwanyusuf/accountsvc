ALTER TABLE "account_roles" ADD CONSTRAINT fk_account_roles_a_key FOREIGN KEY("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "account_roles" ADD CONSTRAINT fk_account_roles_r_key FOREIGN KEY("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;