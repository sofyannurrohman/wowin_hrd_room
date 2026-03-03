import { defineStore } from 'pinia'
import client from '@/api/client'

export interface ModuleGroup {
  id: string
  name: string
  startIndex: number // inclusive index into allQuestions
  endIndex: number   // inclusive index into allQuestions
}

export const useExamStore = defineStore('exam', {
  state: () => ({
    participantId: localStorage.getItem('participantId') || null,
    sessionId: null as string | null,
    sessionName: localStorage.getItem('sessionName') || '',
    modules: [] as any[],
    moduleGroups: [] as ModuleGroup[],   // one entry per module
    currentModuleIndex: 0,
    allQuestions: [] as any[],           // all questions across all modules, flat
    questions: [] as any[],              // questions for the CURRENT module only
    answers: {} as Record<string, any>,
    timeRemaining: 0,
    isSubmitting: false,
  }),
  actions: {
    async joinExam(payload: { token: string, name: string, email: string, age: number, position: string }) {
      const res = await client.post('/exam/join', payload)
      this.participantId = res.data.participant_id
      this.sessionId = res.data.session_id
      this.sessionName = res.data.session_name || ''
      const durationSeconds = (res.data.time_remaining ?? (res.data.duration * 60)) || 0
      this.timeRemaining = durationSeconds
      localStorage.setItem('participantId', res.data.participant_id)
      localStorage.setItem('sessionName', this.sessionName)
      return res.data
    },

    async fetchModulesAndStart() {
      if (!this.sessionId) return
      try {
        // 1. Fetch the module list
        const modRes = await client.get(`/exam/${this.sessionId}/modules`)
        const modules: any[] = modRes.data || []
        this.modules = modules

        if (modules.length === 0) {
          this.questions = []
          this.allQuestions = []
          this.moduleGroups = []
          return
        }

        // 2. Fetch questions for EVERY module using the module's .id
        const groups: ModuleGroup[] = []
        let flat: any[] = []

        for (let i = 0; i < modules.length; i++) {
          const mod = modules[i]
          const moduleId = mod.id ?? mod.module_id  // prefer .id

          if (!moduleId) {
            console.error('[exam] Module at index', i, 'has no id:', JSON.stringify(mod))
            continue
          }

          const qRes = await client.get(`/exam/${this.sessionId}/modules/${moduleId}/questions`)
          const qs: any[] = qRes.data || []

          groups.push({
            id: moduleId,
            name: mod.name || `Modul ${i + 1}`,
            startIndex: flat.length,
            endIndex: flat.length + qs.length - 1,
          })

          flat = flat.concat(qs)
        }

        this.allQuestions = flat
        this.moduleGroups = groups
        this.currentModuleIndex = 0

        // 3. Load only the first module's questions initially
        this.loadCurrentModuleQuestions()
      } catch (e) {
        console.error('Failed to fetch modules:', e)
        this.modules = []
        this.questions = []
        this.allQuestions = []
        this.moduleGroups = []
      }
    },

    /** Called after currentModuleIndex changes to update visible questions */
    loadCurrentModuleQuestions() {
      const group = this.moduleGroups[this.currentModuleIndex]
      if (!group) { this.questions = []; return }
      this.questions = this.allQuestions.slice(group.startIndex, group.endIndex + 1)
    },

    /** Legacy compat – still called from ExamPage when moving to next module */
    async fetchQuestionsForCurrentModule() {
      this.loadCurrentModuleQuestions()
    },

    setAnswer(questionId: string, answer: any) {
      this.answers[questionId] = answer
    },

    async submitExam() {
      if (this.isSubmitting || !this.participantId) return
      this.isSubmitting = true
      try {
        const payload = {
          participant_id: this.participantId,
          answers: Object.entries(this.answers).map(([qId, ans]: [string, any]) => ({
            question_id: qId,
            selected_option_id: ans.selected_option_id || null,
            text_answer: ans.text_answer || null
          }))
        }
        await client.post(`/exam/${this.sessionId}/answers`, payload)
        this.clearExam()
        return true
      } catch (e) {
        console.error('Failed to submit', e)
        return false
      } finally {
        this.isSubmitting = false
      }
    },

    async reportViolation(type: string) {
      if (!this.participantId || !this.sessionId) return
      await client.post('/violations', {
        participant_id: this.participantId,
        session_id: this.sessionId,
        violation_type: type
      }).catch(e => console.error('Failed to report violation', e))
    },

    clearExam() {
      this.participantId = null
      this.sessionId = null
      this.sessionName = ''
      this.modules = []
      this.moduleGroups = []
      this.currentModuleIndex = 0
      this.allQuestions = []
      this.questions = []
      this.answers = {}
      localStorage.removeItem('participantId')
      localStorage.removeItem('sessionName')
    }
  }
})
