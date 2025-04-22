-- +goose Up
CREATE TABLE IF NOT EXISTS feeds (
    id UUID PRIMARY KEY,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP DEFAULT NOW() NOT NULL,
    feed_name VARCHAR(255) NOT NULL,
    feed_url VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_feeds_users_user_id FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS feeds;
