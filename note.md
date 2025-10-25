-- tao file migration
migrate create -ext sql -dir internal/db/migrations -seq users

-- run migration
migrate -path internal/db/migrations -database "postgresql://root:06021998@localhost:5432/postgres?sslmode=disable" up

migrate -path internal/db/migrations -database "postgresql://root:06021998@localhost:5432/postgres?sslmode=disable" up 1


------------------query filter data----------------
-- SELECT * 
-- FROM users 
-- WHERE user_deleted_at IS NULL 
-- AND (
--     sqlc.narg(search)::TEXT IS NULL OR
--     sqlc.narg(search)::TEXT = '' OR
--     user_email ILIKE '%' || sqlc.narg(search) || '%' OR
--     user_fullname ILIKE '%' sqlc.narg(search) || '%'
-- )
-- ORDER BY 
--     CASE
--         WHEN  sqlc.narg(order_by) = 'user_id' AND sqlc.narg(sort) = 'asc' THEN user_id,
--     END DESC,
--     CASE
--         WHEN  sqlc.narg(order_by) = 'user_id' AND sqlc.narg(sort) = 'desc' THEN user_id,
--     END DESC,
--     CASE
--         WHEN  sqlc.narg(order_by) = 'user_created_at' AND sqlc.narg(sort) = 'desc' THEN user_id,
--     END DESC,
--     CASE
--         WHEN  sqlc.narg(order_by) = 'user_created_at' AND sqlc.narg(sort) = 'asc' THEN user_id,
--     END DESC,

-- LIMIT $1 OFFSET $2
-- ;