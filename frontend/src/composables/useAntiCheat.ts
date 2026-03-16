import { ref, onUnmounted } from 'vue'
import { FaceMesh } from '@mediapipe/face_mesh'
import type { Results } from '@mediapipe/face_mesh'
import { useExamStore } from '@/stores/exam'
import { toast } from 'vue-sonner'

// ─── Module-level singleton ───────────────────────────────────────────────────
// MediaPipe's FaceMesh WASM cannot be instantiated more than once per browser
// session. Keeping a single instance here prevents the "Assertion failed: undefined"
// abort that occurs when navigating CameraCheckPage → ExamPage.
let sharedFaceMesh: FaceMesh | null = null
let faceMeshInitializing = false

const getOrCreateFaceMesh = async (): Promise<FaceMesh> => {
  if (sharedFaceMesh) return sharedFaceMesh

  // If another caller is already initializing, wait for it
  if (faceMeshInitializing) {
    await new Promise<void>((resolve) => {
      const wait = setInterval(() => {
        if (!faceMeshInitializing) { clearInterval(wait); resolve() }
      }, 50)
    })
    return sharedFaceMesh!
  }

  faceMeshInitializing = true
  sharedFaceMesh = new FaceMesh({
    locateFile: (file) => `https://cdn.jsdelivr.net/npm/@mediapipe/face_mesh/${file}`
  })

  sharedFaceMesh.setOptions({
    maxNumFaces: 3,
    refineLandmarks: true,
    minDetectionConfidence: 0.5,
    minTrackingConfidence: 0.5
  })

  // Initialize the model eagerly (avoids first-frame stall)
  await sharedFaceMesh.initialize()
  faceMeshInitializing = false
  return sharedFaceMesh
}
// ─────────────────────────────────────────────────────────────────────────────

export function useAntiCheat() {
  const isDetecting = ref(false)
  const faceStatus = ref<'normal' | 'no-face' | 'multi-face' | 'looking-away'>('normal')
  const examStore = useExamStore()

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
    try {
      const faceMesh = await getOrCreateFaceMesh()

      // Replace the onResults callback (safe to call multiple times)
      faceMesh.onResults((results: Results) => {
        if (!isDetecting.value) return

        if (!results.multiFaceLandmarks || results.multiFaceLandmarks.length === 0) {
          faceStatus.value = 'no-face'
          reportViolation('no_face')
        } else if (results.multiFaceLandmarks.length > 1) {
          faceStatus.value = 'multi-face'
          reportViolation('multiple_face')
        } else {
          faceStatus.value = 'normal'
        }
      })

      // Clear any existing detection interval before creating a new one
      if (detectionInterval) clearInterval(detectionInterval)

      // Run inference at ~2 fps
      detectionInterval = window.setInterval(async () => {
        if (videoElement && !videoElement.paused && !videoElement.ended && isDetecting.value) {
          try {
            await faceMesh.send({ image: videoElement })
          } catch (e) {
            // Silently ignore send errors (e.g. if video track ends)
          }
        }
      }, 500)
    } catch (e) {
      console.error('FaceMesh init error:', e)
    }
  }

  const startDetection = () => { isDetecting.value = true }

  const stopDetection = () => {
    isDetecting.value = false
    if (detectionInterval) {
      clearInterval(detectionInterval)
      detectionInterval = null
    }
    // ⚠️  Do NOT call faceMesh.close() — closing the singleton would destroy the
    // WASM instance permanently. We just stop sending frames.
  }

  // ─── Window events (tab switch, fullscreen exit, copy/paste) ─────────────
  const handleVisibilityChange = () => {
    if (document.hidden && isDetecting.value) {
      reportViolation('tab_switch')
      toast.error('PELANGGARAN: Anda membuka tab lain. Aktivitas direkam.')
    }
  }

  const handleFullscreenChange = () => {
    if (!document.fullscreenElement && isDetecting.value) {
      reportViolation('fullscreen_exit')
    }
  }

  const handleCopyPaste = (e: Event) => {
    if (isDetecting.value) e.preventDefault()
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
