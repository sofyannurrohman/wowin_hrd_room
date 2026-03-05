<template>
  <div class="flex h-screen bg-slate-50 font-sans">
    <!-- Sidebar -->
    <aside class="w-64 flex flex-col bg-[#253f86] text-white overflow-hidden shadow-xl z-10">
      <!-- Brand / Logo -->
      <div class="h-16 flex items-center px-6 border-b border-white/10 shrink-0">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-full bg-white flex items-center justify-center text-[#253f86] font-bold">
            <img :src="viteLogo" alt="Logo" class="w-8 h-8" />
          </div>
          <span class="font-bold text-lg tracking-wide">HRD Room</span>
        </div>
      </div>

      <!-- Navigation Menu -->
      <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-8 scrollbar-hide">
        
        <!-- Main Menu -->
        <div class="space-y-2">
          <div class="px-2 text-xs font-semibold text-blue-200/70 uppercase tracking-wider mb-3">
            MAIN MENU
          </div>
          
          <RouterLink
            to="/dashboard"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path === '/dashboard' ? 'bg-blue-600/50 text-white' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <LayoutDashboardIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Dashboard
          </RouterLink>

          <RouterLink
            to="/sessions"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/sessions') ? 'bg-blue-600/50 text-white border-l-4 border-blue-400 -ml-1 pl-2' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <CalendarIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Sessions
          </RouterLink>

          <RouterLink
            to="/modules"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/modules') ? 'bg-blue-600/50 text-white border-l-4 border-blue-400 -ml-1 pl-2' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <LayersIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Module Soal
          </RouterLink>

          <RouterLink
            to="/questions"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/questions') ? 'bg-blue-600/50 text-white' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <FileQuestionIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Bank Soal
          </RouterLink>

          <RouterLink
            to="/participants"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/participants') ? 'bg-blue-600/50 text-white border-l-4 border-blue-400 -ml-1 pl-2' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <UsersRoundIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Participants
          </RouterLink>

          <RouterLink
            to="/job-positions"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/job-positions') ? 'bg-blue-600/50 text-white border-l-4 border-blue-400 -ml-1 pl-2' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <BriefcaseIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Posisi Pekerjaan
          </RouterLink>

          <RouterLink
             to="/admin/users"
             v-if="authStore.isAdmin"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/admin') ? 'bg-blue-600/50 text-white' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <ShieldCheckIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Manajemen Users
          </RouterLink>
          
        </div>

        <!-- System -->
        <div class="space-y-2">
          <div class="px-2 text-xs font-semibold text-blue-200/70 uppercase tracking-wider mb-3">
            SYSTEM
          </div>
          
          <RouterLink
            to="/settings"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/settings') ? 'bg-blue-600/50 text-white' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <SettingsIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Settings
          </RouterLink>

          <RouterLink
            to="/support"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors group"
            :class="$route.path.startsWith('/support') ? 'bg-blue-600/50 text-white' : 'text-blue-100 hover:bg-white/10 hover:text-white'"
          >
            <HelpCircleIcon class="w-5 h-5 opacity-80 group-hover:opacity-100" />
            Support
          </RouterLink>
        </div>
      </nav>

      <!-- User Profile Box -->
      <div class="mt-auto px-4 py-4 bg-[#1e3470] border-t border-white/5 cursor-pointer hover:bg-[#253f86] transition-colors relative" @click="showLogout = !showLogout">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center font-bold text-sm shrink-0">
            {{ authStore.user?.name ? authStore.user.name.substring(0, 2).toUpperCase() : 'HR' }}
          </div>
          <div class="overflow-hidden">
            <p class="font-bold text-sm truncate">{{ authStore.user?.name || 'Loading...' }}</p>
            <p class="text-xs text-blue-200 truncate">{{ authStore.isAdmin ? 'Super Admin' : 'Senior HR Specialist' }}</p>
          </div>
        </div>

        <!-- Mini Logout Popup -->
        <div v-if="showLogout" class="absolute bottom-full left-0 w-full mb-2 px-4 z-50">
          <div class="bg-white rounded-lg shadow-lg overflow-hidden border border-slate-200">
             <button @click.stop="logout" class="w-full text-left px-4 py-3 text-sm text-red-600 hover:bg-slate-50 font-medium flex items-center gap-2">
               <LogOutIcon class="w-4 h-4" /> Sign out
             </button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col min-w-0 bg-slate-50 overflow-y-auto">
      <div class="p-6 md:p-8 flex-1">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import {
  LayoutDashboardIcon,
  CalendarIcon,
  LayersIcon,
  FileQuestionIcon,
  UsersRoundIcon,
  BriefcaseIcon,
  SettingsIcon,
  HelpCircleIcon,
  ShieldCheckIcon,
  LogOutIcon
} from 'lucide-vue-next'
import viteLogo from '@/assets/vite.svg'

const authStore = useAuthStore()
const showLogout = ref(false)

const logout = async () => {
  await authStore.logout()
  showLogout.value = false
}
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
