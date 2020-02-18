CREATE TABLE IF NOT EXISTS rolebindings(
    role_name VARCHAR (50),
    user_id CHAR (36),
    FOREIGN KEY (role_name) REFERENCES roles (name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, user_id)
) ENGINE = INNODB;