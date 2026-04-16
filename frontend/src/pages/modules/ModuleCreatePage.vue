<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div class="flex items-center gap-4">
      <router-link to="/modules" class="p-2 -ml-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100 transition-colors">
        <ArrowLeftIcon class="w-5 h-5" />
      </router-link>
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900">Buat Modul Baru</h1>
        <p class="text-sm text-slate-500 mt-1">Tambah bank soal kedalam modul </p>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      <form @submit.prevent="submit" class="p-6 md:p-8 space-y-6">
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">
              Nama Modul <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors outline-none"
              placeholder="e.g. Frontend Evaluation"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">
              Total Bobot <span class="text-red-500">*</span>
            </label>
            <input
              v-model.number="form.total_weight"
              type="number"
              min="1"
              step="0.01"
              required
              class="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors outline-none"
              placeholder="100"
            />
            <p class="text-xs text-slate-500 mt-1">Total maksimal poin untuk semua soal di dalam modul ini.</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">Deskripsi</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors outline-none resize-y text-sm"
              placeholder="Module description (optional)..."
            ></textarea>
          </div>

          <div class="pt-4 border-t border-slate-100">
            <h3 class="text-sm font-bold text-slate-900 mb-4">Fitur Menghafal (Memorization)</h3>
            
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Durasi Menghafal (Detik)
                </label>
                <input
                  v-model.number="form.memorization_duration"
                  type="number"
                  min="0"
                  class="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors outline-none text-sm"
                  placeholder="e.g. 180 (untuk 3 menit)"
                />
                <p class="text-xs text-slate-500 mt-1">Isi 0 jika tidak memerlukan fase menghafal.</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">
                  Konten yang Harus Diingat
                </label>
                <textarea
                  v-model="form.memorization_content"
                  rows="6"
                  class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors outline-none resize-y text-sm font-mono"
                  placeholder="Masukkan kata-kata atau teks yang harus diingat peserta..."
                ></textarea>
                <p class="text-xs text-slate-500 mt-1">Teks ini akan ditampilkan kepada peserta selama durasi yang ditentukan di atas.</p>
              </div>
            </div>
          </div>
        </div>

        <div class="pt-6 border-t border-slate-100 flex items-center justify-end gap-3">
          <button
            type="button"
            @click="router.push('/modules')"
            class="px-5 py-2.5 text-sm font-medium text-slate-600 hover:bg-slate-50 border border-slate-200 rounded-lg transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="flex items-center gap-2 px-5 py-2.5 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <SaveIcon v-if="!saving" class="w-4 h-4" />
            <Loader2Icon v-else class="w-4 h-4 animate-spin" />
            {{ saving ? 'Saving...' : 'Create Module' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeftIcon, SaveIcon, Loader2Icon } from 'lucide-vue-next'
import client from '@/api/client'
import { toast } from 'vue-sonner'

const router = useRouter()
const saving = ref(false)

const form = ref({
  name: '',
  description: '',
  total_weight: 100,
  memorization_content: '',
  memorization_duration: 0
})

const submit = async () => {
  saving.value = true
  try {
    await client.post('/modules', form.value)
    toast.success('Module successfully created')
    router.push('/modules')
  } catch (error) {
    toast.error('Failed to create module')
    console.error(error)
  } finally {
    saving.value = false
  }
}
</script>
