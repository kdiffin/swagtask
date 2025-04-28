-- name: GetTagsOfTask :many
SELECT tg.ID, tg.name FROM tags tg
INNER JOIN tag_task_relations rel ON tg.ID = rel.tag_id
WHERE rel.task_id = $1;

-- name: CreateTagTaskRelation :exec
INSERT INTO tag_task_relations (task_id, tag_id) VALUES($1, $2);

-- name: DeleteAllTagRelationsForTask :exec
DELETE FROM tag_task_relations WHERE task_id = $1;


-- name: DeleteSingleTagRelation :exec
DELETE FROM tag_task_relations WHERE task_id = $1 AND tag_id = $2;