CREATE TABLE IF NOT EXISTS roles_permissions(
    role_name VARCHAR (50),
    permission_name VARCHAR (50),
    FOREIGN KEY (role_name) REFERENCES roles (name) ON DELETE CASCADE,
    FOREIGN KEY (permission_name) REFERENCES permissions (name) ON DELETE CASCADE,
    PRIMARY KEY (role_name, permission_name)
) ENGINE = INNODB;