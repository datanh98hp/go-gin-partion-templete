<<<<<<< HEAD
-- name: CreateUser :one
INSERT INTO users (user_email, user_fullname, user_password, user_age, user_status,user_level) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING * ;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE user_uuid = $1;

-- name: UpdateUserByUUID :one
UPDATE users 
SET 
user_fullname = COALESCE(sqlc.narg(user_fullname), user_fullname),
user_password =  COALESCE(sqlc.narg(user_password), user_password),
user_age = COALESCE(sqlc.narg(user_age), user_age),
user_status =COALESCE(sqlc.narg(user_status), user_status),
user_level = COALESCE(sqlc.narg(user_level), user_level)
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NULL RETURNING *;

-- name: SoftDeleteUserByUUID :one
UPDATE users 
SET 
user_deleted_at = now()
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NULL RETURNING *;

-- name: RestoreUsers :one
UPDATE users 
SET 
user_deleted_at = NULL
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NOT NULL RETURNING *;

-- name: TrashUsers :one
DELETE FROM users 
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NOT NULL RETURNING *;

-- name: GetUsersIdAsc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL 
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_id ASC
LIMIT $1 OFFSET $2
;

-- name: GetUsersIdDesc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL 
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%' 
    OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_id DESC
LIMIT $1 OFFSET $2
;

-- name: GetUsersCreatedAtAsc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
   sqlc.narg(search)::TEXT IS NULL 
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%' 
    OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_created_at ASC
LIMIT $1 OFFSET $2
;

-- name: GetUsersCreatedAtDesc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL OR
    sqlc.narg(search)::TEXT = '' OR
    user_email ILIKE '%' || sqlc.narg(search) || '%' 
   OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_created_at DESC
LIMIT $1 OFFSET $2
;


-- name: CountUsers :one
SELECT COUNT(*)
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL OR
    sqlc.narg(search)::TEXT = '' OR
    user_email ILIKE '%' || sqlc.narg(search) || '%' 
   OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
);
=======
-- name: CreateUser :one
INSERT INTO users (user_email, user_fullname, user_password, user_age, user_status,user_level) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING * ;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE user_uuid = $1 AND user_deleted_at IS NULL;

-- name: UpdateUserByUUID :one
UPDATE users 
SET 
user_fullname = COALESCE(sqlc.narg(user_fullname), user_fullname),
user_password =  COALESCE(sqlc.narg(user_password), user_password),
user_age = COALESCE(sqlc.narg(user_age), user_age),
user_status =COALESCE(sqlc.narg(user_status), user_status),
user_level = COALESCE(sqlc.narg(user_level), user_level)
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NULL RETURNING *;

-- name: GetUsersDeleted :many
SELECT * FROM users WHERE user_deleted_at IS NOT NULL;

-- name: SoftDeleteUserByUUID :one
UPDATE users 
SET 
user_deleted_at = now()
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NULL RETURNING *;

-- name: RestoreUsers :one
UPDATE users 
SET 
user_deleted_at = NULL
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NOT NULL RETURNING *;

-- name: TrashUsers :one
DELETE FROM users 
WHERE user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NOT NULL RETURNING *;

-- name: GetUsersIdAsc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL 
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_id ASC
LIMIT $1 OFFSET $2
;

-- name: GetUsersIdDesc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL 
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%' 
    OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_id DESC
LIMIT $1 OFFSET $2
;

-- name: GetUsersCreatedAtAsc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
   sqlc.narg(search)::TEXT IS NULL 
    OR sqlc.narg(search)::TEXT = ''
    OR user_email ILIKE '%' || sqlc.narg(search) || '%' 
    OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_created_at ASC
LIMIT $1 OFFSET $2
;

-- name: GetUsersCreatedAtDesc :many
SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
    sqlc.narg(search)::TEXT IS NULL OR
    sqlc.narg(search)::TEXT = '' OR
    user_email ILIKE '%' || sqlc.narg(search) || '%' 
   OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
)
ORDER BY user_created_at DESC
LIMIT $1 OFFSET $2
;


-- name: CountUsers :one
SELECT COUNT(*)
FROM users 
WHERE (
    sqlc.narg(deleted)::bool IS NULL 
    OR  (sqlc.narg(deleted)::bool = TRUE AND user_deleted_at IS NOT NULL )
    OR  (sqlc.narg(deleted)::bool = FALSE AND user_deleted_at IS NULL )
)
AND (
    sqlc.narg(search)::TEXT IS NULL OR
    sqlc.narg(search)::TEXT = '' OR
    user_email ILIKE '%' || sqlc.narg(search) || '%' 
   OR user_fullname ILIKE '%'|| sqlc.narg(search) || '%'
);
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114
