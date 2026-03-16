<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">System Logs (Audit)</h2>
      <div class="flex items-center gap-3">
        <div class="relative w-64">
           <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
           <Input v-model="search" placeholder="Cari aksi atau detail..." class="pl-9" />
        </div>
        <Button variant="outline" @click="loadLogs">Refresh</Button>
      </div>
    </div>



    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[180px] cursor-pointer hover:text-primary transition-colors" @click="toggleSort('created_at')">
              <div class="flex items-center gap-1">
                Waktu
                <ArrowUpDownIcon v-if="sortConfig.key !== 'created_at'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
            <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('action')">
              <div class="flex items-center gap-1">
                Aksi
                <ArrowUpDownIcon v-if="sortConfig.key !== 'action'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
            <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('detail')">
              <div class="flex items-center gap-1">
                Detail
                <ArrowUpDownIcon v-if="sortConfig.key !== 'detail'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
            <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('user_name')">
              <div class="flex items-center gap-1">
                Pengguna
                <ArrowUpDownIcon v-if="sortConfig.key !== 'user_name'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="log in paginatedData" :key="log.id" class="text-sm">
            <TableCell class="text-muted-foreground whitespace-nowrap">{{ new Date(log.created_at).toLocaleString() }}</TableCell>
            <TableCell>
              <Badge variant="outline">{{ log.action.toUpperCase() }}</Badge>
            </TableCell>
            <TableCell class="font-medium">{{ log.detail }}</TableCell>
            <TableCell class="font-mono text-xs">{{ log.user_name || log.user_id || 'System' }}</TableCell>
          </TableRow>
          <TableRow v-if="filteredLogs.length === 0">
             <TableCell colspan="5" class="text-center h-24 text-muted-foreground">{{ logs.length === 0 ? 'Memuat log sistem...' : 'Tidak ada log ditemukan.' }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      
      <div v-if="filteredLogs.length > 0" class="p-4 border-t border-slate-100 bg-slate-50/30">
        <DataTablePagination 
          :total="totalItems"
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
        />
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import client from '@/api/client'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ArrowUpDownIcon, ArrowUpIcon, ArrowDownIcon, SearchIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { Input } from '@/components/ui/input'
import { computed } from 'vue'

const logs = ref<any[]>([])
const search = ref('')

const filteredLogs = computed(() => {
  const q = search.value.toLowerCase()
  return logs.value.filter(l => 
    l.action.toLowerCase().includes(q) || 
    (l.detail || '').toLowerCase().includes(q) ||
    (l.user_name || '').toLowerCase().includes(q)
  )
})

const {
  pageSize,
  currentPage,
  sortConfig,
  toggleSort,
  paginatedData,
  totalItems
} = useDataTable(filteredLogs)

const showError = (message: string) => {
  toast.error(message)
}

const loadLogs = async () => {
  try {
    const res = await client.get('/admin/logs?limit=100')
    logs.value = res.data || []
  } catch(e) {
    showError('Gagal memuat log sistem')
  }
}

onMounted(loadLogs)
</script>
