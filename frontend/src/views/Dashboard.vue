<!--
 Copyright 2026 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Claims Dashboard</h1>
      <div>
        <button
          @click="openModal"
          class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded flex items-center"
        >
          <span>+ New Claim</span>
        </button>
      </div>
    </div>

    <!-- Policy Details Card -->
    <div v-if="user" class="bg-white rounded-lg shadow p-6 mb-8 border-l-4 border-blue-500">
      <h2 class="text-xl font-bold text-gray-800 mb-4">My Policy Details</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div>
          <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Policy Number</p>
          <p class="font-semibold text-lg">{{ user.policy_number }}</p>
        </div>
        <div>
          <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Insured Vehicle</p>
          <p class="font-semibold text-lg">{{ user.auto_year }} {{ user.auto_make }} {{ user.auto_model }}</p>
        </div>
        <div>
          <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Annual Premium</p>
          <p class="font-semibold text-lg">${{ user.policy_annual_premium }}</p>
        </div>
        <div>
          <p class="text-xs font-bold text-gray-500 uppercase tracking-wide">Deductible</p>
          <p class="font-semibold text-lg">${{ user.policy_deductible }}</p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="text-center py-10">
      <p class="text-gray-500">Loading claims...</p>
    </div>

    <div v-else-if="claims.length === 0" class="text-center py-10 bg-white rounded shadow">
      <p class="text-gray-500">No claims found. Create one to get started.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="claim in claims"
        :key="claim.ID"
        class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow duration-200 cursor-pointer border-l-4"
        :class="statusColorClass(claim.status)"
        @click="$router.push(`/claims/${claim.ID}`)"
      >
        <div class="p-5">
          <div class="flex justify-between items-start mb-2">
            <span class="text-xs font-semibold text-gray-500">#{{ claim.ID }}</span>
            <span
              class="px-2 py-1 text-xs font-bold rounded-full uppercase tracking-wide"
              :class="statusBadgeClass(claim.status)"
            >
              {{ claim.status }}
            </span>
          </div>
          <h2 class="text-xl font-bold text-gray-800 mb-1">{{ claim.customer_name }}</h2>
          <p class="text-sm text-gray-600 mb-4">Date: {{ formatDate(claim.accident_date) }}</p>

          <div class="flex justify-between items-center text-sm text-gray-500">
            <span>{{ claim.photos ? claim.photos.length : 0 }} Photos</span>
            <span v-if="claim.action_required" class="text-red-600 font-bold flex items-center">
              <span class="mr-1">⚠</span> Action Required
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- New Claim Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-gray-600/50 overflow-y-auto h-full w-full flex justify-center items-center z-50">
      <div class="bg-white p-8 rounded-lg shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-800">File New Claim</h2>
          <button @click="closeModal" class="text-gray-500 hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="submitClaim">
          <!-- Policy Information -->
          <div class="mb-6 border-b pb-4">
            <h3 class="text-lg font-semibold text-gray-700 mb-4">Policy Information</h3>
            <div class="flex gap-4 items-end">
              <div class="flex-grow">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="policyNumber">
                  Policy Number
                </label>
                <input
                  v-model="form.policyNumber"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-100 cursor-not-allowed"
                  id="policyNumber"
                  type="text"
                  required
                  disabled
                />
              </div>
            </div>
          </div>

          <!-- Incident Details -->
          <div class="mb-6 border-b pb-4">
            <h3 class="text-lg font-semibold text-gray-700 mb-4">Incident Details</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label class="block text-gray-700 text-sm font-bold mb-2" for="accidentDate">
                  Date of Accident
                </label>
                <input
                  v-model="form.accidentDate"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="accidentDate"
                  type="date"
                  required
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-bold mb-2" for="incidentType">
                  Incident Type
                </label>
                <select
                  v-model="form.incidentType"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="incidentType"
                >
                  <option value="">Select Type</option>
                  <option value="Single Vehicle Collision">Single Vehicle Collision</option>
                  <option value="Multi-vehicle Collision">Multi-vehicle Collision</option>
                  <option value="Vehicle Theft">Vehicle Theft</option>
                  <option value="Parked Car">Parked Car</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label class="block text-gray-700 text-sm font-bold mb-2" for="incidentCity">
                  City
                </label>
                <input
                  v-model="form.incidentCity"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="incidentCity"
                  type="text"
                  placeholder="City"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-bold mb-2" for="incidentState">
                  State
                </label>
                <input
                  v-model="form.incidentState"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="incidentState"
                  type="text"
                  placeholder="State (e.g. NY)"
                />
              </div>
            </div>

             <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label class="block text-gray-700 text-sm font-bold mb-2" for="collisionType">
                  Collision Type
                </label>
                <select
                  v-model="form.collisionType"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="collisionType"
                >
                  <option value="">Select Collision</option>
                  <option value="Side Collision">Side Collision</option>
                  <option value="Rear Collision">Rear Collision</option>
                  <option value="Front Collision">Front Collision</option>
                </select>
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-bold mb-2" for="severity">
                  Severity
                </label>
                <select
                  v-model="form.severity"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="severity"
                >
                  <option value="">Select Severity</option>
                  <option value="Minor Damage">Minor Damage</option>
                  <option value="Major Damage">Major Damage</option>
                  <option value="Total Loss">Total Loss</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Description -->
          <div class="mb-6 border-b pb-4">
            <label class="block text-gray-700 text-sm font-bold mb-2" for="description">
              Description
            </label>
            <textarea
              v-model="form.description"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="description"
              rows="4"
              placeholder="Describe what happened..."
            ></textarea>
          </div>

          <!-- Photos -->
          <div class="mb-6">
            <label class="block text-gray-700 text-sm font-bold mb-2" for="photos">
              Upload Photos
            </label>
            <input
              type="file"
              multiple
              accept="image/*"
              @change="handleFiles"
              class="block w-full text-sm text-gray-500
                file:mr-4 file:py-2 file:px-4
                file:rounded-full file:border-0
                file:text-sm file:font-semibold
                file:bg-blue-50 file:text-blue-700
                hover:file:bg-blue-100"
              id="photos"
            />
            <div v-if="form.files.length > 0" class="mt-2 text-sm text-gray-600">
              {{ form.files.length }} files selected
            </div>
          </div>

          <div class="flex items-center justify-end">
            <button
              type="button"
              @click="closeModal"
              class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded mr-4"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline flex items-center"
              :disabled="submitting"
            >
              <span v-if="submitting" class="mr-2">Submitting...</span>
              <span v-else>Submit Claim</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import api from '../api'

const claims = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const user = ref(null)

// Form state
const form = reactive({
  policyNumber: '',
  accidentDate: '',
  description: '',
  incidentCity: '',
  incidentState: '',
  incidentType: '',
  collisionType: '',
  severity: '',
  files: []
})

const fetchClaims = async () => {
  if (!user.value) return

  try {
    const response = await api.get('/api/claims', {
      params: {
        policy_number: user.value.policy_number
      }
    })
    claims.value = response.data
  } catch (error) {
    console.error('Error fetching claims:', error)
  } finally {
    loading.value = false
  }
}

const openModal = () => {
  resetForm()
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  resetForm()
}

const resetForm = () => {
  if (user.value) {
    form.policyNumber = user.value.policy_number
  }
  form.accidentDate = ''
  form.description = ''
  form.incidentCity = ''
  form.incidentState = ''
  form.incidentType = ''
  form.collisionType = ''
  form.severity = ''
  form.files = []
}

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
  formData.append('description', form.description)
  formData.append('accident_date', form.accidentDate)
  formData.append('incident_city', form.incidentCity)
  formData.append('incident_state', form.incidentState)
  formData.append('incident_type', form.incidentType)
  formData.append('collision_type', form.collisionType)
  formData.append('severity', form.severity)

  for (let i = 0; i < form.files.length; i++) {
    formData.append('files', form.files[i])
  }

  try {
    await api.post('/api/claims', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    closeModal()
    fetchClaims()
  } catch (error) {
    console.error('Error creating claim:', error)
    alert('Failed to submit claim. Please try again.')
  } finally {
    submitting.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

const statusColorClass = (status) => {
  switch (status) {
    case 'Simple': return 'border-green-500'
    case 'Complex': return 'border-yellow-500'
    case 'Total Loss': return 'border-red-500'
    default: return 'border-gray-300'
  }
}

const statusBadgeClass = (status) => {
  switch (status) {
    case 'Simple': return 'bg-green-100 text-green-800'
    case 'Complex': return 'bg-yellow-100 text-yellow-800'
    case 'Total Loss': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
      form.policyNumber = user.value.policy_number
      fetchClaims()
    } catch (e) {
      console.error('Invalid user data', e)
    }
  } else {
    loading.value = false
  }
})
</script>
