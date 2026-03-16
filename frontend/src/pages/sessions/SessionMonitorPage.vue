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
      <div class="flex items-center gap-4">
        <div class="relative w-64">
           <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
           <Input v-model="search" placeholder="Cari nama peserta..." class="pl-9" />
        </div>
        <h3 class="text-lg font-bold text-slate-800 flex items-center gap-2">
          <UsersIcon class="w-5 h-5 text-slate-400" />
          {{ onlineCount }} Active / {{ onlineParticipants.length }} Total Participants
        </h3>
      </div>
      <div class="flex bg-slate-100 p-1 rounded-lg">
        <button 
          @click="viewMode = 'grid'"
          class="px-4 py-1.5 text-sm font-semibold rounded-md transition-all"
          :class="viewMode === 'grid' ? 'bg-white shadow-sm text-slate-800' : 'text-slate-500 hover:text-slate-800'"
        >
          Grid View
        </button>
        <button 
          @click="viewMode = 'list'"
          class="px-4 py-1.5 text-sm font-semibold rounded-md transition-all"
          :class="viewMode === 'list' ? 'bg-white shadow-sm text-slate-800' : 'text-slate-500 hover:text-slate-800'"
        >
          List View
        </button>
      </div>
    </div>

    <div class="flex-1 flex gap-6 min-h-0">
      <!-- Participants Grid -->
      <div class="flex-1 overflow-y-auto bg-slate-50/50 rounded-2xl border border-slate-100 p-4">
        <!-- GRID VIEW -->
        <div v-if="viewMode === 'grid'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 pb-10">
          <Card v-for="(p, pid) in participantsStatus" :key="pid" class="overflow-hidden border border-slate-200 shadow-sm rounded-xl flex flex-col group transition-all hover:shadow-md hover:border-blue-200 cursor-pointer relative" :class="getCardBorderClass(p.status)">
            <!-- Video Feed: live frame when streaming, avatar fallback otherwise -->
            <div class="h-36 bg-slate-900 relative w-full overflow-hidden flex items-center justify-center">
               <img v-if="participantFrames[pid]" :src="participantFrames[pid]" class="w-full h-full object-cover" :alt="p.name" />
               <img v-else :src="`https://api.dicebear.com/7.x/initials/svg?seed=${p.name}&backgroundColor=0f172a&textColor=ffffff`" class="w-16 h-16 rounded-full opacity-50" />
               <div class="absolute inset-x-2 top-2 flex justify-between items-start">
                  <div class="p-1.5 bg-black/50 rounded-md text-white backdrop-blur-sm">
                    <MicOffIcon class="w-3.5 h-3.5" v-if="Math.random() > 0.5" />
                    <MicIcon class="w-3.5 h-3.5" v-else />
                  </div>
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
        </div>

        <!-- LIST VIEW -->
        <div v-if="viewMode === 'list'" class="bg-white rounded-xl border border-slate-200 overflow-hidden">
          <Table>
            <TableHeader class="bg-slate-50">
              <TableRow>
                <TableHead class="w-[300px]">Participant</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Safety Level</TableHead>
                <TableHead>Camera Feed</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="p in paginatedData" :key="p.id" class="hover:bg-slate-50/50">
                <TableCell>
                  <div>
                    <p class="font-bold text-slate-800 text-sm">{{ p.name }}</p>
                    <p class="text-[10px] text-slate-500 font-mono">ID: {{ p.id.toString().substring(0,8).toUpperCase() }}</p>
                  </div>
                </TableCell>
                <TableCell>
                  <div class="flex items-center gap-2">
                    <div class="w-2 h-2 rounded-full" :class="p.status === 'online' ? 'bg-green-500' : 'bg-slate-300'"></div>
                    <span class="text-xs font-medium capitalize">{{ p.status || 'Active' }}</span>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge 
                    variant="outline" 
                    :class="p.status === 'multiple_face' ? 'bg-red-50 text-red-700 border-red-200' : 'bg-green-50 text-green-700 border-green-200'"
                    class="text-[10px] font-bold"
                  >
                    {{ p.status === 'multiple_face' ? 'CRITICAL' : 'SECURE' }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div class="w-16 h-10 bg-slate-900 rounded border border-slate-200 overflow-hidden flex items-center justify-center">
                    <img v-if="participantFrames[p.id]" :src="participantFrames[p.id]" class="w-full h-full object-cover" />
                    <component v-else :is="getWebcamBadgeIcon(p.status)" class="w-4 h-4 text-white/30" />
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
          
          <div class="p-4 border-t border-slate-100 bg-slate-50/30">
            <DataTablePagination 
              v-if="participantsList.length > 0"
              :total="totalItems"
              v-model:pageSize="pageSize"
              v-model:currentPage="currentPage"
            />
          </div>
        </div>

        <div v-if="Object.keys(participantsStatus).length === 0" class="py-20 text-center text-slate-400 flex flex-col items-center justify-center">
          <VideoOffIcon class="w-12 h-12 mb-4 opacity-20" />
          <p class="font-semibold">Menunggu peserta bergabung...</p>
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
            <button @click="exportCSV" class="text-[11px] font-bold text-blue-600 hover:text-blue-800 transition-colors flex items-center gap-1">
              <FileSpreadsheetIcon class="w-3 h-3" /> Export CSV
            </button>
          </div>
        </CardHeader>
        
        <CardContent class="flex-1 overflow-y-auto px-3 py-3 space-y-2 bg-slate-50/50">
          <div v-if="allViolations.length === 0" class="p-6 text-center text-xs text-slate-400 font-medium mt-10">
            Scanning in progress.<br/>No violations detected.
          </div>
          
          <div v-for="(v, i) in allViolations.slice(0, 50)" :key="i" class="p-3.5 bg-white border border-slate-100 rounded-xl shadow-[0_2px_10px_-4px_rgba(0,0,0,0.05)] transition-all hover:border-red-200 relative overflow-hidden group">
            <div class="absolute left-0 top-0 bottom-0 w-1" :class="getViolationSeverityColor(v.violation_type)"></div>
            
            <div class="flex gap-3">
              <div class="mt-0.5 text-slate-400">
                <component :is="getViolationIcon(v.violation_type)" class="w-4 h-4" :class="getViolationIconColor(v.violation_type)" />
              </div>
              <div class="flex-1">
                 <div class="flex justify-between items-start mb-0.5">
                   <h5 class="text-xs font-bold text-slate-800">{{ formatViolationType(v.violation_type) }}</h5>
                   <span class="text-[10px] text-slate-400 font-medium">{{ new Date(v.detected_at || v.timestamp).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit', second:'2-digit'}) }}</span>
                 </div>
                 <p class="text-[11px] text-slate-500 leading-tight">Participant: <span class="font-medium text-slate-700">{{ v.participant_name || v.user_name || v.participant_id.substring(0,8) }}</span></p>
                 
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
               <span class="text-[11px] font-bold text-slate-600">{{ criticalCount }} Critical</span>
             </div>
             <div class="flex items-center gap-1.5">
               <div class="w-2 h-2 rounded-full bg-orange-500"></div>
               <span class="text-[11px] font-bold text-slate-600">{{ warningCount }} Warnings</span>
             </div>
           </div>
           <Button 
              variant="outline" 
              class="w-full text-xs font-bold h-9 bg-slate-50 border-slate-200 text-slate-700"
              @click="showAllLogs = true"
            >
              View All Logs
            </Button>
         </div>
      </Card>
    </div>

    <!-- View All Logs Dialog -->
    <Dialog :open="showAllLogs" @update:open="showAllLogs = $event">
      <DialogContent class="max-w-3xl max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>All Violation Logs - {{ sessionTitle }}</DialogTitle>
        </DialogHeader>
        
        <div class="flex-1 overflow-y-auto pr-2 mt-4">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Time</TableHead>
                <TableHead>Participant</TableHead>
                <TableHead>Violation Type</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="(v, i) in allViolations" :key="i">
                <TableCell class="font-mono text-xs">
                  {{ new Date(v.detected_at || v.timestamp).toLocaleTimeString() }}
                </TableCell>
                <TableCell class="font-medium">
                  {{ v.user_name || v.participant_name || v.participant_id.substring(0,8) }}
                </TableCell>
                <TableCell>
                  <div class="flex items-center gap-2">
                    <component :is="getViolationIcon(v.violation_type)" class="w-4 h-4" :class="getViolationIconColor(v.violation_type)" />
                    <span class="text-xs">{{ formatViolationType(v.violation_type) }}</span>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant="outline" :class="v.violation_type === 'multiple_face' ? 'text-red-600 border-red-200' : 'text-orange-600 border-orange-200'">
                    Flagged
                  </Badge>
                </TableCell>
              </TableRow>
              <TableRow v-if="allViolations.length === 0">
                <TableCell colspan="4" class="text-center py-10 text-slate-400">No violations recorded yet.</TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </DialogContent>
    </Dialog>
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
  CheckCircle2Icon, EyeOffIcon, MaximizeIcon, UserXIcon, FileSpreadsheetIcon, SearchIcon
} from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'

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

const allViolations = ref<any[]>([])
const showAllLogs = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')

const criticalCount = computed(() => allViolations.value.filter(v => v.violation_type === 'multiple_face').length)
const warningCount = computed(() => allViolations.value.filter(v => v.violation_type !== 'multiple_face').length)

const participantsList = computed(() => {
  return Object.entries(participantsStatus).map(([id, p]: [string, any]) => ({
    id,
    ...p
  }))
})

const search = ref('')
const filteredParticipants = computed(() => {
  return participantsList.value.filter(p => p.name.toLowerCase().includes(search.value.toLowerCase()))
})

const {
  pageSize,
  currentPage,
  paginatedData,
  totalItems
} = useDataTable(filteredParticipants)

// Sync WS violations to allViolations
import { watch } from 'vue'
watch(() => recentViolations.value, (newVal) => {
  if (newVal.length > 0) {
    const latest = newVal[0]
    // Only add if not already in (avoiding race condition with history fetch)
    const exists = allViolations.value.some(v => 
      (v.id && v.id === latest.id) || 
      (v.timestamp === latest.timestamp && v.participant_id === latest.participant_id && v.violation_type === latest.violation_type)
    )
    if (!exists) {
      allViolations.value.unshift(latest)
    }
  }
}, { deep: true })

const onlineCount = computed(() => {
  return Object.values(participantsStatus).filter((p: any) => p.status === 'online').length
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

    // Seed participantsStatus from DB
    onlineParticipants.value.forEach((part: any) => {
      if (!participantsStatus[part.id]) {
        participantsStatus[part.id] = {
          name: part.user_name || part.user_email || 'Unknown',
          status: part.status || 'active',
        }
      }
    })

    // Fetch violation history
    const v = await client.get(`/sessions/${sessionId}/violations`)
    const history = v.data || []
    // Merge history into allViolations (avoid duplicates)
    history.forEach((h: any) => {
      const exists = allViolations.value.some(v => v.id === h.id)
      if (!exists) allViolations.value.push(h)
    })
    // Sort by time descending
    allViolations.value.sort((a, b) => {
      const timeA = new Date(a.detected_at || a.timestamp).getTime()
      const timeB = new Date(b.detected_at || b.timestamp).getTime()
      return timeB - timeA
    })
  } catch (e) {
    console.error("Failed to load session info for monitor", e)
  }
})

const exportCSV = () => {
  if (allViolations.value.length === 0) {
    alert('No violation data to export')
    return
  }

  const headers = ['Time', 'Participant ID', 'Participant Name', 'Violation Type']
  const rows = allViolations.value.map(v => [
    new Date(v.detected_at || v.timestamp).toLocaleString(),
    v.participant_id,
    v.user_name || v.participant_name || 'Unknown',
    formatViolationType(v.violation_type)
  ])

  let csvContent = "data:text/csv;charset=utf-8," 
    + headers.join(",") + "\n"
    + rows.map(r => r.map(cell => `"${cell}"`).join(",")).join("\n")

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement("a")
  link.setAttribute("href", encodedUri)
  link.setAttribute("download", `violations_${sessionId}_${new Date().toISOString().slice(0,10)}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

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
  if (status === 'offline' || status === 'finished') return 'bg-slate-500/90 text-white'
  return 'bg-green-600/90 text-white'
}

const getWebcamBadgeIcon = (status: string) => {
  if (status === 'multiple_face') return AlertTriangleIcon
  if (status === 'no_face') return EyeOffIcon
  if (status === 'tab_switch') return FocusIcon
  if (status === 'offline' || status === 'finished') return UserXIcon
  return CheckCircle2Icon
}

const getWebcamStatusText = (status: string) => {
  if (status === 'multiple_face') return 'Multiple Faces'
  if (status === 'no_face') return 'No Face Detected'
  if (status === 'tab_switch') return 'Tab Switched'
  if (status === 'offline' || status === 'finished') return 'Finished'
  if (status === 'active') return 'Not Connected'
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
