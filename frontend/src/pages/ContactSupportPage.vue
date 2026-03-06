<template>
  <div class="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8 flex flex-col items-center justify-center">
    <div class="max-w-xl w-full bg-white rounded-2xl shadow-sm p-8 md:p-12">
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4 text-blue-600">
          <HeadphonesIcon class="w-8 h-8" />
        </div>
        <h1 class="text-3xl font-extrabold text-slate-900 tracking-tight mb-2">Contact Support</h1>
        <p class="text-slate-500">Need help with your exam? We're here to assist you.</p>
      </div>

      <div class="space-y-6">
        <div class="bg-slate-50 rounded-xl p-6 border border-slate-100 flex items-start gap-4">
          <MailIcon class="w-6 h-6 text-slate-400 shrink-0 mt-0.5" />
          <div>
            <h3 class="font-bold text-slate-800">Email Support</h3>
            <p class="text-slate-600 text-sm mt-1 mb-2">Send us an email and we'll get back to you within 24 hours.</p>
            <a href="mailto:support@hrdroom.com" class="text-blue-600 font-medium hover:underline">support@hrdroom.com</a>
          </div>
        </div>

        <div class="bg-slate-50 rounded-xl p-6 border border-slate-100 flex items-start gap-4">
          <MessageCircleIcon class="w-6 h-6 text-slate-400 shrink-0 mt-0.5" />
          <div>
            <h3 class="font-bold text-slate-800">WhatsApp Support</h3>
            <p class="text-slate-600 text-sm mt-1 mb-2">Available during business hours (09:00 - 17:00 WIB).</p>
            <a href="https://wa.me/6281234567890" target="_blank" class="text-blue-600 font-medium hover:underline">+62 812 3456 7890</a>
          </div>
        </div>
        
        <div class="bg-amber-50 rounded-xl p-6 border border-amber-100 flex items-start gap-4">
           <AlertCircleIcon class="w-6 h-6 text-amber-500 shrink-0 mt-0.5" />
           <div>
             <h3 class="font-bold text-slate-800">Emergency / Exam Interruption?</h3>
             <p class="text-slate-600 text-sm mt-1">If your exam was interrupted due to a technical issue, please inform the recruiter immediately or mention your token explicitly when contacting support.</p>
           </div>
        </div>
      </div>

      <div class="mt-10 pt-6 border-t border-slate-100 text-center">
        <RouterLink to="/" class="text-slate-600 font-medium hover:text-blue-600 transition-colors inline-block">
          &larr; Back to Login
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { HeadphonesIcon, MailIcon, MessageCircleIcon, AlertCircleIcon } from 'lucide-vue-next'
</script>
