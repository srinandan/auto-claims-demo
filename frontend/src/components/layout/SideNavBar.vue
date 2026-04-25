<template>
  <aside class="hidden lg:flex flex-col fixed left-0 top-0 h-full w-64 bg-[#edeeea] p-6 pt-24 space-y-8 z-40">
    <div class="space-y-1">
      <h2 class="font-headline font-extrabold text-[#0D1B2A] text-lg">Cymbal</h2>
      <p class="text-xs text-[#47607e] uppercase tracking-widest font-body">Premium Concierge</p>
    </div>
    <nav class="flex-1 space-y-2">
      <router-link
        to="/dashboard"
        class="flex items-center gap-3 px-4 py-3 rounded-lg font-body font-medium text-sm transition-transform duration-200 hover:translate-x-1"
        :class="[
          $route.path === '/dashboard' || $route.path.startsWith('/claims/')
            ? 'bg-[#ffffff] text-[#000000] shadow-sm'
            : 'text-[#47607e] hover:bg-[#f3f4f0]'
        ]"
      >
        <span class="material-symbols-outlined">dashboard</span>
        <span>Overview</span>
      </router-link>
      <router-link
        to="/new-claim"
        class="flex items-center gap-3 px-4 py-3 rounded-lg font-body font-medium text-sm transition-transform duration-200 hover:translate-x-1 group"
        :class="[
          $route.path === '/new-claim'
            ? 'bg-[#ffffff] text-[#000000] shadow-sm opacity-80'
            : 'text-[#47607e] hover:bg-[#f3f4f0]'
        ]"
      >
        <span class="material-symbols-outlined" :style="$route.path === '/new-claim' ? 'font-variation-settings: \'FILL\' 1;' : ''">add_circle</span>
        <span>New Claim</span>
      </router-link>
      <router-link
        to="/documents"
        class="flex items-center gap-3 px-4 py-3 rounded-lg font-body font-medium text-sm transition-transform duration-200 hover:translate-x-1 group"
        :class="[
          $route.path === '/documents'
            ? 'bg-[#ffffff] text-[#000000] shadow-sm'
            : 'text-[#47607e] hover:bg-[#f3f4f0]'
        ]"
      >
        <span class="material-symbols-outlined group-hover:opacity-80">description</span>
        <span>Documents</span>
      </router-link>
      <router-link
        to="/settings"
        class="flex items-center gap-3 px-4 py-3 rounded-lg font-body font-medium text-sm transition-transform duration-200 hover:translate-x-1 group"
        :class="[
          $route.path === '/settings'
            ? 'bg-[#ffffff] text-[#000000] shadow-sm'
            : 'text-[#47607e] hover:bg-[#f3f4f0]'
        ]"
      >
        <span class="material-symbols-outlined group-hover:opacity-80">settings</span>
        <span>Settings</span>
      </router-link>
    </nav>
    <div class="mt-auto space-y-6">
      <button class="w-full py-4 px-4 bg-primary text-on-primary rounded-xl font-headline font-bold text-sm shadow-xl flex items-center justify-center gap-2 active:scale-95 transition-all">
        <span class="material-symbols-outlined">emergency</span>
        Emergency Assist
      </button>
      <div class="flex flex-col gap-2 border-t border-outline-variant/10 pt-4">
        <router-link class="flex items-center gap-2 text-xs text-[#47607e] font-body hover:text-primary transition-colors" to="/terms">
          <span class="material-symbols-outlined text-sm">gavel</span> Legal
        </router-link>
        <router-link class="flex items-center gap-2 text-xs text-[#47607e] font-body hover:text-primary transition-colors" to="/privacy">
          <span class="material-symbols-outlined text-sm">shield</span> Privacy
        </router-link>
      </div>
    </div>
  </aside>
</template>

<script setup>
</script>
