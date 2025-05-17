-- name: CreateUser :exec
INSERT INTO users (username, password_hash, default_vault_id) VALUES ($1, $2, $3);

-- name: GetUserCredentials :one
SELECT id, password_hash FROM users WHERE username=$1;

-- name: CreateSession :exec
INSERT INTO sessions (id, user_id) VALUES($1,$2);

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: GetSessionValues :one
SELECT id, user_id FROM sessions WHERE id = $1;

-- name: GetUserInfo :one
SELECT username, path_to_pfp, default_vault_id FROM users WHERE id = $1;