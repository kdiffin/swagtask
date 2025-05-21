-- name: SignUpAndCreateDefaultVault :exec
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