# Order Packs Calculator

A Go-based web service that calculates the optimal number of packs needed to fulfill customer orders.

## Problem Statement

The service solves the pack optimization problem where:
1. Only whole packs can be sent (packs cannot be broken open)
2. Send out the least amount of items to fulfill the order
3. Send out as few packs as possible (within the constraint of rule 2)

## Features

- RESTful API for pack calculation
- Dynamic pack size configuration
- Web UI for easy interaction
- Caching for performance optimization
- Docker support for easy deployment

## Quick Start

### Using Docker (Recommended)

```bash
# Build the Docker image
make docker-build

# Run the container
make docker-run

# Or use docker-compose
docker-compose up
```

The application will be available at `http://localhost:8080`

### Local Development

```bash
# Install dependencies
go mod tidy

# Run tests
make test

# Run the application
make run

# Or directly with Go
go run main.go
```

## API Endpoints

### Calculate Packs
- **POST** `/api/calculate`
- **Body**: `{"items": 263}`
- **Response**: 
```json
{
  "packs": {"250": 1, "500": 1},
  "totalItems": 750,
  "totalPacks": 2
}
```

### Get Pack Sizes
- **GET** `/api/pack-sizes`
- **Response**: `{"packSizes": [5000, 2000, 1000, 500, 250]}`

### Update Pack Sizes
- **POST** `/api/pack-sizes`
- **Body**: `{"packSizes": [250, 500, 1000, 2000, 5000]}`

## Environment Variables

- `PORT`: Server port (default: 8080)
- `PACK_SIZES`: Comma-separated pack sizes (default: 250,500,1000,2000,5000)

## Testing

Run unit tests:
```bash
make test

# With coverage
make test-coverage
```

## Edge Case Verification

The algorithm correctly handles the specified edge case:
- Pack Sizes: [23, 31, 53]
- Amount: 500,000
- Expected Output: {23: 2, 31: 7, 53: 9429}

## Architecture

The application uses:
- Dynamic Programming for optimal pack calculation
- Caching for performance optimization
- Thread-safe operations with sync.RWMutex
- Clean separation of concerns

## Deployment

### Render (Recommended)

1. **Push to GitHub**:
```bash
git init
git add .
git commit -m "feat: complete pack calculator implementation"
git branch -M main
git remote add origin https://github.com/yourusername/pack-calculator.git
git push -u origin main
```

2. **Deploy on Render**:
   - Go to [render.com](https://render.com)
   - Connect your GitHub repository
   - Choose "Web Service"
   - Use Docker environment
   - Set environment variables if needed

3. **Manual Deployment**:
```bash
# The render.yaml file will automatically configure the deployment
```

### Docker Hub

```bash
docker build -t yourusername/pack-calculator .
docker push yourusername/pack-calculator
```

### Local Testing

```bash
# Test locally first
go test -v
go run main.go

# Test with Docker
make docker-build
make docker-run
```

## Git Best Practices

Example commit history:
```bash
git add go.mod
git commit -m "build: fix Go version from 1.25.1 to 1.21"

git add Dockerfile
git commit -m "build: fix Dockerfile Go version"

git add Makefile
git commit -m "build: fix Makefile indentation with tabs"

git add README.md
git commit -m "docs: add comprehensive documentation"

git add main.go
git commit -m "feat: implement pack calculator with dynamic programming"

git add main_test.go
git commit -m "test: add comprehensive unit tests"

git add render.yaml
git commit -m "deploy: add Render deployment configuration"
```

## Live Demo

Once deployed on Render, your app will be available at:
`https://your-app-name.onrender.com`

Test the edge case:
```bash
curl -X POST https://your-app-name.onrender.com/api/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"packSizes": [23, 31, 53]}'

curl -X POST https://your-app-name.onrender.com/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"items": 500000}'
```

## License

MIT