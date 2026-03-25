<template>
  <div class="space-y-8 max-w-[1200px] mx-auto pb-10">
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div>
        <h2 class="text-3xl font-extrabold tracking-tight text-slate-800">Dashboard Overview</h2>
        <p class="text-slate-500 mt-1">Welcome back! mari kita lihat apa yang terjadi di Wowin Assesment Center hari ini.</p>
      </div>
      <div class="flex items-center gap-3 relative">
        <div class="relative">
          <Button 
            variant="outline" 
            size="icon" 
            class="rounded-full bg-white border-slate-200 text-slate-600 transition-all hover:bg-slate-50"
            @click="showNotifications = !showNotifications"
          >
            <BellIcon class="w-5 h-5" />
            <span v-if="notifications.length > 0" class="absolute top-0.5 right-0.5 w-2.5 h-2.5 bg-red-500 border-2 border-white rounded-full"></span>
          </Button>

          <!-- Notifications Dropdown -->
          <div v-if="showNotifications" class="absolute right-0 mt-3 w-80 bg-white rounded-2xl shadow-xl border border-slate-100 z-50 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200">
            <div class="p-4 border-b border-slate-100 bg-slate-50/50 flex items-center justify-between">
              <h3 class="font-bold text-slate-800 text-sm">Notifications</h3>
              <Badge variant="secondary" class="bg-blue-100 text-blue-700 text-[10px]">{{ notifications.length }} New</Badge>
            </div>
            <div class="max-h-[350px] overflow-y-auto">
              <div v-if="notifications.length === 0" class="p-8 text-center">
                <div class="w-12 h-12 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-3">
                  <BellIcon class="w-6 h-6 text-slate-300" />
                </div>
                <p class="text-xs text-slate-500 font-medium">No new notifications</p>
              </div>
              <div v-else>
                <div 
                  v-for="(n, idx) in notifications" 
                  :key="idx"
                  class="p-4 border-b border-slate-50 hover:bg-slate-50 cursor-pointer transition-colors group"
                  @click="router.push(n.link); showNotifications = false"
                >
                  <div class="flex gap-3">
                    <div class="w-8 h-8 rounded-full shrink-0 flex items-center justify-center" :class="n.bgClass">
                      <component :is="n.icon" class="w-4 h-4" :class="n.iconClass" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <p class="text-xs font-bold text-slate-800 leading-tight group-hover:text-blue-600 transition-colors">{{ n.title }}</p>
                      <p class="text-[11px] text-slate-500 mt-1 line-clamp-2">{{ n.message }}</p>
                      <p class="text-[10px] text-slate-400 mt-1.5 font-medium flex items-center gap-1">
                        <ClockIcon class="w-3 h-3" /> {{ n.timeLabel }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="notifications.length > 0" class="p-3 border-t border-slate-50 bg-slate-50/30">
              <button @click="showNotifications = false" class="w-full py-1 text-[11px] font-bold text-slate-500 hover:text-slate-800 transition-colors">
                Close
              </button>
            </div>
          </div>
        </div>

        <Button class="rounded-full bg-blue-600 hover:bg-blue-700 text-white font-semibold flex items-center gap-2" @click="router.push('/sessions/create')">
          <PlusIcon class="w-4 h-4" /> Buat Sesi Exam Baru
        </Button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-slate-500">Active Sessions</h3>
            <div class="w-10 h-10 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
              <MonitorIcon class="w-5 h-5" />
            </div>
          </div>
          <div class="text-4xl font-black text-slate-800 mb-4">{{ stats?.active_sessions || 0 }}</div>
          <Badge variant="secondary" class="bg-green-100 text-green-700 hover:bg-green-100 px-2 py-0.5 text-xs font-semibold rounded shrink-0 w-fit">
            <TrendingUpIcon class="w-3 h-3 mr-1 inline-block" /> +1 dari kemarin
          </Badge>
        </CardContent>
      </Card>

      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-slate-500">Total Peserta Hari Ini</h3>
            <div class="w-10 h-10 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center">
              <UsersIcon class="w-5 h-5" />
            </div>
          </div>
          <div class="text-4xl font-black text-slate-800 mb-4">{{ stats?.total_participants_today || 0 }}</div>
          <Badge variant="outline" class="text-slate-500 border-slate-200 px-2 py-0.5 text-xs font-semibold rounded shrink-0 w-fit">
            <UserCheckIcon class="w-3 h-3 mr-1 inline-block" /> {{ stats?.total_participants_today || 0 }} Telah masuk ke sistem
          </Badge>
        </CardContent>
      </Card>

      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-slate-500">Total Pelamar/Participant</h3>
            <div class="w-10 h-10 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center">
              <UsersIcon class="w-5 h-5" />
            </div>
          </div>
          <div class="text-4xl font-black text-slate-800 mb-4">{{ stats?.total_participants || 0 }}</div>
          <Badge variant="secondary" class="bg-emerald-50 text-emerald-600 hover:bg-emerald-50 px-2 py-0.5 text-xs font-semibold rounded shrink-0 w-fit border border-emerald-100">
            <TrendingUpIcon class="w-3 h-3 mr-1 inline-block" /> Data Keseluruhan
          </Badge>
        </CardContent>
      </Card>
    </div>

    <div class="grid gap-6 md:grid-cols-3">
      <!-- Recent Sessions Table -->
      <Card class="col-span-1 md:col-span-2 shadow-sm border border-slate-100 rounded-2xl bg-white">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-lg font-bold">Sesi Exam Terbaru</CardTitle>
          <a href="#" @click.prevent="router.push('/sessions')" class="text-sm font-semibold text-blue-600 hover:text-blue-800 transition-colors">View All</a>
        </CardHeader>
        <CardContent class="p-0">
          <Table>
            <TableHeader class="bg-slate-50/50">
              <TableRow class="border-b border-slate-100 hover:bg-transparent">
                <TableHead class="text-xs font-semibold text-slate-500 uppercase tracking-wider py-3 px-6 h-auto">Nama Sesi</TableHead>
                <TableHead class="text-xs font-semibold text-slate-500 uppercase tracking-wider py-3 px-2 h-auto">Tanggal</TableHead>
                <TableHead class="text-xs font-semibold text-slate-500 uppercase tracking-wider py-3 px-2 h-auto">Peserta</TableHead>
                <TableHead class="text-xs font-semibold text-slate-500 uppercase tracking-wider py-3 px-2 h-auto">Status</TableHead>
                <TableHead class="text-xs font-semibold text-slate-500 uppercase tracking-wider py-3 px-6 h-auto text-right">Aksi</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="session in recentSessions" :key="session.id" class="border-b border-slate-100 transition-colors hover:bg-slate-50/50">
                <TableCell class="py-4 px-6">
                  <div class="font-bold text-slate-800">{{ session.name }}</div>
                  <div class="text-xs text-slate-500 mt-0.5">Code: {{ getSessionCodeFromName(session.name) }}</div>
                </TableCell>
                <TableCell class="py-4 px-2 text-sm text-slate-600">
                  <div class="whitespace-nowrap">{{ formatDate(session.schedule) }}</div>
                </TableCell>
                <TableCell class="py-4 px-2">
                  <div class="text-sm text-slate-600 font-medium">
                    {{ session.submitted_participant_count || 0 }} / {{ session.total_participant_count || 0 }}
                  </div>
                  <div class="text-[10px] text-slate-400 mt-0.5 uppercase font-bold tracking-tighter">Selesai / Total</div>
                </TableCell>
                <TableCell class="py-4 px-2">
                  <Badge :class="getStatusBadgeClass(session)" class="font-medium text-[11px] px-2.5 py-0.5 rounded-full border-0">
                    {{ getStatusText(session) }}
                  </Badge>
                </TableCell>
                <TableCell class="py-4 px-6 text-right">
                  <Button variant="ghost" size="icon" class="text-slate-400 hover:text-blue-600 hover:bg-blue-50 transition-colors rounded-full" @click="router.push(`/sessions/${session.id}`)">
                    <EyeIcon class="w-4 h-4" />
                  </Button>
                </TableCell>
              </TableRow>
              
              <TableRow v-if="recentSessions.length === 0">
                 <TableCell colspan="5" class="text-center py-10 text-slate-500">
                   Belum ada sesi ujian yang dibuat.
                 </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <!-- Quick Actions -->
      <Card class="col-span-1 shadow-sm border border-slate-100 rounded-2xl bg-white h-fit">
        <CardHeader class="pb-4 border-b border-slate-100">
          <CardTitle class="text-lg font-bold text-slate-800">Quick Actions</CardTitle>
        </CardHeader>
        <CardContent class="p-4 space-y-3">
          
          <button @click="router.push('/sessions/create')" class="w-full flex items-center justify-between p-4 rounded-xl border border-slate-100 bg-white hover:bg-slate-50 hover:border-slate-200 transition-all group shadow-sm text-left">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 group-hover:bg-blue-100 transition-colors">
                <PlusIcon class="w-5 h-5" />
              </div>
              <div>
                <h4 class="font-bold text-sm text-slate-800">Buat Sesi Ujian Baru</h4>
                <p class="text-xs text-slate-500 mt-0.5">Siapkan sesi ujian baru</p>
              </div>
            </div>
            <ChevronRightIcon class="w-4 h-4 text-slate-400 group-hover:text-slate-600" />
          </button>

          <button @click="router.push('/questions')" class="w-full flex items-center justify-between p-4 rounded-xl border border-slate-100 bg-white hover:bg-slate-50 hover:border-slate-200 transition-all group shadow-sm text-left">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-full bg-purple-50 flex items-center justify-center text-purple-600 group-hover:bg-purple-100 transition-colors">
                <FileTextIcon class="w-5 h-5" />
              </div>
              <div>
                <h4 class="font-bold text-sm text-slate-800">Kelola Bank Soal</h4>
                <p class="text-xs text-slate-500 mt-0.5">Tambah atau edit soal</p>
              </div>
            </div>
            <ChevronRightIcon class="w-4 h-4 text-slate-400 group-hover:text-slate-600" />
          </button>

          <button @click="router.push('/sessions')" class="w-full flex items-center justify-between p-4 rounded-xl border border-slate-100 bg-white hover:bg-slate-50 hover:border-slate-200 transition-all group shadow-sm text-left">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-full bg-orange-50 flex items-center justify-center text-orange-600 group-hover:bg-orange-100 transition-colors">
                <KeyIcon class="w-5 h-5" />
              </div>
              <div>
                <h4 class="font-bold text-sm text-slate-800">Generate Token Ujian</h4>
                <p class="text-xs text-slate-500 mt-0.5">Buat token akses</p>
              </div>
            </div>
            <ChevronRightIcon class="w-4 h-4 text-slate-400 group-hover:text-slate-600" />
          </button>

        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import client from '@/api/client'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { 
  BellIcon, PlusIcon, MonitorIcon, UsersIcon,
  TrendingUpIcon, UserCheckIcon, EyeIcon,
  ChevronRightIcon, FileTextIcon, KeyIcon,
  CalendarIcon, ClockIcon, AlertTriangleIcon, MonitorPlayIcon
} from 'lucide-vue-next'

const router = useRouter()
const sessionStore = useSessionStore()

const stats = ref<any>(null)

const recentSessions = computed(() => {
  return [...sessionStore.sessions].slice(0, 4) // Show top 4
})

const showNotifications = ref(false)
const notifications = computed(() => {
  const list: any[] = []
  const now = new Date()

  sessionStore.sessions.forEach(s => {
    const start = new Date(s.schedule)
    const end = new Date(start.getTime() + s.duration_minutes * 60000)
    const status = getStatusText(s)

    // 1. Currently active sessions
    if (status === 'Active') {
      list.push({
        title: 'Sesi Sedang Berlangsung',
        message: `Sesi "${s.name}" sedang aktif sekarang. Klik untuk monitor.`,
        icon: MonitorPlayIcon,
        bgClass: 'bg-green-50',
        iconClass: 'text-green-600',
        timeLabel: 'LIVE',
        link: `/sessions/${s.id}/monitor`
      })
    }

    // 2. Scheduled sessions in the next 24 hours
    const hoursRemaining = (start.getTime() - now.getTime()) / (1000 * 60 * 60)
    if (status === 'Scheduled' && hoursRemaining > 0 && hoursRemaining < 24) {
      list.push({
        title: 'Sesi Akan Datang',
        message: `Sesi "${s.name}" akan dimulai dalam ${Math.floor(hoursRemaining)} jam.`,
        icon: CalendarIcon,
        bgClass: 'bg-blue-50',
        iconClass: 'text-blue-600',
        timeLabel: start.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        link: `/sessions/${s.id}`
      })
    }

    // 3. High violations in completed/active sessions (if data was available, but let's stick to session status)
  })

  return list
})

onMounted(async () => {
  await sessionStore.fetchSessions()
  
  try {
    const res = await client.get('/dashboard/stats')
    stats.value = res.data
  } catch (err) {
    console.error("Failed to load dashboard stats", err)
  }
})

const formatDate = (dateString: string) => {
  const d = new Date(dateString)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const getSessionCodeFromName = (name: string) => {
  // Simple mock to extract acronym/code from name if possible, or build one
  const words = name.split(' ')
  let code = words.map(w => w.charAt(0).toUpperCase()).join('')
  return `${code}-2024`
}

const getStatusText = (session: any) => {
  if (session.is_locked) return 'Locked'
  const now = new Date()
  const start = new Date(session.schedule)
  const end = new Date(start.getTime() + session.duration_minutes * 60000)
  if (now < start) return 'Scheduled'
  if (now > end) return 'Completed'
  return 'Active'
}

const getStatusBadgeClass = (session: any) => {
  const status = getStatusText(session)
  if (status === 'Active') return 'bg-green-100 text-green-700'
  if (status === 'Scheduled') return 'bg-blue-100 text-blue-700'
  if (status === 'Completed') return 'bg-slate-100 text-slate-700'
  return 'bg-red-100 text-red-700'
}
</script>
