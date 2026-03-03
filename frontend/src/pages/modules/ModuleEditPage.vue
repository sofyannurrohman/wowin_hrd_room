<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div class="flex items-center gap-4">
      <router-link to="/modules" class="p-2 -ml-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100 transition-colors">
        <ArrowLeftIcon class="w-5 h-5" />
      </router-link>
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900">Edit Module</h1>
        <p class="text-sm text-slate-500 mt-1">Update module details</p>
      </div>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow-sm border border-slate-200 p-12 text-center text-slate-500">
      Loading module...
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      <form @submit.prevent="submit" class="p-6 md:p-8 space-y-6">
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1">
              Module Name <span class="text-red-500">*</span>
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
            <label class="block text-sm font-medium text-slate-700 mb-1">Description</label>
            <textarea
              v-model="form.description"
              rows="4"
              class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors outline-none resize-y"
              placeholder="Module description (optional)..."
            ></textarea>
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
            {{ saving ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeftIcon, SaveIcon, Loader2Icon } from 'lucide-vue-next'
import client from '@/api/client'

const router = useRouter()
const route = useRoute()
const saving = ref(false)
const loading = ref(true)

const form = ref({
  name: '',
  description: ''
})

const fetchModule = async () => {
  try {
    const res = await client.get(`/modules/${route.params.id}`)
    form.value.name = res.data.module.name
    form.value.description = res.data.module.description
  } catch (error) {
    alert('Failed to load module')
    router.push('/modules')
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  saving.value = true
  try {
    await client.put(`/modules/${route.params.id}`, form.value)
    router.push('/modules')
  } catch (error) {
    alert('Failed to update module')
    console.error(error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchModule()
})
</script>
