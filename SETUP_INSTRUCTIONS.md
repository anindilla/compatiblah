# Setup Instructions

## ✅ Frontend Status
**Frontend is currently running at: http://localhost:5173**

You can open this in your browser now! However, you'll need the backend running for full functionality.

## ⚠️ Backend Issue
The Go backend cannot start because **Go is not installed** on your system.

### To install Go:

**Option 1: Using Homebrew (Recommended)**
```bash
brew install go
```

**Option 2: Manual Installation**
1. Download Go from: https://go.dev/dl/
2. Install the package
3. Restart your terminal

### After Installing Go:

**Option A: Use the startup script**
```bash
./START_SERVERS.sh
```

**Option B: Start manually**
```bash
# Terminal 1 - Backend
export GEMINI_API_KEY=AIzaSyATaLbXQcPDQbYbj-j9Qls-8mLwiFuU9Go
go run backend/main.go

# Terminal 2 - Frontend (already running, but if you restart)
cd frontend
npm run dev
```

## Current Status
- ✅ Frontend: Running at http://localhost:5173
- ❌ Backend: Requires Go installation
- ✅ All code: Complete and ready

## Once Go is installed:
The backend will serve the API at http://localhost:8080 and the frontend will be able to make compatibility assessments.

