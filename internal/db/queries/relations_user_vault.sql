

-- name: CreateCollaboratorVaultRelation :exec
WITH authorized_user AS (
  SELECT 1
  FROM vault_user_relations auth_rel
  WHERE auth_rel.user_id = sqlc.arg('user_id')::UUID
    AND auth_rel.vault_id = sqlc.arg('vault_id')::UUID
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
  SELECT id FROM users WHERE username = sqlc.arg('collaborator_username')
)
DELETE FROM vault_user_relations v
WHERE
  v.vault_id = sqlc.arg('vault_id')::UUID AND v.user_id = (SELECT id FROM user_id_from_username)
  -- owner authorziation
  AND EXISTS (
		SELECT 1 FROM vault_user_relations auth_rel
		WHERE 
			auth_rel.user_id = sqlc.arg('user_id')::UUID
      AND auth_rel.vault_id = sqlc.arg('vault_id')::UUID
			AND auth_rel.role = 'owner'
	);
  