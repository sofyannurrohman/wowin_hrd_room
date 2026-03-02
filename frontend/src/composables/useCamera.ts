import { ref, onMounted, onUnmounted } from 'vue'

export function useCamera() {
  const stream = ref<MediaStream | null>(null)
  const error = ref<string | null>(null)
  const isVideoLoading = ref(true)

  const stopCamera = () => {
    if (stream.value) {
      stream.value.getTracks().forEach(track => track.stop())
      stream.value = null
    }
  }

  const startCamera = async (videoElement: HTMLVideoElement) => {
    try {
      error.value = null
      isVideoLoading.value = true
      
      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'user', width: 640, height: 480 },
        audio: false
      })
      
      stream.value = mediaStream
      videoElement.srcObject = mediaStream
      
      // Wait until video data is loaded
      return new Promise((resolve) => {
        videoElement.onloadedmetadata = () => {
          videoElement.play()
          isVideoLoading.value = false
          resolve(true)
        }
      })
    } catch (err: any) {
      error.value = err.message || 'Kamera tidak dapat diakses.'
      isVideoLoading.value = false
      return false
    }
  }

  onUnmounted(() => {
    stopCamera()
  })

  return {
    stream,
    error,
    isVideoLoading,
    startCamera,
    stopCamera
  }
}
