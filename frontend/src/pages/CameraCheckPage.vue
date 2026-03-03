<template>
  <div class="min-h-screen bg-muted/30 py-10 px-4 flex flex-col items-center">
    <div class="max-w-2xl w-full space-y-8">
      <div class="text-center space-y-2">
        <h1 class="text-3xl font-bold tracking-tight">Persiapan Ujian</h1>
        <p class="text-muted-foreground">
          Ujian ini menggunakan sistem pengawasan kamera. Mohon izinkan akses kamera Anda.
        </p>
      </div>

      <Card class="overflow-hidden border-2 border-border">
        <div class="relative bg-primary aspect-video flex items-center justify-center">
          <video 
            ref="videoRef" 
            class="w-full h-full object-cover" 
            playsinline 
            muted
          ></video>
          
          <div v-if="cameraLoading" class="absolute inset-0 flex items-center justify-center bg-primary/80">
            <span class="text-primary-foreground animate-pulse">Menyiapkan kamera...</span>
          </div>

          <div v-if="cameraError" class="absolute inset-0 flex flex-col items-center justify-center bg-primary p-6 text-center">
            <AlertCircleIcon class="w-12 h-12 text-destructive mb-4" />
            <h3 class="text-primary-foreground font-semibold mb-2">Akses Kamera Ditolak</h3>
            <p class="text-primary-foreground/80 text-sm mx-auto max-w-sm">
              {{ cameraError }} Pastikan Anda telah memberikan izin kamera pada browser.
            </p>
            <Button variant="secondary" class="mt-4" @click="initCam">Coba Lagi</Button>
          </div>
        </div>
      </Card>

      <Card class="rounded-xl border-primary/20 bg-muted/30">
        <CardHeader class="pb-2">
          <CardTitle class="text-lg">Peraturan Ujian</CardTitle>
        </CardHeader>
        <CardContent class="space-y-3 text-sm text-muted-foreground">
          <p class="flex gap-2 items-center"><CheckCircleIcon class="w-5 h-5 text-primary shrink-0" /> Wajah harus terlihat jelas selama ujian berlangsung.</p>
          <p class="flex gap-2 items-center"><CheckCircleIcon class="w-5 h-5 text-primary shrink-0" /> Dilarang membuka tab browser lain.</p>
          <p class="flex gap-2 items-center"><CheckCircleIcon class="w-5 h-5 text-primary shrink-0" /> Ujian akan dilakukan dalam mode layar penuh (Fullscreen).</p>
          <p class="flex gap-2 items-center"><CheckCircleIcon class="w-5 h-5 text-primary shrink-0" /> Dilarang menggunakan fungsi Copy/Paste.</p>
        </CardContent>
      </Card>

      <Button size="lg" class="w-full h-14 text-lg" :disabled="!isCameraReady" @click="startExam">
        Sembunyikan Peringatan & Mulai Ujian
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useExamStore } from '@/stores/exam'
import { useCamera } from '@/composables/useCamera'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { CheckCircleIcon, AlertCircleIcon } from 'lucide-vue-next'

const router = useRouter()
const examStore = useExamStore()
const { startCamera, stopCamera, error: cameraError, isVideoLoading: cameraLoading } = useCamera()

const videoRef = ref<HTMLVideoElement | null>(null)
const isCameraReady = ref(false)

const initCam = async () => {
  if (videoRef.value) {
    const success = await startCamera(videoRef.value)
    isCameraReady.value = success as boolean
  }
}

onMounted(() => {
  if (!examStore.sessionId) {
    router.replace('/join')
    return
  }
  initCam()
})

onUnmounted(() => {
  stopCamera()
})

const startExam = async () => {
  try {
    // Attempt fullscreen
    if (document.documentElement.requestFullscreen) {
      await document.documentElement.requestFullscreen()
    }
  } catch (e) {
    console.warn("Fullscreen permission denied", e)
  }
  
  router.push(`/exam/${examStore.sessionId}`)
}
</script>
