import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

export function useWebSocket(sessionId: string, participantId?: string) {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const sessionEnded = ref(false)
  const recentViolations = ref<any[]>([])
  const authStore = useAuthStore()

  // Using reactive() (NOT ref) so that adding NEW keys to these maps
  // triggers Vue 3 template re-renders. ref<Record<>> does NOT track new keys.
  const participantsStatus = reactive<Record<string, any>>({})
  const participantFrames = reactive<Record<string, string>>({})

  const connect = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = new URL(`${protocol}//${window.location.host}/api/ws/${sessionId}`)

    // Role determination: if a participantId is supplied this is an exam-taker
    // connection (always 'participant'). Otherwise it's an HR dashboard connection.
    // Using participantId presence is more reliable than authStore.isHR because
    // both HR and participants may share the same browser profile during testing.
    const role = participantId ? 'participant' : 'hr'
    url.searchParams.append('role', role)

    if (participantId) {
      url.searchParams.append('participant_id', participantId)
    }
    if (authStore.user?.name) {
      url.searchParams.append('user_name', String(authStore.user.name))
    }

    console.log(`[WS] Connecting as role="${role}" participantId="${participantId}" sessionId="${sessionId}"`)
    console.log(`[WS] Full URL: ${url.toString()}`)

    socket.value = new WebSocket(url.toString())

    socket.value.onopen = () => {
      isConnected.value = true
      console.log(`[WS] Connected as "${role}"`)
    }
    socket.value.onclose = (e) => {
      isConnected.value = false
      console.warn(`[WS] Disconnected. Code=${e.code} Reason="${e.reason}"`)
    }
    socket.value.onerror = (e) => {
      console.error('[WS] Error:', e)
    }

    socket.value.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        console.log(`[WS] Received type="${msg.type}" participant_id="${msg.participant_id}" data_length=${msg.data ? String(msg.data).length : 0}`)

        switch (msg.type) {
          case 'violation': {
            const payload = msg.data || {}
            recentViolations.value.unshift({
              violation_type: payload.violation_type || msg.type,
              participant_id: msg.participant_id || payload.participant_id,
              participant_name: msg.user_name || payload.participant_name,
              timestamp: msg.timestamp,
            })
            if (recentViolations.value.length > 50) recentViolations.value.pop()
            break
          }

          case 'participant_join': {
            participantsStatus[msg.participant_id] = {
              name: msg.user_name,
              status: 'online',
            }
            break
          }

          case 'participant_leave':
          case 'participant_finish': {
            if (participantsStatus[msg.participant_id]) {
              participantsStatus[msg.participant_id].status = 'offline'
            }
            // Clear the last camera frame so the LIVE badge and image disappear
            delete participantFrames[msg.participant_id]
            break
          }

          case 'session_end': {
            sessionEnded.value = true
            break
          }

          case 'camera_frame': {
            if (msg.participant_id && msg.data) {
              console.log(`[WS] Storing camera frame for participant ${msg.participant_id}, data prefix: ${String(msg.data).substring(0, 30)}`)
              // Store the latest frame — triggers template re-render via reactive()
              participantFrames[msg.participant_id] = msg.data

              // Receiving a frame proves the participant is live
              const existing = participantsStatus[msg.participant_id]
              participantsStatus[msg.participant_id] = {
                name: existing?.name || msg.user_name || 'Participant',
                status: 'online',
              }
              console.log('[WS] participantFrames keys:', Object.keys(participantFrames))
              console.log('[WS] participantsStatus keys:', Object.keys(participantsStatus))
            } else {
              console.warn('[WS] camera_frame dropped: participant_id empty or no data', { participant_id: msg.participant_id, hasData: !!msg.data })
            }
            break
          }
        }
      } catch (e) {
        console.error('WS Parse Error', e)
      }
    }
  }

  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
      socket.value = null
    }
  }

  /**
   * Send a camera frame (base64 JPEG data URL) over the WebSocket.
   * Called from ExamPage every 5 seconds.
   */
  const sendCameraFrame = (dataUrl: string) => {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
      console.warn('[WS] sendCameraFrame: socket not open, state=', socket.value?.readyState)
      return
    }
    const payload = JSON.stringify({ type: 'camera_frame', data: dataUrl })
    console.log(`[WS] Sending camera_frame, payload size=${payload.length} bytes`)
    socket.value.send(payload)
  }

  /**
   * Notify HR that this participant has finished the exam.
   * Call this before submitting/unmounting so the HR feed stops immediately.
   */
  const sendParticipantFinish = () => {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) return
    socket.value.send(JSON.stringify({ type: 'participant_finish' }))
  }

  onMounted(() => { connect() })
  onUnmounted(() => { disconnect() })

  return {
    socket,
    isConnected,
    participantsStatus,
    recentViolations,
    sessionEnded,
    participantFrames,
    sendCameraFrame,
    sendParticipantFinish,
    disconnect,
  }
}
