<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Dashboard Overview</h2>
      <Button @click="router.push('/sessions/create')">
        <PlusCircleIcon class="mr-2 h-4 w-4" /> Buka Sesi Baru
      </Button>
    </div>

    <!-- Stats -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0.5 pb-2">
          <CardTitle class="text-sm font-medium">Total Sesi Ujian</CardTitle>
          <CalendarIcon class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ sessionStore.sessions.length }}</div>
        </CardContent>
      </Card>
      
      <!-- Placeholder for more stats -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0.5 pb-2">
          <CardTitle class="text-sm font-medium">Peserta Aktif</CardTitle>
          <UsersIcon class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-green-600">Terpantau Realtime</div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0.5 pb-2">
          <CardTitle class="text-sm font-medium">Keamanan</CardTitle>
          <ShieldAlertIcon class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">Anti-Cheat Aktif</div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0.5 pb-2">
          <CardTitle class="text-sm font-medium">Bank Soal</CardTitle>
          <DatabaseIcon class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">Terintegrasi</div>
        </CardContent>
      </Card>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <Card class="col-span-4">
        <CardHeader>
          <CardTitle>Sesi Terbaru</CardTitle>
          <CardDescription>Manajemen akses cepat sesi ujian Anda</CardDescription>
        </CardHeader>
        <CardContent>
           <div class="space-y-4">
            <div v-for="session in recentSessions" :key="session.id" class="flex items-center justify-between p-4 border rounded-md">
              <div class="space-y-1">
                <p class="font-medium lead text-base">{{ session.title }}</p>
                <p class="text-sm text-muted-foreground">{{ new Date(session.start_time).toLocaleString() }} - {{ session.participant_limit }} Peserta</p>
              </div>
              <div class="flex gap-2">
                <Button variant="outline" size="sm" @click="router.push(`/sessions/${session.id}`)">Detail</Button>
                <Button variant="default" size="sm" @click="router.push(`/sessions/${session.id}/monitor`)">Live Monitor</Button>
              </div>
            </div>
            <div v-if="recentSessions.length === 0" class="text-center text-muted-foreground py-8">
              Belum ada sesi ujian dibuat.
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { CalendarIcon, UsersIcon, ShieldAlertIcon, DatabaseIcon, PlusCircleIcon } from 'lucide-vue-next'

const router = useRouter()
const sessionStore = useSessionStore()

const recentSessions = computed(() => {
  return [...sessionStore.sessions].slice(0, 5) // Show top 5
})

onMounted(async () => {
  await sessionStore.fetchSessions()
})
</script>
