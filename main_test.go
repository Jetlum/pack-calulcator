package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "testing"
)

func TestNewPackCalculator(t *testing.T) {
    pc := NewPackCalculator()
    
    if pc == nil {
        t.Fatal("NewPackCalculator returned nil")
    }
    
    expectedSizes := []int{5000, 2000, 1000, 500, 250}
    sizes := pc.GetPackSizes()
    
    if !reflect.DeepEqual(sizes, expectedSizes) {
        t.Errorf("Expected default sizes %v, got %v", expectedSizes, sizes)
    }
}

func TestSetAndGetPackSizes(t *testing.T) {
    pc := NewPackCalculator()
    newSizes := []int{100, 200, 300}
    
    pc.SetPackSizes(newSizes)
    
    // Should be sorted in descending order
    expected := []int{300, 200, 100}
    result := pc.GetPackSizes()
    
    if !reflect.DeepEqual(result, expected) {
        t.Errorf("Expected sizes %v, got %v", expected, result)
    }
}

func TestCalculatePacksBasicCases(t *testing.T) {
    tests := []struct {
        name       string
        orderSize  int
        packSizes  []int
        wantPacks  map[int]int
        wantTotal  int
    }{
        {
            name:      "Single pack exact match",
            orderSize: 250,
            packSizes: []int{250, 500, 1000},
            wantPacks: map[int]int{250: 1},
            wantTotal: 250,
        },
        {
            name:      "Need next larger pack",
            orderSize: 251,
            packSizes: []int{250, 500, 1000},
            wantPacks: map[int]int{500: 1},
            wantTotal: 500,
        },
        {
            name:      "Combination of packs",
            orderSize: 501,
            packSizes: []int{250, 500, 1000},
            wantPacks: map[int]int{250: 1, 500: 1},
            wantTotal: 750,
        },
        {
            name:      "Large order",
            orderSize: 12001,
            packSizes: []int{250, 500, 1000, 2000, 5000},
            wantPacks: map[int]int{250: 1, 2000: 1, 5000: 2},
            wantTotal: 12250,
        },
        {
            name:      "Zero order",
            orderSize: 0,
            packSizes: []int{250, 500},
            wantPacks: map[int]int{},
            wantTotal: 0,
        },
        {
            name:      "Negative order",
            orderSize: -10,
            packSizes: []int{250, 500},
            wantPacks: map[int]int{},
            wantTotal: 0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pc := NewPackCalculator()
            pc.SetPackSizes(tt.packSizes)
            
            result := pc.CalculatePacks(tt.orderSize)
            
            if !reflect.DeepEqual(result.Packs, tt.wantPacks) {
                t.Errorf("Packs mismatch: got %v, want %v", result.Packs, tt.wantPacks)
            }
            
            if result.TotalItems != tt.wantTotal {
                t.Errorf("Total items mismatch: got %d, want %d", result.TotalItems, tt.wantTotal)
            }
        })
    }
}

func TestCalculatePacksEdgeCase(t *testing.T) {
    // Critical edge case from requirements
    pc := NewPackCalculator()
    pc.SetPackSizes([]int{23, 31, 53})
    
    result := pc.CalculatePacks(500000)
    
    expectedPacks := map[int]int{23: 2, 31: 7, 53: 9429}
    expectedTotal := 23*2 + 31*7 + 53*9429 // 500000
    
    if !reflect.DeepEqual(result.Packs, expectedPacks) {
        t.Errorf("Edge case failed: got packs %v, want %v", result.Packs, expectedPacks)
    }
    
    if result.TotalItems != expectedTotal {
        t.Errorf("Edge case total mismatch: got %d, want %d", result.TotalItems, expectedTotal)
    }
}

func TestCaching(t *testing.T) {
    pc := NewPackCalculator()
    
    // First call - should calculate
    result1 := pc.CalculatePacks(500)
    
    // Second call - should use cache
    result2 := pc.CalculatePacks(500)
    
    if !reflect.DeepEqual(result1, result2) {
        t.Error("Cached results don't match")
    }
    
    // Change pack sizes - should clear cache
    pc.SetPackSizes([]int{100, 200})
    result3 := pc.CalculatePacks(500)
    
    // Should have different result with different pack sizes
    if reflect.DeepEqual(result1.Packs, result3.Packs) {
        t.Error("Cache wasn't cleared when pack sizes changed")
    }
}

func TestHTTPHandlers(t *testing.T) {
    pc := NewPackCalculator()
    
    t.Run("Calculate endpoint", func(t *testing.T) {
        reqBody := bytes.NewBuffer([]byte(`{"items": 263}`))
        req := httptest.NewRequest("POST", "/api/calculate", reqBody)
        req.Header.Set("Content-Type", "application/json")
        
        w := httptest.NewRecorder()
        pc.handleCalculate(w, req)
        
        if w.Code != http.StatusOK {
            t.Errorf("Expected status 200, got %d", w.Code)
        }
        
        var response CalculationResponse
        if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }
        
        if response.TotalItems == 0 {
            t.Error("Expected non-zero total items")
        }
    })
    
    t.Run("Get pack sizes endpoint", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/pack-sizes", nil)
        w := httptest.NewRecorder()
        
        pc.handleGetPackSizes(w, req)
        
        if w.Code != http.StatusOK {
            t.Errorf("Expected status 200, got %d", w.Code)
        }
    })
    
    t.Run("Update pack sizes endpoint", func(t *testing.T) {
        reqBody := bytes.NewBuffer([]byte(`{"packSizes": [100, 200, 300]}`))
        req := httptest.NewRequest("POST", "/api/update-pack-sizes", reqBody)
        req.Header.Set("Content-Type", "application/json")
        
        w := httptest.NewRecorder()
        pc.handleUpdatePackSizes(w, req)
        
        if w.Code != http.StatusOK {
            t.Errorf("Expected status 200, got %d", w.Code)
        }
        
        // Verify sizes were updated
        sizes := pc.GetPackSizes()
        expected := []int{300, 200, 100}
        if !reflect.DeepEqual(sizes, expected) {
            t.Errorf("Sizes not updated correctly: got %v, want %v", sizes, expected)
        }
    })
}

func TestInvalidRequests(t *testing.T) {
    pc := NewPackCalculator()
    
    t.Run("Invalid JSON", func(t *testing.T) {
        reqBody := bytes.NewBuffer([]byte(`invalid json`))
        req := httptest.NewRequest("POST", "/api/calculate", reqBody)
        w := httptest.NewRecorder()
        
        pc.handleCalculate(w, req)
        
        if w.Code != http.StatusBadRequest {
            t.Errorf("Expected status 400, got %d", w.Code)
        }
    })
    
    t.Run("Invalid pack sizes", func(t *testing.T) {
        reqBody := bytes.NewBuffer([]byte(`{"packSizes": [100, -50, 300]}`))
        req := httptest.NewRequest("POST", "/api/update-pack-sizes", reqBody)
        w := httptest.NewRecorder()
        
        pc.handleUpdatePackSizes(w, req)
        
        if w.Code != http.StatusBadRequest {
            t.Errorf("Expected status 400, got %d", w.Code)
        }
    })
}