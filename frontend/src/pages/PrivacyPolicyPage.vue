<template>
  <div class="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto bg-white rounded-2xl shadow-sm p-8 md:p-12">
      <div class="mb-8">
        <h1 class="text-3xl font-extrabold text-slate-900 tracking-tight mb-2">Privacy Policy</h1>
        <p class="text-slate-500">Last updated: {{ new Date().toLocaleDateString() }}</p>
      </div>

      <div class="prose prose-slate max-w-none">
        <p>At Wowin Assesment Center, accessible from our application, one of our main priorities is the privacy of our visitors and users. This Privacy Policy document contains types of information that is collected and recorded by Wowin Assesment Center and how we use it.</p>
        
        <p>If you have additional questions or require more information about our Privacy Policy, do not hesitate to contact us.</p>
        
        <h2 class="text-xl font-bold mt-8 mb-4">Information we collect</h2>
        <p>The personal information that you are asked to provide, and the reasons why you are asked to provide it, will be made clear to you at the point we ask you to provide your personal information.</p>
        <p>If you contact us directly, we may receive additional information about you such as your name, email address, phone number, the contents of the message and/or attachments you may send us, and any other information you may choose to provide.</p>
        <p>When you register for an Account or a testing session, we may ask for your contact information, including items such as name, email address, and age.</p>

        <h2 class="text-xl font-bold mt-8 mb-4">How we use your information</h2>
        <p>We use the information we collect in various ways, including to:</p>
        <ul class="list-disc pl-5 mb-4 space-y-1 text-slate-700">
          <li>Provide, operate, and maintain our website and application</li>
          <li>Improve, personalize, and expand our application</li>
          <li>Understand and analyze how you use our application</li>
          <li>Develop new products, services, features, and functionality</li>
          <li>Communicate with you, either directly or through one of our partners, including for customer service, to provide you with updates and other information relating to the platform</li>
          <li>Send you emails</li>
          <li>Find and prevent fraud</li>
        </ul>

        <h2 class="text-xl font-bold mt-8 mb-4">Log Files</h2>
        <p>Wowin Assesment Center follows a standard procedure of using log files. These files log visitors when they visit the application. All hosting companies do this and a part of hosting services' analytics. The information collected by log files include internet protocol (IP) addresses, browser type, Internet Service Provider (ISP), date and time stamp, referring/exit pages, and possibly the number of clicks. These are not linked to any information that is personally identifiable. The purpose of the information is for analyzing trends, administering the site, tracking users' movement on the application, and gathering demographic information.</p>

        <h2 class="text-xl font-bold mt-8 mb-4">Data Security during Exams</h2>
        <p>For participants taking exams via Wowin Assesment Center, we employ additional security measures such as webcam monitoring and browser focus tracking to ensure exam integrity. This data is reviewed solely for the purpose of validating the exam attempt and is accessible only to authorized HR personnel reviewing your application.</p>
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
