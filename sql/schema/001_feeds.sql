-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    createdAt TIMESTAMP NOT NULL, 
    updatedAt TIMESTAMP NOT NULL,
    feed_name VARCHAR(255),
    feed_url VARCHAR(255) UNIQUE,
    user_id  VARCHAR(255),

    CONSTRAINT fk_feeds_users_user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id)
                ON DELETE CASCADE -- if a user is deleted all of their feed will also be deleted by the database server automatically
);

-- +goose Down 
DROP TABLE feeds;
