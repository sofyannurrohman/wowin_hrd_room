<template>
  <div class="flex items-center justify-between px-2 py-4">
    <div class="flex-1 text-sm text-slate-500">
      Showing {{ from }} to {{ to }} of {{ total }} entries
    </div>
    <div class="flex items-center space-x-6 lg:space-x-8">
      <div class="flex items-center space-x-2">
        <p class="text-sm font-medium">Rows per page</p>
        <select
          :value="pageSize"
          @change="$emit('update:pageSize', Number(($event.target as HTMLSelectElement).value))"
          class="h-8 w-[70px] rounded-md border border-slate-200 bg-white px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option v-for="size in [5, 10, 20, 30, 40, 50]" :key="size" :value="size">
            {{ size }}
          </option>
        </select>
      </div>
      <div class="flex w-[100px] items-center justify-center text-sm font-medium">
        Page {{ currentPage }} of {{ totalPages }}
      </div>
      <div class="flex items-center space-x-2">
        <Button
          variant="outline"
          class="hidden h-8 w-8 p-0 lg:flex"
          :disabled="currentPage === 1"
          @click="$emit('update:currentPage', 1)"
        >
          <span class="sr-only">Go to first page</span>
          <ChevronsLeftIcon class="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          class="h-8 w-8 p-0"
          :disabled="currentPage === 1"
          @click="$emit('update:currentPage', currentPage - 1)"
        >
          <span class="sr-only">Go to previous page</span>
          <ChevronLeftIcon class="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          class="h-8 w-8 p-0"
          :disabled="currentPage === totalPages"
          @click="$emit('update:currentPage', currentPage + 1)"
        >
          <span class="sr-only">Go to next page</span>
          <ChevronRightIcon class="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          class="hidden h-8 w-8 p-0 lg:flex"
          :disabled="currentPage === totalPages"
          @click="$emit('update:currentPage', totalPages)"
        >
          <span class="sr-only">Go to last page</span>
          <ChevronsRightIcon class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { 
  ChevronLeftIcon, 
  ChevronRightIcon, 
  ChevronsLeftIcon, 
  ChevronsRightIcon 
} from 'lucide-vue-next'

const props = defineProps<{
  total: number
  pageSize: number
  currentPage: number
}>()

defineEmits<{
  (e: 'update:pageSize', value: number): void
  (e: 'update:currentPage', value: number): void
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize) || 1)
const from = computed(() => props.total === 0 ? 0 : (props.currentPage - 1) * props.pageSize + 1)
const to = computed(() => Math.min(props.currentPage * props.pageSize, props.total))
</script>
