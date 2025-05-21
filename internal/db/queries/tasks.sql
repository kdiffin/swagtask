-- READ
-- name: GetTasksWithTagRelations :many
WITH t_with_author AS (
  SELECT t.id, t.name, t.idea, t.completed, t.user_id, t.vault_id, t.created_at, t.updated_at,
         u.path_to_pfp, u.username
  FROM tasks t
  JOIN users u ON t_with_author.user_id = u.id
)
SELECT t_with_author.name, t_with_author.idea, t_with_author.ID, t_with_author.vault_id, t_with_author.completed, t_with_author.user_id, t_with_author.created_at, t_with_author.updated_at,
		tg.ID AS tag_id, tg.name AS tag_name, tg.user_id AS tag_user_id,
    	t_with_author.path_to_pfp AS author_path_to_pfp, t_with_author.username AS author_username
FROM t_with_author
LEFT JOIN tag_task_relations rel 
	ON t_with_author.ID = rel.task_id 
LEFT JOIN tags tg 
	ON tg.ID = rel.tag_id
WHERE
  	-- authorization, checks if user is inside of this vault
	EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	)
ORDER BY t_with_author.created_at DESC; 

-- name: GetTaskWithTagRelations :many 
WITH t_with_author AS (
  SELECT t.id, t.name, t.idea, t.completed, t.user_id, t.vault_id, t.created_at, t.updated_at,
         u.path_to_pfp, u.username
  FROM tasks t
  JOIN users u ON t.user_id = u.id
)
SELECT t_with_author.name, t_with_author.idea, t_with_author.ID, t_with_author.vault_id, t_with_author.completed, t_with_author.user_id, t_with_author.created_at, t_with_author.updated_at,
		tg.ID AS tag_id, tg.name AS tag_name, tg.user_id AS tag_user_id,
    	t_with_author.path_to_pfp AS author_path_to_pfp, t_with_author.username AS author_username
FROM t_with_author
LEFT JOIN tag_task_relations rel
	ON t_with_author.ID = rel.task_id
LEFT JOIN tags tg 
	ON rel.tag_id = tg.ID
WHERE t_with_author.ID = sqlc.arg('id')::UUID 
  	-- authorization, checks if user is inside of this vault
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
)
ORDER BY t_with_author.created_at DESC;

-- name: GetFilteredTasks :many
WITH t_with_author AS (
  SELECT t.id, t.name, t.idea, t.completed, t.user_id, t.vault_id, t.created_at, t.updated_at,
         u.path_to_pfp, u.username
  FROM tasks t
  JOIN users u ON t.user_id = u.id
)
SELECT t_with_author.name, t_with_author.idea, t_with_author.ID, t_with_author.vault_id, t_with_author.completed, t_with_author.user_id, t_with_author.created_at, t_with_author.updated_at,
		tg.ID AS tag_id, tg.name AS tag_name, tg.user_id AS tag_user_id,
    	t_with_author.path_to_pfp AS author_path_to_pfp, t_with_author.username AS author_username
FROM t_with_author
LEFT JOIN tag_task_relations rel ON rel.task_id = t_with_author.ID
LEFT JOIN tags tg ON tg.ID = rel.tag_id
WHERE
	-- where the name is like the task filter (if the filter exists)
	(t_with_author.name ILIKE '%' || COALESCE(sqlc.narg('task_name')::text, t_with_author.name) || '%')
	AND
	-- if the tag filter exists, return the rows of the tasks who have a relation to that tag  
	(sqlc.narg('tag_name')::text IS NULL OR 
		EXISTS (
			SELECT 1
			FROM tag_task_relations r2
			-- to get tag id from name
			JOIN tags tg2 
				ON tg2.name = sqlc.narg('tag_name')::text
			WHERE r2.task_id = t_with_author.ID AND r2.tag_id = tg2.id
		)
	)
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	)
ORDER BY t_with_author.created_at DESC; 

-- TODO: reimplement
-- name: GetPreviousTaskDetails :one
SELECT name, id FROM tasks 
WHERE created_at < $1
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	)  
ORDER BY created_at DESC LIMIT 1;

-- name: GetNextTaskDetails :one
SELECT name, id FROM tasks 
WHERE created_at > $1
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	)
ORDER BY created_at ASC LIMIT 1;

-- CREATE
-- name: CreateTask :one
WITH authorized_user AS (
  SELECT 1
  FROM vault_user_relations
  WHERE user_id = sqlc.arg('user_id')::UUID
    AND vault_id = sqlc.arg('vault_id')::UUID
    AND (role = 'owner' OR role = 'collaborator')
)
INSERT INTO tasks (name, idea, user_id, vault_id) 
SELECT sqlc.arg('name'), sqlc.arg('idea'), sqlc.arg('user_id')::UUID, sqlc.arg('vault_id')::UUID
FROM authorized_user RETURNING id;

-- UPDATE
-- name: ToggleTaskCompletion :exec
UPDATE tasks SET completed = NOT completed
WHERE 
	id = sqlc.arg('id')::UUID
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	    AND (role = 'owner' OR role = 'collaborator')
);	

-- name: UpdateTask :exec
UPDATE tasks
SET
  name = COALESCE(sqlc.narg('name'), name),
  idea = COALESCE(sqlc.narg('idea'), idea)
WHERE 
	id = sqlc.arg('id')::UUID
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	    AND (role = 'owner' OR role = 'collaborator')
);	

-- DELETE 
-- name: DeleteTask :exec
DELETE FROM tasks t
WHERE 
	t.id = $1
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
		AND (role = 'owner' OR role = 'collaborator')
); 
