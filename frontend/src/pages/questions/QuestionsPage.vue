<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Bank Soal</h2>
      <div class="flex gap-2">
        <Button variant="outline" @click="isImportModalOpen = true">
          <UploadIcon class="mr-2 h-4 w-4" /> Import CSV
        </Button>
        <Button @click="router.push('/questions/create')">
          <PlusIcon class="mr-2 h-4 w-4" /> Tambah Soal
        </Button>
      </div>
    </div>

    <Alert v-if="alertMessage" :variant="alertVariant">
      <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
      <AlertDescription>{{ alertMessage }}</AlertDescription>
    </Alert>

    <Dialog v-model:open="isImportModalOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Import Bank Soal (CSV)</DialogTitle>
          <DialogDescription>
            Unggah file CSV dengan format yang sesuai. Anda dapat mendownload template terlebih dahulu.
          </DialogDescription>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <Button variant="outline" as="a" href="/template-soal.csv" download>
            <DownloadIcon class="mr-2 h-4 w-4" /> Download Template CSV
          </Button>
          <div class="flex flex-col gap-2">
            <label for="csvFile" class="text-sm font-medium">Pilih File CSV</label>
            <input
              id="csvFile"
              type="file"
              accept=".csv"
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              @change="handleFileChange"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="isImportModalOpen = false" :disabled="isImporting">Batal</Button>
          <Button @click="submitImport" :disabled="!importFile || isImporting">
            {{ isImporting ? 'Mengimport...' : 'Import' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

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
import { PlusIcon, EditIcon, TrashIcon, UploadIcon, DownloadIcon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const router = useRouter()
const questions = ref<any[]>([])
const modules = ref<any[]>([])
const alertMessage = ref('')
const alertVariant = ref<'default' | 'destructive'>('default')

const isImportModalOpen = ref(false)
const importFile = ref<File | null>(null)
const isImporting = ref(false)

const showSuccess = (message: string) => {
  alertVariant.value = 'default'
  alertMessage.value = message
}

const showError = (message: string) => {
  alertVariant.value = 'destructive'
  alertMessage.value = message
}

const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    importFile.value = target.files[0]
  }
}

const submitImport = async () => {
  if (!importFile.value) return

  isImporting.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)

    const res = await client.post('/questions/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    showSuccess(res.data.message || 'Import berhasil')
    isImportModalOpen.value = false
    importFile.value = null
    await loadQuestions()
  } catch (e: any) {
    showError(e.response?.data?.error || 'Gagal mengimport soal')
  } finally {
    isImporting.value = false
  }
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
