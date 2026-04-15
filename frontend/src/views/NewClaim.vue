<template>
  <section class="max-w-5xl mx-auto px-6 py-12 lg:py-20">
    <!-- Header Section: Atmospheric Statement -->
    <div class="mb-16">
      <h1 class="text-4xl lg:text-5xl font-extrabold text-on-surface tracking-tight mb-4 font-headline">New Claim Submission</h1>
      <p class="text-on-surface-variant max-w-xl text-lg leading-relaxed font-body">
          Take a deep breath. Our concierge team is ready to guide you through this process. We'll handle the complexities so you can focus on what matters.
      </p>
    </div>

    <!-- Multi-step Progress Indicator -->
    <div class="flex items-center space-x-12 mb-20 overflow-x-auto pb-4 no-scrollbar">
      <div class="flex items-center space-x-4 shrink-0">
          <span class="w-10 h-10 rounded-full bg-primary text-on-primary flex items-center justify-center font-bold">1</span>
          <span class="font-headline font-bold text-on-surface">Incident Details</span>
      </div>
      <div class="flex items-center space-x-4 shrink-0 opacity-40">
          <span class="w-10 h-10 rounded-full bg-surface-container-highest text-on-surface flex items-center justify-center font-bold">2</span>
          <span class="font-headline font-bold">Vehicle Info</span>
      </div>
      <div class="flex items-center space-x-4 shrink-0 opacity-40">
          <span class="w-10 h-10 rounded-full bg-surface-container-highest text-on-surface flex items-center justify-center font-bold">3</span>
          <span class="font-headline font-bold">Evidence</span>
      </div>
      <div class="flex items-center space-x-4 shrink-0 opacity-40">
          <span class="w-10 h-10 rounded-full bg-surface-container-highest text-on-surface flex items-center justify-center font-bold">4</span>
          <span class="font-headline font-bold">Final Review</span>
      </div>
    </div>

    <form @submit.prevent="submitClaim">
      <!-- Form Content: Bento-style Layout -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">

        <!-- Left Column: Primary Fields -->
        <div class="lg:col-span-7 space-y-12">
            <div class="bg-surface-container-low p-8 rounded-[1.5rem]">
              <h2 class="text-xl font-bold mb-8 flex items-center space-x-2 font-headline">
                  <span class="material-symbols-outlined text-secondary">event_note</span>
                  <span>When & Where?</span>
              </h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div class="flex flex-col space-y-2">
                    <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Date of Incident</label>
                    <input
                      v-model="form.date"
                      type="date"
                      required
                      class="bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 rounded-md outline-none transition-all font-body"
                    />
                    <p class="text-[10px] text-on-surface-variant font-body">Choose the exact date the accident occurred.</p>
                  </div>
                  <div class="flex flex-col space-y-2">
                    <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Approximate Time</label>
                    <input
                      type="time"
                      class="bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 rounded-md outline-none transition-all font-body"
                    />
                    <p class="text-[10px] text-on-surface-variant font-body">An estimate is fine if exact time is unknown.</p>
                  </div>
                  <div class="md:col-span-2 flex flex-col space-y-2 relative">
                    <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Location</label>
                    <div class="relative">
                        <input
                          v-model="form.location"
                          @input="onLocationInput"
                          type="text"
                          placeholder="Street, City, State"
                          required
                          class="w-full bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 pl-12 rounded-md outline-none transition-all font-body"
                        />
                        <span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-secondary">location_on</span>
                    </div>

                    <!-- Address Suggestions Dropdown -->
                    <ul v-if="addressSuggestions.length > 0" class="absolute top-full left-0 right-0 mt-1 bg-surface-container-lowest border border-outline-variant/15 rounded-md shadow-lg z-10 max-h-60 overflow-y-auto">
                        <li
                            v-for="(suggestion, index) in addressSuggestions"
                            :key="index"
                            @click="selectAddress(suggestion)"
                            class="p-3 hover:bg-surface-container-low cursor-pointer font-body text-sm border-b border-outline-variant/10 last:border-0"
                        >
                            {{ suggestion.formattedAddress }}
                        </li>
                    </ul>
                  </div>
              </div>
            </div>

            <div class="bg-surface-container-low p-8 rounded-[1.5rem]">
              <h2 class="text-xl font-bold mb-8 flex items-center space-x-2 font-headline">
                  <span class="material-symbols-outlined text-secondary">description</span>
                  <span>Incident Description</span>
              </h2>
              <div class="flex flex-col space-y-2">
                  <label class="text-xs font-bold uppercase tracking-widest text-secondary font-label">Tell us what happened</label>
                  <textarea
                    v-model="form.description"
                    rows="5"
                    required
                    placeholder="Briefly describe the circumstances of the incident..."
                    class="bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 rounded-md outline-none transition-all resize-none font-body"
                  ></textarea>
                  <p class="text-[10px] text-on-surface-variant font-body">Include details about road conditions, other vehicles involved, and any immediate actions taken.</p>
              </div>
            </div>
        </div>

        <!-- Right Column: Visual Evidence/Context -->
        <div class="lg:col-span-5 space-y-8">
            <div
              class="bg-surface-container-lowest border-2 border-dashed border-outline-variant/30 rounded-[1.5rem] p-10 flex flex-col items-center text-center space-y-6 hover:bg-surface-container-high/30 transition-colors cursor-pointer relative"
            >
              <input
                type="file"
                @change="handleFileUpload"
                multiple
                accept="image/*"
                required
                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
              />
              <div class="w-16 h-16 bg-surface-container-highest rounded-full flex items-center justify-center pointer-events-none">
                  <span class="material-symbols-outlined text-primary text-3xl">add_a_photo</span>
              </div>
              <div class="pointer-events-none">
                  <h3 class="font-bold text-lg font-headline">Upload Accident Photos</h3>
                  <p class="text-sm text-on-surface-variant mt-2 px-4 font-body">Drag and drop images or click to browse. Clear photos of damage help us process your claim faster.</p>
              </div>
              <button type="button" class="bg-surface-container-highest text-on-surface px-6 py-2 rounded-md font-bold text-sm pointer-events-none font-headline">
                Select Files
              </button>

              <!-- File Preview List -->
              <div v-if="form.photos.length > 0" class="mt-4 w-full">
                <p class="text-sm font-bold text-secondary mb-2">{{ form.photos.length }} files selected:</p>
                <ul class="text-xs text-left text-on-surface-variant space-y-1 max-h-32 overflow-y-auto">
                  <li v-for="(file, index) in form.photos" :key="index" class="truncate">
                    {{ file.name }}
                  </li>
                </ul>
              </div>
            </div>

            <div class="relative h-64 rounded-[1.5rem] overflow-hidden shadow-sm group">
              <img
                alt="Safety first visual"
                class="w-full h-full object-cover"
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuA_gmVTAzDAxhgITTVNHd-fBYDoieXEBDs2V0sO_CXF2TevFXw9vsGfLTKGDrtSQSgO_2ja5cWxwEfOM5vdz0zo0mO39934RxfS46UHKXe2Qh3uweRSqqRQc6hYtRhWmhn1wmrfuIkdbD-vXPqUinftQPGZwvKknElnutpVM0UJ4uEFjRxA2c5uM6I-kjP4QOjXQ0826s1Hfbcyta5oJk5Ktnrl5hEdsDp1OF4CyobzLmZWWD7cGzLHP9dfRLwdcxTdWc6U1EWmxTY"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent flex items-end p-6">
                  <div class="text-white">
                    <span class="bg-tertiary-fixed text-on-tertiary-fixed text-[10px] px-2 py-1 rounded font-bold uppercase tracking-tighter font-label">Pro Tip</span>
                    <p class="text-sm mt-2 font-medium font-body">Ensure everyone is safe before documenting damage.</p>
                  </div>
              </div>
            </div>

            <!-- Trust Factor Card -->
            <div class="bg-primary-container p-8 rounded-[1.5rem] text-on-primary-container">
              <span class="material-symbols-outlined text-tertiary-fixed text-4xl mb-4" style="font-variation-settings: 'FILL' 1;">verified_user</span>
              <h4 class="text-white font-bold text-lg mb-2 font-headline">Secure & Private</h4>
              <p class="text-sm opacity-80 leading-relaxed font-body">Your data is encrypted using banking-grade security standards. Our concierge team reviews each claim personally within 24 hours.</p>
            </div>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="mt-16 pt-12 border-t border-outline-variant/15 flex flex-col md:flex-row justify-between items-center space-y-6 md:space-y-0">
        <router-link to="/dashboard" class="text-secondary font-bold flex items-center space-x-2 hover:text-primary transition-colors font-headline">
            <span class="material-symbols-outlined">arrow_back</span>
            <span>Save for later & Exit</span>
        </router-link>
        <div class="flex items-center space-x-4 w-full md:w-auto">
            <router-link to="/dashboard" class="flex-1 md:flex-none px-8 py-4 bg-surface-container-highest text-on-surface rounded-md font-bold hover:bg-surface-dim transition-colors text-center font-headline">
              Cancel
            </router-link>
            <button
              type="submit"
              :disabled="submitting"
              class="flex-1 md:flex-none px-12 py-4 bg-gradient-to-br from-primary to-primary-container text-on-primary rounded-md font-bold shadow-lg flex items-center justify-center space-x-2 hover:opacity-90 transition-opacity font-headline disabled:opacity-50"
            >
              <span>{{ submitting ? 'Submitting...' : 'Submit Claim' }}</span>
              <span class="material-symbols-outlined" v-if="!submitting">arrow_forward</span>
            </button>
        </div>
      </div>
    </form>
  </section>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref(null)
const submitting = ref(false)
const addressSuggestions = ref([])
let debounceTimeout = null

const form = reactive({
  policyNumber: '',
  customerName: '',
  date: '',
  description: '',
  location: '',
  photos: []
})

const fetchAddressSuggestions = async (query) => {
  if (!query || query.length < 3) {
    addressSuggestions.value = []
    return
  }

  try {
    const response = await axios.post('/api/resolve-address', { address: query })

    // Parse Google Maps Grounding Lite MCP response
    if (response.data && response.data.places) {
        addressSuggestions.value = response.data.places.map(place => {
            // Handle various possible response formats from MCP
            const address = place.formattedAddress || place.formatted_address || (place.displayName && place.displayName.text) || place.summary || "Unknown Address"
            return {
                formattedAddress: address,
                placeId: place.id || place.placeId || place.place_id,
            }
        })
    } else {
        addressSuggestions.value = []
    }
  } catch (error) {
    console.error('Error fetching address suggestions:', error)
    addressSuggestions.value = []
  }
}

const onLocationInput = () => {
  clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    fetchAddressSuggestions(form.location)
  }, 300)
}

const selectAddress = (suggestion) => {
  form.location = suggestion.formattedAddress
  addressSuggestions.value = []
}

const handleFileUpload = (event) => {
  form.photos = Array.from(event.target.files)
}

const submitClaim = async () => {
  submitting.value = true
  const formData = new FormData()
  formData.append('policy_number', form.policyNumber)
  formData.append('customer_name', form.customerName)
  formData.append('accident_date', form.date)
  formData.append('description', form.description)
  // We send the full location as incident_city to avoid DB schema changes,
  // or we can just send it and let the backend save it.
  formData.append('incident_city', form.location)

  form.photos.forEach(photo => {
    formData.append('files', photo)
  })

  try {
    await axios.post('/api/claims', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    router.push('/dashboard')
  } catch (error) {
    console.error('Error submitting claim:', error)
    alert('Failed to submit claim. Please try again.')
    submitting.value = false
  }
}

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
      form.policyNumber = user.value.policy_number
      form.customerName = `${user.value.first_name} ${user.value.last_name}`
    } catch (e) {
      console.error('Invalid user data', e)
    }
  } else {
    router.push('/')
  }
})
</script>
