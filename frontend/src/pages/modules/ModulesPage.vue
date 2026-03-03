<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-slate-900">Module Soal</h1>
        <p class="text-sm text-slate-500 mt-1">Manage question banks and modules</p>
      </div>
      <RouterLink
        to="/modules/create"
        class="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium transition-colors bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        <PlusIcon class="w-4 h-4" /> Create Module
      </RouterLink>
    </div>

    <!-- Modules List -->
    <div class="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm">
      <div v-if="loading" class="p-8 text-center text-slate-500">
        Loading modules...
      </div>
      
      <div v-else-if="modules.length === 0" class="p-12 text-center flex flex-col items-center">
        <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mb-4 text-slate-400">
          <LayersIcon class="w-8 h-8" />
        </div>
        <h3 class="text-lg font-medium text-slate-900 mb-1">No Modules Yet</h3>
        <p class="text-slate-500 max-w-sm">Get started by creating your first question module to organize your exam questions.</p>
        <RouterLink to="/modules/create" class="mt-6 text-blue-600 font-medium hover:underline text-sm">
          Create new module
        </RouterLink>
      </div>

      <div v-else class="divide-y divide-slate-100">
        <div
          v-for="mod in modules"
          :key="mod.id"
          class="p-6 transition-colors hover:bg-slate-50 flex items-start gap-4"
        >
          <div class="w-10 h-10 rounded bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
            <LayersIcon class="w-5 h-5" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-1">
              <h3 class="font-semibold text-slate-900 truncate pr-4 text-lg">
                {{ mod.name }}
              </h3>
              <div class="flex gap-2 shrink-0">
                <RouterLink
                  :to="`/modules/edit/${mod.id}`"
                  class="text-slate-400 hover:text-blue-600 transition-colors p-1"
                >
                  <Edit2Icon class="w-4 h-4" />
                </RouterLink>
                <button
                  @click="deleteModule(mod.id)"
                  class="text-slate-400 hover:text-red-600 transition-colors p-1"
                >
                  <Trash2Icon class="w-4 h-4" />
                </button>
              </div>
            </div>
            <p class="text-sm text-slate-500 leading-relaxed mb-3">
              {{ mod.description || 'No description provided.' }}
            </p>
            <div class="flex items-center gap-4 text-xs font-medium text-slate-500">
              <div class="flex items-center gap-1.5 bg-slate-100 px-2.5 py-1 rounded-md">
                <CalendarIcon class="w-3.5 h-3.5 text-slate-400" />
                {{ new Date(mod.created_at).toLocaleDateString() }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { PlusIcon, LayersIcon, Edit2Icon, Trash2Icon, CalendarIcon } from 'lucide-vue-next'
import client from '@/api/client'

interface Module {
  id: string
  name: string
  description: string
  created_at: string
}

const modules = ref<Module[]>([])
const loading = ref(true)

const fetchModules = async () => {
  try {
    const res = await client.get('/modules')
    modules.value = res.data.modules || []
  } catch (error) {
    console.error('Failed to parse modules:', error)
  } finally {
    loading.value = false
  }
}

const deleteModule = async (id: string) => {
  if (!confirm('Are you sure you want to delete this module? This cannot be undone.')) return
  try {
    await client.delete(`/modules/${id}`)
    modules.value = modules.value.filter((m: Module) => m.id !== id)
  } catch (error) {
    alert('Failed to delete module')
  }
}

onMounted(() => {
  fetchModules()
})
</script>
