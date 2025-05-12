package resolver_test

import (
	"context"
	"errors"
	"hivemind/graph/model"
	"hivemind/graph/resolver"
	"hivemind/internal/storage/mocks"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.CreatePostMock.Return(nil)

		res := resolver.NewResolver(mockStorage)
		post, err := res.CreatePost(ctx, "Test Title", "Test Content", "alice")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if post.Title != "Test Title" || post.Content != "Test Content" || post.Author != "alice" {
			t.Errorf("unexpected post values: %+v", post)
		}
	})

	t.Run("CreatePost fails", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.CreatePostMock.Return(errors.New("db failure"))

		res := resolver.NewResolver(mockStorage)
		_, err := res.CreatePost(ctx, "Test Title", "Test Content", "alice")

		if err == nil || err.Error() != "db failure" {
			t.Errorf("expected 'db failure' error, got: %v", err)
		}
	})
}

func TestToggleComments(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post123",
			Author:          "alice",
			CommentsEnabled: true,
		}, nil)
		mockStorage.ToggleCommentsMock.Return(&model.Post{
			ID:              "post123",
			Author:          "alice",
			CommentsEnabled: false,
		}, nil)

		res := resolver.NewResolver(mockStorage)
		post, err := res.ToggleComments(ctx, "post123", false, "alice")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if post.CommentsEnabled != false {
			t.Errorf("expected comments to be disabled, got: %v", post.CommentsEnabled)
		}
	})

	t.Run("post not found", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(nil, errors.New("not found"))

		res := resolver.NewResolver(mockStorage)
		_, err := res.ToggleComments(ctx, "missing", false, "alice")

		if err == nil || err.Error() != "not found" {
			t.Errorf("expected 'not found' error, got: %v", err)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post123",
			Author:          "alice",
			CommentsEnabled: true,
		}, nil)

		res := resolver.NewResolver(mockStorage)
		_, err := res.ToggleComments(ctx, "post123", false, "bob")

		if err == nil || err.Error() != "only the author of the post can toggle comments" {
			t.Errorf("expected 'unauthorized' error, got: %v", err)
		}
	})
}

func TestPosts(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostsMock.Return([]*model.Post{
			{
				ID:              "post123",
				Title:           "Test Title",
				Content:         "Test Content",
				Author:          "alice",
				CommentsEnabled: true,
				CreatedAt:       time.Now(),
			},
		}, nil)

		res := resolver.NewResolver(mockStorage)
		posts, err := res.Posts(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(posts) != 1 || posts[0].Title != "Test Title" {
			t.Errorf("unexpected posts: %+v", posts)
		}
	})

	t.Run("GetPosts fails", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostsMock.Return(nil, errors.New("db failure"))

		res := resolver.NewResolver(mockStorage)
		_, err := res.Posts(ctx)

		if err == nil || err.Error() != "db failure" {
			t.Errorf("expected 'db failure' error, got: %v", err)
		}
	})
}

func TestPost(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post123",
			Title:           "Test Title",
			Content:         "Test Content",
			Author:          "alice",
			CommentsEnabled: true,
			CreatedAt:       time.Now(),
		}, nil)

		res := resolver.NewResolver(mockStorage)
		post, err := res.Post(ctx, "post123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if post.ID != "post123" || post.Title != "Test Title" {
			t.Errorf("unexpected post values: %+v", post)
		}
	})

	t.Run("post not found", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(nil, errors.New("not found"))

		res := resolver.NewResolver(mockStorage)
		_, err := res.Post(ctx, "missing")

		if err == nil || err.Error() != "not found" {
			t.Errorf("expected 'not found' error, got: %v", err)
		}
	})
}
