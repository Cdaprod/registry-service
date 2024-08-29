package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"github.com/Cdaprod/registry-service/internal/storage"
)

func SetupRoutes(r *mux.Router, store *storage.MemoryStore, logger *zap.Logger) {
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