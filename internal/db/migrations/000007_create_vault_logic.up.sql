CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- creates vaults and vault-user relation
CREATE TABLE
    vaults (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        locked BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now ()
    );

CREATE TYPE vault_rel_role_type AS ENUM ('owner', 'collaborator', 'viewer');

CREATE TABLE
    vault_user_relations (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        vault_id UUID NOT NULL,
        user_id INTEGER NOT NULL,
        role vault_rel_role_type NOT NULL
    );

ALTER TABLE vault_user_relations ADD CONSTRAINT unique_vault_user_pair UNIQUE (vault_id, user_id);

ALTER TABLE vault_user_relations ADD CONSTRAINT fk_vault_id FOREIGN KEY (vault_id) REFERENCES vaults (id) ON DELETE CASCADE;

ALTER TABLE vault_user_relations ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

-- adds vault data to users
ALTER TABLE tasks
ADD COLUMN vault_id UUID NOT NULL;

ALTER TABLE tasks ADD CONSTRAINT fk_tasks_vault FOREIGN KEY (vault_id) REFERENCES vaults (id) ON DELETE CASCADE;

ALTER TABLE tags
ADD COLUMN vault_id UUID NOT NULL;

ALTER TABLE tags ADD CONSTRAINT fk_tags_vault FOREIGN KEY (vault_id) REFERENCES vaults (id) ON DELETE CASCADE;