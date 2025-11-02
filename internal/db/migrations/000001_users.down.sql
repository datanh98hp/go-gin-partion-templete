<<<<<<< HEAD
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
=======
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
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114
drop table if exists users cascade;