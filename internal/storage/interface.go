package storage

import (
	"context"

	"hivemind/graph/model"
)

//go:generate minimock -i hivemind/internal/storage.Storage -o ./mocks -s "_mock.go"

type Storage interface {
	// Post
	CreatePost(ctx context.Context, post *model.Post) error
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	ToggleComments(ctx context.Context, postID string, enabled bool, author string) (*model.Post, error)

	// Comment
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error)
	GetReplies(ctx context.Context, parentID string, limit, offset int) ([]*model.Comment, error)
}
