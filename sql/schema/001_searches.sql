-- +goose Up
CREATE TABLE searches(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	description TEXT NOT NULL,
	arvix_url TEXT NOT NULL, -- this will be something like: "http://rss.arxiv.org/rss/cs.LG+cs.PL"
	results_per_fetch INTEGER,
	last_fetch_date TIMESTAMP,
	CONSTRAINT results_per_fetch_limit CHECK (results_per_fetch < 2000)
);

-- +goose Down
DROP TABLE searches;
