CREATE TABLE IF NOT EXISTS user_keys(
    id CHAR (36) NOT NULL,
    user_id CHAR (36) NOT NULL,
    public_key_data bytea NOT NULL,
    valid_after_time timestamp NULL DEFAULT NULL,
    valid_before_time timestamp NULL DEFAULT NULL,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY(id)
);