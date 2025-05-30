-- this is for the UI whenever we wanna add a relation from a task to a tag
-- name: GetAllTaskOptions :many
SELECT name, id from tasks 
WHERE 
    vault_id = sqlc.arg('vault_id')::UUID
  	-- authorization, checks if user is inside of this vault
	AND EXISTS(
		SELECT 1 FROM vault_user_relations v_u_rel
		WHERE v_u_rel.user_id = sqlc.arg('user_id')::UUID 
		AND v_u_rel.vault_id = sqlc.arg('vault_id')::UUID
	);

-- name: CreateTagTaskRelation :exec
INSERT INTO tag_task_relations (task_id, tag_id)
VALUES($1, $2);

-- name: DeleteTagTaskRelation :exec
DELETE FROM tag_task_relations
WHERE task_id = $1 AND tag_id = $2;
