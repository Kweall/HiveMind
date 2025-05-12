package resolver_test

import (
	"context"
	"errors"
	"hivemind/graph/model"
	"hivemind/graph/resolver"
	"hivemind/internal/storage/mocks"
	"testing"
)

func TestCreateComment(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post123",
			CommentsEnabled: true,
		}, nil)

		mockStorage.CreateCommentMock.Return(nil)

		res := resolver.NewResolver(mockStorage)
		comment, err := res.CreateComment(ctx, "post123", nil, "Test comment", "textik")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if comment.PostID != "post123" || comment.Author != "textik" || comment.Content != "Test comment" {
			t.Errorf("unexpected comment values: %+v", comment)
		}
	})

	t.Run("post not found", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(nil, errors.New("not found"))

		res := resolver.NewResolver(mockStorage)
		_, err := res.CreateComment(ctx, "missing", nil, "Test comment", "bob")

		if err == nil || err.Error() != "not found" {
			t.Errorf("expected 'not found' error, got: %v", err)
		}
	})

	t.Run("comments disabled", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post456",
			CommentsEnabled: false,
		}, nil)

		res := resolver.NewResolver(mockStorage)
		_, err := res.CreateComment(ctx, "post456", nil, "Test comment", "bob")

		if err == nil || err.Error() != "commenting is disabled for this post" {
			t.Errorf("expected 'commenting is disabled' error, got: %v", err)
		}
	})

	t.Run("comment too long", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post789",
			CommentsEnabled: true,
		}, nil)

		longComment := make([]byte, 2001)
		for i := range longComment {
			longComment[i] = 'a'
		}

		res := resolver.NewResolver(mockStorage)
		_, err := res.CreateComment(ctx, "post789", nil, string(longComment), "bob")

		if err == nil || err.Error() != "comment too long" {
			t.Errorf("expected 'comment too long' error, got: %v", err)
		}
	})

	t.Run("CreateComment fails", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		mockStorage.GetPostByIDMock.Return(&model.Post{
			ID:              "post321",
			CommentsEnabled: true,
		}, nil)

		mockStorage.CreateCommentMock.Return(errors.New("db failure"))

		res := resolver.NewResolver(mockStorage)
		_, err := res.CreateComment(ctx, "post321", nil, "Test comment", "alice")

		if err == nil || err.Error() != "db failure" {
			t.Errorf("expected 'db failure' error, got: %v", err)
		}
	})
}
