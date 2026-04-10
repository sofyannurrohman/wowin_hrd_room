<template>
  <div class="space-y-6 w-full max-w-[1200px] mx-auto pb-10">
    <!-- Header/Breadcrumb -->
    <div class="flex items-center gap-4">
      <Button variant="ghost" size="icon" @click="router.back()" class="rounded-full h-10 w-10">
        <ArrowLeftIcon class="w-5 h-5" />
      </Button>
      <div>
        <div class="flex items-center gap-2 text-sm text-slate-500 mb-1">
          <RouterLink to="/modules" class="hover:text-blue-600">Module Soal</RouterLink>
          <ChevronRightIcon class="w-4 h-4" />
          <span>Detail Modul</span>
        </div>
        <h1 class="text-3xl font-extrabold tracking-tight text-slate-900" v-if="module">{{ module.name }}</h1>
        <div v-else class="h-9 w-64 bg-slate-100 animate-pulse rounded"></div>
      </div>
    </div>

    <div v-if="loading" class="space-y-6">
       <Card class="p-6 h-32 animate-pulse bg-slate-50 border-0"></Card>
       <Card class="p-6 h-64 animate-pulse bg-slate-50 border-0"></Card>
    </div>

    <div v-else-if="module" class="space-y-6">
      <!-- Module Info Card -->
      <Card class="overflow-hidden border-slate-200">
        <CardContent class="p-0">
          <div class="grid grid-cols-1 md:grid-cols-3 divide-y md:divide-y-0 md:divide-x divide-slate-100">
            <div class="p-6 space-y-1">
              <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">Deskripsi Modul</p>
              <p class="text-slate-600 text-sm leading-relaxed">{{ module.description || 'Tidak ada deskripsi tersedia.' }}</p>
            </div>
            <div class="p-6 space-y-1">
              <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">Metrik & Bobot</p>
              <div class="flex items-center gap-2">
                <span class="text-2xl font-black text-blue-600">{{ module.total_weight }}</span>
                <span class="text-sm text-slate-500 font-medium">Total Point Maksimal</span>
              </div>
              <p class="text-[10px] text-slate-400 mt-1">Total bobot semua pertanyaan tidak boleh melebihi nilai ini.</p>
            </div>
            <div class="p-6 space-y-1">
              <p class="text-xs font-bold text-slate-400 uppercase tracking-wider">Statistik Soal</p>
              <div class="flex items-center gap-4">
                <div class="flex flex-col">
                   <span class="text-2xl font-black text-slate-800">{{ questions.length }}</span>
                   <span class="text-xs text-slate-500">Jumlah Soal</span>
                </div>
                <div class="w-px h-8 bg-slate-100"></div>
                <div class="flex flex-col">
                   <span class="text-2xl font-black" :class="totalWeightUsed > module.total_weight ? 'text-red-500' : 'text-green-600'">{{ totalWeightUsed }}</span>
                   <span class="text-xs text-slate-500">Bobot Terpakai</span>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Questions List Header -->
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-bold text-slate-800">Daftar Pertanyaan</h2>
          <p class="text-sm text-slate-500">Kumpulan soal yang terhubung ke modul ini</p>
        </div>
        <Button @click="router.push({ name: 'questionsCreate', query: { moduleId: module.id } })" class="bg-blue-600 hover:bg-blue-700">
          <PlusIcon class="w-4 h-4 mr-2" /> Tambah Soal ke Modul
        </Button>
      </div>

      <!-- Questions Table -->
      <Card class="border-slate-200 overflow-hidden shadow-sm">
        <Table>
          <TableHeader class="bg-slate-50/50">
            <TableRow>
              <TableHead class="w-[40%] font-bold text-slate-600">Pertanyaan</TableHead>
              <TableHead class="font-bold text-slate-600">Tipe</TableHead>
              <TableHead class="font-bold text-slate-600">Kunci Jawaban</TableHead>
              <TableHead class="font-bold text-slate-600">Bobot</TableHead>
              <TableHead class="font-bold text-slate-600 text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="q in questions" :key="q.id" class="hover:bg-slate-50/50 transition-colors group">
              <TableCell>
                <div class="flex items-start gap-3">
                   <div v-if="q.image_url" class="w-12 h-12 rounded bg-slate-100 shrink-0 border border-slate-200 overflow-hidden">
                      <img :src="q.image_url" class="w-full h-full object-cover" />
                   </div>
                   <p class="text-sm font-medium text-slate-800 line-clamp-2 leading-relaxed">{{ q.content }}</p>
                </div>
              </TableCell>
              <TableCell>
                <Badge variant="outline" class="capitalize font-medium text-slate-500 border-slate-200">
                  {{ q.type.replace('_', ' ') }}
                </Badge>
              </TableCell>
              <TableCell>
                <div v-if="q.type === 'essay'" class="text-xs text-orange-600 font-bold bg-orange-50 px-2 py-1 rounded inline-block">
                  REVIEW MANUAL
                </div>
                <div v-else class="text-xs font-bold text-green-600">
                  <div v-for="opt in q.options?.filter((o: any) => o.is_correct)" :key="opt.id">
                    {{ opt.content }}
                  </div>
                  <span v-if="!q.options?.some((o: any) => o.is_correct)" class="text-slate-300 font-normal italic">Tidak diset</span>
                </div>
              </TableCell>
              <TableCell>
                <div class="flex items-center gap-1.5">
                  <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                  <span class="font-bold text-slate-700 text-sm">{{ q.weight }}</span>
                </div>
              </TableCell>
              <TableCell class="text-right">
                <div class="flex justify-end gap-1">
                  <Button variant="ghost" size="icon" class="h-8 w-8 text-slate-400 hover:text-blue-600 hover:bg-blue-50" @click="router.push(`/questions/edit/${q.id}`)">
                    <EditIcon class="w-4 h-4" />
                  </Button>
                  <Button variant="ghost" size="icon" class="h-8 w-8 text-slate-400 hover:text-red-600 hover:bg-red-50" @click="deleteQuestion(q.id)">
                    <TrashIcon class="w-4 h-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
            <TableRow v-if="questions.length === 0">
              <TableCell colspan="5" class="h-32 text-center text-slate-400 italic">
                Belum ada pertanyaan dalam modul ini.
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Card>
    </div>

    <div v-else class="text-center py-20 bg-slate-50 rounded-2xl border-2 border-dashed border-slate-200">
       <LayersIcon class="w-12 h-12 text-slate-300 mx-auto mb-4" />
       <h3 class="text-lg font-bold text-slate-800">Modul tidak ditemukan</h3>
       <p class="text-slate-500 mb-6">ID modul mungkin salah atau telah dihapus.</p>
       <Button @click="router.push('/modules')" variant="outline">Kembali ke Daftar Modul</Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import client from '@/api/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { 
  ArrowLeftIcon, 
  ChevronRightIcon, 
  PlusIcon, 
  EditIcon, 
  TrashIcon, 
  LayersIcon 
} from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const route = useRoute()
const router = useRouter()
const id = route.params.id as string

const module = ref<any>(null)
const questions = ref<any[]>([])
const loading = ref(true)

const totalWeightUsed = computed(() => {
  return questions.value.reduce((acc, q) => acc + (q.weight || 0), 0)
})

const fetchData = async () => {
  loading.value = true
  try {
    const [resM, resQ] = await Promise.all([
      client.get(`/modules/${id}`),
      client.get(`/modules/${id}/questions`)
    ])
    module.value = resM.data.module
    questions.value = resQ.data || []
  } catch (error) {
    console.error('Failed to fetch module detail:', error)
    toast.error('Gagal memuat detail modul')
  } finally {
    loading.value = false
  }
}

const deleteQuestion = async (qId: string) => {
  if (!confirm('Apakah anda yakin ingin menghapus soal ini?')) return
  try {
    await client.delete(`/questions/${qId}`)
    toast.success('Soal berhasil dihapus')
    questions.value = questions.value.filter(q => q.id !== qId)
  } catch (err) {
    toast.error('Gagal menghapus soal')
  }
}

onMounted(fetchData)
</script>
