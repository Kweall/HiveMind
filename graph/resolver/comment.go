package resolver

import (
	"context"
	"errors"
	"hivemind/graph/model"
	"time"
)

func (r *Resolver) CreateComment(ctx context.Context, postID string, parentID *string, content, author string) (*model.Comment, error) {
	post, err := r.Storage.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if !post.CommentsEnabled {
		return nil, errors.New("commenting is disabled for this post")
	}

	if len(content) > 2000 {
		return nil, errors.New("comment too long")
	}

	comment := &model.Comment{
		ID:        GenerateID(),
		PostID:    postID,
		ParentID:  parentID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now(),
		Replies:   []*model.Comment{},
	}

	if err := r.Storage.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	r.NotifySubscribers(postID, comment)

	return comment, nil
}

func (r *Resolver) NotifySubscribers(postID string, comment *model.Comment) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ch := range r.subscribers[postID] {
		select {
		case ch <- comment:
		default:
		}
	}
}
