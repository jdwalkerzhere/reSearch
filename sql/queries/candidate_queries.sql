-- name: CreateCandidate :one
INSERT INTO candidates (
  id,
  created_at,
  updated_at,
  name,
  linkedin_url,
  github_url
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetCandidateWithCategories :many
SELECT 
  c.*,
  cc.arxiv_category
FROM candidates c
LEFT JOIN candidate_categories cc ON c.id = cc.candidate_id
WHERE c.id = ?;

-- name: ListCandidatesBySearch :many
SELECT 
  c.*,
  cs.relevance_score,
  cs.notes
FROM candidates c
JOIN candidate_searches cs ON c.id = cs.candidate_id
WHERE cs.search_id = ?
ORDER BY cs.relevance_score DESC
LIMIT ?
OFFSET ?;

-- name: UpdateCandidate :one
UPDATE candidates
SET 
  updated_at = ?,
  name = ?,
  linkedin_url = ?,
  github_url = ?
WHERE id = ?
RETURNING *;

-- name: GetCandidatesByCategory :many
SELECT 
  c.*
FROM candidates c
JOIN candidate_categories cc ON c.id = cc.candidate_id
WHERE cc.arxiv_category = ?
ORDER BY c.created_at DESC
LIMIT ?
OFFSET ?;

-- name: AddCandidateCategory :one
INSERT INTO candidate_categories (
  id,
  created_at,
  updated_at,
  candidate_id,
  arxiv_category
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: LinkCandidateToSearch :one
INSERT INTO candidate_searches (
  id,
  created_at,
  updated_at,
  candidate_id,
  search_id,
  relevance_score,
  notes
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateCandidateRelevance :one
UPDATE candidate_searches
SET
  updated_at = ?,
  relevance_score = ?,
  notes = ?
WHERE candidate_id = ? AND search_id = ?
RETURNING *;

-- name: CheckCandidateExists :one
SELECT EXISTS (
  SELECT 1 FROM candidates
  WHERE linkedin_url = ? OR github_url = ? OR name = ?
) AS candidate_exists;

-- name: DeleteCandidateCategory :exec
DELETE FROM candidate_categories
WHERE candidate_id = ? AND arxiv_category = ?;

-- name: FindCandidateByName :many
SELECT * FROM candidates
WHERE LOWER(name) LIKE LOWER('%' || ? || '%')
ORDER BY created_at DESC
LIMIT ?
OFFSET ?;

