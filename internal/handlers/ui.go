package handlers

import (
	"net/http"
)

// handleIndex serves the web UI
func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `<!DOCTYPE html>
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
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
