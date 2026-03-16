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
        <div class="mb-4">
           <div class="relative w-64">
              <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
              <Input v-model="search" placeholder="Cari nama atau email..." class="pl-9" />
           </div>
        </div>


        <Table>
          <TableHeader>
            <TableRow>
              <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('user_name')">
                <div class="flex items-center gap-1">
                  Peserta
                  <ArrowUpDownIcon v-if="sortConfig.key !== 'user_name'" class="w-3.5 h-3.5" />
                  <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                  <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
                </div>
              </TableHead>
              <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('user_email')">
                <div class="flex items-center gap-1">
                  Email
                  <ArrowUpDownIcon v-if="sortConfig.key !== 'user_email'" class="w-3.5 h-3.5" />
                  <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                  <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
                </div>
              </TableHead>
              <TableHead>Status Ujian</TableHead>
              <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('total_score')">
                <div class="flex items-center gap-1">
                  Skor
                  <ArrowUpDownIcon v-if="sortConfig.key !== 'total_score'" class="w-3.5 h-3.5" />
                  <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                  <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
                </div>
              </TableHead>
              <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('grading_status')">
                <div class="flex items-center gap-1">
                  Status Review
                  <ArrowUpDownIcon v-if="sortConfig.key !== 'grading_status'" class="w-3.5 h-3.5" />
                  <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                  <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
                </div>
              </TableHead>
              <TableHead class="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="res in paginatedData" :key="res.id">
              <TableCell class="font-medium">{{ res.user_name || '-' }}</TableCell>
              <TableCell>{{ res.user_email || '-' }}</TableCell>
              <TableCell>Disubmit</TableCell>
              <TableCell class="font-bold text-primary">{{ res.total_score }}</TableCell>
              <TableCell>
                <Badge :variant="res.grading_status === 'completed' ? 'default' : 'secondary'">
                  {{ res.grading_status === 'completed' ? 'Selesai' : 'Menunggu Review' }}
                </Badge>
              </TableCell>
              <TableCell class="text-right flex items-center justify-end gap-2">
                <Button size="sm" variant="outline" @click="router.push(`/sessions/${sessionId}/results/${res.participant_id}/answers`)" v-if="res.grading_status !== 'completed'">Review Jawaban</Button>
                <Button size="sm" variant="ghost" @click="router.push(`/sessions/${sessionId}/results/${res.participant_id}/answers`)" v-else>Lihat Jawaban</Button>
                <Button size="sm" @click="submitFinalize(res)" v-if="res.grading_status !== 'completed'">Finalisasi</Button>
              </TableCell>
            </TableRow>
            <TableRow v-if="filteredResults.length === 0">
              <TableCell colspan="6" class="text-center h-24">{{ results.length === 0 ? 'Belum ada hasil yang masuk. Pastikan peserta telah Submit ujian.' : 'Tidak ada hasil ditemukan.' }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
        
        <div v-if="filteredResults.length > 0" class="mt-4">
          <DataTablePagination 
            :total="totalItems"
            v-model:pageSize="pageSize"
            v-model:currentPage="currentPage"
          />
        </div>
      </CardContent>
    </Card>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowLeftIcon, ArrowUpDownIcon, ArrowUpIcon, ArrowDownIcon, SearchIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { Input } from '@/components/ui/input'
import { computed } from 'vue'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const results = ref<any[]>([])
const search = ref('')

const filteredResults = computed(() => {
  const q = search.value.toLowerCase()
  return results.value.filter(res => 
    (res.user_name || '').toLowerCase().includes(q) || 
    (res.user_email || '').toLowerCase().includes(q)
  )
})

const {
  pageSize,
  currentPage,
  sortConfig,
  toggleSort,
  paginatedData,
  totalItems
} = useDataTable(filteredResults)

const showSuccess = (message: string) => {
  toast.success(message)
}

const showError = (message: string) => {
  toast.error(message)
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
