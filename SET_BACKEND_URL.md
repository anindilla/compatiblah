# 🚀 Quick Fix: Set Your Backend URL

## Option 1: Edit config.json (Easiest - No Redeploy Needed!)

1. **Edit this file**: `frontend/public/config.json`
2. **Replace** `YOUR-BACKEND-URL-HERE` with your actual Railway backend URL
3. **Example**: `"apiUrl": "https://compatiblah-production.up.railway.app"`
4. **Commit and push** - Vercel will auto-redeploy

## Option 2: Browser Console (Instant Fix - No Code Changes!)

1. **Open your deployed site**: https://compatiblah.vercel.app
2. **Open browser console** (F12 → Console tab)
3. **Paste and run**:
   ```javascript
   window.__API_URL__ = 'https://YOUR-BACKEND-URL.railway.app';
   location.reload();
   ```
4. **Replace** `YOUR-BACKEND-URL` with your actual Railway URL
5. **Done!** The app will now work

## Option 3: Vercel Environment Variable

1. Go to Vercel → Your Project → Settings → Environment Variables
2. Add: `VITE_API_URL` = `https://YOUR-BACKEND-URL.railway.app`
3. Redeploy (WITHOUT build cache)
4. Done!

---

**Your backend URL should look like**: `https://compatiblah-production.up.railway.app`

Find it in Railway dashboard → Your Service → Settings → Domain

