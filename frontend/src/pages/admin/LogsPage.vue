<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold tracking-tight">System Logs (Audit)</h2>
      <Button variant="outline" @click="loadLogs">Refresh</Button>
    </div>

    <Alert v-if="alertMessage" :variant="alertVariant">
      <AlertTitle>{{ alertVariant === 'destructive' ? 'Error' : 'Berhasil' }}</AlertTitle>
      <AlertDescription>{{ alertMessage }}</AlertDescription>
    </Alert>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[180px]">Waktu</TableHead>
            <TableHead>Aksi</TableHead>
            <TableHead>Detail</TableHead>
            <TableHead>Pengguna</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="log in logs" :key="log.id" class="text-sm">
            <TableCell class="text-muted-foreground whitespace-nowrap">{{ new Date(log.created_at).toLocaleString() }}</TableCell>
            <TableCell>
              <Badge variant="outline">{{ log.action.toUpperCase() }}</Badge>
            </TableCell>
            <TableCell class="font-medium">{{ log.detail }}</TableCell>
            <TableCell class="font-mono text-xs">{{ log.user_name || log.user_id || 'System' }}</TableCell>
          </TableRow>
          <TableRow v-if="logs.length === 0">
             <TableCell colspan="5" class="text-center h-24 text-muted-foreground">Memuat log sistem...</TableCell>
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
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

const logs = ref<any[]>([])
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

const loadLogs = async () => {
  try {
    const res = await client.get('/admin/logs?limit=100')
    logs.value = res.data || []
  } catch(e) {
    showError('Gagal memuat log sistem')
  }
}

onMounted(loadLogs)
</script>
