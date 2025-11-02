# Render Manual Setup (If Build Failed)

Since you've already created the service in Render, if the build is failing, here's how to configure it manually:

## Render Dashboard Configuration

1. **Go to your Render service dashboard**
   - Click on your service (compatiblah or compatiblah-backend)

2. **Settings Tab - Configure these:**

   **Basic Settings:**
   - **Name**: `compatiblah-backend` (or whatever you named it)
   - **Root Directory**: Leave **EMPTY** (not `backend`)
   - **Environment**: `Go`
   - **Region**: Choose closest to you
   - **Branch**: `main` (or your main branch)

   **Build & Deploy:**
   - **Build Command**: `go mod download && cd backend && go build -o app`
   - **Start Command**: `cd backend && ./app`
   - **Plan**: `Free`

   **Health Check:**
   - **Health Check Path**: `/health`

3. **Environment Tab - Add:**
   - **Key**: `GEMINI_API_KEY`
   - **Value**: Your Gemini API key
   - Click **Save Changes**

4. **Manual Deploy:**
   - Go to **Manual Deploy** tab
   - Click **Deploy latest commit**
   - Watch the build logs

## Check Build Logs

1. Click on your service → **Logs** tab
2. Look for error messages
3. Common errors:

   **"go.mod not found"**
   - Fix: Make sure Root Directory is **EMPTY** (not `backend`)
   - Go builds from root where `go.mod` is located

   **"package not found" or import errors**
   - Fix: Make sure build command includes `go mod download`
   - Try: `go mod tidy` locally, commit, and redeploy

   **"cannot find package compatiblah/backend/..."**
   - Fix: Build must run from root directory where `go.mod` is
   - The `cd backend` in build command is correct, but must start from root

## Test Build Locally

To verify it will work on Render, test locally:

```bash
cd /Users/dilleuh/Coding/Compatiblah
go mod download
cd backend
go build -o app
./app
```

If this works locally, Render should work too with the correct settings.

## After Successful Build

Once build succeeds:
1. Your backend will be live at: `https://compatiblah.onrender.com`
2. Test: Visit `https://compatiblah.onrender.com/health`
3. Should return: `{"status":"ok"}`
4. Frontend will auto-detect this URL (already configured!)

