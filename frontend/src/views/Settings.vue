<template>
  <div class="max-w-5xl mx-auto px-8 py-12">
    <!-- Header -->
    <header class="mb-12">
      <h1 class="font-headline font-extrabold text-2xl text-[#0D1B2A] tracking-tight">Account Settings</h1>
    </header>

    <div class="space-y-16">
      <!-- Section 1: Profile -->
      <section>
        <div class="mb-8">
            <h2 class="text-3xl font-headline font-extrabold text-on-surface tracking-tight mb-2">Personal Information</h2>
            <p class="text-on-surface-variant font-body max-w-2xl">Manage your concierge profile details and how we contact you for policy updates and claims.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 items-start">
            <!-- Profile Card -->
            <div class="md:col-span-1 bg-surface-container-lowest rounded-xl p-8 shadow-[0_12px_32px_rgba(16,27,48,0.04)]">
              <div class="relative w-32 h-32 mx-auto mb-6 group">
                  <img
                    alt="Profile Large"
                    class="w-full h-full object-cover rounded-full"
                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuAzZ8mQYgcNnugTdatsLWQUnSy5A8pVscnWrxN5ZYaffAtnzGvJDJqryjag8DworCC2g2y7DP8sql2LQub4adxP3F_HUdHZHok1NUFSlHiqEEO4lUDWJQAzXhPaiOwoKU5QX689MANJwRhp-6DRa-zFZkoCwYHwgSM2c4IgqQfSNmH_i9cJIzRMoYKK16UXy0fjo1dKDeaHrSWpUJpLyA_i43p5r7cXMlXhgsNIb3JHzDaZYPHRl6pEcs0Q0v6hGxyiNgcAuXlV9Ng"
                  />
                  <button class="absolute bottom-0 right-0 bg-primary text-white p-2 rounded-full shadow-lg scale-90 group-hover:scale-100 transition-transform">
                    <span class="material-symbols-outlined text-sm">edit</span>
                  </button>
              </div>
              <div class="text-center">
                  <h3 class="font-headline font-bold text-lg text-on-surface">{{ user?.first_name }} {{ user?.last_name }}</h3>
                  <p class="text-xs font-label uppercase tracking-widest text-secondary mt-1">Policy Holder</p>
              </div>
            </div>

            <!-- Form Area -->
            <div class="md:col-span-2 space-y-6">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                  <div class="space-y-2">
                    <label class="text-xs font-label uppercase tracking-widest text-secondary ml-1">First Name</label>
                    <input
                        type="text"
                        :value="user?.first_name"
                        class="w-full bg-surface-container-lowest border border-outline-variant/15 rounded-lg px-4 py-3 focus:border-primary focus:ring-0 transition-all font-body text-on-surface"
                    />
                  </div>
                  <div class="space-y-2">
                    <label class="text-xs font-label uppercase tracking-widest text-secondary ml-1">Last Name</label>
                    <input
                        type="text"
                        :value="user?.last_name"
                        class="w-full bg-surface-container-lowest border border-outline-variant/15 rounded-lg px-4 py-3 focus:border-primary focus:ring-0 transition-all font-body text-on-surface"
                    />
                  </div>
              </div>

              <div class="space-y-2">
                  <label class="text-xs font-label uppercase tracking-widest text-secondary ml-1">Email Address</label>
                  <input
                    type="email"
                    :value="user?.email"
                    class="w-full bg-surface-container-lowest border border-outline-variant/15 rounded-lg px-4 py-3 focus:border-primary focus:ring-0 transition-all font-body text-on-surface"
                  />
              </div>

              <div class="flex justify-end pt-4">
                  <button class="bg-primary text-on-primary px-8 py-3 rounded-md font-headline font-bold text-sm tracking-tight hover:shadow-lg transition-all active:scale-95">
                      Update Profile
                  </button>
              </div>
            </div>
        </div>
      </section>

      <!-- Section 2: Notifications & Security (Asymmetric Layout) -->
      <section class="grid grid-cols-1 lg:grid-cols-2 gap-12">

        <!-- Notifications -->
        <div class="space-y-8">
            <div>
              <h2 class="text-2xl font-headline font-bold text-on-surface tracking-tight mb-2">Notification Preferences</h2>
              <p class="text-sm text-on-surface-variant font-body">Stay updated on your claim progress and renewals.</p>
            </div>
            <div class="bg-surface-container-low rounded-xl p-2 space-y-1">
              <div class="flex items-center justify-between p-4 bg-surface-container-lowest rounded-lg">
                  <div class="flex items-center space-x-4">
                    <span class="material-symbols-outlined text-secondary">mail</span>
                    <div>
                        <div class="font-medium text-on-surface font-body">Email Updates</div>
                        <div class="text-xs text-on-surface-variant font-body">Monthly reports and policy changes</div>
                    </div>
                  </div>
                  <div class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" checked class="sr-only peer">
                    <div class="w-11 h-6 bg-surface-container-highest peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                  </div>
              </div>

              <div class="flex items-center justify-between p-4 bg-surface-container-lowest rounded-lg">
                  <div class="flex items-center space-x-4">
                    <span class="material-symbols-outlined text-secondary">sms</span>
                    <div>
                        <div class="font-medium text-on-surface font-body">SMS Alerts</div>
                        <div class="text-xs text-on-surface-variant font-body">Instant claim status changes</div>
                    </div>
                  </div>
                  <div class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" class="sr-only peer">
                    <div class="w-11 h-6 bg-surface-container-highest peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                  </div>
              </div>
            </div>
        </div>

        <!-- Security -->
        <div class="space-y-8">
            <div>
              <h2 class="text-2xl font-headline font-bold text-on-surface tracking-tight mb-2">Security & Access</h2>
              <p class="text-sm text-on-surface-variant font-body">Protect your premium assets and information.</p>
            </div>

            <div class="space-y-4">
              <div class="group cursor-pointer flex items-center justify-between p-5 bg-surface-container-lowest rounded-xl border border-outline-variant/10 hover:border-primary/30 transition-all shadow-[0_4px_12px_rgba(16,27,48,0.02)]">
                  <div class="flex items-center space-x-4">
                    <div class="w-10 h-10 bg-tertiary-fixed flex items-center justify-center rounded-full">
                        <span class="material-symbols-outlined text-on-tertiary-fixed">lock_reset</span>
                    </div>
                    <div>
                        <div class="font-bold text-on-surface font-body">Change Password</div>
                        <div class="text-xs text-on-surface-variant font-body">Last updated 4 months ago</div>
                    </div>
                  </div>
                  <span class="material-symbols-outlined text-secondary opacity-30 group-hover:opacity-100 transition-opacity">chevron_right</span>
              </div>

              <div class="group cursor-pointer flex items-center justify-between p-5 bg-surface-container-lowest rounded-xl border border-outline-variant/10 hover:border-primary/30 transition-all shadow-[0_4px_12px_rgba(16,27,48,0.02)]">
                  <div class="flex items-center space-x-4">
                    <div class="w-10 h-10 bg-secondary-fixed flex items-center justify-center rounded-full">
                        <span class="material-symbols-outlined text-on-secondary-fixed">key</span>
                    </div>
                    <div>
                        <div class="font-bold text-on-surface font-body">Biometric Login</div>
                        <div class="text-xs text-on-surface-variant font-body">FaceID or Fingerprint enabled</div>
                    </div>
                  </div>
                  <span class="material-symbols-outlined text-secondary opacity-30 group-hover:opacity-100 transition-opacity">chevron_right</span>
              </div>
            </div>
        </div>
      </section>

      <!-- Section 3: Advanced (Magazine Style) -->
      <section class="bg-primary-container text-on-primary-container rounded-[2rem] p-12 overflow-hidden relative">
        <div class="relative z-10 max-w-xl">
            <span class="inline-block px-3 py-1 bg-tertiary-fixed text-on-tertiary-fixed text-[10px] font-bold uppercase tracking-widest rounded-full mb-6 font-label">Concierge Benefit</span>
            <h2 class="text-4xl font-headline font-extrabold tracking-tight mb-6 leading-tight">Digital Vault Management</h2>
            <p class="text-on-primary-container/80 font-body mb-8 leading-relaxed">Securely store title deeds, estate documents, and high-value asset certificates with 256-bit encryption. Accessible only through multi-signature authentication.</p>

            <button class="bg-surface-container-lowest text-primary px-8 py-4 rounded-md font-headline font-bold text-sm tracking-tight hover:scale-105 transition-transform flex items-center space-x-3">
              <span>Open Digital Vault</span>
              <span class="material-symbols-outlined">arrow_forward</span>
            </button>
        </div>

        <!-- Decorative Background Image -->
        <div class="absolute right-[-10%] top-[-10%] w-2/3 h-120 opacity-20 transform rotate-12 pointer-events-none">
            <img
              alt="Technology background"
              class="w-full h-full object-cover rounded-[2rem]"
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuBzA6wJCl3vx2wAQeusKhAAyxn4MLpYxpGpn_PQm6QRrZQLHLajVHdHJVcoBhtGvrkZC5dtv20tVtSgYmghpnrpdTvS63bSBCx_NVjhouBJ_OesyyYU4AZ_xJ-XayloN6aMsWvc6o6N6o4qSbxNS34Hm81enq0Zy0uQV9D7EyonLAi7Znyf7ToYIuje7d9wXfg5sBY8V_Gm7PGSYxvYYivxo7oFDI--9nlq9o5-zB-5bqoEjx3JMCVSTXHoy0K4uVnwVzooPwVehyM"
            />
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref(null)

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
    } catch (e) {
      console.error('Invalid user data', e)
    }
  } else {
    router.push('/')
  }
})
</script>
