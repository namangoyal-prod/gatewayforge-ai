.PHONY: help setup dev build test clean docker-up docker-down db-setup db-migrate

# Default target
help:
	@echo "GatewayForge AI - Make Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup          - Initial project setup (dependencies, database, etc.)"
	@echo "  make db-setup       - Create database and run migrations"
	@echo "  make db-migrate     - Run database migrations"
	@echo ""
	@echo "Development:"
	@echo "  make dev            - Start development environment (all services)"
	@echo "  make dev-backend    - Start backend API only"
	@echo "  make dev-frontend   - Start frontend only"
	@echo "  make dev-worker     - Start Temporal worker only"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up      - Start all services with Docker Compose"
	@echo "  make docker-down    - Stop all Docker Compose services"
	@echo "  make docker-rebuild - Rebuild and restart Docker services"
	@echo "  make docker-logs    - View logs from all services"
	@echo ""
	@echo "Build:"
	@echo "  make build          - Build all components"
	@echo "  make build-backend  - Build backend API"
	@echo "  make build-frontend - Build frontend"
	@echo ""
	@echo "Test:"
	@echo "  make test           - Run all tests"
	@echo "  make test-backend   - Run backend tests"
	@echo "  make test-frontend  - Run frontend tests"
	@echo ""
	@echo "Clean:"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make clean-all      - Clean everything including Docker volumes"

# Setup
setup:
	@echo "Setting up GatewayForge AI..."
	cp .env.example .env
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Creating database..."
	make db-setup
	@echo "Setup complete! Edit .env with your configuration."

# Database
db-setup:
	@echo "Creating database..."
	createdb gatewayforge || true
	@echo "Running migrations..."
	psql gatewayforge < backend/database/schema.sql
	@echo "Database setup complete!"

db-migrate:
	@echo "Running migrations..."
	psql gatewayforge < backend/database/schema.sql

# Development
dev:
	@echo "Starting development environment..."
	make docker-up

dev-backend:
	@echo "Starting backend API..."
	cd backend/api && go run main.go

dev-frontend:
	@echo "Starting frontend..."
	cd frontend && npm run dev

dev-worker:
	@echo "Starting Temporal worker..."
	cd backend/orchestration && go run worker/main.go

# Docker
docker-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d
	@echo "Services started!"
	@echo "API: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@echo "Temporal UI: http://localhost:8088"
	@echo "Grafana: http://localhost:3000"
	@echo "MinIO: http://localhost:9001"

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-rebuild:
	@echo "Rebuilding Docker services..."
	docker-compose down
	docker-compose build
	docker-compose up -d

docker-logs:
	docker-compose logs -f

# Build
build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	cd backend && go build -o bin/api ./api
	cd backend && go build -o bin/worker ./orchestration/worker

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Test
test: test-backend test-frontend

test-backend:
	@echo "Running backend tests..."
	cd backend && go test ./...

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm test

# Lint
lint:
	@echo "Running linters..."
	cd backend && golangci-lint run
	cd frontend && npm run lint

# Format
format:
	@echo "Formatting code..."
	cd backend && go fmt ./...
	cd frontend && npm run format

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.cache

clean-all: clean
	@echo "Cleaning Docker volumes..."
	docker-compose down -v
	rm -rf uploads/*

# Production
deploy-prod:
	@echo "Deploying to production..."
	kubectl apply -f k8s/production/

deploy-staging:
	@echo "Deploying to staging..."
	kubectl apply -f k8s/staging/

# Monitoring
monitor:
	@echo "Opening monitoring dashboards..."
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Prometheus: http://localhost:9090"
	@echo "Temporal UI: http://localhost:8088"

# Database backup
db-backup:
	@echo "Backing up database..."
	pg_dump gatewayforge > backup-$(shell date +%Y%m%d-%H%M%S).sql
	@echo "Backup complete!"

# Database restore
db-restore:
	@echo "Restoring database from backup..."
	@read -p "Enter backup file name: " backup; \
	psql gatewayforge < $$backup

# Generate API documentation
api-docs:
	@echo "Generating API documentation..."
	cd backend && swag init -g api/main.go

# Run security scan
security-scan:
	@echo "Running security scan..."
	cd backend && gosec ./...
	cd frontend && npm audit

# Performance test
perf-test:
	@echo "Running performance tests..."
	k6 run tests/performance/load-test.js

# Claude Desktop MCP
mcp-check:
	@echo "Checking Claude Desktop MCP health..."
	cd backend/cmd/mcp-health-check && go run main.go

mcp-setup:
	@echo "Setting up Claude Desktop integration..."
	@echo "Opening setup guide..."
	open docs/CLAUDE_DESKTOP_SETUP.md || cat docs/CLAUDE_DESKTOP_SETUP.md
