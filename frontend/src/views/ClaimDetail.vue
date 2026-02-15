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
  <div class="p-6 h-screen flex flex-col box-border">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6 bg-white p-4 rounded shadow shrink-0">
      <div>
        <h1 class="text-2xl font-bold flex items-center">
          Claim #{{ claim.ID }}
          <span :class="['ml-3 px-3 py-1 text-sm rounded-full', statusBadgeClass(claim.status)]">
            {{ claim.status }}
          </span>
        </h1>
        <p class="text-gray-600">{{ claim.customer_name }} - {{ claim.policy_number }}</p>
      </div>
      <div class="space-x-2">
        <button
          v-if="claim.status === 'New'"
          @click="analyzeClaim"
          class="bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"
          :disabled="analyzing"
        >
          {{ analyzing ? 'Analyzing...' : 'Analyze with AI' }}
        </button>
      </div>
    </div>

    <!-- Repair Shops Modal -->
    <div v-if="showRepairShops" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl max-h-[80vh] flex flex-col">
        <div class="p-4 border-b flex justify-between items-center">
          <h2 class="text-xl font-bold">Nearby Repair Shops</h2>
          <button @click="showRepairShops = false" class="text-gray-500 hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-4 overflow-y-auto flex-grow">
          <div v-if="loadingShops" class="text-center py-8">
            <p>Searching for the best shops near you...</p>
          </div>

          <div v-else-if="repairShops.length === 0" class="text-center py-8 text-gray-500">
            No repair shops found.
          </div>

          <div v-else class="space-y-4">
            <div v-for="(shop, index) in repairShops" :key="index" class="border rounded p-4 hover:bg-gray-50">
              <div class="flex justify-between items-start">
                <h3 class="font-bold text-lg">{{ shop.name }}</h3>
                <span v-if="shop.rating" class="bg-yellow-100 text-yellow-800 text-xs px-2 py-1 rounded flex items-center">
                  ★ {{ shop.rating }}
                </span>
              </div>
              <p class="text-gray-600 text-sm mt-1">{{ shop.address }}</p>
              <p v-if="shop.phone" class="text-gray-600 text-sm">{{ shop.phone }}</p>
              <p v-if="shop.reasoning" class="text-sm mt-2 text-gray-700 italic">"{{ shop.reasoning }}"</p>
            </div>
          </div>
        </div>

        <div class="p-4 border-t bg-gray-50 text-right">
          <button @click="showRepairShops = false" class="bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded">
            Close
          </button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="text-center py-10 flex-grow">
      <p>Loading details...</p>
    </div>

    <div v-else class="flex flex-grow gap-6 overflow-hidden h-full">
      <!-- Left: Thumbnails -->
      <div class="w-1/6 bg-white p-4 rounded shadow overflow-y-auto shrink-0">
        <h3 class="font-bold mb-4">Photos</h3>
        <div
          v-for="photo in claim.photos"
          :key="photo.ID"
          class="mb-4 cursor-pointer border-2 transition-colors relative"
          :class="selectedPhoto?.ID === photo.ID ? 'border-blue-500' : 'border-transparent'"
          @click="selectPhoto(photo)"
        >
          <img
            :src="getImageSrc(photo.url)"
            class="w-full h-24 object-cover rounded"
          />
          <div class="text-xs text-center mt-1 text-gray-500">
             <span v-if="photo.analysis_result?.ID" class="text-green-600 font-bold">✓ Analyzed</span>
             <span v-else>Pending</span>
          </div>
        </div>
      </div>

      <!-- Center: Active Image -->
      <div class="flex-grow bg-gray-900 rounded shadow flex items-center justify-center relative overflow-hidden p-4">
        <div v-if="selectedPhoto" class="relative inline-block" ref="imageContainer">
          <img
            :src="getImageSrc(selectedPhoto.url)"
            class="max-w-full max-h-[80vh] object-contain block"
            ref="imageRef"
            @load="onImageLoad"
          />

          <!-- Bounding Boxes Overlay -->
          <div
            v-if="detections.length > 0"
            class="absolute top-0 left-0 w-full h-full pointer-events-none"
          >
             <div
               v-for="(det, idx) in detections"
               :key="idx"
               class="absolute border-2 border-red-500 bg-transparent"
               :style="{
                 left: (det.box[0] * 100) + '%',
                 top: (det.box[1] * 100) + '%',
                 width: ((det.box[2] - det.box[0]) * 100) + '%',
                 height: ((det.box[3] - det.box[1]) * 100) + '%'
               }"
             >
               <span class="absolute -top-6 left-0 bg-red-600 text-white text-xs px-1 rounded whitespace-nowrap z-10">
                 {{ det.label }} ({{ Math.round(det.score * 100) }}%)
               </span>
             </div>
          </div>
        </div>
        <div v-else class="text-white">Select a photo</div>
      </div>

      <!-- Right: Analysis & Estimate -->
      <div class="w-1/4 bg-white p-4 rounded shadow overflow-y-auto shrink-0">
        <!-- Photo Analysis -->
        <div v-if="selectedPhoto && selectedPhoto.analysis_result?.ID" class="mb-6">
          <h3 class="font-bold text-lg mb-2">Image Analysis</h3>

          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="bg-gray-50 p-3 rounded">
              <span class="text-xs text-gray-500 block">Quality</span>
              <span class="font-bold" :class="qualityColor(selectedPhoto.analysis_result.quality_score)">
                {{ selectedPhoto.analysis_result.quality_score }}
              </span>
            </div>
            <div class="bg-gray-50 p-3 rounded">
              <span class="text-xs text-gray-500 block">Severity</span>
              <span class="font-bold">{{ selectedPhoto.analysis_result.severity || 'N/A' }}</span>
            </div>
          </div>

          <div class="mb-4">
            <h4 class="font-semibold text-sm mb-1">Detected Damage</h4>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="part in parseParts(selectedPhoto.analysis_result.parts_detected)"
                :key="part"
                class="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded capitalize"
              >
                {{ part.replace('_', ' ') }}
              </span>
              <span v-if="parseParts(selectedPhoto.analysis_result.parts_detected).length === 0" class="text-sm text-gray-500">None detected</span>
            </div>
          </div>
        </div>
        <div v-else-if="selectedPhoto" class="mb-6 text-gray-500 italic">
          No analysis data available. <br/>Click "Analyze with AI" above.
        </div>

        <hr class="my-4" />

        <!-- Claim Estimate -->
        <div v-if="claim.estimates && claim.estimates.length > 0">
          <h3 class="font-bold text-lg mb-2">Estimate ({{ claim.estimates[0].source }})</h3>

          <div class="bg-gray-50 p-4 rounded border mb-6">
            <div class="flex justify-between font-bold text-lg mb-2">
              <span>Total</span>
              <span>${{ claim.estimates[0].total_amount.toFixed(2) }}</span>
            </div>
             <!-- Simple Items List -->
             <div class="text-xs text-gray-600 space-y-1 mt-2">
               <div v-for="(item, idx) in parseItems(claim.estimates[0].items)" :key="idx" class="flex justify-between">
                 <span>{{ item.part }}</span>
                 <span>${{ item.cost }}</span>
               </div>
             </div>
          </div>

          <!-- Decision Actions -->
          <div v-if="claim.status === 'Assessed'" class="space-y-3">
            <h4 class="font-semibold text-sm">Claim Decision</h4>
            <button
              @click="updateStatus('Approved')"
              class="w-full bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded flex justify-center items-center"
              :disabled="updating"
            >
              <span>Accept Estimate</span>
            </button>
            <button
              @click="updateStatus('Review Required')"
              class="w-full bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-2 px-4 rounded flex justify-center items-center"
              :disabled="updating"
            >
              <span>File Exception (Request Review)</span>
            </button>
          </div>

          <div v-if="claim.status === 'Approved'" class="p-4 bg-green-100 text-green-800 rounded border border-green-200">
            <p class="font-bold text-center">Estimate Accepted</p>
            <p class="text-xs text-center mt-1">Payment processing will begin shortly.</p>
          </div>

          <div v-if="claim.status === 'Review Required'" class="p-4 bg-yellow-100 text-yellow-800 rounded border border-yellow-200">
            <p class="font-bold text-center">Exception Filed</p>
            <p class="text-xs text-center mt-1">An agent will review your claim.</p>
          </div>

        </div>

        <!-- Repair Shop Search -->
        <div v-if="['Assessed', 'Approved', 'Review Required'].includes(claim.status)" class="mt-6 pt-4 border-t">
          <button
            @click="findRepairShops"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded flex justify-center items-center"
          >
            <span class="mr-2">🔧</span> Find Repair Shops
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import axios from 'axios'
import { useRoute } from 'vue-router'

const route = useRoute()
const claim = ref({})
const loading = ref(true)
const analyzing = ref(false)
const updating = ref(false)
const selectedPhoto = ref(null)
const imageRef = ref(null)
const showRepairShops = ref(false)
const repairShops = ref([])
const loadingShops = ref(false)

// Detections for the selected photo
const detections = computed(() => {
  if (!selectedPhoto.value || !selectedPhoto.value.analysis_result?.detections) return []
  try {
    const d = JSON.parse(selectedPhoto.value.analysis_result.detections)
    return Array.isArray(d) ? d : []
  } catch (e) {
    return []
  }
})

const fetchClaim = async () => {
  try {
    const response = await axios.get(`/api/claims/${route.params.id}`)
    claim.value = response.data
    // Select first photo by default if none selected
    if (!selectedPhoto.value && claim.value.photos && claim.value.photos.length > 0) {
      selectedPhoto.value = claim.value.photos[0]
    } else if (selectedPhoto.value) {
        // Refresh selected photo data from new claim data
        const fresh = claim.value.photos.find(p => p.ID === selectedPhoto.value.ID)
        if (fresh) selectedPhoto.value = fresh
    }
  } catch (error) {
    console.error('Error fetching claim:', error)
  } finally {
    loading.value = false
  }
}

const analyzeClaim = async () => {
  analyzing.value = true
  try {
    const response = await axios.post(`/api/claims/${route.params.id}/analyze`)
    claim.value = response.data
    // Update selected photo with new data
    if (selectedPhoto.value) {
       const fresh = claim.value.photos.find(p => p.ID === selectedPhoto.value.ID)
       if (fresh) selectedPhoto.value = fresh
    }
  } catch (error) {
    console.error('Error analyzing claim:', error)
    alert('Analysis failed')
  } finally {
    analyzing.value = false
  }
}

const findRepairShops = async () => {
  showRepairShops.value = true
  loadingShops.value = true
  repairShops.value = []

  try {
    const response = await axios.post(`/api/claims/${route.params.id}/repair-shops`)
    if (response.data && response.data.shops) {
      repairShops.value = response.data.shops
    }
  } catch (error) {
    console.error('Error finding repair shops:', error)
    alert('Failed to find repair shops')
    showRepairShops.value = false
  } finally {
    loadingShops.value = false
  }
}

const updateStatus = async (status) => {
  updating.value = true
  try {
    const response = await axios.put(`/api/claims/${route.params.id}`, { status })
    claim.value.status = response.data.status
  } catch (error) {
    console.error('Error updating claim:', error)
    alert('Update failed')
  } finally {
    updating.value = false
  }
}

const selectPhoto = (photo) => {
  selectedPhoto.value = photo
}

const getImageSrc = (url) => {
  if (url.startsWith('http')) return url
  if (url.startsWith('uploads/')) return `/${url}`
  // Using placehold.co for dummy images
  return `https://placehold.co/800x600/EEE/31343C?text=${url}`
}

const parseParts = (partsStr) => {
  if (!partsStr) return []
  try {
    // Check if it is a JSON array string
    if (partsStr.startsWith('[')) {
        return JSON.parse(partsStr)
    }
    // Else comma separated
    return partsStr.split(',').filter(p => p)
  } catch (e) {
    return []
  }
}

const parseItems = (itemsStr) => {
    if (!itemsStr) return []
    try {
        return JSON.parse(itemsStr)
    } catch (e) {
        return []
    }
}

const statusBadgeClass = (status) => {
  switch (status) {
    case 'Simple': return 'bg-green-100 text-green-800'
    case 'Complex': return 'bg-yellow-100 text-yellow-800'
    case 'Total Loss': return 'bg-red-100 text-red-800'
    case 'Approved': return 'bg-green-100 text-green-800'
    case 'Review Required': return 'bg-yellow-100 text-yellow-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const qualityColor = (score) => {
    if (score === 'Good') return 'text-green-600'
    if (score === 'Blurry' || score === 'Dark') return 'text-red-600'
    return 'text-gray-600'
}

const onImageLoad = () => {
    // Trigger reactivity if needed, but CSS should handle the layout
}

onMounted(() => {
  fetchClaim()
})
</script>
