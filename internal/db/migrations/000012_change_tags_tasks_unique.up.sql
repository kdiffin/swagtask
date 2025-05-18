ALTER TABLE tasks
DROP CONSTRAINT unique_user_idea_pair;

ALTER TABLE tasks ADD CONSTRAINT unique_vault_idea_pair UNIQUE (vault_id, idea);

ALTER TABLE tags
DROP CONSTRAINT unique_user_name_pair;

ALTER TABLE tags ADD CONSTRAINT unique_vault_name_pair UNIQUE (vault_id, name);