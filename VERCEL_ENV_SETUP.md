# Quick Fix for CORS and API URL Errors

## The Problem
1. **CORS Error**: Backend is blocking requests from Vercel
2. **Wrong API URL**: Frontend is calling `localhost:8080` instead of your backend URL

## Solution Steps

### Step 1: Deploy Backend (if not already done)

The CORS fix is already in the code, but you need to deploy the backend first.

**Railway Deployment:**
1. Go to https://railway.app
2. New Project → Deploy from GitHub
3. Select `compatiblah` repo
4. Add Service → Empty Service
5. Settings:
   - Root Directory: `backend`
   - Build Command: `go mod download`
   - Start Command: `go run main.go`
6. Variables:
   - `GEMINI_API_KEY` = your API key
7. Deploy → Copy the URL (e.g., `https://compatiblah-production.up.railway.app`)

### Step 2: Set Environment Variable in Vercel

1. Go to Vercel Dashboard → Your Project → Settings → Environment Variables
2. Add new variable:
   - **Key**: `VITE_API_URL`
   - **Value**: `https://your-backend-url.railway.app` (use the URL from Step 1)
   - **Environments**: Check all (Production, Preview, Development)
3. Click Save

### Step 3: Redeploy Frontend

1. In Vercel Dashboard → Deployments
2. Click the three dots (⋯) on latest deployment
3. Click "Redeploy"
4. Wait for deployment to complete

### Step 4: Verify

1. Open your Vercel app URL
2. Open browser console (F12)
3. Try submitting a form
4. Should see API calls going to your Railway backend URL (not localhost)

## If Still Having Issues

### Check Backend Logs
- Railway Dashboard → Your Service → Logs
- Should see "Server starting on port..." message
- Check for any errors

### Test Backend Directly
- Visit: `https://your-backend-url.railway.app/health`
- Should return: `{"status":"ok"}`

### Verify Environment Variable
- In Vercel → Settings → Environment Variables
- Make sure `VITE_API_URL` is set and redeploy

## Notes

- The backend now allows **all origins** by default (safe for production)
- You can restrict it by setting `CORS_ORIGINS` environment variable in Railway if needed
- The frontend will use `VITE_API_URL` if set, otherwise defaults to `localhost:8080` (for local dev)

