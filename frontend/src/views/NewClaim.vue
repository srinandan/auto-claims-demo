<template>
  <div class="p-6 max-w-5xl mx-auto">
    <!-- Header Section -->
    <div class="mb-12">
      <h1 class="text-3xl lg:text-4xl font-extrabold text-gray-800 tracking-tight mb-3">New Claim Submission</h1>
      <p class="text-gray-600 max-w-xl text-md leading-relaxed">
        Take a deep breath. Our concierge team is ready to guide you through this process. We'll handle the complexities so you can focus on what matters.
      </p>
    </div>

    <!-- Multi-step Progress Indicator -->
    <div class="flex items-center space-x-12 mb-12 overflow-x-auto pb-4 no-scrollbar">
      <div class="flex items-center space-x-4 shrink-0">
        <span class="w-10 h-10 rounded-full bg-teal-600 text-white flex items-center justify-center font-bold">1</span>
        <span class="font-bold text-gray-800">Incident Details</span>
      </div>
      <div class="flex items-center space-x-4 shrink-0 opacity-40">
        <span class="w-10 h-10 rounded-full bg-gray-200 text-gray-800 flex items-center justify-center font-bold">2</span>
        <span class="font-bold text-gray-800">Vehicle Info</span>
      </div>
      <div class="flex items-center space-x-4 shrink-0 opacity-40">
        <span class="w-10 h-10 rounded-full bg-gray-200 text-gray-800 flex items-center justify-center font-bold">3</span>
        <span class="font-bold text-gray-800">Evidence</span>
      </div>
      <div class="flex items-center space-x-4 shrink-0 opacity-40">
        <span class="w-10 h-10 rounded-full bg-gray-200 text-gray-800 flex items-center justify-center font-bold">4</span>
        <span class="font-bold text-gray-800">Final Review</span>
      </div>
    </div>

    <!-- Form Content -->
    <form @submit.prevent="submitClaim" class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      <!-- Left Column: Primary Fields -->
      <div class="lg:col-span-7 space-y-8">
        <div class="bg-white shadow p-8 rounded-xl border-l-4 border-teal-500">
          <h2 class="text-xl font-bold mb-6 flex items-center space-x-2 text-gray-800">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-teal-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span>When &amp; Where?</span>
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="flex flex-col space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-gray-500">Date of Incident</label>
              <input v-model="form.accidentDate" class="bg-gray-50 border border-gray-300 focus:ring-2 focus:ring-teal-500 focus:border-teal-500 p-3 rounded-md outline-none transition-all" type="date" required />
              <p class="text-[10px] text-gray-400">Choose the exact date the accident occurred.</p>
            </div>
            <div class="flex flex-col space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-gray-500">Approximate Time</label>
              <input v-model="form.approximateTime" class="bg-gray-50 border border-gray-300 focus:ring-2 focus:ring-teal-500 focus:border-teal-500 p-3 rounded-md outline-none transition-all" type="time" />
              <p class="text-[10px] text-gray-400">An estimate is fine if exact time is unknown.</p>
            </div>
            <div class="md:col-span-2 flex flex-col space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-gray-500">Location</label>
              <div class="relative">
                <input v-model="form.location" class="w-full bg-gray-50 border border-gray-300 focus:ring-2 focus:ring-teal-500 focus:border-teal-500 p-3 pl-10 rounded-md outline-none transition-all" placeholder="City, State" type="text" />
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white shadow p-8 rounded-xl border-l-4 border-teal-500">
          <h2 class="text-xl font-bold mb-6 flex items-center space-x-2 text-gray-800">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-teal-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>Incident Description</span>
          </h2>
          <div class="flex flex-col space-y-2">
            <label class="text-xs font-bold uppercase tracking-widest text-gray-500">Tell us what happened</label>
            <textarea v-model="form.description" class="bg-gray-50 border border-gray-300 focus:ring-2 focus:ring-teal-500 focus:border-teal-500 p-3 rounded-md outline-none transition-all resize-none" placeholder="Briefly describe the circumstances of the incident..." rows="5"></textarea>
            <p class="text-[10px] text-gray-400">Include details about road conditions, other vehicles involved, and any immediate actions taken.</p>
          </div>
        </div>
      </div>

      <!-- Right Column: Visual Evidence/Context -->
      <div class="lg:col-span-5 space-y-8">
        <div class="bg-gray-50 border-2 border-dashed border-gray-300 rounded-xl p-8 flex flex-col items-center text-center space-y-4 hover:bg-gray-100 transition-colors">
          <div class="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-teal-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <div>
            <h3 class="font-bold text-lg text-gray-800">Upload Accident Photos</h3>
            <p class="text-sm text-gray-500 mt-2">Drag and drop images or click to browse. Clear photos of damage help us process your claim faster.</p>
          </div>
          <div class="mt-4 w-full">
            <input type="file" multiple accept="image/*" @change="handleFiles" class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-teal-50 file:text-teal-700 hover:file:bg-teal-100" />
            <div v-if="form.files.length > 0" class="mt-2 text-sm text-gray-600 text-left">
                {{ form.files.length }} files selected
            </div>
          </div>
        </div>

        <div class="relative h-48 rounded-xl overflow-hidden shadow-sm group bg-gray-800 flex items-center justify-center">
             <!-- Placeholder for image in design, using a simple icon and text to represent the "Pro Tip" card -->
             <div class="absolute inset-0 bg-gradient-to-t from-black/80 to-black/30 flex items-end p-6">
                <div class="text-white">
                  <span class="bg-teal-200 text-teal-900 text-[10px] px-2 py-1 rounded font-bold uppercase tracking-tighter">Pro Tip</span>
                  <p class="text-sm mt-2 font-medium">Ensure everyone is safe before documenting damage.</p>
                </div>
            </div>
        </div>

        <!-- Trust Factor Card -->
        <div class="bg-gray-800 p-6 rounded-xl text-white">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-teal-400 mb-4" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <h4 class="font-bold text-lg mb-2">Secure &amp; Private</h4>
          <p class="text-sm text-gray-300 leading-relaxed">Your data is encrypted using banking-grade security standards. Our concierge team reviews each claim personally within 24 hours.</p>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="lg:col-span-12 mt-8 pt-8 border-t border-gray-200 flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
        <button type="button" @click="$router.push('/dashboard')" class="text-gray-500 font-bold flex items-center space-x-2 hover:text-gray-800 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>Cancel &amp; Exit</span>
        </button>
        <div class="flex items-center space-x-4 w-full md:w-auto">
          <button type="button" @click="$router.push('/dashboard')" class="flex-1 md:flex-none px-6 py-3 bg-gray-100 text-gray-700 rounded-md font-bold hover:bg-gray-200 transition-colors">Cancel</button>
          <button type="submit" :disabled="submitting" class="flex-1 md:flex-none px-8 py-3 bg-teal-600 text-white rounded-md font-bold shadow-md flex items-center justify-center space-x-2 hover:bg-teal-700 transition-colors disabled:opacity-50">
             <span v-if="submitting" class="mr-2">Submitting...</span>
             <template v-else>
                <span>Continue to Step 2</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
             </template>
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const user = ref(null)
const submitting = ref(false)

const form = reactive({
  policyNumber: '',
  accidentDate: '',
  approximateTime: '',
  location: '',
  description: '',
  files: []
})

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
      form.policyNumber = user.value.policy_number
    } catch (e) {
      console.error('Invalid user data', e)
    }
  }
})

const handleFiles = (event) => {
  form.files = Array.from(event.target.files)
}

const submitClaim = async () => {
  if (!form.policyNumber || !form.accidentDate) {
    alert('Please fill in required fields')
    return
  }

  submitting.value = true

  const formData = new FormData()
  formData.append('policy_number', form.policyNumber)
  formData.append('customer_name', user.value ? `${user.value.first_name} ${user.value.last_name}` : 'Valued Customer')

  let fullDescription = form.description
  if (form.approximateTime) {
      fullDescription = `[Time: ${form.approximateTime}] ${fullDescription}`
  }
  formData.append('description', fullDescription)
  formData.append('accident_date', form.accidentDate)

  let incidentCity = form.location
  let incidentState = ''
  if (form.location.includes(',')) {
      const parts = form.location.split(',')
      incidentCity = parts[0].trim()
      incidentState = parts[1].trim()
  }

  formData.append('incident_city', incidentCity)
  formData.append('incident_state', incidentState)
  formData.append('incident_type', '')
  formData.append('collision_type', '')
  formData.append('severity', '')

  for (let i = 0; i < form.files.length; i++) {
    formData.append('files', form.files[i])
  }

  try {
    await axios.post('/api/claims', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    router.push('/dashboard')
  } catch (error) {
    console.error('Error creating claim:', error)
    alert('Failed to submit claim. Please try again.')
  } finally {
    submitting.value = false
  }
}
</script>