CREATE TABLE IF NOT EXISTS accounts(
    id CHAR (36) PRIMARY KEY,
    project_id VARCHAR (50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL
) ENGINE = INNODB;