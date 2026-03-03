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
          path: '/participants',
          name: 'participantsList',
          component: () => import('../pages/dashboard/ParticipantsPage.vue')
        },
        {
          path: '/modules',
          name: 'modulesList',
          component: () => import('../pages/modules/ModulesPage.vue')
        },
        {
          path: '/modules/create',
          name: 'moduleCreate',
          component: () => import('../pages/modules/ModuleCreatePage.vue')
        },
        {
          path: '/modules/edit/:id',
          name: 'moduleEdit',
          component: () => import('../pages/modules/ModuleEditPage.vue')
        },
        {
          path: '/questions',
          name: 'questionsList',
          component: () => import('../pages/questions/QuestionsPage.vue')
        },
        {
          path: '/questions/create',
          name: 'questionsCreate',
          component: () => import('../pages/questions/QuestionCreatePage.vue')
        },
        {
          path: '/questions/edit/:id',
          name: 'questionsEdit',
          component: () => import('../pages/questions/QuestionEditPage.vue')
        },
        {
          path: '/job-positions',
          name: 'jobPositionsList',
          component: () => import('../pages/job_positions/JobPositionsPage.vue')
        },
        {
          path: '/job-positions/create',
          name: 'jobPositionsCreate',
          component: () => import('../pages/job_positions/JobPositionCreatePage.vue')
        },
        {
          path: '/job-positions/edit/:id',
          name: 'jobPositionsEdit',
          component: () => import('../pages/job_positions/JobPositionEditPage.vue')
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

router.beforeEach(async (to, _from, next) => {
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
