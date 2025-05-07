-- name: CreateUser :exec
INSERT INTO users (username, password_hash) VALUES ($1, $2);

-- name: GetUserCredentials :one
SELECT id, password_hash FROM users WHERE username=$1;