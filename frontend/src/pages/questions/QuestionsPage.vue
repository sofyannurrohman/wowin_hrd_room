<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Bank Soal</h2>
      <Button @click="router.push('/questions/create')">
        <PlusIcon class="mr-2 h-4 w-4" /> Tambah Soal
      </Button>
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[400px]">Konten Soal</TableHead>
            <TableHead>Tipe Soal</TableHead>
            <TableHead>Lampiran Gambar</TableHead>
            <TableHead>Sesi (Opsional)</TableHead>
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
                <ImageIcon v-if="q.image_url" class="w-4 h-4 text-green-500" />
                <span v-else class="text-slate-400">-</span>
            </TableCell>
            <TableCell>{{ q.session_id ? 'Terkait Sesi Tertentu' : 'Global Bank Soal' }}</TableCell>
            <TableCell class="text-right">
              <Button variant="ghost" size="sm">Edit</Button>
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
import { PlusIcon, ImageIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const router = useRouter()
const questions = ref<any[]>([])

onMounted(async () => {
  try {
    const res = await client.get('/questions')
    questions.value = res.data || []
  } catch(e) {
    toast.error('Gagal memuat bank soal')
  }
})
</script>
