# CORS Troubleshooting

If you're still seeing CORS errors after the fix:

## Verify Backend is Deployed

1. **Check Render Deployment Status**
   - Go to Render Dashboard → Your Service
   - Check "Events" tab - should show latest deployment
   - Verify it deployed after the CORS fix commit

2. **Test Backend Directly**
   - Open browser console
   - Run: `fetch('https://compatiblah.onrender.com/health')`
   - Check Response Headers - should see `Access-Control-Allow-Origin: *`

3. **Test OPTIONS Preflight**
   - Open browser console
   - Run:
     ```javascript
     fetch('https://compatiblah.onrender.com/api/assess', {
       method: 'OPTIONS',
       headers: { 'Content-Type': 'application/json' }
     }).then(r => console.log('OPTIONS:', r.status, r.headers.get('Access-Control-Allow-Origin')))
     ```
   - Should return status 204 and show CORS headers

## Manual Render Redeploy

If auto-deploy didn't happen:
1. Render Dashboard → Your Service
2. **Manual Deploy** tab
3. Click **"Deploy latest commit"**
4. Wait for build to complete

## Check Render Logs

1. Render Dashboard → Your Service → **Logs** tab
2. Look for:
   - Build errors
   - Runtime errors
   - Any CORS-related messages

## Alternative: Manual CORS Headers

If CORS middleware still doesn't work, we can add manual middleware. Let me know if you need this.

## Verify Service is Live

- Render Dashboard → Service status should be **"Live"** (not "Sleeping" or "Build Failed")
- If "Sleeping", first request will wake it up (takes 30-60 seconds)

