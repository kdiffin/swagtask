-- Drop constraint first, then drop column
ALTER TABLE tasks
DROP CONSTRAINT fk_tasks_user;

ALTER TABLE tasks
DROP COLUMN user_id;

ALTER TABLE tags
DROP CONSTRAINT fk_tags_user;

ALTER TABLE tags
DROP COLUMN user_id;