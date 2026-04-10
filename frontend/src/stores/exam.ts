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
    sessionId: localStorage.getItem('examSessionId') || null,
    sessionName: localStorage.getItem('sessionName') || '',
    modules: [] as any[],
    moduleGroups: [] as ModuleGroup[],
    currentModuleIndex: 0,
    allQuestions: [] as any[],
    questions: [] as any[],
    answers: (() => {
      try {
        const saved = localStorage.getItem('examAnswers')
        return saved ? JSON.parse(saved) : {}
      } catch { return {} }
    })() as Record<string, any>,
    // Restore timeRemaining from the persisted end timestamp so reload doesn't trigger auto-submit
    timeRemaining: (() => {
      const end = Number(localStorage.getItem('examEndTimestamp') || 0)
      if (!end) return 0
      return Math.max(0, Math.floor((end - Date.now()) / 1000))
    })(),
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
      // Persist end timestamp so timer survives page reload
      const endTs = Date.now() + durationSeconds * 1000
      localStorage.setItem('participantId', res.data.participant_id)
      localStorage.setItem('examSessionId', res.data.session_id)
      localStorage.setItem('sessionName', this.sessionName)
      localStorage.setItem('examEndTimestamp', String(endTs))
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
      } catch (e: any) {
        // Reset state but re-throw so callers (ExamPage) can detect if the
        // session was deleted mid-exam and redirect accordingly.
        this.modules = []
        this.questions = []
        this.allQuestions = []
        this.moduleGroups = []
        throw e
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
      // Persist to localStorage for instant restore on page reload
      try {
        localStorage.setItem('examAnswers', JSON.stringify(this.answers))
      } catch { /* ignore quota errors */ }
    },

    /** Fetch already-saved answers from backend and merge into local store */
    async fetchAnswersFromBackend() {
      if (!this.participantId) return
      try {
        const res = await client.get(`/exam/answers/${this.participantId}`)
        const backendAnswers: any[] = res.data || []
        for (const a of backendAnswers) {
          const qId: string = a.question_id
          if (!qId) continue
          // Only overwrite local if backend has a non-null value
          const existing = this.answers[qId]
          const merged: Record<string, any> = { ...(existing || {}) }
          if (a.selected_option_id) merged.selected_option_id = a.selected_option_id
          if (a.text_answer)        merged.text_answer        = a.text_answer
          if (Object.keys(merged).length > 0) {
            this.answers[qId] = merged
          }
        }
        // Re-persist merged result
        localStorage.setItem('examAnswers', JSON.stringify(this.answers))
      } catch (e) {
        console.warn('[exam] fetchAnswersFromBackend failed:', e)
      }
    },

    /** Silently push all current answers to the backend autosave endpoint */
    async autoSaveToBackend() {
      if (!this.participantId || !this.sessionId) return
      const answersArr = Object.entries(this.answers).map(([qId, ans]: [string, any]) => ({
        question_id: qId,
        selected_option_id: ans.selected_option_id || null,
        text_answer: ans.text_answer || null
      }))
      if (answersArr.length === 0) return
      try {
        await client.post(`/exam/${this.sessionId}/answers/autosave`, {
          participant_id: this.participantId,
          answers: answersArr
        })
      } catch (e) {
        console.warn('[exam] autoSaveToBackend failed:', e)
      }
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

    async reportViolation(type: string, proofImage?: string) {
      if (!this.participantId || !this.sessionId) return
      await client.post('/violations', {
        participant_id: this.participantId,
        session_id: this.sessionId,
        violation_type: type,
        proof_image: proofImage
      }).catch(e => console.error('Failed to report violation', e))
    },

    async uploadMonitoringSnapshot(imageData: string) {
      if (!this.participantId || !this.sessionId) return
      await client.post(`/sessions/${this.sessionId}/monitoring`, {
        participant_id: this.participantId,
        image_data: imageData
      }).catch(e => console.error('Failed to upload monitoring photo', e))
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
      this.timeRemaining = 0
      localStorage.removeItem('participantId')
      localStorage.removeItem('examSessionId')
      localStorage.removeItem('sessionName')
      localStorage.removeItem('examEndTimestamp')
      localStorage.removeItem('examAnswers')
    }
  }
})
