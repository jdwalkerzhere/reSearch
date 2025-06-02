-- +goose Up
CREATE TABLE candidate_searches(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	candidate_id UUID NOT NULL,
	search_id UUID NOT NULL,
	relevance_score FLOAT,
	notes TEXT,
	FOREIGN KEY(candidate_id) REFERENCES candidates(id),
	FOREIGN KEY(search_id) REFERENCES searches(id),
	UNIQUE(candidate_id, search_id)
);

-- +goose Down
DROP TABLE candidate_searches;

