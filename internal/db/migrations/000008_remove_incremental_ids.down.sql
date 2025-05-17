-- gippity made
-- Recreate sequences
CREATE SEQUENCE IF NOT EXISTS tag_task_relations_id_seq;
CREATE SEQUENCE IF NOT EXISTS tags_id_seq;
CREATE SEQUENCE IF NOT EXISTS tasks_id_seq;
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS sessions_id_seq;

-- Drop all relevant foreign keys before altering types
ALTER TABLE tag_task_relations DROP CONSTRAINT IF EXISTS fk_tag;
ALTER TABLE tag_task_relations DROP CONSTRAINT IF EXISTS fk_task;
ALTER TABLE tags DROP CONSTRAINT IF EXISTS fk_tags_user;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_user;
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS fk_sessions_user;
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS fk_user_id;
ALTER TABLE vault_user_relations DROP CONSTRAINT IF EXISTS fk_user_id;

-- Tag-Task table
ALTER TABLE tag_task_relations
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN tag_id DROP DEFAULT,
    ALTER COLUMN task_id DROP DEFAULT;

ALTER TABLE tag_task_relations
    ALTER COLUMN id SET DATA TYPE integer USING (nextval('tag_task_relations_id_seq')),
    ALTER COLUMN tag_id SET DATA TYPE integer USING (nextval('tags_id_seq')),
    ALTER COLUMN task_id SET DATA TYPE integer USING (nextval('tasks_id_seq'));

ALTER TABLE tag_task_relations
    ALTER COLUMN id SET DEFAULT nextval('tag_task_relations_id_seq'),
    ALTER COLUMN tag_id SET DEFAULT nextval('tags_id_seq'),
    ALTER COLUMN task_id SET DEFAULT nextval('tasks_id_seq');

-- Tags table
ALTER TABLE tags
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE tags
    ALTER COLUMN id SET DATA TYPE integer USING (nextval('tags_id_seq')),
    ALTER COLUMN user_id SET DATA TYPE integer USING (nextval('users_id_seq'));

ALTER TABLE tags
    ALTER COLUMN id SET DEFAULT nextval('tags_id_seq'),
    ALTER COLUMN user_id SET DEFAULT nextval('users_id_seq');

-- Tasks table
ALTER TABLE tasks
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE tasks
    ALTER COLUMN id SET DATA TYPE integer USING (nextval('tasks_id_seq')),
    ALTER COLUMN user_id SET DATA TYPE integer USING (nextval('users_id_seq'));

ALTER TABLE tasks
    ALTER COLUMN id SET DEFAULT nextval('tasks_id_seq'),
    ALTER COLUMN user_id SET DEFAULT nextval('users_id_seq');

-- Users table
ALTER TABLE users
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE users
    ALTER COLUMN id SET DATA TYPE integer USING (nextval('users_id_seq'));

ALTER TABLE users
    ALTER COLUMN id SET DEFAULT nextval('users_id_seq');

-- Sessions table
ALTER TABLE sessions
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE sessions
    ALTER COLUMN id SET DATA TYPE integer USING (nextval('sessions_id_seq')),
    ALTER COLUMN user_id SET DATA TYPE integer USING (nextval('users_id_seq'));

ALTER TABLE sessions
    ALTER COLUMN id SET DEFAULT nextval('sessions_id_seq'),
    ALTER COLUMN user_id SET DEFAULT nextval('users_id_seq');

-- Vault-User table
ALTER TABLE vault_user_relations
    ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE vault_user_relations
    ALTER COLUMN user_id SET DATA TYPE integer USING (nextval('users_id_seq'));

ALTER TABLE vault_user_relations
    ALTER COLUMN user_id SET DEFAULT nextval('users_id_seq');

-- Clean up orphaned sessions before restoring foreign keys
DELETE FROM sessions
WHERE user_id IS NOT NULL
  AND user_id NOT IN (SELECT id FROM users);

-- Restore all foreign keys
ALTER TABLE tag_task_relations ADD CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id);
ALTER TABLE tag_task_relations ADD CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks (id);
ALTER TABLE tags ADD CONSTRAINT fk_tags_user FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE sessions ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE vault_user_relations ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id);

-- Clean up: (sessions_id_seq may not exist if not previously used, so ignore errors)
DO $$
BEGIN
    BEGIN
        CREATE SEQUENCE IF NOT EXISTS sessions_id_seq;
    EXCEPTION WHEN duplicate_table THEN
        -- do nothing
    END;
END$$;