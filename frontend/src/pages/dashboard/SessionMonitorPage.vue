<template>
  <div class="h-[calc(100vh-6rem)] flex flex-col pt-0 space-y-4">
    <!-- Custom Top Header for Monitor -->
    <div class="bg-[#1e3470] text-white -mx-6 -mt-6 px-6 py-4 flex items-center justify-between shadow-md z-10">
      <div class="flex items-center gap-6">
        <Button variant="ghost" size="icon" @click="router.back()" class="text-white hover:bg-white/10">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <div>
          <h2 class="text-xl font-bold tracking-tight flex items-center gap-3">
             {{ sessionTitle }}
          </h2>
          <p class="text-blue-200 text-xs mt-0.5 font-medium tracking-wide">Session ID: #{{ sessionId.substring(0,8).toUpperCase() }}</p>
        </div>
      </div>
      
      <div class="flex items-center gap-8">
        <div class="text-right">
           <p class="text-[10px] text-blue-200 font-bold uppercase tracking-widest mb-0.5">Time Remaining</p>
           <p class="text-2xl font-black leading-none tracking-widest">{{ formattedTimeRemaining }}</p>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-2 px-3 py-1.5 rounded-full border border-blue-400/30 bg-[#253f86] text-sm font-semibold">
            <div class="w-2 h-2 rounded-full" :class="isConnected ? 'bg-green-400 animate-pulse' : 'bg-red-500'"></div>
            Live Monitoring
          </div>
          <Button
            variant="destructive"
            class="bg-red-600 hover:bg-red-700 font-bold tracking-wide rounded-lg flex items-center gap-2"
            @click="endSession"
            :disabled="endingSession"
          >
            <StopCircleIcon class="w-4 h-4" /> End Session
          </Button>
        </div>
      </div>
    </div>

    <div class="flex items-center justify-between pt-4">
      <h3 class="text-lg font-bold text-slate-800 flex items-center gap-2">
        <UsersIcon class="w-5 h-5 text-slate-400" />
        {{ onlineCount }} Active Participants
      </h3>
      <div class="flex bg-slate-100 p-1 rounded-lg">
        <button class="px-4 py-1.5 text-sm font-semibold rounded-md bg-white shadow-sm text-slate-800">Grid View</button>
        <button class="px-4 py-1.5 text-sm font-semibold rounded-md text-slate-500 hover:text-slate-800">List View</button>
      </div>
    </div>

    <div class="flex-1 flex gap-6 min-h-0">
      <!-- Participants Grid -->
      <div class="flex-1 overflow-y-auto bg-slate-50/50 rounded-2xl border border-slate-100 p-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 pb-10">
          
          <Card v-for="(p, pid) in participantsStatus" :key="pid" class="overflow-hidden border border-slate-200 shadow-sm rounded-xl flex flex-col group transition-all hover:shadow-md hover:border-blue-200 cursor-pointer relative" :class="getCardBorderClass(p.status)">
            
            <!-- Video Feed: live frame when streaming, avatar fallback otherwise -->
            <div class="h-36 bg-slate-900 relative w-full overflow-hidden flex items-center justify-center">
               <!-- Live camera frame -->
               <img
                 v-if="participantFrames[pid]"
                 :src="participantFrames[pid]"
                 class="w-full h-full object-cover"
                 :alt="p.name"
               />
               <!-- Fallback avatar -->
               <img
                 v-else
                 :src="`https://api.dicebear.com/7.x/initials/svg?seed=${p.name}&backgroundColor=0f172a&textColor=ffffff`"
                 class="w-16 h-16 rounded-full opacity-50"
               />
               <div class="absolute inset-x-2 top-2 flex justify-between items-start">
                  <div class="p-1.5 bg-black/50 rounded-md text-white backdrop-blur-sm">
                    <MicOffIcon class="w-3.5 h-3.5" v-if="Math.random() > 0.5" />
                    <MicIcon class="w-3.5 h-3.5" v-else />
                  </div>
                  <!-- Live pulsing badge when streaming -->
                  <div v-if="participantFrames[pid]" class="flex items-center gap-1 bg-red-600/90 text-white text-[10px] font-bold px-2 py-0.5 rounded-md backdrop-blur-sm">
                    <span class="w-1.5 h-1.5 rounded-full bg-white animate-pulse"></span> LIVE
                  </div>
                  <Badge v-else :class="getWebcamBadgeStyle(p.status)" class="font-bold text-[10px] px-2 py-0.5 rounded-md border-0 uppercase flex items-center gap-1 backdrop-blur-md shadow-sm">
                     <component :is="getWebcamBadgeIcon(p.status)" class="w-3 h-3" />
                     {{ getWebcamStatusText(p.status) }}
                  </Badge>
               </div>
            </div>

            <CardContent class="p-3.5 bg-white flex flex-col gap-1 relative z-10">
              <div class="flex items-center justify-between">
                <h4 class="font-bold text-sm text-slate-800 truncate pr-2 w-full">{{ p.name }}</h4>
                <div class="w-2 h-2 rounded-full shrink-0" :class="p.status === 'online' ? 'bg-green-500' : 'bg-slate-300'"></div>
              </div>
              <div class="flex items-center justify-between">
                <p class="text-[11px] text-slate-500 font-medium">ID: {{ pid.toString().substring(0,5).toUpperCase() }}</p>
                <span class="text-[10px] font-bold text-red-600" v-if="p.status === 'multiple_faces'">Critical Alert</span>
              </div>
            </CardContent>
          </Card>

          <div v-if="Object.keys(participantsStatus).length === 0" class="col-span-full py-20 text-center text-slate-400 flex flex-col items-center justify-center">
            <VideoOffIcon class="w-12 h-12 mb-4 opacity-20" />
            <p class="font-semibold">Menunggu peserta bergabung...</p>
          </div>
        </div>
      </div>

      <!-- Violations Log Sidebar -->
      <Card class="w-80 flex flex-col min-h-0 shadow-sm border border-slate-200 rounded-2xl bg-white shrink-0">
        <CardHeader class="pb-3 border-b border-slate-100 pt-5 px-5">
          <div class="flex items-start justify-between">
            <div>
              <CardTitle class="text-base font-bold text-slate-800">Violation Log</CardTitle>
              <CardDescription class="text-xs text-slate-500 mt-0.5">Real-time alerts stream</CardDescription>
            </div>
            <a href="#" class="text-[11px] font-bold text-blue-600 hover:text-blue-800 transition-colors">Export CSV</a>
          </div>
        </CardHeader>
        
        <CardContent class="flex-1 overflow-y-auto px-3 py-3 space-y-2 bg-slate-50/50">
          <div v-if="recentViolations.length === 0" class="p-6 text-center text-xs text-slate-400 font-medium mt-10">
            Scanning in progress.<br/>No violations detected.
          </div>
          
          <div v-for="(v, i) in recentViolations" :key="i" class="p-3.5 bg-white border border-slate-100 rounded-xl shadow-[0_2px_10px_-4px_rgba(0,0,0,0.05)] transition-all hover:border-red-200 relative overflow-hidden group">
            <div class="absolute left-0 top-0 bottom-0 w-1" :class="getViolationSeverityColor(v.violation_type)"></div>
            
            <div class="flex gap-3">
              <div class="mt-0.5 text-slate-400">
                <component :is="getViolationIcon(v.violation_type)" class="w-4 h-4" :class="getViolationIconColor(v.violation_type)" />
              </div>
              <div class="flex-1">
                 <div class="flex justify-between items-start mb-0.5">
                   <h5 class="text-xs font-bold text-slate-800">{{ formatViolationType(v.violation_type) }}</h5>
                   <span class="text-[10px] text-slate-400 font-medium">{{ new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit', second:'2-digit'}) }}</span>
                 </div>
                 <p class="text-[11px] text-slate-500 leading-tight">Participant: <span class="font-medium text-slate-700">{{ v.participant_name || v.participant_id.substring(0,8) }}</span></p>
                 
                 <div class="mt-2 flex gap-2" v-if="v.violation_type === 'multiple_face' || v.violation_type === 'tab_switch'">
                    <Badge variant="destructive" class="text-[9px] px-1.5 py-0 rounded bg-red-100 text-red-700 hover:bg-red-200 border-0 uppercase font-bold tracking-wider">Flag</Badge>
                    <button class="text-[9px] font-bold text-slate-500 border border-slate-200 rounded px-1.5 py-0 hover:bg-slate-50 transition-colors">View Clip</button>
                 </div>
              </div>
            </div>
          </div>
        </CardContent>

        <div class="p-4 border-t border-slate-100 bg-white rounded-b-2xl">
           <div class="flex gap-4 items-center mb-3">
             <div class="flex items-center gap-1.5">
               <div class="w-2 h-2 rounded-full bg-red-500"></div>
               <span class="text-[11px] font-bold text-slate-600">3 Critical</span>
             </div>
             <div class="flex items-center gap-1.5">
               <div class="w-2 h-2 rounded-full bg-orange-500"></div>
               <span class="text-[11px] font-bold text-slate-600">5 Warnings</span>
             </div>
           </div>
           <Button variant="outline" class="w-full text-xs font-bold h-9 bg-slate-50 border-slate-200 text-slate-700">
             View All Logs
           </Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useWebSocket } from '@/composables/useWebSocket'
import client from '@/api/client'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { 
  ArrowLeftIcon, UsersIcon, StopCircleIcon, VideoOffIcon,
  MicIcon, MicOffIcon, ShieldAlertIcon, AlertTriangleIcon, FocusIcon,
  CheckCircle2Icon, EyeOffIcon, MaximizeIcon
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const sessionTitle = ref('Loading Session...')
const sessionParticipantsCount = ref(0)
const onlineParticipants = ref<any[]>([])

// Timer based on real session schedule
const timeRemaining = ref(0)
let timerInterval: number | null = null

// Init WebSocket connection
const { isConnected, participantsStatus, recentViolations, participantFrames } = useWebSocket(sessionId)

const onlineCount = computed(() => {
  return Object.values(participantsStatus.value).filter(p => p.status === 'online').length
})

const formattedTimeRemaining = computed(() => {
  const total = Math.max(timeRemaining.value, 0)
  const h = Math.floor(total / 3600).toString().padStart(2, '0')
  const m = Math.floor((total % 3600) / 60).toString().padStart(2, '0')
  const s = (total % 60).toString().padStart(2, '0')
  return `${h} : ${m} : ${s}`
})

onMounted(async () => {
  try {
    const s = await client.get(`/sessions/${sessionId}`)
    sessionTitle.value = s.data.name
    sessionParticipantsCount.value = s.data.max_participants

    // Compute remaining time from schedule + duration
    const schedule = new Date(s.data.schedule)
    const endTime = new Date(schedule.getTime() + s.data.duration_minutes * 60000)
    timeRemaining.value = Math.max(Math.floor((endTime.getTime() - Date.now()) / 1000), 0)

    if (timerInterval) clearInterval(timerInterval)
    timerInterval = window.setInterval(() => {
      if (timeRemaining.value > 0) timeRemaining.value--
    }, 1000)

    const p = await client.get(`/sessions/${sessionId}/participants`)
    onlineParticipants.value = p.data || []

    // Seed participantsStatus from DB (status stays as-is — only WS joins set 'online')
    onlineParticipants.value.forEach((part: any) => {
      if (!participantsStatus.value[part.id]) {
        participantsStatus.value[part.id] = {
          name: part.user_name || part.user_email || 'Unknown',
          status: part.status || 'active',
        }
      }
    })
  } catch (e) {
    console.error("Failed to load session info for monitor", e)
  }
})

const getCardBorderClass = (status: string) => {
  if (status === 'multiple_face') return 'border-red-500 shadow-[0_0_0_1px_rgba(239,68,68,1)]'
  if (status === 'no_face' || status === 'tab_switch') return 'border-orange-500 shadow-[0_0_0_1px_rgba(249,115,22,1)]'
  return ''
}

const endingSession = ref(false)
const endSession = async () => {
  if (endingSession.value) return
  if (!confirm('Akhiri sesi ujian ini? Peserta tidak akan dapat melanjutkan ujian.')) return

  endingSession.value = true
  try {
    await client.put(`/sessions/${sessionId}/lock`, { locked: true })
  } catch (e) {
    console.error('Failed to end session', e)
  } finally {
    endingSession.value = false
  }
}

const getWebcamBadgeStyle = (status: string) => {
  if (status === 'multiple_face') return 'bg-red-600/90 text-white'
  if (status === 'no_face') return 'bg-orange-600/90 text-white'
  if (status === 'tab_switch') return 'bg-yellow-500/90 text-slate-900'
  return 'bg-green-600/90 text-white'
}

const getWebcamBadgeIcon = (status: string) => {
  if (status === 'multiple_face') return AlertTriangleIcon
  if (status === 'no_face') return EyeOffIcon
  if (status === 'tab_switch') return FocusIcon
  return CheckCircle2Icon
}

const getWebcamStatusText = (status: string) => {
  if (status === 'multiple_face') return 'Multiple Faces'
  if (status === 'no_face') return 'No Face Detected'
  if (status === 'tab_switch') return 'Tab Switched'
  return 'Face Detected'
}

const formatViolationType = (type: string) => {
  const map: Record<string, string> = {
    'no_face': 'No Face Detected',
    'multiple_face': 'Multiple Faces',
    'tab_switch': 'Tab Switched',
    'fullscreen_exit': 'Fullscreen Exit'
  }
  return map[type] || type
}

const getViolationSeverityColor = (type: string) => {
  if (type === 'multiple_face') return 'bg-red-500'
  if (type === 'no_face') return 'bg-orange-500'
  return 'bg-yellow-400'
}

const getViolationIcon = (type: string) => {
  if (type === 'multiple_face') return AlertTriangleIcon
  if (type === 'no_face') return EyeOffIcon
  if (type === 'tab_switch') return FocusIcon
  if (type === 'fullscreen_exit') return MaximizeIcon
  return ShieldAlertIcon
}

const getViolationIconColor = (type: string) => {
  if (type === 'multiple_face') return 'text-red-500'
  if (type === 'no_face') return 'text-orange-500'
  return 'text-yellow-500'
}
</script>
