<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold tracking-tight">Tambah Posisi Pekerjaan</h2>
    </div>

    <Card>
      <CardContent class="pt-6 space-y-6">


          <div class="space-y-4">
            <div class="space-y-2">
              <Label for="name">Nama Posisi</Label>
              <Input id="name" v-model="form.name" placeholder="Misal: Software Engineer" />
            </div>

            <div class="space-y-2">
              <Label for="description">Deskripsi Pekerjaan</Label>
              <Textarea id="description" v-model="form.description" placeholder="Jelaskan peran ini..." class="h-32" />
            </div>

            <div class="space-y-2">
              <Label for="requirements">Kualifikasi/Persyaratan</Label>
              <Textarea id="requirements" v-model="form.requirements" placeholder="Daftar persyaratan (bisa dipisah baris)..." class="h-32" />
            </div>
            
            <div class="flex items-center space-x-2 pt-2">
              <Checkbox id="is_active" v-model:checked="form.is_active" />
              <div class="grid gap-1.5 leading-none">
                <label for="is_active" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  Posisi Aktif (Buka Lowongan)
                </label>
                <p class="text-sm text-muted-foreground"> Posisi yang aktif akan muncul di pilihan halaman registrasi peserta. </p>
              </div>
            </div>
          </div>

        <div class="flex justify-end gap-2 pt-4 border-t">
          <Button variant="outline" @click="router.back()">Kembali</Button>
          <Button @click="submit" :disabled="loading">{{ loading ? 'Menyimpan...' : 'Simpan Posisi' }}</Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import client from '@/api/client'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'vue-sonner'

const router = useRouter()
const form = ref({
  name: '',
  description: '',
  requirements: '',
  is_active: true
})
const loading = ref(false)

const showSuccess = (message: string) => {
  toast.success(message)
}

const showError = (message: string) => {
  toast.error(message)
}

const submit = async () => {
  if (!form.value.name) {
    showError('Nama posisi harus diisi!')
    return
  }

  loading.value = true
  
  try {
    await client.post('/job-positions', {
      name: form.value.name,
      description: form.value.description,
      requirements: form.value.requirements,
      is_active: form.value.is_active
    })
    showSuccess('Posisi berhasil disimpan')
    setTimeout(() => {
      router.push('/job-positions')
    }, 1500)
  } catch (error: any) {
    showError(error.response?.data?.error || 'Gagal menyimpan posisi')
  } finally {
    loading.value = false
  }
}
</script>
