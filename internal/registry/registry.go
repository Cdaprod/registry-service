package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Registerable is an interface for any item that can be registered
type Registerable interface {
	GetID() string
	GetType() string
}

// Registry interface defines methods for managing registerable items
type Registry interface {
	Register(item Registerable) error
	Get(id string) (Registerable, bool)
	Unregister(id string) error
	List() []Registerable
	ListByType(itemType string) []Registerable
}

// CentralRegistry provides a thread-safe implementation of the Registry interface
type CentralRegistry struct {
	mu    sync.RWMutex
	items map[string]Registerable
}

func NewCentralRegistry() *CentralRegistry {
	return &CentralRegistry{
		items: make(map[string]Registerable),
	}
}

func (r *CentralRegistry) Register(item Registerable) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[item.GetID()]; exists {
		return fmt.Errorf("item already registered: %s", item.GetID())
	}
	r.items[item.GetID()] = item
	return nil
}

func (r *CentralRegistry) Get(id string) (Registerable, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.items[id]
	return item, exists
}

func (r *CentralRegistry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[id]; !exists {
		return fmt.Errorf("item not found: %s", id)
	}
	delete(r.items, id)
	return nil
}

func (r *CentralRegistry) List() []Registerable {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var items []Registerable
	for _, item := range r.items {
		items = append(items, item)
	}
	return items
}

func (r *CentralRegistry) ListByType(itemType string) []Registerable {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var items []Registerable
	for _, item := range r.items {
		if item.GetType() == itemType {
			items = append(items, item)
		}
	}
	return items
}

// RegistryServer provides HTTP handlers for interacting with the registry
type RegistryServer struct {
	registry Registry
}

func NewRegistryServer(registry Registry) *RegistryServer {
	return &RegistryServer{registry: registry}
}

func (s *RegistryServer) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Implementation for handling registration via HTTP
}

func (s *RegistryServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Implementation for handling get requests via HTTP
}

func (s *RegistryServer) HandleList(w http.ResponseWriter, r *http.Request) {
	items := s.registry.List()
	json.NewEncoder(w).Encode(items)
}

// SetupRoutes configures the HTTP routes for the registry server
func (s *RegistryServer) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", s.HandleRegister)
	mux.HandleFunc("/get", s.HandleGet)
	mux.HandleFunc("/list", s.HandleList)
	return mux
}

// // Example usage
// func main() {
// 	registry := NewCentralRegistry()
// 	server := NewRegistryServer(registry)

//     // Initialize BuiltinLoader and load built-in plugins
//     builtinLoader := builtins.NewBuiltinLoader(reg, "pkg/plugins/")
//     if err := builtinLoader.LoadAll(); err != nil {
//         fmt.Printf("Error loading built-ins: %v\n", err)
//     }

// 	// Start the HTTP server
// 	http.ListenAndServe(":8080", server.SetupRoutes())
// }
