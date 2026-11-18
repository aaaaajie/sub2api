import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import ProfileView from '@/views/user/ProfileView.vue'

const {
  fetchPublicSettingsMock,
  refreshUserMock,
  authState
} = vi.hoisted(() => ({
  fetchPublicSettingsMock: vi.fn(),
  refreshUserMock: vi.fn(),
  authState: {
    user: null as Record<string, unknown> | null,
    refreshUser: vi.fn()
  }
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => authState
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    fetchPublicSettings: fetchPublicSettingsMock
  })
}))

vi.mock('@/utils/format', () => ({
  formatDate: () => 'April 2026'
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

describe('ProfileView', () => {
  beforeEach(() => {
    refreshUserMock.mockReset()
    fetchPublicSettingsMock.mockReset()
    refreshUserMock.mockResolvedValue(undefined)
    authState.refreshUser = refreshUserMock
    authState.user = {
      id: 1,
      username: 'alice',
      email: 'alice@example.com',
      role: 'user',
      balance: 10,
      concurrency: 2,
      status: 'active',
      allowed_groups: null,
      balance_notify_enabled: true,
      balance_notify_threshold: null,
      balance_notify_extra_emails: [],
      created_at: '2026-04-20T00:00:00Z',
      updated_at: '2026-04-20T00:00:00Z'
    }
    fetchPublicSettingsMock.mockResolvedValue({
      contact_info: '',
      balance_low_notify_enabled: false,
      balance_low_notify_threshold: 0,
      linuxdo_oauth_enabled: true,
      wechat_oauth_enabled: true,
      wechat_oauth_open_enabled: true,
      wechat_oauth_mp_enabled: false,
      oidc_oauth_enabled: true,
      oidc_oauth_provider_name: 'OIDC'
    })
  })

  function renderPage() {
    return mount(ProfileView, {
      global: { stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        ProfileInfoCard: { template: '<div data-testid="summary" />' },
        ProfileIdentityBindingsSection: { template: '<div data-testid="bindings" />' },
        ProfileBalanceNotifyCard: { template: '<div data-testid="notifications" />' },
        ProfilePasswordForm: { template: '<input data-testid="password-draft" />' },
        ProfileTotpCard: { template: '<div data-testid="totp" />' },
        ProfilePasskeyCard: { template: '<div data-testid="passkey" />' }
      } }
    })
  }

  it('opens security settings first and preserves drafts between sections', async () => {
    const wrapper = renderPage()
    await flushPromises()
    expect(wrapper.find('[data-testid="summary"]').exists()).toBe(true)
    expect(wrapper.get('[data-testid="totp"]').isVisible()).toBe(true)
    const safetyInstance = wrapper.get('[data-testid="totp"]').element
    await wrapper.get('[aria-controls="profile-password-editor"]').trigger('click')
    await wrapper.get('[data-testid="password-draft"]').setValue('unsaved draft')
    await wrapper.get('#profile-tab-notifications').trigger('click')
    await wrapper.get('#profile-tab-security').trigger('click')
    expect(wrapper.get('[data-testid="totp"]').element).toBe(safetyInstance)
    expect((wrapper.get('[data-testid="password-draft"]').element as HTMLInputElement).value).toBe('unsaved draft')
    wrapper.unmount()
  })

  it('shows disabled notifications explicitly and supports keyboard tab navigation', async () => {
    const wrapper = renderPage()
    await flushPromises()
    await wrapper.get('#profile-tab-security').trigger('keydown', { key: 'End' })
    expect(wrapper.get('#profile-tab-notifications').attributes('aria-selected')).toBe('true')
    expect(wrapper.get('#profile-panel-notifications').text()).toContain('profile.settings.notificationsUnavailable')
    expect(wrapper.find('[data-testid="notifications"]').exists()).toBe(false)
    wrapper.unmount()
  })

  it('reports settings failures and retries without pretending features are disabled', async () => {
    fetchPublicSettingsMock.mockRejectedValueOnce(new Error('offline'))
    const wrapper = renderPage()
    await flushPromises()
    await wrapper.get('#profile-tab-security').trigger('click')
    expect(wrapper.get('[role="alert"]').text()).toContain('profile.settings.loadError')
    expect(wrapper.find('[data-testid="passkey"]').exists()).toBe(false)
    await wrapper.get('[role="alert"] button').trigger('click')
    await flushPromises()
    expect(wrapper.find('[role="alert"]').exists()).toBe(false)
    expect(wrapper.get('[data-testid="passkey"]').isVisible()).toBe(true)
    wrapper.unmount()
  })
})
