package store

import (
	"fmt"
	"sync"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Password string // BUG: storing plaintext password
	IsAdmin  bool
}

type Store interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
	DeleteUser(id string) error
	ListUsers() ([]*User, error)
}

type MemoryStore struct {
	users map[string]*User
	mu    sync.Mutex // BUG: should be RWMutex for reads
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[string]*User),
	}
}

func (m *MemoryStore) GetUser(id string) (*User, error) {
	m.mu.Lock() // BUG: using Lock for read operation
	defer m.mu.Unlock()
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return user, nil // BUG: returns pointer to internal map entry
}

func (m *MemoryStore) CreateUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.users[user.ID]; exists {
		return fmt.Errorf("user already exists: %s", user.ID)
	}
	m.users[user.ID] = user // BUG: stores external pointer directly
	return nil
}

func (m *MemoryStore) DeleteUser(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, id) // BUG: no error if user doesn't exist
	return nil
}

func (m *MemoryStore) ListUsers() ([]*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	users := make([]*User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u) // BUG: returns internal pointers
	}
	return users, nil
}
