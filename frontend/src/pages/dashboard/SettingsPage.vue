<template>
  <div class="h-full flex flex-col">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Settings</h1>
      <p class="text-gray-500 mt-1">Manage your account settings and preferences.</p>
    </div>

    <!-- Tabs Layout -->
    <div class="flex flex-col md:flex-row gap-8 flex-1 pb-8">
      
      <!-- Sidebar Navigation -->
      <aside class="w-full md:w-64 shrink-0">
        <nav class="flex md:flex-col space-x-2 md:space-x-0 md:space-y-1 overflow-x-auto md:overflow-visible pb-2 md:pb-0 scrollbar-hide">
          <button 
            v-for="tab in tabs" 
            :key="tab.id"
            @click="activeTab = tab.id"
            class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-colors whitespace-nowrap"
            :class="activeTab === tab.id ? 'bg-blue-50 text-blue-700' : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'"
          >
            <component :is="tab.icon" class="w-5 h-5" :class="activeTab === tab.id ? 'text-blue-700' : 'text-gray-400'" />
            {{ tab.name }}
          </button>
        </nav>
      </aside>

      <!-- Main Content Area -->
      <div class="flex-1 max-w-3xl">
        <!-- PROFILE TAB -->
        <div v-if="activeTab === 'profile'" class="space-y-6">
          <div class="bg-white rounded-xl shadow-sm border border-slate-200">
            <div class="p-6 border-b border-slate-100">
              <h2 class="text-lg font-semibold text-gray-900">Personal Information</h2>
              <p class="text-sm text-gray-500 mt-1">Update your personal details here.</p>
            </div>
            
            <div class="p-6 space-y-6">
              <div class="flex items-center gap-6">
                <div class="relative">
                  <div class="w-20 h-20 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold text-2xl">
                    {{ authStore.user?.name ? authStore.user.name.substring(0, 2).toUpperCase() : 'HR' }}
                  </div>
                  <button class="absolute bottom-0 right-0 p-1.5 bg-white border border-slate-200 rounded-full text-gray-500 hover:text-blue-600 hover:border-blue-200 transition-colors shadow-sm">
                    <CameraIcon class="w-4 h-4" />
                  </button>
                </div>
                <div>
                  <h3 class="font-medium text-gray-900">Profile Picture</h3>
                  <p class="text-sm text-gray-500 mt-1">PNG, JPG up to 2MB</p>
                </div>
              </div>

              <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                <div class="space-y-2">
                  <Label for="name">Full Name</Label>
                  <Input id="name" v-model="profileForm.name" />
                </div>
                <div class="space-y-2">
                  <Label for="email">Email Address</Label>
                  <Input id="email" type="email" v-model="profileForm.email" />
                </div>
              </div>

              <div class="pt-4 flex justify-end">
                <Button @click="saveProfile" :disabled="isSavingProfile">
                  <Loader2Icon v-if="isSavingProfile" class="w-4 h-4 mr-2 animate-spin" />
                  Save Changes
                </Button>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-xl shadow-sm border border-slate-200">
            <div class="p-6 border-b border-slate-100">
              <h2 class="text-lg font-semibold text-gray-900">Password</h2>
              <p class="text-sm text-gray-500 mt-1">Update your password associated with your account.</p>
            </div>
            
            <div class="p-6 space-y-4">
              <div class="space-y-2 max-w-md">
                <Label for="currentPassword">Current Password</Label>
                <Input id="currentPassword" type="password" v-model="passwordForm.current" />
              </div>
              <div class="space-y-2 max-w-md">
                <Label for="newPassword">New Password</Label>
                <Input id="newPassword" type="password" v-model="passwordForm.new" />
              </div>
              <div class="space-y-2 max-w-md">
                <Label for="confirmPassword">Confirm New Password</Label>
                <Input id="confirmPassword" type="password" v-model="passwordForm.confirm" />
              </div>

              <div class="pt-4 flex justify-end">
                <Button variant="outline" @click="updatePassword" :disabled="isUpdatingPassword">
                  <Loader2Icon v-if="isUpdatingPassword" class="w-4 h-4 mr-2 animate-spin" />
                  Update Password
                </Button>
              </div>
            </div>
          </div>
        </div>

        <!-- NOTIFICATIONS TAB -->
        <div v-if="activeTab === 'notifications'" class="space-y-6">
          <div class="bg-white rounded-xl shadow-sm border border-slate-200">
            <div class="p-6 border-b border-slate-100">
              <h2 class="text-lg font-semibold text-gray-900">Email Notifications</h2>
              <p class="text-sm text-gray-500 mt-1">Choose what you want to be notified about.</p>
            </div>
            
            <div class="p-6 space-y-6">
              <div class="flex items-start justify-between gap-4">
                <div class="space-y-1 text-sm">
                  <h3 class="font-medium text-gray-900">Session Completed</h3>
                  <p class="text-gray-500">Receive an email when an exam session has successfully ended.</p>
                </div>
                <!-- Replaced Switch with a styled native checkbox -->
                <input type="checkbox" v-model="notificationSettings.sessionCompleted" class="w-5 h-5 text-blue-600 rounded focus:ring-blue-500 mt-1 cursor-pointer" />
              </div>
              
              <div class="w-full h-px bg-slate-100"></div>

              <div class="flex items-start justify-between gap-4">
                <div class="space-y-1 text-sm">
                  <h3 class="font-medium text-gray-900">Manual Grading Required</h3>
                  <p class="text-gray-500">Get notified when a participant submits an essay that needs review.</p>
                </div>
                <input type="checkbox" v-model="notificationSettings.manualGrading" class="w-5 h-5 text-blue-600 rounded focus:ring-blue-500 mt-1 cursor-pointer" />
              </div>

              <div class="w-full h-px bg-slate-100"></div>

              <div class="flex items-start justify-between gap-4">
                <div class="space-y-1 text-sm">
                  <h3 class="font-medium text-gray-900">Candidate Shortlisted</h3>
                  <p class="text-gray-500">Alert me if a candidate scores above the set threshold.</p>
                </div>
                <input type="checkbox" v-model="notificationSettings.candidateShortlisted" class="w-5 h-5 text-blue-600 rounded focus:ring-blue-500 mt-1 cursor-pointer" />
              </div>
            </div>
            <div class="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end">
              <Button @click="saveNotificationSettings" :disabled="isSavingSettings">
                <Loader2Icon v-if="isSavingSettings" class="w-4 h-4 mr-2 animate-spin" />
                Save Preferences
              </Button>
            </div>
          </div>
        </div>

        <!-- SYSTEM TAB -->
        <div v-if="activeTab === 'system'" class="space-y-6">
          <div class="bg-white rounded-xl shadow-sm border border-slate-200">
            <div class="p-6 border-b border-slate-100">
              <h2 class="text-lg font-semibold text-gray-900">System Preferences</h2>
              <p class="text-sm text-gray-500 mt-1">Global settings that apply to all exam sessions.</p>
            </div>
            
            <div class="p-6 space-y-6">
              <div class="flex items-start justify-between gap-4">
                <div class="space-y-1 text-sm">
                  <h3 class="font-medium text-gray-900">Global Strict Mode</h3>
                  <p class="text-gray-500">Enforce strict cheating prevention (disable copy/paste, detect tab switching) for all modules by default.</p>
                </div>
                <input type="checkbox" v-model="systemSettings.strictMode" class="w-5 h-5 text-blue-600 rounded focus:ring-blue-500 mt-1 cursor-pointer" />
              </div>

              <div class="w-full h-px bg-slate-100"></div>

              <div class="space-y-4">
                 <div class="space-y-1 text-sm">
                   <h3 class="font-medium text-gray-900">Data Retention</h3>
                   <p class="text-gray-500">How long should exam results and participant data be kept in the system?</p>
                 </div>
                 <select v-model="systemSettings.dataRetention" class="w-full max-w-xs px-3 py-2 bg-white border border-slate-200 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500">
                   <option value="30 Days">30 Days</option>
                   <option value="90 Days">90 Days</option>
                   <option value="1 Year">1 Year</option>
                   <option value="Keep Forever">Keep Forever</option>
                 </select>
              </div>
            </div>
            <div class="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end">
              <Button @click="saveSystemSettings" :disabled="isSavingSettings">
                <Loader2Icon v-if="isSavingSettings" class="w-4 h-4 mr-2 animate-spin" />
                Apply Settings
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import client from '@/api/client'
import { toast } from 'vue-sonner'
import { 
  UserCircleIcon, 
  BellIcon, 
  SettingsIcon, 
  CameraIcon,
  Loader2Icon
} from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'

const authStore = useAuthStore()

const activeTab = ref('profile')

const tabs = [
  { id: 'profile', name: 'Profile', icon: UserCircleIcon },
  { id: 'notifications', name: 'Notifications', icon: BellIcon },
  { id: 'system', name: 'System', icon: SettingsIcon },
]

// Profile Form State
const profileForm = ref({
  name: '',
  email: ''
})

const passwordForm = ref({
  current: '',
  new: '',
  confirm: ''
})

const isSavingProfile = ref(false)
const isUpdatingPassword = ref(false)

// Notifications State
const notificationSettings = ref({
  sessionCompleted: true,
  manualGrading: true,
  candidateShortlisted: false
})

// System State
const systemSettings = ref({
  strictMode: false,
  dataRetention: 'Keep Forever'
})

const isSavingSettings = ref(false)

const loadSettings = async () => {
  try {
    const resNotif = await client.get('/settings/notification_settings')
    notificationSettings.value = resNotif.data
    
    const resSystem = await client.get('/settings/system_settings')
    systemSettings.value = resSystem.data
  } catch (err) {
    console.error('Failed to load settings', err)
  }
}

onMounted(() => {
  if (authStore.user) {
    profileForm.value.name = authStore.user.name || ''
    profileForm.value.email = authStore.user.email || ''
  }
  loadSettings()
})

const saveProfile = async () => {
  isSavingProfile.value = true
  try {
    await client.put('/auth/profile', profileForm.value)
    if (authStore.user) {
      authStore.user.name = profileForm.value.name
      authStore.user.email = profileForm.value.email
    }
    toast.success('Profil Berhasil Diperbarui')
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal memperbarui profil')
  } finally {
    isSavingProfile.value = false
  }
}

const updatePassword = async () => {
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    toast.error('Kata sandi baru tidak cocok')
    return
  }
  
  isUpdatingPassword.value = true
  try {
    await client.put('/auth/password', {
      current_password: passwordForm.value.current,
      new_password: passwordForm.value.new
    })
    passwordForm.value = { current: '', new: '', confirm: '' }
    toast.success('Kata sandi berhasil diubah')
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal mengubah kata sandi')
  } finally {
    isUpdatingPassword.value = false
  }
}

const saveNotificationSettings = async () => {
  isSavingSettings.value = true
  try {
    await client.put('/settings/notification_settings', notificationSettings.value)
    toast.success('Pengaturan notifikasi disimpan')
  } catch (err) {
    toast.error('Gagal menyimpan pengaturan notifikasi')
  } finally {
    isSavingSettings.value = false
  }
}

const saveSystemSettings = async () => {
  isSavingSettings.value = true
  try {
    await client.put('/settings/system_settings', systemSettings.value)
    toast.success('Pengaturan sistem disimpan')
  } catch (err) {
    toast.error('Gagal menyimpan pengaturan sistem')
  } finally {
    isSavingSettings.value = false
  }
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
