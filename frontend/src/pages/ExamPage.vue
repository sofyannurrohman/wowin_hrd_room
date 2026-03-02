<template>
  <div class="min-h-screen bg-slate-50 flex flex-col">
    <!-- Header: Timer & Status -->
    <header class="bg-white border-b sticky top-0 z-10 px-6 py-4 flex items-center justify-between shadow-sm">
      <div class="flex items-center gap-4">
        <h1 class="font-bold text-xl text-primary">Sesi Ujian Aktif</h1>
        <Badge :variant="wsConnected ? 'default' : 'destructive'">
          {{ wsConnected ? 'Terhubung' : 'Terputus' }}
        </Badge>
      </div>
      <div class="flex items-center gap-6">
        <div class="text-right">
          <div class="text-sm font-medium text-slate-500 uppercase tracking-wider">Sisa Waktu</div>
          <div class="text-2xl font-mono font-bold" :class="{'text-destructive': isTimeLow}">
            {{ formattedTime }}
          </div>
        </div>
        <Button variant="default" @click="confirmSubmit" :disabled="examStore.isSubmitting">
          {{ examStore.isSubmitting ? 'Mengirim...' : 'Selesai & Kumpul' }}
        </Button>
      </div>
    </header>

    <main class="flex-1 max-w-7xl w-full mx-auto p-6 grid grid-cols-1 lg:grid-cols-4 gap-8 relative">
      <!-- Main Content: Questions -->
      <div class="lg:col-span-3 space-y-6">
        <template v-if="examStore.questions.length > 0">
          <Card v-for="(q, index) in examStore.questions" :key="q.id" class="shadow-sm">
            <CardHeader class="flex flex-row gap-4">
              <div class="flex-shrink-0 w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center font-bold">
                {{ index + 1 }}
              </div>
              <div class="space-y-4 w-full">
                <p class="text-lg leading-relaxed whitespace-pre-wrap">{{ q.content }}</p>
                <img v-if="q.image_url" :src="q.image_url" class="rounded max-h-64 object-contain border" />
              </div>
            </CardHeader>
            <CardContent class="ml-12 pt-0">
              <!-- Objective: MCQ or True/False -->
              <div v-if="q.type === 'multiple_choice' || q.type === 'true_false'" class="space-y-3">
                <label 
                  v-for="opt in q.options" 
                  :key="opt.id"
                  class="flex items-center p-4 border rounded-lg cursor-pointer transition-colors hover:bg-slate-50"
                  :class="{'bg-primary/5 border-primary': examStore.answers[q.id]?.selected_option_id === opt.id}"
                >
                  <input 
                    type="radio" 
                    :name="`q_${q.id}`" 
                    :value="opt.id"
                    :checked="examStore.answers[q.id]?.selected_option_id === opt.id"
                    @change="setAnswer(q.id, { selected_option_id: opt.id })"
                    class="w-4 h-4 text-primary bg-gray-100 border-gray-300 focus:ring-primary"
                  />
                  <span class="ml-3 text-slate-700">{{ opt.text }}</span>
                </label>
              </div>
              
              <!-- Subjective: Essay or Short Answer -->
              <div v-else class="space-y-3">
                <textarea 
                  class="w-full p-4 border rounded-lg min-h-[120px] focus:ring-2 focus:ring-primary focus:border-transparent outline-none"
                  placeholder="Tulis jawaban Anda di sini..."
                  :value="examStore.answers[q.id]?.text_answer"
                  @input="(e) => setAnswer(q.id, { text_answer: (e.target as HTMLTextAreaElement).value })"
                ></textarea>
              </div>
            </CardContent>
          </Card>
        </template>
        <div v-else class="text-center py-20 text-slate-500">
          Memuat soal ujian...
        </div>
      </div>

      <!-- Sidebar: Camera & Monitoring -->
      <div class="space-y-6">
        <div class="sticky top-24 space-y-4">
          <Card class="overflow-hidden border-2" :class="statusColorClass">
            <div class="bg-slate-900 aspect-video relative">
              <video ref="videoRef" class="w-full h-full object-cover" playsinline muted></video>
              <div class="absolute bottom-2 right-2 flex gap-2">
                <Badge variant="secondary" class="shadow-sm font-mono text-xs">
                  {{ faceStatus.toUpperCase() }}
                </Badge>
              </div>
            </div>
          </Card>
          
          <div class="text-sm text-slate-500 flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 rounded-full" :class="wsConnected ? 'bg-green-500' : 'bg-red-500'"></div>
              <span>Sistem Terhubung</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 rounded-full" :class="faceStatus === 'normal' ? 'bg-green-500' : 'bg-destructive'"></div>
              <span>Pelacakan Wajah Aktif</span>
            </div>
            <p class="mt-2 text-xs border-t pt-2">
              Pengawasan ketat aktif. Setiap pelanggaran akan dicatat otomatis.
            </p>
          </div>
        </div>
      </div>
    </main>

    <!-- Dialog for confirm submit -->
    <Dialog v-model:open="showSubmitDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Kumpulkan Ujian?</DialogTitle>
          <DialogDescription>
            Pastikan Anda telah mengisi semua jawaban. Anda tidak dapat kembali ke halaman ini setelah mengumpulkan jawaban.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="showSubmitDialog = false">Batal</Button>
          <Button variant="default" @click="doSubmit" :disabled="examStore.isSubmitting">
            {{ examStore.isSubmitting ? 'Progress...' : 'Ya, Kumpulkan' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useExamStore } from '@/stores/exam'
import { useCamera } from '@/composables/useCamera'
import { useAntiCheat } from '@/composables/useAntiCheat'
import { useWebSocket } from '@/composables/useWebSocket'

import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { toast } from 'vue-sonner'

const router = useRouter()
const route = useRoute()
const examStore = useExamStore()

const sessionId = route.params.sessionId as string
const showSubmitDialog = ref(false)

// Hardware & Anti-cheat
const videoRef = ref<HTMLVideoElement | null>(null)
const { startCamera, stopCamera } = useCamera()
const { faceStatus, initFaceDetection, startDetection, stopDetection, setupBrowserAntiCheat, removeBrowserAntiCheat } = useAntiCheat()
const { isConnected: wsConnected } = useWebSocket(sessionId, examStore.participantId || undefined)

// Timer
const timeRemaining = ref(7200) // Default 2 hours, should be fetched from session setup ideally
let timerInterval: number | null = null

const formattedTime = computed(() => {
  const h = Math.floor(timeRemaining.value / 3600).toString().padStart(2, '0')
  const m = Math.floor((timeRemaining.value % 3600) / 60).toString().padStart(2, '0')
  const s = (timeRemaining.value % 60).toString().padStart(2, '0')
  return `${h}:${m}:${s}`
})
const isTimeLow = computed(() => timeRemaining.value < 300) // Less than 5 mins

const statusColorClass = computed(() => {
  switch(faceStatus.value) {
    case 'no-face': return 'border-destructive'
    case 'multi-face': return 'border-orange-500'
    default: return 'border-border'
  }
})

onMounted(async () => {
  if (!examStore.participantId || !examStore.sessionId || examStore.sessionId !== sessionId) {
    router.replace('/join')
    return
  }

  // Init Data
  await examStore.fetchQuestions()

  // Init Hardware
  if (videoRef.value) {
    await startCamera(videoRef.value)
    await initFaceDetection(videoRef.value)
    startDetection()
  }

  // Security
  setupBrowserAntiCheat()

  // Timer run
  timerInterval = window.setInterval(() => {
    if (timeRemaining.value > 0) {
      timeRemaining.value--
    } else {
      // Auto submit
      clearInterval(timerInterval!)
      toast.error('Waktu habis! Jawaban otomatis dikirim.')
      doSubmit()
    }
  }, 1000)
})

onUnmounted(() => {
  stopCamera()
  stopDetection()
  removeBrowserAntiCheat()
  if (timerInterval) clearInterval(timerInterval)
})

// Answers interact
const setAnswer = (questionId: string, payload: any) => {
  examStore.setAnswer(questionId, payload)
}

const confirmSubmit = () => {
  showSubmitDialog.value = true
}

const doSubmit = async () => {
  showSubmitDialog.value = false
  const success = await examStore.submitExam()
  if (success) {
    router.push('/exam-finished')
  } else {
    toast.error('Gagal mengirim jawaban. Silakan coba lagi.')
  }
}

watch(faceStatus, (newStat) => {
  if (newStat === 'no-face') toast.error("Peringatan: Wajah tidak terdeteksi!")
  if (newStat === 'multi-face') toast.warning("Peringatan: Terdeteksi lebih dari satu wajah!")
})
</script>
