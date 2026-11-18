import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount, RouterLinkStub } from '@vue/test-utils'
import { defineComponent } from 'vue'
import HomeView from '../HomeView.vue'

const { appStore, authStore, routeState, router, playRetractMock, getPublicActivityTrendMock } = vi.hoisted(() => ({
  appStore: {
    cachedPublicSettings: {} as Record<string, unknown>,
    siteName: 'Fallback site',
    siteLogo: '',
    docUrl: '',
    publicSettingsLoaded: true,
    fetchPublicSettings: vi.fn(),
  },
  authStore: {
    isAuthenticated: false,
    isAdmin: false,
    user: null as { email?: string } | null,
    checkAuth: vi.fn(),
  },
  routeState: { path: '/home', query: {} as Record<string, unknown> },
  router: { replace: vi.fn(), push: vi.fn() },
  playRetractMock: vi.fn().mockResolvedValue(undefined),
  getPublicActivityTrendMock: vi.fn(),
}))

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
    useRoute: () => routeState,
    useRouter: () => router,
  }
})

vi.mock('@/stores', () => ({
  useAppStore: () => appStore,
  useAuthStore: () => authStore,
}))

vi.mock('@/api/public', () => ({
  getPublicActivityTrend: getPublicActivityTrendMock,
}))

vi.mock('@/utils/featureFlags', () => ({
  FeatureFlags: { modelPlaza: 'modelPlaza' },
  isFeatureFlagEnabled: () => false,
}))

const TokenConsumptionHeroStub = defineComponent({
  props: {
    trendData: { type: Array, default: () => [] },
    presentation: { type: String, default: '' },
  },
  setup(props, { expose }) {
    expose({ playRetractAnimation: playRetractMock })
    return props
  },
  template: '<div data-testid="token-hero" :data-presentation="presentation" :data-points="trendData.length" :data-values="JSON.stringify(trendData.map((point) => point.total_tokens))" />',
})

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key }),
  }
})

function mountHome(settings: Record<string, unknown> = {}) {
  appStore.cachedPublicSettings = {
    site_name: 'Test site',
    site_subtitle: 'Test subtitle',
    ...settings,
  }

  return mount(HomeView, {
    global: {
      stubs: {
        RouterLink: RouterLinkStub,
        BaseDialog: {
          props: ['show'],
          template: '<div v-if="show" data-testid="login-dialog"><slot /></div>',
        },
        LoginView: {
          props: ['embedded', 'redirectTo'],
          emits: ['success'],
          template: '<div data-testid="login-form" :data-redirect="redirectTo"><button data-testid="complete-login" @click="$emit(\'success\')" /></div>',
        },
        UserBottomNav: { template: '<div data-testid="user-bottom-nav" />' },
        TokenConsumptionHero: TokenConsumptionHeroStub,
        Icon: { template: '<span data-testid="icon" />' },
      },
    },
  })
}

describe('HomeView unified landing page', () => {
  beforeEach(() => {
    authStore.isAuthenticated = false
    authStore.isAdmin = false
    authStore.user = null
    authStore.checkAuth.mockClear()
    appStore.fetchPublicSettings.mockClear()
    router.replace.mockReset()
    router.push.mockReset()
    playRetractMock.mockReset().mockResolvedValue(undefined)
    getPublicActivityTrendMock.mockReset().mockResolvedValue({
      trend: Array.from({ length: 10 }, (_, index) => ({
        date: new Date(Date.UTC(2026, 10, 15, index)).toISOString(),
        activity_percent: index === 7 ? 43 : 0,
      })),
    })
    routeState.path = '/home'
    routeState.query = {}
    localStorage.clear()
    vi.spyOn(window, 'matchMedia').mockReturnValue({ matches: false } as MediaQueryList)
  })

  it('renders the real public activity trend without using legacy home settings', async () => {
    const wrapper = mountHome({
      home_content: '<section id="legacy-home">Old content</section>',
      compact_home_enabled: true,
    })

    await flushPromises()

    expect(wrapper.find('[data-testid="token-hero"]').exists()).toBe(true)
    expect(wrapper.find('#legacy-home').exists()).toBe(false)
    const chart = wrapper.get('[data-testid="token-hero"]')
    expect(chart.attributes('data-presentation')).toBe('landing')
    expect(chart.attributes('data-points')).toBe('10')
    expect(JSON.parse(chart.attributes('data-values') ?? '[]')[7]).toBe(43)
    expect(getPublicActivityTrendMock).toHaveBeenCalledOnce()
  })

  it('renders the product capability sections and lets guests continue from the closing action', async () => {
    const wrapper = mountHome()

    expect(wrapper.get('#home-access-title').exists()).toBe(true)
    expect(wrapper.get('#home-management-title').exists()).toBe(true)
    expect(wrapper.get('#home-visibility-title').exists()).toBe(true)

    await wrapper.get('[data-testid="home-final-cta"]').trigger('click')

    expect(router.replace).toHaveBeenCalledWith({
      path: '/home',
      query: { auth: 'login', redirect: '/data-dashboard' },
    })
  })

  it('links authenticated visitors from the closing action to the dashboard', () => {
    authStore.isAuthenticated = true
    const wrapper = mountHome()
    const finalLink = wrapper.findAllComponents(RouterLinkStub).find((link) =>
      link.attributes('data-testid') === 'home-final-cta'
    )

    expect(finalLink?.props('to')).toBe('/data-dashboard')
  })

  it('opens the login dialog route with the data dashboard as the default return target', async () => {
    const wrapper = mountHome()

    await wrapper.get('[data-testid="home-login-button"]').trigger('click')

    expect(router.replace).toHaveBeenCalledWith({
      path: '/home',
      query: { auth: 'login', redirect: '/data-dashboard' },
    })
  })

  it('shows the same home with the login dialog for direct /login visits', () => {
    routeState.path = '/login'
    routeState.query = { redirect: '/keys' }
    const wrapper = mountHome()

    expect(wrapper.get('[data-testid="login-dialog"]').exists()).toBe(true)
    expect(wrapper.get('[data-testid="login-form"]').attributes('data-redirect')).toBe('/keys')
  })

  it('retracts the activity line before entering the requested destination after login', async () => {
    routeState.query = { auth: 'login', redirect: '/keys' }
    const wrapper = mountHome()

    await wrapper.get('[data-testid="complete-login"]').trigger('click')
    await flushPromises()

    expect(playRetractMock).toHaveBeenCalledOnce()
    expect(router.replace).toHaveBeenCalledWith('/keys')
  })

  it('replays the home transition after an OAuth callback when session state matches', async () => {
    authStore.isAuthenticated = true
    routeState.query = { auth_transition: '1', redirect: '/keys' }
    sessionStorage.setItem('sub2api:home-login-transition', '/keys')
    const wrapper = mountHome()
    await flushPromises()

    expect(playRetractMock).toHaveBeenCalledOnce()
    expect(router.replace).toHaveBeenCalledWith('/keys')
    wrapper.unmount()
  })

  it('replaces the guest login button with the authenticated floating navigation', () => {
    authStore.isAuthenticated = true
    const wrapper = mountHome()

    expect(wrapper.find('[data-testid="user-bottom-nav"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="home-login-button"]').exists()).toBe(false)
  })
})
