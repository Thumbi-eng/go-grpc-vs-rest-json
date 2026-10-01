package repository

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"grpc-auth-service/internal/db"
	"grpc-auth-service/internal/model"
)

var (
	// ErrUserNotFound is returned when a requested user does not exist
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailAlreadyExists is returned when attempting to register with an email already taken
	ErrEmailAlreadyExists = errors.New("email is already registered")
)

// UserRepository defines the standard CRUD interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, page, limit int32) ([]*model.User, int32, error)
	Update(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id string) error
}

type memoryUserRepository struct {
	db *db.MemoryDB
}

// NewUserRepository creates a new UserRepository backed by our thread-safe MemoryDB
func NewUserRepository(database *db.MemoryDB) UserRepository {
	return &memoryUserRepository{
		db: database,
	}
}

// Create inserts a new user record into the store
func (r *memoryUserRepository) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	normalizedEmail := strings.ToLower(strings.TrimSpace(input.Email))

	// Validate email uniqueness
	for _, u := range r.db.Users {
		if strings.EqualFold(u.Email, normalizedEmail) {
			return nil, ErrEmailAlreadyExists
		}
	}

	newUser := &model.User{
		ID:        uuid.NewString(),
		Username:  strings.TrimSpace(input.Username),
		Email:     normalizedEmail,
		Password:  input.Password,
		CreatedAt: time.Now().UTC(),
	}

	r.db.Users[newUser.ID] = newUser
	return cloneUser(newUser), nil
}

// GetByID finds a single user by their unique UUID
func (r *memoryUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	user, exists := r.db.Users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return cloneUser(user), nil
}

// GetByEmail finds a single user by email address (case-insensitive)
func (r *memoryUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	for _, user := range r.db.Users {
		if strings.EqualFold(user.Email, normalizedEmail) {
			return cloneUser(user), nil
		}
	}

	return nil, ErrUserNotFound
}

// List returns a paginated list of users, sorted newest to oldest
func (r *memoryUserRepository) List(ctx context.Context, page, limit int32) ([]*model.User, int32, error) {
	r.db.Mu.RLock()
	defer r.db.Mu.RUnlock()

	// Collect all users into a slice
	all := make([]*model.User, 0, len(r.db.Users))
	for _, u := range r.db.Users {
		all = append(all, cloneUser(u))
	}

	// Sort users newest first (deterministic order)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := int32(len(all))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	startIndex := (page - 1) * limit
	if startIndex >= total {
		return []*model.User{}, total, nil
	}

	endIndex := startIndex + limit
	if endIndex > total {
		endIndex = total
	}

	return all[startIndex:endIndex], total, nil
}

// Update updates fields on an existing user record
func (r *memoryUserRepository) Update(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	existing, exists := r.db.Users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	// If email is changing, verify the new email is not already taken by someone else
	if input.Email != nil {
		newEmail := strings.ToLower(strings.TrimSpace(*input.Email))
		if !strings.EqualFold(existing.Email, newEmail) {
			for _, u := range r.db.Users {
				if u.ID != id && strings.EqualFold(u.Email, newEmail) {
					return nil, ErrEmailAlreadyExists
				}
			}
			existing.Email = newEmail
		}
	}

	if input.Username != nil && strings.TrimSpace(*input.Username) != "" {
		existing.Username = strings.TrimSpace(*input.Username)
	}

	if input.Password != nil && *input.Password != "" {
		existing.Password = *input.Password
	}

	return cloneUser(existing), nil
}

// Delete removes a user record from the store
func (r *memoryUserRepository) Delete(ctx context.Context, id string) error {
	r.db.Mu.Lock()
	defer r.db.Mu.Unlock()

	if _, exists := r.db.Users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.db.Users, id)
	return nil
}

// cloneUser creates a defensive copy of a User pointer to prevent concurrency leaks
func cloneUser(u *model.User) *model.User {
	if u == nil {
		return nil
	}
	copied := *u
	return &copied
}
