<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Claims Dashboard</h1>
      <div>
        <button
          @click="showModal = true"
          class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded flex items-center"
        >
          <span>+ New Claim</span>
        </button>
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
          <p class="text-sm text-gray-600 mb-4">Policy: {{ claim.policy_number }}</p>

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
    <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full flex justify-center items-center z-50">
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
                  @blur="verifyPolicy"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                  id="policyNumber"
                  type="text"
                  placeholder="Enter Policy Number"
                  required
                />
              </div>
              <button
                type="button"
                @click="verifyPolicy"
                class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mb-[1px]"
                :disabled="verifying"
              >
                {{ verifying ? 'Verifying...' : 'Verify' }}
              </button>
            </div>
            <div v-if="policyVerified" class="mt-2 text-green-600 text-sm font-semibold">
              ✓ Policy Verified: {{ policyDetails.auto_year }} {{ policyDetails.auto_make }} {{ policyDetails.auto_model }}
            </div>
            <div v-if="policyError" class="mt-2 text-red-600 text-sm font-semibold">
              ⚠ {{ policyError }}
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
import axios from 'axios'

const claims = ref([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const verifying = ref(false)
const policyVerified = ref(false)
const policyError = ref('')
const policyDetails = ref({})

const form = reactive({
  policyNumber: '',
  customerName: '', // Will be derived or manual
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
  try {
    const response = await axios.get('/api/claims')
    claims.value = response.data
  } catch (error) {
    console.error('Error fetching claims:', error)
  } finally {
    loading.value = false
  }
}

const verifyPolicy = async () => {
  if (!form.policyNumber) return
  verifying.value = true
  policyVerified.value = false
  policyError.value = ''
  policyDetails.value = {}

  try {
    const response = await axios.get(`/api/policies/${form.policyNumber}`)
    policyDetails.value = response.data
    policyVerified.value = true
    // Auto-fill some data if needed, or just show confirmation
    form.customerName = "Insured Customer" // Placeholder as CSV has no name
  } catch (error) {
    console.error('Policy verification failed:', error)
    policyError.value = 'Policy not found. Please check the number.'
  } finally {
    verifying.value = false
  }
}

const handleFiles = (event) => {
  form.files = Array.from(event.target.files)
}

const closeModal = () => {
  showModal.value = false
  resetForm()
}

const resetForm = () => {
  form.policyNumber = ''
  form.customerName = ''
  form.accidentDate = ''
  form.description = ''
  form.incidentCity = ''
  form.incidentState = ''
  form.incidentType = ''
  form.collisionType = ''
  form.severity = ''
  form.files = []
  policyVerified.value = false
  policyError.value = ''
  policyDetails.value = {}
}

const submitClaim = async () => {
  if (!form.policyNumber || !form.accidentDate) {
    alert('Please fill in required fields')
    return
  }

  submitting.value = true

  const formData = new FormData()
  formData.append('policy_number', form.policyNumber)
  formData.append('customer_name', form.customerName || 'Valued Customer')
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
    await axios.post('/api/claims', formData, {
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
  fetchClaims()
})
</script>
