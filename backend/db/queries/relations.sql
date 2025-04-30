-- name: GetTagsOfTask :many
SELECT tg.ID, tg.name FROM tags tg
INNER JOIN tag_task_relations rel ON tg.ID = rel.tag_id
WHERE rel.task_id = $1;

-- this is for the UI whenever we wanna add a relation from a task to a tag
-- name: GetAllTaskOptions :many
SELECT name, id from tasks ORDER BY id DESC;


-- name: CreateTagTaskRelation :exec
INSERT INTO tag_task_relations (task_id, tag_id) VALUES($1, $2);

-- name: DeleteAllTaskRelations :exec
DELETE FROM tag_task_relations WHERE task_id = $1;

-- name: DeleteAllTagRelations :exec
DELETE FROM tag_task_relations WHERE tag_id = $1;

-- name: DeleteTagTaskRelation :exec
DELETE FROM tag_task_relations WHERE task_id = $1 AND tag_id = $2;
