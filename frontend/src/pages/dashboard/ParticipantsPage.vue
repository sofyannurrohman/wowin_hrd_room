<template>
  <div class="space-y-6 max-w-[1200px] mx-auto pb-10">

    <!-- Toast Notification -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 -translate-y-3"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-3"
    >
      <div
        v-if="toast.show"
        class="fixed top-5 right-5 z-[9999] flex items-start gap-3 px-5 py-4 rounded-2xl shadow-2xl min-w-[280px] max-w-sm border"
        :class="toast.type === 'success'
          ? 'bg-green-50 border-green-200 text-green-800'
          : 'bg-red-50 border-red-200 text-red-800'"
      >
        <!-- Icon -->
        <div class="mt-0.5 shrink-0">
          <svg v-if="toast.type === 'success'" class="w-5 h-5 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
          <svg v-else class="w-5 h-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        <!-- Message -->
        <div class="flex-1">
          <p class="font-bold text-sm">{{ toast.type === 'success' ? 'Berhasil' : 'Gagal' }}</p>
          <p class="text-sm mt-0.5 opacity-80">{{ toast.message }}</p>
        </div>
        <!-- Close -->
        <button @click="toast.show = false" class="text-current opacity-50 hover:opacity-100 transition-opacity mt-0.5">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </Transition>

    <!-- Header -->
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div>
        <h2 class="text-3xl font-extrabold tracking-tight text-slate-800">Participant Management</h2>
        <p class="text-slate-500 mt-1">Kelola semua peserta ujian dalam sistem.</p>
      </div>
      <div class="flex items-center gap-3 w-full md:w-auto flex-wrap">
        <!-- Search -->
        <div class="relative w-full md:w-[240px]">
          <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            v-model="searchQuery"
            placeholder="Cari nama atau email..."
            class="w-full pl-9 pr-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <!-- Status Filter -->
        <select
          v-model="filterStatus"
          class="px-4 py-2 text-sm border border-slate-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-slate-700 font-medium"
        >
          <option value="">Semua Status</option>
          <option value="Active">Active</option>
          <option value="Inactive">Inactive</option>
        </select>
        <!-- Add -->
        <Button @click="openAdd" class="h-10 px-4 rounded-xl bg-blue-600 hover:bg-blue-700 text-white font-semibold flex items-center gap-2 whitespace-nowrap">
          <UserPlusIcon class="w-4 h-4" /> Tambah Peserta
        </Button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid gap-4 grid-cols-2 lg:grid-cols-4">
      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white hover:shadow-md transition-all">
        <CardContent class="p-5">
          <div class="flex justify-between items-start mb-3">
            <h3 class="font-semibold text-[11px] text-slate-500 uppercase tracking-widest">Total Peserta</h3>
            <div class="w-8 h-8 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
              <UsersIcon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.total }}</div>
        </CardContent>
      </Card>
      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white hover:shadow-md transition-all">
        <CardContent class="p-5">
          <div class="flex justify-between items-start mb-3">
            <h3 class="font-semibold text-[11px] text-slate-500 uppercase tracking-widest">Sudah Ujian</h3>
            <div class="w-8 h-8 rounded-lg bg-green-50 text-green-600 flex items-center justify-center">
              <UserCheckIcon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.hasSession }}</div>
        </CardContent>
      </Card>
      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white hover:shadow-md transition-all">
        <CardContent class="p-5">
          <div class="flex justify-between items-start mb-3">
            <h3 class="font-semibold text-[11px] text-slate-500 uppercase tracking-widest">Rata-rata Skor</h3>
            <div class="w-8 h-8 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center">
              <BarChart2Icon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.avgScore }}<span class="text-base font-semibold text-slate-400">%</span></div>
        </CardContent>
      </Card>
      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white hover:shadow-md transition-all">
        <CardContent class="p-5">
          <div class="flex justify-between items-start mb-3">
            <h3 class="font-semibold text-[11px] text-slate-500 uppercase tracking-widest">Belum Ujian</h3>
            <div class="w-8 h-8 rounded-lg bg-orange-50 text-orange-600 flex items-center justify-center">
              <ClockIcon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.noSession }}</div>
        </CardContent>
      </Card>
    </div>

    <!-- Table -->
    <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20 text-slate-400">
        <svg class="w-6 h-6 animate-spin mr-2 text-blue-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
        </svg>
        Memuat data peserta...
      </div>

      <Table v-else>
        <TableHeader>
          <TableRow class="border-b border-slate-100 hover:bg-transparent">
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6">Nama</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Email</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap">Sesi Terakhir</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap">Avg Score</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Status</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6 text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="p in paginatedParticipants"
            :key="p.id"
            class="border-b border-slate-100 transition-colors hover:bg-slate-50/50 group"
          >
            <TableCell class="py-4 px-6">
              <div class="flex items-center gap-3">
                <img
                  :src="`https://api.dicebear.com/7.x/initials/svg?seed=${p.name || p.email}&backgroundColor=1e3a8a&textColor=ffffff`"
                  class="w-9 h-9 rounded-full shrink-0"
                />
                <div>
                  <div class="font-bold text-slate-800 text-sm leading-none">{{ p.name || '—' }}</div>
                  <div class="text-[11px] text-slate-400 mt-0.5 font-mono">CAND-{{ (p.id || '').substring(0, 8).toUpperCase() }}</div>
                </div>
              </div>
            </TableCell>
            <TableCell class="py-4 px-4 text-sm text-slate-500">{{ p.email || '—' }}</TableCell>
            <TableCell class="py-4 px-4 text-sm">
              <div v-if="p.last_session_name" class="text-slate-800 font-semibold text-[13px] max-w-[160px] truncate">{{ p.last_session_name }}</div>
              <div v-else class="text-slate-400 text-[12px] italic">Belum pernah ujian</div>
              <div v-if="p.last_session_name" class="text-slate-400 text-[11px] mt-0.5">{{ formatDate(p.last_session_date) }}</div>
            </TableCell>
            <TableCell class="py-4 px-4">
              <div class="flex items-center gap-2">
                <span class="font-bold text-slate-800 text-sm w-9">{{ Math.round(p.avg_score || 0) }}%</span>
                <div class="w-16 h-1.5 bg-slate-100 rounded-full overflow-hidden">
                  <div class="h-full rounded-full transition-all" :class="getScoreBarColor(p.avg_score || 0)" :style="{ width: `${Math.min(100, p.avg_score || 0)}%` }"></div>
                </div>
              </div>
            </TableCell>
            <TableCell class="py-4 px-4">
              <Badge :class="getStatusBadgeClass(p.status)" class="font-medium text-[11px] px-2.5 py-0.5 rounded-full border-0">
                {{ p.status || 'Active' }}
              </Badge>
            </TableCell>
            <TableCell class="py-4 px-6 text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon" class="h-8 w-8 rounded-lg hover:bg-blue-50 text-slate-400 hover:text-blue-600" @click="openEdit(p)" title="Edit">
                  <PencilIcon class="w-3.5 h-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="h-8 w-8 rounded-lg hover:bg-red-50 text-slate-400 hover:text-red-600" @click="confirmDelete(p)" title="Hapus">
                  <Trash2Icon class="w-3.5 h-3.5" />
                </Button>
              </div>
            </TableCell>
          </TableRow>

          <TableRow v-if="filteredParticipants.length === 0">
            <TableCell colspan="6" class="text-center py-20 text-slate-400">
              <UsersIcon class="w-10 h-10 mx-auto mb-3 opacity-20" />
              <p class="font-semibold">Tidak ada peserta ditemukan.</p>
              <p class="text-xs mt-1">Coba ubah filter atau tambahkan peserta baru.</p>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination -->
      <div v-if="filteredParticipants.length > 0 && !loading" class="p-4 border-t border-slate-100 flex items-center justify-between text-sm text-slate-500 bg-slate-50/30">
        <div>Showing <strong>{{ (currentPage - 1) * perPage + 1 }}</strong> – <strong>{{ Math.min(currentPage * perPage, filteredParticipants.length) }}</strong> of <strong>{{ filteredParticipants.length }}</strong></div>
        <div class="flex gap-1 items-center">
          <Button variant="outline" size="sm" class="h-8 border-slate-200 text-slate-600 bg-white" :disabled="currentPage === 1" @click="currentPage--">Previous</Button>
          <Button
            v-for="page in totalPages" :key="page"
            size="sm"
            :variant="page === currentPage ? 'default' : 'outline'"
            :class="page === currentPage ? 'h-8 w-8 p-0 bg-blue-600 hover:bg-blue-700 text-white font-bold' : 'h-8 w-8 p-0 border-slate-200 bg-white text-slate-600'"
            @click="currentPage = page"
          >{{ page }}</Button>
          <Button variant="outline" size="sm" class="h-8 border-slate-200 text-slate-600 bg-white" :disabled="currentPage === totalPages" @click="currentPage++">Next</Button>
        </div>
      </div>
    </Card>

    <!-- ─── Add / Edit Modal ─────────────────────────────────────────── -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 space-y-5 z-10">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-extrabold text-slate-800">{{ editingId ? 'Edit Peserta' : 'Tambah Peserta Baru' }}</h3>
            <button @click="showModal = false" class="w-8 h-8 rounded-full hover:bg-slate-100 flex items-center justify-center text-slate-400">✕</button>
          </div>

          <!-- errors are now shown as toast notifications -->

          <div class="space-y-4">
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Nama Lengkap *</label>
              <input v-model="form.name" type="text" placeholder="Nama peserta" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Email *</label>
              <input v-model="form.email" type="email" placeholder="email@perusahaan.com" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div v-if="!editingId">
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Password *</label>
              <input v-model="form.password" type="password" placeholder="Minimal 8 karakter" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
          </div>

          <div class="flex gap-3 pt-2">
            <Button variant="outline" class="flex-1 rounded-xl" @click="showModal = false">Batal</Button>
            <Button class="flex-1 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold" :disabled="saving" @click="saveParticipant">
              {{ saving ? 'Menyimpan...' : (editingId ? 'Simpan Perubahan' : 'Tambah Peserta') }}
            </Button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ─── Delete Confirm Modal ─────────────────────────────────────── -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="showDeleteModal = false">
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6 z-10 text-center space-y-4">
          <div class="w-14 h-14 rounded-full bg-red-100 flex items-center justify-center mx-auto">
            <Trash2Icon class="w-6 h-6 text-red-500" />
          </div>
          <h3 class="text-lg font-extrabold text-slate-800">Hapus Peserta?</h3>
          <p class="text-sm text-slate-500">Akun <strong>{{ deletingParticipant?.name }}</strong> akan dihapus permanen dari sistem.</p>
          <div class="flex gap-3">
            <Button variant="outline" class="flex-1 rounded-xl" @click="showDeleteModal = false">Batal</Button>
            <Button variant="destructive" class="flex-1 rounded-xl font-bold" :disabled="deleting" @click="deleteParticipant">
              {{ deleting ? 'Menghapus...' : 'Hapus' }}
            </Button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import client from '@/api/client'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import {
  SearchIcon, UserPlusIcon, UsersIcon, UserCheckIcon,
  BarChart2Icon, ClockIcon, PencilIcon, Trash2Icon
} from 'lucide-vue-next'

// ─── Toast ────────────────────────────────────────────────────────
const toast = ref({ show: false, type: 'success' as 'success' | 'error', message: '' })
let toastTimer: ReturnType<typeof setTimeout> | null = null
const showToast = (type: 'success' | 'error', message: string) => {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value = { show: true, type, message }
  toastTimer = setTimeout(() => { toast.value.show = false }, 4000)
}

// ─── State ────────────────────────────────────────────────────────
const participants = ref<any[]>([])
const loading = ref(true)
const searchQuery = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const perPage = 10

// ─── Stats ────────────────────────────────────────────────────────
const stats = computed(() => {
  const total = participants.value.length
  const hasSession = participants.value.filter(p => p.last_session_name && p.last_session_name !== '').length
  const noSession = total - hasSession
  const withScores = participants.value.filter(p => (p.avg_score || 0) > 0)
  const avgScore = withScores.length
    ? Math.round(withScores.reduce((s, p) => s + (p.avg_score || 0), 0) / withScores.length)
    : 0
  return { total, hasSession, noSession, avgScore }
})

// ─── Data Fetching ────────────────────────────────────────────────
const loadParticipants = async () => {
  loading.value = true
  try {
    const res = await client.get('/participants')
    participants.value = res.data || []
  } catch (err) {
    console.error('Failed to fetch participants', err)
  } finally {
    loading.value = false
  }
}

onMounted(loadParticipants)

// ─── Filtering & Pagination ───────────────────────────────────────
const filteredParticipants = computed(() => {
  const q = searchQuery.value.toLowerCase()
  let list = participants.value

  if (q) {
    list = list.filter(p =>
      (p.name || '').toLowerCase().includes(q) ||
      (p.email || '').toLowerCase().includes(q)
    )
  }
  if (filterStatus.value) {
    list = list.filter(p => (p.status || 'Active') === filterStatus.value)
  }
  return list
})

const totalPages = computed(() => Math.max(1, Math.ceil(filteredParticipants.value.length / perPage)))

const paginatedParticipants = computed(() => {
  const start = (currentPage.value - 1) * perPage
  return filteredParticipants.value.slice(start, start + perPage)
})

// Reset page on filter change
watch([searchQuery, filterStatus], () => { currentPage.value = 1 })

// ─── CRUD ─────────────────────────────────────────────────────────
const showModal = ref(false)
const editingId = ref<string | null>(null)
const form = ref({ name: '', email: '', password: '' })
const saving = ref(false)

const openAdd = () => {
  editingId.value = null
  form.value = { name: '', email: '', password: '' }
  showModal.value = true
}

const openEdit = (p: any) => {
  editingId.value = p.id
  form.value = { name: p.name || '', email: p.email || '', password: '' }
  showModal.value = true
}

const saveParticipant = async () => {
  if (!form.value.name.trim()) { showToast('error', 'Nama tidak boleh kosong.'); return }
  if (!form.value.email.trim()) { showToast('error', 'Email tidak boleh kosong.'); return }
  if (!editingId.value && !form.value.password.trim()) { showToast('error', 'Password tidak boleh kosong.'); return }

  saving.value = true
  try {
    if (editingId.value) {
      await client.put(`/participants/${editingId.value}`, {
        name: form.value.name,
        email: form.value.email,
      })
      showToast('success', `Data peserta ${form.value.name} berhasil diperbarui.`)
    } else {
      await client.post('/auth/register', {
        name: form.value.name,
        email: form.value.email,
        password: form.value.password,
        role: 'participant',
      })
      showToast('success', `Peserta ${form.value.name} berhasil ditambahkan.`)
    }
    showModal.value = false
    await loadParticipants()
  } catch (e: any) {
    showToast('error', e?.response?.data?.error || 'Gagal menyimpan peserta.')
  } finally {
    saving.value = false
  }
}

// Delete
const showDeleteModal = ref(false)
const deletingParticipant = ref<any>(null)
const deleting = ref(false)

const confirmDelete = (p: any) => {
  deletingParticipant.value = p
  showDeleteModal.value = true
}

const deleteParticipant = async () => {
  if (!deletingParticipant.value) return
  deleting.value = true
  const name = deletingParticipant.value.name
  try {
    await client.delete(`/participants/${deletingParticipant.value.id}`)
    showDeleteModal.value = false
    showToast('success', `Peserta ${name} berhasil dihapus.`)
    await loadParticipants()
  } catch (e: any) {
    showDeleteModal.value = false
    showToast('error', e?.response?.data?.error || 'Gagal menghapus peserta.')
  } finally {
    deleting.value = false
  }
}

// ─── Helpers ──────────────────────────────────────────────────────
const formatDate = (d: string) => {
  if (!d || d.includes('1970')) return ''
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const getStatusBadgeClass = (status: string) => {
  const s = (status || 'Active').toLowerCase()
  if (s === 'active') return 'bg-green-100 text-green-700'
  if (s === 'inactive') return 'bg-slate-100 text-slate-600'
  if (s === 'suspended') return 'bg-red-100 text-red-600'
  return 'bg-blue-100 text-blue-700'
}

const getScoreBarColor = (score: number) => {
  if (score >= 80) return 'bg-green-500'
  if (score >= 60) return 'bg-blue-500'
  if (score >= 40) return 'bg-orange-400'
  return 'bg-red-400'
}
</script>
