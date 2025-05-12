package resolver

import (
	"context"
	"errors"
	"hivemind/graph/model"
	"time"
)

func (r *Resolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.Storage.GetPosts(ctx)
}

func (r *Resolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return r.Storage.GetPostByID(ctx, id)
}

func (r *Resolver) CreatePost(ctx context.Context, title, content, author string) (*model.Post, error) {
	post := &model.Post{
		ID:              GenerateID(),
		Title:           title,
		Content:         content,
		Author:          author,
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
		Comments:        []*model.Comment{},
	}
	if err := r.Storage.CreatePost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *Resolver) ToggleComments(ctx context.Context, postID string, enabled bool, author string) (*model.Post, error) {
	post, err := r.Storage.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.Author != author {
		return nil, errors.New("only the author of the post can toggle comments")
	}

	return r.Storage.ToggleComments(ctx, postID, enabled, author)
}
