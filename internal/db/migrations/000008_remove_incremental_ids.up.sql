-- Tag-Task table
ALTER TABLE tag_task_relations
DROP CONSTRAINT fk_tag;

ALTER TABLE tag_task_relations
DROP CONSTRAINT fk_task;

ALTER TABLE tag_task_relations
ALTER COLUMN tag_id
DROP DEFAULT,
ALTER COLUMN tag_id
SET
    DATA TYPE UUID USING (gen_random_uuid ()),
ALTER COLUMN tag_id
SET DEFAULT gen_random_uuid (),
ALTER COLUMN task_id
DROP DEFAULT,
ALTER COLUMN task_id
SET
    DATA TYPE UUID USING (gen_random_uuid ()),
ALTER COLUMN task_id
SET DEFAULT gen_random_uuid ();

-- Tags table
ALTER TABLE tags
DROP CONSTRAINT fk_tags_user;

ALTER TABLE tags
ALTER COLUMN id
DROP DEFAULT,
ALTER COLUMN id
SET
    DATA TYPE UUID USING (gen_random_uuid ()),
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();

-- Tasks table
ALTER TABLE tasks
DROP CONSTRAINT fk_tasks_user;

ALTER TABLE tasks
ALTER COLUMN id
DROP DEFAULT,
ALTER COLUMN id
SET
    DATA TYPE UUID USING (gen_random_uuid ()),
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();

-- User table 
ALTER TABLE sessions
DROP CONSTRAINT sessions_user_id_fkey;

ALTER TABLE vault_user_relations
DROP CONSTRAINT fk_user_id;

ALTER TABLE users
ALTER COLUMN id
DROP DEFAULT,
ALTER COLUMN id
SET
    DATA TYPE UUID USING (gen_random_uuid ()),
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();

-- Now convert all referencing user_id columns to UUID
ALTER TABLE tags
ALTER COLUMN user_id
DROP DEFAULT,
ALTER COLUMN user_id
SET
    DATA TYPE UUID USING (gen_random_uuid ());

ALTER TABLE tasks
ALTER COLUMN user_id
DROP DEFAULT,
ALTER COLUMN user_id
SET
    DATA TYPE UUID USING (gen_random_uuid ());

ALTER TABLE sessions
ALTER COLUMN user_id
DROP DEFAULT,
ALTER COLUMN user_id
SET
    DATA TYPE UUID USING (gen_random_uuid ());

ALTER TABLE vault_user_relations
ALTER COLUMN user_id
DROP DEFAULT,
ALTER COLUMN user_id
SET
    DATA TYPE UUID USING (gen_random_uuid ());

-- Now add all foreign keys back
ALTER TABLE tags ADD CONSTRAINT fk_user_tags FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE tasks ADD CONSTRAINT fk_user_tasks FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE sessions ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE vault_user_relations ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id);

-- Tag-Task id column
ALTER TABLE tag_task_relations
ALTER COLUMN id
DROP DEFAULT,
ALTER COLUMN id
SET
    DATA TYPE UUID USING (gen_random_uuid ()),
ALTER COLUMN id
SET DEFAULT gen_random_uuid ();

-- Add back tag_task_relations foreign keys
ALTER TABLE tag_task_relations ADD CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id);

ALTER TABLE tag_task_relations ADD CONSTRAINT fk_task FOREIGN KEY (task_id) REFERENCES tasks (id);

-- removing the sequence stuff
DROP SEQUENCE IF EXISTS tag_task_relations_id_seq;

DROP SEQUENCE IF EXISTS tags_id_seq;

DROP SEQUENCE IF EXISTS tasks_id_seq;

DROP SEQUENCE IF EXISTS users_id_seq;

DROP SEQUENCE IF EXISTS sessions_id_seq;

DROP SEQUENCE IF EXISTS users_id_seq1;