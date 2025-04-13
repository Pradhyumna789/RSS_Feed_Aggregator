-- name: CreateFeed :one
INSERT INTO feeds(id, createdAt, updatedAt, feed_name, feed_url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


