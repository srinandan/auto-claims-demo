<template>
  <div class="pt-8 pb-20 px-8 max-w-5xl mx-auto selection:bg-tertiary-fixed">
    <!-- Hero Section -->
    <div class="mb-20">
      <div class="inline-flex items-center px-3 py-1 bg-tertiary-fixed text-on-tertiary-fixed-variant rounded-full text-xs font-bold tracking-widest uppercase mb-6">
        Security & Trust
      </div>
      <h1 class="text-6xl font-extrabold text-on-surface leading-tight tracking-tighter mb-6">Privacy Policy</h1>
      <div class="flex items-center gap-4 text-on-surface-variant">
        <span class="w-12 h-[1px] bg-outline-variant"></span>
        <p class="font-label text-sm uppercase tracking-widest font-medium">Last Updated: October 24, 2024</p>
      </div>
    </div>

    <!-- Asymmetric Layout Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-16">
      <!-- Sticky Sidebar Navigation (Anchor) -->
      <aside class="lg:col-span-3 hidden lg:block">
        <div class="sticky top-32 space-y-4">
          <p class="text-xs font-bold uppercase tracking-[0.2em] text-secondary mb-8">On This Page</p>
          <nav class="flex flex-col space-y-3">
            <a class="text-on-surface hover:text-primary transition-colors font-medium border-l-2 border-primary pl-4 py-1" href="#introduction">Introduction</a>
            <a class="text-on-surface-variant hover:text-primary transition-colors pl-4 py-1" href="#data-collection">Data Collection</a>
            <a class="text-on-surface-variant hover:text-primary transition-colors pl-4 py-1" href="#usage">How We Use Data</a>
            <a class="text-on-surface-variant hover:text-primary transition-colors pl-4 py-1" href="#sharing">Information Sharing</a>
            <a class="text-on-surface-variant hover:text-primary transition-colors pl-4 py-1" href="#security">Data Security</a>
            <a class="text-on-surface-variant hover:text-primary transition-colors pl-4 py-1" href="#rights">Your Privacy Rights</a>
          </nav>

          <div class="mt-12 p-6 bg-surface-container rounded-xl">
            <span class="material-symbols-outlined text-primary mb-4" data-icon="lock" style="font-variation-settings: 'FILL' 1;">lock</span>
            <h4 class="font-headline font-bold text-sm mb-2">Concierge Support</h4>
            <p class="text-xs text-on-surface-variant leading-relaxed">Need help understanding our privacy practices?</p>
            <button class="mt-4 text-xs font-bold text-primary flex items-center gap-1 group">
              Contact a Specialist
              <span class="material-symbols-outlined text-sm group-hover:translate-x-1 transition-transform" data-icon="arrow_forward">arrow_forward</span>
            </button>
          </div>
        </div>
      </aside>

      <!-- Content Area -->
      <article class="lg:col-span-9 space-y-16">
        <section id="introduction" class="scroll-mt-32">
          <h2 class="text-3xl font-bold mb-6 text-on-surface">1. Introduction</h2>
          <div class="space-y-4 text-lg leading-relaxed text-on-surface-variant font-light">
            <p>At Cymbal Insurance, we believe that privacy is the ultimate luxury. Our Digital Concierge service is built on a foundation of absolute transparency and unyielding security. This policy outlines how we handle your personal information with the same care and precision we bring to our bespoke insurance solutions.</p>
            <p>By using the Cymbal Sovereign platform, you entrust us with sensitive details. We honor that trust by implementing industry-leading data governance protocols designed to protect your identity and your legacy.</p>
          </div>
        </section>

        <div class="w-full h-80 rounded-2xl overflow-hidden shadow-2xl">
          <img alt="Cybersecurity representation" class="w-full h-full object-cover" data-alt="high-tech secure server room with glowing blue led lights and data patterns, professional architectural photography, cinematic lighting" src="https://lh3.googleusercontent.com/aida-public/AB6AXuBh0bMNKZ34S39cNs2tv1oT1W24kIimfsHSWbjAuRhYWpc2r8FwlyOsd8m4nBwRVu-cj9sI_T4tmE9epq8YaqNV_isvbgfOBcJ0ZhpN9N2pFfsySXiZD-r-zM6TacQSH0tEk-b0lpY1JoTxFgQO5h7kvb-sYBulLFR92zV6SyiL0Cg3iqRIoMJdmSKYgYm7ppUOqW_axcm_Xzw8AKNO-TibrLvTFzVw2Kv1v2QJEmHR769vrm9pmtcHWIclAhfXjzo6xyUIe4fq3N8"/>
        </div>

        <section id="data-collection" class="scroll-mt-32">
          <h2 class="text-3xl font-bold mb-8 text-on-surface">2. Information We Collect</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="bg-surface-container-lowest p-8 rounded-xl shadow-sm border border-outline-variant/10">
              <span class="material-symbols-outlined text-secondary mb-4" data-icon="person">person</span>
              <h3 class="font-bold text-xl mb-3">Direct Identifiers</h3>
              <p class="text-on-surface-variant leading-relaxed">Name, address, contact information, and government-issued identification required for policy underwriting and verification.</p>
            </div>
            <div class="bg-surface-container-lowest p-8 rounded-xl shadow-sm border border-outline-variant/10">
              <span class="material-symbols-outlined text-secondary mb-4" data-icon="analytics">analytics</span>
              <h3 class="font-bold text-xl mb-3">Behavioral Data</h3>
              <p class="text-on-surface-variant leading-relaxed">Interaction logs with the Cymbal Sovereign platform, including claim filing frequency and feature usage to improve concierge services.</p>
            </div>
          </div>
        </section>

        <section id="usage" class="scroll-mt-32">
          <div class="bg-on-primary-fixed text-white p-12 rounded-3xl relative overflow-hidden">
            <div class="relative z-10">
              <h2 class="text-3xl font-bold mb-6">3. Data Usage Purpose</h2>
              <p class="text-on-primary-container text-lg mb-8 max-w-2xl leading-relaxed">Your data is never "productized." We utilize information exclusively to enhance your insurance experience and provide proactive risk management.</p>
              <ul class="space-y-4">
                <li class="flex items-start gap-3">
                  <span class="material-symbols-outlined text-tertiary-fixed mt-1" data-icon="check_circle" style="font-variation-settings: 'FILL' 1;">check_circle</span>
                  <span class="text-on-primary-container">Personalized policy refinement based on real-time asset valuation.</span>
                </li>
                <li class="flex items-start gap-3">
                  <span class="material-symbols-outlined text-tertiary-fixed mt-1" data-icon="check_circle" style="font-variation-settings: 'FILL' 1;">check_circle</span>
                  <span class="text-on-primary-container">Expedited claim processing through automated AI analysis of submitted documentation.</span>
                </li>
                <li class="flex items-start gap-3">
                  <span class="material-symbols-outlined text-tertiary-fixed mt-1" data-icon="check_circle" style="font-variation-settings: 'FILL' 1;">check_circle</span>
                  <span class="text-on-primary-container">Fraud prevention and anti-money laundering compliance checks.</span>
                </li>
              </ul>
            </div>
            <div class="absolute right-0 top-0 w-1/3 h-full opacity-10 bg-gradient-to-l from-white to-transparent pointer-events-none"></div>
          </div>
        </section>

        <section id="security" class="scroll-mt-32">
          <h2 class="text-3xl font-bold mb-6 text-on-surface">4. Security Infrastructure</h2>
          <div class="prose prose-slate max-w-none text-on-surface-variant text-lg leading-relaxed space-y-6">
            <p>We employ <strong>AES-256 bit encryption</strong> for all data at rest and TLS 1.3 for all data in transit. Our security operations center (SOC) monitors the Sovereign infrastructure 24/7/365 to identify and mitigate emerging threats before they can impact our clients.</p>
            <p>Furthermore, we implement strict physical security at our data centers, ensuring that your digital footprint is protected by multiple layers of physical and logical barriers.</p>
          </div>
        </section>

        <section id="rights" class="scroll-mt-32">
          <h2 class="text-3xl font-bold mb-8 text-on-surface">5. Your Global Privacy Rights</h2>
          <div class="bg-surface-container-low p-1 rounded-2xl">
            <div class="bg-surface-container-lowest p-8 rounded-xl">
              <div class="flex flex-col md:flex-row gap-8 items-center">
                <div class="flex-1">
                  <p class="text-on-surface-variant leading-relaxed">Regardless of your jurisdiction, Cymbal provides the right to access, rectify, or erase your personal data. You may also object to processing or request data portability at any time through your Concierge dashboard.</p>
                </div>
                <button class="px-8 py-4 bg-primary text-white rounded-lg font-bold tracking-tight hover:opacity-90 transition-opacity whitespace-nowrap">
                  Manage Privacy Settings
                </button>
              </div>
            </div>
          </div>
        </section>
      </article>
    </div>
  </div>
</template>

<script setup>
</script>