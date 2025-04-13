-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, user_name)
VALUES(
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: DeleteUser :exec
TRUNCATE TABLE users;

-- name: GetUser :one
SELECT user_name FROM users WHERE user_name = $1;

-- name: GetUserByName :one
SELECT id FROM users WHERE user_name = $1;
