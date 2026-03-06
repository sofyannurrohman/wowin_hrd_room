<template>
  <div class="space-y-6 max-w-[1200px] mx-auto pb-10">
    <!-- Header -->
    <div class="flex items-center gap-4">
      <Button variant="ghost" class="h-10 w-10 p-0 rounded-full" @click="router.back()">
        <ArrowLeftIcon class="w-5 h-5 text-slate-500" />
      </Button>
      <div>
        <h2 class="text-3xl font-extrabold tracking-tight text-slate-800">Detail Peserta</h2>
        <p class="text-slate-500 mt-1">Informasi lengkap profil dan riwayat ujian peserta.</p>
      </div>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-20 text-slate-400">
      <svg class="w-8 h-8 animate-spin mr-3 text-blue-500" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
      </svg>
      Memuat detail peserta...
    </div>

    <div v-else-if="!profile" class="py-20 text-center">
      <UserXIcon class="w-12 h-12 text-slate-300 mx-auto mb-3" />
      <h3 class="text-lg font-bold text-slate-700">Peserta tidak ditemukan</h3>
      <p class="text-slate-500">Data peserta mungkin telah dihapus atau ID tidak valid.</p>
      <Button class="mt-4" @click="router.push('/participants')">Kembali ke Daftar</Button>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Profile Card -->
      <Card class="col-span-1 border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden">
        <CardContent class="p-6">
          <div class="flex flex-col items-center text-center pb-6 border-b border-slate-100">
            <img
              :src="`https://api.dicebear.com/7.x/initials/svg?seed=${profile.name || profile.email}&backgroundColor=1e3a8a&textColor=ffffff`"
              class="w-24 h-24 rounded-full mb-4 shadow-md border-4 border-white"
            />
            <h3 class="text-xl font-bold text-slate-800">{{ profile.name }}</h3>
            <p class="text-sm font-medium text-blue-600 mt-1">{{ profile.applied_position || 'Posisi Umum' }}</p>
            <Badge class="mt-3 bg-green-100 text-green-700 hover:bg-green-100 border-0 font-semibold px-3 py-1">Active</Badge>
          </div>

          <div class="space-y-4 pt-6">
            <div class="flex items-start gap-4">
              <MailIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div>
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-0.5">Email</div>
                <div class="text-sm font-medium text-slate-800 break-all">{{ profile.email }}</div>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <PhoneIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div>
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-0.5">WhatsApp</div>
                <div class="text-sm font-medium text-slate-800">
                  {{ profile.whatsapp_number || 'Tidak tersedia' }}
                </div>
              </div>
            </div>

             <div class="flex items-start gap-4">
              <UserIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div>
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-0.5">Umur</div>
                <div class="text-sm font-medium text-slate-800">
                  {{ profile.age ? `${profile.age} Tahun` : 'Tidak tersedia' }}
                </div>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <GraduationCapIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div>
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-0.5">Pendidikan Terakhir</div>
                <div class="text-sm font-medium text-slate-800">
                  {{ profile.last_education || 'Tidak tersedia' }}
                </div>
              </div>
            </div>

            <div class="flex items-start gap-4">
              <MapPinIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div>
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-0.5">Alamat Lengkap</div>
                <div class="text-sm font-medium text-slate-800 leading-relaxed">
                  {{ profile.address || 'Tidak tersedia' }}
                </div>
              </div>
            </div>
            
            <div class="flex items-start gap-4">
              <BanknoteIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div>
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-0.5">Ekspektasi Gaji</div>
                <div class="text-sm font-medium text-slate-800">
                  {{ profile.expected_salary ? `Rp ${profile.expected_salary}` : 'Tidak tersedia' }}
                </div>
              </div>
            </div>
            
            <div class="flex items-start gap-4 pt-4 border-t border-slate-100">
              <FileTextIcon class="w-5 h-5 text-slate-400 mt-0.5 shrink-0" />
              <div class="w-full">
                <div class="text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-2">Curriculum Vitae (CV)</div>
                <Button v-if="profile.cv_url" variant="outline" class="w-full justify-center flex items-center gap-2 border-blue-200 bg-blue-50 text-blue-700 hover:bg-blue-100 font-bold" @click="downloadCV(profile.cv_url)">
                  <DownloadIcon class="w-4 h-4" /> Lihat / Unduh CV
                </Button>
                <div v-else class="text-sm font-medium text-slate-400 italic bg-slate-50 p-3 rounded-lg text-center border border-slate-100">
                  Tidak ada CV yang diunggah
                </div>
              </div>
            </div>

          </div>
        </CardContent>
      </Card>

      <!-- History Table -->
      <div class="col-span-1 lg:col-span-2 space-y-6">
        <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden">
          <div class="p-6 border-b border-slate-100 flex items-center justify-between">
            <div>
              <h3 class="text-lg font-bold text-slate-800 flex items-center gap-2">
                <HistoryIcon class="w-5 h-5 text-blue-600" /> Riwayat Ujian
              </h3>
              <p class="text-sm text-slate-500 mt-1">Daftar semua sesi ujian yang diikuti.</p>
            </div>
            <div class="w-10 h-10 rounded-full bg-blue-50 text-blue-600 flex items-center justify-center font-bold text-sm">
              {{ history.length }}
            </div>
          </div>
          
          <Table v-if="history.length > 0">
            <TableHeader>
              <TableRow class="border-b border-slate-100 hover:bg-transparent">
                <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6">Sesi Ujian</TableHead>
                <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Tanggal</TableHead>
                <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Status & Skor</TableHead>
                <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6 text-right">Detil</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="h in history"
                :key="h.session_id"
                class="border-b border-slate-100 transition-colors hover:bg-slate-50/50 group"
              >
                <TableCell class="py-4 px-6">
                  <div class="font-bold text-slate-800 text-sm max-w-[200px] truncate" :title="h.session_name">{{ h.session_name }}</div>
                </TableCell>
                <TableCell class="py-4 px-4">
                  <div class="text-sm text-slate-600 font-medium whitespace-nowrap">
                    {{ formatDate(h.schedule) }}
                  </div>
                </TableCell>
                <TableCell class="py-4 px-4">
                  <div class="flex flex-col gap-1.5 items-start">
                    <Badge :class="getStatusBadgeClass(h.status)" class="font-medium text-[10px] px-2 py-0 border-0">
                      {{ (h.status || 'Active').toUpperCase() }}
                    </Badge>
                    <div v-if="h.score !== null && h.score !== undefined" class="flex items-center gap-2 w-full mt-1">
                      <span class="font-bold text-slate-800 text-sm w-8">{{ Math.round(h.score) }}%</span>
                      <div class="w-16 h-1.5 bg-slate-100 rounded-full overflow-hidden shrink-0">
                        <div class="h-full rounded-full transition-all" :class="getScoreBarColor(h.score)" :style="{ width: `${Math.min(100, h.score)}%` }"></div>
                      </div>
                    </div>
                    <div v-else class="text-xs text-slate-400 italic">Belum ada skor</div>
                  </div>
                </TableCell>
                <TableCell class="py-4 px-6 text-right">
                  <Button 
                    v-if="h.status === 'finished' || h.score !== null"
                    variant="outline" size="sm" 
                    class="h-8 rounded-lg border-blue-200 text-blue-600 hover:bg-blue-50 font-semibold text-xs transition-colors"
                    @click="viewAnswers(h.session_id, h.participant_id)"
                  >
                    Lihat Jawaban
                  </Button>
                  <span v-else class="text-xs text-slate-300 italic">—</span>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>

          <div v-else class="py-16 text-center text-slate-400">
            <CalendarXIcon class="w-12 h-12 mx-auto mb-3 opacity-20" />
            <p class="font-bold text-slate-600">Belum ada riwayat ujian</p>
            <p class="text-sm mt-1">Peserta ini belum ditugaskan atau mengikuti ujian apapun.</p>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import client from '@/api/client'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import {
  ArrowLeftIcon, MailIcon, PhoneIcon, UserIcon, MapPinIcon, BanknoteIcon,
  GraduationCapIcon, FileTextIcon, DownloadIcon, HistoryIcon, CalendarXIcon, UserXIcon
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const participantId = route.params.id as string

const loading = ref(true)
const profile = ref<any>(null)
const history = ref<any[]>([])

const loadData = async () => {
  if (!participantId) return
  loading.value = true
  try {
    const res = await client.get(`/participants/${participantId}`)
    profile.value = res.data.profile
    history.value = res.data.history || []
  } catch (err: any) {
    console.error('Failed to load participant detail', err)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const base = new Date(dateString)
  if (base.getFullYear() === 1970) return '—'
  return base.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

const getStatusBadgeClass = (status: string) => {
  const s = (status || '').toLowerCase()
  if (s === 'finished') return 'bg-green-100 text-green-700'
  if (s === 'active') return 'bg-blue-100 text-blue-700'
  if (s === 'disconnected') return 'bg-red-100 text-red-700'
  return 'bg-slate-100 text-slate-600'
}

const getScoreBarColor = (score: number) => {
  if (score >= 80) return 'bg-green-500'
  if (score >= 60) return 'bg-blue-500'
  if (score >= 40) return 'bg-orange-400'
  return 'bg-red-400'
}

const viewAnswers = (sessionId: string, spId: string) => {
  router.push(`/sessions/${sessionId}/results/${spId}/answers`)
}

const downloadCV = (url: string) => {
  const backendUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const baseUrl = backendUrl.replace('/api', '')
  window.open(`${baseUrl}${url}`, '_blank')
}
</script>
