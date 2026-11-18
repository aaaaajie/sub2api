<template>
  <AppLayout>
    <div data-testid="profile-shell" class="profile-settings">
      <ProfileInfoCard :user="user" :oidc-provider-name="oidcOAuthProviderName" />
      <div v-if="profileError" class="profile-feedback" role="alert">
        {{ t('profile.settings.profileError') }}
        <button type="button" class="underline underline-offset-4" @click="refreshProfile">{{ t('common.tryAgain') }}</button>
      </div>
      <nav class="profile-tabs" role="tablist" :aria-label="t('profile.settings.navigation')">
        <button v-for="(tab, index) in tabs" :id="`profile-tab-${tab.id}`" :key="tab.id"
          type="button" role="tab" :aria-selected="activeTab === tab.id"
          :aria-controls="`profile-panel-${tab.id}`" :tabindex="activeTab === tab.id ? 0 : -1"
          @click="selectTab(tab.id)" @keydown="navigateTabs($event, index)">{{ tab.label }}</button>
      </nav>

      <div v-if="visitedTabs.has('security')" id="profile-panel-security" v-show="activeTab === 'security'" role="tabpanel" aria-labelledby="profile-tab-security" tabindex="0" class="profile-panel">
        <template v-if="settingsReady">
          <div class="profile-section profile-bindings" data-testid="profile-auth-bindings-panel">
            <ProfileIdentityBindingsSection :user="user" :linuxdo-enabled="linuxdoOAuthEnabled"
              :dingtalk-enabled="dingtalkOAuthEnabled" :oidc-enabled="oidcOAuthEnabled" :oidc-provider-name="oidcOAuthProviderName"
              :wechat-enabled="wechatOAuthEnabled" :wechat-open-enabled="wechatOAuthOpenEnabled" :wechat-mp-enabled="wechatOAuthMPEnabled" embedded compact />
          </div>
          <section class="profile-section">
            <div class="profile-setting-row">
              <div><h3>{{ t('profile.changePassword') }}</h3><p>{{ t('profile.settings.passwordHint') }}</p></div>
              <button type="button" class="btn btn-secondary btn-sm" :aria-expanded="passwordExpanded" aria-controls="profile-password-editor" @click="passwordExpanded = !passwordExpanded">{{ passwordExpanded ? t('common.cancel') : t('common.edit') }}</button>
            </div>
            <div id="profile-password-editor" v-show="passwordExpanded" class="profile-password-editor"><ProfilePasswordForm embedded /></div>
          </section>
          <div class="profile-section profile-security-section"><ProfileTotpCard /></div>
          <div class="profile-section profile-security-section"><ProfilePasskeyCard :enabled="passkeyEnabled" /></div>
        </template>
        <p v-else-if="settingsLoading" class="profile-feedback" role="status">{{ t('common.loading') }}</p>
      </div>

      <div v-if="visitedTabs.has('notifications')" id="profile-panel-notifications" v-show="activeTab === 'notifications'" role="tabpanel" aria-labelledby="profile-tab-notifications" tabindex="0" class="profile-panel">
        <template v-if="settingsReady">
          <div v-if="user && balanceLowNotifyEnabled" class="profile-security-section">
            <ProfileBalanceNotifyCard :enabled="user.balance_notify_enabled ?? true" :threshold="user.balance_notify_threshold"
              :extra-emails="user.balance_notify_extra_emails ?? []" :system-default-threshold="systemDefaultThreshold" :user-email="user.email" />
          </div>
          <div v-else-if="!balanceLowNotifyEnabled" class="profile-setting-row">
            <div><h3>{{ t('profile.balanceNotify.title') }}</h3><p>{{ t('profile.settings.notificationsUnavailable') }}</p></div>
            <span class="profile-unavailable">{{ t('profile.settings.unavailable') }}</span>
          </div>
        </template>
        <p v-else-if="settingsLoading" class="profile-feedback" role="status">{{ t('common.loading') }}</p>
      </div>
      <div v-if="settingsError" class="profile-feedback" role="alert">
        {{ t('profile.settings.loadError') }}
        <button type="button" class="underline underline-offset-4" @click="loadSettings">{{ t('common.tryAgain') }}</button>
      </div>
      <details v-if="contactInfo" class="profile-support">
        <summary>{{ t('common.contactSupport') }}</summary>
        <p>{{ contactInfo }}</p>
      </details>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import ProfileBalanceNotifyCard from '@/components/user/profile/ProfileBalanceNotifyCard.vue'
import ProfileIdentityBindingsSection from '@/components/user/profile/ProfileIdentityBindingsSection.vue'
import ProfileInfoCard from '@/components/user/profile/ProfileInfoCard.vue'
import ProfilePasswordForm from '@/components/user/profile/ProfilePasswordForm.vue'
import ProfileTotpCard from '@/components/user/profile/ProfileTotpCard.vue'
import ProfilePasskeyCard from '@/components/user/profile/ProfilePasskeyCard.vue'
import { isWeChatWebOAuthEnabled } from '@/api/auth'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const user = computed(() => authStore.user)

type ProfileTab = 'security' | 'notifications'
const activeTab = ref<ProfileTab>('security')
const visitedTabs = ref(new Set<ProfileTab>(['security']))
const passwordExpanded = ref(false)
const settingsLoading = ref(true)
const settingsError = ref(false)
const settingsReady = ref(false)
const profileError = ref(false)
const tabs = computed(() => (['security', 'notifications'] as const).map(id => ({ id, label: t(`profile.settings.tabs.${id}`) })))
function selectTab(tab: ProfileTab) {
  activeTab.value = tab
  visitedTabs.value.add(tab)
}
function navigateTabs(event: KeyboardEvent, index: number) {
  let next: number
  if (event.key === 'ArrowRight') next = (index + 1) % tabs.value.length
  else if (event.key === 'ArrowLeft') next = (index + tabs.value.length - 1) % tabs.value.length
  else if (event.key === 'Home') next = 0
  else if (event.key === 'End') next = tabs.value.length - 1
  else return
  event.preventDefault()
  selectTab(tabs.value[next]!.id)
  document.getElementById(`profile-tab-${tabs.value[next]!.id}`)?.focus()
}
async function refreshProfile() {
  profileError.value = false
  try { await authStore.refreshUser() } catch { profileError.value = true }
}
const contactInfo = ref('')
const balanceLowNotifyEnabled = ref(false)
const systemDefaultThreshold = ref(0)
const linuxdoOAuthEnabled = ref(false)
const dingtalkOAuthEnabled = ref(false)
const wechatOAuthEnabled = ref(false)
const wechatOAuthOpenEnabled = ref<boolean | undefined>(undefined)
const wechatOAuthMPEnabled = ref<boolean | undefined>(undefined)
const oidcOAuthEnabled = ref(false)
const oidcOAuthProviderName = ref('OIDC')
const passkeyEnabled = ref(false)

async function loadSettings() {
  settingsLoading.value = true
  settingsError.value = false
  try {
      const settings = await appStore.fetchPublicSettings()
      if (!settings) throw new Error('Public settings unavailable')
      contactInfo.value = settings.contact_info || ''
      balanceLowNotifyEnabled.value = settings.balance_low_notify_enabled ?? false
      systemDefaultThreshold.value = settings.balance_low_notify_threshold ?? 0
      linuxdoOAuthEnabled.value = settings.linuxdo_oauth_enabled ?? false
      dingtalkOAuthEnabled.value = settings.dingtalk_oauth_enabled ?? false
      wechatOAuthEnabled.value = isWeChatWebOAuthEnabled(settings)
      wechatOAuthOpenEnabled.value = typeof settings.wechat_oauth_open_enabled === 'boolean'
        ? settings.wechat_oauth_open_enabled
        : undefined
      wechatOAuthMPEnabled.value = typeof settings.wechat_oauth_mp_enabled === 'boolean'
        ? settings.wechat_oauth_mp_enabled
        : undefined
      oidcOAuthEnabled.value = settings.oidc_oauth_enabled ?? false
      oidcOAuthProviderName.value = settings.oidc_oauth_provider_name || 'OIDC'
      passkeyEnabled.value = settings.passkey_enabled === true
      settingsReady.value = true
  } catch {
    settingsError.value = true
  } finally {
    settingsLoading.value = false
  }
}
onMounted(() => {
  void refreshProfile()
  void loadSettings()
})
</script>

<style scoped>
.profile-settings { max-width: 880px; margin: 0 auto; padding-bottom: 100px; }
.profile-tabs { display: flex; gap: 5px; padding: 5px; width: fit-content; margin-bottom: 24px; border-radius: 14px; @apply bg-gray-100/70 dark:bg-white/[0.035]; }
.profile-tabs button { padding: 9px 20px; border-radius: 10px; font-size: 12px; font-weight: 500; transition: color 160ms, background-color 160ms, box-shadow 160ms; @apply text-gray-500 dark:text-gray-400; }
.profile-tabs button:hover { @apply text-gray-900 dark:text-gray-100; }
.profile-tabs button[aria-selected='true'] { @apply bg-white text-gray-900 dark:bg-white/[0.08] dark:text-gray-100; box-shadow: 0 1px 4px rgb(0 0 0 / 4%); }
.profile-tabs button:focus-visible, .profile-support summary:focus-visible { outline: 2px solid currentColor; outline-offset: 3px; }
.profile-panel { padding: 30px; border: 1px solid; border-radius: 22px; animation: profile-panel-enter 180ms ease-out; @apply border-gray-200/70 bg-white/60 dark:border-white/[0.07] dark:bg-dark-950; }
.profile-section + .profile-section { margin-top: 28px; padding-top: 28px; border-top: 1px solid; @apply border-gray-200/60 dark:border-white/[0.06]; }
.profile-setting-row { display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.profile-setting-row h3 { font-size: 14px; font-weight: 600; @apply text-gray-900 dark:text-gray-100; }
.profile-setting-row p { margin-top: 6px; font-size: 12px; line-height: 1.7; @apply text-gray-500 dark:text-gray-400; }
.profile-password-editor { max-width: 420px; margin-top: 22px; }

.profile-security-section :deep(.card) { background: transparent; border: 0; border-radius: 0; box-shadow: none; }
.profile-security-section :deep(.card > div:first-child) { padding: 0 0 18px; border: 0; }
.profile-security-section :deep(.card > div:nth-child(2)) { padding: 0; }
.profile-security-section :deep(h2) { font-size: 14px; font-weight: 600; }
.profile-security-section :deep(h2 + p) { font-size: 12px; line-height: 1.7; }
.profile-bindings :deep(.space-y-4 > .rounded-2xl) { border: 0; border-radius: 0; padding: 16px 0; box-shadow: none; background: transparent; }
.profile-bindings :deep(.space-y-4 > .rounded-2xl + .rounded-2xl) { border-top: 1px solid; @apply border-gray-200/60 dark:border-white/[0.06]; }
.profile-settings :deep(.input) { font-size: 13px; @apply bg-transparent border-gray-200 dark:bg-transparent dark:border-white/10; }
.profile-settings :deep(.btn-secondary) { @apply bg-transparent border-gray-200 dark:bg-transparent dark:border-white/10; }
.profile-settings :deep(.btn-secondary:hover) { @apply bg-gray-100 dark:bg-white/[0.05]; }
.profile-bindings :deep(.badge-gray) { @apply bg-gray-100 text-gray-500 dark:bg-white/[0.05] dark:text-gray-400; }
.profile-settings :deep(.input-label) { font-size: 12px; }
.profile-settings :deep(.btn) { font-size: 12px; }
.profile-feedback { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; padding: 16px 0; font-size: 12px; @apply text-gray-500 dark:text-gray-400; }
.profile-unavailable { flex-shrink: 0; border-radius: 999px; padding: 5px 10px; font-size: 11px; @apply bg-gray-100 text-gray-500 dark:bg-white/[0.05] dark:text-gray-400; }
.profile-support { margin-top: 24px; font-size: 12px; @apply text-gray-500 dark:text-gray-400; }
.profile-support summary { cursor: pointer; width: fit-content; }
.profile-support p { margin-top: 12px; white-space: pre-wrap; overflow-wrap: anywhere; }
@keyframes profile-panel-enter { from { opacity: 0; } to { opacity: 1; } }
@media (max-width: 640px) {
  .profile-panel { padding: 22px 18px; border-radius: 18px; }
  .profile-tabs { width: 100%; }
  .profile-tabs button { flex: 1; padding: 9px 8px; }
}
@media (prefers-reduced-motion: reduce) { .profile-panel { animation: none; } .profile-tabs button { transition: none; } }
</style>
