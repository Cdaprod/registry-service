package api

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/Cdaprod/registry-service/internal/registry"
    "github.com/Cdaprod/registry-service/internal/storage"
    "github.com/gorilla/mux"
    "go.uber.org/zap"
)

type Handler struct {
    store  *storage.MemoryStorage
    logger *zap.Logger
}

func NewHandler(store *storage.MemoryStorage, logger *zap.Logger) *Handler {
    return &Handler{
        store:  store,
        logger: logger,
    }
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to the Registry Service!"))
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func (h *Handler) ServeDocs(w http.ResponseWriter, r *http.Request) {
    // Implement API documentation serving (e.g., Swagger UI)
    w.Write([]byte("API Documentation (To be implemented)"))
}


func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
    var item registry.Item
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        h.logger.Error("Failed to decode request body", zap.Error(err))
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    createdItem, err := h.store.CreateItem(&item)
    if err != nil {
        h.logger.Error("Failed to create item", zap.Error(err))
        http.Error(w, "Failed to create item", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdItem)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]

    item, err := h.store.GetItem(id)
    if err != nil {
        h.logger.Error("Failed to get item", zap.Error(err))
        http.Error(w, "Item not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(item)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]

    var item registry.Item
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        h.logger.Error("Failed to decode request body", zap.Error(err))
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    item.ID = id
    updatedItem, err := h.store.UpdateItem(&item)
    if err != nil {
        h.logger.Error("Failed to update item", zap.Error(err))
        http.Error(w, "Failed to update item", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedItem)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]

    if err := h.store.DeleteItem(id); err != nil {
        h.logger.Error("Failed to delete item", zap.Error(err))
        http.Error(w, "Failed to delete item", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    var items []registry.Registerable

    if limit > 0 || offset > 0 {
        // Use ListPaginated if limit or offset is specified
        items = h.store.ListPaginated(limit, offset)
    } else {
        // Use List if no pagination is specified
        items = h.store.List()
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

// Add a new method for error responses
func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
    h.respondWithJSON(w, code, map[string]string{"error": message})
}

// Add a new method for JSON responses
func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}