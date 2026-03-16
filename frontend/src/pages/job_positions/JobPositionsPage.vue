<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900">Posisi Pekerjaan</h1>
        <p class="text-sm text-slate-500 mt-1">Kelola daftar posisi yang sedang aktif dibuka</p>
      </div>
      <div class="flex items-center gap-3 w-full md:w-auto">
        <div class="relative w-64">
           <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
           <input
             v-model="search"
             placeholder="Cari nama posisi..."
             class="w-full pl-9 pr-4 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
           />
        </div>
        <RouterLink
          to="/job-positions/create"
          class="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium transition-colors bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          <PlusIcon class="w-4 h-4" /> Tambah Posisi
        </RouterLink>
      </div>
    </div>

    <!-- Job Positions List -->
    <div class="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm">
      <div v-if="loading" class="p-8 text-center text-slate-500">
        Memuat data...
      </div>
      
      <div v-else-if="filteredPositions.length === 0" class="p-12 text-center flex flex-col items-center">
        <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mb-4 text-slate-400">
          <BriefcaseIcon class="w-8 h-8" />
        </div>
        <h3 class="text-lg font-medium text-slate-900 mb-1">Belum Ada Posisi</h3>
        <p class="text-slate-500 max-w-sm">Tambahkan posisi pekerjaan pertama Anda untuk ditampilkan di halaman Join peserta.</p>
        <RouterLink to="/job-positions/create" class="mt-6 text-blue-600 font-medium hover:underline text-sm">
          Tambah Posisi Baru
        </RouterLink>
      </div>

      <div v-else class="divide-y divide-slate-100">
        <div
          v-for="pos in paginatedData"
          :key="pos.id"
          class="p-6 transition-colors hover:bg-slate-50 flex items-start gap-4 justify-between"
        >
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 rounded bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
              <BriefcaseIcon class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-slate-900 truncate pr-4 text-lg">
                {{ pos.name }}
              </h3>
              <p v-if="pos.description" class="text-xs text-slate-400 mt-1 line-clamp-1 italic">
                {{ pos.description }}
              </p>
              <div class="flex items-center gap-4 text-xs font-medium text-slate-500 mt-2">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-medium" :class="pos.is_active ? 'bg-green-100 text-green-700' : 'bg-slate-100 text-slate-600'">
                  {{ pos.is_active ? 'Aktif' : 'Non-Aktif' }}
                </span>
                <div class="flex items-center gap-1.5 bg-slate-100 px-2.5 py-1 rounded-md">
                  <CalendarIcon class="w-3.5 h-3.5 text-slate-400" />
                  {{ new Date(pos.created_at).toLocaleDateString() }}
                </div>
              </div>
            </div>
          </div>
          
          <div class="flex gap-2 shrink-0">
            <RouterLink
              :to="`/job-positions/edit/${pos.id}`"
              class="text-slate-400 hover:text-blue-600 transition-colors p-2 rounded hover:bg-blue-50"
            >
              <Edit2Icon class="w-4 h-4" />
            </RouterLink>
            <button
              @click="deletePosition(pos.id)"
              class="text-slate-400 hover:text-red-600 transition-colors p-2 rounded hover:bg-red-50"
            >
              <Trash2Icon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
      
      <div v-if="filteredPositions.length > 0" class="p-4 border-t border-slate-100 bg-slate-50/30">
        <DataTablePagination 
          :total="totalItems"
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { PlusIcon, BriefcaseIcon, Edit2Icon, Trash2Icon, CalendarIcon, SearchIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import client from '@/api/client'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'

interface JobPosition {
  id: string
  name: string
  description: string
  requirements: string
  is_active: boolean
  created_at: string
}

const positions = ref<JobPosition[]>([])
const loading = ref(true)
const search = ref('')

const filteredPositions = computed(() => {
  return positions.value.filter(p => p.name.toLowerCase().includes(search.value.toLowerCase()))
})

const {
  pageSize,
  currentPage,
  paginatedData,
  totalItems
} = useDataTable(filteredPositions)

const fetchPositions = async () => {
  try {
    const res = await client.get('/job-positions')
    positions.value = res.data.positions || []
  } catch (error) {
    console.error('Failed to fetch positions:', error)
  } finally {
    loading.value = false
  }
}

const deletePosition = async (id: string) => {
  if (!confirm('Hapus posisi ini secara permanen?')) return
  try {
    await client.delete(`/job-positions/${id}`)
    positions.value = positions.value.filter((p: JobPosition) => p.id !== id)
    toast.success('Posisi berhasil dihapus')
  } catch (error: any) {
    toast.error(error.response?.data?.error || 'Gagal menghapus posisi')
  }
}

onMounted(() => {
  fetchPositions()
})
</script>
