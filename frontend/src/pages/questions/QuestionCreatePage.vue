<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <div class="flex items-center gap-4">
      <Button variant="ghost" size="icon" @click="router.back()">
        <ArrowLeftIcon class="w-5 h-5" />
      </Button>
      <h2 class="text-3xl font-bold tracking-tight">Buat Soal Baru</h2>
    </div>

    <Card>
      <CardContent class="pt-6 space-y-6">
        <form @submit.prevent="submit" class="space-y-6">
          <div class="space-y-2">
            <Label for="type">Tipe Soal</Label>
            <select v-model="form.type" id="type" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
              <option value="multiple_choice">Pilihan Ganda</option>
              <option value="true_false">Benar / Salah</option>
              <option value="short_answer">Isian Singkat</option>
              <option value="essay">Esai / Psikologi (Review Manual)</option>
            </select>
          </div>

          <div class="space-y-2">
            <Label for="content">Pertanyaan</Label>
            <Textarea id="content" v-model="form.content" class="min-h-[100px]" required />
          </div>

          <div class="space-y-2">
            <Label for="image">Gambar Lampiran (Opsional)</Label>
            <Input id="image" type="file" @change="onFileChange" accept="image/*" class="cursor-pointer"/>
          </div>

          <div class="space-y-2">
            <Label>Lampirkan ke Modul Soal (Opsional)</Label>
            <select v-model="form.module_id" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2">
              <option :value="''">-- Jadikan Bank Soal Global --</option>
              <option v-for="m in modules" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
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

          <div class="flex justify-end pt-6 border-t">
            <Button type="submit" :disabled="loading">
              {{ loading ? 'Menyimpan...' : 'Simpan Soal' }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import client from '@/api/client'

import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { ArrowLeftIcon, TrashIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const router = useRouter()
const modules = ref<any[]>([])

const loading = ref(false)

const form = ref({
  module_id: '',
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

onMounted(async () => {
    try {
      const res = await client.get('/modules')
      modules.value = res.data.modules || []
    } catch(e) {
      toast.error('Gagal memuat daftar modul')
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
  try {
    const formData = new FormData()
    formData.append('content', form.value.content)
    formData.append('type', form.value.type)
    if (form.value.module_id) formData.append('module_id', form.value.module_id)
    if (fileObj.value) formData.append('image', fileObj.value)

    // Append options as JSON string based on type
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
    // essay has no fixed options

    formData.append('options', JSON.stringify(finalOptions))

    await client.post('/questions', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })

    toast.success('Soal berhasil ditambahkan')
    router.push('/questions')
  } catch(e: any) {
    toast.error(e.response?.data?.error || 'Gagal menambahkan soal')
  } finally {
    loading.value = false
  }
}
</script>
