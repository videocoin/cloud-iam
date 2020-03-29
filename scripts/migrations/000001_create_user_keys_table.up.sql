CREATE TABLE IF NOT EXISTS user_keys(
    id CHAR (36) NOT NULL,
    user_id CHAR (36) NOT NULL,
    public_key_data blob NOT NULL,
    valid_after_time timestamp NULL DEFAULT NULL,
    valid_before_time timestamp NULL DEFAULT NULL,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY(id)
) ENGINE = INNODB;

CREATE TRIGGER max_keys
BEFORE INSERT ON user_keys
FOR EACH ROW
BEGIN
    DECLARE cnt INT;

    SELECT count(*) INTO cnt FROM user_keys;

    IF cnt = 3 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You can store only 3 records.';
    END IF;
END;