package registry

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	Version   int64                  `json:"version"`
}

type ItemStore struct {
	mu    sync.RWMutex
	items map[string]*Item
}

func NewItemStore() *ItemStore {
	return &ItemStore{
		items: make(map[string]*Item),
	}
}

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

func (s *ItemStore) GetItem(id string) (*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, fmt.Errorf("item not found: %s", id)
	}

	return item, nil
}

func (s *ItemStore) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, id)
	return nil
}

func (s *ItemStore) ListItems() []*Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
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
