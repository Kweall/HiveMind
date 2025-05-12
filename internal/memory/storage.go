package memory

import (
	"context"
	"errors"
	"sync"

	"hivemind/graph/model"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	posts    map[string]*model.Post
	comments map[string]*model.Comment
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    make(map[string]*model.Post),
		comments: make(map[string]*model.Comment),
	}
}

func (m *MemoryStorage) CreatePost(ctx context.Context, post *model.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.posts[post.ID] = post
	return nil
}

func (m *MemoryStorage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	posts := []*model.Post{}
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	return posts, nil
}

func (m *MemoryStorage) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	post, ok := m.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *MemoryStorage) ToggleComments(ctx context.Context, postID string, enabled bool, author string) (*model.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	post, ok := m.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	post.CommentsEnabled = enabled
	return post, nil
}

func (m *MemoryStorage) CreateComment(ctx context.Context, comment *model.Comment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.posts[comment.PostID]; !ok {
		return errors.New("post not found")
	}
	m.comments[comment.ID] = comment
	if comment.ParentID != nil {
		parent, ok := m.comments[*comment.ParentID]
		if !ok {
			return errors.New("parent comment not found")
		}
		parent.Replies = append(parent.Replies, comment)
	} else {
		post := m.posts[comment.PostID]
		post.Comments = append(post.Comments, comment)
	}
	return nil
}

func (m *MemoryStorage) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	post, ok := m.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	comments := post.Comments
	if offset > len(comments) {
		return []*model.Comment{}, nil
	}
	end := offset + limit
	if end > len(comments) {
		end = len(comments)
	}
	return comments[offset:end], nil
}

func (m *MemoryStorage) GetReplies(ctx context.Context, parentID string, limit, offset int) ([]*model.Comment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	parent, ok := m.comments[parentID]
	if !ok {
		return nil, errors.New("parent comment not found")
	}
	replies := parent.Replies
	if offset > len(replies) {
		return []*model.Comment{}, nil
	}
	end := offset + limit
	if end > len(replies) {
		end = len(replies)
	}
	return replies[offset:end], nil
}
