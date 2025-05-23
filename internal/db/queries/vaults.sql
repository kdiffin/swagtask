-- name: GetVaultsWithCollaborators :many
SELECT v.name, v.description, v.ID, v.locked, v.kind, v.created_at, v.updated_at, rel.role,
		us.username AS collaborator_username, us.path_to_pfp AS collaborator_path_to_pfp
FROM vaults v
JOIN vault_user_relations rel 
	ON v.id = rel.vault_id 
JOIN users us
	ON us.id = rel.user_id
WHERE
    EXISTS (
        SELECT
            1
        FROM
            vault_user_relations r_user_vault_check
        WHERE
            r_user_vault_check.vault_id = v.id 
            AND r_user_vault_check.user_id = sqlc.arg('user_id')::UUID
    )
ORDER BY
    v.created_at DESC, 
    us.username ASC;  

-- name: GetVaultWithCollaborators :many
SELECT v.name, v.description, v.ID, v.locked, v.kind, v.created_at, v.updated_at, rel.role,
		us.username AS collaborator_username, us.path_to_pfp AS collaborator_path_to_pfp
FROM vaults v
LEFT JOIN vault_user_relations rel 
	ON v.id = rel.vault_id 
LEFT JOIN users us
	ON us.id = rel.user_id
WHERE
	EXISTS (
        SELECT
            1
        FROM
            vault_user_relations r_user_vault_check
        WHERE
            r_user_vault_check.vault_id = v.id 
            AND r_user_vault_check.user_id = sqlc.arg('user_id')::UUID
			AND r_user_vault_check.vault_id = sqlc.arg('vault_id')::UUID
			
    )
ORDER BY v.created_at DESC;

-- TODO: authenticate
-- CREATE
-- name: CreateVault :one
WITH 
	new_vault AS (
		INSERT INTO vaults (name, description, kind) 
		VALUES ($1, $2, 'collaborative') 
		RETURNING id
	)
	INSERT INTO vault_user_relations (vault_id, user_id, role)
		VALUES(
			(SELECT id from new_vault), $3, 'owner'
		)
		RETURNING (SELECT id from new_vault);


-- name: DeleteVault :exec
DELETE FROM vaults v
WHERE 
	v.id = $1
	-- authorization part, checks if person is owner of vault
	AND EXISTS (
		SELECT 1 FROM vault_user_relations rel
		WHERE 
			rel.user_id = sqlc.arg('user_id')::UUID
			AND rel.vault_id = $1
			AND rel.role = 'owner'
	);

-- name: UpdateVault :exec
UPDATE vaults v
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  locked = sqlc.arg('locked')
WHERE 
	v.id = $1
	-- authorization part, checks if person is owner
	AND EXISTS (
		SELECT 1 FROM vault_user_relations rel
		WHERE 
			rel.user_id = sqlc.arg('user_id')::UUID
			AND rel.vault_id = $1
			AND rel.role = 'owner'
);