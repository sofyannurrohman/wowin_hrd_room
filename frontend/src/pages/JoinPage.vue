<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <Card class="w-full max-w-md shadow-lg border-t-4 border-t-primary">
      <CardHeader class="space-y-1 text-center pb-8 pt-6">
        <CardTitle class="text-3xl font-bold tracking-tight text-primary">HRD Room</CardTitle>
        <CardDescription class="text-base mt-2">Masuk ke Sesi Ujian Anda</CardDescription>
      </CardHeader>
      <CardContent class="space-y-6">
        <div class="space-y-3">
          <Label for="token" class="text-sm font-semibold text-slate-700">Exam Token</Label>
          <Input 
            id="token" 
            v-model="token" 
            placeholder="Masukkan token 48 karakter" 
            class="h-12 text-center tracking-widest font-mono bg-slate-100"
            autocomplete="off"
            @keyup.enter="handleJoin"
          />
        </div>
        
        <Alert v-if="errorMsg" variant="destructive" class="animate-in fade-in slide-in-from-top-2">
          <AlertTitle>Akses Ditolak</AlertTitle>
          <AlertDescription>{{ errorMsg }}</AlertDescription>
        </Alert>

        <Button class="w-full h-12 text-lg font-medium transition-all hover:scale-[1.02]" @click="handleJoin" :disabled="loading || !token">
          {{ loading ? 'Memverifikasi...' : 'Mulai Ujian' }}
        </Button>
      </CardContent>
      <div class="px-6 pb-6 text-center text-xs text-slate-500">
        Pastikan Anda memiliki koneksi internet yang stabil dan webcam yang berfungsi dengan baik sebelum memulai.
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useExamStore } from '@/stores/exam'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const router = useRouter()
const examStore = useExamStore()

const token = ref('')
const errorMsg = ref('')
const loading = ref(false)

const handleJoin = async () => {
  if (!token.value.trim()) return
  
  errorMsg.value = ''
  loading.value = true
  try {
    // Attempt to join using the token
    await examStore.joinExam(token.value.trim())
    // Success, proceed to camera check
    router.push('/camera-check')
  } catch (err: any) {
    errorMsg.value = err.response?.data?.error || 'Gagal masuk sesi ujian. Token tidak valid.'
  } finally {
    loading.value = false
  }
}
</script>
