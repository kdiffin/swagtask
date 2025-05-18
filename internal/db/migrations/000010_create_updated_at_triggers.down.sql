DROP TRIGGER IF EXISTS set_vaults_updated_at ON vaults;

DROP TRIGGER IF EXISTS set_users_updated_at ON users;

DROP TRIGGER IF EXISTS set_tasks_updated_at ON tasks;

DROP TRIGGER IF EXISTS set_tags_updated_at ON tags;

DROP FUNCTION IF EXISTS set_updated_at ();