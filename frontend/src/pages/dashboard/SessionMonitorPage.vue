<template>
  <div class="h-[calc(100vh-6rem)] flex flex-col pt-0 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <h2 class="text-2xl font-bold tracking-tight">Live Monitor: {{ sessionTitle }}</h2>
        <Badge :variant="isConnected ? 'default' : 'destructive'" class="ml-2">
            Status: {{ isConnected ? 'Terhubung (WebSocket)' : 'Terputus' }}
        </Badge>
      </div>
      <div class="text-sm font-medium border px-4 py-2 rounded-lg bg-background shadow-sm">
        <span class="text-green-600 font-bold">{{ onlineCount }}</span> / {{ sessionParticipantsCount }} Peserta Aktif
      </div>
    </div>

    <div class="flex-1 grid grid-cols-1 lg:grid-cols-4 gap-6 min-h-0">
      <!-- Participants Grid -->
      <Card class="lg:col-span-3 flex flex-col min-h-0 shadow-sm">
        <CardHeader class="pb-2 border-b bg-slate-50/50">
          <CardTitle class="text-lg">Kamera Peserta (Real-time)</CardTitle>
        </CardHeader>
        <CardContent class="flex-1 overflow-y-auto p-4 bg-slate-100">
          <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-4">
            <div v-for="(p, pid) in participantsStatus" :key="pid" 
                 class="bg-white rounded-lg shadow-sm border p-3 flex flex-col items-center gap-3 transition-shadow hover:shadow-md">
              <div class="w-16 h-16 rounded-full flex items-center justify-center font-bold text-xl relative"
                   :class="p.status === 'online' ? 'bg-primary/10 text-primary' : 'bg-slate-100 text-slate-400'">
                {{ p.name?.charAt(0).toUpperCase() || 'P' }}
                <div class="absolute bottom-0 right-0 w-4 h-4 rounded-full border-2 border-white"
                     :class="p.status === 'online' ? 'bg-green-500' : 'bg-slate-300'"></div>
              </div>
              <div class="text-center w-full">
                <p class="font-medium text-sm truncate w-full" :title="p.name">{{ p.name }}</p>
                <Badge variant="outline" class="mt-1 text-[10px] w-full justify-center">{{ p.status === 'online' ? 'Mengikuti Ujian' : 'Offline' }}</Badge>
              </div>
            </div>
            <div v-if="Object.keys(participantsStatus).length === 0" class="col-span-full py-12 text-center text-slate-400 font-medium">
              Menunggu peserta bergabung...
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Violations Log -->
      <Card class="flex flex-col min-h-0 shadow-sm border-orange-200">
        <CardHeader class="pb-2 border-b bg-orange-50/50">
          <CardTitle class="text-lg text-orange-800 flex items-center gap-2">
            <AlertTriangleIcon class="w-5 h-5"/>
            Log Pelanggaran
          </CardTitle>
        </CardHeader>
        <CardContent class="flex-1 overflow-y-auto p-0">
          <div v-if="recentViolations.length === 0" class="p-6 text-center text-sm text-slate-500">
            Belum ada pelanggaran terdeteksi.
          </div>
          <div class="divide-y">
            <div v-for="(v, i) in recentViolations" :key="i" class="p-3 bg-red-50/50 hover:bg-red-50 transition-colors">
              <div class="flex justify-between items-start mb-1">
                <span class="text-xs font-bold text-red-600">{{ formatViolationType(v.violation_type) }}</span>
                <span class="text-[10px] text-slate-400">{{ new Date().toLocaleTimeString() }}</span>
              </div>
              <p class="text-xs text-slate-700">Oleh: <span class="font-medium">{{ v.participant_name || v.participant_id }}</span></p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useWebSocket } from '@/composables/useWebSocket'
import client from '@/api/client'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ArrowLeftIcon, AlertTriangleIcon } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const sessionTitle = ref('Memuat...')
const sessionParticipantsCount = ref(0) // total participants registered
const onlineParticipants = ref<any[]>([])

// Init WebSocket connection
const { isConnected, participantsStatus, recentViolations } = useWebSocket(sessionId)

const onlineCount = computed(() => {
  return Object.values(participantsStatus.value).filter(p => p.status === 'online').length
})

onMounted(async () => {
  try {
    const s = await client.get(`/sessions/${sessionId}`)
    sessionTitle.value = s.data.title
    sessionParticipantsCount.value = s.data.participant_limit
    
    // Fetch initial list of participants already in session (Optional depending on API design)
    const p = await client.get(`/sessions/${sessionId}/participants`)
    onlineParticipants.value = p.data || []
    
    // Merge into local status
    onlineParticipants.value.forEach((part: any) => {
      if (!participantsStatus.value[part.id]) {
        participantsStatus.value[part.id] = { name: part.user?.name || 'Unknown', status: part.status }
      }
    })
  } catch (e) {
    console.error("Failed to load session info for monitor", e)
  }
})

const formatViolationType = (type: string) => {
  const map: Record<string, string> = {
    'no_face': 'Wajah tidak terdeteksi',
    'multiple_face': 'Lebih dari satu wajah',
    'tab_switch': 'Membuka Tab Lain',
    'fullscreen_exit': 'Keluar Layar Penuh'
  }
  return map[type] || type
}
</script>
