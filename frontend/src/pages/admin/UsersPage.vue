<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Manajemen User (HR)</h2>
      <Button @click="showCreateDialog = true">
        <UserPlusIcon class="mr-2 h-4 w-4" /> Tambah HR
      </Button>
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('name')">
              <div class="flex items-center gap-1">
                Nama
                <ArrowUpDownIcon v-if="sortConfig.key !== 'name'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
            <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('email')">
              <div class="flex items-center gap-1">
                Email
                <ArrowUpDownIcon v-if="sortConfig.key !== 'email'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
            <TableHead class="cursor-pointer hover:text-primary transition-colors" @click="toggleSort('role_id')">
              <div class="flex items-center gap-1">
                Role ID
                <ArrowUpDownIcon v-if="sortConfig.key !== 'role_id'" class="w-3.5 h-3.5" />
                <ArrowUpIcon v-else-if="sortConfig.direction === 'asc'" class="w-3.5 h-3.5 text-primary" />
                <ArrowDownIcon v-else class="w-3.5 h-3.5 text-primary" />
              </div>
            </TableHead>
            <TableHead class="text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="user in paginatedData" :key="user.id">
            <TableCell class="font-medium">{{ user.name || (user.role_id === 1 ? 'Super Admin' : user.role_id === 2 ? 'HR' : 'Peserta') }}</TableCell>
            <TableCell>{{ user.email }}</TableCell>
            <TableCell>
              <Badge :variant="user.role_id === 1 ? 'default' : user.role_id === 2 ? 'secondary' : 'outline'">
                {{ user.role_id === 1 ? 'Super Admin' : user.role_id === 2 ? 'HR' : 'Peserta' }}
              </Badge>
            </TableCell>
            <TableCell class="text-right flex justify-end gap-1">
              <Button variant="ghost" size="icon" @click="openEditDialog(user)">
                <EditIcon class="w-4 h-4 text-primary" />
              </Button>
              <Button variant="ghost" size="icon" @click="deleteUser(user.id)" :disabled="user.role_id === 1">
                <TrashIcon class="w-4 h-4 text-destructive" />
              </Button>
            </TableCell>
          </TableRow>
          <TableRow v-if="users.length === 0">
             <TableCell colspan="4" class="text-center h-24">Memuat pengguna...</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      
      <DataTablePagination 
        v-if="users.length > 0"
        :total="totalItems"
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
      />
    </Card>

    <!-- Create User Dialog -->
    <Dialog v-model:open="showCreateDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tambah Akun HR Baru</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>Nama Lengkap</Label>
            <Input v-model="form.name" type="text" placeholder="John Doe"/>
          </div>
          <div class="space-y-2">
            <Label>Email Akun</Label>
            <Input v-model="form.email" type="email" placeholder="hr.baru@hrdroom.com"/>
          </div>
          <div class="space-y-2">
            <Label>Password Default</Label>
            <Input v-model="form.password" type="password"/>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showCreateDialog = false">Batal</Button>
          <Button @click="createUser" :disabled="creating">Simpan</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Edit User Dialog -->
    <Dialog v-model:open="showEditDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Pengguna</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label>Nama</Label>
            <Input v-model="editForm.name" type="text" placeholder="Nama Lengkap"/>
          </div>
          <div class="space-y-2">
            <Label>Email Akun</Label>
            <Input v-model="editForm.email" type="email" placeholder="hr.baru@hrdroom.com"/>
          </div>
          <div class="space-y-2">
            <Label>Role ID (1=Super Admin, 2=HR, 3=Peserta)</Label>
            <Input v-model.number="editForm.role_id" type="number"/>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showEditDialog = false">Batal</Button>
          <Button @click="submitEditUser" :disabled="updating">Simpan</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
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
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'vue-sonner'
import { UserPlusIcon, TrashIcon, EditIcon, ArrowUpDownIcon, ArrowUpIcon, ArrowDownIcon } from 'lucide-vue-next'

const users = ref<any[]>([])

const {
  pageSize,
  currentPage,
  sortConfig,
  toggleSort,
  paginatedData,
  totalItems
} = useDataTable(users)

const showCreateDialog = ref(false)
const creating = ref(false)

const showSuccess = (message: string) => {
  toast.success(message)
}

const showError = (message: string) => {
  toast.error(message)
}

const form = ref({
  name: '',
  email: '',
  password: '',
  role_id: 2 // 2 is typically HR
})

const showEditDialog = ref(false)
const editingUser = ref<any>(null)
const updating = ref(false)
const editForm = ref({
  name: '',
  email: '',
  role_id: 2
})

const openEditDialog = (user: any) => {
  editingUser.value = user
  editForm.value = {
    name: user.name || '',
    email: user.email || '',
    role_id: user.role_id || 2
  }
  showEditDialog.value = true
}

const submitEditUser = async () => {
  if (!editingUser.value) return
  updating.value = true
  try {
    await client.put(`/admin/users/${editingUser.value.id}`, editForm.value)
    toast.success('Pengguna berhasil diperbarui')
    showEditDialog.value = false
    await loadUsers()
  } catch(e: any) {
    toast.error(e.response?.data?.error || 'Gagal memperbarui pengguna')
  } finally {
    updating.value = false
  }
}

const loadUsers = async () => {
  try {
    const res = await client.get('/admin/users')
    // Filter out participants (role_id 3) so this page only shows Admins and HRs
    users.value = (res.data || []).filter((u: any) => u.role_id !== 3)
  } catch(e) {
    showError('Gagal memuat daftar pengguna')
  }
}

onMounted(loadUsers)

const createUser = async () => {
  creating.value = true
  try {
    // Only super admin can register new HR
    await client.post('/auth/register', form.value)
    showSuccess('Pengguna HR berhasil ditambahkan')
    showCreateDialog.value = false
    form.value.name = ''
    form.value.email = ''
    form.value.password = ''
    await loadUsers()
  } catch(e: any) {
    showError(e.response?.data?.error || 'Gagal menambahkan pengguna')
  } finally {
    creating.value = false
  }
}

const deleteUser = async (id: string) => {
  if (!confirm('Hapus pengguna ini?')) return
  try {
    await client.delete(`/admin/users/${id}`)
    showSuccess('Pengguna berhasil dihapus')
    await loadUsers()
  } catch(e) {
    showError('Gagal menghapus pengguna')
  }
}
</script>
