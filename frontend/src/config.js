// Configuration helper for API URL
// This automatically detects the environment and uses the appropriate backend URL

let cachedApiUrl = null;

export async function getApiUrl() {
  // Return cached value if already loaded
  if (cachedApiUrl) {
    return cachedApiUrl;
  }

  // Priority order:
  // 1. Window override (for manual testing)
  // 2. Config file (runtime config)
  // 3. Vite env var (build-time)
  // 4. Auto-detect based on hostname
  
  const windowOverride = window.__API_URL__;
  if (windowOverride) {
    cachedApiUrl = windowOverride;
    return cachedApiUrl;
  }

  // Try to load from config.json (runtime config)
  try {
    const response = await fetch('/config.json?' + Date.now());
    if (response.ok) {
      const config = await response.json();
      if (config.apiUrl && config.apiUrl.trim() !== '') {
        cachedApiUrl = config.apiUrl.trim();
        return cachedApiUrl;
      }
    }
  } catch (e) {
    // Config file not found or invalid, continue to next option
  }

  // Check Vite env var (build-time)
  const viteApiUrl = import.meta.env.VITE_API_URL;
  if (viteApiUrl && viteApiUrl.trim() !== '') {
    cachedApiUrl = viteApiUrl.trim();
    return cachedApiUrl;
  }

  // If we're in production (Vercel) and still no URL configured, show error
  const hostname = window.location.hostname;
  const isProduction = hostname.includes('vercel.app') || hostname.includes('vercel.com');
  
  if (isProduction && !hostname.includes('localhost') && !hostname.includes('127.0.0.1')) {
    // Production but no backend URL configured
    console.error('❌ Backend API URL not configured!');
    console.error('📝 Please do ONE of the following:');
    console.error('   1. Set VITE_API_URL in Vercel environment variables and redeploy');
    console.error('   2. Edit frontend/public/config.json with your backend URL and commit');
    console.error('   3. Set window.__API_URL__ in browser console (temporary fix)');
    
    // Don't use localhost in production - it will fail anyway
    // Return empty string so fetch will fail with a clear error
    cachedApiUrl = '';
    return cachedApiUrl;
  }

  // Default: localhost for local development
  cachedApiUrl = 'http://localhost:8080';
  return cachedApiUrl;
}

// Synchronous version for immediate use (returns default, will be updated async)
export function getApiUrlSync() {
  return cachedApiUrl || import.meta.env.VITE_API_URL || 'http://localhost:8080';
}

