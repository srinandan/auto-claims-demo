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
  <div class="min-h-screen bg-surface text-on-surface antialiased">
    <template v-if="user">
      <TopNavBar />
      <SideNavBar />
      <main class="lg:ml-64 pt-24 min-h-screen">
        <router-view />
      </main>
      <Footer />
    </template>
    <template v-else>
      <router-view />
    </template>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import TopNavBar from './components/layout/TopNavBar.vue'
import SideNavBar from './components/layout/SideNavBar.vue'
import Footer from './components/layout/Footer.vue'

const route = useRoute()
const user = ref(null)

const checkUser = () => {
  const userData = localStorage.getItem('user')
  if (userData) {
    try {
      user.value = JSON.parse(userData)
    } catch (e) {
      user.value = null
    }
  } else {
    user.value = null
  }
}

watch(() => route.path, () => {
  checkUser()
})

onMounted(() => {
  checkUser()
})
</script>
