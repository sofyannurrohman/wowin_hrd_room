<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4 py-6 font-sans">
    <div class="w-full max-w-5xl flex flex-col gap-6">
      <!-- Top header (Smaller) -->
      <div class="flex items-center justify-between px-2">
        <div class="flex items-center gap-2">
          <img :src="viteLogo" alt="Logo" class="w-10 h-10" />
          <span class="font-black text-xl text-slate-800 tracking-tight">Assesment Center</span>
        </div>
        <div class="hidden md:flex items-center gap-1 text-[11px] font-bold text-slate-400 uppercase tracking-widest">
          <div class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
          System Live
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
        <!-- Main Form (Left) -->
        <Card class="md:col-span-7 shadow-xl border-0 rounded-3xl bg-white overflow-hidden">
          <div class="bg-blue-600 h-2 w-full"></div>
          <CardHeader class="pb-4">
            <CardTitle class="text-2xl font-black text-slate-900">Laman Ujian</CardTitle>
            <CardDescription class="text-sm">
              Silakan masukkan token ujian Anda untuk memulai Test.
            </CardDescription>
          </CardHeader>
          
          <CardContent class="space-y-6">
            <div class="p-5 bg-blue-50/50 rounded-2xl border border-blue-100/50">
              <Label for="token" class="text-[10px] font-black text-blue-700 tracking-widest mb-2 block uppercase">
                Token Ujian <span class="text-red-500">*</span>
              </Label>
              <div class="relative">
                <Input 
                  id="token" 
                  v-model="form.token" 
                  placeholder="Contoh: HPV6cfzDQUe..." 
                  class="h-12 border-blue-200 focus-visible:ring-blue-500 bg-white rounded-xl"
                  autocomplete="off"
                />
              </div>
              <p class="text-[10px] text-slate-400 mt-2 italic font-medium">Token dapat dilihat pada email undangan.</p>
            </div>

            <Alert v-if="errorMsg" variant="destructive" class="bg-red-50 border-red-100 text-red-700 rounded-xl py-3 px-4">
              <AlertDescription class="text-xs font-semibold">{{ errorMsg }}</AlertDescription>
            </Alert>

            <Button 
              class="w-full h-14 text-lg font-black bg-blue-600 hover:bg-blue-700 text-white rounded-2xl transition-all shadow-lg shadow-blue-600/20 active:scale-95 flex items-center justify-center gap-2" 
              @click="handleJoin" 
              :disabled="loading || !isFormValid"
            >
              {{ loading ? 'Memverifikasi...' : 'Mulai Sekarang' }}
              <ArrowRightIcon v-if="!loading" class="w-5 h-5" />
            </Button>

            <div class="pt-4 border-t border-slate-50">
              <button 
                @click="router.push('/career')" 
                class="w-full py-3 text-xs font-bold text-slate-500 hover:text-blue-600 transition-colors flex items-center justify-center gap-2"
              >
                <BriefcaseIcon class="w-4 h-4" />
                Lihat Peluang Karir & Lowongan
              </button>
            </div>
          </CardContent>
        </Card>

        <!-- Side Content (Right) -->
        <div class="md:col-span-5 space-y-6">
          <!-- Compact Info Card -->
          <Card class="shadow-sm border-0 rounded-3xl bg-white p-1">
            <CardContent class="p-5 space-y-5">
              <div class="flex items-center gap-3 border-b border-slate-50 pb-3">
                 <div class="w-8 h-8 rounded-lg bg-green-100 flex items-center justify-center text-green-600">
                    <ShieldCheckIcon class="w-5 h-5" />
                 </div>
                 <h3 class="font-black text-slate-800 text-sm italic tracking-tight">Aturan & Informasi</h3>
              </div>

              <!-- Compact Rules -->
              <div class="grid grid-cols-2 gap-3">
                <div class="p-3 bg-slate-50 rounded-2xl border border-slate-100 flex flex-col gap-1">
                   <div class="text-green-600 mb-1"><WifiIcon class="w-4 h-4" /></div>
                   <p class="font-bold text-[10px] text-slate-700">Internet Stabil</p>
                </div>
                <div class="p-3 bg-slate-50 rounded-2xl border border-slate-100 flex flex-col gap-1">
                   <div class="text-blue-600 mb-1"><VideoIcon class="w-4 h-4" /></div>
                   <p class="font-bold text-[10px] text-slate-700">Webcam Aktif</p>
                </div>
                <div class="p-3 bg-slate-50 rounded-2xl border border-slate-100 flex flex-col gap-1">
                   <div class="text-orange-600 mb-1"><AlertCircleIcon class="w-4 h-4" /></div>
                   <p class="font-bold text-[10px] text-slate-700">Dilarang Browsing</p>
                </div>
                <div class="p-3 bg-slate-50 rounded-2xl border border-slate-100 flex flex-col gap-1">
                   <div class="text-purple-600 mb-1"><MoonIcon class="w-4 h-4" /></div>
                   <p class="font-bold text-[10px] text-slate-700">Tempat Tenang</p>
                </div>
              </div>

              <div class="space-y-3 pt-2">
                <!-- IST -->
                <div class="p-3 rounded-2xl border border-blue-50 bg-blue-50/20 group hover:bg-blue-50 transition-colors">
                   <div class="flex items-center gap-2 mb-1.5">
                      <div class="w-1.5 h-1.5 rounded-full bg-blue-500"></div>
                      <span class="text-[11px] font-black text-blue-700 uppercase tracking-tight">IST Assessment</span>
                   </div>
                   <p class="text-[10px] text-slate-500 leading-relaxed px-1">
                     Intelligence Structure Test mengukur struktur intelegensi dan potensi akademik dalam berbagai aspek kognitif.
                   </p>
                </div>

                <!-- DISC -->
                <div class="p-3 rounded-2xl border border-purple-50 bg-purple-50/10 group hover:bg-purple-50 transition-colors">
                   <div class="flex items-center gap-2 mb-1.5">
                      <div class="w-1.5 h-1.5 rounded-full bg-purple-500"></div>
                      <span class="text-[11px] font-black text-purple-700 uppercase tracking-tight">DISC Assessment</span>
                   </div>
                   <p class="text-[10px] text-slate-500 leading-relaxed px-1">
                     Memahami gaya perilaku dan kecenderungan interaksi Anda di lingkungan kerja profesional.
                   </p>
                </div>

                <!-- Skill -->
                <div class="p-3 rounded-2xl border border-emerald-50 bg-emerald-50/10 group hover:bg-emerald-50 transition-colors">
                   <div class="flex items-center gap-2 mb-1.5">
                      <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
                      <span class="text-[11px] font-black text-emerald-700 uppercase tracking-tight">Skill Based Hiring</span>
                   </div>
                   <p class="text-[10px] text-slate-500 leading-relaxed px-1">
                     Evaluasi kemampuan teknis dan keterampilan spesifik yang relevan dengan posisi yang dilamar.
                   </p>
                </div>
              </div>

              <div class="p-3 bg-amber-50/50 rounded-2xl border border-amber-100 flex gap-2 items-start mt-2">
                 <AlertCircleIcon class="w-4 h-4 text-amber-600 shrink-0" />
                 <p class="text-[10px] text-amber-800 leading-relaxed font-medium">
                   Pastikan Anda memiliki waktu luang minimal 60-120 menit untuk menyelesaikan seluruh rangkaian tes.
                 </p>
              </div>

              <div class="pt-4 mt-2 border-t border-slate-50 flex items-center justify-between">
                <div class="flex items-center gap-2">
                   <img src="https://api.dicebear.com/7.x/notionists/svg?seed=Support" alt="Support" class="w-8 h-8 rounded-full bg-orange-50" />
                   <div class="text-left">
                     <p class="text-[10px] text-slate-400 font-medium">Bantuan?</p>
                     <RouterLink to="/contact-support" class="text-[11px] font-bold text-blue-600">Hubungi Kami</RouterLink>
                   </div>
                </div>
                <div class="text-[10px] text-slate-300 font-mono">v2.4.0</div>
              </div>
            </CardContent>
          </Card>

      </div>
    </div>

    <!-- Default Footer -->
      <div class="w-full mt-12 flex flex-col md:flex-row items-center justify-between text-xs text-slate-500 pb-10 border-t border-slate-100 pt-8 px-2">
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
import { Alert, AlertDescription } from '@/components/ui/alert'
import viteLogo from '@/assets/vite.svg'
import { 
  ArrowRightIcon, 
  BriefcaseIcon, 
  ShieldCheckIcon, 
  WifiIcon, 
  VideoIcon, 
  AlertCircleIcon, 
  MoonIcon 
} from 'lucide-vue-next'

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
