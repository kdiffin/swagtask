-- name: GetAllTagsDesc :many
SELECT * FROM tags
WHERE vault_id = sqlc.arg('vault_id')::UUID
-- authorization, checks if user is inside of this vault
  AND EXISTS(
    SELECT 1 FROM vault_user_relations v_u_rel
      WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
      AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
  )
ORDER BY created_at DESC;

-- name: CreateTag :exec
WITH authorized_user AS (
  SELECT 1
  FROM vault_user_relations
  WHERE user_id = sqlc.arg('user_id')::UUID
    AND vault_id = sqlc.arg('vault_id')::UUID
    AND (role = 'owner' OR role = 'collaborator')
)
INSERT INTO tags (name, user_id, vault_id)
SELECT sqlc.arg('name'), sqlc.arg('user_id')::UUID, sqlc.arg('vault_id')::UUID
FROM authorized_user RETURNING id;

-- name: GetTagWithTaskRelations :many
WITH tg_author AS (
  SELECT tg.id, tg.name, tg.user_id, tg.vault_id, tg.created_at, tg.updated_at,
         u.path_to_pfp, u.username
  FROM tags tg
  JOIN users u ON tg.user_id = u.id
)
SELECT tg_author.id, tg_author.name, tg_author.user_id, tg_author.vault_id, tg_author.created_at, tg_author.updated_at,
       t.id AS task_id, t.name AS task_name, t.user_id AS task_user_id,
       tg_author.path_to_pfp AS author_path_to_pfp, tg_author.username AS author_username
FROM tg_author
LEFT JOIN tag_task_relations rel ON tg_author.id = rel.tag_id
LEFT JOIN tasks t ON t.id = rel.task_id
WHERE tg_author.id = sqlc.arg('id')::UUID 
  AND tg_author.vault_id = sqlc.arg('vault_id')::UUID
  -- authorization, checks if user is inside of this vault
  AND EXISTS(
    SELECT 1 FROM vault_user_relations v_u_rel
      WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
      AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
  )
ORDER BY tg_author.created_at DESC;

-- name: GetTagsWithTaskRelations :many
WITH tg_author AS (
  SELECT tg.id, tg.name, tg.user_id, tg.vault_id, tg.created_at, tg.updated_at,
         u.path_to_pfp, u.username
  FROM tags tg
  JOIN users u ON tg.user_id = u.id
)
SELECT tg_author.id, tg_author.name, tg_author.user_id, tg_author.vault_id, tg_author.created_at, tg_author.updated_at,
       t.id AS task_id, t.name AS task_name, t.user_id AS task_user_id,
       tg_author.path_to_pfp AS author_path_to_pfp, tg_author.username AS author_username
FROM tg_author
LEFT JOIN tag_task_relations rel ON tg_author.id = rel.tag_id
LEFT JOIN tasks t ON t.id = rel.task_id
WHERE tg_author.vault_id = sqlc.arg('vault_id')::UUID
  -- authorization, checks if user is inside of this vault
  AND EXISTS(
    SELECT 1 FROM vault_user_relations v_u_rel
      WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
      AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
  )
ORDER BY tg_author.created_at DESC;

-- name: DeleteTag :exec
DELETE FROM tags tg
WHERE tg.id = $1 AND tg.vault_id = sqlc.arg('vault_id')::UUID
  -- authorization, checks if user is inside of this vault
  AND EXISTS(
    SELECT 1 FROM vault_user_relations v_u_rel
      WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
      AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
      AND (role = 'owner' OR role = 'collaborator')
  );

-- name: UpdateTag :exec
UPDATE tags 
SET name = sqlc.arg('name')::TEXT
WHERE id = sqlc.arg('id')::UUID AND vault_id = sqlc.arg('vault_id')::UUID
  -- authorization, checks if user is inside of this vault
 AND EXISTS(
    SELECT 1 FROM vault_user_relations v_u_rel
      WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
      AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
      AND (role = 'owner' OR role = 'collaborator')
  );

