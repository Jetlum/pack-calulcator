package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sort"
    "strconv"
    "strings"
    "sync"
)

// PackSize represents a single pack size option
type PackSize struct {
    Size int `json:"size"`
}

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

// PackCalculator handles the pack calculation logic
type PackCalculator struct {
    packSizes []int
    mu        sync.RWMutex
    cache     map[string]*CalculationResponse
}

// NewPackCalculator creates a new pack calculator with default pack sizes
func NewPackCalculator() *PackCalculator {
    return &PackCalculator{
        packSizes: []int{250, 500, 1000, 2000, 5000},
        cache:     make(map[string]*CalculationResponse),
    }
}

// SetPackSizes updates the available pack sizes
func (pc *PackCalculator) SetPackSizes(sizes []int) {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    
    // Sort sizes in descending order for optimization
    sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
    pc.packSizes = sizes
    // Clear cache when pack sizes change
    pc.cache = make(map[string]*CalculationResponse)
    
    log.Printf("Pack sizes updated: %v", sizes)
}

// GetPackSizes returns the current pack sizes
func (pc *PackCalculator) GetPackSizes() []int {
    pc.mu.RLock()
    defer pc.mu.RUnlock()
    
    sizes := make([]int, len(pc.packSizes))
    copy(sizes, pc.packSizes)
    return sizes
}

// CalculatePacks determines the optimal pack configuration for an order
// Rules:
// 1. Only whole packs can be sent
// 2. Send out the least amount of items to fulfill the order
// 3. Send out as few packs as possible
func (pc *PackCalculator) CalculatePacks(orderSize int) *CalculationResponse {
    pc.mu.RLock()
    defer pc.mu.RUnlock()
    
    // Check cache first
    cacheKey := fmt.Sprintf("%d_%v", orderSize, pc.packSizes)
    if cached, ok := pc.cache[cacheKey]; ok {
        return cached
    }
    
    // Edge case: if order is 0 or negative
    if orderSize <= 0 {
        result := &CalculationResponse{
            Packs:      make(map[int]int),
            TotalItems: 0,
            TotalPacks: 0,
        }
        pc.cache[cacheKey] = result
        return result
    }
    
    // Dynamic programming approach to find optimal solution
    // dp[i] stores the optimal solution for i items
    maxSize := orderSize + pc.packSizes[0] // Add largest pack size as buffer
    dp := make([]dpEntry, maxSize+1)
    
    // Initialize dp array
    for i := range dp {
        dp[i] = dpEntry{
            totalItems: -1, // -1 means impossible
            packs:      make(map[int]int),
        }
    }
    
    // Base case: 0 items needs 0 packs
    dp[0] = dpEntry{
        totalItems: 0,
        totalPacks: 0,
        packs:      make(map[int]int),
    }
    
    // Fill dp table
    for i := 1; i <= maxSize; i++ {
        for _, packSize := range pc.packSizes {
            if packSize <= i && dp[i-packSize].totalItems != -1 {
                newTotalItems := dp[i-packSize].totalItems + packSize
                newTotalPacks := dp[i-packSize].totalPacks + 1
                
                // Check if this is a better solution
                if dp[i].totalItems == -1 || 
                   (i >= orderSize && newTotalItems < dp[i].totalItems) ||
                   (i >= orderSize && newTotalItems == dp[i].totalItems && newTotalPacks < dp[i].totalPacks) {
                    // Copy the previous solution and add current pack
                    newPacks := make(map[int]int)
                    for k, v := range dp[i-packSize].packs {
                        newPacks[k] = v
                    }
                    newPacks[packSize]++
                    
                    dp[i] = dpEntry{
                        totalItems: newTotalItems,
                        totalPacks: newTotalPacks,
                        packs:      newPacks,
                    }
                }
            }
        }
    }
    
    // Find the optimal solution (least items >= orderSize)
    bestSolution := dpEntry{totalItems: -1}
    for i := orderSize; i <= maxSize; i++ {
        if dp[i].totalItems != -1 {
            if bestSolution.totalItems == -1 || 
               dp[i].totalItems < bestSolution.totalItems ||
               (dp[i].totalItems == bestSolution.totalItems && dp[i].totalPacks < bestSolution.totalPacks) {
                bestSolution = dp[i]
            }
        }
    }
    
    result := &CalculationResponse{
        Packs:      bestSolution.packs,
        TotalItems: bestSolution.totalItems,
        TotalPacks: bestSolution.totalPacks,
    }
    
    // Cache the result
    pc.cache[cacheKey] = result
    
    return result
}

// dpEntry represents an entry in the dynamic programming table
type dpEntry struct {
    totalItems int
    totalPacks int
    packs      map[int]int
}

// HTTP Handlers

// handleCalculate handles pack calculation requests
func (pc *PackCalculator) handleCalculate(w http.ResponseWriter, r *http.Request) {
    // Enable CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req CalculationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    result := pc.CalculatePacks(req.Items)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// handleUpdatePackSizes handles pack size updates
func (pc *PackCalculator) handleUpdatePackSizes(w http.ResponseWriter, r *http.Request) {
    // Enable CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var update PackSizeUpdate
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate pack sizes
    for _, size := range update.PackSizes {
        if size <= 0 {
            http.Error(w, "Pack sizes must be positive", http.StatusBadRequest)
            return
        }
    }
    
    pc.SetPackSizes(update.PackSizes)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success":   true,
        "packSizes": update.PackSizes,
    })
}

// handleGetPackSizes returns current pack sizes
func (pc *PackCalculator) handleGetPackSizes(w http.ResponseWriter, r *http.Request) {
    // Enable CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    sizes := pc.GetPackSizes()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "packSizes": sizes,
    })
}

// handleIndex serves a simple HTML UI
func handleIndex(w http.ResponseWriter, r *http.Request) {
    html := `
<!DOCTYPE html>
<html>
<head>
    <title>Order Packs Calculator</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
        }
        .section {
            background: white;
            padding: 20px;
            margin: 20px 0;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        input[type="text"], input[type="number"] {
            width: 100%;
            padding: 8px;
            margin: 5px 0;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            margin: 5px 0;
        }
        button:hover {
            background-color: #45a049;
        }
        .result {
            margin-top: 20px;
            padding: 15px;
            background-color: #f0f0f0;
            border-radius: 4px;
        }
        .pack-input {
            display: inline-block;
            width: 80px;
            margin: 5px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 10px;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #4CAF50;
            color: white;
        }
    </style>
</head>
<body>
    <h1>Order Packs Calculator</h1>
    
    <div class="section">
        <h2>Pack Sizes</h2>
        <div id="packSizes"></div>
        <button onclick="updatePackSizes()">Submit pack sizes change</button>
    </div>
    
    <div class="section">
        <h2>Calculate packs for order</h2>
        <label>Items: <input type="number" id="orderItems" value="263" min="1"></label>
        <button onclick="calculatePacks()">Calculate</button>
        <div id="result"></div>
    </div>
    
    <script>
        let currentPackSizes = [];
        
        async function loadPackSizes() {
            const response = await fetch('/api/pack-sizes');
            const data = await response.json();
            currentPackSizes = data.packSizes || [];
            renderPackSizes();
        }
        
        function renderPackSizes() {
            const container = document.getElementById('packSizes');
            container.innerHTML = currentPackSizes.map((size, index) => 
                '<input type="number" class="pack-input" id="pack' + index + '" value="' + size + '">'
            ).join('') + '<button onclick="addPackSize()">+</button>';
        }
        
        function addPackSize() {
            currentPackSizes.push(100);
            renderPackSizes();
        }
        
        async function updatePackSizes() {
            const sizes = [];
            for (let i = 0; i < currentPackSizes.length; i++) {
                const value = parseInt(document.getElementById('pack' + i).value);
                if (value > 0) sizes.push(value);
            }
            
            const response = await fetch('/api/pack-sizes', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({packSizes: sizes})
            });
            
            if (response.ok) {
                alert('Pack sizes updated successfully!');
                loadPackSizes();
            }
        }
        
        async function calculatePacks() {
            const items = parseInt(document.getElementById('orderItems').value);
            
            const response = await fetch('/api/calculate', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({items: items})
            });
            
            const data = await response.json();
            
            let html = '<h3>Result</h3><table><tr><th>Pack</th><th>Quantity</th></tr>';
            for (const [size, quantity] of Object.entries(data.packs || {})) {
                html += '<tr><td>' + size + '</td><td>' + quantity + '</td></tr>';
            }
            html += '</table>';
            html += '<p><strong>Total Items:</strong> ' + data.totalItems + '</p>';
            html += '<p><strong>Total Packs:</strong> ' + data.totalPacks + '</p>';
            
            document.getElementById('result').innerHTML = html;
        }
        
        // Load pack sizes on page load
        loadPackSizes();
    </script>
</body>
</html>
    `
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
}

func main() {
    // Initialize calculator
    calculator := NewPackCalculator()
    
    // Check for custom pack sizes from environment or command line
    if packSizesEnv := os.Getenv("PACK_SIZES"); packSizesEnv != "" {
        sizes := []int{}
        for _, s := range strings.Split(packSizesEnv, ",") {
            if size, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
                sizes = append(sizes, size)
            }
        }
        if len(sizes) > 0 {
            calculator.SetPackSizes(sizes)
        }
    }
    
    // Set up routes
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/api/calculate", calculator.handleCalculate)
    http.HandleFunc("/api/pack-sizes", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            calculator.handleGetPackSizes(w, r)
        } else if r.Method == "POST" {
            calculator.handleUpdatePackSizes(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Printf("Initial pack sizes: %v", calculator.GetPackSizes())
    
    // Test edge case on startup
    log.Println("Testing edge case: Pack Sizes: [23, 31, 53], Amount: 500000")
    calculator.SetPackSizes([]int{23, 31, 53})
    result := calculator.CalculatePacks(500000)
    log.Printf("Edge case result: %+v", result)
    
    // Reset to default pack sizes
    calculator.SetPackSizes([]int{250, 500, 1000, 2000, 5000})
    
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}