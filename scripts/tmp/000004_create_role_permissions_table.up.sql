CREATE TABLE IF NOT EXISTS role_permissions(
    role_id CHAR (36) NOT NULL,
    permission_id CHAR (36) NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
) ENGINE = INNODB;