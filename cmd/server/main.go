package main

import (
    "context"
    "log"
    "mime"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "strings"
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

// setCorrectMIMEType sets the MIME type for static files based on their extension.
func setCorrectMIMEType(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ext := filepath.Ext(r.URL.Path)
        if mimeType := mime.TypeByExtension(ext); mimeType != "" {
            w.Header().Set("Content-Type", mimeType)
        } else {
            w.Header().Set("Content-Type", "application/octet-stream")
        }
        w.Header().Set("X-Content-Type-Options", "nosniff") // Prevent MIME type sniffing
        next.ServeHTTP(w, r)
    })
}

// initializeServer sets up and starts the HTTP server with all configurations.
func initializeServer(router http.Handler, bindAddr string, l *zap.Logger) *http.Server {
    l.Info("Starting server", zap.String("bind_address", bindAddr))

    server := &http.Server{
        Addr:    bindAddr,
        Handler: router,
    }

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            l.Fatal("Failed to start server", zap.Error(err))
        }
    }()

    l.Info("Server is ready to handle requests")
    return server
}

// handleGracefulShutdown gracefully shuts down the server on receiving a termination signal.
func handleGracefulShutdown(server *http.Server, l *zap.Logger) {
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

// main is the entry point for the application.
func main() {
    // Initialize logger
    l, err := logger.NewLogger()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer l.Sync()

    // Initialize in-memory storage
    memoryStorage := storage.NewMemoryStorage()

    // Load built-in plugins
    builtinLoader := builtins.NewBuiltinLoader(memoryStorage, "pkg/plugins/")
    if err := builtinLoader.LoadAll(); err != nil {
        l.Fatal("Error loading built-ins", zap.Error(err))
    }

    // Set up router using mux
    r := mux.NewRouter()
    api.SetupRoutes(r, memoryStorage, l)

    // Serve static files from the web/build directory with correct MIME types
    fs := http.FileServer(http.Dir("./web/build"))
    r.PathPrefix("/static/").Handler(setCorrectMIMEType(http.StripPrefix("/static/", fs)))

    // Serve index.html for any non-static file requests (fallback for React Router)
    r.PathPrefix("/").Handler(setCorrectMIMEType(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./web/build/index.html")
    })))

    // Determine port and bind address
    port := os.Getenv("PORT")
    if port == "" {
        port = "7777"
    }
    bindAddr := "0.0.0.0:" + port

    // Set up CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},  // Allow all origins for staged environment
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

    // Wrap router with CORS handler
    handler := c.Handler(r)

    // Start the HTTP server
    server := initializeServer(handler, bindAddr, l)

    // Handle graceful shutdown
    handleGracefulShutdown(server, l)
}