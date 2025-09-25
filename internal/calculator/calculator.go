package calculator

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"calculator/internal/algorithm"
	"calculator/internal/models"
)

// Service handles pack calculation business logic
type Service struct {
	packSizes []int
	mu        sync.RWMutex
	cache     *Cache
	algo      *algorithm.Calculator
}

// NewService creates a new pack calculator service
func NewService(defaultPackSizes []int) *Service {
	service := &Service{
		cache: NewCache(),
		algo:  algorithm.NewCalculator(),
	}
	service.SetPackSizes(defaultPackSizes)
	return service
}

// SetPackSizes updates the available pack sizes
func (s *Service) SetPackSizes(sizes []int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Sort sizes in descending order for optimization
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	s.packSizes = sizes

	// Clear cache when pack sizes change
	s.cache.Clear()

	log.Printf("Pack sizes updated: %v", sizes)
}

// GetPackSizes returns the current pack sizes
func (s *Service) GetPackSizes() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sizes := make([]int, len(s.packSizes))
	copy(sizes, s.packSizes)
	return sizes
}

// CalculatePacks determines the optimal pack configuration for an order
func (s *Service) CalculatePacks(orderSize int) *models.CalculationResponse {
	s.mu.RLock()
	packSizes := make([]int, len(s.packSizes))
	copy(packSizes, s.packSizes)
	s.mu.RUnlock()

	// Check cache first
	if cached, ok := s.cache.Get(orderSize, packSizes); ok {
		return cached
	}

	// Calculate using algorithm
	result := s.algo.CalculateOptimalPacks(orderSize, packSizes)

	// Cache the result
	s.cache.Set(orderSize, packSizes, result)

	return result
}

// ValidatePackSizes checks if pack sizes are valid
func (s *Service) ValidatePackSizes(sizes []int) error {
	for _, size := range sizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be positive, got: %d", size)
		}
	}
	return nil
}
