#!/bin/bash

echo "🚀 Starting GatewayForge AI without Docker"
echo "=========================================="
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Install with: brew install go"
    exit 1
fi
echo "✅ Go: $(go version)"

# Check Node
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Install with: brew install node"
    exit 1
fi
echo "✅ Node: $(node --version)"

# Check PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Install with: brew install postgresql@14"
    exit 1
fi
echo "✅ PostgreSQL: $(psql --version | head -1)"

echo ""
echo "🗄️ Setting up database..."

# Start PostgreSQL if not running
brew services start postgresql@14 2>/dev/null || brew services restart postgresql@14

sleep 2

# Create database
echo "Creating database 'gatewayforge'..."
createdb gatewayforge 2>/dev/null || echo "Database already exists"

# Run migrations
echo "Running database migrations..."
psql gatewayforge < backend/database/schema.sql

echo ""
echo "📦 Installing dependencies..."

# Backend dependencies
cd backend/api
echo "Installing Go dependencies..."
go mod download

cd ../..

# Frontend dependencies
cd frontend
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "=========================================="
echo "🎯 Starting services..."
echo "=========================================="
echo ""
echo "Backend will start on: http://localhost:8080"
echo "Frontend will start on: http://localhost:5173"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Start backend in background
cd backend/api
export DATABASE_URL="postgres://$(whoami)@localhost:5432/gatewayforge?sslmode=disable"
export PORT=8080
export CLAUDE_AUTH_MODE=desktop
export CLAUDE_MCP_ENDPOINT=http://localhost:52828
export ENV=development

echo "🔧 Starting backend API..."
go run main.go &
BACKEND_PID=$!

cd ../..

# Wait for backend to start
sleep 5

# Start frontend
cd frontend
echo "🎨 Starting frontend..."
npm run dev &
FRONTEND_PID=$!

cd ..

echo ""
echo "=========================================="
echo "✅ GatewayForge AI is running!"
echo "=========================================="
echo ""
echo "Access the platform:"
echo "  • Frontend: http://localhost:5173"
echo "  • Backend API: http://localhost:8080"
echo "  • API Health: http://localhost:8080/api/v1/health"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for user interrupt
trap "echo ''; echo '🛑 Stopping services...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT

wait
