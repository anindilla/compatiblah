# Backend API Endpoints

Your backend is running at: `https://compatiblah.onrender.com`

## Available Endpoints

### Health Check
- **GET** `/health`
- Returns: `{"status":"ok"}`
- Use this to verify backend is running

### Root
- **GET** `/`
- Returns API information and available endpoints

### Assessment Endpoints

1. **Create Assessment**
   - **POST** `/api/assess`
   - Body: JSON with `person1` and `person2` data
   - Returns: Assessment results with scores and explanations

2. **Get Assessment by ID**
   - **GET** `/api/assessment/:id`
   - Returns: Specific assessment by ID

3. **Get All Assessments**
   - **GET** `/api/assessments`
   - Returns: List of all assessments (limited)

## Testing

### Test Health Check
```bash
curl https://compatiblah.onrender.com/health
```
Should return: `{"status":"ok"}`

### Test Root
```bash
curl https://compatiblah.onrender.com/
```
Should return API information

### Test Assessment (from browser console on Vercel site)
The frontend automatically calls `/api/assess` when you submit the form.

## Troubleshooting

**404 on root `/`**
- This is now fixed - root endpoint returns API info

**404 on `/api/assess`**
- Check that service is "Live" in Render dashboard
- Check Render logs for errors
- Verify GEMINI_API_KEY is set

**CORS errors**
- Backend is configured to allow all origins
- Should work automatically with Vercel frontend

