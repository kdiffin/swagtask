-- name: CreateDefaultVault :one
BEGIN;
INSERT INTO vaults (name, description, kind) VALUES($1, $2, 'default') RETURNING id;
INSERT INTO vault_user_relations (vault_id, role, user_id) VALUES(id, 'owner', $3);
COMMIT;
