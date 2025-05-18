-- name: GetAllTagsDesc :many
SELECT * FROM tags
WHERE user_id = $1 AND vault_id = $2
ORDER BY id DESC;

-- name: CreateTag :exec
INSERT INTO tags (name, user_id, vault_id)
VALUES ($1, $2, $3);

-- name: GetTagWithTaskRelations :many
WITH tg_u AS (
  SELECT tg.id, tg.name, tg.user_id, tg.vault_id, tg.created_at, tg.updated_at,
         u.path_to_pfp, u.username
  FROM tags tg
  JOIN users u ON tg.user_id = u.id
)
SELECT tg_u.id, tg_u.name, tg_u.user_id, tg_u.vault_id, tg_u.created_at, tg_u.updated_at,
       t.id AS task_id, t.name AS task_name, t.user_id AS task_user_id,
       tg_u.path_to_pfp AS author_path_to_pfp, tg_u.username AS author_username
FROM tg_u
LEFT JOIN tag_task_relations rel ON tg_u.id = rel.tag_id
LEFT JOIN tasks t ON t.id = rel.task_id
WHERE tg_u.id = sqlc.arg('id')::UUID AND tg_u.user_id = sqlc.arg('user_id')::UUID AND tg_u.vault_id = sqlc.arg('vault_id')::UUID;

-- name: GetTagsWithTaskRelations :many
WITH tg_u AS (
  SELECT tg.id, tg.name, tg.user_id, tg.vault_id, tg.created_at, tg.updated_at,
         u.path_to_pfp, u.username
  FROM tags tg
  JOIN users u ON tg.user_id = u.id
)
SELECT tg_u.id, tg_u.name, tg_u.user_id, tg_u.vault_id, tg_u.created_at, tg_u.updated_at,
       t.id AS task_id, t.name AS task_name, t.user_id AS task_user_id,
       tg_u.path_to_pfp AS author_path_to_pfp, tg_u.username AS author_username
FROM tg_u
LEFT JOIN tag_task_relations rel ON tg_u.id = rel.tag_id
LEFT JOIN tasks t ON t.id = rel.task_id
WHERE tg_u.user_id = sqlc.arg('user_id')::UUID  AND tg_u.vault_id = sqlc.arg('vault_id')::UUID;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1 AND user_id = $2 AND vault_id = $3;

-- name: UpdateTag :exec
UPDATE tags
SET name = $1
WHERE id = $2 AND user_id = $3 AND vault_id = $4;
