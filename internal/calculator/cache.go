package calculator

import (
	"fmt"
	"sync"

	"calculator/internal/models"
)

// Cache provides thread-safe caching for calculation results
type Cache struct {
	mu    sync.RWMutex
	cache map[string]*models.CalculationResponse
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]*models.CalculationResponse),
	}
}

// Get retrieves a cached result
func (c *Cache) Get(orderSize int, packSizes []int) (*models.CalculationResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.generateKey(orderSize, packSizes)
	result, exists := c.cache[key]
	return result, exists
}

// Set stores a result in the cache
func (c *Cache) Set(orderSize int, packSizes []int, result *models.CalculationResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(orderSize, packSizes)
	c.cache[key] = result
}

// Clear removes all cached entries
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*models.CalculationResponse)
}

// generateKey creates a cache key from order size and pack sizes
func (c *Cache) generateKey(orderSize int, packSizes []int) string {
	return fmt.Sprintf("%d_%v", orderSize, packSizes)
}
