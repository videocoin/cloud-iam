CREATE TABLE IF NOT EXISTS roles(
    name VARCHAR (50),
    title VARCHAR (50) NOT NULL UNIQUE,
    description VARCHAR (100),
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY(name)
) ENGINE = INNODB;