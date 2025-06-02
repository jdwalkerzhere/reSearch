-- name: GetTopSearches :many
SELECT 
  s.*,
  COUNT(DISTINCT cs.candidate_id) AS candidate_count,
  AVG(cs.relevance_score) AS avg_relevance
FROM searches s
LEFT JOIN candidate_searches cs ON s.id = cs.search_id
GROUP BY s.id
ORDER BY candidate_count DESC, avg_relevance DESC
LIMIT ?;

-- name: GetCandidateCategoryDistribution :many
SELECT 
  cc.arxiv_category,
  COUNT(DISTINCT cc.candidate_id) AS candidate_count
FROM candidate_categories cc
GROUP BY cc.arxiv_category
ORDER BY candidate_count DESC;

-- name: GetSearchStatistics :one
SELECT 
  s.id,
  s.description,
  s.arvix_url,
  COUNT(DISTINCT a.id) AS total_articles,
  COUNT(DISTINCT cs.candidate_id) AS total_candidates,
  MIN(a.fetched_at) AS first_fetch_time,
  s.last_fetch_date,
  AVG(cs.relevance_score) AS avg_relevance_score
FROM searches s
LEFT JOIN articles a ON s.id = a.search_id
LEFT JOIN candidate_searches cs ON s.id = cs.search_id
WHERE s.id = ?
GROUP BY s.id, s.description, s.arvix_url, s.last_fetch_date;

-- name: GetCandidateGrowthOverTime :many
SELECT 
  strftime('%Y-%m-%d', c.created_at) AS day,
  COUNT(DISTINCT c.id) AS new_candidates
FROM candidates c
JOIN candidate_searches cs ON c.id = cs.candidate_id
WHERE cs.search_id = ?
GROUP BY day
ORDER BY day;

-- name: GetSearchCoverageByCategory :many
SELECT 
  cc.arxiv_category,
  COUNT(DISTINCT cs.candidate_id) AS candidate_count
FROM candidate_categories cc
JOIN candidate_searches cs ON cc.candidate_id = cs.candidate_id
WHERE cs.search_id = ?
GROUP BY cc.arxiv_category
ORDER BY candidate_count DESC;

-- name: GetCandidateDiscoverySource :many
SELECT 
  c.name as candidate_name,
  a.article_title,
  a.article_url,
  ca.created_at as discovery_date,
  s.description as search_description
FROM candidate_articles ca
JOIN candidates c ON ca.candidate_id = c.id
JOIN articles a ON ca.article_id = a.id
JOIN searches s ON a.search_id = s.id
WHERE ca.candidate_id = ?
ORDER BY ca.created_at DESC;

-- name: GetMostProductiveSearches :many
SELECT 
  s.description,
  s.id,
  COUNT(DISTINCT ca.candidate_id) as discovered_candidates,
  COUNT(DISTINCT a.id) as processed_articles,
  CAST(COUNT(DISTINCT ca.candidate_id) AS FLOAT) / NULLIF(COUNT(DISTINCT a.id), 0) as discovery_rate
FROM searches s
JOIN articles a ON s.id = a.search_id
LEFT JOIN candidate_articles ca ON a.id = ca.article_id
GROUP BY s.id, s.description
HAVING COUNT(DISTINCT a.id) > 0
ORDER BY discovery_rate DESC
LIMIT ?;

-- name: GetRecentCandidateDiscoveries :many
SELECT
  c.id as candidate_id,
  c.name as candidate_name,
  a.article_title,
  ca.created_at as discovery_date,
  s.description as search_description
FROM candidate_articles ca
JOIN candidates c ON ca.candidate_id = c.id
JOIN articles a ON ca.article_id = a.id
JOIN searches s ON a.search_id = s.id
ORDER BY ca.created_at DESC
LIMIT ?;

-- name: GetSearchFetchHistory :many
SELECT
  s.id as search_id,
  s.description as search_description,
  a.fetched_at,
  COUNT(a.id) as articles_fetched,
  COUNT(DISTINCT ca.candidate_id) as candidates_discovered
FROM searches s
JOIN articles a ON s.id = a.search_id
LEFT JOIN candidate_articles ca ON a.id = ca.article_id
WHERE s.id = ?
GROUP BY s.id, s.description, DATE(a.fetched_at)
ORDER BY a.fetched_at DESC;

