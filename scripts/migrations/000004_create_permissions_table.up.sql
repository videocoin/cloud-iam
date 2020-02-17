CREATE TABLE IF NOT EXISTS permissions(
    name VARCHAR (50),
    description VARCHAR (100),
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY(name)
) ENGINE = INNODB;