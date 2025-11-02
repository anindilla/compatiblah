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

  // Auto-detect: if on Vercel, assume backend is on same domain with /api proxy
  // OR if hostname is vercel, use a default backend URL pattern
  const hostname = window.location.hostname;
  
  if (hostname.includes('vercel.app') || hostname.includes('vercel.com')) {
    // If deployed on Vercel but no backend URL set, we need to detect it
    // For now, use a common pattern or prompt user
    // Default to trying the same domain with /api prefix (if proxied)
    const currentOrigin = window.location.origin;
    
    // Check if we're in production (not localhost)
    if (!hostname.includes('localhost') && !hostname.includes('127.0.0.1')) {
      // In production, we'll try to auto-detect
      // For Railway deployments, we'll need to set this manually via config.json
      // But let's provide a helpful error message
      console.warn('⚠️ Backend API URL not configured. Please set it in /config.json or VITE_API_URL env var.');
      // Don't fail completely, try a reasonable default
      cachedApiUrl = currentOrigin.replace(/^https?:\/\//, '').replace(/\.vercel\.app.*/, '') + '.railway.app';
      if (!cachedApiUrl.startsWith('http')) {
        cachedApiUrl = 'https://' + cachedApiUrl;
      }
      console.warn('⚠️ Using auto-detected backend URL:', cachedApiUrl);
      return cachedApiUrl;
    }
  }

  // Default: localhost for local development
  cachedApiUrl = 'http://localhost:8080';
  return cachedApiUrl;
}

// Synchronous version for immediate use (returns default, will be updated async)
export function getApiUrlSync() {
  return cachedApiUrl || import.meta.env.VITE_API_URL || 'http://localhost:8080';
}

