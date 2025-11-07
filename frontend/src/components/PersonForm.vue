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
  mbti: props.modelValue?.mbti || ''
})

// Track if we're updating from parent to prevent infinite loop
let isUpdatingFromParent = false

watch(personData, (newVal) => {
  if (!isUpdatingFromParent) {
    emit('update:modelValue', newVal)
  }
}, { deep: true })

watch(() => props.modelValue, (newVal) => {
  if (!newVal) {
    return
  }

  isUpdatingFromParent = true

  const incoming = {
    name: newVal.name || '',
    mbti: newVal.mbti || ''
  }

  if (personData.value.name !== incoming.name || personData.value.mbti !== incoming.mbti) {
    personData.value = incoming
  }

  setTimeout(() => {
    isUpdatingFromParent = false
  }, 0)
}, { deep: true })
</script>




