<template>
  <div class="min-h-screen bg-gray-50 flex flex-col font-sans">
    <!-- Hidden canvas for frame capture (not visible in UI) -->
    <canvas ref="captureCanvas" width="320" height="240" class="hidden"></canvas>

    <!-- ─── Header ──────────────────────────────────────────────────── -->
    <header class="bg-white border-b sticky top-0 z-20 shadow-sm">
      <div class="max-w-[1400px] mx-auto px-6 py-3 flex items-center gap-6">

        <!-- Session info -->
        <div class="flex-shrink-0 min-w-0">
          <h1 class="font-bold text-gray-900 text-base leading-tight truncate">{{ examStore.sessionName || 'Sesi Ujian Aktif' }}</h1>
          <p class="text-xs text-gray-400 mt-0.5">Session ID: <span class="font-mono">#{{ shortSessionId }}</span></p>
        </div>

        <!-- Progress bar (center) -->
        <div class="flex-1 flex flex-col gap-1 min-w-0 max-w-xs mx-auto">
          <div class="flex justify-between text-xs text-gray-500 font-medium">
            <span>
              <span class="text-gray-400">{{ currentModuleTitle }}</span>
              &nbsp;·&nbsp; Soal {{ currentQuestionIndex + 1 }} dari {{ totalQuestions }}
            </span>
            <span>{{ completionPercent }}% Completed</span>
          </div>
          <div class="h-2 rounded-full bg-gray-200 overflow-hidden">
            <div
              class="h-full rounded-full bg-blue-600 transition-all duration-500"
              :style="{ width: completionPercent + '%' }"
            ></div>
          </div>
        </div>

        <!-- Timer + Finish button -->
        <div class="flex items-center gap-4 flex-shrink-0 ml-auto">
          <div
            class="flex items-center gap-2 font-mono text-xl font-bold tracking-widest px-4 py-1.5 rounded-lg"
            :class="isTimeLow ? 'text-red-600 bg-red-50' : 'text-gray-800 bg-gray-100'"
          >
            <svg class="w-5 h-5" :class="isTimeLow ? 'text-red-500' : 'text-blue-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
            </svg>
            {{ formattedTime }}
          </div>
          <button
            @click="handlePrimaryAction"
            :disabled="examStore.isSubmitting || showingTransition"
            class="flex items-center gap-2 px-5 py-2.5 bg-gray-900 hover:bg-gray-700 text-white text-sm font-semibold rounded-lg transition-colors disabled:opacity-50"
          >
            {{ isLastModule ? (examStore.isSubmitting ? 'Mengirim...' : 'Finish Test') : 'Lanjut Modul' }}
            <svg v-if="!examStore.isSubmitting" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
          </button>
        </div>

      </div>
    </header>

    <!-- ─── Module Transition Modal ──────────────────────────────── -->
    <div v-if="showingTransition" class="flex-1 flex items-center justify-center p-12">
      <div class="bg-white rounded-2xl shadow-lg border border-gray-200 p-12 flex flex-col items-center gap-5 max-w-lg w-full text-center">
        <!-- Step indicator -->
        <div class="flex items-center gap-2 text-xs font-semibold text-blue-600 bg-blue-50 px-4 py-1.5 rounded-full">
          Modul {{ examStore.currentModuleIndex }} selesai &nbsp;·&nbsp; {{ examStore.moduleGroups.length - examStore.currentModuleIndex }} modul tersisa
        </div>
        <div class="w-16 h-16 rounded-2xl bg-blue-100 flex items-center justify-center">
          <svg class="w-8 h-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
          </svg>
        </div>
        <div>
          <p class="text-sm text-gray-400 mb-1">Modul berikutnya</p>
          <h2 class="text-2xl font-bold text-gray-900">{{ nextModuleTitle }}</h2>
        </div>
        <p class="text-gray-500 text-sm max-w-sm">
          Anda akan mulai mengerjakan modul ujian berikutnya. Pastikan Anda sudah siap sebelum melanjutkan.
        </p>
        <button @click="startNextModule" class="mt-2 px-10 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition-colors text-base shadow-sm">
          Mulai Modul Sekarang
        </button>
      </div>
    </div>

    <!-- ─── Main 3-column layout ──────────────────────────────────── -->
    <main v-else class="flex-1 max-w-[1400px] w-full mx-auto px-4 py-5 grid grid-cols-[210px_1fr_230px] gap-5">

      <!-- Left: Question Navigation -->
      <aside class="sticky top-[73px] self-start space-y-3">
        <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-4">
          <div class="flex items-center gap-2 mb-3">
            <svg class="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/>
            </svg>
            <span class="font-semibold text-sm text-gray-800">Daftar Soal</span>
          </div>
          <!-- Legend -->
          <div class="flex gap-3 mb-3 text-[11px] text-gray-500">
            <div class="flex items-center gap-1">
              <span class="w-3 h-3 rounded-sm bg-blue-100 border border-blue-300 inline-block"></span>
              <span>telah dijawab</span>
            </div>
            <div class="flex items-center gap-1">
              <span class="w-3 h-3 rounded-sm bg-yellow-100 border border-yellow-400 inline-block"></span>
              <span>ditandai</span>
            </div>
          </div>
          <!-- Grid -->
          <div class="grid grid-cols-5 gap-1.5">
            <button
              v-for="(q, idx) in examStore.questions"
              :key="q.id"
              @click="goToQuestion(idx)"
              class="relative h-8 w-full rounded-md text-xs font-medium border transition-all"
              :class="navButtonClass(q, idx)"
            >
              {{ idx + 1 }}
              <!-- flag dot for marked -->
              <span
                v-if="markedForReview.has(q.id)"
                class="absolute -top-1 -right-1 w-2.5 h-2.5 rounded-full bg-yellow-400 border border-white"
              ></span>
            </button>
          </div>
        </div>
      </aside>

      <!-- Center: Current Question -->
      <section>
        <template v-if="currentQuestion">
          <div class="bg-white rounded-xl border border-gray-200 shadow-sm">
            <!-- Question header -->
            <div class="flex items-center justify-between px-6 pt-5 pb-2">
              <span class="px-3 py-1 text-xs font-semibold bg-gray-100 text-gray-600 rounded-full border border-gray-200">
                {{ questionTypeBadge(currentQuestion.type) }}
              </span>
              <span class="text-xs text-gray-400 font-mono">ID: {{ shortQuestionId(currentQuestion.id) }}</span>
            </div>

            <!-- Question text -->
            <div class="px-6 py-4">
              <p class="text-base font-semibold text-gray-900 leading-relaxed whitespace-pre-wrap">{{ currentQuestion.content }}</p>
              <img
                v-if="currentQuestion.image_url"
                :src="currentQuestion.image_url"
                class="mt-4 rounded-lg max-h-56 object-contain border border-gray-200"
                alt="Gambar soal"
              />
            </div>

            <!-- Options (MCQ / True-False) -->
            <div v-if="currentQuestion.type === 'Multiple Choice' || currentQuestion.type === 'True/False'" class="px-6 pb-4 space-y-3">
              <label
                v-for="opt in currentQuestion.options"
                :key="opt.id"
                class="flex items-start gap-4 p-4 border rounded-xl cursor-pointer transition-all"
                :class="
                  examStore.answers[currentQuestion.id]?.selected_option_id === opt.id
                    ? 'bg-blue-50 border-blue-500'
                    : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50/30'
                "
              >
                <!-- Custom radio -->
                <span class="flex-shrink-0 mt-0.5 w-5 h-5 rounded-full border-2 flex items-center justify-center transition-colors"
                  :class="examStore.answers[currentQuestion.id]?.selected_option_id === opt.id ? 'border-blue-600 bg-blue-600' : 'border-gray-300'"
                >
                  <span v-if="examStore.answers[currentQuestion.id]?.selected_option_id === opt.id" class="w-2 h-2 rounded-full bg-white"></span>
                </span>
                <input
                  type="radio"
                  :name="`q_${currentQuestion.id}`"
                  :value="opt.id"
                  :checked="examStore.answers[currentQuestion.id]?.selected_option_id === opt.id"
                  @change="setAnswer(currentQuestion.id, { selected_option_id: opt.id })"
                  class="sr-only"
                />
                <div>
                  <span class="text-sm text-gray-800">{{ opt.content }}</span>
                </div>
                <!-- checkmark when selected -->
                <svg
                  v-if="examStore.answers[currentQuestion.id]?.selected_option_id === opt.id"
                  class="ml-auto flex-shrink-0 w-5 h-5 text-blue-600"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"
                >
                  <path d="M20 6L9 17l-5-5"/>
                </svg>
              </label>
            </div>

            <!-- Essay / Short Answer -->
            <div v-else class="px-6 pb-4">
              <textarea
                class="w-full min-h-[140px] p-3 text-sm border border-gray-200 rounded-xl resize-y focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Tulis jawaban Anda di sini..."
                :value="examStore.answers[currentQuestion.id]?.text_answer ?? ''"
                @input="(e) => setAnswer(currentQuestion.id, { text_answer: (e.target as HTMLTextAreaElement).value })"
              ></textarea>
            </div>
          </div>

          <!-- Navigation bar -->
          <div class="mt-4 flex items-center justify-between">
            <button
              @click="prevQuestion"
              :disabled="currentQuestionIndex === 0"
              class="flex items-center gap-2 px-5 py-2.5 border border-gray-300 bg-white hover:bg-gray-50 text-gray-700 text-sm font-medium rounded-lg transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
              Previous
            </button>

            <button
              @click="toggleMarkForReview"
              class="flex items-center gap-2 px-5 py-2.5 text-sm font-medium rounded-lg transition-colors"
              :class="markedForReview.has(currentQuestion.id)
                ? 'bg-yellow-100 border border-yellow-400 text-yellow-800'
                : 'border border-gray-300 bg-white text-gray-600 hover:bg-gray-50'"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
              </svg>
              {{ markedForReview.has(currentQuestion.id) ? 'Ditandai' : 'Mark for Review' }}
            </button>

            <button
              @click="nextQuestion"
              :disabled="currentQuestionIndex >= totalQuestions - 1"
              class="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
            >
              Next Question
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
            </button>
          </div>
        </template>

        <!-- Loading -->
        <div v-else class="flex items-center justify-center h-64 bg-white rounded-xl border border-gray-200">
          <div class="flex flex-col items-center gap-3 text-gray-400">
            <svg class="w-10 h-10 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
            </svg>
            <span class="text-sm font-medium">Memuat soal ujian...</span>
          </div>
        </div>
      </section>

      <!-- Right: Proctoring -->
      <aside class="sticky top-[73px] self-start space-y-4">

        <!-- Camera feed -->
        <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
          <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
            <span class="text-xs font-bold tracking-widest text-gray-500 uppercase">Proctoring</span>
            <span class="flex items-center gap-1.5 text-xs font-semibold text-green-600">
              <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span> Live
            </span>
          </div>
          <div class="relative bg-gray-800 aspect-video" :class="['border-2', statusBorderClass]">
            <video ref="videoRef" class="w-full h-full object-cover" playsinline muted></video>
            <!-- Monitoring badge -->
            <div class="absolute bottom-2.5 left-1/2 -translate-x-1/2 flex items-center gap-1.5 bg-black/60 text-white text-[11px] font-semibold px-3 py-1 rounded-full">
              <span class="w-2 h-2 rounded-full bg-green-400"></span> Monitoring Active
            </div>
          </div>
          <p class="text-[11px] text-gray-400 px-4 py-2 text-center leading-snug">
            Your camera feed is being analyzed for exam integrity. Please keep your face within the frame.
          </p>
        </div>

        <!-- Exam Rules -->
        <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-4">
          <div class="flex items-center gap-2 mb-3">
            <svg class="w-4 h-4 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            <span class="font-semibold text-sm text-gray-800">Exam Rules</span>
          </div>
          <ul class="space-y-2.5">
            <li class="flex items-start gap-2.5">
              <svg class="w-4 h-4 text-red-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
              <span class="text-xs text-gray-600 leading-snug">Do not switch browser tabs or minimize the window.</span>
            </li>
            <li class="flex items-start gap-2.5">
              <svg class="w-4 h-4 text-red-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
              <span class="text-xs text-gray-600 leading-snug">No mobile devices allowed in the camera frame.</span>
            </li>
            <li class="flex items-start gap-2.5">
              <svg class="w-4 h-4 text-red-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
              <span class="text-xs text-gray-600 leading-snug">Only one face should be visible at all times.</span>
            </li>
          </ul>
        </div>

        <!-- Need Help -->
        <div class="bg-blue-50 rounded-xl border border-blue-200 p-4">
          <p class="font-semibold text-sm text-blue-900 mb-1">Need Help?</p>
          <p class="text-xs text-blue-700 leading-snug mb-3">If you encounter technical issues, contact the proctor immediately.</p>
          <button class="w-full py-2 px-4 bg-white border border-blue-300 hover:bg-blue-50 text-blue-700 text-xs font-semibold rounded-lg transition-colors">
            Chat with Support
          </button>
        </div>

      </aside>
    </main>

    <!-- ─── Submit Dialog ──────────────────────────────────────────── -->
    <div
      v-if="showSubmitDialog"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
      @click.self="showSubmitDialog = false"
    >
      <div class="bg-white rounded-2xl shadow-xl p-8 max-w-md w-full mx-4">
        <h2 class="text-xl font-bold text-gray-900 mb-2">Kumpulkan Ujian?</h2>
        <p class="text-gray-500 text-sm mb-6">Pastikan Anda telah mengisi semua jawaban. Anda tidak dapat kembali ke halaman ini setelah mengumpulkan jawaban.</p>
        <div class="flex gap-3 justify-end">
          <button @click="showSubmitDialog = false" class="px-5 py-2.5 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50">Batal</button>
          <button @click="doSubmit" :disabled="examStore.isSubmitting" class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg disabled:opacity-50">
            {{ examStore.isSubmitting ? 'Progress...' : 'Ya, Kumpulkan' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted, watch } from 'vue'
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router'
import { useExamStore } from '@/stores/exam'
import { useCamera } from '@/composables/useCamera'
import { useAntiCheat } from '@/composables/useAntiCheat'
import { useWebSocket } from '@/composables/useWebSocket'
import { toast } from 'vue-sonner'

const router = useRouter()
const route = useRoute()
const examStore = useExamStore()

const sessionId = route.params.sessionId as string
const showSubmitDialog = ref(false)
const showingTransition = ref(false)

// ─── Module helpers ───────────────────────────────────────────────
const isLastModule = computed(() => examStore.currentModuleIndex >= examStore.moduleGroups.length - 1)
const isLastQuestionInModule = computed(() => currentQuestionIndex.value >= totalQuestions.value - 1)

const currentModuleTitle = computed(() =>
  examStore.moduleGroups[examStore.currentModuleIndex]?.name || `Modul ${examStore.currentModuleIndex + 1}`
)
const nextModuleTitle = computed(() =>
  examStore.moduleGroups[examStore.currentModuleIndex + 1]?.name || `Modul ${examStore.currentModuleIndex + 2}`
)

// ─── Question navigation ──────────────────────────────────────────
const currentQuestionIndex = ref(0)
const markedForReview = ref<Set<string>>(new Set())

const totalQuestions = computed(() => examStore.questions.length)
const currentQuestion = computed(() => examStore.questions[currentQuestionIndex.value] ?? null)

// Global progress counts answered across ALL modules
const answeredCount = computed(() => {
  let count = 0
  for (const q of examStore.allQuestions) {
    const ans = examStore.answers[q.id]
    if (!ans) continue
    if (ans.selected_option_id || (ans.text_answer && String(ans.text_answer).trim())) count++
  }
  return count
})

const completionPercent = computed(() => {
  const total = examStore.allQuestions.length || examStore.questions.length
  if (!total) return 0
  return Math.round((answeredCount.value / total) * 100)
})

const shortSessionId = computed(() => {
  if (!sessionId) return ''
  return sessionId.slice(0, 8).toUpperCase()
})

const shortQuestionId = (id: string) => {
  if (!id) return ''
  return 'Q-' + id.slice(0, 4).toUpperCase()
}

const questionTypeBadge = (type: string) => {
  const map: Record<string, string> = {
    'Multiple Choice': 'Multiple Choice',
    'True/False': 'True / False',
    'Short Answer': 'Short Answer',
    'Psychological': 'Psychological',
    'multiple_choice': 'Multiple Choice',
    'true_false': 'True / False',
    'short_answer': 'Short Answer',
    'psychological': 'Psychological',
  }
  return map[type] ?? type
}

const isQuestionAnswered = (q: any) => {
  const ans = examStore.answers[q.id]
  if (!ans) return false
  if (ans.selected_option_id) return true
  if (ans.text_answer && String(ans.text_answer).trim().length > 0) return true
  return false
}



const navButtonClass = (q: any, idx: number) => {
  const isCurrent = currentQuestionIndex.value === idx
  const isAnswered = isQuestionAnswered(q)
  const isMarked = markedForReview.value.has(q.id)

  if (isCurrent) return 'bg-blue-600 text-white border-blue-600'
  if (isMarked) return 'bg-yellow-100 border-yellow-400 text-yellow-800'
  if (isAnswered) return 'bg-blue-100 border-blue-300 text-blue-700'
  return 'bg-white border-gray-200 text-gray-600 hover:border-blue-300'
}

const goToQuestion = (idx: number) => { currentQuestionIndex.value = idx }
const prevQuestion = () => { if (currentQuestionIndex.value > 0) currentQuestionIndex.value-- }
const nextQuestion = () => { if (currentQuestionIndex.value < totalQuestions.value - 1) currentQuestionIndex.value++ }

const toggleMarkForReview = () => {
  if (!currentQuestion.value) return
  const id = currentQuestion.value.id
  const next = new Set(markedForReview.value)
  if (next.has(id)) { next.delete(id) } else { next.add(id) }
  markedForReview.value = next
}

const setAnswer = (questionId: string, payload: any) => { examStore.setAnswer(questionId, payload) }

// ─── Camera & Anti-cheat ─────────────────────────────────────────
const videoRef = ref<HTMLVideoElement | null>(null)
const captureCanvas = ref<HTMLCanvasElement | null>(null)
const { startCamera, stopCamera } = useCamera()
const { faceStatus, initFaceDetection, startDetection, stopDetection, setupBrowserAntiCheat, removeBrowserAntiCheat } = useAntiCheat()
const { sessionEnded, sendCameraFrame, sendParticipantFinish } = useWebSocket(sessionId, examStore.participantId || undefined)

let frameCaptureInterval: number | null = null

const captureAndSendFrame = () => {
  const video = videoRef.value
  const canvas = captureCanvas.value
  if (!video || !canvas || video.readyState < 2) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
  sendCameraFrame(canvas.toDataURL('image/jpeg', 0.5))
}

const statusBorderClass = computed(() => {
  if (faceStatus.value === 'no-face' || faceStatus.value === 'multi-face') return 'border-red-500'
  return 'border-green-500'
})

// ─── Timer ────────────────────────────────────────────────────────
const timeRemaining = ref(examStore.timeRemaining || 0)
let timerInterval: number | null = null

const formattedTime = computed(() => {
  const h = Math.floor(timeRemaining.value / 3600).toString().padStart(2, '0')
  const m = Math.floor((timeRemaining.value % 3600) / 60).toString().padStart(2, '0')
  const s = (timeRemaining.value % 60).toString().padStart(2, '0')
  return `${h}:${m}:${s}`
})
const isTimeLow = computed(() => timeRemaining.value < 300)

// ─── Navigation Guards ───────────────────────────────────────────
// Warn participants when they try to close / reload the browser tab
const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (examStore.isSubmitting) return
  e.preventDefault()
  e.returnValue = ''
}
window.addEventListener('beforeunload', handleBeforeUnload)

// Warn participants when they try to navigate away within the SPA
onBeforeRouteLeave((_to, _from) => {
  if (examStore.isSubmitting) return true
  const confirmed = window.confirm(
    'Anda yakin ingin keluar dari halaman ujian? Progres jawaban Anda sudah tersimpan otomatis.'
  )
  return confirmed
})

// ─── Lifecycle ───────────────────────────────────────────────────
onMounted(async () => {
  if (!examStore.participantId || !examStore.sessionId || examStore.sessionId !== sessionId) {
    router.replace('/join')
    return
  }

  // Sync timer from store (which restores from localStorage on page reload)
  timeRemaining.value = examStore.timeRemaining

  await examStore.fetchModulesAndStart()

  if (videoRef.value) {
    await startCamera(videoRef.value)
    await initFaceDetection(videoRef.value)
    startDetection()
    frameCaptureInterval = window.setInterval(captureAndSendFrame, 5000)
  }

  setupBrowserAntiCheat()

  // Only start countdown if there is actually time remaining (avoids instant submit on reload)
  if (timeRemaining.value > 0) {
    timerInterval = window.setInterval(() => {
      if (timeRemaining.value > 0) {
        timeRemaining.value--
      } else {
        clearInterval(timerInterval!)
        toast.error('Waktu habis! Jawaban otomatis dikirim.')
        doSubmit()
      }
    }, 1000)
  } else if (timeRemaining.value === 0) {
    // Only auto-submit if an end timestamp was actually set and it's genuinely expired
    const storedEnd = Number(localStorage.getItem('examEndTimestamp') || 0)
    if (storedEnd && storedEnd < Date.now()) {
      toast.error('Waktu ujian telah habis.')
      await doSubmit()
    }
    // Otherwise: first join without timestamp yet — just wait for timer to be set normally
  }
})

onUnmounted(() => {
  stopCamera()
  stopDetection()
  removeBrowserAntiCheat()
  if (timerInterval) clearInterval(timerInterval)
  if (frameCaptureInterval) clearInterval(frameCaptureInterval)
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

// ─── Actions ──────────────────────────────────────────────────────
const handlePrimaryAction = () => {
  if (isLastModule.value) {
    // On the last module, Finish Test ends the exam
    showSubmitDialog.value = true
  } else if (isLastQuestionInModule.value) {
    // Last question of the current module → show transition to next module
    examStore.currentModuleIndex++
    examStore.loadCurrentModuleQuestions()
    currentQuestionIndex.value = 0
    markedForReview.value = new Set()
    showingTransition.value = true
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } else {
    // Still within module — just advance to next question
    nextQuestion()
  }
}

const startNextModule = () => {
  showingTransition.value = false
  currentQuestionIndex.value = 0
}

const doSubmit = async () => {
  showSubmitDialog.value = false
  // Notify HR immediately that this participant has finished
  sendParticipantFinish()
  // Stop frame capture before submitting
  if (frameCaptureInterval) { clearInterval(frameCaptureInterval); frameCaptureInterval = null }
  const success = await examStore.submitExam()
  if (success) {
    router.push('/exam-finished')
  } else {
    toast.error('Gagal mengirim jawaban. Silakan coba lagi.')
  }
}

watch(faceStatus, (newStat) => {
  if (newStat === 'no-face') toast.error('Peringatan: Wajah tidak terdeteksi!')
  if (newStat === 'multi-face') toast.warning('Peringatan: Terdeteksi lebih dari satu wajah!')
})

watch(sessionEnded, async (ended) => {
  if (!ended) return
  if (examStore.isSubmitting) return
  toast.error('Sesi ujian telah diakhiri oleh HR. Jawaban Anda akan dikirim.')
  await doSubmit()
})
</script>
