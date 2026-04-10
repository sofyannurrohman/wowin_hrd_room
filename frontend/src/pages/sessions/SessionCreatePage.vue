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
        <Alert v-if="alertMessage" :variant="alertVariant">
          <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
          <AlertDescription>{{ alertMessage }}</AlertDescription>
        </Alert>

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

        <div class="space-y-4 pt-4 border-t">
          <Label>Pilih Modul Soal</Label>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
             <div v-for="m in modules" :key="m.id" class="flex items-start space-x-3 p-3 border rounded-lg bg-card hover:bg-slate-50 transition-colors cursor-pointer" @click="toggleModule(m.id)">
                 <input type="checkbox" :id="'mod-' + m.id" :checked="selectedModules.includes(m.id)" class="mt-1 w-4 h-4 rounded border-slate-300 accent-blue-600" />
                 <div class="space-y-1.5 leading-none flex-1">
                     <label :for="'mod-' + m.id" class="font-medium cursor-pointer block select-none">{{ m.name }}</label>
                     <p class="text-xs text-muted-foreground line-clamp-2 select-none">{{ m.description || 'Tidak ada deskripsi' }}</p>
                 </div>
             </div>
          </div>
          <p class="text-xs text-muted-foreground mt-2">Peserta akan mengerjakan modul secara berurutan sesuai urutan Anda memilih box di atas.</p>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import client from '@/api/client'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { ArrowLeftIcon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

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
const alertMessage = ref('')
const alertVariant = ref<'default' | 'destructive'>('default')

const modules = ref<any[]>([])
const selectedModules = ref<string[]>([])

const toggleModule = (id: string) => {
  const index = selectedModules.value.indexOf(id)
  if (index === -1) {
    selectedModules.value.push(id)
  } else {
    selectedModules.value.splice(index, 1)
  }
}

onMounted(async () => {
  try {
    const res = await client.get('/modules')
    modules.value = res.data.modules || []
  } catch(e) {
    showError('Gagal memuat daftar modul')
  }
})

const showSuccess = (message: string) => {
  alertVariant.value = 'default'
  alertMessage.value = message
}

const showError = (message: string) => {
  alertVariant.value = 'destructive'
  alertMessage.value = message
}

const submit = async () => {
  if (!form.value.title) {
    showError('Nama sesi harus diisi')
    return
  }
  if (selectedModules.value.length === 0) {
    showError('Silakan pilih minimal satu modul soal')
    return
  }
  loading.value = true
  alertMessage.value = ''
  try {
    const payload = {
      name: form.value.title,
      schedule: new Date(form.value.start_time).toISOString(),
      duration_minutes: form.value.duration_minutes,
      max_participants: form.value.participant_limit,
      randomize_questions: false,
      show_score: false,
      module_ids: selectedModules.value
    }
    await sessionStore.createSession(payload)
    showSuccess('Sesi berhasil dibuat')
    router.push('/sessions')
  } catch (e: any) {
    showError(e.response?.data?.error || 'Gagal membuat sesi')
  } finally {
    loading.value = false
  }
}
</script>
