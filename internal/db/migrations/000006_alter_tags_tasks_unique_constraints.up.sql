ALTER TABLE tasks
DROP CONSTRAINT tasks_idea_key;

ALTER TABLE tasks ADD CONSTRAINT unique_user_idea_pair UNIQUE (user_id, idea);

ALTER TABLE tags
DROP CONSTRAINT tags_name_key;

ALTER TABLE tags ADD CONSTRAINT unique_user_name_pair UNIQUE (user_id, name);