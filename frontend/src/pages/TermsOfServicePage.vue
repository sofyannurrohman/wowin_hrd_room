<template>
  <div class="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto bg-white rounded-2xl shadow-sm p-8 md:p-12">
      <div class="mb-8">
        <h1 class="text-3xl font-extrabold text-slate-900 tracking-tight mb-2">Terms of Service</h1>
        <p class="text-slate-500">Last updated: {{ new Date().toLocaleDateString() }}</p>
      </div>

      <div class="prose prose-slate max-w-none">
        <p>Welcome to HRD Room!</p>

        <p>These terms and conditions outline the rules and regulations for the use of PT Wowin Purnomo Putera's Website and Application, located at our platform.</p>

        <p>By accessing this website we assume you accept these terms and conditions. Do not continue to use HRD Room if you do not agree to take all of the terms and conditions stated on this page.</p>

        <h2 class="text-xl font-bold mt-8 mb-4">License</h2>
        <p>Unless otherwise stated, PT Wowin Purnomo Putera and/or its licensors own the intellectual property rights for all material on HRD Room. All intellectual property rights are reserved. You may access this from HRD Room for your own personal use subjected to restrictions set in these terms and conditions.</p>

        <p>You must not:</p>
        <ul class="list-disc pl-5 mb-4 space-y-1 text-slate-700">
          <li>Republish material from HRD Room</li>
          <li>Sell, rent or sub-license material from HRD Room</li>
          <li>Reproduce, duplicate or copy material from HRD Room</li>
          <li>Redistribute content from HRD Room</li>
        </ul>

        <h2 class="text-xl font-bold mt-8 mb-4">Exam Integrity and Conduct</h2>
        <p>As a candidate taking an exam through the HRD Room platform, you agree to the following code of conduct:</p>
          <ul class="list-disc pl-5 mb-4 space-y-1 text-slate-700">
          <li>You will not use external assistance, including but not limited to search engines, other applications, or other individuals during the exam unless explicitly permitted.</li>
          <li>You consent to the monitoring of your exam session via webcam and browser activity tracking to ensure fairness and prevent cheating.</li>
          <li>You will not capture, copy, or distribute any exam questions or materials.</li>
          <li>Any attempt to manipulate the system or violate these rules may result in immediate disqualification and invalidation of your exam results.</li>
        </ul>

        <h2 class="text-xl font-bold mt-8 mb-4">Disclaimer</h2>
        <p>To the maximum extent permitted by applicable law, we exclude all representations, warranties and conditions relating to our website and the use of this website. Nothing in this disclaimer will:</p>
        <ul class="list-disc pl-5 mb-4 space-y-1 text-slate-700">
          <li>limit or exclude our or your liability for death or personal injury;</li>
          <li>limit or exclude our or your liability for fraud or fraudulent misrepresentation;</li>
          <li>limit any of our or your liabilities in any way that is not permitted under applicable law; or</li>
          <li>exclude any of our or your liabilities that may not be excluded under applicable law.</li>
        </ul>
        <p>The limitations and prohibitions of liability set in this Section and elsewhere in this disclaimer: (a) are subject to the preceding paragraph; and (b) govern all liabilities arising under the disclaimer, including liabilities arising in contract, in tort and for breach of statutory duty.</p>
        <p>As long as the website and the information and services on the website are provided free of charge, we will not be liable for any loss or damage of any nature.</p>
      </div>

      <div class="mt-12 pt-8 border-t border-slate-100">
        <RouterLink to="/" class="text-blue-600 font-medium hover:underline flex items-center gap-2">
          &larr; Back to Home
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
