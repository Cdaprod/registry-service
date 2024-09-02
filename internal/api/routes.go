package api

import (
    "net/http"
    "encoding/json"

    "github.com/Cdaprod/registry-service/internal/storage"
    "github.com/gorilla/mux"
    "go.uber.org/zap"
)

func SetupRoutes(r *mux.Router, store *storage.MemoryStorage, logger *zap.Logger) {
    handler := NewHandler(store, logger)

    // API versioning
    v1 := r.PathPrefix("/api/v1").Subrouter()

    // Items endpoints
    v1.HandleFunc("/items", handler.CreateItem).Methods("POST")
    v1.HandleFunc("/items", handler.ListItems).Methods("GET")
    v1.HandleFunc("/items/{id}", handler.GetItem).Methods("GET")
    v1.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
    v1.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE")

    // New routes for RegistryDashboard
    v1.HandleFunc("/registries", handler.ListRegistries).Methods("GET")
    v1.HandleFunc("/registry/{name}/list", handler.ListRegistryItems).Methods("GET")

    // Health check endpoint
    r.HandleFunc("/health", handler.HealthCheck).Methods("GET")

    // Documentation endpoint (consider implementing Swagger/OpenAPI)
    r.HandleFunc("/docs", handler.ServeDocs).Methods("GET")

    // Root handler
    r.HandleFunc("/", handler.HomeHandler).Methods("GET")

    // Middleware for logging, CORS, etc.
    r.Use(loggingMiddleware(logger))
    r.Use(corsMiddleware)

    // Serve static files from the web/build directory
    fs := http.FileServer(http.Dir("./web/build"))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    // Serve index.html for any other routes
    r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/build/index.html")
    })
}

func (h *Handler) ListRegistries(w http.ResponseWriter, r *http.Request) {
    // For now, let's return a mock list of registries
    registries := []string{"Registry1", "Registry2", "Registry3"}
    json.NewEncoder(w).Encode(registries)
}

func (h *Handler) ListRegistryItems(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    registryName := vars["name"]

    // Use the new method to list items by registry name
    items := h.store.ListByRegistryName(registryName)

    json.NewEncoder(w).Encode(items)
}

func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Info("Received request", 
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("remote_addr", r.RemoteAddr))
            next.ServeHTTP(w, r)
        })
    }
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins for now
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}