<template>
  <div class="bg-white dark:bg-gray-800 rounded-2xl sm:rounded-3xl shadow-2xl p-4 sm:p-6 md:p-8 space-y-4 sm:space-y-5 border-2 border-orange-200/50 dark:border-gray-700 backdrop-blur-sm">
    <h2 class="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold text-orange-700 dark:text-orange-400 mb-3 sm:mb-4 md:mb-6 border-b-2 border-orange-200 dark:border-orange-800 pb-1.5 sm:pb-2 md:pb-3 transition-all duration-300">
      <span class="mr-1.5 sm:mr-2">{{ personData.name && personData.name.trim() ? '👤' : '👥' }}</span>
      <span>{{ personData.name && personData.name.trim() ? personData.name : title }}</span>
      <span v-if="personData.name && personData.name.trim()" class="text-sm sm:text-base md:text-lg font-medium text-orange-500 dark:text-orange-400 opacity-70 ml-1.5 sm:ml-2">
        ({{ personData.mbti || 'MBTI' }})
      </span>
    </h2>

    <!-- Name -->
    <div>
      <label 
        for="name" 
        class="block text-sm font-semibold text-gray-800 dark:text-gray-200 mb-2"
      >
        ✏️ Name <span class="text-red-500 font-bold">*</span>
      </label>
      <input
        id="name"
        v-model="personData.name"
        type="text"
        required
        class="w-full px-3 sm:px-4 py-2 sm:py-3 text-sm sm:text-base border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white transition-all"
        placeholder="Enter name"
        aria-required="true"
      />
    </div>

    <!-- MBTI Type -->
    <div>
      <label 
        for="mbti" 
        class="block text-sm font-semibold text-gray-800 dark:text-gray-200 mb-2"
      >
        🧠 MBTI Type <span class="text-red-500 font-bold">*</span>
      </label>
      <select
        id="mbti"
        v-model="personData.mbti"
        required
        class="w-full px-3 sm:px-4 py-2 sm:py-3 text-sm sm:text-base border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white transition-all"
        aria-required="true"
      >
        <option value="">Select MBTI Type</option>
        <option value="INTJ">INTJ - Architect</option>
        <option value="INTP">INTP - Thinker</option>
        <option value="ENTJ">ENTJ - Commander</option>
        <option value="ENTP">ENTP - Debater</option>
        <option value="INFJ">INFJ - Advocate</option>
        <option value="INFP">INFP - Mediator</option>
        <option value="ENFJ">ENFJ - Protagonist</option>
        <option value="ENFP">ENFP - Campaigner</option>
        <option value="ISTJ">ISTJ - Logistician</option>
        <option value="ISFJ">ISFJ - Protector</option>
        <option value="ESTJ">ESTJ - Executive</option>
        <option value="ESFJ">ESFJ - Consul</option>
        <option value="ISTP">ISTP - Virtuoso</option>
        <option value="ISFP">ISFP - Adventurer</option>
        <option value="ESTP">ESTP - Entrepreneur</option>
        <option value="ESFP">ESFP - Entertainer</option>
      </select>
    </div>

    <!-- Social Media Fields -->
    <div class="space-y-2">
      <label class="block text-sm font-semibold text-gray-800 dark:text-gray-200 mb-2">
        📱 Social Media
      </label>
      <div
        v-for="(social, index) in personData.socialMedia"
        :key="index"
        class="flex gap-2"
      >
        <select
          v-model="social.platform"
          class="flex-1 px-3 sm:px-4 py-2 sm:py-3 text-sm sm:text-base border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white transition-all"
        >
          <option value="">Select Platform</option>
          <option value="Instagram">Instagram</option>
          <option value="Twitter/X">Twitter/X</option>
          <option value="LinkedIn">LinkedIn</option>
          <option value="Facebook">Facebook</option>
          <option value="TikTok">TikTok</option>
          <option value="YouTube">YouTube</option>
          <option value="Snapchat">Snapchat</option>
          <option value="Other">Other</option>
        </select>
        <input
          v-model="social.handle"
          type="text"
          class="flex-2 px-3 sm:px-4 py-2 sm:py-3 text-sm sm:text-base border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white transition-all"
          placeholder="Handle or URL"
        />
        <button
          v-if="personData.socialMedia.length > 0"
          type="button"
          @click="removeSocialMedia(index)"
          class="px-3 sm:px-4 py-1.5 sm:py-2 text-sm sm:text-base bg-red-500 text-white rounded-xl hover:bg-red-600 focus:ring-2 focus:ring-red-400 transition-all font-semibold"
          aria-label="Remove social media entry"
        >
          ×
        </button>
      </div>
      <button
        type="button"
        @click="addSocialMedia"
        class="text-xs sm:text-sm font-semibold text-orange-600 dark:text-orange-400 hover:text-orange-700 dark:hover:text-orange-300 hover:underline focus:ring-2 focus:ring-orange-400 rounded px-2 py-1 transition-colors"
      >
        + Add Social Media
      </button>
    </div>

    <!-- Dynamic Parameters -->
    <div class="space-y-2">
      <label class="block text-sm font-semibold text-gray-800 dark:text-gray-200 mb-2">
        ℹ️ Additional Information
      </label>
      <div
        v-for="(param, index) in personData.additionalParams"
        :key="index"
        class="flex gap-2"
      >
        <input
          v-model="param.key"
          type="text"
          class="flex-1 px-3 sm:px-4 py-2 sm:py-3 text-sm sm:text-base border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white transition-all"
          placeholder="Parameter name (e.g., Zodiac, DISC)"
        />
        <input
          v-model="param.value"
          type="text"
          class="flex-1 px-3 sm:px-4 py-2 sm:py-3 text-sm sm:text-base border-2 border-gray-300 dark:border-gray-600 rounded-xl focus:ring-2 focus:ring-orange-500 focus:border-orange-500 dark:bg-gray-700 dark:text-white transition-all"
          placeholder="Value"
        />
        <button
          v-if="personData.additionalParams.length > 0"
          type="button"
          @click="removeParameter(index)"
          class="px-3 sm:px-4 py-1.5 sm:py-2 text-sm sm:text-base bg-red-500 text-white rounded-xl hover:bg-red-600 focus:ring-2 focus:ring-red-400 transition-all font-semibold"
          aria-label="Remove parameter"
        >
          ×
        </button>
      </div>
      <button
        type="button"
        @click="addParameter"
        class="text-xs sm:text-sm font-semibold text-orange-600 dark:text-orange-400 hover:text-orange-700 dark:hover:text-orange-300 hover:underline focus:ring-2 focus:ring-orange-400 rounded px-2 py-1 transition-colors"
      >
        + Add New Parameter
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  modelValue: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const personData = ref({
  name: props.modelValue?.name || '',
  mbti: props.modelValue?.mbti || '',
  socialMedia: props.modelValue?.socialMedia || [],
  additionalParams: props.modelValue?.additionalParams || []
})

// Track if we're updating from parent to prevent infinite loop
let isUpdatingFromParent = false

watch(personData, (newVal) => {
  if (!isUpdatingFromParent) {
    emit('update:modelValue', newVal)
  }
}, { deep: true })

watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    isUpdatingFromParent = true
    // Only update if values actually changed to prevent loops
    const currentJson = JSON.stringify(personData.value)
    const newJson = JSON.stringify({
      name: newVal.name || '',
      mbti: newVal.mbti || '',
      socialMedia: newVal.socialMedia || [],
      additionalParams: newVal.additionalParams || []
    })
    
    if (currentJson !== newJson) {
      personData.value = {
        name: newVal.name || '',
        mbti: newVal.mbti || '',
        socialMedia: newVal.socialMedia || [],
        additionalParams: newVal.additionalParams || []
      }
    }
    // Reset flag after Vue's reactivity has processed
    setTimeout(() => {
      isUpdatingFromParent = false
    }, 0)
  }
}, { deep: true })

function addSocialMedia() {
  personData.value.socialMedia.push({ platform: '', handle: '' })
}

function removeSocialMedia(index) {
  personData.value.socialMedia.splice(index, 1)
}

function addParameter() {
  personData.value.additionalParams.push({ key: '', value: '' })
}

function removeParameter(index) {
  personData.value.additionalParams.splice(index, 1)
}
</script>




