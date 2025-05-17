-- change the schemas
ALTER TABLE users
ADD COLUMN default_vault_id UUID NOT NULL;

ALTER TABLE users ADD CONSTRAINT fk_default_vault_id FOREIGN KEY (default_vault_id) REFERENCES vaults (id);

CREATE TYPE vault_role_type AS ENUM ('default', 'collaborative', 'public');

ALTER TABLE vaults
ADD COLUMN kind vault_role_type NOT NULL DEFAULT 'collaborative';