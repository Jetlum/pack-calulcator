package models

// CalculationRequest represents the request to calculate packs for an order
type CalculationRequest struct {
	Items int `json:"items"`
}

// CalculationResponse represents the response with pack quantities
type CalculationResponse struct {
	Packs      map[int]int `json:"packs"`
	TotalItems int         `json:"totalItems"`
	TotalPacks int         `json:"totalPacks"`
}

// PackSizeUpdate represents a request to update pack sizes
type PackSizeUpdate struct {
	PackSizes []int `json:"packSizes"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	PackSizes []int       `json:"packSizes,omitempty"`
}

// DPEntry represents an entry in the dynamic programming table
type DPEntry struct {
	TotalItems int
	TotalPacks int
	Packs      map[int]int
}
