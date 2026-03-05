<template>
  <div class="h-full flex flex-col max-w-5xl mx-auto">
    <!-- Header -->
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Support & Help</h1>
        <p class="text-gray-500 mt-1">Get assistance, read documentation, or contact the IT team.</p>
      </div>
      <!-- Search Placeholder -->
      <div class="relative w-64 hidden md:block">
        <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <Input placeholder="Search help articles..." class="pl-9 bg-white" />
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 pb-8">
      
      <!-- Main Content (FAQs) -->
      <div class="lg:col-span-2 space-y-8">
        <section>
          <h2 class="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <BookOpenIcon class="w-5 h-5 text-blue-600" />
            Frequently Asked Questions
          </h2>
          
          <div class="space-y-4">
            <div 
              v-for="(faq, index) in faqs" 
              :key="index"
              class="bg-white border border-slate-200 rounded-xl overflow-hidden transition-all duration-200"
              :class="activeFaq === index ? 'shadow-md border-blue-200' : 'hover:border-slate-300'"
            >
              <button 
                @click="activeFaq = activeFaq === index ? null : index"
                class="w-full text-left px-6 py-4 flex items-center justify-between font-medium text-gray-900 hover:bg-slate-50 transition-colors"
              >
                {{ faq.question }}
                <ChevronDownIcon 
                  class="w-5 h-5 text-gray-400 transition-transform duration-200"
                  :class="activeFaq === index ? 'rotate-180 text-blue-600' : ''"
                />
              </button>
              <div 
                v-show="activeFaq === index"
                class="px-6 pb-4 pt-1 text-sm text-gray-600 leading-relaxed border-t border-slate-100"
              >
                {{ faq.answer }}
              </div>
            </div>
          </div>
        </section>
      </div>

      <!-- Sidebar (Contact & Status) -->
      <div class="space-y-6">
        <!-- Contact Card -->
        <div class="bg-blue-600 rounded-xl p-6 text-white shadow-md relative overflow-hidden group">
          <div class="absolute -right-4 -top-4 w-24 h-24 bg-white/10 rounded-full blur-xl group-hover:bg-white/20 transition-all duration-500"></div>
          
          <HeadphonesIcon class="w-8 h-8 mb-4 text-blue-200" />
          <h3 class="text-lg font-bold mb-2">Need more help?</h3>
          <p class="text-blue-100 text-sm mb-6 leading-relaxed">
            Can't find the answer you're looking for? Contact our IT Support team directly.
          </p>
          
          <a href="mailto:it-support@company.com" class="w-full">
            <Button variant="secondary" class="w-full bg-white text-blue-600 hover:bg-blue-50">
              <MailIcon class="w-4 h-4 mr-2" />
              Contact Support
            </Button>
          </a>
        </div>

        <!-- System Status Card -->
        <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
          <h3 class="font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <ServerIcon class="w-4 h-4 text-gray-500" />
            System Status
          </h3>
          
          <div class="space-y-4">
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">All Systems</span>
              <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                Operational
              </span>
            </div>
            
            <div class="w-full h-px bg-slate-100"></div>
            
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">App Version</span>
              <span class="font-mono text-gray-900 bg-gray-100 px-2 py-0.5 rounded text-xs">v1.2.0</span>
            </div>
            
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Last Updated</span>
              <span class="text-gray-900">Today, 08:00 AM</span>
            </div>
          </div>
        </div>

        <!-- Useful Links Card -->
        <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-sm">
          <h3 class="font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <ExternalLinkIcon class="w-4 h-4 text-gray-500" />
            Useful Resources
          </h3>
          
          <ul class="space-y-3 text-sm">
            <li v-for="(link, i) in links" :key="i">
              <a href="#" class="text-blue-600 hover:underline flex items-center gap-2">
                <FileTextIcon class="w-4 h-4 text-blue-400" />
                {{ link }}
              </a>
            </li>
          </ul>
        </div>
      </div>
      
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { 
  SearchIcon,
  BookOpenIcon,
  ChevronDownIcon,
  HeadphonesIcon,
  MailIcon,
  ServerIcon,
  ExternalLinkIcon,
  FileTextIcon
} from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const activeFaq = ref<number | null>(0)

const faqs = [
  {
    question: "How do I create a new Module Soal?",
    answer: "Go to the 'Module Soal' page and click the 'Create Module' button. You will need to provide a title, description, and the duration for the exam. After creating the module, you can attach questions to it from the Bank Soal."
  },
  {
    question: "A candidate's laptop died during the exam. How do they resume?",
    answer: "Their session is automatically saved. If they log back in using the same exam code and participant token on a different device, they can continue exactly from where they left off without losing their previous answers."
  },
  {
    question: "Why can't a participant join the session yet?",
    answer: "Participants can only join an exam once the scheduled Start Time has been reached. Please check the session details to ensure the start time is correct. If they still cannot join, verify that their participant token matches."
  },
  {
    question: "How does manual grading for essay questions work?",
    answer: "After a session is completed, go to the Session Details -> Monitor -> View Answers for a participant. For essay questions, you will see a text box to input a score and provide feedback. The total score will be recalculated automatically."
  },
  {
    question: "Can I export candidate results to Excel?",
    answer: "Yes. Navigate to the Session Results page. There will be an 'Export to CSV/Excel' button at the top right of the table which allows you to download all participant scores and final statuses."
  }
]

const links = [
  "HR Training Manual PDF",
  "Interview Assessment Rubric",
  "Company Code of Conduct",
  "Privacy Policy"
]
</script>
