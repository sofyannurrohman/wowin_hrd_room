<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">System Logs (Audit)</h2>
      <Button variant="outline" @click="loadLogs">Refresh</Button>
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[180px]">Waktu</TableHead>
            <TableHead>Level</TableHead>
            <TableHead>Event</TableHead>
            <TableHead>Pengguna ID</TableHead>
            <TableHead>IP Address</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="log in logs" :key="log.id" class="text-sm">
            <TableCell class="text-slate-500 whitespace-nowrap">{{ new Date(log.created_at).toLocaleString() }}</TableCell>
            <TableCell>
              <Badge :variant="log.level === 'error' ? 'destructive' : (log.level === 'warning' ? 'outline' : 'secondary')">
                {{ log.level.toUpperCase() }}
              </Badge>
            </TableCell>
            <TableCell class="font-medium">{{ log.event }}</TableCell>
            <TableCell class="font-mono text-xs">{{ log.user_id || 'System' }}</TableCell>
            <TableCell class="font-mono text-xs">{{ log.ip_address || '-' }}</TableCell>
          </TableRow>
          <TableRow v-if="logs.length === 0">
             <TableCell colspan="5" class="text-center h-24 text-slate-500">Memuat log sistem...</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import client from '@/api/client'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { toast } from 'vue-sonner'

const logs = ref<any[]>([])

const loadLogs = async () => {
  try {
    const res = await client.get('/admin/logs?limit=100')
    logs.value = res.data || []
  } catch(e) {
    toast.error('Gagal memuat log sistem')
  }
}

onMounted(loadLogs)
</script>
