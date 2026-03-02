<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Manajemen Pengguna (HR)</h2>
      <Button @click="showCreateDialog = true">
        <UserPlusIcon class="mr-2 h-4 w-4" /> Tambah HR
      </Button>
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Nama</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Role ID</TableHead>
            <TableHead class="text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="user in users" :key="user.id">
            <TableCell class="font-medium">{{ user.name || 'Admin / HR' }}</TableCell>
            <TableCell>{{ user.email }}</TableCell>
            <TableCell>
              <Badge :variant="user.role_id === 1 ? 'default' : 'outline'">{{ user.role_id === 1 ? 'Super Admin' : 'HR' }}</Badge>
            </TableCell>
            <TableCell class="text-right">
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
    </Card>

    <!-- Create User Dialog -->
    <Dialog v-model:open="showCreateDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tambah Akun HR Baru</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 py-4">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import client from '@/api/client'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { UserPlusIcon, TrashIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const users = ref<any[]>([])
const showCreateDialog = ref(false)
const creating = ref(false)

const form = ref({
  email: '',
  password: '',
  role_id: 2 // 2 is typically HR
})

const loadUsers = async () => {
  try {
    const res = await client.get('/admin/users')
    users.value = res.data || []
  } catch(e) {
    toast.error('Gagal memuat daftar pengguna')
  }
}

onMounted(loadUsers)

const createUser = async () => {
  creating.value = true
  try {
    // Only super admin can register new HR
    await client.post('/auth/register', form.value)
    toast.success('Pengguna HR berhasil ditambahkan')
    showCreateDialog.value = false
    form.value.email = ''
    form.value.password = ''
    await loadUsers()
  } catch(e: any) {
    toast.error(e.response?.data?.error || 'Gagal menambahkan pengguna')
  } finally {
    creating.value = false
  }
}

const deleteUser = async (id: string) => {
  if (!confirm('Hapus pengguna ini?')) return
  try {
    await client.delete(`/admin/users/${id}`)
    toast.success('Pengguna berhasil dihapus')
    await loadUsers()
  } catch(e) {
    toast.error('Gagal menghapus pengguna')
  }
}
</script>
