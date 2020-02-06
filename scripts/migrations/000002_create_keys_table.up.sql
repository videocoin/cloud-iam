CREATE TABLE IF NOT EXISTS account_keys(
    id CHAR (36) NOT NULL,
    account_id CHAR (36) NOT NULL,
    private_key_data VARBINARY (3000) NOT NULL,
    public_key_data VARBINARY (3000) NOT NULL,
    valid_after_time timestamp NULL DEFAULT NULL,
    valid_before_time timestamp NULL DEFAULT NULL,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
) ENGINE = INNODB;