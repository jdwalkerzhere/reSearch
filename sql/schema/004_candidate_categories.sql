-- +goose Up
CREATE TABLE candidate_categories(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	candidate_id UUID NOT NULL,
	arxiv_category TEXT NOT NULL,
	FOREIGN KEY(candidate_id) REFERENCES candidates(id),
	UNIQUE (candidate_id, arxiv_category)
);

-- +goose Down
DROP TABLE candidate_categories;

