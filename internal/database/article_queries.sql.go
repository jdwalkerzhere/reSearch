// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: article_queries.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const checkArticleExists = `-- name: CheckArticleExists :one
SELECT EXISTS (
  SELECT 1 FROM articles
  WHERE article_url = ?
) AS article_exists
`

func (q *Queries) CheckArticleExists(ctx context.Context, articleUrl string) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkArticleExists, articleUrl)
	var article_exists int64
	err := row.Scan(&article_exists)
	return article_exists, err
}

const countArticlesBySearch = `-- name: CountArticlesBySearch :one
SELECT COUNT(*) as article_count
FROM articles
WHERE search_id = ?
`

func (q *Queries) CountArticlesBySearch(ctx context.Context, searchID interface{}) (int64, error) {
	row := q.db.QueryRowContext(ctx, countArticlesBySearch, searchID)
	var article_count int64
	err := row.Scan(&article_count)
	return article_count, err
}

const createArticle = `-- name: CreateArticle :one
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
RETURNING id, fetched_at, article_url, article_title, article_summary, article_authors, search_id
`

type CreateArticleParams struct {
	ID             interface{}
	FetchedAt      time.Time
	ArticleUrl     string
	ArticleTitle   string
	ArticleSummary string
	ArticleAuthors string
	SearchID       interface{}
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error) {
	row := q.db.QueryRowContext(ctx, createArticle,
		arg.ID,
		arg.FetchedAt,
		arg.ArticleUrl,
		arg.ArticleTitle,
		arg.ArticleSummary,
		arg.ArticleAuthors,
		arg.SearchID,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.FetchedAt,
		&i.ArticleUrl,
		&i.ArticleTitle,
		&i.ArticleSummary,
		&i.ArticleAuthors,
		&i.SearchID,
	)
	return i, err
}

const deleteArticle = `-- name: DeleteArticle :exec
DELETE FROM articles
WHERE id = ?
`

func (q *Queries) DeleteArticle(ctx context.Context, id interface{}) error {
	_, err := q.db.ExecContext(ctx, deleteArticle, id)
	return err
}

const getArticleByID = `-- name: GetArticleByID :one
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id FROM articles
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetArticleByID(ctx context.Context, id interface{}) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticleByID, id)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.FetchedAt,
		&i.ArticleUrl,
		&i.ArticleTitle,
		&i.ArticleSummary,
		&i.ArticleAuthors,
		&i.SearchID,
	)
	return i, err
}

const getArticleByURL = `-- name: GetArticleByURL :one
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id FROM articles
WHERE article_url = ?
LIMIT 1
`

func (q *Queries) GetArticleByURL(ctx context.Context, articleUrl string) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticleByURL, articleUrl)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.FetchedAt,
		&i.ArticleUrl,
		&i.ArticleTitle,
		&i.ArticleSummary,
		&i.ArticleAuthors,
		&i.SearchID,
	)
	return i, err
}

const getArticlesFetchedBetween = `-- name: GetArticlesFetchedBetween :many
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id
FROM articles
WHERE 
  search_id = ? AND
  fetched_at BETWEEN ? AND ?
ORDER BY fetched_at DESC
`

func (q *Queries) GetArticlesFetchedBetween(ctx context.Context, searchID interface{}) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, getArticlesFetchedBetween, searchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.FetchedAt,
			&i.ArticleUrl,
			&i.ArticleTitle,
			&i.ArticleSummary,
			&i.ArticleAuthors,
			&i.SearchID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecentArticlesWithSearchInfo = `-- name: GetRecentArticlesWithSearchInfo :many
SELECT 
  a.id, a.fetched_at, a.article_url, a.article_title, a.article_summary, a.article_authors, a.search_id,
  s.description as search_description,
  s.arvix_url as search_url
FROM articles a
JOIN searches s ON a.search_id = s.id
ORDER BY a.fetched_at DESC
LIMIT ?
`

type GetRecentArticlesWithSearchInfoRow struct {
	ID                interface{}
	FetchedAt         time.Time
	ArticleUrl        string
	ArticleTitle      string
	ArticleSummary    string
	ArticleAuthors    string
	SearchID          interface{}
	SearchDescription string
	SearchUrl         string
}

func (q *Queries) GetRecentArticlesWithSearchInfo(ctx context.Context, limit int64) ([]GetRecentArticlesWithSearchInfoRow, error) {
	rows, err := q.db.QueryContext(ctx, getRecentArticlesWithSearchInfo, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRecentArticlesWithSearchInfoRow
	for rows.Next() {
		var i GetRecentArticlesWithSearchInfoRow
		if err := rows.Scan(
			&i.ID,
			&i.FetchedAt,
			&i.ArticleUrl,
			&i.ArticleTitle,
			&i.ArticleSummary,
			&i.ArticleAuthors,
			&i.SearchID,
			&i.SearchDescription,
			&i.SearchUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listArticlesByAuthor = `-- name: ListArticlesByAuthor :many
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id FROM articles
WHERE article_authors LIKE '%' || ? || '%'
ORDER BY fetched_at DESC
LIMIT ?
OFFSET ?
`

type ListArticlesByAuthorParams struct {
	Column1 sql.NullString
	Limit   int64
	Offset  int64
}

func (q *Queries) ListArticlesByAuthor(ctx context.Context, arg ListArticlesByAuthorParams) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listArticlesByAuthor, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.FetchedAt,
			&i.ArticleUrl,
			&i.ArticleTitle,
			&i.ArticleSummary,
			&i.ArticleAuthors,
			&i.SearchID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listArticlesBySearch = `-- name: ListArticlesBySearch :many
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id FROM articles
WHERE search_id = ?
ORDER BY fetched_at DESC
LIMIT ?
OFFSET ?
`

type ListArticlesBySearchParams struct {
	SearchID interface{}
	Limit    int64
	Offset   int64
}

func (q *Queries) ListArticlesBySearch(ctx context.Context, arg ListArticlesBySearchParams) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listArticlesBySearch, arg.SearchID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.FetchedAt,
			&i.ArticleUrl,
			&i.ArticleTitle,
			&i.ArticleSummary,
			&i.ArticleAuthors,
			&i.SearchID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecentArticles = `-- name: ListRecentArticles :many
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id FROM articles
ORDER BY fetched_at DESC
LIMIT ?
`

func (q *Queries) ListRecentArticles(ctx context.Context, limit int64) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listRecentArticles, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.FetchedAt,
			&i.ArticleUrl,
			&i.ArticleTitle,
			&i.ArticleSummary,
			&i.ArticleAuthors,
			&i.SearchID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchArticlesByKeyword = `-- name: SearchArticlesByKeyword :many
SELECT id, fetched_at, article_url, article_title, article_summary, article_authors, search_id
FROM articles
WHERE 
  LOWER(article_title) LIKE LOWER('%' || ? || '%') OR
  LOWER(article_summary) LIKE LOWER('%' || ? || '%')
ORDER BY fetched_at DESC
LIMIT ?
OFFSET ?
`

type SearchArticlesByKeywordParams struct {
	Column1 sql.NullString
	Column2 sql.NullString
	Limit   int64
	Offset  int64
}

func (q *Queries) SearchArticlesByKeyword(ctx context.Context, arg SearchArticlesByKeywordParams) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, searchArticlesByKeyword,
		arg.Column1,
		arg.Column2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.FetchedAt,
			&i.ArticleUrl,
			&i.ArticleTitle,
			&i.ArticleSummary,
			&i.ArticleAuthors,
			&i.SearchID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateArticleDetails = `-- name: UpdateArticleDetails :one
UPDATE articles
SET
  article_title = ?,
  article_summary = ?,
  article_authors = ?
WHERE id = ?
RETURNING id, fetched_at, article_url, article_title, article_summary, article_authors, search_id
`

type UpdateArticleDetailsParams struct {
	ArticleTitle   string
	ArticleSummary string
	ArticleAuthors string
	ID             interface{}
}

func (q *Queries) UpdateArticleDetails(ctx context.Context, arg UpdateArticleDetailsParams) (Article, error) {
	row := q.db.QueryRowContext(ctx, updateArticleDetails,
		arg.ArticleTitle,
		arg.ArticleSummary,
		arg.ArticleAuthors,
		arg.ID,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.FetchedAt,
		&i.ArticleUrl,
		&i.ArticleTitle,
		&i.ArticleSummary,
		&i.ArticleAuthors,
		&i.SearchID,
	)
	return i, err
}
