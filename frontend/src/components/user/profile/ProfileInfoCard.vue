<template>
  <section data-testid="profile-overview-hero" class="profile-identity">
    <button type="button" class="profile-identity__avatar" :aria-label="t('profile.avatar.title')" @click="editor = 'avatar'">
      <img v-if="avatarUrl" :src="avatarUrl" :alt="displayName" />
      <span v-else>{{ avatarInitial }}</span>
      <span class="profile-avatar-edit" aria-hidden="true"><Icon name="edit" size="sm" /></span>
    </button>
    <div class="min-w-0 flex-1">
      <div class="profile-name">
        <h2 class="truncate text-xl font-semibold tracking-tight text-gray-900 dark:text-gray-100">{{ displayName }}</h2>
        <button type="button" class="profile-name-edit" :aria-label="t('profile.editProfile')" @click="editor = 'username'"><Icon name="edit" size="sm" /></button>
      </div>
      <p v-if="primaryEmailDisplay && primaryEmailDisplay !== displayName" class="mt-1 truncate text-sm text-gray-500 dark:text-gray-400">{{ primaryEmailDisplay }}</p>
      <div class="mt-3 flex flex-wrap gap-x-5 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
        <span data-testid="profile-overview-metric-member-since">{{ t('profile.memberSince') }} <span class="ml-1 text-gray-700 dark:text-gray-300">{{ memberSinceLabel }}</span></span>
        <span data-testid="profile-overview-metric-concurrency">{{ t('profile.concurrencyLimit') }} <span class="ml-1 font-mono text-gray-700 dark:text-gray-300">{{ user?.concurrency ?? 0 }}</span></span>
      </div>
    </div>
    <BaseDialog :show="editor !== null" :title="editor === 'avatar' ? t('profile.avatar.title') : t('profile.editProfile')" width="narrow" @close="editor = null">
      <ProfileIdentityEditor v-if="editor" :kind="editor" :user="user" :oidc-provider-name="oidcProviderName" @saved="editor = null" />
    </BaseDialog>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ProfileIdentityEditor from './ProfileIdentityEditor.vue'
import { useI18n } from 'vue-i18n'
import type { User } from '@/types'
import { useProfileIdentity } from './useProfileIdentity'
const props = withDefaults(defineProps<{ user: User | null; oidcProviderName?: string }>(), { oidcProviderName: 'OIDC' })
const editor = ref<'avatar' | 'username' | null>(null)
const { t } = useI18n()
const { avatarUrl, avatarInitial, displayName, primaryEmailDisplay, memberSinceLabel } = useProfileIdentity(() => props.user, () => 'OIDC')
</script>

<style scoped>
.profile-identity { display: flex; align-items: center; gap: 22px; padding: 12px 4px 30px; }
.profile-identity__avatar { position: relative; display: flex; align-items: center; justify-content: center; flex-shrink: 0; width: 68px; height: 68px; border-radius: 22px; font-size: 24px; font-weight: 500; overflow: hidden; @apply bg-gray-900 text-gray-100 dark:bg-gray-100 dark:text-gray-900; }
.profile-identity__avatar img { width: 100%; height: 100%; object-fit: cover; }
.profile-avatar-edit { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; background: rgb(0 0 0 / 48%); color: white; opacity: 0; transition: opacity 160ms ease; }
.profile-identity__avatar:hover .profile-avatar-edit, .profile-identity__avatar:focus-visible .profile-avatar-edit { opacity: 1; }
.profile-name { display: flex; align-items: center; gap: 8px; width: fit-content; max-width: 100%; }
.profile-name-edit { display: flex; align-items: center; justify-content: center; flex-shrink: 0; width: 30px; height: 30px; border-radius: 9px; opacity: 0; transition: opacity 160ms ease, background-color 160ms ease; @apply text-gray-500 dark:text-gray-400; }
.profile-name:hover .profile-name-edit, .profile-name:focus-within .profile-name-edit { opacity: 1; }
.profile-name-edit:hover { @apply bg-gray-100 dark:bg-white/[0.06]; }
.profile-name-edit:focus-visible, .profile-identity__avatar:focus-visible { outline: 2px solid currentColor; outline-offset: 3px; }
@media (hover: none), (pointer: coarse) {
  .profile-name-edit { opacity: 1; }
  .profile-avatar-edit { inset: auto 0 0; height: 24px; opacity: 1; }
}
@media (prefers-reduced-motion: reduce) { .profile-avatar-edit, .profile-name-edit { transition: none; } }
@media (max-width: 480px) { .profile-identity { gap: 16px; } .profile-identity__avatar { width: 56px; height: 56px; border-radius: 18px; } }
</style>
