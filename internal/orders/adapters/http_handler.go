package adapters

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mmorall/booksonline/internal/orders"

	"github.com/google/uuid"
)

type HTTPHandler struct {
	svc       orders.Service
	adminUser string
	adminPass string
}

func NewHTTPHandler(svc orders.Service, adminUser, adminPass string) *HTTPHandler {
	return &HTTPHandler{
		svc:       svc,
		adminUser: adminUser,
		adminPass: adminPass,
	}
}

func (h *HTTPHandler) adminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != h.adminUser || pass != h.adminPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", h.handleCreateOrder)
	mux.HandleFunc("GET /orders", h.adminAuthMiddleware(h.handleListOrders))
	mux.HandleFunc("GET /orders/{id}", h.adminAuthMiddleware(h.handleGetOrder))
}

type CreateOrderRequest struct {
	CustomerEmail string         `json:"customer_email"`
	Items         map[string]int `json:"items"`
}

func (h *HTTPHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.CustomerEmail == "" {
		http.Error(w, "customer_email is required", http.StatusBadRequest)
		return
	}

	domainItems := make(map[uuid.UUID]int)
	for k, v := range req.Items {
		prodID, err := uuid.Parse(k)
		if err != nil {
			http.Error(w, "invalid product ID format: "+k, http.StatusBadRequest)
			return
		}
		domainItems[prodID] = v
	}

	order, err := h.svc.PlaceOrder(r.Context(), req.CustomerEmail, domainItems)
	if err != nil {
		if errors.Is(err, orders.ErrInvalidOrder) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *HTTPHandler) handleListOrders(w http.ResponseWriter, r *http.Request) {
	orderList, err := h.svc.ListOrders(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orderList); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *HTTPHandler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	order, err := h.svc.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
