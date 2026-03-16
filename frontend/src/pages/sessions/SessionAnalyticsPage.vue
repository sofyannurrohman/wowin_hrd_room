<template>
  <div class="space-y-6 max-w-7xl mx-auto pb-10">

    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()" class="rounded-full">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <div>
          <h2 class="text-2xl font-extrabold tracking-tight text-slate-800">Analitik Hasil & Kandidat Terbaik</h2>
          <p class="text-sm text-slate-500 mt-0.5">Sesi: <span class="font-semibold text-slate-700">{{ sessionName }}</span></p>
        </div>
      </div>
      <div class="flex items-center gap-3">
        <Badge variant="outline" class="text-xs font-semibold border-slate-200 text-slate-600 px-3 py-1 rounded-full">
          {{ results.length }} Peserta Selesai
        </Badge>
        <Button variant="outline" size="sm" @click="exportCSV" class="gap-2 text-xs font-bold">
          <DownloadIcon class="w-3.5 h-3.5" /> Export CSV
        </Button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="flex flex-col items-center gap-3 text-slate-400">
        <svg class="w-8 h-8 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
        </svg>
        <span class="text-sm font-medium">Memuat data hasil...</span>
      </div>
    </div>

    <template v-else>
      <!-- Summary Stats -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div v-for="stat in summaryStats" :key="stat.label" class="bg-white rounded-2xl border border-slate-100 shadow-sm p-5 hover:shadow-md transition-all">
          <div class="flex items-center justify-between mb-2">
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">{{ stat.label }}</p>
            <div class="w-8 h-8 rounded-lg flex items-center justify-center" :class="stat.iconBg">
              <component :is="stat.icon" class="w-4 h-4" :class="stat.iconColor" />
            </div>
          </div>
          <p class="text-3xl font-black text-slate-800">{{ stat.value }}</p>
          <p class="text-xs text-slate-500 mt-1">{{ stat.sublabel }}</p>
        </div>
      </div>

      <!-- Score Distribution Visual -->
      <div class="bg-white border border-slate-100 rounded-2xl shadow-sm p-6">
        <h3 class="text-base font-bold text-slate-800 mb-4 flex items-center gap-2">
          <BarChart3Icon class="w-5 h-5 text-blue-500" />
          Distribusi Skor
        </h3>
        <div class="flex items-end gap-2 h-32">
          <div
            v-for="bucket in scoreBuckets"
            :key="bucket.label"
            class="flex-1 flex flex-col items-center gap-1"
          >
            <span class="text-[10px] font-bold text-slate-700">{{ bucket.count }}</span>
            <div
              class="w-full rounded-t-md transition-all duration-500"
              :class="bucket.color"
              :style="{ height: bucket.count > 0 ? `${Math.max(8, (bucket.count / maxBucketCount) * 108)}px` : '4px' }"
            ></div>
            <span class="text-[10px] text-slate-500 font-medium">{{ bucket.label }}</span>
          </div>
        </div>
        <div class="flex gap-4 mt-4 flex-wrap">
          <div class="flex items-center gap-1.5 text-xs text-slate-600">
            <span class="w-3 h-3 rounded-sm bg-red-300 inline-block"></span> 0-39 (Kurang)
          </div>
          <div class="flex items-center gap-1.5 text-xs text-slate-600">
            <span class="w-3 h-3 rounded-sm bg-orange-300 inline-block"></span> 40-59 (Cukup)
          </div>
          <div class="flex items-center gap-1.5 text-xs text-slate-600">
            <span class="w-3 h-3 rounded-sm bg-blue-300 inline-block"></span> 60-79 (Baik)
          </div>
          <div class="flex items-center gap-1.5 text-xs text-slate-600">
            <span class="w-3 h-3 rounded-sm bg-green-400 inline-block"></span> 80-100 (Sangat Baik)
          </div>
        </div>
      </div>

      <!-- Filter + Search -->
      <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3">
        <div class="relative flex-1">
          <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            v-model="search"
            type="text"
            placeholder="Cari nama peserta..."
            class="w-full pl-9 pr-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <select
          v-model="filterStatus"
          class="px-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-slate-700 font-medium"
        >
          <option value="">Semua Status</option>
          <option value="completed">Selesai</option>
          <option value="pending_review">Perlu Review</option>
        </select>
        <select
          v-model="sortBy"
          class="px-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-slate-700 font-medium"
        >
          <option value="score_desc">Skor Tertinggi</option>
          <option value="score_asc">Skor Terendah</option>
          <option value="name">Nama A-Z</option>
          <option value="time_desc">Terbaru Submit</option>
        </select>
      </div>

      <!-- Candidate Ranking Table -->
      <div class="bg-white border border-slate-100 rounded-2xl shadow-sm overflow-hidden">
        <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
          <h3 class="text-base font-bold text-slate-800 flex items-center gap-2">
            <TrophyIcon class="w-5 h-5 text-amber-500" />
            Ranking Kandidat
          </h3>
          <span class="text-xs text-slate-500">{{ filteredResults.length }} ditampilkan</span>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-slate-50/70 border-b border-slate-100">
              <tr>
                <th class="px-6 py-3 text-left text-[11px] font-bold text-slate-500 uppercase tracking-wider">Rank</th>
                <th class="px-4 py-3 text-left text-[11px] font-bold text-slate-500 uppercase tracking-wider">Peserta</th>
                <th class="px-4 py-3 text-left text-[11px] font-bold text-slate-500 uppercase tracking-wider">Email</th>
                <th class="px-4 py-3 text-center text-[11px] font-bold text-slate-500 uppercase tracking-wider">Skor</th>
                <th class="px-4 py-3 text-center text-[11px] font-bold text-slate-500 uppercase tracking-wider">Nilai Bar</th>
                <th class="px-4 py-3 text-center text-[11px] font-bold text-slate-500 uppercase tracking-wider">Status</th>
                <th class="px-4 py-3 text-center text-[11px] font-bold text-slate-500 uppercase tracking-wider">Pelanggaran</th>
                <th class="px-4 py-3 text-center text-[11px] font-bold text-slate-500 uppercase tracking-wider">Rekomendasi</th>
                <th class="px-6 py-3 text-right text-[11px] font-bold text-slate-500 uppercase tracking-wider">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr
                v-for="(res, idx) in paginatedData"
                :key="res.id"
                class="hover:bg-slate-50/50 transition-colors"
                :class="idx === 0 && sortBy === 'score_desc' ? 'bg-amber-50/40' : ''"
              >
                <!-- Rank -->
                <td class="px-6 py-4">
                  <div class="flex items-center justify-center w-8 h-8 rounded-full font-black text-sm"
                    :class="getRankStyle((currentPage - 1) * pageSize + idx + 1)">
                    {{ (currentPage - 1) * pageSize + idx + 1 }}
                  </div>
                </td>
                
                <!-- Name -->
                <td class="px-4 py-4">
                  <div class="flex items-center gap-3">
                    <img
                      :src="`https://api.dicebear.com/7.x/initials/svg?seed=${res.user_name}&backgroundColor=1e3a8a&textColor=ffffff`"
                      class="w-8 h-8 rounded-full shrink-0"
                    />
                    <div>
                      <p class="text-sm font-bold text-slate-800 leading-tight">{{ res.user_name || 'Unknown' }}</p>
                      <p class="text-[10px] text-slate-400 font-mono">{{ res.participant_id?.substring(0, 8).toUpperCase() }}</p>
                    </div>
                  </div>
                </td>

                <!-- Email -->
                <td class="px-4 py-4">
                  <p class="text-sm text-slate-500">{{ res.user_email || '-' }}</p>
                </td>

                <!-- Score -->
                <td class="px-4 py-4 text-center">
                  <span class="text-xl font-black" :class="getScoreColor(res.total_score)">
                    {{ res.total_score?.toFixed(1) || '0.0' }}
                  </span>
                </td>

                <!-- Score Bar -->
                <td class="px-4 py-4">
                  <div class="w-32 h-2.5 bg-slate-100 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full transition-all duration-700"
                      :class="getScoreBarColor(res.total_score)"
                      :style="{ width: `${Math.min(100, (res.total_score / maxScore) * 100)}%` }"
                    ></div>
                  </div>
                </td>

                <!-- Status -->
                <td class="px-4 py-4 text-center">
                  <Badge
                    :class="res.grading_status === 'completed' ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'"
                    class="text-[10px] font-bold px-2 py-0.5 rounded-full border-0"
                  >
                    {{ res.grading_status === 'completed' ? 'Selesai' : 'Perlu Review' }}
                  </Badge>
                </td>

                <!-- Violations -->
                <td class="px-4 py-4 text-center">
                  <span
                    class="text-sm font-bold"
                    :class="(violationsByParticipant[res.participant_id] || 0) > 3 ? 'text-red-600' : (violationsByParticipant[res.participant_id] || 0) > 0 ? 'text-orange-500' : 'text-slate-400'"
                  >
                    {{ violationsByParticipant[res.participant_id] || 0 }}
                  </span>
                </td>

                <!-- Recommendation -->
                <td class="px-4 py-4 text-center">
                  <Badge :class="getRecommendationStyle(res)" class="text-[10px] font-bold px-2.5 py-1 rounded-full border-0">
                    {{ getRecommendation(res) }}
                  </Badge>
                </td>

                <!-- Actions -->
                <td class="px-6 py-4 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <Button
                      v-if="res.grading_status !== 'completed'"
                      size="sm"
                      variant="outline"
                      class="text-xs h-7 px-3 font-bold"
                      @click="finalizeScore(res)"
                      :disabled="finalizing === res.id"
                    >
                      {{ finalizing === res.id ? 'Saving...' : 'Finalisasi' }}
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      class="text-xs h-7 px-2 text-slate-500 hover:text-blue-600"
                      @click="router.push(`/sessions/${sessionId}/results/${res.participant_id}/answers`)"
                    >
                      <EyeIcon class="w-3.5 h-3.5" />
                    </Button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredResults.length === 0">
                <td colspan="9" class="text-center py-16 text-slate-400">
                  <TrophyIcon class="w-10 h-10 mx-auto mb-3 opacity-20" />
                  <p class="font-semibold">Belum ada hasil yang masuk.</p>
                  <p class="text-xs mt-1">Pastikan peserta telah menyelesaikan ujian.</p>
                </td>
              </tr>
            </tbody>
          </table>
          
          <div v-if="filteredResults.length > 0" class="p-4 border-t border-slate-100 bg-slate-50/30">
            <DataTablePagination 
              :total="totalItems"
              v-model:pageSize="pageSize"
              v-model:currentPage="currentPage"
            />
          </div>
        </div>
      </div>

      <!-- Top 3 Podium (shown when sorted by score) -->
      <div v-if="sortBy === 'score_desc' && filteredResults.length >= 3" class="bg-gradient-to-br from-slate-900 to-indigo-900 rounded-2xl p-8 text-white">
        <h3 class="text-base font-bold text-white/80 mb-6 flex items-center gap-2">
          <TrophyIcon class="w-5 h-5 text-amber-400" />
          Top 3 Kandidat
        </h3>
        <div class="flex items-end justify-center gap-6">
          <!-- 2nd place -->
          <div class="flex flex-col items-center gap-2 flex-1">
            <img :src="`https://api.dicebear.com/7.x/initials/svg?seed=${filteredResults[1]?.user_name}&backgroundColor=334155&textColor=ffffff`" class="w-14 h-14 rounded-full border-2 border-slate-500" />
            <p class="text-sm font-bold text-center leading-tight">{{ filteredResults[1]?.user_name }}</p>
            <div class="w-full bg-slate-600/50 rounded-t-xl py-5 flex flex-col items-center gap-1">
              <span class="text-2xl font-black text-slate-300">2</span>
              <span class="text-lg font-black text-white">{{ filteredResults[1]?.total_score?.toFixed(1) }}</span>
            </div>
          </div>
          <!-- 1st place -->
          <div class="flex flex-col items-center gap-2 flex-1">
            <div class="relative">
              <img :src="`https://api.dicebear.com/7.x/initials/svg?seed=${filteredResults[0]?.user_name}&backgroundColor=1e40af&textColor=ffffff`" class="w-18 h-18 rounded-full border-4 border-amber-400 w-20 h-20" />
              <div class="absolute -top-2 left-1/2 -translate-x-1/2 text-xl">👑</div>
            </div>
            <p class="text-sm font-bold text-center leading-tight text-amber-300">{{ filteredResults[0]?.user_name }}</p>
            <div class="w-full bg-amber-500/30 rounded-t-xl py-8 flex flex-col items-center gap-1 border-t-2 border-amber-400">
              <span class="text-3xl font-black text-amber-400">1</span>
              <span class="text-xl font-black text-white">{{ filteredResults[0]?.total_score?.toFixed(1) }}</span>
            </div>
          </div>
          <!-- 3rd place -->
          <div class="flex flex-col items-center gap-2 flex-1">
            <img :src="`https://api.dicebear.com/7.x/initials/svg?seed=${filteredResults[2]?.user_name}&backgroundColor=44403c&textColor=ffffff`" class="w-14 h-14 rounded-full border-2 border-amber-700/60" />
            <p class="text-sm font-bold text-center leading-tight">{{ filteredResults[2]?.user_name }}</p>
            <div class="w-full bg-amber-800/30 rounded-t-xl py-3 flex flex-col items-center gap-1">
              <span class="text-2xl font-black text-amber-700">3</span>
              <span class="text-lg font-black text-white">{{ filteredResults[2]?.total_score?.toFixed(1) }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  ArrowLeftIcon, DownloadIcon, SearchIcon,
  TrophyIcon, BarChart3Icon, EyeIcon, UsersIcon, ClipboardListIcon
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const sessionName = ref('Loading...')
const results = ref<any[]>([])
const violations = ref<any[]>([])
const loading = ref(true)
const search = ref('')
const filterStatus = ref('')
const sortBy = ref('score_desc')
const finalizing = ref<string | null>(null)

onMounted(async () => {
  try {
    const [sRes, rRes, vRes] = await Promise.all([
      client.get(`/sessions/${sessionId}`),
      client.get(`/sessions/${sessionId}/results`),
      client.get(`/sessions/${sessionId}/violations`)
    ])
    sessionName.value = sRes.data.name
    results.value = rRes.data || []
    violations.value = vRes.data || []
  } catch (e) {
    console.error('Failed to load analytics', e)
  } finally {
    loading.value = false
  }
})

const violationsByParticipant = computed(() => {
  const map: Record<string, number> = {}
  for (const v of violations.value) {
    map[v.participant_id] = (map[v.participant_id] || 0) + 1
  }
  return map
})

const maxScore = computed(() => {
  return Math.max(...results.value.map(r => r.total_score || 0), 1)
})

const avgScore = computed(() => {
  if (!results.value.length) return 0
  return results.value.reduce((s, r) => s + (r.total_score || 0), 0) / results.value.length
})

const summaryStats = computed(() => [
  {
    label: 'Total Peserta',
    value: results.value.length,
    sublabel: 'Telah menyelesaikan ujian',
    icon: UsersIcon,
    iconBg: 'bg-blue-50',
    iconColor: 'text-blue-500',
  },
  {
    label: 'Skor Tertinggi',
    value: maxScore.value.toFixed(1),
    sublabel: filteredResults.value[0]?.user_name || '-',
    icon: TrophyIcon,
    iconBg: 'bg-amber-50',
    iconColor: 'text-amber-500',
  },
  {
    label: 'Rata-rata Skor',
    value: avgScore.value.toFixed(1),
    sublabel: 'Dari semua peserta',
    icon: BarChart3Icon,
    iconBg: 'bg-green-50',
    iconColor: 'text-green-500',
  },
  {
    label: 'Perlu Review',
    value: results.value.filter(r => r.grading_status !== 'completed').length,
    sublabel: 'Jawaban esai manual',
    icon: ClipboardListIcon,
    iconBg: 'bg-orange-50',
    iconColor: 'text-orange-500',
  },
])

const scoreBuckets = computed(() => {
  const buckets = [
    { label: '0-9', min: 0, max: 10, count: 0, color: 'bg-red-300' },
    { label: '10-19', min: 10, max: 20, count: 0, color: 'bg-red-300' },
    { label: '20-29', min: 20, max: 30, count: 0, color: 'bg-red-300' },
    { label: '30-39', min: 30, max: 40, count: 0, color: 'bg-red-300' },
    { label: '40-49', min: 40, max: 50, count: 0, color: 'bg-orange-300' },
    { label: '50-59', min: 50, max: 60, count: 0, color: 'bg-orange-300' },
    { label: '60-69', min: 60, max: 70, count: 0, color: 'bg-blue-300' },
    { label: '70-79', min: 70, max: 80, count: 0, color: 'bg-blue-300' },
    { label: '80-89', min: 80, max: 90, count: 0, color: 'bg-green-400' },
    { label: '90-100', min: 90, max: 101, count: 0, color: 'bg-green-600' },
  ]
  for (const r of results.value) {
    const score = r.total_score || 0
    for (const bucket of buckets) {
      if (score >= bucket.min && score < bucket.max) { bucket.count++; break }
    }
  }
  return buckets
})

const maxBucketCount = computed(() => Math.max(...scoreBuckets.value.map(b => b.count), 1))

const sortedResults = computed(() => {
  const list = [...results.value]
  if (sortBy.value === 'score_desc') list.sort((a, b) => (b.total_score || 0) - (a.total_score || 0))
  if (sortBy.value === 'score_asc') list.sort((a, b) => (a.total_score || 0) - (b.total_score || 0))
  if (sortBy.value === 'name') list.sort((a, b) => (a.user_name || '').localeCompare(b.user_name || ''))
  if (sortBy.value === 'time_desc') list.sort((a, b) => new Date(b.submitted_at).getTime() - new Date(a.submitted_at).getTime())
  return list
})

const filteredResults = computed(() => {
  let list = sortedResults.value
  if (filterStatus.value) list = list.filter(r => r.grading_status === filterStatus.value)
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(r => (r.user_name || '').toLowerCase().includes(q) || (r.user_email || '').toLowerCase().includes(q))
  }
  return list
})

const {
  pageSize,
  currentPage,
  paginatedData,
  totalItems
} = useDataTable(filteredResults)

const getScoreColor = (score: number) => {
  if (score >= 80) return 'text-green-600'
  if (score >= 60) return 'text-blue-600'
  if (score >= 40) return 'text-orange-500'
  return 'text-red-600'
}

const getScoreBarColor = (score: number) => {
  if (score >= 80) return 'bg-green-500'
  if (score >= 60) return 'bg-blue-500'
  if (score >= 40) return 'bg-orange-400'
  return 'bg-red-400'
}

const getRankStyle = (rank: number) => {
  if (rank === 1) return 'bg-amber-100 text-amber-700 ring-2 ring-amber-400'
  if (rank === 2) return 'bg-slate-100 text-slate-600'
  if (rank === 3) return 'bg-orange-50 text-orange-700'
  return 'bg-slate-50 text-slate-500'
}

const getRecommendation = (res: any) => {
  const score = res.total_score || 0
  const viols = violationsByParticipant.value[res.participant_id] || 0
  if (score >= 80 && viols === 0) return '⭐ Sangat Disarankan'
  if (score >= 80 && viols <= 2) return '✅ Disarankan'
  if (score >= 60 && viols <= 3) return '🟡 Pertimbangkan'
  if (score >= 40) return '⚠️ Kurang'
  return '❌ Tidak Disarankan'
}

const getRecommendationStyle = (res: any) => {
  const score = res.total_score || 0
  const viols = violationsByParticipant.value[res.participant_id] || 0
  if (score >= 80 && viols === 0) return 'bg-green-100 text-green-800'
  if (score >= 80) return 'bg-teal-100 text-teal-800'
  if (score >= 60) return 'bg-yellow-100 text-yellow-800'
  if (score >= 40) return 'bg-orange-100 text-orange-700'
  return 'bg-red-100 text-red-700'
}

const finalizeScore = async (res: any) => {
  finalizing.value = res.id
  try {
    await client.post(`/results/${res.id}/finalize`)
    const rRes = await client.get(`/sessions/${sessionId}/results`)
    results.value = rRes.data || []
  } catch (e) {
    console.error('Failed to finalize', e)
  } finally {
    finalizing.value = null
  }
}

const exportCSV = () => {
  const rows = [
    ['Rank', 'Nama', 'Email', 'Skor', 'Status', 'Pelanggaran', 'Rekomendasi'],
    ...filteredResults.value.map((r, i) => [
      i + 1,
      r.user_name || '',
      r.user_email || '',
      (r.total_score || 0).toFixed(1),
      r.grading_status,
      violationsByParticipant.value[r.participant_id] || 0,
      getRecommendation(r).replace(/[⭐✅🟡⚠️❌]/g, '').trim(),
    ])
  ]
  const csv = rows.map(r => r.join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `results-analytics-${sessionId.substring(0, 8)}.csv`
  a.click()
}
</script>
