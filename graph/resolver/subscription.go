package resolver

import (
	"context"
	"hivemind/graph/model"
)

func (r *Resolver) Subscribe(postID string) <-chan *model.Comment {
	ch := make(chan *model.Comment, 1)
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.subscribers == nil {
		r.subscribers = make(map[string][]chan *model.Comment)
	}

	r.subscribers[postID] = append(r.subscribers[postID], ch)
	return ch
}

func (r *Resolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	commentChan := r.Subscribe(postID)
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		defer r.mu.Unlock()
		subs := r.subscribers[postID]
		for i, ch := range subs {
			if ch == commentChan {
				r.subscribers[postID] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
	}()
	return commentChan, nil
}
