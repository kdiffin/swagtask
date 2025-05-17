-- name: GetAllTagsDesc :many
SELECT * FROM tags WHERE user_id = $1 AND vault_id = $2 ORDER BY id DESC;

-- name: CreateTag :exec
INSERT INTO tags (name, user_id, vault_id) VALUES($1, $2, $3);

-- name: GetTagWithTaskRelations :many
SELECT tg.ID, tg.name, tg.user_id, tg.created_at, tg.updated_at,
    t.ID AS task_id, t.name AS task_name, t.user_id AS task_user_id
    FROM tags tg
    LEFT JOIN tag_task_relations rel 
        ON tg.ID = rel.tag_id
    LEFT JOIN tasks t 
        ON t.ID = rel.task_id
    WHERE tg.id = $1 AND tg.user_id = $2 AND tg.vault_id = $3;



-- name: GetTagsWithTaskRelations :many
SELECT tg.ID, tg.name, tg.user_id, tg.created_at, tg.updated_at,
    t.ID AS task_id, t.name AS task_name, t.user_id AS task_user_id
    FROM tags tg
    LEFT JOIN tag_task_relations rel 
        ON tg.ID = rel.tag_id
    LEFT JOIN tasks t 
        ON t.ID = rel.task_id
    WHERE tg.user_id = $1 AND tg.vault_id = $2;
-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1 AND user_id = $2 AND vault_id = $3;
-- name: UpdateTag :exec
UPDATE tags SET name = $1 WHERE id = $2 AND user_id = $3 AND vault_id = $4;