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
          v-if="claim.status === 'New' || claim.status === 'Simple' || claim.status === 'Complex'"
          @click="analyzeClaim"
          class="bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50"
          :disabled="analyzing"
        >
          {{ analyzing ? 'Analyzing...' : 'Analyze with AI' }}
        </button>
        <button class="bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded">
          Request Retake
        </button>
        <button class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
          Approve Estimate
        </button>
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
          <!-- We position this absolutely over the image container.
               Since the container shrinks to fit the image (inline-block),
               absolute positioning 0,0,100%,100% should match the image EXACTLY.
          -->
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
              <span v-if="parseParts(selectedPhoto.analysis_result.parts_detected).length === 0" class="text-sm text-gray-500">None</span>
            </div>
          </div>
        </div>
        <div v-else-if="selectedPhoto" class="mb-6 text-gray-500 italic">
          No analysis data available. <br/>Click "Analyze with AI" above.
        </div>

        <hr class="my-4" />

        <!-- Claim Estimate -->
        <div>
          <h3 class="font-bold text-lg mb-2">Assessment</h3>

          <div class="mb-4">
            <span class="text-sm text-gray-500">Total Loss Probability</span>
            <div class="w-full bg-gray-200 rounded-full h-2.5 mt-1">
              <div
                class="bg-red-600 h-2.5 rounded-full"
                :style="{ width: (claim.total_loss_probability * 100) + '%' }"
              ></div>
            </div>
            <div class="text-right text-xs mt-1">{{ (claim.total_loss_probability * 100).toFixed(1) }}%</div>
          </div>

          <div v-if="claim.estimates && claim.estimates.length > 0">
            <h4 class="font-semibold text-sm mb-2">Draft Estimate ({{ claim.estimates[0].source }})</h4>
            <div class="bg-gray-50 p-4 rounded border">
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

            <!-- Repair Shop Comparison Placeholder -->
             <div class="mt-4 p-3 border border-dashed border-gray-400 rounded bg-yellow-50">
               <h5 class="font-bold text-sm text-yellow-800">Repair Shop Comparison</h5>
               <p class="text-xs text-gray-600 mt-1">Submit external estimate to compare. (Optional)</p>
             </div>
          </div>
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
const selectedPhoto = ref(null)
const imageRef = ref(null)

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
