package resolver

import (
	"hivemind/graph/model"
	"hivemind/internal/storage"
	"sync"
)

type Resolver struct {
	Storage     storage.Storage
	subscribers map[string][]chan *model.Comment
	mu          sync.Mutex
}

func NewResolver(storage storage.Storage) *Resolver {
	return &Resolver{
		Storage:     storage,
		subscribers: make(map[string][]chan *model.Comment),
	}
}
