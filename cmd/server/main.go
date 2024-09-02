package main

import (
    "log"
    "os"
    "os/signal"
    "net/http"
    "context"
    "syscall"
    "time"

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

    // Create MemoryStorage instance
    memoryStorage := storage.NewMemoryStorage()

    // Initialize BuiltinLoader and load built-in plugins
    builtinLoader := builtins.NewBuiltinLoader(memoryStorage, "pkg/plugins/")
    if err := builtinLoader.LoadAll(); err != nil {
        l.Fatal("Error loading built-ins", zap.Error(err))
    }

    // Set up router using mux
    r := mux.NewRouter()
    api.SetupRoutes(r, memoryStorage, l)

    // Serve static files from the web/build directory
    fs := http.FileServer(http.Dir("./web/build"))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    // Serve index.html for any other routes
    r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/build/index.html")
    })

    // Determine port and bind address
    port := os.Getenv("PORT")
    if port == "" {
        port = "7777"
    }
    bindAddr := "0.0.0.0:" + port

    // Set up CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},  // Allow all origins for staged environment
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

    // Wrap our router with the CORS handler
    handler := c.Handler(r)

    // Start server
    l.Info("Starting server", 
        zap.String("port", port),
        zap.String("bind_address", bindAddr))
    
    server := &http.Server{
        Addr:    bindAddr,
        Handler: handler,
    }

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            l.Fatal("Failed to start server", zap.Error(err))
        }
    }()

    l.Info("Server is ready to handle requests")

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    l.Info("Server is shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        l.Fatal("Server forced to shutdown", zap.Error(err))
    }

    l.Info("Server has shut down gracefully")
}