# Vercel Deployment Guide

## Current Issue
You're getting a 404 error because Vercel needs to know:
1. Where the frontend code is located
2. How to build it
3. Where the backend API is (for production)

## Fix Steps

### Option 1: Deploy from Frontend Folder (Recommended)

1. **In Vercel Dashboard:**
   - Go to your project settings
   - Change "Root Directory" to `frontend`
   - Update "Build Command" to: `npm install && npm run build`
   - Update "Output Directory" to: `dist`
   - Save and redeploy

### Option 2: Use vercel.json (Current Setup)

The `vercel.json` file I created should work, but you need to:

1. **In Vercel Dashboard:**
   - Go to Settings → General
   - Make sure "Root Directory" is set to: `.` (root of repo)
   - The vercel.json should handle the rest

2. **Set Environment Variable:**
   - Go to Settings → Environment Variables
   - Add: `VITE_API_URL` = `https://your-backend-url.com`
   - (You'll need to deploy the backend first - see below)

### Deploy Backend

The backend needs to be deployed separately. Here's a quick guide:

#### Railway (Easiest)
1. Go to https://railway.app
2. New Project → Deploy from GitHub repo
3. Select your repo
4. Add Service → Select backend folder
5. Set environment variable: `GEMINI_API_KEY`
6. Deploy
7. Copy the URL (e.g., `https://your-app.railway.app`)
8. Use this URL for `VITE_API_URL` in Vercel

#### Alternative: Render
1. Go to https://render.com
2. New Web Service
3. Connect GitHub repo
4. Root Directory: `backend`
5. Build Command: `go mod download`
6. Start Command: `go run main.go`
7. Set environment variable: `GEMINI_API_KEY`
8. Deploy and copy URL

### After Backend is Deployed

1. Go to Vercel project settings
2. Environment Variables → Add:
   - Name: `VITE_API_URL`
   - Value: `https://your-backend-url.railway.app` (or your backend URL)
   - Environment: Production, Preview, Development (select all)
3. Redeploy your frontend

### Verify It Works

After redeploying, your app should work! The frontend will call your backend API using the `VITE_API_URL` environment variable.

## Troubleshooting

- **404 Error**: Make sure Root Directory is set correctly in Vercel
- **API Calls Fail**: Check that `VITE_API_URL` is set and backend is deployed
- **Build Fails**: Make sure you're building from the `frontend` directory
