<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Bank Soal</h2>
      <Button @click="router.push('/questions/create')">
        <PlusIcon class="mr-2 h-4 w-4" /> Tambah Soal
      </Button>
    </div>

    <Alert v-if="alertMessage" :variant="alertVariant">
      <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
      <AlertDescription>{{ alertMessage }}</AlertDescription>
    </Alert>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[400px]">Konten Soal</TableHead>
            <TableHead>Tipe Soal</TableHead>
            <TableHead>Lampiran Gambar</TableHead>
            <TableHead>Modul (Opsional)</TableHead>
            <TableHead class="text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="q in questions" :key="q.id">
            <TableCell class="font-medium line-clamp-2">{{ q.content }}</TableCell>
            <TableCell>
              <Badge variant="outline">{{ q.type.replace('_', ' ').toUpperCase() }}</Badge>
            </TableCell>
            <TableCell>
              <div v-if="q.image_url" class="relative w-16 h-12 rounded overflow-hidden border border-slate-200">
                <img :src="q.image_url" class="absolute inset-0 w-full h-full object-cover" alt="Lampiran" />
              </div>
              <span v-else class="text-muted-foreground">-</span>
            </TableCell>
            <TableCell>{{ getModuleName(q.module_id) }}</TableCell>
            <TableCell class="text-right flex justify-end gap-1">
              <Button variant="ghost" size="icon" @click="router.push(`/questions/edit/${q.id}`)">
                <EditIcon class="w-4 h-4 text-primary" />
              </Button>
              <Button variant="ghost" size="icon" @click="deleteQuestion(q.id)">
                <TrashIcon class="w-4 h-4 text-destructive" />
              </Button>
            </TableCell>
          </TableRow>
          <TableRow v-if="questions.length === 0">
             <TableCell colspan="5" class="text-center h-24">Belum ada soal dibuat.</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import client from '@/api/client'

import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { PlusIcon, EditIcon, TrashIcon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const router = useRouter()
const questions = ref<any[]>([])
const modules = ref<any[]>([])
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

const loadQuestions = async () => {
  try {
    const [resQ, resM] = await Promise.all([
      client.get('/questions'),
      client.get('/modules')
    ])
    questions.value = resQ.data || []
    modules.value = resM.data.modules || []
  } catch(e) {
    showError('Gagal memuat bank soal dan modul')
  }
}

const getModuleName = (moduleId: string) => {
  if (!moduleId) return 'Global Bank Soal'
  const mod = modules.value.find(m => m.id === moduleId)
  return mod ? mod.name : 'Modul Tidak Diketahui'
}

onMounted(loadQuestions)

const deleteQuestion = async (id: string) => {
  if (!confirm('Hapus soal ini dari bank soal?')) return
  try {
    await client.delete(`/questions/${id}`)
    showSuccess('Soal berhasil dihapus')
    await loadQuestions()
  } catch (e: any) {
    showError(e.response?.data?.error || 'Gagal menghapus soal')
  }
}
</script>
