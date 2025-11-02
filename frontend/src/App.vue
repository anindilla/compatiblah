<template>
  <div class="min-h-screen bg-gradient-to-br from-orange-50/80 via-amber-50/60 to-orange-100/80 dark:from-gray-950 dark:via-gray-900 dark:to-gray-950 backdrop-blur-3xl flex flex-col">
    <div class="container mx-auto px-4 sm:px-6 py-6 sm:py-10 max-w-7xl flex-1">
      <!-- Header -->
      <header class="text-center mb-6 sm:mb-8 md:mb-12">
        <h1 class="text-3xl sm:text-5xl md:text-6xl lg:text-7xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-orange-600 via-amber-600 to-orange-500 mb-2 sm:mb-3 md:mb-4 tracking-tighter">
          Compatiblah
        </h1>
        <p class="text-sm sm:text-base md:text-lg lg:text-xl text-gray-800 dark:text-gray-200 font-medium max-w-2xl mx-auto leading-relaxed px-2">
          Discover compatibility between two people as friends 👥, coworkers 💼, and partners 💕
        </p>
      </header>

      <!-- Assessment Form -->
      <div>
        <form @submit.prevent="submitAssessment" class="space-y-4 sm:space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6 md:gap-8">
            <PersonForm
              v-model="person1"
              title="Person 1"
            />
            <PersonForm
              v-model="person2"
              title="Person 2"
            />
          </div>

          <div class="flex flex-col sm:flex-row justify-center gap-3 sm:gap-4">
            <button
              type="submit"
              :disabled="loading || !isFormValid"
              class="w-full sm:w-auto px-8 sm:px-10 py-3 sm:py-4 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-xl font-semibold text-base sm:text-lg shadow-lg hover:from-orange-600 hover:to-orange-700 disabled:opacity-50 disabled:cursor-not-allowed focus:ring-4 focus:ring-orange-300 transition-all transform hover:scale-105 disabled:transform-none"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>Assessing...</span>
                <span>🔮</span>
              </span>
              <span v-else class="flex items-center justify-center gap-2">
                <span>✨</span>
                <span>Assess Compatibility</span>
              </span>
            </button>
            <button
              type="button"
              @click="resetForm"
              :disabled="loading"
              class="w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 bg-gray-500 hover:bg-gray-600 text-white rounded-xl font-semibold text-base sm:text-lg shadow-lg disabled:opacity-50 disabled:cursor-not-allowed focus:ring-4 focus:ring-gray-300 transition-all transform hover:scale-105 disabled:transform-none"
            >
              <span class="flex items-center justify-center gap-2">
                <span>🔄</span>
                <span>Reset</span>
              </span>
            </button>
          </div>
        </form>

        <!-- Error Message -->
        <div v-if="error" class="mt-4 sm:mt-6 bg-red-50 dark:bg-red-900/30 border-2 border-red-300 dark:border-red-700 text-red-800 dark:text-red-200 px-4 sm:px-6 py-3 sm:py-4 rounded-xl shadow-md" role="alert">
          <p class="font-medium text-sm sm:text-base">⚠️ {{ error }}</p>
        </div>

        <!-- Results -->
        <div v-if="results" class="mt-8">
          <CompatibilityResults
            :results="results"
          />
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import PersonForm from './components/PersonForm.vue'
import CompatibilityResults from './components/CompatibilityResults.vue'
import Footer from './components/Footer.vue'
import { getApiUrl } from './config.js'

const apiUrl = getApiUrl()
const person1 = ref({
  name: '',
  mbti: '',
  socialMedia: [],
  additionalParams: []
})
const person2 = ref({
  name: '',
  mbti: '',
  socialMedia: [],
  additionalParams: []
})
const loading = ref(false)
const error = ref(null)
const results = ref(null)
const assessmentId = ref(null)

const isFormValid = computed(() => {
  return person1.value.name && person1.value.mbti && person2.value.name && person2.value.mbti
})

function resetForm() {
  // Reset all form data
  person1.value = {
    name: '',
    mbti: '',
    socialMedia: [],
    additionalParams: []
  }
  person2.value = {
    name: '',
    mbti: '',
    socialMedia: [],
    additionalParams: []
  }
  // Clear results and errors
  results.value = null
  error.value = null
  assessmentId.value = null
  // Scroll to top
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function submitAssessment() {
  if (!isFormValid.value) {
    error.value = 'Please fill in all required fields (name and MBTI type for both people).'
    return
  }

  loading.value = true
  error.value = null
  results.value = null
  assessmentId.value = null

  try {
    const response = await fetch(`${apiUrl}/api/assess`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        person1: person1.value,
        person2: person2.value
      })
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    
    // Log for debugging
    console.log('Received data:', data)
    
    // Normalize explanations to structured format with subcategories and bullets
    const normalizeExplanation = (explanation) => {
      if (!explanation) {
        return { sections: [] }
      }
      
      // If already structured with subcategories (new format), return as is
      if (explanation.sections && Array.isArray(explanation.sections)) {
        // Check if sections have subcategories (new format)
        const hasSubcategories = explanation.sections.some(s => s.subcategories && Array.isArray(s.subcategories))
        if (hasSubcategories) {
          return explanation
        }
        
        // If sections exist but have old 'content' field, convert to subcategories
        return {
          sections: explanation.sections.map(section => ({
            heading: section.heading || 'Compatibility Analysis',
            subcategories: convertContentToSubcategories(section.content || '')
          }))
        }
      }
      
      // If it's a string (oldest format), convert to structured
      if (typeof explanation === 'string') {
        return {
          sections: [{
            heading: 'Compatibility Analysis',
            subcategories: convertContentToSubcategories(explanation)
          }]
        }
      }
      
      // Fallback
      return { sections: [] }
    }
    
    // Convert content string to subcategories with bullets
    const convertContentToSubcategories = (content) => {
      if (!content) {
        return [{
          title: 'Compatibility Analysis',
          bullets: [{ text: 'Analysis details are being prepared.' }]
        }]
      }
      
      // Split into sentences/paragraphs and create bullets
      const sentences = content.split(/[.!?]+/).filter(s => s.trim().length > 0)
      const bullets = sentences.map(s => ({ text: s.trim() + (s.trim().match(/[.!?]$/) ? '' : '.') }))
      
      // Group bullets into subcategories (3-5 bullets each)
      const subcategories = []
      const bulletsPerSubcat = 4
      
      for (let i = 0; i < bullets.length; i += bulletsPerSubcat) {
        const subcatBullets = bullets.slice(i, i + bulletsPerSubcat)
        if (subcatBullets.length > 0) {
          subcategories.push({
            title: i === 0 ? 'Key Insights' : i < bullets.length / 2 ? 'Additional Points' : 'Considerations',
            bullets: subcatBullets
          })
        }
      }
      
      return subcategories.length > 0 ? subcategories : [{
        title: 'Compatibility Analysis',
        bullets: [{ text: content }]
      }]
    }
    
    results.value = {
      friend_score: data.friend_score,
      coworker_score: data.coworker_score,
      partner_score: data.partner_score,
      overall_score: data.overall_score,
      friend_explanation: normalizeExplanation(data.friend_explanation),
      coworker_explanation: normalizeExplanation(data.coworker_explanation),
      partner_explanation: normalizeExplanation(data.partner_explanation)
    }
    assessmentId.value = data.id
    
    console.log('Normalized results:', results.value)

    // Scroll to results
    setTimeout(() => {
      window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' })
    }, 100)
  } catch (err) {
    error.value = err.message || 'Failed to assess compatibility. Please try again.'
    console.error('Assessment error:', err)
    // Log stack trace for debugging
    if (err.stack) {
      console.error('Error stack:', err.stack)
    }
  } finally {
    loading.value = false
  }
}

</script>
