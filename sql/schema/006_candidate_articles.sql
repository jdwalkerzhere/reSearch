-- +goose Up
CREATE TABLE candidate_articles(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	candidate_id UUID NOT NULL,
	article_id UUID NOT NULL,
	FOREIGN KEY(candidate_id) REFERENCES candidates(id),
	FOREIGN KEY(article_id) REFERENCES articles(id),
	UNIQUE(candidate_id, article_id)
);

-- +goose Down
DROP TABLE candidate_articles;

