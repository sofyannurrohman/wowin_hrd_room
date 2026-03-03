<template>
  <div class="space-y-6 max-w-[1200px] mx-auto pb-10">
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div>
        <h2 class="text-3xl font-extrabold tracking-tight text-slate-800">Participant Management</h2>
        <p class="text-slate-500 mt-1">Manage all candidates, track their progress, and invite new users.</p>
      </div>
      <div class="flex items-center gap-3 w-full md:w-auto">
        <div class="relative w-full md:w-[280px]">
          <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <Input 
            v-model="searchQuery"
            placeholder="Search by name or email..." 
            class="pl-9 h-10 border-slate-200 focus-visible:ring-blue-500 bg-white"
          />
        </div>
        <Button variant="outline" class="h-10 px-4 rounded-lg bg-white border-slate-200 text-slate-600 font-semibold flex items-center gap-2">
          <FilterIcon class="w-4 h-4" /> Filter
        </Button>
        <Button class="h-10 px-4 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-semibold flex items-center gap-2 whitespace-nowrap">
          <UserPlusIcon class="w-4 h-4" /> Invite Participant
        </Button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-[11px] text-slate-500 uppercase tracking-widest">Total Candidates</h3>
            <div class="w-8 h-8 rounded bg-blue-50 text-blue-600 flex items-center justify-center">
              <UsersIcon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.total }}</div>
        </CardContent>
      </Card>

      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-[11px] text-slate-500 uppercase tracking-widest">Active Today</h3>
            <div class="w-8 h-8 rounded bg-green-50 text-green-600 flex items-center justify-center">
              <UserCheckIcon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.activeToday }}</div>
        </CardContent>
      </Card>

      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-[11px] text-slate-500 uppercase tracking-widest">Avg Score</h3>
            <div class="w-8 h-8 rounded bg-purple-50 text-purple-600 flex items-center justify-center">
              <BarChart2Icon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">{{ stats.avgScore }}%</div>
        </CardContent>
      </Card>

      <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden transition-all hover:shadow-md">
        <CardContent class="p-6">
          <div class="flex justify-between items-start mb-4">
            <h3 class="font-medium text-[11px] text-slate-500 uppercase tracking-widest">Pending Invites</h3>
            <div class="w-8 h-8 rounded bg-orange-50 text-orange-600 flex items-center justify-center">
              <MailIcon class="w-4 h-4" />
            </div>
          </div>
          <div class="text-3xl font-black text-slate-800">12</div>
        </CardContent>
      </Card>
    </div>

    <!-- Participants Table -->
    <Card class="border border-slate-100 shadow-sm rounded-2xl bg-white overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow class="border-b border-slate-100 hover:bg-transparent">
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6">Name</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Email</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap">Last Session Attended</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4 whitespace-nowrap">Total Score (Avg)</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-4">Account Status</TableHead>
            <TableHead class="text-[11px] font-bold text-slate-500 uppercase tracking-wider py-4 px-6 text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="p in filteredParticipants" :key="p.id" class="border-b border-slate-100 transition-colors hover:bg-slate-50/50 group">
            <TableCell class="py-4 px-6">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm bg-blue-100 text-blue-700 shrink-0">
                  {{ p.name.substring(0, 2).toUpperCase() }}
                </div>
                <div>
                  <div class="font-bold text-slate-800 text-sm leading-none">{{ p.name }}</div>
                  <div class="text-[11px] text-slate-500 mt-1">ID: CAND-{{ p.id.substring(0,4).toUpperCase() }}</div>
                </div>
              </div>
            </TableCell>
            <TableCell class="py-4 px-4 text-sm text-slate-500">
              {{ p.email }}
            </TableCell>
            <TableCell class="py-4 px-4 text-sm">
              <div class="text-slate-800 font-bold text-[13px] max-w-[180px] truncate" v-if="p.last_session_name">{{ p.last_session_name }}</div>
              <div class="text-slate-800 font-bold text-[13px] text-muted-foreground" v-else>Belum pernah ikut ujian</div>
              <div class="text-slate-500 text-[11px] mt-0.5" v-if="p.last_session_name">{{ formatDateMonthDay(p.last_session_date) }}</div>
            </TableCell>
            <TableCell class="py-4 px-4">
              <div class="flex items-center gap-3">
                <div class="font-bold text-slate-800 w-10">{{ Math.round(p.avg_score) }}%</div>
                <div class="w-20 h-1.5 bg-slate-100 rounded-full overflow-hidden">
                  <div class="h-full bg-green-500 rounded-full" :style="{ width: `${p.avg_score}%` }"></div>
                </div>
              </div>
            </TableCell>
            <TableCell class="py-4 px-4">
              <Badge :class="getStatusBadgeClass(p.status)" class="font-medium text-[11px] px-2.5 py-0.5 rounded-full border-0">
                {{ p.status }}
              </Badge>
            </TableCell>
            <TableCell class="py-4 px-6 text-right">
              <div class="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <!-- Dropdown mock logic or action icons -->
              </div>
            </TableCell>
          </TableRow>
          
          <TableRow v-if="filteredParticipants.length === 0">
             <TableCell colspan="6" class="text-center py-20 text-slate-500">
               <p v-if="loading">Memuat data kandidat...</p>
               <p v-else>Tidak ada kandidat.</p>
             </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      
      <div v-if="filteredParticipants.length > 0" class="p-4 border-t border-slate-100 flex items-center justify-between text-sm text-slate-500 bg-slate-50/30">
        <div>Showing <strong>1</strong> to <strong>5</strong> of <strong>{{ filteredParticipants.length }}</strong> entries</div>
        <div class="flex gap-1 items-center">
           <Button variant="outline" size="sm" class="h-8 border-slate-200 text-slate-600 bg-white" disabled>Previous</Button>
           <Button size="sm" variant="default" class="h-8 w-8 p-0 bg-blue-600 hover:bg-blue-700 text-white font-bold">1</Button>
           <Button size="sm" variant="outline" class="h-8 w-8 p-0 border-slate-200 bg-white text-slate-600">2</Button>
           <Button size="sm" variant="outline" class="h-8 w-8 p-0 border-slate-200 bg-white text-slate-600">3</Button>
           <span class="px-2">...</span>
           <Button variant="outline" size="sm" class="h-8 border-slate-200 text-slate-600 bg-white">Next</Button>
        </div>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import client from '@/api/client'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { SearchIcon, FilterIcon, UserPlusIcon, UsersIcon, UserCheckIcon, BarChart2Icon, MailIcon } from 'lucide-vue-next'

const participants = ref<any[]>([])
const loading = ref(true)
const searchQuery = ref('')
const stats = ref({
  total: 0,
  activeToday: 0,
  avgScore: 0
})

onMounted(async () => {
  try {
    const res = await client.get('/participants')
    participants.value = res.data || []
    
    // Mock calculating dashboard level metrics for Participants page
    if (participants.value.length > 0) {
      stats.value.total = participants.value.length
      stats.value.activeToday = Math.min(participants.value.length, 86) // mock
      
      const totalScore = participants.value.reduce((acc: number, curr: any) => acc + curr.avg_score, 0)
      stats.value.avgScore = Math.round(totalScore / participants.value.length)
    }
  } catch (err) {
    console.error("Failed to fetch participants", err)
  } finally {
    loading.value = false
  }
})

const filteredParticipants = computed(() => {
  return participants.value.filter(p => {
    return p.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
           p.email.toLowerCase().includes(searchQuery.value.toLowerCase())
  })
})

const formatDateMonthDay = (dateString: string) => {
  if (dateString.includes('1970')) return ''
  const d = new Date(dateString)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const getStatusBadgeClass = (status: string) => {
  const statusLower = status.toLowerCase()
  if (statusLower === 'active') return 'bg-green-50 text-green-700 font-bold border border-green-200'
  if (statusLower === 'inactive') return 'bg-slate-50 text-slate-600 font-bold border border-slate-200'
  if (statusLower === 'suspended') return 'bg-red-50 text-red-600 font-bold border border-red-200'
  if (statusLower === 'pending') return 'bg-yellow-50 text-yellow-700 font-bold border border-yellow-200'
  return 'bg-blue-50 text-blue-700 border border-blue-200'
}
</script>
