import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      redirect: '/join'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../pages/LoginPage.vue')
    },
    {
      path: '/join',
      name: 'join',
      component: () => import('../pages/JoinPage.vue')
    },
    {
      path: '/camera-check',
      name: 'cameraCheck',
      component: () => import('../pages/CameraCheckPage.vue')
    },
    {
      path: '/exam/:sessionId',
      name: 'exam',
      component: () => import('../pages/ExamPage.vue'),
      meta: { requiresExamSession: true }
    },
    {
      path: '/exam-finished',
      name: 'examFinished',
      component: () => import('../pages/ExamFinishedPage.vue')
    },
    // HR / Admin Dashboard
    {
      path: '/dashboard',
      component: () => import('../layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true, requiresHR: true },
      children: [
        {
          path: '',
          name: 'dashboardOverview',
          component: () => import('../pages/dashboard/DashboardPage.vue')
        },
        {
          path: '/sessions',
          name: 'sessionsList',
          component: () => import('../pages/dashboard/SessionsPage.vue')
        },
        {
          path: '/sessions/create',
          name: 'sessionCreate',
          component: () => import('../pages/dashboard/SessionCreatePage.vue')
        },
        {
          path: '/sessions/:id',
          name: 'sessionDetail',
          component: () => import('../pages/dashboard/SessionDetailPage.vue')
        },
        {
          path: '/sessions/:id/monitor',
          name: 'sessionMonitor',
          component: () => import('../pages/dashboard/SessionMonitorPage.vue')
        },
        {
          path: '/sessions/:id/results',
          name: 'sessionResults',
          component: () => import('../pages/dashboard/SessionResultsPage.vue')
        },
        {
          path: '/questions',
          name: 'questionsList',
          component: () => import('../pages/questions/QuestionsPage.vue')
        },
        {
          path: '/admin/users',
          name: 'adminUsers',
          component: () => import('../pages/admin/UsersPage.vue'),
          meta: { requiresAdmin: true }
        }
      ]
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // if not authenticated and it requires auth, fetch profile or redirect
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    await authStore.fetchUser()
    if (!authStore.isAuthenticated) {
      return next({ name: 'login' })
    }
  }

  if (to.meta.requiresHR && !authStore.isHR) {
    return next({ name: 'join' }) // fallback to public if not HR
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return next({ name: 'dashboardOverview' })
  }

  next()
})

export default router
