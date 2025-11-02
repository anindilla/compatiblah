# 🚨 CRITICAL: Fix Vercel Environment Variable

## The Problem
Your app is still calling `localhost:8080` because the `VITE_API_URL` environment variable is **NOT set in Vercel**, or it wasn't set **BEFORE the last build**.

## ⚠️ IMPORTANT: Vite Environment Variables Must Be Set BEFORE Build

Vite embeds environment variables **at build time**, not runtime. If you set the env var after building, you MUST redeploy.

## ✅ Step-by-Step Fix

### 1. Get Your Backend URL
First, make sure your backend is deployed and you have the URL (e.g., `https://compatiblah-production.up.railway.app`)

### 2. Set Environment Variable in Vercel

1. Go to **Vercel Dashboard** → Your Project (`compatiblah`)
2. Click **Settings** (gear icon)
3. Click **Environment Variables** (left sidebar)
4. Click **Add New**
5. Fill in:
   - **Key**: `VITE_API_URL`
   - **Value**: `https://your-backend-url.railway.app` (replace with your actual backend URL)
   - **Environment**: Check **ALL THREE**:
     - ☑️ Production
     - ☑️ Preview  
     - ☑️ Development
6. Click **Save**

### 3. **IMPORTANT: Redeploy Immediately**

1. Go to **Deployments** tab
2. Find your latest deployment
3. Click the **three dots (⋯)** menu
4. Click **Redeploy**
5. **Select "Use existing Build Cache"** = OFF (unchecked)
6. Click **Redeploy**

⚠️ **You MUST redeploy after setting the environment variable!**

### 4. Verify It Worked

After redeploy completes:

1. Open your Vercel app URL
2. Open Browser Console (F12 → Console tab)
3. Look for: `🔧 API Configuration:` log
4. Check that `final` shows your backend URL (NOT localhost:8080)
5. Try submitting a form

## 🔍 Debugging

If it still doesn't work:

### Check Environment Variable
1. In Vercel → Settings → Environment Variables
2. Verify `VITE_API_URL` is listed and has the correct value
3. Make sure it's enabled for Production environment

### Check Build Logs
1. In Vercel → Deployments → Click on latest deployment
2. Check Build Logs
3. Look for any errors or warnings

### Test Backend Directly
1. Visit: `https://your-backend-url.railway.app/health`
2. Should return: `{"status":"ok"}`
3. If this fails, your backend isn't deployed correctly

### Manual Runtime Override (Temporary)
If you need to test immediately, you can temporarily add this to your browser console on the deployed site:

```javascript
window.__API_URL__ = 'https://your-backend-url.railway.app';
location.reload();
```

This will use the runtime override, but you should still fix the environment variable properly.

## 📋 Checklist

- [ ] Backend is deployed and accessible
- [ ] `VITE_API_URL` is set in Vercel with correct backend URL
- [ ] Environment variable is enabled for Production
- [ ] Frontend has been **redeployed** after setting env var
- [ ] Build logs show no errors
- [ ] Browser console shows correct API URL (not localhost)

## 🆘 Still Not Working?

1. **Double-check the backend URL**: Make sure it's accessible and returns `{"status":"ok"}` at `/health`
2. **Check CORS**: Backend should allow all origins (already fixed in code)
3. **Verify env var name**: Must be exactly `VITE_API_URL` (case-sensitive)
4. **Check build output**: Make sure the build completed successfully

