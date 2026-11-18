<template>
  <div class="min-h-screen bg-gray-50 dark:bg-dark-950">
    <!-- Admin navigation stays in the sidebar; user navigation uses the floating bar. -->
    <AppSidebar v-if="isAdminRoute" />

    <!-- Main Content Area -->
    <div
      class="relative min-h-screen transition-all duration-300"
      :class="isAdminRoute ? (sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64') : ''"
    >
      <!-- Header -->
      <AppHeader
        class="app-page-header"
        :show-sidebar-toggle="isAdminRoute"
        :class="animateLoginHeader && isDataDashboardRoute ? 'dashboard-entry-header' : undefined"
      />

      <!-- Main Content -->
      <main
        class="app-page-content p-4 md:p-6 lg:p-8"
        :class="animateLoginHeader && isDataDashboardRoute ? 'dashboard-entry-main' : undefined"
        :style="showBottomNav ? { paddingBottom: 'calc(8rem + env(safe-area-inset-bottom, 0px))' } : undefined"
      >
        <slot />
      </main>
    </div>

    <UserBottomNav v-if="showBottomNav" />
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { consumeLoginEntry } from '@/utils/dashboard-entry'
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import UserBottomNav from './UserBottomNav.vue'

const animateLoginHeader = consumeLoginEntry()
const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.user?.role === 'admin')
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
const isDataDashboardRoute = computed(() => route.path === '/data-dashboard')
const showBottomNav = computed(() => !isAdminRoute.value)

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
})

defineExpose({ replayTour })
</script>

<style scoped>
.dashboard-entry-header {
  animation: dashboard-header-enter 620ms cubic-bezier(0.2, 0.7, 0.2, 1) both;
}

.dashboard-entry-main {
  animation: dashboard-main-enter 700ms cubic-bezier(0.2, 0.7, 0.2, 1) 90ms both;
}

@keyframes dashboard-header-enter {
  from { opacity: 0; transform: translateY(-8px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes dashboard-main-enter {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (prefers-reduced-motion: reduce) {
  .dashboard-entry-header,
  .dashboard-entry-main {
    animation: none;
  }
}
</style>
