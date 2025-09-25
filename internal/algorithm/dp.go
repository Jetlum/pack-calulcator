package algorithm

import (
	"calculator/internal/models"
)

// Calculator handles the dynamic programming logic for pack optimization
type Calculator struct{}

// NewCalculator creates a new algorithm calculator
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateOptimalPacks determines the optimal pack configuration for an order
// Rules:
// 1. Only whole packs can be sent
// 2. Send out the least amount of items to fulfill the order
// 3. Send out as few packs as possible
func (c *Calculator) CalculateOptimalPacks(orderSize int, packSizes []int) *models.CalculationResponse {
	// Edge case: if order is 0 or negative
	if orderSize <= 0 {
		return &models.CalculationResponse{
			Packs:      make(map[int]int),
			TotalItems: 0,
			TotalPacks: 0,
		}
	}

	// Dynamic programming approach to find optimal solution
	maxSize := orderSize + packSizes[0] // Add largest pack size as buffer
	dp := make([]models.DPEntry, maxSize+1)

	// Initialize dp array
	for i := range dp {
		dp[i] = models.DPEntry{
			TotalItems: -1, // -1 means impossible
			Packs:      make(map[int]int),
		}
	}

	// Base case: 0 items needs 0 packs
	dp[0] = models.DPEntry{
		TotalItems: 0,
		TotalPacks: 0,
		Packs:      make(map[int]int),
	}

	// Fill dp table
	for i := 1; i <= maxSize; i++ {
		for _, packSize := range packSizes {
			if packSize <= i && dp[i-packSize].TotalItems != -1 {
				newTotalItems := dp[i-packSize].TotalItems + packSize
				newTotalPacks := dp[i-packSize].TotalPacks + 1

				// Check if this is a better solution
				if dp[i].TotalItems == -1 ||
					(i >= orderSize && newTotalItems < dp[i].TotalItems) ||
					(i >= orderSize && newTotalItems == dp[i].TotalItems && newTotalPacks < dp[i].TotalPacks) {
					// Copy the previous solution and add current pack
					newPacks := make(map[int]int)
					for k, v := range dp[i-packSize].Packs {
						newPacks[k] = v
					}
					newPacks[packSize]++

					dp[i] = models.DPEntry{
						TotalItems: newTotalItems,
						TotalPacks: newTotalPacks,
						Packs:      newPacks,
					}
				}
			}
		}
	}

	// Find the optimal solution (least items >= orderSize)
	bestSolution := models.DPEntry{TotalItems: -1}
	for i := orderSize; i <= maxSize; i++ {
		if dp[i].TotalItems != -1 {
			if bestSolution.TotalItems == -1 ||
				dp[i].TotalItems < bestSolution.TotalItems ||
				(dp[i].TotalItems == bestSolution.TotalItems && dp[i].TotalPacks < bestSolution.TotalPacks) {
				bestSolution = dp[i]
			}
		}
	}

	return &models.CalculationResponse{
		Packs:      bestSolution.Packs,
		TotalItems: bestSolution.TotalItems,
		TotalPacks: bestSolution.TotalPacks,
	}
}
