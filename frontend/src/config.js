// Configuration helper for API URL
// Vite environment variables are embedded at build time
// For Vercel, set VITE_API_URL in environment variables BEFORE building

export function getApiUrl() {
  // Check for Vite env var (set at build time)
  const viteApiUrl = import.meta.env.VITE_API_URL;
  
  // Check for runtime config (useful for debugging)
  const runtimeApiUrl = window.__API_URL__;
  
  // Priority: runtime > vite env > localhost fallback
  const apiUrl = runtimeApiUrl || viteApiUrl || 'http://localhost:8080';
  
  // Log for debugging (only in development or if explicitly enabled)
  if (import.meta.env.DEV || window.__DEBUG_API__) {
    console.log('🔧 API Configuration:', {
      viteEnv: viteApiUrl,
      runtime: runtimeApiUrl,
      final: apiUrl,
      mode: import.meta.env.MODE,
    });
  }
  
  return apiUrl;
}

