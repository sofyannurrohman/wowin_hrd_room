<template>
  <div class="space-y-6 max-w-[1200px] mx-auto pb-10">
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div>
        <h2 class="text-3xl font-extrabold tracking-tight text-slate-800">Sessions Management</h2>
        <p class="text-slate-500 mt-1">Manage and monitor sesi ujian.</p>
      </div>
      <div class="flex items-center gap-3 w-full md:w-auto">
        <div class="relative w-full md:w-[280px]">
          <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <Input v-model="searchQuery" placeholder="Search sessions..." class="pl-9 h-10 border-slate-200 focus-visible:ring-blue-500 bg-white" />
        </div>
        <Button class="h-10 px-4 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-semibold flex items-center gap-2 whitespace-nowrap" @click="router.push('/sessions/create')">
          <PlusIcon class="w-4 h-4" /> Buat Sesi
        </Button>
      </div>
    </div>

    <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden">
      <div class="p-4 border-b border-slate-100 bg-slate-50/50 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="flex flex-wrap items-center gap-4 w-full sm:w-auto">
          <div class="flex items-center gap-2">
            <FilterIcon class="w-4 h-4 text-slate-400" />
            <span class="text-sm font-semibold text-slate-600">Filter by:</span>
          </div>
          <Select v-model="statusFilter">
            <SelectTrigger class="w-[160px] h-9 bg-white border-slate-200 text-sm">
              <SelectValue placeholder="All Statuses" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="scheduled">Scheduled</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <Button variant="ghost" size="sm" class="text-slate-500 hover:text-slate-800" @click="clearFilters">
          Clear Filters
        </Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow class="border-b border-slate-100 hover:bg-transparent">
            <TableHead class="text-xs font-bold text-slate-500 uppercase tracking-wider py-4 px-6 cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('name')">
              <div class="flex items-center gap-1">
                Session Name
                <ArrowUpDownIcon v-if="sortConfig.key !== 'name'" class="w-3 h-3" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
              </div>
            </TableHead>
            <TableHead class="text-xs font-bold text-slate-500 uppercase tracking-wider py-4 px-4 cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('schedule')">
              <div class="flex items-center gap-1">
                Start Date
                <ArrowUpDownIcon v-if="sortConfig.key !== 'schedule'" class="w-3 h-3" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
              </div>
            </TableHead>
            <TableHead class="text-xs font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap">End Date</TableHead>
            <TableHead class="text-xs font-bold text-slate-500 uppercase tracking-wider py-4 px-4 text-center cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('total_participant_count')">
              <div class="flex flex-col items-center gap-0.5">
                Participants
                <div class="flex items-center gap-1">
                  <ArrowUpDownIcon v-if="sortConfig.key !== 'total_participant_count'" class="w-3 h-3" />
                  <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                  <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
                </div>
              </div>
            </TableHead>
            <TableHead class="text-xs font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Status</TableHead>
            <TableHead class="text-xs font-bold text-slate-500 uppercase tracking-wider py-4 px-6 text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="s in paginatedData" :key="s.id" class="border-b border-slate-100 transition-colors hover:bg-slate-50/50 group">
            <TableCell class="py-4 px-6">
              <div class="font-bold text-slate-800 text-sm">{{ s.name }}</div>
              <div class="text-[12px] text-slate-400 mt-0.5 font-mono">{{ s.id.slice(0, 8).toUpperCase() }}</div>
            </TableCell>
            <TableCell class="py-4 px-4 text-sm">
              <div class="text-slate-800 font-medium">{{ formatDateMonthDay(s.schedule) }}</div>
              <div class="text-slate-500 text-xs mt-0.5">{{ formatTime(s.schedule) }}</div>
            </TableCell>
            <TableCell class="py-4 px-4 text-sm">
              <div class="text-slate-800 font-medium">{{ formatDateMonthDay(getEndDate(s)) }}</div>
              <div class="text-slate-500 text-xs mt-0.5">{{ formatTime(getEndDate(s)) }}</div>
            </TableCell>
            <TableCell class="py-4 px-4 text-center">
              <div class="inline-flex flex-col items-center">
                <div class="inline-flex items-center justify-center bg-blue-50 text-blue-600 font-bold text-xs h-6 px-3 rounded-full">
                  {{ s.submitted_participant_count || 0 }} / {{ s.total_participant_count || 0 }}
                </div>
                <span class="text-[10px] text-slate-400 mt-0.5">selesai / total</span>
              </div>
            </TableCell>
            <TableCell class="py-4 px-4">
              <Badge :class="getStatusBadgeClass(getStatusText(s))" class="font-medium text-[11px] px-2.5 py-1 rounded-full border-0 inline-flex items-center gap-1.5">
                <div class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(getStatusText(s))"></div>
                {{ getStatusText(s) }}
              </Badge>
            </TableCell>
            <TableCell class="py-4 px-6 text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon" class="text-slate-400 hover:text-blue-600 hover:bg-blue-50 h-8 w-8 rounded-full" @click="router.push(`/sessions/${s.id}`)">
                  <EyeIcon class="w-4 h-4" />
                </Button>
                <Button variant="ghost" size="icon" class="text-slate-400 hover:text-blue-600 hover:bg-blue-50 h-8 w-8 rounded-full" @click="openEditModal(s)">
                  <Edit2Icon class="w-4 h-4" />
                </Button>
                <Button variant="ghost" size="icon" class="text-slate-400 hover:text-red-600 hover:bg-red-50 h-8 w-8 rounded-full" @click="confirmDelete(s)">
                  <Trash2Icon class="w-4 h-4" />
                </Button>
              </div>
            </TableCell>
          </TableRow>

          <TableRow v-if="paginatedData.length === 0">
            <TableCell colspan="6" class="text-center h-48 text-slate-500">
              <div class="flex flex-col items-center justify-center space-y-3">
                <div class="w-12 h-12 rounded-full bg-slate-100 flex items-center justify-center text-slate-400">
                  <SearchIcon class="w-6 h-6" />
                </div>
                <p>Tidak ada sesi yang ditemukan.</p>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <div class="p-4 border-t border-slate-100 bg-slate-50/30">
        <DataTablePagination 
          v-if="filteredSessions.length > 0"
          :total="totalItems"
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
        />
      </div>
    </Card>

    <!-- ─── Edit Modal ─────────────────────────────────────────────── -->
    <Teleport to="body">
      <div v-if="editModal.open" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="editModal.open = false">
        <div class="bg-white rounded-2xl shadow-xl p-8 max-w-lg w-full mx-4 space-y-5">
          <div class="flex items-center justify-between">
            <h2 class="text-xl font-bold text-slate-900">Edit Sesi</h2>
            <button @click="editModal.open = false" class="text-slate-400 hover:text-slate-700">
              <XIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="space-y-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-slate-700">Nama Sesi</label>
              <Input v-model="editForm.name" placeholder="Nama sesi" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-slate-700">Jadwal</label>
                <Input v-model="editForm.schedule" type="datetime-local" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-slate-700">Durasi (menit)</label>
                <Input v-model.number="editForm.duration_minutes" type="number" min="1" />
              </div>
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-slate-700">Maks Peserta</label>
              <Input v-model.number="editForm.max_participants" type="number" min="1" />
            </div>
            
            <div class="space-y-2 pt-2 border-t">
              <label class="text-sm font-bold text-slate-800">Modul Soal (Wajib)</label>
              <div class="grid grid-cols-1 gap-2 max-h-[160px] overflow-y-auto pr-2">
                <label v-for="m in allModules" :key="m.id" class="flex items-center gap-3 p-2 border rounded-lg hover:bg-slate-50 cursor-pointer transition-colors">
                  <input type="checkbox" :value="m.id" v-model="selectedModuleIds" class="w-4 h-4 rounded text-blue-600" />
                  <div class="flex-1">
                    <p class="text-sm font-medium text-slate-800">{{ m.name }}</p>
                    <p class="text-[10px] text-slate-500 line-clamp-1">{{ m.description || 'Tidak ada deskripsi' }}</p>
                  </div>
                </label>
              </div>
              <p class="text-[10px] text-slate-400">Pilih modul yang akan dikerjakan oleh peserta dalam sesi ini.</p>
            </div>

            <div class="flex items-center gap-6">
              <label class="flex items-center gap-2 text-sm text-slate-700 cursor-pointer">
                <input type="checkbox" v-model="editForm.randomize_questions" class="rounded" />
                Acak Soal
              </label>
              <label class="flex items-center gap-2 text-sm text-slate-700 cursor-pointer">
                <input type="checkbox" v-model="editForm.show_score" class="rounded" />
                Tampilkan Skor
              </label>
            </div>
          </div>
          <p v-if="editError" class="text-sm text-red-500">{{ editError }}</p>
          <div class="flex justify-end gap-3 pt-2">
            <Button variant="outline" @click="editModal.open = false">Batal</Button>
            <Button :disabled="editModal.saving" @click="saveEdit" class="bg-blue-600 hover:bg-blue-700 text-white">
              {{ editModal.saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
            </Button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ─── Delete Confirmation Dialog ────────────────────────────── -->
    <Teleport to="body">
      <div v-if="deleteDialog.open" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="deleteDialog.open = false">
        <div class="bg-white rounded-2xl shadow-xl p-8 max-w-md w-full mx-4">
          <div class="flex items-start gap-4 mb-6">
            <div class="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center flex-shrink-0 mt-0.5">
              <Trash2Icon class="w-5 h-5 text-red-600" />
            </div>
            <div>
              <h2 class="text-lg font-bold text-slate-900">Hapus Sesi</h2>
              <p class="text-sm text-slate-500 mt-1">
                Apakah Anda yakin ingin menghapus sesi <strong class="text-slate-800">{{ deleteDialog.session?.name }}</strong>?
                Semua data terkait (token, jawaban, pelanggaran) akan ikut terhapus dan tidak bisa dikembalikan.
              </p>
            </div>
          </div>
          <div class="flex justify-end gap-3">
            <Button variant="outline" @click="deleteDialog.open = false">Batal</Button>
            <Button :disabled="deleteDialog.deleting" @click="doDelete" class="bg-red-600 hover:bg-red-700 text-white">
              {{ deleteDialog.deleting ? 'Menghapus...' : 'Ya, Hapus' }}
            </Button>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import client from '@/api/client'
import { toast } from 'vue-sonner'
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'

import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { SearchIcon, PlusIcon, FilterIcon, EyeIcon, Edit2Icon, Trash2Icon, XIcon, ArrowUpDownIcon, ArrowUpIcon, ArrowDownIcon } from 'lucide-vue-next'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

const router = useRouter()
const sessionStore = useSessionStore()

const searchQuery = ref('')
const statusFilter = ref('all')

const filteredSessions = computed(() => {
  return sessionStore.sessions.filter((s: any) => {
    const matchesSearch = s.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    if (!matchesSearch) return false
    if (statusFilter.value !== 'all') {
      const status = getStatusText(s).toLowerCase()
      if (status !== statusFilter.value) return false
    }
    return true
  })
})

const allModules = ref<any[]>([])
const selectedModuleIds = ref<string[]>([])

onMounted(async () => {
  await sessionStore.fetchSessions()
  try {
    const res = await client.get('/modules')
    allModules.value = res.data.modules || []
  } catch(e) {
    console.error('Gagal memuat modul:', e)
  }
})

const {
  pageSize,
  currentPage,
  sortConfig,
  toggleSort,
  paginatedData,
  totalItems
} = useDataTable(filteredSessions)

const clearFilters = () => {
  searchQuery.value = ''
  statusFilter.value = 'all'
}

// ─── Edit Modal ───────────────────────────────────────────────────
const editModal = reactive({ open: false, saving: false, sessionId: '' })
const editError = ref('')
const editForm = reactive({
  name: '',
  schedule: '',
  duration_minutes: 60,
  max_participants: 50,
  randomize_questions: false,
  show_score: true,
})

const toLocalDatetime = (iso: string) => {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const openEditModal = async (session: any) => {
  editModal.sessionId = session.id
  editForm.name = session.name
  editForm.schedule = toLocalDatetime(session.schedule)
  editForm.duration_minutes = session.duration_minutes
  editForm.max_participants = session.max_participants
  editForm.randomize_questions = session.randomize_questions
  editForm.show_score = session.show_score
  editError.value = ''
  
  // Fetch assigned modules
  selectedModuleIds.value = []
  try {
    const res = await client.get(`/sessions/${session.id}/modules`)
    const assigned = res.data || []
    selectedModuleIds.value = assigned.map((m: any) => m.id)
  } catch(e) {
    toast.error('Gagal mengambil data modul sesi.')
  }
  
  editModal.open = true
}

const saveEdit = async () => {
  if (!editForm.name.trim()) { editError.value = 'Nama sesi tidak boleh kosong.'; return }
  if (selectedModuleIds.value.length === 0) {
    editError.value = 'Minimal satu modul soal harus dipilih.'
    return
  }

  editModal.saving = true
  editError.value = ''
  try {
    await sessionStore.updateSession(editModal.sessionId, {
      name: editForm.name,
      schedule: new Date(editForm.schedule).toISOString(),
      duration_minutes: editForm.duration_minutes,
      max_participants: editForm.max_participants,
      randomize_questions: editForm.randomize_questions,
      show_score: editForm.show_score,
      module_ids: selectedModuleIds.value,
    })
    editModal.open = false
    toast.success('Sesi berhasil diperbarui.')
  } catch (e: any) {
    editError.value = e?.response?.data?.error || 'Gagal memperbarui sesi.'
  } finally {
    editModal.saving = false
  }
}

// ─── Delete Dialog ────────────────────────────────────────────────
const deleteDialog = reactive({ open: false, deleting: false, session: null as any })

const confirmDelete = (session: any) => {
  deleteDialog.session = session
  deleteDialog.open = true
}

const doDelete = async () => {
  if (!deleteDialog.session) return
  deleteDialog.deleting = true
  try {
    await sessionStore.deleteSession(deleteDialog.session.id)
    deleteDialog.open = false
    toast.success('Sesi berhasil dihapus.')
  } catch (e: any) {
    toast.error(e?.response?.data?.error || 'Gagal menghapus sesi.')
  } finally {
    deleteDialog.deleting = false
  }
}

// ─── Helpers ──────────────────────────────────────────────────────
const getEndDate = (session: any) => {
  const start = new Date(session.schedule)
  return new Date(start.getTime() + session.duration_minutes * 60000).toISOString()
}

const formatDateMonthDay = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
}

const getStatusText = (session: any) => {
  if (session.is_locked) return 'Locked'
  const now = new Date()
  const start = new Date(session.schedule)
  const end = new Date(start.getTime() + session.duration_minutes * 60000)
  if (now < start) return 'Scheduled'
  if (now > end) return 'Completed'
  return 'Active'
}

const getStatusBadgeClass = (status: string) => {
  if (status === 'Active') return 'bg-green-100 text-green-700 hover:bg-green-100'
  if (status === 'Scheduled') return 'bg-blue-100 text-blue-700 hover:bg-blue-100'
  if (status === 'Completed') return 'bg-slate-100 text-slate-700 hover:bg-slate-100'
  return 'bg-red-100 text-red-700 hover:bg-red-100'
}

const getStatusDotClass = (status: string) => {
  if (status === 'Active') return 'bg-green-500'
  if (status === 'Scheduled') return 'bg-blue-500'
  if (status === 'Completed') return 'bg-slate-500'
  return 'bg-red-500'
}
</script>
