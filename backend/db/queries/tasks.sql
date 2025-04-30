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
SELECT t.ID, t.name, t.idea, t.completed, tg.ID AS tag_id, tg.name AS tag_name
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

-- name: GetFilteredTasks :many
SELECT t.ID, t.name, t.idea, t.completed, tg.ID as tag_id, tg.name AS tag_name
FROM tasks t
LEFT JOIN tag_task_relations rel ON rel.task_id = t.ID
LEFT JOIN tags tg ON tg.ID = rel.tag_id
WHERE
	-- where the name is like the task filter (if the filter exists)
	(t.name ILIKE '%' || COALESCE(sqlc.narg('task_name')::text, t.name) || '%')
	AND
	-- if the tag filter exists, return the rows of the tasks who have a relation to that tag  
	(sqlc.narg('tag_name')::text IS NULL OR 
		EXISTS (
			SELECT 1
			FROM tag_task_relations r2
			-- to get tag id from name
			JOIN tags tg2 
				ON tg2.name = sqlc.narg('tag_name')::text
			WHERE r2.task_id = t.ID AND r2.tag_id = tg2.id 
		)
	);


-- DELETE 
-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;