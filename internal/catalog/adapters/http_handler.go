package adapters

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mmorall/booksonline/internal/catalog"

	"github.com/google/uuid"
)

type HTTPHandler struct {
	svc catalog.Service
}

func NewHTTPHandler(svc catalog.Service) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", h.handleListProducts)
	mux.HandleFunc("GET /products/{id}", h.handleGetProduct)
}

func (h *HTTPHandler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts(r.Context())
	if err != nil {
		// TODO: log the raw error and return a generic message
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *HTTPHandler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid product ID format", http.StatusBadRequest)
		return
	}

	product, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, catalog.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
