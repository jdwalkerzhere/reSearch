-- +goose Up
CREATE TABLE articles(
	id UUID PRIMARY KEY,
	fetched_at TIMESTAMP NOT NULL,
	article_url TEXT NOT NULL,
	article_title TEXT NOT NULL,
	article_summary TEXT NOT NULL,
	article_authors TEXT NOT NULL,
	search_id UUID NOT NULL,
	FOREIGN KEY(search_id) REFERENCES searches(id)
);

-- +goose Down
DROP TABLE articles;
