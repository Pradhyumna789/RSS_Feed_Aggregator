-- +goose Up 
CREATE TABLE feed_follow(
    id UUID PRIMARY KEY,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP DEFAULT NOW() NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL, 

    CONSTRAINT fk_feed_follow_users_user_id FOREIGN KEY (user_id) 
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_feed_follow_feeds_feed_id FOREIGN KEY (feed_id) 
        REFERENCES feeds (id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS feed_follow;