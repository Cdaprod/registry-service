package main

import (
    "log"
    "net/http"
    "os"

    "github.com/Cdaprod/registry-service/internal/api"
    "github.com/Cdaprod/registry-service/internal/storage"
    "github.com/Cdaprod/registry-service/pkg/builtins"
    "github.com/Cdaprod/registry-service/pkg/logger"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    l, err := logger.NewLogger()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer l.Sync()

    // Create MemoryStorage instance and its adapter
    memoryStorage := storage.NewMemoryStorage()
    
    // Create handler using MemoryStorage directly
    handler := api.NewHandler(memoryStorage, logger)

    // If needed for other components that expect the Registry interface:
    adapter := storage.NewMemoryStorageAdapter(memoryStorage)

    // Initialize BuiltinLoader and load built-in plugins
    builtinLoader := builtins.NewBuiltinLoader(adapter, "pkg/plugins/")
    if err := builtinLoader.LoadAll(); err != nil {
        l.Fatal("Error loading built-ins", zap.Error(err))
    }

    // Set up router using mux
    r := mux.NewRouter()
    api.SetupRoutes(r, adapter, l)

    // Serve static files from the web/build directory
    staticDir := "./web/build"
    fs := http.FileServer(http.Dir(staticDir))
    r.PathPrefix("/").Handler(fs)

    // Set up CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

    // Determine port
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Start server
    l.Info("Starting server", zap.String("port", port))
    log.Fatal(http.ListenAndServe(":"+port, c.Handler(r)))
}