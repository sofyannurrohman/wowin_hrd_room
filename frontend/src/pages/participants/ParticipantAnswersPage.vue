<template>
  <div class="space-y-6 max-w-5xl mx-auto pb-10">
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="icon" @click="router.back()" class="rounded-full">
          <ArrowLeftIcon class="w-5 h-5" />
        </Button>
        <div>
          <h2 class="text-2xl font-bold tracking-tight text-slate-800">
            Jawaban Peserta: <span class="text-blue-600">{{ participantName }}</span>
          </h2>
          <p class="text-sm text-slate-500 mt-0.5">Sesi: <span class="font-semibold text-slate-700">{{ sessionName }}</span></p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-20">
      <div class="flex flex-col items-center gap-3 text-slate-400">
        <svg class="w-8 h-8 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
        </svg>
        <span class="text-sm font-medium">Memuat data jawaban...</span>
      </div>
    </div>

    <div v-else class="space-y-6">
      <Card v-for="(ans, index) in combinedAnswers" :key="ans.id || index" class="overflow-hidden border-slate-200">
        <CardHeader class="bg-slate-50 border-b border-slate-100 flex flex-row items-start justify-between py-4">
          <div>
            <div class="flex items-center gap-2 mb-2">
              <Badge variant="outline" class="text-xs bg-white">{{ index + 1 }}</Badge>
              <Badge :class="getTypeColor(ans.question?.type)" class="text-xs">{{ ans.question?.type || 'Unknown' }}</Badge>
            </div>
            <p class="text-sm font-medium text-slate-800">{{ ans.question?.content || 'Pertanyaan tidak ditemukan' }}</p>
            <img
                v-if="ans.question?.image_url"
                :src="ans.question?.image_url"
                class="mt-3 rounded-lg max-h-48 object-contain border border-slate-200"
                alt="Gambar soal"
            />
          </div>
          <div v-if="ans.question?.requires_manual_review" class="flex flex-col items-end gap-2">
            <Badge variant="secondary" class="bg-amber-100 text-amber-800 border-0">Manual Review</Badge>
            <div class="text-xs font-semibold text-slate-500">Skor: <span class="text-slate-900">{{ ans.score || 0 }}</span> / {{ ans.question?.weight || 0 }}</div>
          </div>
          <div v-else class="flex flex-col items-end gap-2">
            <Badge :class="ans.is_correct ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" class="border-0">
              {{ ans.is_correct ? 'Benar' : 'Salah' }}
            </Badge>
            <div class="text-xs font-semibold text-slate-500">Skor: <span class="text-slate-900">{{ ans.score || 0 }}</span> / {{ ans.question?.weight || 0 }}</div>
          </div>
        </CardHeader>
        
        <CardContent class="py-4">
          <!-- Text Answer -->
          <div v-if="['Short Answer', 'Psychological', 'short_answer', 'essay'].includes(ans.question?.type)">
            <div class="mb-4">
              <Label class="text-xs text-slate-500 mb-1 block">Kunci Jawaban / Referensi:</Label>
              <div class="p-3 bg-blue-50/50 rounded-lg border border-blue-100 text-sm text-slate-700 whitespace-pre-wrap">
                <!-- Short Answer uses correct Option as expected text if available -->
                {{ ans.question?.options?.find(o => o.is_correct)?.content || '(Manual review diperlukan - tidak ada kunci)' }}
              </div>
            </div>
            
            <Label class="text-xs text-slate-500 mb-1 block">Jawaban Peserta:</Label>
            <div class="p-3 bg-slate-50 rounded-lg border border-slate-100 text-sm text-slate-700 whitespace-pre-wrap">
              {{ ans.text_answer || '(Tidak dijawab)' }}
            </div>
          </div>

          <!-- Options Answer -->
          <div v-else class="space-y-2 mt-2">
            <div 
              v-for="opt in ans.question?.options" 
              :key="opt.id"
              class="p-3 rounded-lg border text-sm flex items-center justify-between"
              :class="getOptionStyle(opt, ans.selected_option_id)"
            >
              <span>{{ opt.content }}</span>
              <CheckCircleIcon v-if="opt.is_correct" class="w-4 h-4 text-green-600" />
              <XCircleIcon v-if="opt.id === ans.selected_option_id && !opt.is_correct" class="w-4 h-4 text-red-500" />
            </div>
          </div>

          <!-- HR Notes & Manual Scoring -->
          <div v-if="ans.question?.requires_manual_review && reviewData[ans.id]" class="mt-6 pt-4 border-t border-slate-100 space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div class="md:col-span-3 space-y-2">
                <Label class="text-xs text-slate-500">Catatan HR (Opsional)</Label>
                <Textarea v-model="getReview(ans.id).notes" placeholder="Berikan catatan review..." class="min-h-[80px] text-sm" />
              </div>
              <div class="space-y-2">
                <Label class="text-xs text-slate-500">Nilai (Max: {{ ans.question?.weight }})</Label>
                <Input
                  type="number"
                  v-model.number="getReview(ans.id).score"
                  min="0"
                  :max="ans.question?.weight"
                  :class="['text-sm', scoreErrors[ans.id] ? 'border-red-400 focus-visible:ring-red-400' : '']"
                  @input="validateScore(ans.id, ans.question?.weight)"
                />
                <p v-if="scoreErrors[ans.id]" class="text-xs text-red-500 font-medium">
                  {{ scoreErrors[ans.id] }}
                </p>
                <Button @click="submitReview(ans.id)" class="w-full mt-2" size="sm" :disabled="saving === ans.id || !!scoreErrors[ans.id]">
                  {{ saving === ans.id ? 'Loading...' : 'Simpan Nilai' }}
                </Button>
              </div>
            </div>
          </div>

        </CardContent>
      </Card>
      
      <div v-if="combinedAnswers.length === 0" class="text-center py-16 text-slate-400 bg-white rounded-xl border border-slate-100">
        <p>Tidak ada data jawaban untuk peserta ini.</p>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { toast } from 'vue-sonner'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Input } from '@/components/ui/input'
import { ArrowLeftIcon, CheckCircleIcon, XCircleIcon } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const sessionId = route.params.id as string
const participantId = route.params.participantId as string

const loading = ref(true)
const sessionName = ref('Loading...')
const participantName = ref('Loading...')
const answers = ref<any[]>([])
const modules = ref<any[]>([])
const allQuestions = ref<any[]>([])

const reviewData = reactive<Record<string, { score: number, notes: string }>>({})
const scoreErrors = reactive<Record<string, string>>({})
const saving = ref<string | null>(null)

/** Type-safe accessor – always returns a defined review entry (initialised in loadData) */
const getReview = (id: string): { score: number; notes: string } => {
  if (!reviewData[id]) reviewData[id] = { score: 0, notes: '' }
  return reviewData[id] as { score: number; notes: string }
}

const loadData = async () => {
  loading.value = true
  try {
    // Fetch session details
    const resS = await client.get(`/sessions/${sessionId}`)
    sessionName.value = resS.data.name

    // Fetch participant name from results
    const resR = await client.get(`/sessions/${sessionId}/results`)
    const results = resR.data || []
    const participantRes = results.find((r: any) => r.participant_id === participantId)
    if (participantRes) {
      participantName.value = participantRes.user_name || 'Peserta'
    } else {
      participantName.value = 'Peserta'
    }

    // Fetch answers for this participant
    // Since HR only route is GetParticipantAnswers at /results/:participantId/answers
    const resA = await client.get(`/results/${participantId}/answers`)
    answers.value = resA.data || []

    // To display the questions, we need to fetch all questions for this session.
    // Fetch modules first
    const resM = await client.get(`/sessions/${sessionId}/modules`)
    modules.value = resM.data || []

    // Then fetch questions for each module
    let allQs: any[] = []
    for (const m of modules.value) {
      const resQ = await client.get(`/modules/${m.id}/questions`)
      allQs = allQs.concat(resQ.data || [])
    }
    allQuestions.value = allQs

    // Initialize review data mapping for manual review questions
    answers.value.forEach(ans => {
      reviewData[ans.id] = {
        score: ans.score || 0,
        notes: ans.hr_notes || ''
      }
    })

  } catch (e) {
    console.error('Failed to load answers data', e)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

const combinedAnswers = computed(() => {
  // Map answers to their corresponding questions
  return answers.value.map(ans => {
    const question = allQuestions.value.find(q => q.id === ans.question_id)
    return { ...ans, question }
  }).filter(ans => ans.question) // only show answers that have questions
})

const getTypeColor = (type: string) => {
  switch (type) {
    case 'Multiple Choice':
    case 'multiple_choice':
      return 'bg-blue-100 text-blue-700 border-blue-200'
    case 'True/False':
    case 'true_false':
      return 'bg-purple-100 text-purple-700 border-purple-200'
    case 'Short Answer':
    case 'short_answer':
      return 'bg-amber-100 text-amber-700 border-amber-200'
    case 'Psychological':
    case 'essay':
      return 'bg-teal-100 text-teal-700 border-teal-200'
    default: return 'bg-slate-100 text-slate-700 border-slate-200'
  }
}

const getOptionStyle = (opt: any, selectedId: string) => {
  if (opt.is_correct) {
    return 'bg-green-50 border-green-200 text-green-900 border-2'
  }
  if (opt.id === selectedId) {
    return 'bg-red-50 border-red-200 text-red-900 border-2'
  }
  return 'bg-white border-slate-100 text-slate-600'
}

const validateScore = (answerId: string, maxWeight: number | undefined) => {
  if (maxWeight === undefined) return
  const score = reviewData[answerId]?.score ?? 0
  if (score < 0) {
    scoreErrors[answerId] = 'Nilai tidak boleh negatif.'
  } else if (score > maxWeight) {
    scoreErrors[answerId] = `Nilai melebihi batas maksimal (${maxWeight}).`
  } else {
    delete scoreErrors[answerId]
  }
}

const submitReview = async (answerId: string) => {
  if (!reviewData[answerId]) return

  // Find the question to get the max weight
  const ans = combinedAnswers.value.find(a => a.id === answerId)
  const maxWeight = ans?.question?.weight ?? Infinity
  const score = reviewData[answerId].score

  if (score < 0) {
    toast.error('Nilai tidak boleh negatif.')
    return
  }
  if (score > maxWeight) {
    toast.error(`Nilai tidak boleh melebihi batas maksimal (${maxWeight}).`)
    return
  }

  saving.value = answerId
  try {
    await client.put(`/results/${answerId}/review`, {
      answer_id: answerId,
      score,
      notes: reviewData[answerId].notes
    })
    // Successfully saved, update local state
    const ansIndex = answers.value.findIndex(a => a.id === answerId)
    if (ansIndex !== -1) {
      answers.value[ansIndex].score = score
      answers.value[ansIndex].hr_notes = reviewData[answerId].notes
    }
    toast.success('Nilai berhasil disimpan.')
    delete scoreErrors[answerId]
  } catch (e) {
    console.error('Failed to submit review', e)
    toast.error('Gagal menyimpan nilai. Silakan coba lagi.')
  } finally {
    saving.value = null
  }
}
</script>
