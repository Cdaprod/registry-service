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

// Store adds or updates an Item in the storage
func (ms *MemoryStorage) Store(item *registry.Item) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if existing, ok := ms.items[item.ID]; ok {
		existing.Update(item.Name, item.Version, item.Metadata)
	} else {
		ms.items[item.ID] = item
	}

	return nil
}

// Get retrieves an Item from the storage
func (ms *MemoryStorage) Get(id string) (*registry.Item, error) {
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

// List returns all non-deleted Items in the storage
func (ms *MemoryStorage) List() []*registry.Item {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []*registry.Item
	for _, item := range ms.items {
		if !item.IsDeleted() {
			result = append(result, item)
		}
	}

	return result
}

// Delete soft-deletes an Item in the storage
func (ms *MemoryStorage) Delete(id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	item, ok := ms.items[id]
	if !ok {
		return errors.New("item not found")
	}

	item.SoftDelete()
	return nil
}

// Restore removes the soft-delete mark from an Item
func (ms *MemoryStorage) Restore(id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	item, ok := ms.items[id]
	if !ok {
		return errors.New("item not found")
	}

	item.Restore()
	return nil
}
