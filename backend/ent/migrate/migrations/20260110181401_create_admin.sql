-- +goose Up
INSERT INTO
	users (id, username, sessions_valid_from)
VALUES
	(lower(hex(randomblob(16))), "admin", DATE("now"));

-- +goose Down
DELETE FROM users
WHERE
	username = "admin"
LIMIT
	1;