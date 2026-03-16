<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <!-- Success State -->
    <div v-if="submitted" class="max-w-md w-full bg-white rounded-3xl shadow-xl border border-slate-100 p-8 text-center space-y-6">
      <div class="w-20 h-20 bg-green-50 rounded-full flex items-center justify-center mx-auto mb-2">
        <svg class="w-10 h-10 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <h2 class="text-2xl font-black text-slate-800">Lamaran Terkirim!</h2>
      <p class="text-slate-500 leading-relaxed text-sm">
        Terima kasih, <strong>{{ form.name }}</strong>. Data Anda telah kami terima dengan baik. 
        Tim HRD kami akan meninjau lamaran Anda, dan <strong>Link Ujian akan dikirimkan langsung ke email Anda</strong> jika Anda lolos ke tahap selanjutnya.
      </p>
      <img :src="viteLogo" alt="Logo" class="w-10 h-10 mx-auto opacity-50 grayscale mt-6" />
    </div>

    <!-- Application Form -->
    <div v-else class="max-w-2xl w-full bg-white rounded-3xl shadow-xl border border-slate-100 overflow-hidden">
      <!-- Header -->
      <div class="bg-blue-600 p-8 text-center relative overflow-hidden">
        <div class="absolute top-4 left-4 z-20">
          <Button variant="ghost" class="text-white hover:bg-white/10" @click="() => $router.push('/join')">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
            Kembali
          </Button>
        </div>
        <div class="absolute inset-0 bg-gradient-to-br from-blue-500 to-indigo-700 opacity-90 leading-none"></div>
        <div class="relative z-10 flex flex-col items-center">
          <div class="w-16 h-16 bg-white/20 backdrop-blur-md rounded-2xl flex items-center justify-center mb-4">
            <img :src="viteLogo" alt="Logo" class="w-16 h-16" />
          </div>
          <h1 class="text-3xl font-black text-white tracking-tight">Portal Rekrutmen</h1>
          <p class="text-blue-100 mt-2 font-medium">Buka Peluang Karirmu Bersama Kami</p>
        </div>
      </div>

      <!-- Form Body -->
      <div class="p-8 md:p-10 space-y-6">
        <Alert v-if="errorMsg" variant="destructive" class="mb-6 border-red-200 bg-red-50 text-red-800">
          <AlertTitle class="font-bold">Gagal Mengirim</AlertTitle>
          <AlertDescription class="opacity-90">{{ errorMsg }}</AlertDescription>
        </Alert>

        <form @submit.prevent="submitApplication" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Name -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Nama Lengkap *</Label>
              <Input v-model="form.name" type="text" placeholder="Sesuai KTP" required class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl" />
            </div>

            <!-- Email -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Email Aktif *</Label>
              <Input v-model="form.email" type="email" placeholder="email@anda.com" required class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl" />
            </div>

            <!-- Age -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Umur *</Label>
              <Input v-model="form.age" type="number" placeholder="Contoh: 25" required min="17" max="65" class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl" />
            </div>

            <!-- WhatsApp Number -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Nomor HP (WhatsApp) *</Label>
              <Input v-model="form.whatsapp_number" type="tel" placeholder="Contoh: 08123456789" required class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl" />
            </div>

            <!-- Position -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Posisi yang Dilamar *</Label>
              <Select v-model="form.position" required>
                <SelectTrigger class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl w-full">
                  <SelectValue placeholder="Pilih Posisi..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem v-for="pos in positions" :key="pos.id" :value="pos.name">{{ pos.name }}</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div class="space-y-2">
            <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Alamat Tempat Tinggal Sekarang *</Label>
            <textarea v-model="form.address" placeholder="Masukkan alamat lengkap" required class="w-full h-24 p-4 bg-slate-50 border border-slate-200 focus:bg-white transition-colors text-base rounded-xl resize-none outline-none focus:ring-2 focus:ring-blue-500"></textarea>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Last Education -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Pendidikan Terakhir *</Label>
              <Select v-model="form.last_education" required>
                <SelectTrigger class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl w-full">
                  <SelectValue placeholder="Pilih Pendidikan..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="SMA/SMK">SMA / SMK Sederajat</SelectItem>
                    <SelectItem value="D3">Diploma (D3)</SelectItem>
                    <SelectItem value="D4">Diploma (D4)</SelectItem>
                    <SelectItem value="S1">Sarjana (S1)</SelectItem>
                    <SelectItem value="S2">Magister (S2)</SelectItem>
                    <SelectItem value="S3">Doktor (S3)</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>

            <!-- Expected Salary -->
            <div class="space-y-2">
              <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Gaji yang Diharapkan (Rp) *</Label>
              <Input v-model="form.expected_salary" type="text" placeholder="Contoh: 5.000.000" required class="h-12 bg-slate-50 border-slate-200 focus:bg-white transition-colors text-base rounded-xl" />
            </div>
          </div>

          <!-- CV Upload -->
          <div class="space-y-2">
            <Label class="text-xs font-bold text-slate-600 uppercase tracking-wider">Upload CV (PDF) *</Label>
            <div class="relative group">
              <input type="file" ref="cvInput" @change="handleFileDrop" accept="application/pdf" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" required />
              <div class="w-full border-2 border-dashed border-slate-200 rounded-2xl p-6 flex flex-col items-center justify-center bg-slate-50 group-hover:bg-blue-50 group-hover:border-blue-300 transition-colors" :class="{'bg-blue-50 border-blue-400': cvFile}">
                <svg v-if="!cvFile" class="w-8 h-8 text-slate-400 mb-2 group-hover:text-blue-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <svg v-else class="w-8 h-8 text-blue-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <div class="text-sm font-semibold" :class="cvFile ? 'text-blue-700' : 'text-slate-600'">
                  {{ cvFile ? cvFile.name : 'Klik atau Tarik File CV ke Sini' }}
                </div>
                <div class="text-xs text-slate-400 mt-1">Format PDF, maksimal 1MB</div>
              </div>
            </div>
          </div>

          <!-- Divider -->
          <hr class="border-slate-100 my-8" />

          <!-- Submit Button -->
          <Button type="submit" :disabled="loading" class="w-full h-14 text-lg font-bold rounded-xl bg-blue-600 hover:bg-blue-700 text-white shadow-xl shadow-blue-600/20 active:scale-[0.98] transition-all">
            <span v-if="loading" class="flex items-center justify-center gap-2">
              <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
              </svg>
              Mengirim Lamaran...
            </span>
            <span v-else>Kirim Lamaran Sekarang</span>
          </Button>
          
          <p class="text-center text-xs text-slate-400 mt-4">
            Dengan menekan tombol di atas, Anda melamar secara resmi. Pastikan data yang Anda isi sudah benar.
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import client from '@/api/client'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import viteLogo from '@/assets/vite.svg'

const submitted = ref(false)
const loading = ref(false)
const errorMsg = ref('')

const positions = ref<any[]>([])

const form = ref({
  name: '',
  email: '',
  age: '',
  whatsapp_number: '',
  position: '',
  expected_salary: '',
  address: '',
  last_education: ''
})

const cvInput = ref<HTMLInputElement | null>(null)
const cvFile = ref<File | null>(null)

onMounted(async () => {
  try {
    const res = await client.get('/job-positions/active')
    positions.value = res.data.positions || []
  } catch (err) {
    console.error('Failed to load job positions', err)
  }
})

const handleFileDrop = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    if (file.type !== 'application/pdf') {
       errorMsg.value = 'File CV harus berformat PDF.'
       target.value = ''
       cvFile.value = null
       return
    }
    if (file.size > 10 * 1024 * 1024) { // 10MB
       errorMsg.value = 'Ukuran CV maksimal 10MB.'
       target.value = ''
       cvFile.value = null
       return
    }
    errorMsg.value = ''
    cvFile.value = file
  }
}

const submitApplication = async () => {
  if (!cvFile.value) {
    errorMsg.value = 'Harap unggah CV Anda.'
    return
  }
  if (!form.value.position) {
    errorMsg.value = 'Pilih posisi yang dilamar.'
    return
  }
  if (!form.value.last_education) {
    errorMsg.value = 'Pilih pendidikan terakhir.'
    return
  }

  loading.value = true
  errorMsg.value = ''

  const formData = new FormData()
  formData.append('name', form.value.name)
  formData.append('email', form.value.email)
  formData.append('age', form.value.age.toString())
  formData.append('whatsapp_number', form.value.whatsapp_number)
  formData.append('applied_position', form.value.position)
  formData.append('expected_salary', form.value.expected_salary)
  formData.append('address', form.value.address)
  formData.append('last_education', form.value.last_education)
  formData.append('cv', cvFile.value)

  try {
    // Send public multipart form without authentication token needed
    await client.post('/auth/apply', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    submitted.value = true
  } catch (err: any) {
    errorMsg.value = err.response?.data?.error || 'Terjadi kesalahan saat mengirim lamaran. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}
</script>
