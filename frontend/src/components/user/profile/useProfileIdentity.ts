import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { User, UserAuthBindingStatus, UserAuthProvider, UserProfileSourceContext } from '@/types'

export function useProfileIdentity(user: () => User | null, oidcProviderName: () => string) {
  const { t } = useI18n()
  function normalizeBindingStatus(binding: boolean | UserAuthBindingStatus | undefined): boolean | null {
    if (typeof binding === 'boolean') {
      return binding
    }
    if (!binding) {
      return null
    }
    if (typeof binding.bound === 'boolean') {
      return binding.bound
    }
    return Boolean(binding.provider_subject || binding.issuer || binding.provider_key)
  }

  function isEmailBound(user: User | null | undefined): boolean {
    if (typeof user?.email_bound === 'boolean') {
      return user.email_bound
    }

    const nested = user?.auth_bindings?.email ?? user?.identity_bindings?.email
    const normalized = normalizeBindingStatus(nested)
    return normalized ?? false
  }

  const avatarUrl = computed(() => user()?.avatar_url?.trim() || '')
  const displayName = computed(() => user()?.username?.trim() || user()?.email?.trim() || t('profile.user'))
  const primaryEmailDisplay = computed(() => {
    const email = user()?.email?.trim() || ''
    if (!email) {
      return ''
    }
    if (email.endsWith('.invalid') && !isEmailBound(user())) {
      return ''
    }
    return email
  })
  const avatarInitial = computed(() => displayName.value.charAt(0).toUpperCase() || 'U')
  const memberSinceLabel = computed(() => {
    const raw = user()?.created_at?.trim()
    if (!raw) {
      return '-'
    }

    const date = new Date(raw)
    if (Number.isNaN(date.getTime())) {
      return '-'
    }

    return new Intl.DateTimeFormat(undefined, {
      year: 'numeric',
      month: 'short',
    }).format(date)
  })

  const providerLabels = computed<Record<UserAuthProvider, string>>(() => ({
    email: t('profile.authBindings.providers.email'),
    linuxdo: t('profile.authBindings.providers.linuxdo'),
    dingtalk: t('profile.authBindings.providers.dingtalk'),
    oidc: t('profile.authBindings.providers.oidc', { providerName: oidcProviderName() }),
    wechat: t('profile.authBindings.providers.wechat'),
    github: 'GitHub',
    google: 'Google'
  }))

  function normalizeProvider(value: string): UserAuthProvider | null {
    const normalized = value.trim().toLowerCase()
    if (
      normalized === 'email' ||
      normalized === 'linuxdo' ||
      normalized === 'dingtalk' ||
      normalized === 'wechat' ||
      normalized === 'github' ||
      normalized === 'google'
    ) {
      return normalized
    }
    if (normalized === 'oidc' || normalized.startsWith('oidc:') || normalized.startsWith('oidc/')) {
      return 'oidc'
    }
    return null
  }

  function readObjectString(source: Record<string, unknown>, ...keys: string[]): string {
    for (const key of keys) {
      const value = source[key]
      if (typeof value === 'string' && value.trim()) {
        return value.trim()
      }
    }
    return ''
  }

  function resolveThirdPartySource(
    rawSource: string | UserProfileSourceContext | null | undefined
  ): { provider: UserAuthProvider; label: string } | null {
    if (!rawSource) {
      return null
    }

    if (typeof rawSource === 'string') {
      const provider = normalizeProvider(rawSource)
      if (!provider || provider === 'email') {
        return null
      }
      return {
        provider,
        label: providerLabels.value[provider]
      }
    }

    const sourceRecord = rawSource as Record<string, unknown>
    const provider = normalizeProvider(
      readObjectString(sourceRecord, 'provider', 'source', 'provider_type', 'auth_provider')
    )
    if (!provider || provider === 'email') {
      return null
    }

    const explicitLabel = readObjectString(
      sourceRecord,
      'provider_label',
      'label',
      'provider_name',
      'providerName'
    )

    return {
      provider,
      label: explicitLabel || providerLabels.value[provider]
    }
  }

  const sourceHints = computed(() => {
    const currentUser = user()
    if (!currentUser) {
      return []
    }

    const hints: Array<{ key: string; text: string }> = []
    const avatarSource = resolveThirdPartySource(
      currentUser.profile_sources?.avatar ?? currentUser.avatar_source
    )
    const usernameSource = resolveThirdPartySource(
      currentUser.profile_sources?.username ??
        currentUser.profile_sources?.display_name ??
        currentUser.profile_sources?.nickname ??
        currentUser.display_name_source ??
        currentUser.username_source ??
        currentUser.nickname_source
    )

    if (avatarSource) {
      hints.push({
        key: 'avatar',
        text: t('profile.authBindings.source.avatar', { providerName: avatarSource.label })
      })
    }

    if (usernameSource) {
      hints.push({
        key: 'username',
        text: t('profile.authBindings.source.username', { providerName: usernameSource.label })
      })
    }

    return hints
  })
  return { avatarUrl, avatarInitial, displayName, primaryEmailDisplay, memberSinceLabel, sourceHints }
}
