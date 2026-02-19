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
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-teal-600 text-white p-4 shadow-md">
      <div class="container mx-auto flex justify-between items-center">
        <router-link to="/" class="font-bold text-xl flex items-center gap-2">
          <img src="/cymbal-logo.svg" alt="Cymbal Logo" class="h-8 w-8" />
          <span>Cymbal Insurance</span>
        </router-link>

        <div v-if="user" class="flex items-center gap-4">
          <span class="text-sm font-medium opacity-90">
            Welcome, {{ user.first_name }} {{ user.last_name }}
          </span>
          <button
            @click="logout"
            class="bg-teal-700 hover:bg-teal-800 text-white px-3 py-1 rounded text-sm transition-colors"
          >
            Logout
          </button>
        </div>
      </div>
    </nav>
    <main class="container mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const user = ref(null)

const checkUser = () => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
    } catch (e) {
      user.value = null
    }
  } else {
    user.value = null
  }
}

watch(() => route.path, () => {
  checkUser()
})

onMounted(() => {
  checkUser()
})

const logout = () => {
  localStorage.removeItem('user')
  user.value = null
  router.push('/')
}
</script>
