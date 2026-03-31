<template>
  <div class="flex-1 px-8 py-12 max-w-7xl mx-auto w-full">
    <!-- Page Header -->
    <section class="mb-16">
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-6">
          <div class="max-w-2xl">
            <div class="mb-4">
                <span class="bg-tertiary-fixed text-on-tertiary-fixed-variant px-3 py-1 rounded-full text-xs font-bold tracking-wider uppercase font-label">Document Center</span>
            </div>
            <h1 class="text-5xl font-headline font-extrabold text-on-surface tracking-tight mb-4">Secure Repository</h1>
            <p class="text-on-surface-variant text-lg leading-relaxed font-body">Access your insurance contracts, policy amendments, and claim reports in one centralized, high-security space.</p>
          </div>
          <div class="flex items-center gap-3">
            <button class="bg-surface-container-highest text-on-surface px-6 py-3 rounded-lg font-semibold text-sm hover:translate-y-[-2px] transition-transform duration-300 font-body">
                Download All
            </button>
            <button class="bg-primary text-white px-6 py-3 rounded-lg font-semibold text-sm flex items-center gap-2 hover:shadow-xl transition-all duration-300 font-body">
                <span class="material-symbols-outlined text-sm">upload</span>
                Submit New File
            </button>
          </div>
      </div>
    </section>

    <!-- Search and Filter Bar -->
    <section class="mb-10 p-1 bg-surface-container-low rounded-xl">
      <div class="flex flex-col md:flex-row gap-4 p-4 bg-surface-container-lowest rounded-lg shadow-sm">
          <div class="flex-1 relative">
            <span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-outline">search</span>
            <input
                type="text"
                class="w-full pl-12 pr-4 py-3 bg-surface-container-low border-none rounded-lg text-on-surface focus:ring-2 focus:ring-primary transition-all placeholder:text-outline/60 font-body outline-none"
                placeholder="Search by document name, policy ID or claim number..."
            />
          </div>
          <div class="flex gap-2">
            <button class="flex items-center gap-2 px-5 py-3 bg-surface-container-low text-secondary font-medium rounded-lg hover:bg-surface-container-high transition-colors font-body">
                <span class="material-symbols-outlined text-lg">filter_list</span>
                Filter
            </button>
            <button class="flex items-center gap-2 px-5 py-3 bg-surface-container-low text-secondary font-medium rounded-lg hover:bg-surface-container-high transition-colors font-body">
                <span class="material-symbols-outlined text-lg">sort</span>
                Sort
            </button>
          </div>
      </div>
    </section>

    <!-- Document Bento Grid/List -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-8">

      <!-- Category Sidebar -->
      <div class="md:col-span-3 space-y-6">
          <h3 class="font-headline font-bold text-sm text-secondary uppercase tracking-widest px-2">Folders</h3>
          <nav class="space-y-1">
            <button class="w-full flex items-center justify-between px-4 py-3 rounded-lg bg-surface-container-high text-on-surface font-semibold font-body">
                <span class="flex items-center gap-3"><span class="material-symbols-outlined text-primary">folder</span>All Files</span>
                <span class="text-xs opacity-60">24</span>
            </button>
            <button class="w-full flex items-center justify-between px-4 py-3 rounded-lg text-on-surface-variant hover:bg-surface-container-low transition-colors font-body">
                <span class="flex items-center gap-3"><span class="material-symbols-outlined">policy</span>Policies</span>
                <span class="text-xs opacity-60">4</span>
            </button>
            <button class="w-full flex items-center justify-between px-4 py-3 rounded-lg text-on-surface-variant hover:bg-surface-container-low transition-colors font-body">
                <span class="flex items-center gap-3"><span class="material-symbols-outlined">assignment</span>Claims</span>
                <span class="text-xs opacity-60">12</span>
            </button>
            <button class="w-full flex items-center justify-between px-4 py-3 rounded-lg text-on-surface-variant hover:bg-surface-container-low transition-colors font-body">
                <span class="flex items-center gap-3"><span class="material-symbols-outlined">receipt_long</span>Billing</span>
                <span class="text-xs opacity-60">8</span>
            </button>
          </nav>
      </div>

      <!-- Document List -->
      <div class="md:col-span-9 space-y-4">
          <!-- Document Row 1 -->
          <div class="group flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 bg-surface-container-lowest rounded-xl hover:shadow-md transition-all duration-300 gap-4">
            <div class="flex items-center gap-6">
                <div class="w-14 h-14 shrink-0 bg-secondary-fixed flex items-center justify-center rounded-lg text-on-secondary-container">
                  <span class="material-symbols-outlined text-3xl">picture_as_pdf</span>
                </div>
                <div>
                  <h4 class="font-headline font-bold text-lg text-on-surface">Auto Policy ID #{{ user?.policy_number || 'CYM-2024-99' }}</h4>
                  <div class="flex items-center gap-4 mt-1 font-body">
                      <span class="text-sm text-on-surface-variant">Modified Jan 12, 2024</span>
                      <span class="w-1 h-1 rounded-full bg-outline-variant"></span>
                      <span class="text-sm text-on-surface-variant">2.4 MB</span>
                  </div>
                </div>
            </div>
            <div class="flex items-center gap-2 self-end sm:self-auto">
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">visibility</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">download</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">more_vert</span>
                </button>
            </div>
          </div>

          <!-- Document Row 2 -->
          <div class="group flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 bg-surface-container-lowest rounded-xl hover:shadow-md transition-all duration-300 gap-4">
            <div class="flex items-center gap-6">
                <div class="w-14 h-14 shrink-0 bg-tertiary-fixed flex items-center justify-center rounded-lg text-on-tertiary-fixed-variant">
                  <span class="material-symbols-outlined text-3xl">description</span>
                </div>
                <div>
                  <div class="flex items-center gap-2">
                      <h4 class="font-headline font-bold text-lg text-on-surface">Claim Incident Report</h4>
                      <span class="text-[10px] bg-tertiary text-white px-2 py-0.5 rounded font-bold uppercase font-label">New</span>
                  </div>
                  <div class="flex items-center gap-4 mt-1 font-body">
                      <span class="text-sm text-on-surface-variant">Modified Today, 10:24 AM</span>
                      <span class="w-1 h-1 rounded-full bg-outline-variant"></span>
                      <span class="text-sm text-on-surface-variant">842 KB</span>
                  </div>
                </div>
            </div>
            <div class="flex items-center gap-2 self-end sm:self-auto">
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">visibility</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">download</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">more_vert</span>
                </button>
            </div>
          </div>

          <!-- Document Row 3 -->
          <div class="group flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 bg-surface-container-lowest rounded-xl hover:shadow-md transition-all duration-300 gap-4">
            <div class="flex items-center gap-6">
                <div class="w-14 h-14 shrink-0 bg-surface-container-highest flex items-center justify-center rounded-lg text-on-surface-variant">
                  <span class="material-symbols-outlined text-3xl">image</span>
                </div>
                <div>
                  <h4 class="font-headline font-bold text-lg text-on-surface">Vehicle Damage Photos (4)</h4>
                  <div class="flex items-center gap-4 mt-1 font-body">
                      <span class="text-sm text-on-surface-variant">Modified Dec 28, 2023</span>
                      <span class="w-1 h-1 rounded-full bg-outline-variant"></span>
                      <span class="text-sm text-on-surface-variant">12.1 MB</span>
                  </div>
                </div>
            </div>
            <div class="flex items-center gap-2 self-end sm:self-auto">
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">visibility</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">download</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">more_vert</span>
                </button>
            </div>
          </div>

          <!-- Document Row 4 -->
          <div class="group flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 bg-surface-container-lowest rounded-xl hover:shadow-md transition-all duration-300 opacity-75 gap-4">
            <div class="flex items-center gap-6">
                <div class="w-14 h-14 shrink-0 bg-surface-container-highest flex items-center justify-center rounded-lg text-on-surface-variant">
                  <span class="material-symbols-outlined text-3xl">article</span>
                </div>
                <div>
                  <h4 class="font-headline font-bold text-lg text-on-surface">Expired Policy - FY2023</h4>
                  <div class="flex items-center gap-4 mt-1 font-body">
                      <span class="text-sm text-on-surface-variant">Modified Jan 01, 2023</span>
                      <span class="w-1 h-1 rounded-full bg-outline-variant"></span>
                      <span class="text-sm text-on-surface-variant">1.9 MB</span>
                  </div>
                </div>
            </div>
            <div class="flex items-center gap-2 self-end sm:self-auto">
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">visibility</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">download</span>
                </button>
                <button class="p-3 text-secondary hover:bg-surface-container-low rounded-full transition-colors active:scale-90">
                  <span class="material-symbols-outlined">more_vert</span>
                </button>
            </div>
          </div>
      </div>
    </div>

    <!-- Storage Analytics / Promotion Card -->
    <section class="mt-20">
      <div class="relative overflow-hidden rounded-2xl bg-primary p-12 text-white">
          <div class="relative z-10 grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
            <div>
                <h2 class="text-3xl font-headline font-extrabold mb-4">Optimized Storage</h2>
                <p class="text-on-primary-container text-lg mb-8 leading-relaxed font-body">Your account currently uses 124MB of 1GB secure cloud storage. Need more space for high-resolution claim assets?</p>
                <button class="bg-white text-primary px-8 py-3 rounded-lg font-bold hover:bg-surface-container-high transition-colors font-body">
                    Upgrade Storage
                </button>
            </div>
            <div class="flex justify-center md:justify-end">
                <div class="w-48 h-48 rounded-full border-8 border-primary-container flex items-center justify-center">
                  <div class="text-center">
                      <div class="text-4xl font-extrabold font-headline">12%</div>
                      <div class="text-xs font-label uppercase tracking-widest text-on-primary-container">Used</div>
                  </div>
                </div>
            </div>
          </div>
          <!-- Decorative Graphic Elements -->
          <div class="absolute -bottom-20 -right-20 w-80 h-80 bg-primary-container/30 rounded-full blur-3xl"></div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const user = ref(null)

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
    } catch (e) {
      console.error('Invalid user data', e)
    }
  }
})
</script>
