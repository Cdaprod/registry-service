package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/Cdaprod/registry-service/internal/api"
	"github.com/Cdaprod/registry-service/internal/registry"
	"github.com/Cdaprod/registry-service/pkg/builtins"
	"github.com/Cdaprod/registry-service/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

//go:embed ../../web/build/*
var webUI embed.FS

func main() {
	// Initialize logger
	l := logger.NewLogger()
	defer l.Sync()

	// Create RegistryService
	rs := registry.NewRegistryService(l)

	// Add specific registries
	rs.AddRegistry("apis", registry.NewBaseRegistry(l))
	rs.AddRegistry("components", registry.NewBaseRegistry(l))

	// Initialize BuiltinLoader and load built-in plugins
	builtinLoader := builtins.NewBuiltinLoader(rs, "pkg/plugins/")
	if err := builtinLoader.LoadAll(); err != nil {
		l.Fatal("Error loading built-ins", zap.Error(err))
	}

	// Set up router
	r := mux.NewRouter()
	api.SetupRoutes(r, rs)

	// Serve static files for web UI
	webUIFS, err := fs.Sub(webUI, "web/build")
	if err != nil {
		l.Fatal("Failed to create sub filesystem", zap.Error(err))
	}

	r.PathPrefix("/").Handler(http.FileServer(http.FS(webUIFS)))

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
