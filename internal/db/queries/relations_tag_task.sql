-- this is for the UI whenever we wanna add a relation from a task to a tag
-- name: GetAllTaskOptions :many
SELECT name, id from tasks WHERE user_id = $1 ORDER BY id DESC;

-- name: CreateTagTaskRelation :exec
INSERT INTO tag_task_relations (task_id, tag_id)
VALUES($1, $2);

-- name: DeleteTagTaskRelation :exec
DELETE FROM tag_task_relations
WHERE task_id = $1 AND tag_id = $2;
