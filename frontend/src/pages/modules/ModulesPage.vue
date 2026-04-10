<template>
  <div class="space-y-6 w-full max-w-[1200px] mx-auto pb-10">
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900">Module Soal</h1>
        <p class="text-sm text-slate-500 mt-1">Manage question banks and modules</p>
      </div>
      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
        <div class="relative w-64">
           <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
           <input
             v-model="search"
             placeholder="Cari nama modul..."
             class="w-full pl-9 pr-4 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
           />
        </div>
        <RouterLink
          to="/modules/create"
          class="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium transition-colors bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          <PlusIcon class="w-4 h-4" /> Buat Module Baru
        </RouterLink>
      </div>
    </div>

    <!-- Modules List -->
    <div class="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm">
      <div v-if="loading" class="p-8 text-center text-slate-500">
        Loading modules...
      </div>
      
      <div v-else-if="filteredModules.length === 0" class="p-12 text-center flex flex-col items-center">
        <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mb-4 text-slate-400">
          <LayersIcon class="w-8 h-8" />
        </div>
        <h3 class="text-lg font-medium text-slate-900 mb-1">Belum ada modul</h3>
        <p class="text-slate-500 max-w-sm">Mulai buat modul soal untuk mengelompokan soal ujian.</p>
        <RouterLink to="/modules/create" class="mt-6 text-blue-600 font-medium hover:underline text-sm">
          Buat Modul Baru
        </RouterLink>
      </div>

      <div v-else class="divide-y divide-slate-100">
        <div
          v-for="mod in paginatedData"
          :key="mod.id"
          class="p-6 transition-colors hover:bg-slate-50 flex items-start gap-4"
        >
          <div class="w-10 h-10 rounded bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
            <LayersIcon class="w-5 h-5" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-1 w-full relative">
              <RouterLink :to="`/modules/${mod.id}`" class="font-semibold text-slate-900 hover:text-blue-600 truncate pr-4 text-lg inline-block">
                {{ mod.name }}
              </RouterLink>
              <div class="absolute right-0 top-0 flex gap-2">
                <RouterLink
                  :to="`/modules/${mod.id}`"
                  class="text-slate-400 hover:text-blue-600 transition-colors p-1"
                >
                  <EyeIcon class="w-4 h-4" />
                </RouterLink>
                <RouterLink
                  :to="`/modules/edit/${mod.id}`"
                  class="text-slate-400 hover:text-blue-600 transition-colors p-1"
                >
                  <Edit2Icon class="w-4 h-4" />
                </RouterLink>
                <button
                  @click="deleteModule(mod.id)"
                  class="text-slate-400 hover:text-red-600 transition-colors p-1"
                >
                  <Trash2Icon class="w-4 h-4" />
                </button>
              </div>
            </div>
            <div class="text-xs text-slate-400 font-mono mb-2">ID: {{ mod.id.substring(0,8) }}</div>
            <p class="text-sm text-slate-500 leading-relaxed mb-4 line-clamp-2">
              {{ mod.description || 'No description provided.' }}
            </p>
            <div class="flex items-center gap-4 text-xs font-medium text-slate-500">
              <div class="flex items-center gap-1.5 bg-slate-100 px-2.5 py-1 rounded-md">
                <CalendarIcon class="w-3.5 h-3.5 text-slate-400" />
                {{ new Date(mod.created_at).toLocaleDateString() }}
              </div>
              <div class="flex items-center gap-1.5 bg-blue-50 text-blue-700 px-2.5 py-1 rounded-md font-bold">
                Total Weight: {{ mod.total_weight }}
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="filteredModules.length > 0" class="p-4 border-t border-slate-100 bg-slate-50/30">
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
import { PlusIcon, LayersIcon, Edit2Icon, Trash2Icon, CalendarIcon, SearchIcon, EyeIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import client from '@/api/client'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'

interface Module {
  id: string
  name: string
  description: string
  total_weight: number
  created_at: string
}

const modules = ref<Module[]>([])
const loading = ref(true)
const search = ref('')

const filteredModules = computed(() => {
  return modules.value.filter(m => m.name.toLowerCase().includes(search.value.toLowerCase()))
})

const {
  pageSize,
  currentPage,
  paginatedData,
  totalItems
} = useDataTable(filteredModules)

const fetchModules = async () => {
  try {
    const res = await client.get('/modules')
    modules.value = res.data.modules || []
  } catch (error) {
    console.error('Failed to parse modules:', error)
  } finally {
    loading.value = false
  }
}

const deleteModule = async (id: string) => {
  if (!confirm('Are you sure you want to delete this module? This cannot be undone.')) return
  try {
    await client.delete(`/modules/${id}`)
    modules.value = modules.value.filter((m: Module) => m.id !== id)
    toast.success('Modul berhasil dihapus')
  } catch (error: any) {
    toast.error(error.response?.data?.error || 'Gagal menghapus modul')
  }
}

onMounted(() => {
  fetchModules()
})
</script>
