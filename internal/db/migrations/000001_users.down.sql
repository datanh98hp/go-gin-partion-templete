-- drop triggers first
DROP TRIGGER IF EXISTS set_update_user_updated_at ON users;

-- drop trigger functions
DROP FUNCTION IF EXISTS update_user_updated_at_collumn ();

-- drop indexes (optional)
DROP INDEX IF EXISTS user_level_idx;

DROP INDEX IF EXISTS user_status_idx;

DROP INDEX IF EXISTS user_created_at_idx;

DROP INDEX IF EXISTS user_deleted_at_idx;

DROP INDEX IF EXISTS user_email_status_idx;

-- drop table
drop table if exists users cascade;