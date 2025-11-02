# Compatiblah 🧡

A compatibility assessment web application that analyzes compatibility between two people as friends, coworkers, and partners using AI-powered insights.

## 🚀 Features

- **AI-Powered Analysis**: Uses Google Gemini AI to provide comprehensive compatibility assessments
- **Three Contexts**: Evaluates compatibility as friends, coworkers, and partners
- **Structured Results**: Detailed analysis with bullet points and subcategories
- **Privacy-First**: Only stores assessment results, no personal data
- **Responsive Design**: Beautiful orange-themed UI that works on all devices
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

### Backend Deployment Options

The backend needs to be deployed separately. Options include:

1. **Railway** (Recommended)
   - Connect your GitHub repo
   - Set `GEMINI_API_KEY` environment variable
   - Deploy the `backend` folder

2. **Render**
   - Create a new Web Service
   - Point to the `backend` folder
   - Set environment variables

3. **Fly.io**
   - Install Fly CLI
   - Run `fly launch` in the backend directory
   - Set secrets: `fly secrets set GEMINI_API_KEY=your_key`

4. **Heroku**
   - Create Procfile in backend: `web: go run main.go`
   - Set config vars: `GEMINI_API_KEY`

**After deploying the backend**, update the `VITE_API_URL` in your Vercel project settings to point to your backend URL.

## 🔐 Environment Variables

### Frontend (Vercel)
- `VITE_API_URL`: Backend API URL (e.g., `https://your-backend.railway.app`)

### Backend
- `GEMINI_API_KEY`: Your Google Gemini API key
- `PORT`: Server port (default: 8080)

## 📁 Project Structure

```
compatiblah/
├── backend/
│   ├── db/           # Database operations
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Data models
│   ├── services/     # Business logic (Gemini API)
│   └── main.go       # Entry point
├── frontend/
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── App.vue      # Main app component
│   │   └── main.js      # Entry point
│   └── public/          # Static assets
└── vercel.json       # Vercel configuration
```

## 🎨 Features in Detail

- **MBTI Compatibility**: Analyzes personality types for compatibility
- **Social Media Integration**: Optional social media profile analysis
- **Additional Parameters**: Support for zodiac signs, DISC, and custom parameters
- **Detailed Explanations**: Structured analysis with subcategories and bullet points
- **Star Ratings**: Visual 5-star rating system for each compatibility context

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is open source and available under the MIT License.

## 👤 Author

Vibe-coded by [dilleuh](https://anindilla.com)

---

Made with ❤️ and ☕
