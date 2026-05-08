<template>
  <div class="typing-test-container select-none">
    <!-- Header Stats -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex gap-8">
        <div class="stat-item">
          <span class="text-xs font-bold uppercase tracking-widest text-slate-400 block mb-1">WPM</span>
          <span class="text-3xl font-black text-blue-600 font-mono leading-none">{{ wpm }}</span>
        </div>
        <div class="stat-item">
          <span class="text-xs font-bold uppercase tracking-widest text-slate-400 block mb-1">Accuracy</span>
          <span class="text-3xl font-black text-emerald-500 font-mono leading-none">{{ accuracy }}%</span>
        </div>
      </div>
      
      <div v-if="timeLeft !== null" class="flex items-center gap-3 px-6 py-3 bg-red-500/10 rounded-2xl border border-red-500/20">
        <span class="text-xs font-bold text-red-500 uppercase tracking-widest">Time Left</span>
        <span class="text-3xl font-black text-red-600 font-mono leading-none">{{ formattedTime }}</span>
      </div>
    </div>

    <!-- Side-by-Side Typing Area -->
    <div class="grid grid-cols-2 gap-6 h-[400px]">
      <!-- Left: Target Text -->
      <div 
        class="relative text-base leading-relaxed font-mono tracking-tight p-6 bg-slate-50 rounded-2xl border border-slate-200 overflow-y-auto shadow-inner"
      >
        <div class="whitespace-pre-wrap break-words text-slate-300">
          <span
            v-for="(char, index) in targetText"
            :key="index"
            :class="getCharClass(index)"
          >{{ char }}</span>
        </div>
      </div>

      <!-- Right: Typing Input -->
      <div class="relative p-6 bg-white rounded-2xl border-2 border-blue-100 shadow-sm focus-within:border-blue-400 transition-all flex flex-col">
        <textarea
          ref="inputRef"
          class="w-full h-full text-base leading-relaxed font-mono tracking-tight text-slate-700 bg-transparent resize-none focus:outline-none"
          v-model="userInput"
          @input="handleInput"
          @keydown="handleKeyDown"
          spellcheck="false"
          autofocus
        ></textarea>
      </div>
    </div>

    <div class="mt-4 flex items-center gap-2 text-slate-400 italic text-xs justify-center">
      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      Ketikkan teks sesuai dengan soal yang ada di sebelah kiri.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'

const props = defineProps<{
  targetText: string
  timerLimit?: number // in seconds
}>()

const emit = defineEmits<{
  (e: 'complete', results: { wpm: number; accuracy: number; raw: string }): void
  (e: 'update', results: { wpm: number; accuracy: number; raw: string }): void
  (e: 'timeUp'): void
}>()

const userInput = ref('')
const inputRef = ref<HTMLTextAreaElement | null>(null)
const isFocused = ref(false)
const startTime = ref<number | null>(null)
const elapsedTime = ref(0)
const totalErrors = ref(0)
const isFinished = ref(false)
const timeLeft = ref<number | null>(props.timerLimit || null)

let timerInterval: number | null = null

const wpm = computed(() => {
  if (elapsedTime.value === 0) return 0
  const words = userInput.value.length / 5
  const mins = elapsedTime.value / 60
  return Math.round(words / mins)
})

const accuracy = computed(() => {
  if (userInput.value.length === 0) return 100
  const correct = userInput.value.length - totalErrors.value
  return Math.max(0, Math.round((correct / userInput.value.length) * 100))
})

const formattedTime = computed(() => {
  if (timeLeft.value === null) return '00:00'
  const m = Math.floor(timeLeft.value / 60).toString().padStart(2, '0')
  const s = (timeLeft.value % 60).toString().padStart(2, '0')
  return `${m}:${s}`
})

const getCharClass = (index: number) => {
  if (index >= userInput.value.length) return 'text-slate-300'
  const isCorrect = userInput.value[index] === props.targetText[index]
  return isCorrect ? 'text-slate-900 bg-emerald-100/50' : 'text-red-500 bg-red-100'
}

const handleInput = () => {
  if (!startTime.value) {
    startTest()
  }

  // Calculate errors
  let errors = 0
  for (let i = 0; i < userInput.value.length; i++) {
    if (userInput.value[i] !== props.targetText[i]) {
      errors++
    }
  }
  totalErrors.value = errors

  emit('update', {
    wpm: wpm.value,
    accuracy: accuracy.value,
    raw: userInput.value
  })

  if (userInput.value.length >= props.targetText.length) {
    finishTest()
  }
}

const handleKeyDown = (e: KeyboardEvent) => {
  // Prevent default behavior for some keys if needed
  if (isFinished.value) e.preventDefault()
}

const startTest = () => {
  startTime.value = Date.now()
  timerInterval = window.setInterval(() => {
    elapsedTime.value = (Date.now() - startTime.value!) / 1000
    
    if (timeLeft.value !== null && timeLeft.value > 0) {
      timeLeft.value--
      if (timeLeft.value <= 0) {
        finishTest()
        emit('timeUp')
      }
    }
  }, 1000)
}

const finishTest = () => {
  isFinished.value = true
  if (timerInterval) clearInterval(timerInterval)
  emit('complete', {
    wpm: wpm.value,
    accuracy: accuracy.value,
    raw: userInput.value
  })
}

const focusInput = () => {
  inputRef.value?.focus()
  isFocused.value = true
}


onMounted(() => {
  focusInput()
  window.addEventListener('blur', () => isFocused.value = false)
  window.addEventListener('focus', () => isFocused.value = true)
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

// Sync props if they change
watch(() => props.targetText, () => {
  userInput.value = ''
  startTime.value = null
  elapsedTime.value = 0
  totalErrors.value = 0
  isFinished.value = false
  timeLeft.value = props.timerLimit || null
  if (timerInterval) clearInterval(timerInterval)
  nextTick(focusInput)
})
</script>

<style scoped>
.stat-item {
  @apply flex flex-col;
}
</style>
