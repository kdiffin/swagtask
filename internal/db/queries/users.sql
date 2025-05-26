

-- name: SignUpAndCreateDefaultVault :exec
WITH
  new_vault AS (
    INSERT INTO vaults (name, description, kind)
    VALUES ('Default', 'This is your default vault. Only you can access this.', 'default')
    RETURNING id
  ),
  new_user AS (
    INSERT INTO users (username, password_hash, default_vault_id, path_to_pfp)
    VALUES ($1, $2, (SELECT id FROM new_vault), $3)
    RETURNING id
  )
    INSERT INTO vault_user_relations (vault_id, role, user_id)
    VALUES (
      (SELECT id FROM new_vault),
      'owner',
      (SELECT id FROM new_user)
    );


-- name: SignUpAndCreateDefaultVaultNoPfp :exec
WITH
  new_vault AS (
    INSERT INTO vaults (name, description, kind)
    VALUES ('Default', 'This is your default vault. Only you can access this.', 'default')
    RETURNING id
  ),
  new_user AS (
    INSERT INTO users (username, password_hash, default_vault_id)
    VALUES ($1, $2, (SELECT id FROM new_vault))
    RETURNING id
  )
    INSERT INTO vault_user_relations (vault_id, role, user_id)
    VALUES (
      (SELECT id FROM new_vault),
      'owner',
      (SELECT id FROM new_user)
    );

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