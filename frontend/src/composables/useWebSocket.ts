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

  const connect = () => {
    // Protocol ws:// or wss:// based on current origin
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = new URL(`${protocol}//${window.location.host}/api/ws`)
    
    url.searchParams.append('session_id', sessionId)
    if (authStore.user) {
      url.searchParams.append('role', authStore.user.role.toString())
    } else {
      url.searchParams.append('role', 'Participant')
    }
    
    if (participantId) {
      url.searchParams.append('participant_id', participantId)
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
        const data = JSON.parse(event.data)
        
        if (data.type === 'participant_joined') {
          participantsStatus.value[data.participant_id] = { ...data, status: 'online' }
        } else if (data.type === 'participant_left') {
           if(participantsStatus.value[data.participant_id]) {
               participantsStatus.value[data.participant_id].status = 'offline'
           }
        } else if (data.type === 'violation') {
          recentViolations.value.unshift(data)
          if (recentViolations.value.length > 50) recentViolations.value.pop()
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
    disconnect
  }
}
