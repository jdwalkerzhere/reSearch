-- +goose Up
CREATE TABLE candidates(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	linkedin_url TEXT,
	github_url TEXT
);

-- +goose Down
DROP TABLE candidates;
