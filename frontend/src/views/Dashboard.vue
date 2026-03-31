<template>
  <div class="px-6 lg:px-12 pb-20">
    <!-- Dashboard Header -->
    <header class="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12">
      <div>
        <span class="text-secondary font-label font-semibold text-xs uppercase tracking-[0.2em] mb-2 block" v-if="user">
          Welcome back, {{ user.first_name }}
        </span>
        <h1 class="text-4xl md:text-5xl font-headline font-extrabold text-on-surface tracking-tight">
          {{ activeClaim ? 'Your claim is moving forward.' : 'Your policies are active.' }}
        </h1>
      </div>
      <router-link
        to="/new-claim"
        class="group relative inline-flex items-center justify-center px-8 py-4 font-headline font-bold text-on-primary bg-gradient-to-br from-primary to-primary-container rounded-lg overflow-hidden transition-all hover:shadow-2xl active:scale-95"
      >
        <span class="relative flex items-center gap-2">
          <span class="material-symbols-outlined" style="font-variation-settings: 'FILL' 1;">add_circle</span>
          Start New Claim
        </span>
      </router-link>
    </header>

    <!-- Bento Grid Layout -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-6">

      <!-- Claim Status Hero Card (Large Asymmetric) -->
      <div v-if="activeClaim" class="md:col-span-8 bg-surface-container-lowest rounded-[2rem] p-1 lg:p-1.5 flex flex-col md:flex-row items-stretch overflow-hidden shadow-[0_12px_32px_rgba(16,27,48,0.04)] h-full">
        <div class="md:w-1/2 relative min-h-[300px] overflow-hidden rounded-[2rem]">
          <img
            alt="Vehicle detail"
            class="absolute inset-0 w-full h-full object-cover"
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuCVNkDdGOtmS6iIiUoNGCSvDSSHwhcBlnXQ7-nihGQc114f8VckZCDTdTo7WRKMjjPyysBoFOzmRmJbDzp0N_FiizQFJUp6MZMao1aHdpcV61ybELABmY0plrJ6poOYogw5s55zAhRlzVY2EiGIqAfVphJDm1LlSgWpbiVMLaIzk1kbT4GG9HJcHguQok_bOYIuGFWk5OeuROSjzo6JxIoaVlr3teGFZMqWNc_702ACcwpRyoossNZ7_e8EKy_RQHb2uUY2a5j9D2A"
          />
          <div class="absolute inset-0 bg-primary/20 backdrop-grayscale-[0.5]"></div>
          <div class="absolute bottom-8 left-8">
            <span class="inline-flex items-center gap-2 px-3 py-1 bg-tertiary-fixed text-on-tertiary-fixed-variant rounded-full text-xs font-bold font-label uppercase tracking-wider">
              <span class="material-symbols-outlined text-sm" style="font-variation-settings: 'FILL' 1;">check_circle</span>
              {{ activeClaim.status }}
            </span>
          </div>
        </div>
        <div class="md:w-1/2 p-8 lg:p-12 flex flex-col justify-center">
          <h3 class="text-2xl font-headline font-bold text-on-surface mb-4">Claim #{{ activeClaim.ID }}</h3>
          <p class="text-on-surface-variant leading-relaxed mb-8 font-body">
            Your recent claim from {{ formatDate(activeClaim.accident_date) }} is currently being verified by our Concierge team. We expect a resolution soon.
          </p>
          <div class="flex items-center gap-4">
            <div class="flex-1 h-1.5 bg-surface-container rounded-full overflow-hidden">
              <div class="h-full bg-primary rounded-full transition-all duration-500" :style="{ width: getProgressPercentage(activeClaim.status) + '%' }"></div>
            </div>
            <span class="text-xs font-bold text-primary">{{ getProgressPercentage(activeClaim.status) }}% Complete</span>
          </div>
          <router-link
            :to="`/claims/${activeClaim.ID}`"
            class="mt-8 text-secondary font-bold text-sm flex items-center gap-2 hover:text-primary transition-colors font-body"
          >
            View Details <span class="material-symbols-outlined text-sm">arrow_forward</span>
          </router-link>
        </div>
      </div>

      <!-- Policy Info Card (Fallback if no active claims) -->
      <div v-else-if="user" class="md:col-span-8 bg-surface-container-lowest rounded-[2rem] p-1 lg:p-1.5 flex flex-col md:flex-row items-stretch overflow-hidden shadow-[0_12px_32px_rgba(16,27,48,0.04)] h-full">
         <div class="w-full p-8 lg:p-12 flex flex-col justify-center">
            <h3 class="text-2xl font-headline font-bold text-on-surface mb-4">My Policy Details</h3>
            <div class="grid grid-cols-2 gap-6 mt-4">
              <div>
                <p class="text-xs font-bold text-on-surface-variant uppercase tracking-wide font-label mb-1">Policy Number</p>
                <p class="font-headline font-bold text-lg">{{ user.policy_number }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-on-surface-variant uppercase tracking-wide font-label mb-1">Insured Vehicle</p>
                <p class="font-headline font-bold text-lg">{{ user.auto_year }} {{ user.auto_make }} {{ user.auto_model }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-on-surface-variant uppercase tracking-wide font-label mb-1">Annual Premium</p>
                <p class="font-headline font-bold text-lg">${{ user.policy_annual_premium }}</p>
              </div>
              <div>
                <p class="text-xs font-bold text-on-surface-variant uppercase tracking-wide font-label mb-1">Deductible</p>
                <p class="font-headline font-bold text-lg">${{ user.policy_deductible }}</p>
              </div>
            </div>
         </div>
      </div>

      <!-- Recent Activity Panel -->
      <div class="md:col-span-4 bg-surface-container-low rounded-[1.5rem] p-8 flex flex-col">
        <div class="flex items-center justify-between mb-8">
          <h3 class="text-lg font-headline font-bold">Recent Updates</h3>
          <span class="material-symbols-outlined text-on-surface-variant">history</span>
        </div>

        <div v-if="loading" class="text-center text-sm py-4">Loading...</div>
        <div v-else-if="claims.length === 0" class="text-center text-sm text-on-surface-variant py-4 font-body">No recent activity.</div>

        <div v-else class="space-y-6 overflow-y-auto pr-2 flex-1">
          <div v-for="(claim, index) in claims.slice(0, 3)" :key="'update-'+claim.ID" class="flex gap-4 items-start">
            <div class="mt-1 w-2 h-2 rounded-full flex-shrink-0" :class="index === 0 ? 'bg-primary' : 'bg-outline-variant'"></div>
            <div>
              <p class="text-sm font-semibold font-body">Claim {{ claim.status }}</p>
              <p class="text-xs text-on-surface-variant font-body mt-1">Status updated for Claim #{{ claim.ID }}</p>
            </div>
          </div>
        </div>
        <button class="mt-8 w-full py-3 text-xs font-bold text-secondary uppercase tracking-widest font-label border border-outline-variant/20 rounded-lg hover:bg-surface-container transition-colors">
          All notifications
        </button>
      </div>

      <!-- Existing Claims Table (Modern List) -->
      <div class="md:col-span-12 mt-8">
        <div class="flex items-center justify-between mb-8">
          <h2 class="text-2xl font-headline font-bold">Claim History</h2>
          <div class="flex gap-2">
            <button class="px-4 py-2 bg-surface-container-lowest text-xs font-bold rounded-lg shadow-sm border border-outline-variant/5 hover:bg-surface-container transition-colors">All</button>
          </div>
        </div>

        <div v-if="loading" class="text-center py-10 text-on-surface-variant">
          Loading claims...
        </div>
        <div v-else-if="claims.length === 0" class="text-center py-10 bg-surface-container-lowest rounded-xl shadow-sm text-on-surface-variant">
          No claims found.
        </div>
        <div v-else class="bg-surface-container-lowest rounded-[1.5rem] overflow-hidden shadow-[0_12px_32px_rgba(16,27,48,0.02)]">
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-surface-container-low/50">
                  <th class="px-8 py-5 text-[10px] font-bold uppercase tracking-widest text-secondary font-label">Claim Reference</th>
                  <th class="px-8 py-5 text-[10px] font-bold uppercase tracking-widest text-secondary font-label">Vehicle</th>
                  <th class="px-8 py-5 text-[10px] font-bold uppercase tracking-widest text-secondary font-label">Date filed</th>
                  <th class="px-8 py-5 text-[10px] font-bold uppercase tracking-widest text-secondary font-label">Est. Settlement</th>
                  <th class="px-8 py-5 text-[10px] font-bold uppercase tracking-widest text-secondary font-label">Status</th>
                  <th class="px-8 py-5"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-outline-variant/10">
                <tr
                  v-for="claim in claims"
                  :key="claim.ID"
                  class="group hover:bg-surface-container-low/30 transition-colors cursor-pointer"
                  @click="$router.push(`/claims/${claim.ID}`)"
                >
                  <td class="px-8 py-6 font-semibold text-sm font-body">#{{ claim.ID }}</td>
                  <td class="px-8 py-6 text-sm text-on-surface-variant font-body">
                    {{ user?.auto_year }} {{ user?.auto_make }} {{ user?.auto_model }}
                  </td>
                  <td class="px-8 py-6 text-sm text-on-surface-variant font-body">{{ formatDate(claim.accident_date) }}</td>
                  <td class="px-8 py-6 text-sm font-headline font-bold">
                    {{ claim.repair_estimate ? '$'+claim.repair_estimate : 'Pending' }}
                  </td>
                  <td class="px-8 py-6">
                    <span
                      class="inline-block px-2.5 py-1 text-[10px] font-bold rounded-full uppercase tracking-wider"
                      :class="statusBadgeClass(claim.status)"
                    >
                      {{ claim.status }}
                    </span>
                  </td>
                  <td class="px-8 py-6 text-right">
                    <button class="material-symbols-outlined text-secondary hover:text-primary p-1" @click.stop="$router.push(`/claims/${claim.ID}`)">
                      chevron_right
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const claims = ref([])
const user = ref(null)
const loading = ref(true)

// Computed property to find the most recent active claim for the hero card
const activeClaim = computed(() => {
  if (claims.value.length === 0) return null
  // Simple check for claims that are not Approved or Rejected
  const active = claims.value.find(c => !['Approved', 'Rejected'].includes(c.status))
  return active || claims.value[0] // fallback to most recent
})

const fetchClaims = async () => {
  try {
    const response = await axios.get(`/api/claims?policy_number=${user.value.policy_number}`)
    // Assume API returns array, order descending by ID
    claims.value = response.data.sort((a, b) => b.ID - a.ID)
  } catch (error) {
    console.error('Error fetching claims:', error)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

const statusBadgeClass = (status) => {
  switch (status) {
    case 'New': return 'bg-tertiary-fixed text-on-tertiary-fixed-variant'
    case 'Assessed': return 'bg-secondary-fixed text-on-secondary-fixed'
    case 'Approved': return 'bg-surface-container-highest text-secondary'
    case 'Rejected': return 'bg-error-container text-on-error-container'
    case 'Total Loss': return 'bg-tertiary text-white'
    default: return 'bg-surface-container-highest text-secondary'
  }
}

const getProgressPercentage = (status) => {
  switch (status) {
    case 'New': return 25
    case 'Assessed': return 65
    case 'Approved': return 100
    case 'Rejected': return 100
    case 'Total Loss': return 100
    default: return 0
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
    router.push('/')
  }
})
</script>
