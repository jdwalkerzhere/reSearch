-- name: CreateSearch :one
INSERT INTO searches (
  id, 
  created_at, 
  updated_at, 
  description, 
  arvix_url, 
  results_per_fetch,
  last_fetch_date
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetSearchByID :one
SELECT * FROM searches
WHERE id = ?
LIMIT 1;

-- name: ListAllSearches :many
SELECT * FROM searches
ORDER BY created_at DESC;

-- name: ListRecentSearches :many
SELECT * FROM searches
ORDER BY updated_at DESC
LIMIT ?;

-- name: UpdateSearch :one
UPDATE searches
SET 
  updated_at = ?,
  description = ?,
  arvix_url = ?,
  results_per_fetch = ?,
  last_fetch_date = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSearch :exec
DELETE FROM searches
WHERE id = ?;

-- name: GetLatestFetchTime :one
SELECT MAX(fetched_at) as latest_fetch
FROM articles
WHERE search_id = ?;

-- name: GetSearchWithStats :one
SELECT 
  s.*,
  COUNT(DISTINCT a.id) AS article_count,
  COUNT(DISTINCT cs.candidate_id) AS candidate_count,
  s.last_fetch_date,
  MIN(a.fetched_at) AS first_fetch_time
FROM searches s
LEFT JOIN articles a ON s.id = a.search_id
LEFT JOIN candidate_searches cs ON s.id = cs.search_id
WHERE s.id = ?
GROUP BY s.id;

-- name: ListActiveSearches :many
SELECT 
  s.*,
  COUNT(a.id) AS article_count
FROM searches s
LEFT JOIN articles a ON s.id = a.search_id
GROUP BY s.id
ORDER BY CASE WHEN s.last_fetch_date IS NULL THEN 0 ELSE 1 END DESC, s.last_fetch_date DESC, s.created_at DESC
LIMIT ?
OFFSET ?;

-- name: GetSearchesByArxivCategory :many
SELECT s.*
FROM searches s
WHERE s.arvix_url LIKE '%' || ? || '%'
ORDER BY s.created_at DESC;

-- name: SearchByDescription :many
SELECT *
FROM searches
WHERE LOWER(description) LIKE LOWER('%' || ? || '%')
ORDER BY created_at DESC;

-- name: GetSearchesWithoutRecentFetches :many
SELECT s.*
FROM searches s
WHERE s.last_fetch_date IS NULL OR s.last_fetch_date < ?
ORDER BY 
  CASE WHEN s.last_fetch_date IS NULL THEN 0 ELSE 1 END ASC,
  s.last_fetch_date ASC
LIMIT ?;

-- name: UpdateSearchFetchRate :one
UPDATE searches
SET
  updated_at = ?,
  results_per_fetch = ?
WHERE id = ?
RETURNING *;

-- name: IncrementSearchFetchResults :one
UPDATE searches
SET
  updated_at = ?,
  results_per_fetch = MIN(results_per_fetch + ?, 1999)
WHERE id = ?
RETURNING *;

-- name: UpdateSearchLastFetchDate :one
UPDATE searches
SET
  updated_at = ?,
  last_fetch_date = ?
WHERE id = ?
RETURNING *;

-- name: DecrementSearchFetchResults :one
UPDATE searches
SET
  updated_at = ?,
  results_per_fetch = MAX(results_per_fetch - ?, 10)
WHERE id = ?
RETURNING *;

