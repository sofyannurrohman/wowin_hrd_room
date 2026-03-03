<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold tracking-tight">Edit Posisi Pekerjaan</h2>
    </div>

    <Card>
      <CardContent class="pt-6 space-y-6">
        <Alert v-if="alertMessage" :variant="alertVariant">
          <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
          <AlertDescription>{{ alertMessage }}</AlertDescription>
        </Alert>

        <div v-if="fetching" class="text-muted-foreground py-10 text-center">
          Memuat data posisi...
        </div>

        <div v-else class="space-y-4">
          <div class="space-y-2">
            <Label for="name">Nama Posisi</Label>
            <Input id="name" v-model="form.name" placeholder="Misal: Software Engineer" />
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

        <div class="flex justify-end gap-2 pt-4 border-t" v-if="!fetching">
          <Button variant="outline" @click="router.back()">Kembali</Button>
          <Button @click="submit" :disabled="loading">{{ loading ? 'Menyimpan...' : 'Simpan Perubahan' }}</Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Checkbox } from '@/components/ui/checkbox'

const router = useRouter()
const route = useRoute()
const positionId = route.params.id as string

const form = ref({
  name: '',
  is_active: true
})

const loading = ref(false)
const fetching = ref(true)
const alertMessage = ref('')
const alertVariant = ref<'default' | 'destructive'>('default')

const showSuccess = (message: string) => {
  alertVariant.value = 'default'
  alertMessage.value = message
}

const showError = (message: string) => {
  alertVariant.value = 'destructive'
  alertMessage.value = message
}

onMounted(async () => {
  try {
    const res = await client.get(`/job-positions/${positionId}`)
    const p = res.data.position
    form.value.name = p.name
    form.value.is_active = p.is_active
  } catch (e: any) {
    showError(e.response?.data?.error || 'Gagal memuat posisi')
  } finally {
    fetching.value = false
  }
})

const submit = async () => {
  if (!form.value.name) {
    showError('Nama posisi harus diisi!')
    return
  }

  loading.value = true
  alertMessage.value = ''
  
  try {
    await client.put(`/job-positions/${positionId}`, {
      name: form.value.name,
      is_active: form.value.is_active
    })
    showSuccess('Perubahan posisi berhasil disimpan')
    setTimeout(() => {
      router.push('/job-positions')
    }, 1500)
  } catch (error: any) {
    showError(error.response?.data?.error || 'Gagal update posisi')
  } finally {
    loading.value = false
  }
}
</script>
