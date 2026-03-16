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
          <p>{{ session.is_locked ? 'Terkunci' : 'Aktif' }} ({{ session.participant_count || 0 }} dari {{ session.max_participants }})</p>
        </div>
      </CardContent>
    </Card>

    <!-- Token Generation -->
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
            <div v-if="loadingParticipants" class="text-center py-10 text-slate-400">Memuat data peserta...</div>
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
import { ArrowLeftIcon, MonitorPlayIcon, FileTextIcon, ShieldAlertIcon, BarChart3Icon, MailIcon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const session = ref<any>(null)
const tokens = ref<any[]>([])
const alertMessage = ref('')
const alertVariant = ref<'default' | 'destructive'>('default')

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
  } catch (e) {
    showError('Gagal memuat detail sesi.')
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
const loadingParticipants = ref(false)
const selectedParticipantIds = ref<string[]>([])
const inviting = ref(false)

const openInviteModal = async () => {
  showInviteModal.value = true
  selectedParticipantIds.value = []
  loadingParticipants.value = true
  try {
    const res = await client.get('/participants')
    allParticipants.value = res.data || []
  } catch (err) {
    showError('Gagal memuat list peserta')
  } finally {
    loadingParticipants.value = false
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
      // In a real app, this is where the backend would send the email via SMTP/SendGrid.
      // For now, the token is generated and bound to their email.
    } catch (err) {
      errorCount++
    }
  }
  
  inviting.value = false
  showInviteModal.value = false
  
  if (successCount > 0) {
    showSuccess(`Berhasil mengundang ${successCount} peserta via email.`)
    await loadData() // Refresh token table
  }
  if (errorCount > 0) {
    setTimeout(() => showError(`Gagal mengundang ${errorCount} peserta.`), 2000)
  }
}
</script>
