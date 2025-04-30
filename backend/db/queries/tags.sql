-- name: GetAllTagsDesc :many
SELECT * FROM tags ORDER BY id DESC;

-- name: CreateTag :exec
INSERT INTO tags (name) VALUES($1);

-- name: GetTagWithTaskRelations :many
SELECT tg.ID, tg.name, t.ID AS task_id, t.name AS task_name 
    FROM tags tg
    LEFT JOIN tag_task_relations rel 
        ON tg.ID = rel.tag_id
    LEFT JOIN tasks t 
        ON t.ID = rel.task_id
    WHERE tg.id = $1;



-- name: GetTagsWithTaskRelations :many
SELECT tg.ID, tg.name, t.ID AS task_id, t.name AS task_name 
    FROM tags tg
    LEFT JOIN tag_task_relations rel 
        ON tg.ID = rel.tag_id
    LEFT JOIN tasks t 
        ON t.ID = rel.task_id;


-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: UpdateTag :exec
UPDATE tags SET name = $1 WHERE id = $2;