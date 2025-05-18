-- TODO: fix this broken thing
-- 1. Drop foreign keys first
ALTER TABLE tag_task_relations DROP CONSTRAINT IF EXISTS fk_tag;
ALTER TABLE tag_task_relations DROP CONSTRAINT IF EXISTS fk_task;
ALTER TABLE tags DROP CONSTRAINT IF EXISTS fk_user_tags;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_user_tasks;
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS fk_sessions_user;
ALTER TABLE vault_user_relations DROP CONSTRAINT IF EXISTS fk_user_id;

-- 2. Recreate sequences before conversion
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS tags_id_seq;
CREATE SEQUENCE IF NOT EXISTS tasks_id_seq;
CREATE SEQUENCE IF NOT EXISTS tag_task_relations_id_seq;
CREATE SEQUENCE IF NOT EXISTS sessions_id_seq;

-- 3. Convert tables from UUID back to integer
ALTER TABLE users
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id SET DATA TYPE integer USING (id::text)::integer;

ALTER TABLE tags
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id SET DATA TYPE integer USING (id::text)::integer,
    ALTER COLUMN user_id DROP DEFAULT,
    ALTER COLUMN user_id SET DATA TYPE integer USING (user_id::text)::integer;

ALTER TABLE tasks
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id SET DATA TYPE integer USING (id::text)::integer,
    ALTER COLUMN user_id DROP DEFAULT,
    ALTER COLUMN user_id SET DATA TYPE integer USING (user_id::text)::integer;

ALTER TABLE sessions
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id SET DATA TYPE integer USING (id::text)::integer,
    ALTER COLUMN user_id DROP DEFAULT,
    ALTER COLUMN user_id SET DATA TYPE integer USING (user_id::text)::integer;

ALTER TABLE vault_user_relations
    ALTER COLUMN user_id DROP DEFAULT,
    ALTER COLUMN user_id SET DATA TYPE integer USING (user_id::text)::integer;

ALTER TABLE tag_task_relations
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id SET DATA TYPE integer USING (id::text)::integer,
    ALTER COLUMN tag_id DROP DEFAULT,
    ALTER COLUMN tag_id SET DATA TYPE integer USING (tag_id::text)::integer,
    ALTER COLUMN task_id DROP DEFAULT,
    ALTER COLUMN task_id SET DATA TYPE integer USING (task_id::text)::integer;

-- 4. Set default values to use sequences
ALTER TABLE users ALTER COLUMN id SET DEFAULT nextval('users_id_seq');
ALTER TABLE tags ALTER COLUMN id SET DEFAULT nextval('tags_id_seq');
ALTER TABLE tasks ALTER COLUMN id SET DEFAULT nextval('tasks_id_seq');
ALTER TABLE tag_task_relations ALTER COLUMN id SET DEFAULT nextval('tag_task_relations_id_seq');
ALTER TABLE sessions ALTER COLUMN id SET DEFAULT nextval('sessions_id_seq');

-- 5. Set sequence current values with explicit cast to integer to avoid type mismatch
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users)::integer, 1));
SELECT setval('tags_id_seq', COALESCE((SELECT MAX(id) FROM tags)::integer, 1));
SELECT setval('tasks_id_seq', COALESCE((SELECT MAX(id) FROM tasks)::integer, 1));
SELECT setval('tag_task_relations_id_seq', COALESCE((SELECT MAX(id) FROM tag_task_relations)::integer, 1));
SELECT setval('sessions_id_seq', COALESCE((SELECT MAX(id) FROM sessions)::integer, 1));

-- 6. Restore foreign key constraints using the original names
ALTER TABLE tags ADD CONSTRAINT fk_tags_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE tag_task_relations ADD CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags(id);
ALTER TABLE tag_task_relations ADD CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks(id);
ALTER TABLE sessions ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE vault_user_relations ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);