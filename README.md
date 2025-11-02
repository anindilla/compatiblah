# Compatiblah

A modern web application that assesses compatibility between two people across three contexts: as friends, coworkers, and romantic partners. Powered by Google's Gemini AI.

## Features

- **Comprehensive Compatibility Assessment**: Get detailed compatibility scores and explanations for three relationship contexts
- **Dynamic Form Fields**: Add custom parameters like zodiac signs, DISC personality types, or any other information
- **Social Media Integration**: Include social media profiles for a more complete assessment
- **Assessment History**: View and share previous compatibility assessments
- **Modern, Responsive Design**: Beautiful UI that works perfectly on mobile and desktop
- **Accessible**: Built with accessibility in mind, following WCAG guidelines

## Tech Stack

- **Frontend**: Vue.js 3 + Vite + Tailwind CSS
- **Backend**: Go (Gin framework)
- **Database**: SQLite
- **AI**: Google Gemini API

## Prerequisites

- Node.js (v20.11.1 or higher recommended)
- Go (1.21 or higher)
- Google Gemini API key

## Setup

### Backend Setup

1. Navigate to the project root:
```bash
cd /Users/dilleuh/Coding/Compatiblah
```

2. Install Go dependencies:
```bash
go mod download
```

3. Set environment variables:
```bash
export GEMINI_API_KEY=AIzaSyATaLbXQcPDQbYbj-j9Qls-8mLwiFuU9Go
export PORT=8080
```

Or create a `.env` file in the backend directory (note: for production, use proper environment variable management).

4. Run the backend server:
```bash
go run backend/main.go
```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file (optional, defaults to `http://localhost:8080`):
```bash
echo "VITE_API_URL=http://localhost:8080" > .env
```

4. Run the development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

## Usage

1. Open your browser and navigate to `http://localhost:5173`
2. Fill in the information for Person 1:
   - Name (required)
   - MBTI Type (required)
   - Social Media profiles (optional)
   - Additional parameters like zodiac sign, DISC type, etc. (optional)
3. Fill in the information for Person 2
4. Click "Assess Compatibility"
5. View the results with detailed explanations and star ratings
6. Share or save your assessment using the generated ID

## API Endpoints

### POST /api/assess
Submit two people's data for compatibility assessment.

**Request Body:**
```json
{
  "person1": {
    "name": "John Doe",
    "mbti": "INTJ",
    "socialMedia": [
      {"platform": "Instagram", "handle": "@johndoe"}
    ],
    "additionalParams": [
      {"key": "Zodiac", "value": "Leo"}
    ]
  },
  "person2": {
    "name": "Jane Smith",
    "mbti": "ENFP",
    "socialMedia": [],
    "additionalParams": []
  }
}
```

**Response:**
```json
{
  "id": "uuid",
  "friend_score": 4,
  "coworker_score": 3,
  "partner_score": 5,
  "overall_score": 4,
  "friend_explanation": "...",
  "coworker_explanation": "...",
  "partner_explanation": "..."
}
```

### GET /api/assessment/:id
Retrieve a saved assessment by ID.

### GET /api/assessments
Get a list of all assessments (limited to 100 most recent).

## Project Structure

```
/
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── PersonForm.vue
│   │   │   ├── CompatibilityResults.vue
│   │   │   └── AssessmentHistory.vue
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css
│   ├── package.json
│   └── vite.config.js
├── backend/
│   ├── main.go
│   ├── handlers/
│   │   └── assessment.go
│   ├── models/
│   │   └── assessment.go
│   ├── services/
│   │   └── gemini.go
│   └── db/
│       └── database.go
├── go.mod
└── README.md
```

## Development

### Running in Development Mode

1. Start the backend server in one terminal:
```bash
go run backend/main.go
```

2. Start the frontend dev server in another terminal:
```bash
cd frontend && npm run dev
```

### Building for Production

**Frontend:**
```bash
cd frontend
npm run build
```

The built files will be in `frontend/dist/`

**Backend:**
```bash
go build -o compatiblah-server backend/main.go
```

## Environment Variables

### Backend
- `GEMINI_API_KEY`: Your Google Gemini API key (required)
- `PORT`: Server port (default: 8080)

### Frontend
- `VITE_API_URL`: Backend API URL (default: http://localhost:8080)

## License

MIT

