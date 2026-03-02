<template>
  <div class="flex h-screen bg-muted/40">
    <aside class="w-64 flex flex-col border-r bg-background">
      <div class="h-14 flex items-center px-4 border-b font-semibold text-lg">
        HRD Room
      </div>
      <nav class="flex-1 overflow-y-auto py-4">
        <div class="px-3 space-y-1">
          <RouterLink to="/dashboard" class="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-muted" :class="{ 'bg-muted': $route.path === '/dashboard' }">
            Dashboard
          </RouterLink>
          <RouterLink to="/sessions" class="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-muted" :class="{ 'bg-muted': $route.path.startsWith('/sessions') }">
            Sessions
          </RouterLink>
          <RouterLink to="/questions" class="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-muted" :class="{ 'bg-muted': $route.path.startsWith('/questions') }">
            Questions Bank
          </RouterLink>
          <div v-if="authStore.isAdmin" class="mt-4 pt-4 border-t px-3 pb-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Admin
          </div>
          <RouterLink v-if="authStore.isAdmin" to="/admin/users" class="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-muted" :class="{ 'bg-muted': $route.path.startsWith('/admin/users') }">
            Users
          </RouterLink>
        </div>
      </nav>
      <div class="p-4 border-t">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium truncate">{{ authStore.user?.name }}</span>
          <Button variant="ghost" size="sm" @click="logout">Logout</Button>
        </div>
      </div>
    </aside>
    <main class="flex-1 flex flex-col overflow-hidden">
      <header class="h-14 flex items-center px-6 border-b bg-background justify-between">
        <h1 class="text-lg font-semibold">{{ currentRouteName }}</h1>
      </header>
      <div class="flex-1 overflow-y-auto p-6">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const currentRouteName = computed(() => {
  return route.name?.toString().replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase()) || 'Dashboard'
})

const logout = () => {
  authStore.logout()
}
</script>
