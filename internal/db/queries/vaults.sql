-- name: CreateDefaultVault :one
INSERT INTO vaults (name, description, kind) VALUES($1, $2, 'default') RETURNING id;