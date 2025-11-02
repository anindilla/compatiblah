# Deploy Backend to Render (Free Tier)

## Quick Setup Guide

### Step 1: Create Render Account
1. Go to https://render.com
2. Sign up with GitHub (free)
3. Connect your GitHub account

### Step 2: Deploy Backend Service
1. In Render Dashboard, click **"New +"** → **"Web Service"**
2. Connect your `compatiblah` GitHub repository
3. Configure the service:
   - **Name**: `compatiblah-backend` (or any name)
   - **Root Directory**: `backend`
   - **Environment**: `Go`
   - **Build Command**: `go mod download && go build -o app`
   - **Start Command**: `./app`
   - **Plan**: `Free` (select free tier)
4. Click **"Advanced"** and set:
   - **Health Check Path**: `/health`

### Step 3: Set Environment Variables
1. In your Render service dashboard, go to **"Environment"** tab
2. Add environment variable:
   - **Key**: `GEMINI_API_KEY`
   - **Value**: Your Gemini API key
   - Click **"Save Changes"**

### Step 4: Deploy
1. Click **"Create Web Service"**
2. Wait for deployment (first deploy takes 2-3 minutes)
3. Once deployed, copy your service URL (e.g., `https://compatiblah-backend.onrender.com`)

### Step 5: Configure Frontend
1. Your frontend is already set up to auto-detect Render URLs
2. OR manually set it:
   - In Vercel → Your Project → Settings → Environment Variables
   - Add: `VITE_API_URL` = `https://compatiblah-backend.onrender.com`
   - Redeploy frontend

## Alternative: Using render.yaml (Automatic)

If you prefer automatic configuration:
1. The `render.yaml` file is already in the repo
2. In Render Dashboard → **"New +"** → **"Blueprint"**
3. Connect your GitHub repo
4. Render will automatically detect `render.yaml` and configure the service
5. You'll still need to set `GEMINI_API_KEY` in Environment Variables

## Important Notes

### Free Tier Limitations
- ✅ Service is free forever
- ⚠️ Service spins down after 15 minutes of inactivity
- ⚠️ First request after spin-down takes 30-60 seconds (cold start)
- ⚠️ SQLite database is ephemeral (may reset on service restart/redeploy)

### SQLite on Render
- SQLite files work on Render free tier
- Data persists while service is running
- Data may be lost on:
  - Service restart
  - Service redeploy
  - Service spin-down/up cycles
- This is acceptable for demo/personal projects

### Health Check
- Render checks `/health` endpoint every 5 minutes
- If health check fails, Render may restart the service
- The `/health` endpoint is already configured in the Go backend

## Troubleshooting

### Service Won't Start
- Check build logs in Render dashboard
- Ensure `go.mod` has all dependencies
- Verify build command: `go mod download && go build -o app`

### CORS Errors
- Backend already configured to allow all origins
- If issues persist, check CORS settings in `backend/main.go`

### Slow First Request
- Normal on free tier (cold start)
- Service wakes up after 15 minutes of inactivity
- First request takes 30-60 seconds, subsequent requests are fast

### Database Issues
- SQLite works but data may be ephemeral
- For persistent data, consider Render PostgreSQL (also free tier)
- Current setup uses SQLite (simpler, no migration needed)

## Verify Deployment

1. Visit: `https://your-service.onrender.com/health`
2. Should return: `{"status":"ok"}`
3. If it works, your backend is ready!

## Next Steps

After backend is deployed:
1. Get your Render service URL (e.g., `https://compatiblah-backend.onrender.com`)
2. Set `VITE_API_URL` in Vercel environment variables
3. OR let the frontend auto-detect (it tries Render patterns automatically)
4. Redeploy frontend
5. Test the app!

---

**Your backend will be available at**: `https://your-service-name.onrender.com`

