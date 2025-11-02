// Configuration helper for API URL
// This automatically detects the environment and uses the appropriate backend URL

let cachedApiUrl = null;

export async function getApiUrl() {
  // Return cached value if already loaded
  if (cachedApiUrl) {
    return cachedApiUrl;
  }

  // Priority order:
  // 1. localStorage (user-set, persistent)
  // 2. Window override (for manual testing)
  // 3. Config file (runtime config)
  // 4. Vite env var (build-time)
  // 5. Auto-detect common patterns
  
  // Check localStorage first (persists across reloads)
  const storedUrl = localStorage.getItem('backend_api_url');
  if (storedUrl && storedUrl.trim() !== '' && !storedUrl.includes('localhost')) {
    cachedApiUrl = storedUrl.trim();
    return cachedApiUrl;
  }
  
  const windowOverride = window.__API_URL__;
  if (windowOverride && windowOverride.trim() !== '') {
    cachedApiUrl = windowOverride.trim();
    localStorage.setItem('backend_api_url', cachedApiUrl);
    return cachedApiUrl;
  }

  // Try to load from config.json (runtime config)
  try {
    const response = await fetch('/config.json?' + Date.now());
    if (response.ok) {
      const config = await response.json();
      if (config.apiUrl && config.apiUrl.trim() !== '' && !config.apiUrl.includes('YOUR-BACKEND')) {
        cachedApiUrl = config.apiUrl.trim();
        localStorage.setItem('backend_api_url', cachedApiUrl);
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
    localStorage.setItem('backend_api_url', cachedApiUrl);
    return cachedApiUrl;
  }

  // Auto-detect: Try common Railway URL patterns based on hostname
  const hostname = window.location.hostname;
  const isProduction = hostname.includes('vercel.app') || hostname.includes('vercel.com');
  
  if (isProduction && !hostname.includes('localhost') && !hostname.includes('127.0.0.1')) {
    // Try common Railway patterns
    const projectName = hostname.replace(/\.vercel\.app.*/, '').replace(/-/g, '');
    const possibleUrls = [
      `https://compatiblah-production.up.railway.app`,
      `https://compatiblah.railway.app`,
      `https://${projectName}-production.up.railway.app`,
      `https://${projectName}.railway.app`,
    ];
    
    // Return first pattern (user will see config dialog if wrong)
    cachedApiUrl = possibleUrls[0];
    console.warn('⚠️ Using auto-detected backend URL:', cachedApiUrl);
    console.warn('📝 If this is wrong, a configuration dialog will appear.');
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

