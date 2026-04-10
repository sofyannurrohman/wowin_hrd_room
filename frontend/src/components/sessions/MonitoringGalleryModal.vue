<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-8 bg-slate-900/90 backdrop-blur-md" @click.self="emit('close')">
      <div class="bg-white rounded-3xl shadow-2xl w-full max-w-5xl flex flex-col max-h-[90vh] overflow-hidden animate-in zoom-in-95 duration-200">
        <!-- Header -->
        <div class="px-8 py-6 border-b border-slate-100 flex items-center justify-between shrink-0 bg-white">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-blue-100 flex items-center justify-center text-blue-600">
              <CameraIcon class="w-6 h-6" />
            </div>
            <div>
              <h3 class="text-xl font-extrabold text-slate-800 tracking-tight">Monitoring Snapshot: {{ participantName }}</h3>
              <p class="text-sm text-slate-500 font-medium">Koleksi foto berkala setiap 5 menit & bukti pelanggaran.</p>
            </div>
          </div>
          <button @click="emit('close')" class="p-2 hover:bg-slate-100 rounded-full transition-colors">
            <XIcon class="w-6 h-6 text-slate-400" />
          </button>
        </div>

        <!-- content -->
        <div class="flex-1 overflow-y-auto p-8 bg-slate-50/50">
          <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-4">
            <div class="w-10 h-10 border-4 border-blue-100 border-t-blue-600 rounded-full animate-spin"></div>
            <p class="text-slate-400 font-bold text-sm">Memuat snapshot...</p>
          </div>

          <div v-else-if="photos.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
            <div class="w-20 h-20 bg-slate-100 rounded-full flex items-center justify-center mb-4">
              <CameraOffIcon class="w-10 h-10 text-slate-300" />
            </div>
            <h4 class="text-xl font-bold text-slate-800">Snapshot Tidak Ditemukan</h4>
            <p class="text-slate-500 max-w-sm mt-2 leading-relaxed">
              Sistem tidak menemukan snapshot berkala untuk peserta ini. 
              <br/>
              <span class="text-red-500 font-bold mt-2 block p-3 bg-red-50 rounded-xl border border-red-100">
                 Foto telah dihapus dari sistem setelah 7 hari untuk efisiensi penyimpanan.
              </span>
            </p>
          </div>

          <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            <div 
              v-for="photo in photos" 
              :key="photo.id" 
              class="group relative aspect-square bg-slate-200 rounded-2xl overflow-hidden border-2 border-transparent hover:border-blue-500 transition-all cursor-zoom-in shadow-sm"
              @click="selectedPhoto = photo.photo_url"
            >
              <img :src="photo.photo_url" class="w-full h-full object-cover transition-transform group-hover:scale-110 duration-500" />
              <div class="absolute inset-x-0 bottom-0 p-3 bg-black/60 backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity">
                <p class="text-[10px] font-bold text-white tracking-widest">{{ formatTimestamp(photo.created_at) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-8 py-4 border-t border-slate-100 bg-white flex justify-end shrink-0">
          <Button variant="outline" class="rounded-xl font-bold" @click="emit('close')">Tutup Gallery</Button>
        </div>
      </div>
    </div>

    <!-- Inner Zoom Modal -->
    <div 
      v-if="selectedPhoto" 
      class="fixed inset-0 z-[110] bg-black/95 flex items-center justify-center p-4 cursor-zoom-out"
      @click="selectedPhoto = null"
    >
      <img :src="selectedPhoto" class="max-w-full max-h-full object-contain rounded-lg shadow-2xl" />
      <button class="absolute top-6 right-6 p-3 bg-white/10 hover:bg-white/20 rounded-full text-white backdrop-blur-md transition-all">
        <XIcon class="w-6 h-6" />
      </button>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import client from '@/api/client'
import { Button } from '@/components/ui/button'
import { CameraIcon, CameraOffIcon, XIcon, ImageIcon } from 'lucide-vue-next'

const props = defineProps<{
  isOpen: boolean
  participantId: string
  participantName: string
  sessionId: string
}>()

const emit = defineEmits(['close'])

const photos = ref<any[]>([])
const loading = ref(false)
const selectedPhoto = ref<string | null>(null)

const fetchPhotos = async () => {
  if (!props.participantId || !props.sessionId) return
  loading.value = true
  try {
    const res = await client.get(`/sessions/${props.sessionId}/participants/${props.participantId}/monitoring`)
    photos.value = res.data || []
  } catch (err) {
    console.error('Failed to fetch monitoring photos', err)
  } finally {
    loading.value = false
  }
}

watch(() => props.isOpen, (val) => {
  if (val) fetchPhotos()
})

const formatTimestamp = (ts: string) => {
  return new Date(ts).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
</script>
