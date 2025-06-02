-- name: LinkArticleToCandidate :one
INSERT INTO candidate_articles (
  id,
  created_at,
  updated_at,
  candidate_id,
  article_id
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetCandidatesByArticle :many
SELECT 
  c.*,
  ca.created_at as discovery_date
FROM candidates c
JOIN candidate_articles ca ON c.id = ca.candidate_id
WHERE ca.article_id = ?
ORDER BY ca.created_at DESC;

-- name: GetArticlesByCandidate :many
SELECT 
  a.*,
  ca.created_at as discovery_date
FROM articles a
JOIN candidate_articles ca ON a.id = ca.article_id
WHERE ca.candidate_id = ?
ORDER BY ca.created_at DESC;

-- name: CheckArticleCandidateLinkExists :one
SELECT EXISTS (
  SELECT 1 FROM candidate_articles
  WHERE candidate_id = ? AND article_id = ?
) AS link_exists;

-- name: UnlinkArticleFromCandidate :exec
DELETE FROM candidate_articles
WHERE candidate_id = ? AND article_id = ?;

-- name: GetCandidateDiscoveryArticles :many
SELECT
  a.*,
  s.description as search_description,
  ca.created_at as discovery_date
FROM candidate_articles ca
JOIN articles a ON ca.article_id = a.id
JOIN searches s ON a.search_id = s.id
WHERE ca.candidate_id = ?
ORDER BY ca.created_at DESC;

