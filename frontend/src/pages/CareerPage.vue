<template>
  <div class="min-h-screen bg-slate-50 font-sans">
    <!-- Hero Section -->
    <div class="bg-blue-600 relative overflow-hidden py-16 md:py-24">
      <div class="absolute inset-0 bg-gradient-to-br from-blue-700 to-indigo-900 opacity-90"></div>
      <div class="absolute -top-24 -right-24 w-96 h-96 bg-white/10 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-blue-400/20 rounded-full blur-3xl"></div>
      
      <div class="max-w-6xl mx-auto px-4 relative z-10">
        <div class="flex flex-col items-center text-center space-y-6">
          <div class="w-20 h-20 bg-white/20 backdrop-blur-md rounded-3xl flex items-center justify-center mb-2 animate-bounce-slow">
            <img :src="viteLogo" alt="Logo" class="w-12 h-12" />
          </div>
          <h1 class="text-4xl md:text-6xl font-black text-white tracking-tight leading-tight">
            Bergabung Bersama <br /><span class="text-blue-200">PT Wowin Purnomo Putera</span>
          </h1>
          <p class="text-blue-100 text-lg md:text-xl max-w-2xl font-medium opacity-90">
            Temukan peluang karir terbaik dan kembangkan potensi Anda bersama tim inovatif kami.
          </p>
          <div class="flex gap-4 pt-4">
            <a href="#openings" class="bg-white text-blue-700 px-8 py-4 rounded-2xl font-bold hover:bg-blue-50 transition-all shadow-xl shadow-blue-900/20 active:scale-95">
              Lihat Lowongan
            </a>
            <RouterLink to="/join" class="bg-blue-500/30 backdrop-blur-md text-white border border-white/20 px-8 py-4 rounded-2xl font-bold hover:bg-white/10 transition-all active:scale-95">
              Masuk Ujian
            </RouterLink>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div id="openings" class="max-w-6xl mx-auto px-4 py-16 md:py-24 -mt-10 relative z-20">
      <div v-if="loading" class="flex flex-col items-center justify-center py-20 space-y-4">
        <div class="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-slate-500 font-medium">Memuat lowongan tersedia...</p>
      </div>

      <div v-else-if="positions.length === 0" class="bg-white rounded-3xl p-12 text-center shadow-sm border border-slate-100">
        <div class="w-20 h-20 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="8" y1="12" x2="16" y2="12"></line></svg>
        </div>
        <h3 class="text-xl font-bold text-slate-800">Belum Ada Lowongan Aktif</h3>
        <p class="text-slate-500 mt-2">Silakan cek kembali di lain waktu atau hubungi HR kami.</p>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Sidebar Navigation -->
        <div class="lg:col-span-1 space-y-4">
          <div class="bg-white rounded-3xl p-6 shadow-sm border border-slate-100 sticky top-8">
            <h3 class="text-sm font-bold text-slate-400 uppercase tracking-widest mb-6">Pilih Posisi</h3>
            <div class="space-y-2">
              <button 
                v-for="pos in positions" 
                :key="pos.id"
                @click="selectedPosition = pos"
                class="w-full text-left p-4 rounded-2xl transition-all duration-300 group flex items-center justify-between"
                :class="selectedPosition?.id === pos.id ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'hover:bg-slate-50 text-slate-700'"
              >
                <span class="font-bold truncate">{{ pos.name }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 opacity-0 group-hover:opacity-100 transition-opacity" :class="{'opacity-100': selectedPosition?.id === pos.id}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Job Details -->
        <div class="lg:col-span-2">
          <transition name="fade-slide" mode="out-in">
            <div :key="selectedPosition?.id" class="bg-white rounded-4xl shadow-sm border border-slate-100 overflow-hidden">
              <div class="p-8 md:p-12 space-y-10">
                <!-- Header -->
                <div class="flex flex-col md:flex-row md:items-center justify-between gap-6 border-b border-slate-100 pb-10">
                  <div class="space-y-2">
                    <div class="inline-flex items-center px-3 py-1 rounded-full bg-emerald-50 text-emerald-600 text-xs font-bold uppercase tracking-wider">
                      Full-time • Open
                    </div>
                    <h2 class="text-3xl md:text-4xl font-black text-slate-800 tracking-tight">{{ selectedPosition?.name }}</h2>
                  </div>
                  <Button 
                    @click="applyNow(selectedPosition)"
                    class="h-16 px-10 rounded-2xl bg-blue-600 hover:bg-blue-700 text-white font-bold text-lg shadow-xl shadow-blue-600/20 transition-all hover:scale-105 active:scale-95"
                  >
                    Lamar Sekarang
                  </Button>
                </div>

                <!-- Body -->
                <div class="grid grid-cols-1 md:grid-cols-1 gap-10">
                  <div class="space-y-4">
                    <h4 class="text-lg font-bold text-slate-800 flex items-center gap-3">
                      <span class="w-8 h-8 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center text-sm">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="17" y1="21" x2="7" y2="21"></line><polyline points="7 3 12 8 17 3"></polyline><line x1="12" y1="8" x2="12" y2="21"></line></svg>
                      </span>
                      Tentang Pekerjaan
                    </h4>
                    <p class="text-slate-600 leading-relaxed whitespace-pre-line text-lg">
                      {{ selectedPosition?.description || 'Tidak ada deskripsi tersedia untuk posisi ini.' }}
                    </p>
                  </div>

                  <div class="space-y-4">
                    <h4 class="text-lg font-bold text-slate-800 flex items-center gap-3">
                      <span class="w-8 h-8 rounded-xl bg-orange-50 text-orange-600 flex items-center justify-center text-sm">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 11 12 14 22 4"></polyline><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path></svg>
                      </span>
                      Kualifikasi
                    </h4>
                    <div class="bg-slate-50 rounded-3xl p-8 border border-slate-100">
                      <p class="text-slate-600 leading-relaxed whitespace-pre-line">
                        {{ selectedPosition?.requirements || 'Hubungi HR kami untuk informasi lebih lanjut mengenai kualifikasi.' }}
                      </p>
                    </div>
                  </div>
                </div>

                <!-- Footer Action -->
                <div class="pt-8 border-t border-slate-100 flex flex-col md:flex-row items-center justify-between gap-6">
                  <p class="text-sm text-slate-400">Dipublikasikan pada {{ formatDate(selectedPosition?.created_at) }}</p>
                  <div class="flex gap-4 items-center">
                    <div class="relative">
                      <button @click="showShareMenu = !showShareMenu" class="p-4 rounded-xl text-slate-400 hover:text-blue-600 hover:bg-blue-50 transition-all" title="Bagikan Lowongan">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="18" cy="5" r="3"></circle><circle cx="6" cy="12" r="3"></circle><circle cx="18" cy="19" r="3"></circle><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line></svg>
                      </button>
                      
                      <div v-if="showShareMenu" class="fixed inset-0 z-40" @click="showShareMenu = false"></div>
                      
                      <transition name="fade-slide">
                        <div v-if="showShareMenu" class="absolute bottom-full right-0 mb-3 w-48 bg-white border border-slate-100 shadow-xl rounded-2xl overflow-hidden z-50">
                          <button @click="shareTo('whatsapp', selectedPosition)" class="w-full text-left px-4 py-3 hover:bg-emerald-50 text-slate-700 hover:text-emerald-700 text-sm font-medium flex items-center gap-3 transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M13.601 2.326A7.85 7.85 0 0 0 7.994 0C3.627 0 .068 3.558.064 7.926c0 1.399.366 2.76 1.057 3.965L0 16l4.204-1.102a7.9 7.9 0 0 0 3.79.965h.004c4.368 0 7.926-3.558 7.93-7.93A7.9 7.9 0 0 0 13.6 2.326zM7.994 14.521a6.6 6.6 0 0 1-3.356-.92l-.24-.144-2.494.654.666-2.433-.156-.251a6.56 6.56 0 0 1-1.007-3.505c0-3.626 2.957-6.584 6.591-6.584a6.56 6.56 0 0 1 4.66 1.931 6.56 6.56 0 0 1 1.928 4.66c-.004 3.639-2.961 6.592-6.592 6.592m3.615-4.934c-.197-.099-1.17-.578-1.353-.646-.182-.065-.315-.099-.445.099-.133.197-.513.646-.627.775-.114.133-.232.148-.43.05-.197-.1-.836-.308-1.592-.985-.59-.525-.985-1.175-1.103-1.372-.114-.198-.011-.304.088-.403.087-.088.197-.232.296-.346.1-.114.133-.198.198-.33.065-.134.034-.248-.015-.347-.05-.099-.445-1.076-.612-1.47-.16-.389-.323-.335-.445-.34-.114-.007-.247-.007-.38-.007a.73.73 0 0 0-.529.247c-.182.198-.691.677-.691 1.654s.71 1.916.81 2.049c.098.133 1.394 2.132 3.383 2.992.47.205.84.326 1.129.418.475.152.904.129 1.246.08.38-.058 1.171-.48 1.338-.943.164-.464.164-.86.114-.943-.049-.084-.182-.133-.38-.232"/></svg>
                            WhatsApp
                          </button>
                          <button @click="shareTo('linkedin', selectedPosition)" class="w-full text-left px-4 py-3 hover:bg-blue-50 text-slate-700 hover:text-blue-700 text-sm font-medium flex items-center gap-3 transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M0 1.146C0 .513.526 0 1.175 0h13.65C15.474 0 16 .513 16 1.146v13.708c0 .633-.526 1.146-1.175 1.146H1.175C.526 16 0 15.487 0 14.854zm4.943 12.248V6.169H2.542v7.225zm-1.2-8.212c.837 0 1.358-.554 1.358-1.248-.015-.709-.52-1.248-1.342-1.248S1.2 3.226 1.2 3.934c0 .694.521 1.248 1.327 1.248zm4.908 8.212V9.359c0-.216.016-.432.08-.586.173-.431.568-.878 1.232-.878.869 0 1.216.662 1.216 1.634v3.865h2.401V9.25c0-2.22-1.184-3.252-2.764-3.252-1.274 0-1.845.7-2.165 1.193v.025h-.016l.016-.025V6.169h-2.4c.03.678 0 7.225 0 7.225z"/></svg>
                            LinkedIn
                          </button>
                          <button @click="copyLink(selectedPosition)" class="w-full text-left px-4 py-3 border-t border-slate-100 hover:bg-slate-50 text-slate-700 text-sm font-medium flex items-center gap-3 transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                            Copy Link
                          </button>
                        </div>
                      </transition>
                    </div>
                    <Button variant="outline" @click="applyNow(selectedPosition)" class="h-14 px-8 rounded-2xl font-bold border-slate-200 text-slate-700 hover:bg-slate-50 shadow-sm transition-all hover:-translate-y-0.5">
                      Lamar Sekarang
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="bg-white border-t border-slate-100 py-12">
      <div class="max-w-6xl mx-auto px-4 flex flex-col md:flex-row items-center justify-between gap-8">
        <div class="flex items-center gap-3">
          <img :src="viteLogo" alt="Logo" class="w-8 h-8 opacity-50 contrast-0" />
          <p class="text-sm text-slate-400 font-medium">© 2026 PT Wowin Purnomo Putera. All rights reserved.</p>
        </div>
        <div class="flex gap-8 text-sm font-bold text-slate-400">
          <RouterLink to="/privacy-policy" class="hover:text-blue-600 transition-colors">Privacy Policy</RouterLink>
          <RouterLink to="/terms-of-service" class="hover:text-blue-600 transition-colors">Terms of Service</RouterLink>
          <RouterLink to="/contact-support" class="hover:text-blue-600 transition-colors">Support</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import client from '@/api/client'
import { Button } from '@/components/ui/button'
import viteLogo from '@/assets/vite.svg'
import { toast } from 'vue-sonner'

const router = useRouter()
const loading = ref(true)
const positions = ref<any[]>([])
const selectedPosition = ref<any>(null)
const showShareMenu = ref(false)

onMounted(async () => {
  try {
    const res = await client.get('/job-positions/active')
    positions.value = res.data.positions || []
    if (positions.value.length > 0) {
      selectedPosition.value = positions.value[0]
    }
  } catch (err) {
    console.error('Failed to load job positions', err)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const applyNow = (position: any) => {
  router.push({
    path: '/apply',
    query: { position: position.name }
  })
}

const shareTo = (platform: string, position: any) => {
  if (!position) return
  showShareMenu.value = false
  
  const url = encodeURIComponent(window.location.href)
  const title = encodeURIComponent(`Lowongan: ${position.name} di PT Wowin Purnomo Putera`)
  const text = encodeURIComponent(`Lihat lowongan ${position.name} yang sedang dibuka!`)

  let shareUrl = ''
  if (platform === 'whatsapp') {
    shareUrl = `https://api.whatsapp.com/send?text=${title}%0A${text}%0A${url}`
  } else if (platform === 'linkedin') {
    shareUrl = `https://www.linkedin.com/sharing/share-offsite/?url=${url}`
  }
  
  if (shareUrl) {
    window.open(shareUrl, '_blank', 'noopener,noreferrer')
  }
}

const copyLink = async (position: any) => {
  if (!position) return
  showShareMenu.value = false
  
  const url = window.location.href
  const title = `Lowongan: ${position.name} di PT Wowin Purnomo Putera`
  const text = `Lihat lowongan ${position.name} yang sedang dibuka!`

  try {
    if (navigator.share) {
      await navigator.share({ title, text, url })
    } else {
      await navigator.clipboard.writeText(`${title}\n${text}\n${url}`)
      toast.success('Info lowongan berhasil disalin ke clipboard')
    }
  } catch (err) {
    if ((err as Error).name !== 'AbortError') {
      toast.error('Gagal membagikan lowongan')
    }
  }
}
</script>

<style scoped>
.animate-bounce-slow {
  animation: bounce 3s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

.rounded-4xl {
  border-radius: 2.5rem;
}
</style>
