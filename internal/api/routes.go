package api

import (
	"net/http"

	"github.com/Cdaprod/registry-service/internal/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func SetupRoutes(r *mux.Router, store *storage.MemoryStorage, logger *zap.Logger) { // Updated type name
	handler := NewHandler(store, logger)

	r.HandleFunc("/api/items", handler.CreateItem).Methods("POST")
	r.HandleFunc("/api/items", handler.ListItems).Methods("GET")
	r.HandleFunc("/api/items/{id}", handler.GetItem).Methods("GET")
	r.HandleFunc("/api/items/{id}", handler.UpdateItem).Methods("PUT")
	r.HandleFunc("/api/items/{id}", handler.DeleteItem).Methods("DELETE")

	// Add a simple health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}
