ALTER TABLE tasks DROP CONSTRAINT unique_user_idea_pair;
ALTER TABLE tasks ADD CONSTRAINT tasks_idea_key UNIQUE(idea);

ALTER TABLE tags DROP CONSTRAINT unique_user_name_pair;
ALTER TABLE tags ADD CONSTRAINT  tags_name_key UNIQUE(name);