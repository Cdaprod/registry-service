package storage

import (
    "github.com/Cdaprod/registry-service/internal/registry"
)

// MemoryStorageAdapter adapts MemoryStorage to the Registry interface.
type MemoryStorageAdapter struct {
    storage *MemoryStorage
}

// NewMemoryStorageAdapter creates a new adapter for MemoryStorage.
func NewMemoryStorageAdapter(storage *MemoryStorage) *MemoryStorageAdapter {
    return &MemoryStorageAdapter{storage: storage}
}

// Register implements the Register method of the Registry interface.
func (a *MemoryStorageAdapter) Register(item registry.Registerable) error {
    // Conversion from Registerable to *registry.Item if necessary
    itemObj, ok := item.(*registry.Item)
    if !ok {
        return errors.New("invalid item type")
    }
    // Call CreateItem and handle its return values properly
    _, err := a.storage.CreateItem(itemObj)
    return err // Only return error as expected by the Registry interface
}

func (a *MemoryStorageAdapter) Get(id string) (registry.Registerable, bool) {
    item, err := a.storage.GetItem(id)
    if err != nil {
        return nil, false
    }
    return item, true
}

func (a *MemoryStorageAdapter) Unregister(id string) error {
    return a.storage.DeleteItem(id)
}

func (a *MemoryStorageAdapter) List() []registry.Registerable {
    items, _ := a.storage.ListItems()
    registerables := make([]registry.Registerable, len(items))
    for i, item := range items {
        registerables[i] = item
    }
    return registerables
}

func (a *MemoryStorageAdapter) ListByType(itemType string) []registry.Registerable {
    var result []registry.Registerable
    for _, item := range a.List() {
        if item.GetType() == itemType {
            result = append(result, item)
        }
    }
    return result
}