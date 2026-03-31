<template>
  <nav class="fixed top-0 z-50 w-full bg-[#f9faf5] backdrop-blur-xl flex justify-between items-center px-8 py-4">
    <div class="flex items-center gap-8">
      <span class="text-xl font-extrabold text-[#0D1B2A] font-headline tracking-tight">Cymbal Insurance</span>
      <div class="hidden md:flex space-x-6">
        <router-link
          to="/dashboard"
          class="font-headline font-bold tracking-tight transition-colors"
          :class="$route.path === '/dashboard' ? 'text-[#000000] border-b-2 border-[#000000] pb-1' : 'text-[#44474c] hover:text-[#000000]'"
        >
          Dashboard
        </router-link>
        <router-link
          to="/dashboard"
          class="font-headline font-bold tracking-tight transition-colors text-[#44474c] hover:text-[#000000]"
        >
          My Claims
        </router-link>
        <router-link
          to="/dashboard"
          class="font-headline font-bold tracking-tight transition-colors text-[#44474c] hover:text-[#000000]"
        >
          Policies
        </router-link>
        <a class="text-[#44474c] hover:text-[#000000] transition-colors font-headline font-bold tracking-tight" href="#">
          Support
        </a>
      </div>
    </div>
    <div class="flex items-center space-x-5">
      <button class="material-symbols-outlined text-[#0D1B2A] p-2 hover:bg-[#edeeea] rounded-md transition-all active:scale-95">
        notifications
      </button>
      <button class="material-symbols-outlined text-[#0D1B2A] p-2 hover:bg-[#edeeea] rounded-md transition-all active:scale-95">
        help_outline
      </button>
      <div class="relative group">
        <div class="w-8 h-8 rounded-full bg-surface-container-highest overflow-hidden cursor-pointer">
          <img
            alt="User profile avatar"
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuBU19MZd10p888nT4bYo51vrMqHbyQt6a_vTh0bDyllbHyTPAiMOJlzXGVVobcfT_XB3AL65mdTURAjEGr-WkaOFecnQRZs_qmlrSpOHNma8wViU4eQGkSQD3OOoICV3vjUPdpP724vNDlOWmQlI6XGzYL2-E8hmZoOo3-K124XkYyiK_mjXwwt33d_LHsLB8pxQcTIR4Yy0le3KpvzsqR0ZNM7SuK5P5o52meRUkq7x28qiW9hCF73A9BBMh9zCfXX9_k4ZedDspw"
            class="w-full h-full object-cover"
          />
        </div>
        <div class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 hidden group-hover:block z-50">
          <button @click="handleLogout" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left">
            Logout
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const handleLogout = () => {
  localStorage.removeItem('user')
  router.push('/')
}
</script>
