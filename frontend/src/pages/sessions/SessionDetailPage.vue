<template>
  <div class="space-y-6 max-w-5xl mx-auto" v-if="session">
    <div class="flex items-center gap-4 justify-between">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <h2 class="text-3xl font-bold tracking-tight">{{ session.name }}</h2>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="router.push(`/sessions/${session.id}/violations`)">
          <ShieldAlertIcon class="w-4 h-4 mr-2"/> Review Pelanggaran
        </Button>
        <Button variant="outline" @click="router.push(`/sessions/${session.id}/analytics`)">
          <BarChart3Icon class="w-4 h-4 mr-2"/> Analitik Hasil
        </Button>
        <Button variant="outline" @click="router.push(`/sessions/${session.id}/results`)">
          <FileTextIcon class="w-4 h-4 mr-2"/> Cek Hasil & Evaluasi
        </Button>
        <Button variant="default" @click="router.push(`/sessions/${session.id}/monitor`)">
          <MonitorPlayIcon class="w-4 h-4 mr-2"/> Buka Monitor Live
        </Button>
      </div>
    </div>

    <!-- Info -->
    <Card>
      <CardHeader>
        <CardTitle>Informasi Sesi</CardTitle>
      </CardHeader>
      <CardContent class="grid grid-cols-2 gap-4">
        <div>
          <Label class="text-muted-foreground">Waktu Pelaksanaan</Label>
          <p>{{ new Date(session.schedule).toLocaleString() }} - {{ new Date(new Date(session.schedule).getTime() + session.duration_minutes * 60000).toLocaleString() }}</p>
        </div>
        <div>
          <Label class="text-muted-foreground">Durasi</Label>
          <p>{{ session.duration_minutes }} Menit</p>
        </div>
        <div>
          <Label class="text-muted-foreground">Status / Peserta</Label>
          <p>{{ session.is_locked ? 'Terkunci' : 'Aktif' }} ({{ session.total_participant_count || 0 }} dari {{ session.max_participants }})</p>
        </div>
      </CardContent>
    </Card>

    <!-- Token Management -->
    <Card>
      <CardHeader>
        <CardTitle>Manajemen Token</CardTitle>
        <CardDescription>Buat token ujian baru dan distribusikan ke peserta.</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <Alert v-if="alertMessage" :variant="alertVariant">
          <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
          <AlertDescription>{{ alertMessage }}</AlertDescription>
        </Alert>

        <div class="flex items-center gap-4 pb-4 border-b flex-wrap">
           <Button variant="default" @click="openInviteModal" class="flex items-center gap-2">
             <MailIcon class="w-4 h-4" /> Kirim Undangan via Email
           </Button>
        </div>
        
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Token (Klik untuk copy)</TableHead>
              <TableHead>Sisa Kuota</TableHead>
              <TableHead>Kadaluarsa</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="t in tokens" :key="t.id">
              <TableCell class="font-mono" :class="t.token ? 'cursor-pointer hover:text-primary transition-colors font-bold' : 'text-muted-foreground'" @click="t.token ? copy(t.token) : null">
                {{ t.token || '*** HIDDEN ***' }}
              </TableCell>
              <TableCell>{{ t.max_usage - (t.used_count || 0) }} / {{ t.max_usage }}</TableCell>
              <TableCell>{{ t.expires_at ? new Date(t.expires_at).toLocaleString() : '-' }}</TableCell>
              <TableCell>
                <Badge :variant="t.is_active ? 'default' : 'destructive'">{{ t.is_active ? 'Aktif' : 'Non-Aktif' }}</Badge>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Participants List -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Daftar Peserta Bergabung</CardTitle>
            <CardDescription>Daftar peserta yang Sedang atau Telah mengikuti sesi ini.</CardDescription>
          </div>
          <Button variant="outline" size="sm" @click="loadParticipants" :disabled="loadingParticipants">
            <RefreshCwIcon class="w-4 h-4 mr-2" :class="{ 'animate-spin': loadingParticipants }" /> Segarkan
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Peserta</TableHead>
              <TableHead>Mulai</TableHead>
              <TableHead>Status</TableHead>
              <TableHead class="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="p in sessionParticipants" :key="p.id" class="group">
              <TableCell>
                <div class="flex items-center gap-3">
                  <img :src="p.ktp_selfie_url || `https://api.dicebear.com/7.x/initials/svg?seed=${p.user_name || '?'}&backgroundColor=1e3a8a&textColor=ffffff`" class="w-10 h-10 rounded-full object-cover border border-slate-200" />
                  <div>
                    <div class="font-bold text-slate-800 text-sm">{{ p.user_name || 'Unknown' }}</div>
                    <div class="text-xs text-slate-500">{{ p.user_email }}</div>
                  </div>
                </div>
              </TableCell>
              <TableCell class="text-xs text-slate-500">
                {{ p.joined_at ? new Date(p.joined_at).toLocaleString() : '-' }}
              </TableCell>
              <TableCell>
                <Badge :variant="p.status === 'finished' ? 'secondary' : (p.status === 'active' ? 'default' : 'outline')">
                  {{ p.status === 'active' ? 'Sedang Mengerjakan' : (p.status === 'finished' ? 'Selesai' : p.status) }}
                </Badge>
              </TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-2">
                   <Button variant="outline" size="sm" class="h-8 text-[11px] font-bold gap-1.5" @click="openMonitoring(p)">
                     <CameraIcon class="w-3.5 h-3.5" /> MONITORING
                   </Button>
                </div>
              </TableCell>
            </TableRow>
            <TableRow v-if="sessionParticipants.length === 0">
              <TableCell colspan="4" class="text-center py-10 text-slate-400 font-medium italic">Belum ada peserta yang bergabung ke sesi ini.</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Invite Participants Modal -->
    <Teleport to="body">
      <div v-if="showInviteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="showInviteModal = false">
        <div class="bg-white rounded-2xl shadow-xl w-full max-w-2xl flex flex-col max-h-[85vh]">
          <div class="p-5 border-b flex items-center justify-between shrink-0">
            <div>
              <h3 class="text-xl font-bold text-slate-800">Undang Peserta ke Sesi</h3>
              <p class="text-sm text-slate-500 mt-1">Pilih peserta untuk digenerate token dan dikirimkan Magic Link-nya via Email.</p>
            </div>
            <button @click="showInviteModal = false" class="text-slate-400 hover:bg-slate-100 rounded-full p-2">✕</button>
          </div>
          
          <div class="p-5 overflow-y-auto flex-1 bg-slate-50 relative">
            <div v-if="loadingParticipantsList" class="text-center py-10 text-slate-400">Memuat data peserta...</div>
            <div v-else-if="allParticipants.length === 0" class="text-center py-10 text-slate-400">Tidak ada peserta tersedia. Tembahkan peserta di menu Participant Management.</div>
            <div v-else class="space-y-2">
              <label v-for="p in allParticipants" :key="p.id" class="flex items-center gap-4 p-3 bg-white rounded-xl border border-slate-200 cursor-pointer hover:border-blue-400 transition-colors">
                <input type="checkbox" :value="p.id" v-model="selectedParticipantIds" class="w-5 h-5 rounded text-blue-600 focus:ring-blue-500" />
                <img :src="`https://api.dicebear.com/7.x/initials/svg?seed=${p.name || p.email}&backgroundColor=1e3a8a&textColor=ffffff`" class="w-10 h-10 rounded-full" />
                <div class="flex-1">
                  <div class="font-bold text-slate-800 text-sm">{{ p.name || '—' }}</div>
                  <div class="text-xs text-slate-500">{{ p.email }}</div>
                </div>
                <Badge variant="outline" class="text-xs">{{ p.applied_position || 'Umum' }}</Badge>
              </label>
            </div>
          </div>
          
          <div class="p-5 border-t bg-white flex justify-end gap-3 shrink-0">
            <div class="flex-1 flex items-center text-sm text-slate-500 font-medium">
              {{ selectedParticipantIds.length }} peserta dipilih
            </div>
            <Button variant="outline" class="rounded-xl" @click="showInviteModal = false">Batal</Button>
            <Button class="bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold" @click="generateAndInvite" :disabled="inviting || selectedParticipantIds.length === 0">
              {{ inviting ? 'Memproses...' : 'Kirim Undangan' }}
            </Button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Monitoring Gallery Modal -->
    <MonitoringGalleryModal 
      :is-open="isMonitoringOpen" 
      :participant-id="selectedParticipant?.id" 
      :participant-name="selectedParticipant?.user_name"
      :session-id="sessionId"
      @close="isMonitoringOpen = false"
    />

  </div>
  <div v-else class="text-center py-20 text-muted-foreground">Memuat detail sesi...</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { 
  ArrowLeftIcon, MonitorPlayIcon, FileTextIcon, ShieldAlertIcon, BarChart3Icon, 
  MailIcon, CameraIcon, RefreshCwIcon 
} from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import MonitoringGalleryModal from '@/components/sessions/MonitoringGalleryModal.vue'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const session = ref<any>(null)
const tokens = ref<any[]>([])
const sessionParticipants = ref<any[]>([])
const loadingParticipants = ref(false)
const alertMessage = ref('')
const alertVariant = ref<'default' | 'destructive'>('default')

// Monitoring State
const isMonitoringOpen = ref(false)
const selectedParticipant = ref<any>(null)

const showSuccess = (message: string) => {
  alertVariant.value = 'default'
  alertMessage.value = message
}

const showError = (message: string) => {
  alertVariant.value = 'destructive'
  alertMessage.value = message
}

const loadData = async () => {
  try {
    const resS = await client.get(`/sessions/${sessionId}`)
    session.value = resS.data
    
    const resT = await client.get(`/sessions/${sessionId}/tokens`)
    tokens.value = resT.data || []

    await loadParticipants()
  } catch (e) {
    showError('Gagal memuat detail sesi.')
  }
}

const loadParticipants = async () => {
  loadingParticipants.value = true
  try {
    const res = await client.get(`/sessions/${sessionId}/participants`)
    sessionParticipants.value = res.data || []
  } catch (e) {
    console.error('Failed to load participants', e)
  } finally {
    loadingParticipants.value = false
  }
}

onMounted(loadData)

const copy = (txt: string) => {
  navigator.clipboard.writeText(txt)
  showSuccess('Token tersalin ke clipboard')
}

// ─── Invite Logic ───────────────────────────────────────────────
const showInviteModal = ref(false)
const allParticipants = ref<any[]>([])
const loadingParticipantsList = ref(false)
const selectedParticipantIds = ref<string[]>([])
const inviting = ref(false)

const openInviteModal = async () => {
  showInviteModal.value = true
  selectedParticipantIds.value = []
  loadingParticipantsList.value = true
  try {
    const res = await client.get('/participants')
    allParticipants.value = res.data || []
  } catch (err) {
    showError('Gagal memuat list peserta')
  } finally {
    loadingParticipantsList.value = false
  }
}

const generateAndInvite = async () => {
  if (selectedParticipantIds.value.length === 0) return
  inviting.value = true
  
  const selectedUsers = allParticipants.value.filter(p => selectedParticipantIds.value.includes(p.id))
  let successCount = 0
  let errorCount = 0
  
  const end_time = new Date(new Date(session.value.schedule).getTime() + session.value.duration_minutes * 60000).toISOString()
  
  for (const user of selectedUsers) {
    try {
      await client.post(`/sessions/${sessionId}/tokens`, {
        count: 1,
        max_usage: 1,
        expires_at: end_time,
        bound_email: user.email
      })
      successCount++
    } catch (err) {
      errorCount++
    }
  }
  
  inviting.value = false
  showInviteModal.value = false
  
  if (successCount > 0) {
    showSuccess(`Berhasil mengundang ${successCount} peserta via email.`)
    await loadData() // Refresh token table and participant list
  }
  if (errorCount > 0) {
    setTimeout(() => showError(`Gagal mengundang ${errorCount} peserta.`), 2000)
  }
}

const openMonitoring = (p: any) => {
  selectedParticipant.value = p
  isMonitoringOpen.value = true
}
</script>
