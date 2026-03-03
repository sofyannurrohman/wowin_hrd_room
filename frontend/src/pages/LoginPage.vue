<template>
  <div class="min-h-screen bg-muted/30 flex items-center justify-center p-4">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1 text-center">
        <img :src="viteLogo" alt="Logo" class="mx-auto w-32 h-32" />
        <CardTitle class="text-2xl">Wowin HRD Room</CardTitle>
        <CardDescription>Gunakan Akun HR untuk login ke sistem</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <Label for="email">Email</Label>
          <Input id="email" v-model="email" type="email" placeholder="admin@hrdroom.com" />
        </div>
        <div class="space-y-2">
          <Label for="password">Password</Label>
          <Input id="password" v-model="password" type="password" />
        </div>
        <Alert v-if="errorMsg" variant="destructive">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{{ errorMsg }}</AlertDescription>
        </Alert>
        <Button class="w-full" @click="handleLogin" :disabled="loading">
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </Button>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import client from '@/api/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import viteLogo from '@/assets/vite.svg'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('admin@hrdroom.com')
const password = ref('admin123')
const errorMsg = ref('')
const loading = ref(false)

const handleLogin = async () => {
  errorMsg.value = ''
  loading.value = true
  try {
    const res = await client.post('/auth/login', {
      email: email.value,
      password: password.value
    })
    authStore.setAuth(res.data.user, res.data.token)
    router.push('/dashboard')
  } catch (err: any) {
    errorMsg.value = err.response?.data?.error || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
