import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useExamStore } from '@/stores/exam'

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
      path: '/apply',
      name: 'apply',
      component: () => import('../pages/ApplicationPage.vue')
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
    {
      path: '/privacy-policy',
      name: 'privacyPolicy',
      component: () => import('../pages/PrivacyPolicyPage.vue')
    },
    {
      path: '/terms-of-service',
      name: 'termsOfService',
      component: () => import('../pages/TermsOfServicePage.vue')
    },
    {
      path: '/contact-support',
      name: 'contactSupport',
      component: () => import('../pages/ContactSupportPage.vue')
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
          component: () => import('../pages/sessions/SessionsPage.vue')
        },
        {
          path: '/sessions/create',
          name: 'sessionCreate',
          component: () => import('../pages/sessions/SessionCreatePage.vue')
        },
        {
          path: '/sessions/:id',
          name: 'sessionDetail',
          component: () => import('../pages/sessions/SessionDetailPage.vue')
        },
        {
          path: '/sessions/:id/monitor',
          name: 'sessionMonitor',
          component: () => import('../pages/sessions/SessionMonitorPage.vue')
        },
        {
          path: '/sessions/:id/results',
          name: 'sessionResults',
          component: () => import('../pages/sessions/SessionResultsPage.vue')
        },
        {
          path: '/sessions/:id/violations',
          name: 'sessionViolations',
          component: () => import('../pages/sessions/SessionViolationsPage.vue')
        },
        {
          path: '/sessions/:id/analytics',
          name: 'sessionAnalytics',
          component: () => import('../pages/sessions/SessionAnalyticsPage.vue')
        },
        {
          path: '/sessions/:id/results/:participantId/answers',
          name: 'participantAnswers',
          component: () => import('../pages/participants/ParticipantAnswersPage.vue')
        },
        {
          path: '/participants',
          name: 'participantsList',
          component: () => import('../pages/participants/ParticipantsPage.vue')
        },
        {
          path: '/participants/:id',
          name: 'participantDetail',
          component: () => import('../pages/participants/ParticipantDetailPage.vue')
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
          path: '/settings',
          name: 'settings',
          component: () => import('../pages/dashboard/SettingsPage.vue')
        },
        {
          path: '/support',
          name: 'support',
          component: () => import('../pages/dashboard/SupportPage.vue')
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
  const examStore = useExamStore()
  
  // On page reload: token is in localStorage (isAuthenticated = true) but user profile
  // is null (Pinia resets in-memory). Fetch the profile before checking role-based guards.
  if (authStore.isAuthenticated && !authStore.user) {
    await authStore.fetchUser()
  }

  // If route requires auth and user is still not authenticated (bad/expired token), redirect to login
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next({ name: 'login' })
  }

  if (to.meta.requiresHR && !authStore.isHR) {
    return next({ name: 'join' }) // fallback to public if not HR
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return next({ name: 'dashboardOverview' })
  }

  // ─── Exam session guard ──────────────────────────────────────────────────
  // If user has an active exam session, redirect /join back to the exam
  if (to.name === 'join') {
    const participantId = examStore.participantId || localStorage.getItem('participantId')
    const sessionId = examStore.sessionId || localStorage.getItem('examSessionId')
    if (participantId && sessionId) {
      return next({ name: 'exam', params: { sessionId } })
    }
  }

  // If exam route requires an active session, redirect to /join if missing
  if (to.meta.requiresExamSession) {
    const participantId = examStore.participantId || localStorage.getItem('participantId')
    const sessionId = examStore.sessionId || localStorage.getItem('examSessionId')
    if (!participantId || !sessionId) {
      return next({ name: 'join' })
    }
  }

  next()
})

export default router
