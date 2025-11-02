<template>
  <div v-if="results && (results.friend_score > 0 || results.coworker_score > 0 || results.partner_score > 0)" class="space-y-4 sm:space-y-6 md:space-y-8 animate-fadeIn">
    <!-- Overall Header -->
    <div class="bg-gradient-to-br from-orange-500 via-orange-600 to-amber-600 rounded-2xl sm:rounded-3xl shadow-2xl p-5 sm:p-7 md:p-9 lg:p-10 text-white border-4 border-orange-400/50 backdrop-blur-sm">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-4 mb-3 sm:mb-4 md:mb-6">
        <div class="flex items-center gap-2 sm:gap-3">
          <span class="text-xl sm:text-2xl md:text-3xl lg:text-4xl">🔍</span>
          <h2 class="text-xl sm:text-2xl md:text-3xl lg:text-4xl font-bold tracking-tight">Compatibility Analysis</h2>
        </div>
        <div class="flex items-center gap-2 bg-white/20 px-3 sm:px-4 py-1.5 sm:py-2 rounded-full backdrop-blur-sm self-start sm:self-auto">
          <span class="text-lg sm:text-xl md:text-2xl font-bold">{{ results.overall_score }}</span>
          <span class="text-xs sm:text-sm opacity-90">/ 5</span>
        </div>
      </div>
      <div class="flex flex-col sm:flex-row sm:flex-wrap items-start sm:items-center gap-2 sm:gap-3 md:gap-4">
        <span class="text-sm sm:text-base md:text-lg lg:text-xl font-semibold">Overall Compatibility Score:</span>
        <div class="flex gap-1 sm:gap-1.5">
          <span
            v-for="i in 5"
            :key="i"
            class="text-2xl sm:text-3xl md:text-4xl lg:text-5xl transition-all duration-300"
            :class="i <= results.overall_score ? 'text-amber-200 drop-shadow-2xl scale-110' : 'text-white/20'"
            aria-label="Star rating"
          >
            ★
          </span>
        </div>
      </div>
    </div>

    <!-- Detailed Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 sm:gap-5 md:gap-6">
      <!-- Friend Compatibility -->
      <div class="group bg-white/90 dark:bg-gray-800/90 rounded-2xl sm:rounded-3xl shadow-2xl p-4 sm:p-5 md:p-6 lg:p-7 border-2 border-orange-200/60 dark:border-gray-700/60 hover:border-orange-400 dark:hover:border-orange-600 hover:shadow-3xl transition-all duration-500 hover:-translate-y-1 backdrop-blur-sm">
        <div class="flex items-center justify-between mb-3 sm:mb-4 md:mb-6 pb-2 sm:pb-3 md:pb-4 border-b-2 border-orange-100 dark:border-gray-700">
          <div class="flex items-center gap-1.5 sm:gap-2">
            <span class="text-lg sm:text-xl md:text-2xl">🤝</span>
            <div>
              <h3 class="text-lg sm:text-xl md:text-2xl font-bold text-orange-700 dark:text-orange-400 mb-0.5 sm:mb-1">
                Friendship
              </h3>
              <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">As Friends</p>
            </div>
          </div>
          <div class="flex flex-col items-end gap-0.5 sm:gap-1">
            <div class="flex gap-0.5">
              <span
                v-for="i in 5"
                :key="i"
                class="text-base sm:text-lg transition-all"
                :class="i <= results.friend_score ? 'text-amber-500 scale-110' : 'text-gray-300'"
                aria-label="Friend compatibility rating"
              >
                ★
              </span>
            </div>
            <span class="text-xs sm:text-sm font-bold text-orange-600 dark:text-orange-400">{{ results.friend_score }}/5</span>
          </div>
        </div>
        <div class="space-y-4 sm:space-y-5">
          <template v-if="results.friend_explanation && results.friend_explanation.sections && results.friend_explanation.sections.length > 0">
            <div
              v-for="(section, idx) in results.friend_explanation.sections"
              :key="idx"
              class="space-y-3 sm:space-y-4"
            >
              <h4 class="text-base sm:text-lg md:text-xl font-bold text-orange-600 dark:text-orange-400 border-b border-orange-200 dark:border-orange-800 pb-1.5 sm:pb-2">
                {{ section.heading || 'Compatibility Analysis' }}
              </h4>
              <div v-if="section.subcategories && section.subcategories.length > 0" class="space-y-3 sm:space-y-4 ml-1 sm:ml-2">
                <div
                  v-for="(subcat, subIdx) in section.subcategories"
                  :key="subIdx"
                  class="space-y-1.5 sm:space-y-2"
                >
                  <h5 class="text-sm sm:text-base md:text-lg font-semibold text-gray-700 dark:text-gray-300">
                    {{ subcat.title }}
                  </h5>
                  <ul class="list-disc list-inside space-y-1 sm:space-y-1.5 ml-2 sm:ml-3">
                    <li
                      v-for="(bullet, bulletIdx) in subcat.bullets"
                      :key="bulletIdx"
                      class="text-sm sm:text-base md:text-lg text-gray-700 dark:text-gray-300 leading-relaxed"
                    >
                      {{ bullet.text }}
                    </li>
                  </ul>
                </div>
              </div>
              <p v-else-if="section.content" class="text-xs sm:text-sm md:text-base text-gray-800 dark:text-gray-200 leading-relaxed">
                {{ section.content }}
              </p>
            </div>
          </template>
          <div v-else class="flex items-center justify-center py-4">
            <div class="flex items-center gap-2 text-orange-600 dark:text-orange-400">
              <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-sm sm:text-base">Analyzing...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Coworker Compatibility -->
      <div class="group bg-white/90 dark:bg-gray-800/90 rounded-2xl sm:rounded-3xl shadow-2xl p-4 sm:p-5 md:p-6 lg:p-7 border-2 border-orange-200/60 dark:border-gray-700/60 hover:border-orange-400 dark:hover:border-orange-600 hover:shadow-3xl transition-all duration-500 hover:-translate-y-1 backdrop-blur-sm">
        <div class="flex items-center justify-between mb-3 sm:mb-4 md:mb-6 pb-2 sm:pb-3 md:pb-4 border-b-2 border-orange-100 dark:border-gray-700">
          <div class="flex items-center gap-1.5 sm:gap-2">
            <span class="text-lg sm:text-xl md:text-2xl">💼</span>
            <div>
              <h3 class="text-lg sm:text-xl md:text-2xl font-bold text-orange-700 dark:text-orange-400 mb-0.5 sm:mb-1">
                Workplace
              </h3>
              <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">As Coworkers</p>
            </div>
          </div>
          <div class="flex flex-col items-end gap-0.5 sm:gap-1">
            <div class="flex gap-0.5">
              <span
                v-for="i in 5"
                :key="i"
                class="text-base sm:text-lg transition-all"
                :class="i <= results.coworker_score ? 'text-amber-500 scale-110' : 'text-gray-300'"
                aria-label="Coworker compatibility rating"
              >
                ★
              </span>
            </div>
            <span class="text-xs sm:text-sm font-bold text-orange-600 dark:text-orange-400">{{ results.coworker_score }}/5</span>
          </div>
        </div>
        <div class="space-y-4 sm:space-y-5">
          <template v-if="results.coworker_explanation && results.coworker_explanation.sections && results.coworker_explanation.sections.length > 0">
            <div
              v-for="(section, idx) in results.coworker_explanation.sections"
              :key="idx"
              class="space-y-3 sm:space-y-4"
            >
              <h4 class="text-base sm:text-lg md:text-xl font-bold text-orange-600 dark:text-orange-400 border-b border-orange-200 dark:border-orange-800 pb-1.5 sm:pb-2">
                {{ section.heading || 'Compatibility Analysis' }}
              </h4>
              <div v-if="section.subcategories && section.subcategories.length > 0" class="space-y-3 sm:space-y-4 ml-1 sm:ml-2">
                <div
                  v-for="(subcat, subIdx) in section.subcategories"
                  :key="subIdx"
                  class="space-y-1.5 sm:space-y-2"
                >
                  <h5 class="text-sm sm:text-base md:text-lg font-semibold text-gray-700 dark:text-gray-300">
                    {{ subcat.title }}
                  </h5>
                  <ul class="list-disc list-inside space-y-1 sm:space-y-1.5 ml-2 sm:ml-3">
                    <li
                      v-for="(bullet, bulletIdx) in subcat.bullets"
                      :key="bulletIdx"
                      class="text-sm sm:text-base md:text-lg text-gray-700 dark:text-gray-300 leading-relaxed"
                    >
                      {{ bullet.text }}
                    </li>
                  </ul>
                </div>
              </div>
              <p v-else-if="section.content" class="text-xs sm:text-sm md:text-base text-gray-800 dark:text-gray-200 leading-relaxed">
                {{ section.content }}
              </p>
            </div>
          </template>
          <div v-else class="flex items-center justify-center py-4">
            <div class="flex items-center gap-2 text-orange-600 dark:text-orange-400">
              <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-sm sm:text-base">Analyzing...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Partner Compatibility -->
      <div class="group bg-white/90 dark:bg-gray-800/90 rounded-2xl sm:rounded-3xl shadow-2xl p-4 sm:p-5 md:p-6 lg:p-7 border-2 border-orange-200/60 dark:border-gray-700/60 hover:border-orange-400 dark:hover:border-orange-600 hover:shadow-3xl transition-all duration-500 hover:-translate-y-1 backdrop-blur-sm sm:col-span-2 md:col-span-1">
        <div class="flex items-center justify-between mb-3 sm:mb-4 md:mb-6 pb-2 sm:pb-3 md:pb-4 border-b-2 border-orange-100 dark:border-gray-700">
          <div class="flex items-center gap-1.5 sm:gap-2">
            <span class="text-lg sm:text-xl md:text-2xl">💕</span>
            <div>
              <h3 class="text-lg sm:text-xl md:text-2xl font-bold text-orange-700 dark:text-orange-400 mb-0.5 sm:mb-1">
                Romance
              </h3>
              <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">As Partners</p>
            </div>
          </div>
          <div class="flex flex-col items-end gap-0.5 sm:gap-1">
            <div class="flex gap-0.5">
              <span
                v-for="i in 5"
                :key="i"
                class="text-base sm:text-lg transition-all"
                :class="i <= results.partner_score ? 'text-amber-500 scale-110' : 'text-gray-300'"
                aria-label="Partner compatibility rating"
              >
                ★
              </span>
            </div>
            <span class="text-xs sm:text-sm font-bold text-orange-600 dark:text-orange-400">{{ results.partner_score }}/5</span>
          </div>
        </div>
        <div class="space-y-4 sm:space-y-5">
          <template v-if="results.partner_explanation && results.partner_explanation.sections && results.partner_explanation.sections.length > 0">
            <div
              v-for="(section, idx) in results.partner_explanation.sections"
              :key="idx"
              class="space-y-3 sm:space-y-4"
            >
              <h4 class="text-base sm:text-lg md:text-xl font-bold text-orange-600 dark:text-orange-400 border-b border-orange-200 dark:border-orange-800 pb-1.5 sm:pb-2">
                {{ section.heading || 'Compatibility Analysis' }}
              </h4>
              <div v-if="section.subcategories && section.subcategories.length > 0" class="space-y-3 sm:space-y-4 ml-1 sm:ml-2">
                <div
                  v-for="(subcat, subIdx) in section.subcategories"
                  :key="subIdx"
                  class="space-y-1.5 sm:space-y-2"
                >
                  <h5 class="text-sm sm:text-base md:text-lg font-semibold text-gray-700 dark:text-gray-300">
                    {{ subcat.title }}
                  </h5>
                  <ul class="list-disc list-inside space-y-1 sm:space-y-1.5 ml-2 sm:ml-3">
                    <li
                      v-for="(bullet, bulletIdx) in subcat.bullets"
                      :key="bulletIdx"
                      class="text-sm sm:text-base md:text-lg text-gray-700 dark:text-gray-300 leading-relaxed"
                    >
                      {{ bullet.text }}
                    </li>
                  </ul>
                </div>
              </div>
              <p v-else-if="section.content" class="text-xs sm:text-sm md:text-base text-gray-800 dark:text-gray-200 leading-relaxed">
                {{ section.content }}
              </p>
            </div>
          </template>
          <div v-else class="flex items-center justify-center py-4">
            <div class="flex items-center gap-2 text-orange-600 dark:text-orange-400">
              <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-sm sm:text-base">Analyzing...</span>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  results: {
    type: Object,
    default: null
  }
})
</script>



