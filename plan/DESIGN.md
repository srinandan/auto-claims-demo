<!DOCTYPE html>

<html class="light" lang="en"><head>
<meta charset="utf-8"/>
<meta content="width=device-width, initial-scale=1.0" name="viewport"/>
<title>New Claim | Cymbal Insurance</title>
<script src="https://cdn.tailwindcss.com?plugins=forms,container-queries"></script>
<link href="https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;700;800&amp;family=Inter:wght@400;500;600&amp;display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&amp;display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&amp;display=swap" rel="stylesheet"/>
<script id="tailwind-config">
        tailwind.config = {
            darkMode: "class",
            theme: {
                extend: {
                    colors: {
                        "secondary-fixed-dim": "#afc9ea",
                        "surface-container-lowest": "#ffffff",
                        "outline": "#74777d",
                        "tertiary": "#000000",
                        "outline-variant": "#c4c6cc",
                        "on-secondary-container": "#48617e",
                        "on-tertiary": "#ffffff",
                        "surface-dim": "#d9dad6",
                        "on-surface-variant": "#44474c",
                        "on-primary-fixed": "#101b30",
                        "primary-container": "#101b30",
                        "primary-fixed-dim": "#bbc6e2",
                        "on-tertiary-fixed-variant": "#0e5138",
                        "surface-bright": "#f9faf5",
                        "on-error": "#ffffff",
                        "inverse-surface": "#2f312e",
                        "surface-container-low": "#f3f4f0",
                        "on-background": "#1a1c1a",
                        "secondary-fixed": "#d1e4ff",
                        "inverse-on-surface": "#f0f1ed",
                        "primary-fixed": "#d7e2ff",
                        "background": "#f9faf5",
                        "surface-container": "#edeeea",
                        "tertiary-fixed-dim": "#95d4b3",
                        "tertiary-fixed": "#b1f0ce",
                        "on-tertiary-fixed": "#002114",
                        "secondary": "#47607e",
                        "on-secondary-fixed": "#001d36",
                        "inverse-primary": "#bbc6e2",
                        "surface-container-highest": "#e2e3df",
                        "surface-variant": "#e2e3df",
                        "primary": "#000000",
                        "on-primary-fixed-variant": "#3c475d",
                        "error": "#ba1a1a",
                        "surface": "#f9faf5",
                        "on-tertiary-container": "#539072",
                        "on-error-container": "#93000a",
                        "on-surface": "#1a1c1a",
                        "error-container": "#ffdad6",
                        "surface-tint": "#545e76",
                        "tertiary-container": "#002114",
                        "on-primary-container": "#79849d",
                        "on-primary": "#ffffff",
                        "secondary-container": "#c2dcff",
                        "on-secondary": "#ffffff",
                        "on-secondary-fixed-variant": "#2f4865",
                        "surface-container-high": "#e8e8e4"
                    },
                    fontFamily: {
                        "headline": ["Manrope"],
                        "body": ["Inter"],
                        "label": ["Inter"]
                    },
                    borderRadius: { "DEFAULT": "0.125rem", "lg": "0.25rem", "xl": "0.5rem", "full": "0.75rem" },
                },
            },
        }
    </script>
<style>
        .material-symbols-outlined {
            font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
        }
        body { font-family: 'Inter', sans-serif; }
        h1, h2, h3 { font-family: 'Manrope', sans-serif; }
    </style>
</head>
<body class="bg-surface text-on-surface">
<!-- SideNavBar (Authority: Shared Components JSON) -->
<nav class="hidden lg:flex flex-col fixed left-0 top-0 h-full p-6 space-y-8 bg-[#edeeea] text-[#0D1B2A] h-screen w-64 border-r-0 font-['Inter'] font-medium text-sm">
<div class="flex flex-col space-y-1 mb-8">
<span class="font-['Manrope'] font-extrabold text-[#0D1B2A] text-2xl">Cymbal</span>
<span class="text-secondary text-[10px] uppercase tracking-widest">Premium Concierge</span>
</div>
<div class="flex-1 space-y-2">
<a class="flex items-center space-x-3 px-4 py-3 text-[#47607e] hover:bg-[#f3f4f0] transition-transform duration-200 hover:translate-x-1 group" href="#">
<span class="material-symbols-outlined group-hover:opacity-80">dashboard</span>
<span>Overview</span>
</a>
<!-- Active State for New Claim -->
<a class="flex items-center space-x-3 px-4 py-3 bg-[#ffffff] text-[#000000] rounded-lg shadow-sm transition-opacity opacity-80 group" href="#">
<span class="material-symbols-outlined" style="font-variation-settings: 'FILL' 1;">add_circle</span>
<span>New Claim</span>
</a>
<a class="flex items-center space-x-3 px-4 py-3 text-[#47607e] hover:bg-[#f3f4f0] transition-transform duration-200 hover:translate-x-1 group" href="#">
<span class="material-symbols-outlined group-hover:opacity-80">description</span>
<span>Documents</span>
</a>
<a class="flex items-center space-x-3 px-4 py-3 text-[#47607e] hover:bg-[#f3f4f0] transition-transform duration-200 hover:translate-x-1 group" href="#">
<span class="material-symbols-outlined group-hover:opacity-80">settings</span>
<span>Settings</span>
</a>
</div>
<div class="pt-8 border-t border-outline-variant/15 space-y-4">
<button class="w-full bg-primary text-on-primary py-3 rounded-lg font-bold flex items-center justify-center space-x-2 shadow-sm">
<span class="material-symbols-outlined text-sm">contact_support</span>
<span>Emergency Assist</span>
</button>
<div class="flex flex-col space-y-2 px-4 text-xs text-secondary">
<a class="hover:text-primary" href="#">Legal</a>
<a class="hover:text-primary" href="#">Privacy</a>
</div>
</div>
</nav>
<!-- Main Content Canvas -->
<main class="lg:ml-64 min-h-screen">
<!-- Top Branding/Context (Mobile Only View Trigger) -->
<header class="lg:hidden flex items-center justify-between p-6 bg-surface">
<span class="font-headline font-extrabold text-xl">Cymbal</span>
<span class="material-symbols-outlined">menu</span>
</header>
<section class="max-w-5xl mx-auto px-6 py-12 lg:py-20">
<!-- Header Section: Atmospheric Statement -->
<div class="mb-16">
<h1 class="text-4xl lg:text-5xl font-extrabold text-on-surface tracking-tight mb-4">New Claim Submission</h1>
<p class="text-on-surface-variant max-w-xl text-lg leading-relaxed">
                    Take a deep breath. Our concierge team is ready to guide you through this process. We'll handle the complexities so you can focus on what matters.
                </p>
</div>
<!-- Multi-step Progress Indicator (High-End Asymmetric) -->
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
<!-- Form Content: Bento-style Layout -->
<div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
<!-- Left Column: Primary Fields -->
<div class="lg:col-span-7 space-y-12">
<div class="bg-surface-container-low p-8 rounded-xl">
<h2 class="text-xl font-bold mb-8 flex items-center space-x-2">
<span class="material-symbols-outlined text-secondary">event_note</span>
<span>When &amp; Where?</span>
</h2>
<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
<div class="flex flex-col space-y-2">
<label class="text-xs font-bold uppercase tracking-widest text-secondary">Date of Incident</label>
<input class="bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 rounded-md outline-none transition-all" type="date"/>
<p class="text-[10px] text-on-surface-variant">Choose the exact date the accident occurred.</p>
</div>
<div class="flex flex-col space-y-2">
<label class="text-xs font-bold uppercase tracking-widest text-secondary">Approximate Time</label>
<input class="bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 rounded-md outline-none transition-all" type="time"/>
<p class="text-[10px] text-on-surface-variant">An estimate is fine if exact time is unknown.</p>
</div>
<div class="md:col-span-2 flex flex-col space-y-2">
<label class="text-xs font-bold uppercase tracking-widest text-secondary">Location</label>
<div class="relative">
<input class="w-full bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 pl-12 rounded-md outline-none transition-all" placeholder="Street, City, State" type="text"/>
<span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-secondary">location_on</span>
</div>
</div>
</div>
</div>
<div class="bg-surface-container-low p-8 rounded-xl">
<h2 class="text-xl font-bold mb-8 flex items-center space-x-2">
<span class="material-symbols-outlined text-secondary">description</span>
<span>Incident Description</span>
</h2>
<div class="flex flex-col space-y-2">
<label class="text-xs font-bold uppercase tracking-widest text-secondary">Tell us what happened</label>
<textarea class="bg-surface-container-lowest border-none ring-1 ring-outline-variant/15 focus:ring-2 focus:ring-primary p-4 rounded-md outline-none transition-all resize-none" placeholder="Briefly describe the circumstances of the incident..." rows="5"></textarea>
<p class="text-[10px] text-on-surface-variant">Include details about road conditions, other vehicles involved, and any immediate actions taken.</p>
</div>
</div>
</div>
<!-- Right Column: Visual Evidence/Context -->
<div class="lg:col-span-5 space-y-8">
<div class="bg-surface-container-lowest border-2 border-dashed border-outline-variant/30 rounded-xl p-10 flex flex-col items-center text-center space-y-6 hover:bg-surface-container-high/30 transition-colors cursor-pointer">
<div class="w-16 h-16 bg-surface-container-highest rounded-full flex items-center justify-center">
<span class="material-symbols-outlined text-primary text-3xl">add_a_photo</span>
</div>
<div>
<h3 class="font-bold text-lg">Upload Accident Photos</h3>
<p class="text-sm text-on-surface-variant mt-2 px-4">Drag and drop images or click to browse. Clear photos of damage help us process your claim faster.</p>
</div>
<button class="bg-surface-container-highest text-on-surface px-6 py-2 rounded-md font-bold text-sm">Select Files</button>
</div>
<div class="relative h-64 rounded-xl overflow-hidden shadow-sm group">
<img alt="Safety first visual" class="w-full h-full object-cover" data-alt="Modern high-end silver luxury SUV parked safely on the side of a clean urban street with soft morning sunlight" src="https://lh3.googleusercontent.com/aida-public/AB6AXuA_gmVTAzDAxhgITTVNHd-fBYDoieXEBDs2V0sO_CXF2TevFXw9vsGfLTKGDrtSQSgO_2ja5cWxwEfOM5vdz0zo0mO39934RxfS46UHKXe2Qh3uweRSqqRQc6hYtRhWmhn1wmrfuIkdbD-vXPqUinftQPGZwvKknElnutpVM0UJ4uEFjRxA2c5uM6I-kjP4QOjXQ0826s1Hfbcyta5oJk5Ktnrl5hEdsDp1OF4CyobzLmZWWD7cGzLHP9dfRLwdcxTdWc6U1EWmxTY"/>
<div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent flex items-end p-6">
<div class="text-white">
<span class="bg-tertiary-fixed text-on-tertiary-fixed text-[10px] px-2 py-1 rounded font-bold uppercase tracking-tighter">Pro Tip</span>
<p class="text-sm mt-2 font-medium">Ensure everyone is safe before documenting damage.</p>
</div>
</div>
</div>
<!-- Trust Factor Card -->
<div class="bg-primary-container p-8 rounded-xl text-on-primary-container">
<span class="material-symbols-outlined text-tertiary-fixed text-4xl mb-4" style="font-variation-settings: 'FILL' 1;">verified_user</span>
<h4 class="text-white font-bold text-lg mb-2">Secure &amp; Private</h4>
<p class="text-sm opacity-80 leading-relaxed">Your data is encrypted using banking-grade security standards. Our concierge team reviews each claim personally within 24 hours.</p>
</div>
</div>
</div>
<!-- Footer Actions -->
<div class="mt-16 pt-12 border-t border-outline-variant/15 flex flex-col md:flex-row justify-between items-center space-y-6 md:space-y-0">
<button class="text-secondary font-bold flex items-center space-x-2 hover:text-primary transition-colors">
<span class="material-symbols-outlined">arrow_back</span>
<span>Save for later &amp; Exit</span>
</button>
<div class="flex items-center space-x-4 w-full md:w-auto">
<button class="flex-1 md:flex-none px-8 py-4 bg-surface-container-highest text-on-surface rounded-md font-bold hover:bg-surface-dim transition-colors">Cancel</button>
<button class="flex-1 md:flex-none px-12 py-4 bg-gradient-to-br from-primary to-primary-container text-on-primary rounded-md font-bold shadow-lg flex items-center justify-center space-x-2 hover:opacity-90 transition-opacity">
<span>Continue to Step 2</span>
<span class="material-symbols-outlined">arrow_forward</span>
</button>
</div>
</div>
</section>
</main>
<!-- Footer (Authority: Shared Components JSON) -->
<footer class="w-full border-t-0 mt-20 bg-[#f9faf5] text-[#47607e] font-['Inter'] text-xs uppercase tracking-widest py-12 px-8 flex flex-col md:flex-row justify-between items-center max-w-7xl mx-auto">
<p>© 2024 Cymbal Insurance. A Digital Concierge Experience.</p>
<div class="flex space-x-8 mt-6 md:mt-0">
<a class="hover:text-[#000000] transition-colors" href="#">Support Center</a>
<a class="hover:text-[#000000] transition-colors" href="#">Privacy Policy</a>
<a class="hover:text-[#000000] transition-colors" href="#">Terms of Service</a>
<a class="hover:text-[#000000] transition-colors" href="#">Contact Specialist</a>
</div>
</footer>
</body></html>