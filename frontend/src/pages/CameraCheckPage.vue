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
          
          <div v-if="photoPreview" class="absolute inset-0 z-10 bg-black">
            <img :src="photoPreview" class="w-full h-full object-contain" alt="KTP Selfie Preview" />
            <div class="absolute bottom-4 left-0 right-0 flex justify-center gap-4">
              <Button variant="destructive" @click="retakePhoto">Foto Ulang</Button>
            </div>
          </div>
        </div>
        <div class="p-4 bg-muted text-center border-t border-border flex justify-between items-center px-6">
           <span class="text-sm font-medium text-muted-foreground w-full block" v-if="!photoPreview">Posisikan wajah dan KTP Anda menghadap kamera dengan jelas (Apabila tidak ada kesesuaian data pada proses validasi, maka anda akan di diskualifikasi).</span>
           <Button v-if="isCameraReady && !photoPreview" @click="takePhoto" class="w-full">
             <CameraIcon class="w-4 h-4 mr-2" />
             Ambil Foto
           </Button>
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

      <Button size="lg" class="w-full h-14 text-lg" :disabled="!isCameraReady || !photoPreview || isUploading" @click="startExam">
        {{ isUploading ? 'Mengunggah Foto...' : 'Sembunyikan Peringatan & Mulai Ujian' }}
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
import { CheckCircleIcon, AlertCircleIcon, CameraIcon } from 'lucide-vue-next'
import client from '@/api/client'

const router = useRouter()
const examStore = useExamStore()
const { startCamera, stopCamera, takeSnapshot, error: cameraError, isVideoLoading: cameraLoading } = useCamera()

const videoRef = ref<HTMLVideoElement | null>(null)
const isCameraReady = ref(false)
const photoPreview = ref<string | null>(null)
const isUploading = ref(false)

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

const takePhoto = () => {
  if (videoRef.value) {
    photoPreview.value = takeSnapshot(videoRef.value)
  }
}

const retakePhoto = () => {
  photoPreview.value = null
}

const dataURItoBlob = (dataURI: string) => {
  const parts = dataURI.split(',');
  const byteString = atob(parts[1] || '');
  const mimeString = parts[0]?.split(':')[1]?.split(';')[0] || 'image/jpeg';
  const ab = new ArrayBuffer(byteString.length);
  const ia = new Uint8Array(ab);
  for (let i = 0; i < byteString.length; i++) {
    ia[i] = byteString.charCodeAt(i);
  }
  return new Blob([ab], { type: mimeString });
}

const startExam = async () => {
  if (!photoPreview.value) {
    alert("Silakan ambil foto selfie KTP terlebih dahulu.");
    return;
  }

  isUploading.value = true;
  try {
    const participantId = examStore.participantId || localStorage.getItem('participantId');
    if (!participantId) throw new Error("Participant ID tidak ditemukan");

    const blob = dataURItoBlob(photoPreview.value);
    const formData = new FormData();
    formData.append('selfie', blob, 'selfie.jpg');

    await client.post(`/exam/participant/${participantId}/selfie`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });

    // Attempt fullscreen
    if (document.documentElement.requestFullscreen) {
      await document.documentElement.requestFullscreen()
    }
    
    router.push(`/exam/${examStore.sessionId}`)
  } catch (e) {
    console.warn("Fullscreen or upload error", e)
    alert("Gagal mengunggah foto. Pastikan koneksi stabil.");
  } finally {
    isUploading.value = false;
  }
}
</script>
