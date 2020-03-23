CREATE TABLE IF NOT EXISTS roles(
    id CHAR (36) NOT NULL,
    name VARCHAR (36) NOT NULL UNIQUE,
    title VARCHAR (36) NOT NULL DEFAULT "",
    description VARCHAR (100) NOT NULL DEFAULT "",
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY(id)
) ENGINE = INNODB;