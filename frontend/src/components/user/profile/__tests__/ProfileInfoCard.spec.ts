import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import ProfileInfoCard from '@/components/user/profile/ProfileInfoCard.vue'
import ProfileIdentityEditor from '@/components/user/profile/ProfileIdentityEditor.vue'
import type { User } from '@/types'

vi.mock('vue-router', () => ({
  useRoute: () => ({
    fullPath: '/profile'
  })
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    user: null
  })
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn()
  })
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string>) => {
        if (key === 'profile.accountBalance') return 'Account Balance'
        if (key === 'profile.concurrencyLimit') return 'Concurrency Limit'
        if (key === 'profile.memberSince') return 'Member Since'
        if (key === 'profile.administrator') return 'Administrator'
        if (key === 'profile.user') return 'User'
        if (key === 'profile.authBindings.providers.email') return 'Email'
        if (key === 'profile.authBindings.providers.linuxdo') return 'LinuxDo'
        if (key === 'profile.authBindings.providers.dingtalk') return 'DingTalk'
        if (key === 'profile.authBindings.providers.wechat') return 'WeChat'
        if (key === 'profile.authBindings.providers.oidc') return params?.providerName || 'OIDC'
        if (key === 'profile.authBindings.source.avatar') {
          return `Avatar synced from ${params?.providerName || 'provider'}`
        }
        if (key === 'profile.authBindings.source.username') {
          return `Username synced from ${params?.providerName || 'provider'}`
        }
        return key
      }
    })
  }
})

function createUser(overrides: Partial<User> = {}): User {
  return {
    id: 5,
    username: 'alice',
    email: 'alice@example.com',
    avatar_url: null,
    role: 'user',
    balance: 10,
    concurrency: 2,
    status: 'active',
    allowed_groups: null,
    balance_notify_enabled: true,
    balance_notify_threshold: null,
    balance_notify_extra_emails: [],
    created_at: '2026-04-20T00:00:00Z',
    updated_at: '2026-04-20T00:00:00Z',
    ...overrides
  }
}

describe('ProfileInfoCard', () => {
  it.each([
    { avatar_source: 'dingtalk', username_source: 'dingtalk' },
    { profile_sources: { avatar: { provider: 'dingtalk' }, username: { provider: 'dingtalk' } } },
  ])('shows DingTalk as the synchronized profile source', async (sources) => {
    const wrapper = mount(ProfileIdentityEditor, {
      props: { kind: 'avatar', user: createUser(sources), oidcProviderName: 'OIDC' },
      global: { stubs: { BaseDialog: true, Icon: true } },
    })
    expect(wrapper.text()).toContain('Avatar synced from DingTalk')
    await wrapper.setProps({ kind: 'username' })
    expect(wrapper.text()).toContain('Username synced from DingTalk')
    wrapper.unmount()
  })

  it('renders basic account information inside the new overview shell', () => {
    const wrapper = mount(ProfileInfoCard, {
      props: {
        user: createUser()
      },
      global: {
        stubs: {
          BaseDialog: true, Icon: true
        }
      }
    })

    expect(wrapper.text()).toContain('alice@example.com')
    expect(wrapper.text()).toContain('alice')
    expect(wrapper.find('[data-testid="profile-basics-panel"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="profile-auth-bindings-panel"]').exists()).toBe(false)
  })

  it('renders third-party source hints from profile sources', async () => {
    const wrapper = mount(ProfileIdentityEditor, {
      props: {
        oidcProviderName: 'OIDC',
        kind: 'avatar',
        user: createUser({
          avatar_url: 'https://cdn.example.com/linuxdo.png',
          profile_sources: {
            avatar: { provider: 'linuxdo', source: 'linuxdo' },
            username: { provider: 'linuxdo', source: 'linuxdo' }
          }
        })
      },
      global: {
        stubs: {
          BaseDialog: true, Icon: true
        }
      }
    })

    expect(wrapper.text()).toContain('Avatar synced from LinuxDo')
    await wrapper.setProps({ kind: 'username' })
    expect(wrapper.text()).toContain('Username synced from LinuxDo')
  })

  it('uses the configured OIDC provider name in source hints', () => {
    const wrapper = mount(ProfileIdentityEditor, {
      props: {
        user: createUser({
          profile_sources: {
            username: { provider: 'oidc', source: 'oidc' }
          }
        }),
        oidcProviderName: 'ExampleID',
        kind: 'username'
      },
      global: {
        stubs: {
          BaseDialog: true, Icon: true
        }
      }
    })

    expect(wrapper.text()).toContain('Username synced from ExampleID')
  })

  it('does not display synthetic oauth-only emails as a real bound email', () => {
    const wrapper = mount(ProfileInfoCard, {
      props: {
        user: createUser({
          email: 'legacy-user@oidc-connect.invalid',
          email_bound: false,
          auth_bindings: {
            email: { bound: false }
          }
        })
      },
      global: {
        stubs: {
          BaseDialog: true, Icon: true
        }
      }
    })

    expect(wrapper.text()).not.toContain('legacy-user@oidc-connect.invalid')
  })

  it('does not display synthetic oauth-only emails when only legacy identity bindings mark email as unbound', () => {
    const wrapper = mount(ProfileInfoCard, {
      props: {
        user: createUser({
          email: 'legacy-user@wechat-connect.invalid',
          identity_bindings: {
            email: { bound: false }
          }
        })
      },
      global: {
        stubs: {
          BaseDialog: true, Icon: true
        }
      }
    })

    expect(wrapper.text()).not.toContain('legacy-user@wechat-connect.invalid')
  })

  it('shows a compact summary without duplicate balance or editing panels', () => {
    const wrapper = mount(ProfileInfoCard, {
      props: {
        user: createUser()
      },
      global: {
        stubs: {
          BaseDialog: true, Icon: true
        }
      }
    })

    expect(wrapper.get('[data-testid="profile-overview-hero"]').text()).toContain('alice@example.com')
    expect(wrapper.find('[data-testid="profile-overview-metric-balance"]').exists()).toBe(false)
    expect(wrapper.get('[data-testid="profile-overview-metric-concurrency"]').text()).toContain('Concurrency Limit')
    expect(wrapper.get('[data-testid="profile-overview-metric-member-since"]').text()).toContain('Member Since')
    expect(wrapper.find('[data-testid="profile-info-summary-grid"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="profile-basics-panel"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="profile-auth-bindings-panel"]').exists()).toBe(false)
  })
  it('opens the corresponding editor and closes it after saving', async () => {
    const wrapper = mount(ProfileInfoCard, {
      props: { user: createUser() },
      global: { stubs: {
        Icon: true,
        BaseDialog: { props: ['show'], template: '<div v-if="show"><slot /></div>' },
        ProfileIdentityEditor: { props: ['kind'], emits: ['saved'], template: '<button data-testid="editor" @click="$emit(\'saved\')">{{ kind }}</button>' }
      } }
    })
    expect(wrapper.find('[data-testid="editor"]').exists()).toBe(false)
    await wrapper.get('.profile-identity__avatar').trigger('click')
    expect(wrapper.get('[data-testid="editor"]').text()).toBe('avatar')
    await wrapper.get('[data-testid="editor"]').trigger('click')
    expect(wrapper.find('[data-testid="editor"]').exists()).toBe(false)
    await wrapper.get('.profile-name-edit').trigger('click')
    expect(wrapper.get('[data-testid="editor"]').text()).toBe('username')
    wrapper.unmount()
  })

})
