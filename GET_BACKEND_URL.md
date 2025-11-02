# Get Your Render Backend URL

Since you've already set up Render, here's how to get your backend URL:

## In Render Dashboard:

1. Go to https://dashboard.render.com
2. Click on your **Web Service** (should be named something like `compatiblah-backend`)
3. Look at the top of the page - you'll see your service URL like:
   - `https://compatiblah-backend.onrender.com`
   - OR `https://your-service-name.onrender.com`
4. **Copy this URL** - this is your backend URL!

## Test Your Backend:

Visit: `https://your-service.onrender.com/health`

Should return: `{"status":"ok"}`

If it works, your backend is live! 🎉

## Configure Frontend:

Once you have your backend URL:

**Option 1: Auto-detect (Recommended)**
- The frontend will automatically try to detect your Render backend
- Just reload your Vercel site after backend is deployed

**Option 2: Manual Config**
- Open your Vercel site
- If the config dialog appears, paste your Render URL
- OR set in browser console:
  ```javascript
  window.__API_URL__ = 'https://your-service.onrender.com';
  location.reload();
  ```

**Option 3: Vercel Environment Variable**
- Vercel Dashboard → Your Project → Settings → Environment Variables
- Add: `VITE_API_URL` = `https://your-service.onrender.com`
- Redeploy frontend

## Troubleshooting:

**Backend not responding?**
- Check Render dashboard → Logs tab
- Make sure service is "Live" (not sleeping)
- First request after sleep takes 30-60 seconds (normal for free tier)

**CORS errors?**
- Backend already configured to allow all origins
- Should work automatically

**Need to check service status?**
- Render Dashboard → Your Service → Should show "Live" status
- Check "Logs" tab for any errors

