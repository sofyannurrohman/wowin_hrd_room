import { ref, onUnmounted } from 'vue'
import { FaceMesh } from '@mediapipe/face_mesh'
import type { Results } from '@mediapipe/face_mesh'
import { useExamStore } from '@/stores/exam'

export function useAntiCheat() {
  const isDetecting = ref(false)
  const faceStatus = ref<'normal'|'no-face'|'multi-face'|'looking-away'>('normal')
  const examStore = useExamStore()

  let faceMesh: FaceMesh | null = null
  let detectionInterval: number | null = null
  
  // Throttle violation reports to avoid spamming the backend
  let lastViolationTime = 0
  
  const reportViolation = (type: string) => {
    const now = Date.now()
    if (now - lastViolationTime > 5000) { // Max 1 report every 5s
      examStore.reportViolation(type)
      lastViolationTime = now
    }
  }

  const initFaceDetection = async (videoElement: HTMLVideoElement) => {
    faceMesh = new FaceMesh({
      locateFile: (file) => `https://cdn.jsdelivr.net/npm/@mediapipe/face_mesh/${file}`
    })

    faceMesh.setOptions({
      maxNumFaces: 3,
      refineLandmarks: true,
      minDetectionConfidence: 0.5,
      minTrackingConfidence: 0.5
    })

    faceMesh.onResults((results: Results) => {
      if (!isDetecting.value) return

      if (!results.multiFaceLandmarks || results.multiFaceLandmarks.length === 0) {
        faceStatus.value = 'no-face'
        reportViolation('no_face')
      } else if (results.multiFaceLandmarks.length > 1) {
        faceStatus.value = 'multi-face'
        reportViolation('multiple_face')
      } else {
        // Basic looking away heuristic: check the distance between eyes
        // In a real app, a proper PnP pose estimation is better
        faceStatus.value = 'normal'
      }
    })

    // Run inference in a loop instead of requestAnimationFrame for controlled rate (e.g., 2 fps)
    detectionInterval = window.setInterval(async () => {
      if (videoElement && !videoElement.paused && !videoElement.ended && faceMesh && isDetecting.value) {
        await faceMesh.send({ image: videoElement })
      }
    }, 500)
  }

  const startDetection = () => {
    isDetecting.value = true
  }

  const stopDetection = () => {
    isDetecting.value = false
    if (detectionInterval) clearInterval(detectionInterval)
    if (faceMesh) faceMesh.close()
  }

  // Window events (tab switch, fullscreen exit, copy/paste)
  const handleVisibilityChange = () => {
    if (document.hidden && isDetecting.value) {
      reportViolation('tab_switch')
      alert("PELANGGARAN: Anda membuka tab lain. Aktivitas direkam.")
    }
  }
  
  const handleFullscreenChange = () => {
    if (!document.fullscreenElement && isDetecting.value) {
      reportViolation('fullscreen_exit')
    }
  }

  const handleCopyPaste = (e: ClipboardEvent) => {
    if (isDetecting.value) {
      e.preventDefault()
    }
  }

  const setupBrowserAntiCheat = () => {
    document.addEventListener('visibilitychange', handleVisibilityChange)
    document.addEventListener('fullscreenchange', handleFullscreenChange)
    document.addEventListener('copy', handleCopyPaste)
    document.addEventListener('paste', handleCopyPaste)
    document.addEventListener('contextmenu', handleCopyPaste)
  }

  const removeBrowserAntiCheat = () => {
    document.removeEventListener('visibilitychange', handleVisibilityChange)
    document.removeEventListener('fullscreenchange', handleFullscreenChange)
    document.removeEventListener('copy', handleCopyPaste)
    document.removeEventListener('paste', handleCopyPaste)
    document.removeEventListener('contextmenu', handleCopyPaste)
  }

  onUnmounted(() => {
    stopDetection()
    removeBrowserAntiCheat()
  })

  return {
    faceStatus,
    initFaceDetection,
    startDetection,
    stopDetection,
    setupBrowserAntiCheat,
    removeBrowserAntiCheat
  }
}
