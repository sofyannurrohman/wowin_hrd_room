<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <Card class="w-full max-w-lg text-center shadow-lg border-t-4 border-t-green-500">
      <CardContent class="pt-10 pb-8 space-y-6 flex flex-col items-center">
        <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mb-4">
          <CheckCircle2Icon class="w-10 h-10 text-green-600" />
        </div>
        <h2 class="text-3xl font-bold text-slate-800">Ujian Selesai!</h2>
        <p class="text-slate-500 max-w-sm">
          Jawaban Anda telah berhasil dikirimkan dan tersimpan di sistem kami. Anda dapat menutup halaman ini sekarang.
        </p>
        <Button variant="outline" class="w-full mt-6" @click="backHome">Kembali ke Halaman Awal</Button>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { CheckCircle2Icon } from 'lucide-vue-next'

const router = useRouter()

const backHome = () => {
  router.push('/')
}
</script>
