<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-blue-600 text-white p-4 shadow-md">
      <div class="container mx-auto flex justify-between items-center">
        <router-link to="/" class="font-bold text-xl flex items-center gap-2">
          <span>Claims Intake System</span>
        </router-link>

        <div v-if="user" class="flex items-center gap-4">
          <span class="text-sm font-medium opacity-90">
            Welcome, {{ user.first_name }} {{ user.last_name }}
          </span>
          <button
            @click="logout"
            class="bg-blue-700 hover:bg-blue-800 text-white px-3 py-1 rounded text-sm transition-colors"
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
