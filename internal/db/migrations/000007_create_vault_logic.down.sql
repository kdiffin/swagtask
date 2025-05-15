ALTER TABLE tasks
DROP CONSTRAINT fk_tasks_vault;

ALTER TABLE tasks
DROP COLUMN vault_id;

ALTER TABLE tags
DROP CONSTRAINT fk_tags_vault;

ALTER TABLE tags
DROP COLUMN vault_id;

ALTER TABLE vault_user_relations
DROP CONSTRAINT fk_user_id;

ALTER TABLE vault_user_relations
DROP CONSTRAINT fk_vault_id;

DROP TABLE IF EXISTS vault_user_relations;

DROP TABLE IF EXISTS vaults;