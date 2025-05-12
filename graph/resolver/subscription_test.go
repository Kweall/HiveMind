package resolver_test

import (
	"context"
	"hivemind/graph/model"
	"hivemind/graph/resolver"
	"hivemind/internal/storage/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	t.Run("successful subscription", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		res := resolver.NewResolver(mockStorage)

		postID := "post123"
		commentChan := res.Subscribe(postID)
		comment := &model.Comment{
			ID:        "comment123",
			PostID:    postID,
			Content:   "Test comment",
			Author:    "alice",
			CreatedAt: time.Now(),
		}
		t.Logf("Notifying subscribers with comment: %v", comment)

		res.NotifySubscribers(postID, comment)

		select {
		case receivedComment := <-commentChan:
			t.Logf("Received comment: %v", receivedComment)
			assert.Equal(t, comment, receivedComment, "Received comment should match the sent comment")
		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for comment")
		}
	})

	t.Run("no subscribers", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)

		res := resolver.NewResolver(mockStorage)

		commentChan := res.Subscribe("post123")

		select {
		case <-commentChan:
			t.Error("Received a comment unexpectedly")
		case <-time.After(1 * time.Second):
			t.Log("No comment received as expected")
		}
	})
}

func TestCommentAdded(t *testing.T) {
	t.Run("successful comment added", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)
		res := resolver.NewResolver(mockStorage)

		postID := "post123"
		commentChan, err := res.CommentAdded(context.Background(), postID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		comment := &model.Comment{
			ID:        "comment123",
			PostID:    postID,
			Content:   "Test comment",
			Author:    "alice",
			CreatedAt: time.Now(),
		}

		t.Logf("Notifying subscribers with comment: %v", comment)
		res.NotifySubscribers(postID, comment)

		select {
		case receivedComment := <-commentChan:
			t.Logf("Received comment: %v", receivedComment)
			assert.Equal(t, comment, receivedComment, "Received comment should match the sent comment")
		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for comment")
		}
	})

	t.Run("subscription cancellation", func(t *testing.T) {
		mockStorage := mocks.NewStorageMock(t)
		res := resolver.NewResolver(mockStorage)

		postID := "post123"
		ctx, cancel := context.WithCancel(context.Background())
		commentChan, err := res.CommentAdded(ctx, postID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cancel()

		time.Sleep(100 * time.Millisecond)

		comment := &model.Comment{
			ID:        "comment123",
			PostID:    postID,
			Content:   "Test comment",
			Author:    "alice",
			CreatedAt: time.Now(),
		}
		res.NotifySubscribers(postID, comment)

		select {
		case _, ok := <-commentChan:
			if ok {
				t.Error("Received comment after cancel")
			} else {
				t.Log("Channel closed correctly after cancel")
			}
		case <-time.After(500 * time.Millisecond):
			t.Log("No comment received as expected")
		}
	})
}
