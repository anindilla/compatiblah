# Compatiblah 🧡

A compatibility assessment web application that analyzes compatibility between two people as friends, coworkers, and partners using AI-powered insights.

## 🚀 Features

- **AI-Powered Analysis**: Uses Google Gemini AI to provide comprehensive compatibility assessments
- **Three Contexts**: Evaluates compatibility as friends, coworkers, and partners
- **Progressive Generation**: Results appear progressively as each category is analyzed (friend → coworker → partner)
- **Structured Results**: Detailed analysis with sections, subcategories, and 2-3 concise bullet points each
- **Collapsible Sections**: Expandable/collapsible sections for better content organization
- **Sticky Navigation**: Mobile-only sticky category header that updates as you scroll
- **Privacy-First**: Only stores assessment results, no personal data
- **Responsive Design**: Beautiful orange-themed UI optimized for desktop and mobile
- **Modern Stack**: Vue.js 3 + Go backend

## 📋 Tech Stack

### Frontend
- Vue.js 3
- Vite
- Tailwind CSS
- Modern fonts (Plus Jakarta Sans, Space Grotesk)

### Backend
- Go
- Gin framework
- SQLite database
- Google Gemini API integration

## 🛠️ Local Development

### Prerequisites
- Node.js 20+
- Go 1.21+
- Gemini API Key

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/anindilla/compatiblah.git
   cd compatiblah
   ```

2. **Backend Setup**
   ```bash
   cd backend
   go mod download
   
   # Set your Gemini API key
   export GEMINI_API_KEY=your_api_key_here
   
   # Run the backend
   go run main.go
   ```
   Backend runs on `http://localhost:8080`

3. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   
   # Create .env file (optional, defaults to localhost:8080)
   echo "VITE_API_URL=http://localhost:8080" > .env.local
   
   # Run the frontend
   npm run dev
   ```
   Frontend runs on `http://localhost:5173`

## 🌐 Deployment

### Frontend (Vercel)

The frontend is configured for Vercel deployment. Simply connect your GitHub repository to Vercel and it will automatically deploy.

**Important**: You'll need to set the `VITE_API_URL` environment variable in Vercel pointing to your backend API.

### Backend Deployment (Render - Free Tier)

The backend is deployed on Render's free tier. The Go backend handles sequential category generation for progressive results.

#### Quick Deploy to Render

1. **Create Render Account**
   - Go to https://render.com and sign up with GitHub

2. **Deploy Backend**
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Configure:
     - **Root Directory**: `backend`
     - **Build Command**: `go mod download && go build -o app`
     - **Start Command**: `./app`
     - **Health Check Path**: `/health`
     - **Plan**: Free
   - Add environment variable: `GEMINI_API_KEY` = your API key
   - Deploy!

3. **Get Backend URL**
   - After deployment, copy your service URL (e.g., `https://compatiblah-backend.onrender.com`)

4. **Configure Frontend**
   - Option A: Set `VITE_API_URL` in Vercel environment variables to your Render URL
   - Option B: Let the app auto-detect (it tries Render patterns automatically)
   - Redeploy frontend

**Detailed instructions**: See [DEPLOY_RENDER.md](./DEPLOY_RENDER.md)

#### Alternative: Using render.yaml (Blueprint)

The repository includes `render.yaml` for automatic configuration:
1. In Render → "New +" → "Blueprint"
2. Connect your GitHub repo
3. Render auto-configures from `render.yaml`
4. Set `GEMINI_API_KEY` environment variable
5. Deploy!

#### Free Tier Notes

- Service is free forever
- Spins down after 15 minutes of inactivity (cold start ~30-60s)
- SQLite database works but may reset on service restart (acceptable for demo/personal projects)

## 🔐 Environment Variables

### Frontend (Vercel)
- `VITE_API_URL`: Backend API URL (optional, auto-detects Render URLs)

### Backend (Render)
- `GEMINI_API_KEY`: Your Google Gemini API key (required)
- `PORT`: Server port (automatically set by Render)
- `CORS_ORIGINS`: Not needed (uses AllowAllOrigins)

## 📁 Project Structure

```
compatiblah/
├── backend/
│   ├── db/              # Database operations (SQLite)
│   ├── handlers/         # HTTP handlers (assessments, categories)
│   ├── models/          # Data models (PersonData, Assessment, etc.)
│   ├── services/        # Business logic (Gemini API integration)
│   │   └── gemini.go    # Progressive category assessment
│   └── main.go          # Entry point with CORS middleware
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── CompatibilityResults.vue  # Results with sticky header & collapsible sections
│   │   │   ├── PersonForm.vue            # Input forms
│   │   │   └── Footer.vue                # Footer component
│   │   ├── App.vue                      # Main app with progressive API calls
│   │   ├── config.js                    # Runtime API URL detection
│   │   └── main.js                      # Entry point
│   └── public/
│       └── config.json                  # Runtime backend URL config
├── render.yaml          # Render deployment blueprint
└── vercel.json          # Vercel frontend configuration
```

## 🔄 How It Works

1. **User Input**: Enter names, MBTI types, optional social media, and additional parameters
2. **Sequential Processing**: Frontend makes 3 sequential API calls to `/api/assess/category`:
   - First: Friendship compatibility
   - Second: Workplace compatibility  
   - Third: Romance compatibility
3. **Progressive Display**: Results appear as each category completes, showing loading states per category
4. **Structured Output**: Each category returns:
   - Score (1-5 stars)
   - Explanation with 3+ sections
   - Each section has 2-3 subcategories
   - Each subcategory has 2-3 bullet points
5. **Interactive UI**: Users can collapse/expand sections, mobile users see sticky category header

## 🎨 Features in Detail

- **MBTI Compatibility**: Analyzes personality types for compatibility
- **Social Media Integration**: Optional social media profile analysis (only used if accessible)
- **Additional Parameters**: Support for zodiac signs, DISC, Enneagram, and custom parameters
- **Progressive Loading**: Categories are generated sequentially, showing results as they complete
- **Structured Explanations**: 
  - 3+ sections per category
  - 2-3 subcategories per section
  - 2-3 detailed bullet points per subcategory
- **Interactive UI**:
  - Collapsible sections (all expanded by default)
  - Mobile sticky header showing current category
  - Smooth animations and transitions
- **Star Ratings**: Visual 5-star rating system for each compatibility context
- **Overall Score**: Calculated dynamically as categories complete

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is open source and available under the MIT License.

## 👤 Author

Vibe-coded by [dilleuh](https://anindilla.com)

---

Made with ❤️ and ☕
