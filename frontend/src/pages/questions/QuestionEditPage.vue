<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <div class="flex items-center gap-4">
      <Button variant="ghost" size="icon" @click="router.back()">
        <ArrowLeftIcon class="w-5 h-5" />
      </Button>
      <h2 class="text-3xl font-bold tracking-tight">Edit Soal</h2>
    </div>

    <Card>
      <CardContent class="pt-6 space-y-6">
        <Alert v-if="alertMessage" :variant="alertVariant">
          <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
          <AlertDescription>{{ alertMessage }}</AlertDescription>
        </Alert>

        <form @submit.prevent="submit" class="space-y-6">
          <div class="space-y-2">
            <Label for="type">Tipe Soal</Label>
            <select v-model="form.type" id="type" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
              <option value="multiple_choice">Pilihan Ganda</option>
              <option value="true_false">Benar / Salah</option>
              <option value="short_answer">Isian Singkat</option>
              <option value="essay">Esai / Psikologi (Review Manual)</option>
            </select>
          </div>

          <div class="space-y-2">
            <Label for="content">Pertanyaan</Label>
            <textarea id="content" v-model="form.content" class="w-full p-3 border rounded-md min-h-[100px] outline-none focus:ring-2 focus:ring-primary" required></textarea>
          </div>

          <div class="space-y-2">
            <Label for="image">Gambar Lampiran Baru (Opsional)</Label>
            <div v-if="existingImage" class="mb-2">
              <span class="text-sm text-gray-500 block mb-1">Gambar saat ini:</span>
              <img :src="getImageUrl(existingImage)" class="max-h-32 rounded-md object-contain border" />
            </div>
            <Input id="image" type="file" @change="onFileChange" accept="image/*" class="cursor-pointer"/>
          </div>

          <div class="space-y-2">
            <Label>Sesi Ujian</Label>
            <Input type="text" :value="form.session_id ? 'Terkait Sesi Tertentu' : 'Global Bank Soal'" disabled class="bg-gray-100 dark:bg-gray-800" />
            <p class="text-xs text-gray-500">Pemindahan sesi tidak diizinkan saat edit.</p>
          </div>

          <!-- Options Builder for MCQ -->
          <div v-if="form.type === 'multiple_choice'" class="space-y-4 pt-4 border-t">
            <div class="flex justify-between items-center">
              <Label>Pilihan Jawaban (Tandai yang Benar)</Label>
              <Button type="button" variant="outline" size="sm" @click="addOption">Tambah Pilihan</Button>
            </div>
            
            <div v-for="(opt, index) in options" :key="index" class="flex items-center gap-3">
              <input type="radio" :name="'correctOption'" :value="index" v-model="correctOptionIndex" class="w-5 h-5 accent-primary" required />
              <Input v-model="opt.text" :placeholder="`Pilihan ${index + 1}`" required />
              <Button type="button" variant="ghost" size="icon" @click="removeOption(index)" v-if="options.length > 2">
                <TrashIcon class="w-4 h-4 text-destructive" />
              </Button>
            </div>
          </div>

           <!-- Options Builder for True False -->
           <div v-if="form.type === 'true_false'" class="space-y-4 pt-4 border-t">
            <Label>Pilihan Jawaban Benar / Salah</Label>
            <div class="flex items-center gap-6 mt-2">
              <label class="flex items-center gap-2">
                <input type="radio" value="true" v-model="correctBoolAnswer" class="w-5 h-5 accent-primary" name="tf" required />
                <span>Pernyataan Benar</span>
              </label>
              <label class="flex items-center gap-2">
                <input type="radio" value="false" v-model="correctBoolAnswer" class="w-5 h-5 accent-primary" name="tf" required />
                <span>Pernyataan Salah</span>
              </label>
            </div>
          </div>

          <!-- Correct Keyword for Short Answer -->
          <div v-if="form.type === 'short_answer'" class="space-y-2 pt-4 border-t">
            <Label>Kunci Jawaban Singkat (Exact Match)</Label>
            <Input v-model="correctTextAnswer" placeholder="Budi Oetomo" required/>
          </div>

          <div class="flex justify-end pt-6 border-t gap-2">
             <Button type="button" variant="outline" @click="router.back()">Batal</Button>
            <Button type="submit" :disabled="loading">
              {{ loading ? 'Menyimpan...' : 'Simpan Perubahan' }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import client from '@/api/client'

import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { ArrowLeftIcon, TrashIcon } from 'lucide-vue-next'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const existingImage = ref<string | null>(null)
const alertMessage = ref('')
const alertVariant = ref<'default' | 'destructive'>('default')

const showSuccess = (message: string) => {
  alertVariant.value = 'default'
  alertMessage.value = message
}

const showError = (message: string) => {
  alertVariant.value = 'destructive'
  alertMessage.value = message
}

const form = ref({
  session_id: '',
  content: '',
  type: 'multiple_choice',
})

const fileObj = ref<File | null>(null)

// For multiple choice
const options = ref([{ text: '' }, { text: '' }, { text: '' }, { text: '' }])
const correctOptionIndex = ref<number>(0)

// For true false
const correctBoolAnswer = ref<string>('true')

// For short answer
const correctTextAnswer = ref('')

let questionId = ''

const getImageUrl = (path: string) => {
    return import.meta.env.VITE_API_BASE_URL.replace('/api', '') + path
}

onMounted(async () => {
    questionId = route.params.id as string
    try {
        const res = await client.get(`/questions/${questionId}`)
        const data = res.data
        if (!data) return
        
        form.value.session_id = data.session_id === '00000000-0000-0000-0000-000000000000' ? '' : data.session_id
        form.value.content = data.content
        form.value.type = data.type
        existingImage.value = data.image_url

        if (data.options && data.options.length > 0) {
            if (data.type === 'multiple_choice') {
                options.value = data.options.map((o: any) => ({ text: o.content }))
                correctOptionIndex.value = data.options.findIndex((o: any) => o.is_correct)
            } else if (data.type === 'true_false') {
                const correctOpt = data.options.find((o: any) => o.is_correct)
                if (correctOpt) {
                    correctBoolAnswer.value = correctOpt.content === 'Benar' ? 'true' : 'false'
                }
            } else if (data.type === 'short_answer') {
                const correctOpt = data.options.find((o: any) => o.is_correct)
                if (correctOpt) {
                     correctTextAnswer.value = correctOpt.content
                }
            }
        }
    } catch(e) {
        showError('Gagal memuat detail soal')
    }
})

const onFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    fileObj.value = target.files[0]
  }
}

const addOption = () => options.value.push({ text: '' })
const removeOption = (idx: number) => {
  if (options.value.length > 2) options.value.splice(idx, 1)
  if (correctOptionIndex.value === idx) correctOptionIndex.value = 0
}

const submit = async () => {
  loading.value = true
  alertMessage.value = ''
  try {
    const formData = new FormData()
    formData.append('content', form.value.content)
    formData.append('type', form.value.type)
    if (fileObj.value) formData.append('image', fileObj.value)

    let finalOptions: any[] = []

    if (form.value.type === 'multiple_choice') {
      finalOptions = options.value.map((o, idx) => ({
        content: o.text,
        is_correct: correctOptionIndex.value === idx
      }))
    } else if (form.value.type === 'true_false') {
      finalOptions = [
        { content: 'Benar', is_correct: correctBoolAnswer.value === 'true' },
        { content: 'Salah', is_correct: correctBoolAnswer.value === 'false' }
      ]
    } else if (form.value.type === 'short_answer') {
      finalOptions = [{ content: correctTextAnswer.value, is_correct: true }] // Convention for short answer
    }

    formData.append('options', JSON.stringify(finalOptions))

    await client.put(`/questions/${questionId}`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })

    showSuccess('Soal berhasil diperbarui')
    router.push('/questions')
  } catch(e: any) {
    showError(e.response?.data?.error || 'Gagal memperbarui soal')
  } finally {
    loading.value = false
  }
}
</script>
