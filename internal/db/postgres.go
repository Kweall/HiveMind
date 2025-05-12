package db

import (
	"context"
	"database/sql"
	"errors"

	"hivemind/graph/model"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) CreatePost(ctx context.Context, post *model.Post) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO posts (id, title, content, author, comments_enabled, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		post.ID, post.Title, post.Content, post.Author, post.CommentsEnabled, post.CreatedAt)
	return err
}

func (p *PostgresStorage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT id, title, content, author, comments_enabled, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CommentsEnabled, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (p *PostgresStorage) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	row := p.db.QueryRowContext(ctx, `SELECT id, title, content, author, comments_enabled, created_at FROM posts WHERE id = $1`, id)
	var post model.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CommentsEnabled, &post.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (p *PostgresStorage) ToggleComments(ctx context.Context, postID string, enabled bool, author string) (*model.Post, error) {
	_, err := p.db.ExecContext(ctx, `UPDATE posts SET comments_enabled = $1 WHERE id = $2`, enabled, postID)
	if err != nil {
		return nil, err
	}
	return p.GetPostByID(ctx, postID)
}

func (p *PostgresStorage) CreateComment(ctx context.Context, c *model.Comment) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO comments (id, post_id, parent_id, author, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.PostID, c.ParentID, c.Author, c.Content, c.CreatedAt)
	return err
}

func (p *PostgresStorage) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT id, post_id, parent_id, author, content, created_at FROM comments WHERE post_id = $1 AND parent_id IS NULL ORDER BY created_at ASC LIMIT $2 OFFSET $3`, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Author, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (p *PostgresStorage) GetReplies(ctx context.Context, parentID string, limit, offset int) ([]*model.Comment, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT id, post_id, parent_id, author, content, created_at FROM comments WHERE parent_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Author, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		replies = append(replies, &c)
	}
	return replies, nil
}
