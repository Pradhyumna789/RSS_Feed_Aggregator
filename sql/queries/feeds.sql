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

-- name: GetFeedByURL :one
SELECT id FROM feeds WHERE feed_url = $1;

-- name: GetFeedNameById :one
SELECT feed_name FROM feeds WHERE id = $1; 

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = NOW(), updatedat = NOW();

-- name: GetNextFeedToFetch :many
SELECT id, feed_name, feed_url, last_fetched_at
    FROM feeds
    ORDER BY last_fetched_at DESC NULLS LAST
    LIMIT 5; 
