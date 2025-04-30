-- READ
-- name: GetTasksWithTagRelations :many
SELECT t.name, t.Idea, t.ID, t.completed, tg.ID AS tag_id, tg.name AS tag_name
		FROM tasks t
		LEFT JOIN tag_task_relations rel 
			ON t.ID = rel.task_id 
		LEFT JOIN tags tg 
			ON tg.ID = rel.tag_id
		ORDER BY t.ID DESC;

-- name: GetTaskWithTagRelations :many
SELECT t.ID, t.name, t.Idea, t.completed, tg.ID AS tag_id, tg.name AS tag_name
	FROM tasks t
	LEFT JOIN tag_task_relations rel
		ON t.ID = rel.task_id
	LEFT JOIN tags tg 
		ON rel.tag_id = tg.ID
	WHERE t.ID = $1;


-- CREATE
-- name: CreateTask :one
INSERT INTO tasks (name, idea) VALUES ($1, $2) RETURNING *;


-- UPDATE
-- name: ToggleTaskCompletion :exec
UPDATE tasks SET completed = NOT completed WHERE id = $1;

-- name: UpdateTask :exec
UPDATE tasks
SET
  name = COALESCE(sqlc.narg('name'), name),
  idea = COALESCE(sqlc.narg('idea'), idea)
WHERE id = sqlc.arg('id');

-- DELETE 
-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;