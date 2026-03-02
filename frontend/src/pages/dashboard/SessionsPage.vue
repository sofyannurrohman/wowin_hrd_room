<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">Manajemen Sesi Ujian</h2>
      <Button @click="router.push('/sessions/create')">
        Sesi Baru
      </Button>
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Nama Sesi</TableHead>
            <TableHead>Waktu Mulai</TableHead>
            <TableHead>Waktu Selesai</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Kapasitas</TableHead>
            <TableHead class="text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="s in sessionStore.sessions" :key="s.id">
            <TableCell class="font-medium">{{ s.name }}</TableCell>
            <TableCell>{{ new Date(s.schedule).toLocaleString() }}</TableCell>
            <TableCell>{{ new Date(new Date(s.schedule).getTime() + s.duration_minutes * 60000).toLocaleString() }}</TableCell>
            <TableCell>
              <Badge :variant="getStatusVariant(s)">{{ getStatus(s) }}</Badge>
            </TableCell>
            <TableCell>{{ s.participant_count || 0 }} / {{ s.max_participants }}</TableCell>
            <TableCell class="text-right flex justify-end gap-2">
              <Button variant="outline" size="sm" @click="router.push(`/sessions/${s.id}`)">Detail</Button>
              <Button variant="default" size="sm" @click="router.push(`/sessions/${s.id}/monitor`)">Monitor</Button>
            </TableCell>
          </TableRow>
          
          <TableRow v-if="sessionStore.sessions.length === 0">
             <TableCell colspan="6" class="text-center h-24">Tidak ada sesi ujian yang ditemukan.</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'

import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const router = useRouter()
const sessionStore = useSessionStore()

onMounted(async () => {
  await sessionStore.fetchSessions()
})

const getStatus = (session: any) => {
  if (session.is_locked) return 'Terkunci'
  const now = new Date()
  const start = new Date(session.schedule)
  const end = new Date(start.getTime() + session.duration_minutes * 60000)
  if (now < start) return 'Akan Datang'
  if (now > end) return 'Berakhir'
  return 'Berjalan'
}

const getStatusVariant = (session: any) => {
  if (session.is_locked) return 'destructive'
  const stat = getStatus(session)
  if (stat === 'Berjalan') return 'default'
  if (stat === 'Akan Datang') return 'secondary'
  return 'outline'
}
</script>
