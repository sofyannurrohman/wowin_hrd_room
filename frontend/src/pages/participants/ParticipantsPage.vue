<template>
  <div class="space-y-6 w-full max-w-[1200px] mx-auto pb-10">

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
        <!-- Add & Import -->
        <Button @click="downloadTemplate" variant="outline" class="h-10 px-4 rounded-xl border-slate-200 bg-white text-slate-600 hover:bg-slate-50 font-semibold flex items-center gap-2 whitespace-nowrap">
          <DownloadIcon class="w-4 h-4" /> Template CSV
        </Button>
        <Button @click="triggerCsvUpload" :disabled="uploading" variant="outline" class="h-10 px-4 rounded-xl border-blue-200 bg-blue-50 text-blue-700 hover:bg-blue-100 font-semibold flex items-center gap-2 whitespace-nowrap">
          <UploadIcon v-if="!uploading" class="w-4 h-4" />
          <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
          </svg>
          {{ uploading ? 'Mengunggah...' : 'Import CSV' }}
        </Button>
        <input type="file" accept=".csv" ref="csvInput" class="hidden" @change="handleCsvUpload" />
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
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6 cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('name')">
              <div class="flex items-center gap-1">
                Nama & Posisi
                <ArrowUpDownIcon v-if="sortConfig.key !== 'name'" class="w-3 h-3" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
              </div>
            </TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Kontak & Gaji</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('last_session_name')">
              <div class="flex items-center gap-1">
                Sesi Terakhir
                <ArrowUpDownIcon v-if="sortConfig.key !== 'last_session_name'" class="w-3 h-3" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
              </div>
            </TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('avg_score')">
              <div class="flex items-center gap-1">
                Avg Score
                <ArrowUpDownIcon v-if="sortConfig.key !== 'avg_score'" class="w-3 h-3" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
              </div>
            </TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 cursor-pointer hover:text-blue-600 transition-colors" @click="toggleSort('status')">
              <div class="flex items-center gap-1">
                Status
                <ArrowUpDownIcon v-if="sortConfig.key !== 'status'" class="w-3 h-3" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3 h-3 text-blue-600" />
                <ArrowDownIcon v-else class="w-3 h-3 text-blue-600" />
              </div>
            </TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6 text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="p in paginatedData"
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
                  <div @click="openDetail(p)" class="font-bold text-slate-800 text-sm leading-none cursor-pointer hover:text-blue-600 transition-colors">{{ p.name || '—' }} <span v-if="p.age" class="text-xs text-slate-400 font-normal">({{p.age}} th)</span></div>
                  <div class="text-[11px] text-blue-600 mt-1 font-semibold flex items-center gap-1">
                    {{ p.applied_position || 'Umum' }}
                    <span v-if="p.last_education" class="text-slate-400 text-[10px] ml-1 px-1.5 py-0.5 bg-blue-50 rounded-md">{{ p.last_education }}</span>
                  </div>
                </div>
              </div>
            </TableCell>
            <TableCell class="py-4 px-4">
               <div class="text-sm text-slate-700 font-medium">{{ p.email || '—' }}</div>
               <div v-if="p.whatsapp_number" class="text-xs text-slate-500 mt-0.5 flex items-center gap-1 group">
                 <svg class="w-3 h-3 text-emerald-500" fill="currentColor" viewBox="0 0 24 24"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51a12.8 12.8 0 0 0-.57-.01c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"/></svg>
                 <a :href="getWaLink(p.whatsapp_number)" target="_blank" class="group-hover:text-emerald-600 transition-colors">{{ p.whatsapp_number }}</a>
               </div>
               <div v-if="p.expected_salary" class="text-[11px] text-emerald-600 font-semibold mt-0.5">Ekspektasi: Rp {{ p.expected_salary }}</div>
               <div v-if="p.address" class="text-[10px] text-slate-400 mt-1 max-w-[150px] truncate" :title="p.address">{{ p.address }}</div>
            </TableCell>
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
                <a v-if="p.whatsapp_number" :href="getWaLink(p.whatsapp_number)" target="_blank" class="inline-flex items-center justify-center h-8 w-8 rounded-lg hover:bg-emerald-50 text-slate-400 hover:text-emerald-500 transition-colors" title="Hubungi via WA">
                  <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51a12.8 12.8 0 0 0-.57-.01c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"/></svg>
                </a>
                <Button variant="ghost" size="icon" class="h-8 w-8 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-700" @click="openDetail(p)" title="Detail">
                  <EyeIcon class="w-3.5 h-3.5" />
                </Button>
                <Button v-if="p.cv_url" variant="ghost" size="icon" class="h-8 w-8 rounded-lg hover:bg-emerald-50 text-emerald-500 hover:text-emerald-700" @click="downloadCV(p.cv_url)" title="Unduh CV">
                  <DownloadIcon class="w-3.5 h-3.5" />
                </Button>
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
      <div v-if="filteredParticipants.length > 0 && !loading" class="p-4 border-t border-slate-100 bg-slate-50/30">
        <DataTablePagination 
          :total="totalItems"
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
        />
      </div>
    </Card>

    <!-- ─── Add / Edit Modal ─────────────────────────────────────────── -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-xl p-6 space-y-5 z-10">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-extrabold text-slate-800">{{ editingId ? 'Edit Peserta' : 'Tambah Peserta Baru' }}</h3>
            <button @click="showModal = false" class="w-8 h-8 rounded-full hover:bg-slate-100 flex items-center justify-center text-slate-400">✕</button>
          </div>

          <!-- errors are now shown as toast notifications -->

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="col-span-1 md:col-span-2">
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Nama Lengkap *</label>
              <input v-model="form.name" type="text" placeholder="Nama peserta" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Email *</label>
              <input v-model="form.email" type="email" placeholder="email@perusahaan.com" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">WhatsApp *</label>
              <input v-model="form.whatsapp_number" type="text" placeholder="0812..." class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Posisi Dilamar</label>
              <input v-model="form.applied_position" type="text" placeholder="e.g. Sales" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Umur</label>
              <input v-model.number="form.age" type="number" placeholder="Thn" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Pendidikan Terakhir</label>
              <input v-model="form.last_education" type="text" placeholder="e.g. S1" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div>
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Ekspektasi Gaji</label>
              <input v-model="form.expected_salary" type="text" placeholder="e.g. 5.000.000" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500" />
            </div>
            <div class="col-span-1 md:col-span-2">
              <label class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1 block">Alamat</label>
              <textarea v-model="form.address" rows="2" placeholder="Alamat lengkap" class="w-full px-4 py-2.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"></textarea>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import client from '@/api/client'
import { toast } from "vue-sonner"
import { useDataTable } from '@/composables/useDataTable'
import DataTablePagination from '@/components/shared/DataTablePagination.vue'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import {
  SearchIcon, UserPlusIcon, UsersIcon, UserCheckIcon,
  BarChart2Icon, ClockIcon, PencilIcon, Trash2Icon,
  DownloadIcon, UploadIcon, EyeIcon, ArrowUpDownIcon, ArrowUpIcon, ArrowDownIcon
} from 'lucide-vue-next'

const router = useRouter()

const showToast = (type: 'success' | 'error', message: string) => {
  if (type === 'success') {
    toast.success(message)
  } else {
    toast.error(message)
  }
}

const participants = ref<any[]>([])
const loading = ref(true)
const searchQuery = ref('')
const filterStatus = ref('')

// ─── CSV Import Refs ──────────────────────────────────────────────
const csvInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

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

// ─── CSV Import Logic ──────────────────────────────────────────────
const downloadTemplate = () => {
  const csvContent = "data:text/csv;charset=utf-8,Nama,Email,Umur,Posisi\nJohn Doe,john@example.com,25,Software Engineer\nJane Smith,jane@example.com,28,Product Manager"
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement("a")
  link.setAttribute("href", encodedUri)
  link.setAttribute("download", "template_peserta.csv")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const triggerCsvUpload = () => {
  csvInput.value?.click()
}

const handleCsvUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return
  
  const file = target.files[0] as File
  const formData = new FormData()
  formData.append('file', file)

  uploading.value = true
  try {
    const res = await client.post('/participants/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    // Reset input
    target.value = ''
    
    showToast('success', `${res.data.imported} peserta diimpor. ${res.data.skipped} dilewati.`)
    await loadParticipants()
  } catch (err: any) {
    target.value = ''
    showToast('error', err.response?.data?.error || 'Gagal mengimpor file CSV.')
  } finally {
    uploading.value = false
  }
}

const filteredParticipants = computed(() => {
  const q = searchQuery.value.toLowerCase()
  return participants.value.filter(p => {
    const matchesSearch = (p.name || '').toLowerCase().includes(q) || (p.email || '').toLowerCase().includes(q)
    if (!matchesSearch) return false
    
    if (filterStatus.value && (p.status || 'Active') !== filterStatus.value) {
      return false
    }
    return true
  })
})

const {
  pageSize,
  currentPage,
  sortConfig,
  toggleSort,
  paginatedData,
  totalItems
} = useDataTable(filteredParticipants)

// ─── CRUD ─────────────────────────────────────────────────────────
const showModal = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  email: '',
  age: null as number | null,
  applied_position: '',
  whatsapp_number: '',
  last_education: '',
  expected_salary: '',
  address: ''
})
const saving = ref(false)

const openDetail = (p: any) => {
  router.push(`/participants/${p.id}`)
}

const openAdd = () => {
  editingId.value = null
  form.value = {
    name: '',
    email: '',
    age: null,
    applied_position: '',
    whatsapp_number: '',
    last_education: '',
    expected_salary: '',
    address: ''
  }
  showModal.value = true
}

const openEdit = (p: any) => {
  editingId.value = p.id
  form.value = {
    name: p.name || '',
    email: p.email || '',
    age: p.age || null,
    applied_position: p.applied_position || '',
    whatsapp_number: p.whatsapp_number || '',
    last_education: p.last_education || '',
    expected_salary: p.expected_salary || '',
    address: p.address || ''
  }
  showModal.value = true
}

const saveParticipant = async () => {
  if (!form.value.name.trim()) { showToast('error', 'Nama tidak boleh kosong.'); return }
  if (!form.value.email.trim()) { showToast('error', 'Email tidak boleh kosong.'); return }

  saving.value = true
  try {
    if (editingId.value) {
      await client.put(`/participants/${editingId.value}`, {
        name: form.value.name,
        email: form.value.email,
        age: form.value.age,
        applied_position: form.value.applied_position,
        whatsapp_number: form.value.whatsapp_number,
        last_education: form.value.last_education,
        expected_salary: form.value.expected_salary,
        address: form.value.address
      })
      showToast('success', `Data peserta ${form.value.name} berhasil diperbarui.`)
    } else {
      // Generate a random internal password as users login via magic links/tokens
      const randomPassword = Math.random().toString(36).slice(-10) + Math.random().toString(36).slice(-10)
      
      await client.post('/auth/register', {
        name: form.value.name,
        email: form.value.email,
        password: randomPassword,
        role: 'participant',
        age: form.value.age,
        applied_position: form.value.applied_position,
        whatsapp_number: form.value.whatsapp_number,
        last_education: form.value.last_education,
        expected_salary: form.value.expected_salary,
        address: form.value.address
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
const getWaLink = (phone: string) => {
  if (!phone) return '#';
  let formatted = phone.replace(/\D/g, '');
  if (formatted.startsWith('0')) {
    formatted = '62' + formatted.slice(1);
  }
  return `https://wa.me/${formatted}`;
}

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

const downloadCV = (url: string) => {
  if (!url) return
  
  // Use VITE_API_URL if provided (removing /api for static uploads), 
  // otherwise fallback to a relative path which works via Nginx/Vite proxy.
  const backendUrl = import.meta.env.VITE_API_URL
  if (backendUrl) {
    const baseUrl = backendUrl.replace(/\/api$/, '')
    window.open(`${baseUrl}${url}`, '_blank')
  } else {
    window.open(url, '_blank')
  }
}
</script>
