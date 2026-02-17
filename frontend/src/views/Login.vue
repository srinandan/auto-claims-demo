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
  <div class="min-h-screen flex items-center justify-center bg-gray-100 px-4">
    <div class="max-w-md w-full bg-white rounded-lg shadow-md p-8">
      <h2 class="text-2xl font-bold text-center text-gray-800 mb-6">Customer Login</h2>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="policyNumber">
            Policy Number
          </label>
          <input
            v-model="policyNumber"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
            id="policyNumber"
            type="text"
            placeholder="Enter your Policy Number"
            required
          />
        </div>

        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="password">
            Password
          </label>
          <input
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
            id="password"
            type="password"
            placeholder="******************"
            value="dummy"
          />
          <p class="text-xs text-gray-500 mt-1">Any password will work.</p>
        </div>

        <div v-if="error" class="mb-4 text-red-500 text-sm text-center">
          {{ error }}
        </div>

        <div class="flex items-center justify-between">
          <button
            class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full transition duration-150 ease-in-out"
            type="submit"
            :disabled="loading"
          >
            {{ loading ? 'Signing In...' : 'Sign In' }}
          </button>
        </div>
      </form>

      <div class="mt-4 text-center">
        <p class="text-sm text-gray-600">
          Sample Policy: <span class="font-mono bg-gray-200 px-1 rounded">521585</span>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const router = useRouter()
const policyNumber = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!policyNumber.value) return

  loading.value = true
  error.value = ''

  try {
    const response = await api.get(`/api/policies/${policyNumber.value}`)
    // Store user session
    localStorage.setItem('user', JSON.stringify(response.data))
    router.push('/dashboard')
  } catch (err) {
    console.error(err)
    if (err.response && err.response.status === 404) {
      error.value = 'Policy number not found. Please check and try again.'
    } else {
      error.value = 'An error occurred. Please try again later.'
    }
  } finally {
    loading.value = false
  }
}
</script>
