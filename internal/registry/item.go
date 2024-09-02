package registry

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Item represents an item in the registry with metadata and timestamps
type Item struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	Version   int64                  `json:"version"`
	deleted   bool                   // field to track if the item is deleted
	mu        sync.RWMutex           // mutex for thread-safe operations
}

// GetID returns the ID of the item
func (i *Item) GetID() string {
    return i.ID
}

// GetType returns the type of the item
func (i *Item) GetType() string {
    return i.Type
}

// Update updates the item's name, version, and metadata
func (i *Item) Update(name, itemType string, metadata map[string]interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Name = name
	i.Type = itemType
	i.Metadata = metadata
	i.Version++
	i.UpdatedAt = time.Now()
}

// IsDeleted checks if the item is marked as deleted
func (i *Item) IsDeleted() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.deleted
}

// SoftDelete marks the item as deleted
func (i *Item) SoftDelete() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.deleted = true
	i.UpdatedAt = time.Now()
}

// Restore removes the deleted mark from the item
func (i *Item) Restore() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.deleted = false
	i.UpdatedAt = time.Now()
}

// ItemStore represents an in-memory store for items
type ItemStore struct {
	mu    sync.RWMutex
	items map[string]*Item
}

// NewItemStore creates a new in-memory store for items
func NewItemStore() *ItemStore {
	return &ItemStore{
		items: make(map[string]*Item),
	}
}

// UpsertItem inserts or updates an item in the store
func (s *ItemStore) UpsertItem(item *Item) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.New().String()
		item.CreatedAt = time.Now()
		item.Version = 1
	} else {
		existingItem, exists := s.items[item.ID]
		if exists {
			if item.Version <= existingItem.Version {
				return existingItem, nil // Return existing item if version is not newer
			}
			item.CreatedAt = existingItem.CreatedAt
		} else {
			item.CreatedAt = time.Now()
		}
	}

	item.UpdatedAt = time.Now()
	s.items[item.ID] = item

	return item, nil
}

// GetItem retrieves an item from the store by ID
func (s *ItemStore) GetItem(id string) (*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, fmt.Errorf("item not found: %s", id)
	}

	if item.IsDeleted() {
		return nil, fmt.Errorf("item is deleted")
	}

	return item, nil
}

// DeleteItem removes an item from the store
func (s *ItemStore) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item, exists := s.items[id]; exists {
		item.SoftDelete()
		return nil
	}
	return fmt.Errorf("item not found: %s", id)
}

// RestoreItem restores a soft-deleted item in the store
func (s *ItemStore) RestoreItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item, exists := s.items[id]; exists {
		item.Restore()
		return nil
	}
	return fmt.Errorf("item not found: %s", id)
}

// ListItems returns all non-deleted items in the store
func (s *ItemStore) ListItems() []*Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*Item, 0, len(s.items))
	for _, item := range s.items {
		if !item.IsDeleted() {
			items = append(items, item)
		}
	}

	return items
}

// MarshalJSON implements custom JSON marshaling for Item
func (i *Item) MarshalJSON() ([]byte, error) {
	type Alias Item
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(i),
		CreatedAt: i.CreatedAt.Format(time.RFC3339),
		UpdatedAt: i.UpdatedAt.Format(time.RFC3339),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Item
func (i *Item) UnmarshalJSON(data []byte) error {
	type Alias Item
	aux := &struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	i.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	i.UpdatedAt, err = time.Parse(time.RFC3339, aux.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
