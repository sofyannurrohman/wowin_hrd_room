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

        <div class="flex items-center gap-4 pb-4 border-b">
           <Input v-model.number="tokenCount" type="number" min="1" max="50" class="w-32" />
           <Button @click="generateTokens" :disabled="generating">Generate Token Baru</Button>
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
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { ArrowLeftIcon, MonitorPlayIcon, FileTextIcon, ShieldAlertIcon, BarChart3Icon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const session = ref<any>(null)
const tokens = ref<any[]>([])
const tokenCount = ref(1)
const generating = ref(false)
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

const generateTokens = async () => {
  generating.value = true
  try {
    const end_time = new Date(new Date(session.value.schedule).getTime() + session.value.duration_minutes * 60000).toISOString()
    const res = await client.post(`/sessions/${sessionId}/tokens`, {
      count: tokenCount.value,
      max_usage: 1, // typically 1 participant per token
      expires_at: end_time
    })
    
    showSuccess(`${tokenCount.value} Token berhasil dibuat`)
    
    // Merge new tokens with plainly visible token strings into the list
    if (res.data && res.data.tokens) {
      tokens.value = [...res.data.tokens, ...tokens.value]
    } else {
      await loadData()
    }
  } catch (e) {
    showError('Gagal generate token')
  } finally {
    generating.value = false
  }
}

const copy = (txt: string) => {
  navigator.clipboard.writeText(txt)
  showSuccess('Token tersalin ke clipboard')
}
</script>
