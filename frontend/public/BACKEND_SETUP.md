# Backend URL Configuration

## Automatic Detection
The app will automatically try to detect your backend URL. However, if it can't find it, you need to configure it.

## Option 1: Edit config.json (Easiest - No Redeploy Needed!)

1. In your Vercel deployment, go to your project files
2. Edit `/public/config.json` 
3. Set your backend URL:
```json
{
  "apiUrl": "https://your-backend-url.railway.app"
}
```
4. Commit and push - Vercel will auto-redeploy

## Option 2: Environment Variable (Requires Redeploy)

Set `VITE_API_URL` in Vercel environment variables, then redeploy WITHOUT build cache.

## Option 3: Quick Fix via Browser Console

Open browser console on your deployed site and run:
```javascript
window.__API_URL__ = 'https://your-backend-url.railway.app';
location.reload();
```

