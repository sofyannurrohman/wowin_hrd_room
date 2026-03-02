<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <h2 class="text-3xl font-bold tracking-tight">Hasil Ujian & Evaluasi</h2>
      </div>
      <div class="flex gap-2">
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Daftar Hasil Peserta</CardTitle>
        <CardDescription>Review hasil grading otomatis dan berikan nilai manual untuk esai.</CardDescription>
      </CardHeader>
      <CardContent>
        <Alert v-if="alertMessage" :variant="alertVariant" class="mb-4">
          <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
          <AlertDescription>{{ alertMessage }}</AlertDescription>
        </Alert>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Peserta</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Status Ujian</TableHead>
              <TableHead>Skor</TableHead>
              <TableHead>Status Review</TableHead>
              <TableHead class="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="res in results" :key="res.id">
              <TableCell class="font-medium">{{ res.participant?.user?.name || res.participant_id }}</TableCell>
              <TableCell>{{ res.participant?.user?.email || '-' }}</TableCell>
              <TableCell>Disubmit</TableCell>
              <TableCell class="font-bold text-primary">{{ res.total_score }}</TableCell>
              <TableCell>
                <Badge :variant="res.status === 'Finalized' ? 'default' : 'secondary'">{{ res.status }}</Badge>
              </TableCell>
              <TableCell class="text-right">
                 <Button size="sm" @click="submitFinalize(res)" :disabled="res.status === 'completed'" v-if="res.status !== 'completed'">Finalisasi</Button>
              </TableCell>
            </TableRow>
            <TableRow v-if="results.length === 0">
              <TableCell colspan="6" class="text-center h-24">Belum ada hasil yang masuk. Pastikan peserta telah Submit ujian.</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowLeftIcon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const results = ref<any[]>([])
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

const loadResults = async () => {
  try {
    const res = await client.get(`/sessions/${sessionId}/results`)
    results.value = res.data || []
  } catch (e) {
    showError('Gagal memuat hasil ujian')
  }
}

onMounted(loadResults)

const submitFinalize = async (res: any) => {
  try {
    await client.post(`/results/${res.id}/finalize`)
    showSuccess('Nilai berhasil di-finalisasi')
    await loadResults()
  } catch(e) {
    showError('Gagal finalisasi nilai')
  }
}
</script>
