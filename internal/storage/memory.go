package storage

import (
	"errors"
	"sync"

	"github.com/Cdaprod/registry-service/internal/registry"
)

var _ registry.Registry = (*MemoryStorage)(nil)

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

// Register adds or updates an Item in the storage
// Register adds or updates an Item in the storage
func (ms *MemoryStorage) Register(item registry.Registerable) error {
    ms.mu.Lock()
    defer ms.mu.Unlock()

    itemObj, ok := item.(*registry.Item)
    if !ok {
        return errors.New("invalid item type")
    }

    if itemObj.RegistryName == "" {
        return errors.New("registry name must be set")
    }

    if existing, exists := ms.items[itemObj.ID]; exists {
        existing.Name = itemObj.Name
        existing.Metadata = itemObj.Metadata
        existing.Version++
    } else {
        itemObj.Version = 1
        ms.items[itemObj.ID] = itemObj
    }

    return nil
}

// Get retrieves an item from the storage
func (ms *MemoryStorage) Get(id string) (registry.Registerable, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	item, exists := ms.items[id]
	if !exists || item.IsDeleted() {
		return nil, false
	}
	return item, true
}

// Unregister soft-deletes an Item in the storage
func (ms *MemoryStorage) Unregister(id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	item, ok := ms.items[id]
	if !ok {
		return errors.New("item not found")
	}

	item.SoftDelete()
	return nil
}

// List returns all non-deleted Items in the storage
func (ms *MemoryStorage) List() []registry.Registerable {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []registry.Registerable
	for _, item := range ms.items {
		if !item.IsDeleted() {
			result = append(result, item)
		}
	}

	return result
}

// ListByType returns all non-deleted Items of a specific type
func (ms *MemoryStorage) ListByType(itemType string) []registry.Registerable {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []registry.Registerable
	for _, item := range ms.items {
		if !item.IsDeleted() && item.GetType() == itemType {
			result = append(result, item)
		}
	}

	return result
}

// ListByRegistryName returns all non-deleted Items of a specific registry name
func (ms *MemoryStorage) ListByRegistryName(registryName string) []registry.Registerable {
    ms.mu.RLock()
    defer ms.mu.RUnlock()

    var result []registry.Registerable
    for _, item := range ms.items {
        if !item.IsDeleted() && item.RegistryName == registryName {
            result = append(result, item)
        }
    }

    return result
}

// ListPaginated returns a slice of non-deleted Items with pagination support
func (ms *MemoryStorage) ListPaginated(limit, offset int) []registry.Registerable {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []registry.Registerable
	for _, item := range ms.items {
		if !item.IsDeleted() {
			result = append(result, item)
		}
	}

	// Apply pagination
	if offset > len(result) {
		return []registry.Registerable{}
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end]
}

// Additional methods for compatibility with existing code

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
	return item, ms.Register(item)
}

// GetItem retrieves an Item from the storage
func (ms *MemoryStorage) GetItem(id string) (*registry.Item, error) {
	item, ok := ms.Get(id)
	if !ok {
		return nil, errors.New("item not found")
	}
	return item.(*registry.Item), nil
}

// UpdateItem updates an existing Item in the storage
func (ms *MemoryStorage) UpdateItem(item *registry.Item) (*registry.Item, error) {
	err := ms.Register(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteItem soft-deletes an Item in the storage
func (ms *MemoryStorage) DeleteItem(id string) error {
	return ms.Unregister(id)
}