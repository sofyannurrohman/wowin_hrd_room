<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <h2 class="text-3xl font-bold tracking-tight">Hasil Ujian & Evaluasi</h2>
      </div>
      <div class="flex gap-2">
        <Button variant="secondary" @click="autoGrade" :disabled="grading" class="border">
          <Wand2Icon class="w-4 h-4 mr-2" /> Auto-Grade Pilihan Ganda
        </Button>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Daftar Hasil Peserta</CardTitle>
        <CardDescription>Review hasil grading otomatis dan berikan nilai manual untuk esai.</CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Peserta</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Status Ujian</TableHead>
              <TableHead>Skor</TableHead>
              <TableHead>Status Review</TableHead>
              <TableHead class="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="res in results" :key="res.id">
              <TableCell class="font-medium">{{ res.participant?.user?.name || res.participant_id }}</TableCell>
              <TableCell>{{ res.participant?.user?.email || '-' }}</TableCell>
              <TableCell>Disubmit</TableCell>
              <TableCell class="font-bold text-primary">{{ res.total_score }}</TableCell>
              <TableCell>
                <Badge :variant="res.status === 'Finalized' ? 'default' : 'secondary'">{{ res.status }}</Badge>
              </TableCell>
              <TableCell class="text-right">
                <Button size="sm" @click="openReview(res)">Review Esai</Button>
              </TableCell>
            </TableRow>
            <TableRow v-if="results.length === 0">
              <TableCell colspan="6" class="text-center h-24">Belum ada hasil yang masuk. Pastikan peserta telah Submit ujian.</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Dialog for Manual Review (Simplified) -->
    <Dialog v-model:open="showReviewDialog" class="max-w-4xl">
      <DialogContent class="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Review Manual: {{ currentReview?.participant?.user?.name }}</DialogTitle>
          <DialogDescription>
            Berikan nilai dan catatan HR untuk jawaban subyektif (Esai/Psikologi).
          </DialogDescription>
        </DialogHeader>
        
        <div class="space-y-6 my-4">
          <!-- In a real app, you would fetch the full answers for this participant. 
               This is a simplified assumption where answers are part of the result payload or fetched separately. -->
          <Card class="border-orange-200 bg-orange-50/20">
            <CardHeader>
              <CardTitle class="text-base text-orange-800">HR Notes & Penyesuaian Skor</CardTitle>
            </CardHeader>
            <CardContent class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label>Skor Total Final (0-100)</Label>
                <Input type="number" v-model.number="reviewForm.total_score" min="0" max="100"/>
              </div>
              <div class="space-y-2">
                <Label>Catatan HR</Label>
                <textarea v-model="reviewForm.hr_notes" class="w-full p-2 border rounded resize-y focus:ring-1 focus:ring-primary outline-none" rows="3" placeholder="Contoh: Jawaban esai sangat analitis..."></textarea>
              </div>
            </CardContent>
          </Card>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="showReviewDialog = false">Batal</Button>
          <Button variant="default" @click="submitFinalize">Finalisasi Nilai</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { ArrowLeftIcon, Wand2Icon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string

const results = ref<any[]>([])
const grading = ref(false)

const showReviewDialog = ref(false)
const currentReview = ref<any>(null)
const reviewForm = ref({ total_score: 0, hr_notes: '' })

const loadResults = async () => {
  try {
    const res = await client.get(`/sessions/${sessionId}/results`)
    results.value = res.data || []
  } catch (e) {
    toast.error('Gagal memuat hasil ujian')
  }
}

onMounted(loadResults)

const autoGrade = async () => {
  grading.value = true
  try {
    const res = await client.post(`/sessions/${sessionId}/grade`)
    toast.success(res.data.message || 'Auto grading selesai')
    await loadResults()
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Gagal melakukan grading otomatis')
  } finally {
    grading.value = false
  }
}

const openReview = (res: any) => {
  currentReview.value = res
  reviewForm.value.total_score = res.total_score
  reviewForm.value.hr_notes = res.hr_notes || ''
  showReviewDialog.value = true
}

const submitFinalize = async () => {
  try {
    await client.post(`/results/${currentReview.value.id}/finalize`, {
      total_score: reviewForm.value.total_score,
      hr_notes: reviewForm.value.hr_notes
    })
    toast.success('Nilai berhasil di-finalisasi')
    showReviewDialog.value = false
    await loadResults()
  } catch(e) {
    toast.error('Gagal finalisasi nilai')
  }
}
</script>
