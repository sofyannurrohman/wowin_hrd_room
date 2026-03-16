import { ref, computed, watch, type Ref } from 'vue'

export interface SortConfig {
  key: string
  direction: 'asc' | 'desc'
}

export function useDataTable<T>(data: Ref<T[]>, initialPageSize = 10) {
  const pageSize = ref(initialPageSize)
  const currentPage = ref(1)
  const sortConfig = ref<SortConfig>({
    key: '',
    direction: 'asc'
  })

  const toggleSort = (key: string) => {
    if (sortConfig.value.key === key) {
      sortConfig.value.direction = sortConfig.value.direction === 'asc' ? 'desc' : 'asc'
    } else {
      sortConfig.value.key = key
      sortConfig.value.direction = 'asc'
    }
    // Reset to first page when sorting changes
    currentPage.value = 1
  }

  const sortedData = computed(() => {
    let list = [...data.value]
    if (sortConfig.value.key) {
      list.sort((a: any, b: any) => {
        const key = sortConfig.value.key
        const aVal = a[key]
        const bVal = b[key]

        const aStr = String(aVal || '').toLowerCase()
        const bStr = String(bVal || '').toLowerCase()

        if (aStr < bStr) return sortConfig.value.direction === 'asc' ? -1 : 1
        if (aStr > bStr) return sortConfig.value.direction === 'asc' ? 1 : -1
        return 0
      })
    }
    return list
  })

  const paginatedData = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return sortedData.value.slice(start, end)
  })

  // Reset page when data changes (e.g. search)
  watch(data, () => {
    currentPage.value = 1
  })

  return {
    pageSize,
    currentPage,
    sortConfig,
    toggleSort,
    paginatedData,
    totalItems: computed(() => data.value.length)
  }
}
