-- name: CreateArticle :one
INSERT INTO articles (
  id,
  fetched_at,
  article_url,
  article_title,
  article_summary,
  article_authors,
  search_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetArticleByID :one
SELECT * FROM articles
WHERE id = ?
LIMIT 1;

-- name: GetArticleByURL :one
SELECT * FROM articles
WHERE article_url = ?
LIMIT 1;

-- name: ListArticlesBySearch :many
SELECT * FROM articles
WHERE search_id = ?
ORDER BY fetched_at DESC
LIMIT ?
OFFSET ?;

-- name: ListRecentArticles :many
SELECT * FROM articles
ORDER BY fetched_at DESC
LIMIT ?;

-- name: ListArticlesByAuthor :many
SELECT * FROM articles
WHERE article_authors LIKE '%' || ? || '%'
ORDER BY fetched_at DESC
LIMIT ?
OFFSET ?;

-- name: CheckArticleExists :one
SELECT EXISTS (
  SELECT 1 FROM articles
  WHERE article_url = ?
) AS article_exists;

-- name: GetRecentArticlesWithSearchInfo :many
SELECT 
  a.*,
  s.description as search_description,
  s.arvix_url as search_url
FROM articles a
JOIN searches s ON a.search_id = s.id
ORDER BY a.fetched_at DESC
LIMIT ?;

-- name: CountArticlesBySearch :one
SELECT COUNT(*) as article_count
FROM articles
WHERE search_id = ?;

-- name: SearchArticlesByKeyword :many
SELECT *
FROM articles
WHERE 
  LOWER(article_title) LIKE LOWER('%' || ? || '%') OR
  LOWER(article_summary) LIKE LOWER('%' || ? || '%')
ORDER BY fetched_at DESC
LIMIT ?
OFFSET ?;

-- name: GetArticlesFetchedBetween :many
SELECT *
FROM articles
WHERE 
  search_id = ? AND
  fetched_at BETWEEN ? AND ?
ORDER BY fetched_at DESC;

-- name: DeleteArticle :exec
DELETE FROM articles
WHERE id = ?;

-- name: UpdateArticleDetails :one
UPDATE articles
SET
  article_title = ?,
  article_summary = ?,
  article_authors = ?
WHERE id = ?
RETURNING *;

