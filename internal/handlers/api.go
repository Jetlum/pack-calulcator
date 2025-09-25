package handlers

import (
	"encoding/json"
	"net/http"

	"calculator/internal/calculator"
	"calculator/internal/models"
)

// Handler manages HTTP endpoints
type Handler struct {
	calcService *calculator.Service
}

// NewHandler creates a new HTTP handler
func NewHandler(calcService *calculator.Service) *Handler {
	return &Handler{
		calcService: calcService,
	}
}

// RegisterRoutes sets up HTTP routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Apply middleware
	mux.Handle("/", LoggingMiddleware(CORSMiddleware(http.HandlerFunc(h.handleIndex))))
	mux.Handle("/api/calculate", LoggingMiddleware(CORSMiddleware(http.HandlerFunc(h.handleCalculate))))
	mux.Handle("/api/pack-sizes", LoggingMiddleware(CORSMiddleware(http.HandlerFunc(h.handlePackSizes))))
}

// handleCalculate handles pack calculation requests
func (h *Handler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := h.calcService.CalculatePacks(req.Items)
	h.sendJSON(w, result)
}

// handlePackSizes handles pack size operations (GET/POST)
func (h *Handler) handlePackSizes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGetPackSizes(w, r)
	case "POST":
		h.handleUpdatePackSizes(w, r)
	default:
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetPackSizes returns current pack sizes
func (h *Handler) handleGetPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes := h.calcService.GetPackSizes()
	response := models.APIResponse{
		Success:   true,
		PackSizes: sizes,
	}
	h.sendJSON(w, response)
}

// handleUpdatePackSizes handles pack size updates
func (h *Handler) handleUpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	var update models.PackSizeUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate pack sizes
	if err := h.calcService.ValidatePackSizes(update.PackSizes); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.calcService.SetPackSizes(update.PackSizes)

	response := models.APIResponse{
		Success:   true,
		PackSizes: update.PackSizes,
	}
	h.sendJSON(w, response)
}

// sendJSON sends a JSON response
func (h *Handler) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func (h *Handler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: false,
		Error:   message,
	})
}
