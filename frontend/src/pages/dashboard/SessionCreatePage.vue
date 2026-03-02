<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <div class="flex items-center gap-4">
      <Button variant="ghost" size="icon" @click="router.back()">
        <ArrowLeftIcon class="w-5 h-5" />
      </Button>
      <h2 class="text-3xl font-bold tracking-tight">Buat Sesi Ujian Baru</h2>
    </div>

    <Card>
      <CardContent class="pt-6 space-y-6">
        <div class="space-y-2">
          <Label for="title">Nama Sesi</Label>
          <Input id="title" v-model="form.title" placeholder="e.g. Rekrutmen Programmer Gelombang 1" />
        </div>
        <div class="space-y-2">
          <Label for="description">Deskripsi</Label>
          <Input id="description" v-model="form.description" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="start">Waktu Mulai</Label>
            <Input id="start" v-model="form.start_time" type="datetime-local" />
          </div>
          <div class="space-y-2">
            <Label for="end">Waktu Selesai</Label>
            <Input id="end" v-model="form.end_time" type="datetime-local" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="duration">Durasi (Menit)</Label>
            <Input id="duration" v-model.number="form.duration_minutes" type="number" min="1" max="300" />
          </div>
          <div class="space-y-2">
            <Label for="limit">Batas Peserta Aktif (Maks 20)</Label>
            <Input id="limit" v-model.number="form.participant_limit" type="number" min="1" max="20" />
          </div>
        </div>
        <div class="flex justify-end pt-4">
          <Button @click="submit" :disabled="loading">
            {{ loading ? 'Menyimpan...' : 'Simpan Sesi' }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { ArrowLeftIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const router = useRouter()
const sessionStore = useSessionStore()

const now = new Date()
now.setMinutes(now.getMinutes() - now.getTimezoneOffset())

const form = ref({
  title: '',
  description: '',
  start_time: now.toISOString().slice(0, 16),
  end_time: new Date(now.getTime() + 2 * 60 * 60 * 1000).toISOString().slice(0, 16),
  duration_minutes: 120,
  participant_limit: 20
})

const loading = ref(false)

const submit = async () => {
  if (!form.value.title) {
    toast.error('Nama sesi harus diisi')
    return
  }
  loading.value = true
  try {
    const payload = {
      ...form.value,
      start_time: new Date(form.value.start_time).toISOString(),
      end_time: new Date(form.value.end_time).toISOString()
    }
    await sessionStore.createSession(payload)
    toast.success('Sesi berhasil dibuat')
    router.push('/sessions')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Gagal membuat sesi')
  } finally {
    loading.value = false
  }
}
</script>
