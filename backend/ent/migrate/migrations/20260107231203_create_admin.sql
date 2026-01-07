-- +goose Up
INSERT INTO
	users (username, sessions_valid_from)
VALUES
	("admin", DATE("now"));
-- +goose Down
DELETE FROM users
WHERE
	username = "admin"
LIMIT
	1;