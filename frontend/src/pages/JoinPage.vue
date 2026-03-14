<template>
  <div class="min-h-screen bg-slate-50 flex flex-col items-center justify-center p-4 py-8 md:py-12 font-sans overflow-y-auto">
    <!-- Header Logo -->
    <div class="w-full max-w-5xl mb-6 flex items-center gap-2">
      <div class="text-primary font-bold text-xl flex items-center gap-2">
        <div class="w-36 h-36 rounded-full flex items-center justify-center">
          <img :src="viteLogo" alt="Logo" class="w-16 h-16" />
        </div>
        <span class="text-slate-800">HRD Room</span>
      </div>
    </div>

    <div class="w-full max-w-5xl grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Main Registration Form -->
      <Card class="col-span-1 md:col-span-2 shadow-sm border-0 rounded-2xl bg-white p-2">
        <CardHeader class="space-y-2 pb-6">
          <CardTitle class="text-3xl font-extrabold tracking-tight text-slate-900">Laman Ujian</CardTitle>
          <CardDescription class="text-base text-slate-500">
            Silakan masukkan token ujian Anda untuk memulai Test.
          </CardDescription>
        </CardHeader>
        
        <CardContent class="space-y-8">
          <!-- Activation Code Box -->
          <div class="bg-blue-50/50 rounded-xl p-5 border border-blue-100">
            <Label for="token" class="text-xs font-bold text-slate-700 tracking-wider mb-2 block uppercase">
              Token Ujian <span class="text-red-500">*</span>
            </Label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"></path></svg>
              </div>
              <Input 
                id="token" 
                v-model="form.token" 
                placeholder="Masukkan kode (contoh. HPV6cfzDQUe...)" 
                class="pl-10 h-12 border-blue-200 focus-visible:ring-blue-500 bg-white"
                autocomplete="off"
              />
            </div>
            <p class="text-xs text-slate-500 mt-2 flex items-center gap-1">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
              Disediakan oleh rekruter HR Anda dalam email undangan.
            </p>
          </div>



          <Alert v-if="errorMsg" variant="destructive" class="bg-red-50 border-red-200">
            <AlertTitle>Gagal</AlertTitle>
            <AlertDescription>{{ errorMsg }}</AlertDescription>
          </Alert>

          <Button 
            class="w-full h-14 text-base font-semibold bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center justify-center gap-2" 
            @click="handleJoin" 
            :disabled="loading || !isFormValid"
          >
            {{ loading ? 'Memverifikasi...' : 'Mulai Test Sekarang' }}
            <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
          </Button>

          <div class="pt-6 mt-2 border-t border-slate-100">
            <p class="text-sm text-center text-slate-500 mb-4 font-medium">Belum memiliki token atau belum melamar?</p>
            <Button 
              variant="outline"
              class="w-full h-12 text-sm border-slate-200 text-slate-700 hover:bg-slate-50 hover:text-blue-600 hover:border-blue-200 rounded-xl transition-all flex items-center justify-center gap-2 font-semibold" 
              @click="router.push('/apply')" 
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M22 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
              Lamar Pekerjaan Sekarang
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Sidebar -->
      <div class="col-span-1 space-y-6">
        <Card class="shadow-sm border-0 rounded-2xl bg-white p-2">
          <CardContent class="p-4 md:p-6 space-y-6">
            <div class="space-y-5">
              <h3 class="font-bold text-slate-800 flex items-center gap-2 mb-2">
                <div class="w-6 h-6 rounded-full bg-green-100 flex items-center justify-center text-green-600">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                </div>
                Aturan Test
              </h3>
              
              <div class="flex gap-3 items-start">
                <div class="mt-0.5 rounded-full bg-green-100 p-0.5 text-green-600 shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </div>
                <div>
                  <p class="font-semibold text-sm text-slate-800 text-left">Internet Stabil</p>
                  <p class="text-xs text-slate-500 text-left">Pastikan Internet Anda Stabil</p>
                </div>
              </div>

              <div class="flex gap-3 items-start">
                <div class="mt-0.5 rounded-full bg-green-100 p-0.5 text-green-600 shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </div>
                <div>
                  <p class="font-semibold text-sm text-slate-800 text-left">Aktifkan Webcam</p>
                  <p class="text-xs text-slate-500 text-left">Kamera Anda harus tetap menyala untuk tujuan pengawasan.</p>
                </div>
              </div>

              <div class="flex gap-3 items-start">
                <div class="mt-0.5 rounded-full bg-green-100 p-0.5 text-green-600 shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </div>
                <div>
                  <p class="font-semibold text-sm text-slate-800 text-left">Dilarang Browsing Jawaban</p>
                  <p class="text-xs text-slate-500 text-left">Keluar dari tab penilaian dapat mengakibatkan diskualifikasi.</p>
                </div>
              </div>

              <div class="flex gap-3 items-start">
                <div class="mt-0.5 rounded-full bg-green-100 p-0.5 text-green-600 shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </div>
                <div>
                  <p class="font-semibold text-sm text-slate-800 text-left">Pilih Tempat yang Tenang</p>
                </div>
              </div>
            </div>

            <div class="pt-6 mt-6 border-t border-slate-100 flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center shrink-0 overflow-hidden">
                <img src="https://api.dicebear.com/7.x/notionists/svg?seed=Support" alt="Support" class="w-full h-full object-cover" />
              </div>
              <div class="text-left">
                <p class="text-xs text-slate-500">Butuh bantuan?</p>
                <RouterLink to="/contact-support" class="text-sm font-semibold text-blue-600 hover:underline">Contact Support</RouterLink>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
    
    <!-- Footer layout -->
    <div class="w-full max-w-5xl mt-12 flex flex-col md:flex-row items-center justify-between text-xs text-slate-500 pb-2">
      <div class="mb-4 md:mb-0 text-center md:text-left">
        <p>© 2026 PT Wowin Purnomo Putera.</p>
        <p>All rights reserved.</p>
      </div>
      <div class="flex justify-center md:justify-end gap-6 font-medium">
        <RouterLink to="/privacy-policy" class="hover:text-slate-800">Privacy Policy</RouterLink>
        <RouterLink to="/terms-of-service" class="hover:text-slate-800">Terms of Service</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useExamStore } from '@/stores/exam'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import viteLogo from '@/assets/vite.svg'

const router = useRouter()
const route = useRoute()
const examStore = useExamStore()

const frictionlessToken = computed(() => route.query.token as string | undefined)

const form = ref({
  token: '',
  name: '',
  age: '',
  email: '',
  position: ''
})

const errorMsg = ref('')
const loading = ref(false)

onMounted(async () => {
  // If token is in URL (Magic Link), try to join immediately
  if (frictionlessToken.value) {
    form.value.token = frictionlessToken.value
    loading.value = true
    try {
      await examStore.joinExam({
        token: frictionlessToken.value,
        name: '', email: '', age: 0, position: ''
      })
      router.push('/camera-check')
      return // Success! Skip loading job positions
    } catch (err: any) {
      errorMsg.value = err.response?.data?.error || 'Tautan Anda tidak valid atau token sudah kadaluarsa. Silakan isi form.'
      loading.value = false
    }
  }
})

const isFormValid = computed(() => {
  return form.value.token.trim() !== ''
})

const handleJoin = async () => {
  if (!isFormValid.value && !frictionlessToken.value) return
  
  errorMsg.value = ''
  loading.value = true
  try {
    const payload = {
      token: form.value.token.trim(),
      name: '',
      email: '',
      age: 0,
      position: ''
    }
    await examStore.joinExam(payload)
    router.push('/camera-check')
  } catch (err: any) {
    errorMsg.value = err.response?.data?.error || 'Gagal masuk sesi ujian. Periksa kembali data atau token Anda.'
  } finally {
    loading.value = false
  }
}
</script>
