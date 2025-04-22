-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS(
    INSERT INTO feed_follow(id, createdAt, updatedAt, user_id, feed_id)
    VALUES($1, $2, $3, $4, $5)
    RETURNING *
)

SELECT
    inserted_feed_follow.*,
    feeds.feed_name AS feed_name,
    users.user_name AS user_name
    FROM inserted_feed_follow
    INNER JOIN users ON inserted_feed_follow.user_id = users.id 
    INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT 
    users.user_name, feeds.feed_name FROM feed_follow 
    INNER JOIN users ON feed_follow.user_id = users.id
    INNER JOIN feeds ON feed_follow.feed_id = feeds.id
    WHERE feed_follow.user_id = $1;

-- name: DeleteFeedFollowRecord :exec
DELETE FROM feed_follow
WHERE feed_follow.user_id = $1
AND feed_follow.feed_id = (SELECT feeds.id FROM feeds WHERE feeds.feed_url = $2);
