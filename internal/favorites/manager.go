// ABOUTME: Business logic for managing favorites (add, remove, list, get).
// ABOUTME: Handles ID generation and persists changes to disk.

package favorites

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("favorite not found")
)

// Manager handles favorites operations
type Manager struct {
	store *FavoritesStore
}

// NewManager creates a new Manager, loading existing favorites from disk
func NewManager() (*Manager, error) {
	store, err := Load()
	if err != nil {
		return nil, err
	}
	return &Manager{store: store}, nil
}

// Add adds a new favorite and returns its generated ID
func (m *Manager) Add(fav Favorite) (string, error) {
	// Generate unique ID using timestamp
	fav.ID = fmt.Sprintf("fav_%d", time.Now().UnixNano())
	fav.CreatedAt = time.Now().UTC()

	m.store.Favorites = append(m.store.Favorites, fav)

	err := Save(m.store)
	if err != nil {
		// Rollback on save failure
		m.store.Favorites = m.store.Favorites[:len(m.store.Favorites)-1]
		return "", err
	}

	return fav.ID, nil
}

// Get returns a favorite by ID
func (m *Manager) Get(id string) (*Favorite, error) {
	for i := range m.store.Favorites {
		if m.store.Favorites[i].ID == id {
			return &m.store.Favorites[i], nil
		}
	}
	return nil, ErrNotFound
}

// Remove deletes a favorite by ID
func (m *Manager) Remove(id string) error {
	idx := -1
	for i, fav := range m.store.Favorites {
		if fav.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrNotFound
	}

	// Remove by index
	removed := m.store.Favorites[idx]
	m.store.Favorites = append(m.store.Favorites[:idx], m.store.Favorites[idx+1:]...)

	err := Save(m.store)
	if err != nil {
		// Rollback on save failure
		m.store.Favorites = append(m.store.Favorites[:idx], append([]Favorite{removed}, m.store.Favorites[idx:]...)...)
		return err
	}

	return nil
}

// List returns all favorites
func (m *Manager) List() []Favorite {
	return m.store.Favorites
}
