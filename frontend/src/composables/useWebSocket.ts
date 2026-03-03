import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

export function useWebSocket(sessionId: string, participantId?: string) {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const authStore = useAuthStore()

  // For HR Dashboard: list of participants status
  const participantsStatus = ref<Record<string, any>>({})
  // For HR Dashboard: history of recent violations
  const recentViolations = ref<any[]>([])
  // For both HR and participants: session end flag
  const sessionEnded = ref(false)
  // For HR Dashboard: latest camera frame per participant_id (base64 data URL)
  const participantFrames = ref<Record<string, string>>({})

  const connect = () => {
    // Protocol ws:// or wss:// based on current origin
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = new URL(`${protocol}//${window.location.host}/api/ws/${sessionId}`)

    // Role for websocket layer: hr dashboard vs participant
    const role =
      authStore.isHR || authStore.isAdmin
        ? 'hr'
        : 'participant'

    url.searchParams.append('role', role)

    if (participantId) {
      url.searchParams.append('participant_id', participantId)
    }
    if (authStore.user?.name) {
      url.searchParams.append('user_name', String(authStore.user.name))
    }

    socket.value = new WebSocket(url.toString())

    socket.value.onopen = () => {
      isConnected.value = true
    }

    socket.value.onclose = () => {
      isConnected.value = false
    }

    socket.value.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)

        switch (msg.type) {
          case 'violation': {
            // Backend wraps extra info inside msg.data
            const payload = msg.data || {}
            const violation = {
              violation_type: payload.violation_type || msg.type,
              participant_id: msg.participant_id || payload.participant_id,
              participant_name: msg.user_name || payload.participant_name,
              timestamp: msg.timestamp,
            }
            recentViolations.value.unshift(violation)
            if (recentViolations.value.length > 50) recentViolations.value.pop()
            break
          }
          case 'participant_join': {
            participantsStatus.value[msg.participant_id] = {
              name: msg.user_name,
              status: 'online',
            }
            break
          }
          case 'participant_leave':
          case 'participant_finish': {
            if (participantsStatus.value[msg.participant_id]) {
              participantsStatus.value[msg.participant_id].status = 'offline'
            }
            break
          }
          case 'session_end': {
            sessionEnded.value = true
            break
          }
          case 'camera_frame': {
            // msg.data is the raw base64 JPEG data URL sent by the participant
            if (msg.participant_id && msg.data) {
              participantFrames.value[msg.participant_id] = msg.data
              // Receiving a frame proves the participant is live — always set 'online'
              // This also upgrades DB-seeded participants from 'active' → 'online'
              participantsStatus.value[msg.participant_id] = {
                ...participantsStatus.value[msg.participant_id],
                name: participantsStatus.value[msg.participant_id]?.name || msg.user_name || 'Participant',
                status: 'online',
              }
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
   * Called from the participant exam page every 5 seconds.
   */
  const sendCameraFrame = (dataUrl: string) => {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) return
    const msg = {
      type: 'camera_frame',
      data: dataUrl,
    }
    socket.value.send(JSON.stringify(msg))
  }

  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    socket,
    isConnected,
    participantsStatus,
    recentViolations,
    sessionEnded,
    participantFrames,
    sendCameraFrame,
    disconnect
  }
}
