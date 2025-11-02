<template>
  <div v-if="showConfig" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-6 sm:p-8 max-w-md w-full border-2 border-orange-500">
      <h2 class="text-2xl font-bold text-orange-700 dark:text-orange-400 mb-4">
        ⚙️ Backend Configuration
      </h2>
      <p class="text-gray-700 dark:text-gray-300 mb-4">
        Enter your backend API URL (from Render):
      </p>
      <input
        v-model="backendUrl"
        type="text"
        placeholder="https://compatiblah-backend.onrender.com"
        class="w-full px-4 py-3 border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white mb-4"
        @keyup.enter="saveConfig"
      />
      <div class="flex gap-3">
        <button
          @click="saveConfig"
          class="flex-1 px-6 py-3 bg-orange-500 text-white rounded-xl font-semibold hover:bg-orange-600 transition-colors"
        >
          Save
        </button>
        <button
          @click="tryAutoDetect"
          class="px-4 py-3 bg-gray-500 text-white rounded-xl font-semibold hover:bg-gray-600 transition-colors"
        >
          Auto-detect
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const emit = defineEmits(['configured'])

const showConfig = ref(false)
const backendUrl = ref('')

onMounted(() => {
  // Check if backend URL is missing
  setTimeout(() => {
    const currentUrl = window.__API_URL__ || localStorage.getItem('backend_api_url')
    if (!currentUrl || currentUrl.includes('localhost') || currentUrl.includes('YOUR-BACKEND')) {
      showConfig.value = true
    }
  }, 500)
})

function saveConfig() {
  if (backendUrl.value.trim()) {
    window.__API_URL__ = backendUrl.value.trim()
    localStorage.setItem('backend_api_url', backendUrl.value.trim())
    emit('configured', backendUrl.value.trim())
    showConfig.value = false
    location.reload()
  }
}

async function tryAutoDetect() {
  // Try common Render URL patterns
  const hostname = window.location.hostname
  const projectName = hostname.replace(/\.vercel\.app.*/, '')
  
  const patterns = [
    `https://compatiblah.onrender.com`,
    `https://compatiblah-backend.onrender.com`,
    `https://${projectName}.onrender.com`,
    `https://${projectName}-backend.onrender.com`,
  ]
  
  // Try each pattern
  for (const url of patterns) {
    try {
      const response = await fetch(`${url}/health`, { method: 'GET', signal: AbortSignal.timeout(5000) })
      if (response.ok) {
        backendUrl.value = url
        saveConfig()
        return
      }
    } catch (e) {
      // Continue to next pattern
    }
  }
  
  alert('Could not auto-detect backend URL. Please enter your Render backend URL manually (e.g., https://compatiblah-backend.onrender.com)')
}
</script>

