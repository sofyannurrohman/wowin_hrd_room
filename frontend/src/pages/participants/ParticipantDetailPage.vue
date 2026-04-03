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
                <div class="text-sm font-medium text-slate-800 flex items-center gap-2">
                  <span v-if="!profile.whatsapp_number">Tidak tersedia</span>
                  <template v-else>
                    <a :href="getWaLink(profile.whatsapp_number)" target="_blank" class="hover:text-emerald-600 hover:underline transition-colors">{{ profile.whatsapp_number }}</a>
                    <a :href="getWaLink(profile.whatsapp_number)" target="_blank" class="inline-flex items-center justify-center px-2.5 py-0.5 text-[10px] font-bold rounded-full bg-emerald-50 text-emerald-600 border border-emerald-200 hover:bg-emerald-100 transition-colors">
                      <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 24 24"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51a12.8 12.8 0 0 0-.57-.01c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"/></svg>
                      Hubungi via WA
                    </a>
                  </template>
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
                  <div class="flex justify-end gap-2 items-center">
                    <Button 
                      v-if="h.ktp_selfie_url"
                      variant="outline" size="sm" 
                      class="h-8 rounded-lg border-emerald-200 text-emerald-600 hover:bg-emerald-50 font-semibold text-xs transition-colors"
                      @click="viewSelfie(h.ktp_selfie_url)"
                      title="Lihat Foto KTP/Selfie untuk Sesi Ujian Ini"
                    >
                      <ImageIcon class="w-3 h-3 mr-1" /> KTP
                    </Button>
                    <Button 
                      v-if="h.status === 'finished' || h.score !== null"
                      variant="outline" size="sm" 
                      class="h-8 rounded-lg border-blue-200 text-blue-600 hover:bg-blue-50 font-semibold text-xs transition-colors"
                      @click="viewAnswers(h.session_id, h.participant_id)"
                    >
                      Jawaban
                    </Button>
                    <span v-if="!h.ktp_selfie_url && h.status !== 'finished' && h.score === null" class="text-xs text-slate-300 italic">—</span>
                  </div>
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

    <!-- Selfie Preview Modal -->
    <Teleport to="body">
      <div v-if="previewSelfieUrl" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="previewSelfieUrl = null">
        <div class="absolute inset-0 bg-black/80 backdrop-blur-sm" @click="previewSelfieUrl = null"></div>
        <div class="relative bg-white rounded-lg p-2 z-10 max-w-3xl max-h-[90vh] overflow-hidden flex flex-col items-center">
          <button @click="previewSelfieUrl = null" class="absolute top-4 right-4 w-8 h-8 bg-black/50 text-white rounded-full flex items-center justify-center hover:bg-black/70 transition-colors">✕</button>
          <img :src="previewSelfieUrl" class="w-auto h-auto max-w-full max-h-[85vh] object-contain rounded" alt="KTP Selfie Preview" />
        </div>
      </div>
    </Teleport>
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
  GraduationCapIcon, FileTextIcon, DownloadIcon, HistoryIcon, CalendarXIcon, UserXIcon, ImageIcon
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const participantId = route.params.id as string

const loading = ref(true)
const profile = ref<any>(null)
const history = ref<any[]>([])
const previewSelfieUrl = ref<string | null>(null)

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

const getWaLink = (phone: string) => {
  if (!phone) return '#';
  let formatted = phone.replace(/\D/g, '');
  if (formatted.startsWith('0')) {
    formatted = '62' + formatted.slice(1);
  }
  return `https://wa.me/${formatted}`;
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

const viewSelfie = (url: string) => {
  if (!url) return
  const backendUrl = import.meta.env.VITE_API_URL
  if (backendUrl) {
    const baseUrl = backendUrl.replace(/\/api$/, '')
    previewSelfieUrl.value = `${baseUrl}${url}`
  } else {
    previewSelfieUrl.value = url
  }
}

const downloadCV = (url: string) => {
  if (!url) return
  
  // Use VITE_API_URL if provided (removing /api for static uploads), 
  // otherwise fallback to a relative path which works via Nginx/Vite proxy.
  const backendUrl = import.meta.env.VITE_API_URL
  if (backendUrl) {
    const baseUrl = backendUrl.replace(/\/api$/, '')
    window.open(`${baseUrl}${url}`, '_blank')
  } else {
    window.open(url, '_blank')
  }
}
</script>
