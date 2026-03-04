<template>
  <div class="space-y-6 max-w-7xl mx-auto pb-10">

    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()" class="rounded-full">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <div>
          <h2 class="text-2xl font-extrabold tracking-tight text-slate-800">Violation Review</h2>
          <p class="text-sm text-slate-500 mt-0.5">Sesi: <span class="font-semibold text-slate-700">{{ sessionName }}</span></p>
        </div>
      </div>
      <div class="flex items-center gap-3">
        <Badge variant="outline" class="text-xs font-semibold border-slate-200 text-slate-600 px-3 py-1 rounded-full">
          {{ violations.length }} Total Pelanggaran
        </Badge>
        <Button variant="outline" size="sm" @click="exportCSV" class="gap-2 text-xs font-bold">
          <DownloadIcon class="w-3.5 h-3.5" /> Export CSV
        </Button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div
        v-for="stat in summaryStats"
        :key="stat.label"
        class="bg-white rounded-2xl border border-slate-100 shadow-sm p-5 flex flex-col gap-2 hover:shadow-md transition-all"
      >
        <div class="flex items-center justify-between">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">{{ stat.label }}</p>
          <div class="w-8 h-8 rounded-lg flex items-center justify-center" :class="stat.iconBg">
            <component :is="stat.icon" class="w-4 h-4" :class="stat.iconColor" />
          </div>
        </div>
        <p class="text-3xl font-black text-slate-800">{{ stat.value }}</p>
        <Badge :class="stat.badgeClass" class="w-fit text-[10px] font-bold px-2 py-0.5 rounded-full border-0">
          {{ stat.sublabel }}
        </Badge>
      </div>
    </div>

    <!-- Filter Row -->
    <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3">
      <div class="relative flex-1">
        <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
        <input
          v-model="search"
          type="text"
          placeholder="Cari nama peserta..."
          class="w-full pl-9 pr-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>
      <select
        v-model="filterType"
        class="px-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-slate-700 font-medium"
      >
        <option value="">Semua Jenis</option>
        <option value="multiple_face">Multiple Faces</option>
        <option value="no_face">No Face Detected</option>
        <option value="tab_switch">Tab Switch</option>
        <option value="fullscreen_exit">Fullscreen Exit</option>
        <option value="looking_away">Looking Away</option>
        <option value="camera_off">Camera Off</option>
      </select>
      <select
        v-model="sortBy"
        class="px-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-slate-700 font-medium"
      >
        <option value="time_desc">Waktu Terbaru</option>
        <option value="time_asc">Waktu Terlama</option>
        <option value="participant">Nama Peserta</option>
        <option value="severity">Tingkat Bahaya</option>
      </select>
    </div>

    <!-- Main Content: Two Columns -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

      <!-- Per-Participant Summary (left) -->
      <div class="lg:col-span-1 space-y-3">
        <h3 class="text-sm font-bold text-slate-600 uppercase tracking-wider px-1">Per Peserta</h3>
        <div class="space-y-2">
          <div
            v-for="ps in participantSummary"
            :key="ps.id"
            @click="selectedParticipantId = (selectedParticipantId === ps.id ? null : ps.id)"
            class="bg-white border rounded-xl p-4 cursor-pointer transition-all hover:border-blue-300 hover:shadow-sm"
            :class="selectedParticipantId === ps.id ? 'border-blue-400 shadow-md ring-1 ring-blue-300 bg-blue-50/40' : 'border-slate-100'"
          >
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2">
                <img
                  :src="`https://api.dicebear.com/7.x/initials/svg?seed=${ps.name}&backgroundColor=1e40af&textColor=ffffff`"
                  class="w-8 h-8 rounded-full"
                />
                <div>
                  <p class="text-sm font-bold text-slate-800 leading-tight">{{ ps.name }}</p>
                  <p class="text-[10px] text-slate-400 font-mono">{{ ps.id.substring(0,8).toUpperCase() }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-lg font-black" :class="ps.total >= 5 ? 'text-red-600' : ps.total >= 2 ? 'text-orange-500' : 'text-slate-700'">{{ ps.total }}</p>
                <p class="text-[10px] text-slate-400">pelanggaran</p>
              </div>
            </div>
            <div class="flex flex-wrap gap-1">
              <span
                v-for="(count, type) in ps.byType"
                :key="type"
                class="text-[10px] font-bold px-1.5 py-0.5 rounded-full"
                :class="getTypeColor(String(type))"
              >
                {{ formatType(String(type)) }}: {{ count }}
              </span>
            </div>
          </div>
          <div v-if="participantSummary.length === 0" class="bg-white border border-slate-100 rounded-xl p-8 text-center text-slate-400">
            <ShieldCheckIcon class="w-10 h-10 mx-auto mb-2 opacity-30" />
            <p class="text-sm font-medium">Tidak ada pelanggaran.</p>
          </div>
        </div>
      </div>

      <!-- Violation Timeline (right) -->
      <div class="lg:col-span-2">
        <h3 class="text-sm font-bold text-slate-600 uppercase tracking-wider px-1 mb-3">
          Timeline Pelanggaran
          <span v-if="selectedParticipantId" class="text-blue-600 normal-case font-normal text-xs ml-2">
            — Difilter: {{ participantSummary.find(p => p.id === selectedParticipantId)?.name }}
            <button @click="selectedParticipantId = null" class="ml-1 text-slate-400 hover:text-red-500">✕</button>
          </span>
        </h3>

        <div class="space-y-2">
          <div
            v-for="v in filteredViolations"
            :key="v.id"
            class="bg-white border border-slate-100 rounded-xl p-4 flex items-start gap-4 hover:border-slate-200 hover:shadow-sm transition-all relative overflow-hidden"
          >
            <!-- Severity bar -->
            <div class="absolute left-0 top-0 bottom-0 w-1 rounded-l-xl" :class="getSeverityBar(v.violation_type)"></div>

            <!-- Icon -->
            <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0 ml-1" :class="getTypeBg(v.violation_type)">
              <component :is="getViolationIcon(v.violation_type)" class="w-4.5 h-4.5" :class="getTypeIconColor(v.violation_type)" />
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-2 mb-0.5">
                <h4 class="text-sm font-bold text-slate-800">{{ formatViolationType(v.violation_type) }}</h4>
                <Badge class="text-[10px] border-0 font-bold shrink-0" :class="getSeverityBadge(v.violation_type)">
                  {{ getSeverityLabel(v.violation_type) }}
                </Badge>
              </div>
              <p class="text-xs text-slate-500">
                Peserta: <span class="font-semibold text-slate-700">{{ v.user_name || v.participant_id?.substring(0, 8) }}</span>
              </p>
              <p class="text-[11px] text-slate-400 mt-1 font-mono">{{ formatTime(v.detected_at) }}</p>
            </div>
          </div>

          <div v-if="filteredViolations.length === 0" class="bg-white border border-dashed border-slate-200 rounded-xl p-12 text-center text-slate-400">
            <ShieldCheckIcon class="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p class="font-semibold">Tidak ada pelanggaran ditemukan.</p>
            <p class="text-xs mt-1">Coba ubah filter atau periksa kembali nanti.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  ArrowLeftIcon, DownloadIcon, SearchIcon, ShieldAlertIcon,
  AlertTriangleIcon, EyeOffIcon, MonitorOffIcon, MaximizeIcon,
  ShieldCheckIcon, CameraOffIcon, ScanEyeIcon
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const sessionName = ref('Loading...')
const violations = ref<any[]>([])
const search = ref('')
const filterType = ref('')
const sortBy = ref('time_desc')
const selectedParticipantId = ref<string | null>(null)

onMounted(async () => {
  try {
    const [sRes, vRes] = await Promise.all([
      client.get(`/sessions/${sessionId}`),
      client.get(`/sessions/${sessionId}/violations`)
    ])
    sessionName.value = sRes.data.name
    violations.value = vRes.data || []
  } catch (e) {
    console.error('Failed to load violations', e)
  }
})

const severityOrder: Record<string, number> = {
  multiple_face: 3,
  no_face: 2,
  tab_switch: 2,
  camera_off: 1,
  fullscreen_exit: 1,
  looking_away: 1,
}

const sortedViolations = computed(() => {
  let list = [...violations.value]
  if (sortBy.value === 'time_desc') list.sort((a, b) => new Date(b.detected_at).getTime() - new Date(a.detected_at).getTime())
  if (sortBy.value === 'time_asc') list.sort((a, b) => new Date(a.detected_at).getTime() - new Date(b.detected_at).getTime())
  if (sortBy.value === 'participant') list.sort((a, b) => (a.user_name || '').localeCompare(b.user_name || ''))
  if (sortBy.value === 'severity') list.sort((a, b) => (severityOrder[b.violation_type] || 0) - (severityOrder[a.violation_type] || 0))
  return list
})

const filteredViolations = computed(() => {
  let list = sortedViolations.value
  if (filterType.value) list = list.filter(v => v.violation_type === filterType.value)
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(v => (v.user_name || '').toLowerCase().includes(q) || (v.participant_id || '').toLowerCase().includes(q))
  }
  if (selectedParticipantId.value) {
    list = list.filter(v => v.participant_id === selectedParticipantId.value)
  }
  return list
})

const participantSummary = computed(() => {
  const map: Record<string, { id: string; name: string; total: number; byType: Record<string, number> }> = {}
  for (const v of violations.value) {
    const pid = v.participant_id
    if (!map[pid]) map[pid] = { id: pid, name: v.user_name || pid.substring(0, 8), total: 0, byType: {} }
    map[pid].total++
    map[pid].byType[v.violation_type] = (map[pid].byType[v.violation_type] || 0) + 1
  }
  return Object.values(map).sort((a, b) => b.total - a.total)
})

const summaryStats = computed(() => {
  const critical = violations.value.filter(v => v.violation_type === 'multiple_face').length
  const warnings = violations.value.filter(v => ['no_face', 'tab_switch'].includes(v.violation_type)).length
  const uniqueParticipants = new Set(violations.value.map(v => v.participant_id)).size
  return [
    { label: 'Total Pelanggaran', value: violations.value.length, sublabel: 'Semua jenis', icon: ShieldAlertIcon, iconBg: 'bg-slate-100', iconColor: 'text-slate-500', badgeClass: 'bg-slate-100 text-slate-600' },
    { label: 'Critical', value: critical, sublabel: 'Multiple Faces', icon: AlertTriangleIcon, iconBg: 'bg-red-50', iconColor: 'text-red-500', badgeClass: 'bg-red-100 text-red-700' },
    { label: 'Warnings', value: warnings, sublabel: 'Tab switch / No face', icon: EyeOffIcon, iconBg: 'bg-orange-50', iconColor: 'text-orange-500', badgeClass: 'bg-orange-100 text-orange-700' },
    { label: 'Peserta Terlibat', value: uniqueParticipants, sublabel: 'Dari total peserta', icon: ScanEyeIcon, iconBg: 'bg-blue-50', iconColor: 'text-blue-500', badgeClass: 'bg-blue-100 text-blue-700' },
  ]
})

const formatViolationType = (type: string) => {
  const map: Record<string, string> = {
    no_face: 'No Face Detected',
    multiple_face: 'Multiple Faces',
    tab_switch: 'Tab Switched',
    fullscreen_exit: 'Fullscreen Exit',
    looking_away: 'Looking Away',
    camera_off: 'Camera Off',
  }
  return map[type] || type
}

const formatType = (type: string) => {
  const map: Record<string, string> = {
    no_face: 'No Face',
    multiple_face: 'Multi-Face',
    tab_switch: 'Tab',
    fullscreen_exit: 'Fullscreen',
    looking_away: 'Away',
    camera_off: 'Cam Off',
  }
  return map[type] || type
}

const formatTime = (dt: string) => {
  if (!dt) return '-'
  return new Date(dt).toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'medium' })
}

const getViolationIcon = (type: string) => {
  if (type === 'multiple_face') return AlertTriangleIcon
  if (type === 'no_face') return EyeOffIcon
  if (type === 'tab_switch') return MonitorOffIcon
  if (type === 'fullscreen_exit') return MaximizeIcon
  if (type === 'camera_off') return CameraOffIcon
  return ShieldAlertIcon
}

const getSeverityBar = (type: string) => {
  if (type === 'multiple_face') return 'bg-red-500'
  if (['no_face', 'tab_switch'].includes(type)) return 'bg-orange-400'
  return 'bg-yellow-400'
}

const getTypeBg = (type: string) => {
  if (type === 'multiple_face') return 'bg-red-50'
  if (['no_face', 'tab_switch'].includes(type)) return 'bg-orange-50'
  return 'bg-yellow-50'
}

const getTypeIconColor = (type: string) => {
  if (type === 'multiple_face') return 'text-red-500'
  if (['no_face', 'tab_switch'].includes(type)) return 'text-orange-500'
  return 'text-yellow-500'
}

const getTypeColor = (type: string) => {
  if (type === 'multiple_face') return 'bg-red-100 text-red-700'
  if (['no_face', 'tab_switch'].includes(type)) return 'bg-orange-100 text-orange-700'
  return 'bg-yellow-100 text-yellow-700'
}

const getSeverityBadge = (type: string) => {
  if (type === 'multiple_face') return 'bg-red-100 text-red-700'
  if (['no_face', 'tab_switch'].includes(type)) return 'bg-orange-100 text-orange-700'
  return 'bg-yellow-100 text-yellow-700'
}

const getSeverityLabel = (type: string) => {
  if (type === 'multiple_face') return 'CRITICAL'
  if (['no_face', 'tab_switch'].includes(type)) return 'WARNING'
  return 'INFO'
}

const exportCSV = () => {
  const rows = [
    ['No', 'Peserta', 'Participant ID', 'Jenis Pelanggaran', 'Severity', 'Waktu Deteksi'],
    ...violations.value.map((v, i) => [
      i + 1,
      v.user_name || '',
      v.participant_id,
      formatViolationType(v.violation_type),
      getSeverityLabel(v.violation_type),
      formatTime(v.detected_at)
    ])
  ]
  const csv = rows.map(r => r.join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `violations-${sessionId.substring(0, 8)}.csv`
  a.click()
}
</script>
