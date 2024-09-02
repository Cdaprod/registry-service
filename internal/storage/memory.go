package storage

import (
	"errors"
	"sync"

	"github.com/Cdaprod/registry-service/internal/registry"
)

// MemoryStorage implements in-memory storage for Items
type MemoryStorage struct {
	items map[string]*registry.Item
	mu    sync.RWMutex
}

// NewMemoryStorage creates a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[string]*registry.Item),
	}
}

// List returns a slice of non-deleted Items with pagination support
func (ms *MemoryStorage) List(limit, offset int) []*registry.Item {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []*registry.Item
	for _, item := range ms.items {
		if !item.IsDeleted() {
			result = append(result, item)
		}
	}

	// Apply pagination
	if offset > len(result) {
		return []*registry.Item{}
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end]
}

// ListItems returns all non-deleted Items in the storage without pagination
func (ms *MemoryStorage) ListItems() ([]*registry.Item, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []*registry.Item
	for _, item := range ms.items {
		if !item.IsDeleted() {
			result = append(result, item)
		}
	}

	return result, nil
}

// CreateItem adds an Item to the storage
func (ms *MemoryStorage) CreateItem(item *registry.Item) (*registry.Item, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.items[item.ID]; exists {
		return nil, errors.New("item already exists")
	}

	item.Version = 1
	ms.items[item.ID] = item
	return item, nil
}

<<<<<<< HEAD
// GetAsRegisterable retrieves an item from the storage and returns it as a Registerable interface
func (ms *MemoryStorage) GetAsRegisterable(id string) (registry.Registerable, bool) {
    ms.mu.RLock()
    defer ms.mu.RUnlock()

    item, exists := ms.items[id]
    if !exists || item.IsDeleted() {
        return nil, false
    }
    return item, true
}

// Get retrieves an Item from the storage
func (ms *MemoryStorage) Get(id string) (*registry.Item, error) {
=======
// GetItem retrieves an Item from the storage
func (ms *MemoryStorage) GetItem(id string) (*registry.Item, error) {
>>>>>>> a8896e1 (commit updated go.mod)
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	item, ok := ms.items[id]
	if !ok {
		return nil, errors.New("item not found")
	}

	if item.IsDeleted() {
		return nil, errors.New("item is deleted")
	}

	return item, nil
}

// UpdateItem updates an existing Item in the storage
func (ms *MemoryStorage) UpdateItem(item *registry.Item) (*registry.Item, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	existing, ok := ms.items[item.ID]
	if !ok {
		return nil, errors.New("item not found")
	}

	// Update fields
	existing.Name = item.Name
	existing.Metadata = item.Metadata
	existing.Version++

	return existing, nil
}

// DeleteItem soft-deletes an Item in the storage
func (ms *MemoryStorage) DeleteItem(id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	item, ok := ms.items[id]
	if !ok {
		return errors.New("item not found")
	}

	item.SoftDelete()
	return nil
}
