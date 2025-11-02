#!/bin/bash

# Compatiblah Server Startup Script

echo "🚀 Starting Compatiblah servers..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first:"
    echo "   brew install go"
    echo ""
    echo "   Or download from: https://go.dev/dl/"
    exit 1
fi

# Set Gemini API Key
export GEMINI_API_KEY=AIzaSyATaLbXQcPDQbYbj-j9Qls-8mLwiFuU9Go

# Start backend server
echo "📡 Starting backend server..."
cd "$(dirname "$0")"
go run backend/main.go &
BACKEND_PID=$!
echo "   Backend PID: $BACKEND_PID"
echo "   Backend running on http://localhost:8080"

# Wait a moment for backend to start
sleep 2

# Start frontend server
echo "🎨 Starting frontend server..."
cd frontend
npm run dev &
FRONTEND_PID=$!
echo "   Frontend PID: $FRONTEND_PID"
echo "   Frontend running on http://localhost:5173"

echo ""
echo "✅ Servers started!"
echo ""
echo "📝 To stop servers, run:"
echo "   kill $BACKEND_PID $FRONTEND_PID"
echo ""
echo "🌐 Open http://localhost:5173 in your browser"

# Wait for user interrupt
trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT TERM
wait

