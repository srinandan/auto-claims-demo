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
          @click="$router.push('/new-claim')"
          class="bg-teal-600 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded flex items-center"
        >
          <span>+ New Claim</span>
        </button>
      </div>
    </div>

    <!-- Policy Details Card -->
    <div v-if="user" class="bg-white rounded-lg shadow p-6 mb-8 border-l-4 border-teal-500">
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
            <div class="flex items-center gap-2">
              <span
                class="px-2 py-1 text-xs font-bold rounded-full uppercase tracking-wide"
                :class="statusBadgeClass(claim.status)"
              >
                {{ claim.status }}
              </span>
              <button
                @click.stop="deleteClaim(claim.ID)"
                class="text-gray-400 hover:text-red-600 p-1"
                title="Delete Claim"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
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

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const claims = ref([])
const loading = ref(true)
const user = ref(null)

const fetchClaims = async () => {
  if (!user.value) return

  try {
    const response = await axios.get('/api/claims', {
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

const deleteClaim = async (id) => {
  if (!confirm('Are you sure you want to permanently delete this claim?')) return

  try {
    await axios.delete(`/api/claims/${id}`)
    // Remove from list
    claims.value = claims.value.filter(c => c.ID !== id)
  } catch (error) {
    console.error('Error deleting claim:', error)
    alert('Failed to delete claim')
  }
}

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
      fetchClaims()
    } catch (e) {
      console.error('Invalid user data', e)
    }
  } else {
    loading.value = false
  }
})
</script>
