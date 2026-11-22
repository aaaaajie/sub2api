<template>
  <div
    class="bottom-nav-anchor pointer-events-none flex justify-center"
    :class="props.embedded ? 'w-full' : 'fixed inset-x-0 z-40'"
    :style="props.embedded ? undefined : { bottom: 'calc(env(safe-area-inset-bottom, 0px) + 2rem)' }"
  >
    <nav
      :aria-label="t('nav.primaryNavigation')"
      class="bottom-nav pointer-events-auto flex items-center gap-1"
      :class="props.embedded
        ? 'w-full justify-center sm:gap-2'
        : 'w-[calc(100%_-_2rem)] max-w-[36rem] rounded-[1.2rem] border border-white/55 bg-white/65 p-1 shadow-[0_16px_45px_-18px_rgba(15,23,42,0.18),0_4px_12px_-5px_rgba(15,23,42,0.08)] backdrop-blur-2xl transition-[background-color,border-color,box-shadow,opacity] duration-300 hover:bg-white/80 focus-within:bg-white/80 dark:border-white/[0.05] dark:bg-[#181b1f]/70 dark:shadow-[0_18px_48px_-18px_rgba(0,0,0,0.5)] dark:hover:bg-[#181b1f]/85 dark:focus-within:bg-[#181b1f]/85 sm:w-auto sm:max-w-none sm:gap-2 sm:p-1.5'"
    >
      <RouterLink
        v-for="item in items"
        :key="item.to"
        :to="item.to"
        :data-tour="item.to === '/keys' ? 'user-my-keys' : undefined"
        :aria-current="activePath === item.to ? 'page' : undefined"
        @click="openAdminSidebar(item.to)"
        class="relative flex h-12 sm:h-14 min-h-[3rem] min-w-0 flex-1 flex-col items-center justify-center gap-0.5 rounded-xl px-1 py-1.5 text-[9px] font-medium leading-tight text-gray-600 transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-400/60 motion-reduce:transition-none dark:text-gray-300 dark:focus-visible:ring-gray-500/70 sm:min-w-[6rem] sm:gap-1 sm:px-3 sm:py-2 sm:text-[10px]"
        :class="activePath === item.to
          ? 'bottom-nav-selected text-white dark:text-gray-900'
          : 'hover:bg-gray-100/90 dark:hover:bg-white/[0.08]'"
      >
        <span v-if="activePath === item.to" aria-hidden="true" class="bottom-nav-active pointer-events-none absolute inset-0 rounded-xl bg-gray-900 shadow-[0_7px_18px_-7px_rgba(17,24,39,0.58)] ring-1 ring-black/5 dark:bg-gray-100 dark:shadow-[0_7px_18px_-7px_rgba(0,0,0,0.65)] dark:ring-white/30" />
        <Icon
          class="relative"
          :name="item.icon"
          size="md"
          :stroke-width="route.path === item.to ? 2 : 1.6"
        />
        <span class="relative whitespace-nowrap">{{ t(item.label) }}</span>
      </RouterLink>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { useRoute, RouterLink } from 'vue-router'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import Icon from '@/components/icons/Icon.vue'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'

const props = withDefaults(defineProps<{ embedded?: boolean; activePath?: string }>(), { embedded: false })
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const activePath = computed(() => props.activePath ?? route.path)
const paymentEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.payment))

const coreItems = [
  { to: '/data-dashboard', label: 'nav.dataDashboard', icon: 'chart' },
  { to: '/keys', label: 'nav.apiKeys', icon: 'key' },
  { to: '/profile', label: 'nav.myInfo', icon: 'userCircle' },
] as const

const rechargeItem = { to: '/purchase', label: 'nav.creditRecharge', icon: 'creditCard' } as const
const adminItem = { to: '/admin/dashboard', label: 'nav.adminPortal', icon: 'cog' } as const

const items = computed(() => {
  return [
    ...coreItems,
    ...(!authStore.isSimpleMode && paymentEnabled.value ? [rechargeItem] : []),
    ...(authStore.isAdmin ? [adminItem] : [])
  ]
})

function openAdminSidebar(path: string) {
  if (path === adminItem.to) appStore.setMobileOpen(true)
}

defineExpose({ itemCount: computed(() => items.value.length) })
</script>

<style scoped>
.bottom-nav {
  opacity: 1;
}

@media (hover: hover) and (pointer: fine) {
  .bottom-nav {
    opacity: 0.66;
  }

  .bottom-nav:hover,
  .bottom-nav:focus-within {
    opacity: 1;
  }
}
</style>
