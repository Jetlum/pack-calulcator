package calculator

import (
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	service := NewService([]int{250, 500, 1000})

	if service == nil {
		t.Fatal("NewService returned nil")
	}

	expectedSizes := []int{1000, 500, 250} // Should be sorted descending
	sizes := service.GetPackSizes()

	if !reflect.DeepEqual(sizes, expectedSizes) {
		t.Errorf("Expected default sizes %v, got %v", expectedSizes, sizes)
	}
}

func TestSetAndGetPackSizes(t *testing.T) {
	service := NewService([]int{})
	newSizes := []int{100, 200, 300}

	service.SetPackSizes(newSizes)

	// Should be sorted in descending order
	expected := []int{300, 200, 100}
	result := service.GetPackSizes()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected sizes %v, got %v", expected, result)
	}
}

func TestCalculatePacksEdgeCase(t *testing.T) {
	// Critical edge case from requirements
	service := NewService([]int{23, 31, 53})

	result := service.CalculatePacks(500000)

	expectedPacks := map[int]int{23: 2, 31: 7, 53: 9429}
	expectedTotal := 23*2 + 31*7 + 53*9429 // 500000

	if !reflect.DeepEqual(result.Packs, expectedPacks) {
		t.Errorf("Edge case failed: got packs %v, want %v", result.Packs, expectedPacks)
	}

	if result.TotalItems != expectedTotal {
		t.Errorf("Edge case total mismatch: got %d, want %d", result.TotalItems, expectedTotal)
	}
}

func TestValidatePackSizes(t *testing.T) {
	service := NewService([]int{})

	// Valid sizes
	if err := service.ValidatePackSizes([]int{100, 200, 300}); err != nil {
		t.Errorf("Valid sizes should not return error: %v", err)
	}

	// Invalid sizes
	if err := service.ValidatePackSizes([]int{100, -50, 300}); err == nil {
		t.Error("Invalid sizes should return error")
	}
}
