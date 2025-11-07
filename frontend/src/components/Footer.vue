<template>
  <footer class="w-full mt-12 sm:mt-16 md:mt-20 py-6 sm:py-8 border-t-2 border-orange-200/50 dark:border-gray-700/50">
    <div class="container mx-auto px-4 sm:px-6 max-w-7xl">
      <!-- Privacy Disclaimer -->
      <div class="mb-6 sm:mb-8 p-4 sm:p-5 bg-orange-50/80 dark:bg-gray-800/80 rounded-2xl border-2 border-orange-200/60 dark:border-gray-700/60 backdrop-blur-sm">
        <h3 class="text-sm sm:text-base font-bold text-orange-800 dark:text-orange-300 mb-2 sm:mb-3 flex items-center gap-2">
          <span>🔒</span>
          <span>Privacy & Data Protection</span>
        </h3>
        <p class="text-xs sm:text-sm text-gray-700 dark:text-gray-300 leading-relaxed">
          Your privacy matters to us. The only personal details we ask for—names and MBTI types—are used <strong>solely</strong> to generate compatibility assessments via the Gemini API. We do <strong>NOT</strong> store any personal information in our database. Only the assessment results (scores and explanations) are saved, identified by a unique ID. Your personal data is never stored, shared, sold, or misused. Your data is yours, and we respect that.
        </p>
      </div>

      <!-- Footer Credits -->
      <div class="text-center">
        <p class="text-sm sm:text-base text-gray-600 dark:text-gray-400">
          Vibe-coded by 
          <a 
            href="https://anindilla.com" 
            target="_blank" 
            rel="noopener noreferrer"
            class="text-orange-600 dark:text-orange-400 hover:text-orange-700 dark:hover:text-orange-300 font-semibold underline decoration-2 underline-offset-2 hover:decoration-orange-500 transition-colors"
          >
            dilleuh
          </a>
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup>
// Footer component with privacy disclaimer
</script>
