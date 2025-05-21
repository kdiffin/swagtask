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

-- name: CreateCollaboratorVaultRelation :exec
WITH authorized_user AS (
  SELECT 1
  FROM vault_user_relations
  WHERE user_id = sqlc.arg('user_id')::UUID
    AND vault_id = sqlc.arg('vault_id')::UUID
    AND role = 'owner'
),
user_id_from_username AS (
  SELECT id FROM users WHERE users.username = sqlc.arg('collaborator_username')
)
INSERT INTO vault_user_relations (vault_id, user_id, role)
SELECT sqlc.arg('vault_id'), u.id, sqlc.arg('role')
FROM user_id_from_username u 
WHERE
    EXISTS (SELECT 1 FROM authorized_user);


-- name: DeleteCollaboratorVaultRelation :exec
WITH user_id_from_username AS (
  SELECT id FROM users WHERE users.username = sqlc.arg('collaborator_username')
)
DELETE FROM vault_user_relations v
WHERE
  v.vault_id = $1 AND v.user_id = user_id_from_username.id
  -- owner authorziation
  AND EXISTS (
		SELECT 1 FROM vault_user_relations rel
		WHERE 
			rel.user_id = sqlc.arg('user_id')::UUID
			AND rel.role = 'owner'
	);
  