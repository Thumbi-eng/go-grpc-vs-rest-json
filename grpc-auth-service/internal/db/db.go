package db

import (
	"sync"

	"grpc-auth-service/internal/model"
)

// MemoryDB provides a thread-safe in-memory storage store for our application
// It requires zero external database setup, making it ideal for tutorials and quick demonstrations
type MemoryDB struct {
	Mu    sync.RWMutex
	Users map[string]*model.User
}

// NewMemoryDB initializes a new in-memory database
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		Users: make(map[string]*model.User),
	}
}
