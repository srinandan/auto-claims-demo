<template>
  <footer class="lg:ml-64 bg-[#f9faf5] border-t-0 mt-20">
    <div class="w-full py-12 px-8 flex flex-col md:flex-row justify-between items-center max-w-7xl mx-auto">
      <p class="text-[#47607e] font-body text-xs uppercase tracking-widest mb-4 md:mb-0">
        © 2024 Cymbal Insurance. A Digital Concierge Experience.
      </p>
      <div class="flex flex-wrap justify-center gap-6">
        <a class="text-[#47607e] font-body text-xs uppercase tracking-widest hover:text-[#000000] transition-colors" href="#">Support Center</a>
        <router-link class="text-[#47607e] font-body text-xs uppercase tracking-widest hover:text-[#000000] transition-colors" to="/privacy">Privacy Policy</router-link>
        <router-link class="text-[#47607e] font-body text-xs uppercase tracking-widest hover:text-[#000000] transition-colors" to="/terms">Terms of Service</router-link>
        <a class="text-[#47607e] font-body text-xs uppercase tracking-widest hover:text-[#000000] transition-colors" href="#">Contact Specialist</a>
      </div>
    </div>
  </footer>
</template>

<script setup>
</script>
