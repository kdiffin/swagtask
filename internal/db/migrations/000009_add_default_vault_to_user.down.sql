-- change the schemas
ALTER TABLE users
DROP CONSTRAINT fk_default_vault_id;

ALTER TABLE users
DROP COLUMN default_vault_id;

ALTER TABLE vaults
DROP COLUMN kind;

DROP TYPE vault_role_type;