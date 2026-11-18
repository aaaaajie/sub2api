<template>
  <div
    ref="homePage"
    class="landing-page relative flex min-h-[100svh] flex-col bg-[#fbfcfd] text-gray-950 dark:bg-[#090b0d] dark:text-gray-100"
    :class="{ 'home-auth-success-fade': authTransition === 'nav-morph' }"
  >
    <header class="relative z-20 flex items-center justify-between px-5 py-5 sm:px-8 lg:px-12">
      <RouterLink to="/home" class="flex min-w-0 items-center gap-3 rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500/70">
        <img :src="siteLogo || '/logo.svg'" alt="" class="h-9 w-9 shrink-0 rounded-xl object-contain" />
        <span class="truncate text-sm font-semibold tracking-wide text-gray-800 dark:text-gray-200">{{ siteName }}</span>
      </RouterLink>

      <nav class="flex items-center gap-1.5 sm:gap-2" :aria-label="t('home.publicNavigation')">
        <a
          v-if="docUrl"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="landing-header-action"
          :title="t('home.viewDocs')"
        >
          <Icon name="book" size="md" />
          <span class="sr-only">{{ t('home.viewDocs') }}</span>
        </a>
        <RouterLink
          v-if="showModelPlazaEntry"
          to="/model-plaza"
          class="landing-header-action"
          :title="t('nav.modelPlaza')"
        >
          <Icon name="grid" size="md" />
          <span class="sr-only">{{ t('nav.modelPlaza') }}</span>
        </RouterLink>
        <button
          type="button"
          class="landing-header-action"
          :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          @click="toggleTheme"
        >
          <Icon v-if="isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
          <span class="sr-only">{{ isDark ? t('home.switchToLight') : t('home.switchToDark') }}</span>
        </button>
      </nav>
    </header>

    <main class="relative z-10 flex flex-1 flex-col pb-36 sm:pb-40">
      <section class="landing-hero flex min-h-[calc(100svh-5rem)] flex-col justify-center pb-28 pt-3 sm:pb-32 sm:pt-0" aria-labelledby="home-title">
        <div class="relative z-10 mx-auto mb-4 w-full max-w-5xl px-5 text-center sm:mb-0 sm:px-8">
          <p class="mb-3 text-[10px] font-medium uppercase tracking-[0.24em] text-[#718198] dark:text-[#8391a3]">
            {{ t('home.heroEyebrow') }}
          </p>
          <h1 id="home-title" class="text-3xl font-medium tracking-[-0.055em] text-[#111820] dark:text-white sm:text-5xl lg:text-6xl">
            {{ siteName }}
          </h1>
          <p class="mx-auto mt-4 max-w-xl text-sm leading-6 text-[#687789] dark:text-gray-400 sm:text-base">
            {{ siteSubtitle }}
          </p>
        </div>

        <div class="landing-chart-stage relative left-1/2 mt-1 w-screen -translate-x-1/2 sm:mt-0">
          <TokenConsumptionHero
            ref="chartRef"
            :trend-data="trendData"
            range="recent"
            granularity="hour"
            presentation="landing"
            :loading="trendLoading"
            :error="trendError"
            :auto-reveal="!shouldRunOAuthTransition"
            :initially-revealed="shouldRunOAuthTransition"
            @retry="loadTrend"
          />
        </div>
      </section>

      <section id="home-resources" class="landing-resources" aria-labelledby="home-resources-title" data-reveal>
        <div class="landing-resources-inner">
          <div class="landing-resources-heading">
            <div>
              <p class="landing-section-kicker">{{ t('home.resources.eyebrow') }}</p>
              <h2 id="home-resources-title">{{ t('home.resources.title') }}</h2>
              <p class="landing-resources-description">{{ t('home.resources.description') }}</p>
            </div>
          </div>

          <div class="landing-resource-grid">
            <article class="landing-resource-card landing-resource-card-gpt">
              <div class="landing-resource-card-top">
                <PlatformIcon platform="openai" size="lg" class="landing-resource-logo landing-resource-logo-openai" aria-hidden="true" />
              </div>
              <p class="landing-resource-type">{{ t('home.resources.models.gpt.type') }}</p>
              <h3>ChatGPT</h3>
              <p class="landing-resource-copy">{{ t('home.resources.models.gpt.description') }}</p>
              <div class="landing-resource-rule"><span></span></div>
            </article>

            <article class="landing-resource-card landing-resource-card-ds">
              <div class="landing-resource-card-top">
                <PlatformIcon platform="deepseek" size="lg" class="landing-resource-logo landing-resource-logo-deepseek" aria-hidden="true" />
              </div>
              <p class="landing-resource-type">{{ t('home.resources.models.ds.type') }}</p>
              <h3>DeepSeek</h3>
              <p class="landing-resource-copy">{{ t('home.resources.models.ds.description') }}</p>
              <div class="landing-resource-rule"><span></span></div>
            </article>

            <article class="landing-resource-card landing-resource-card-secai">
              <div class="landing-resource-card-top">
                <Icon name="shield" size="xl" class="landing-resource-logo landing-resource-logo-secai" aria-hidden="true" />
              </div>
              <p class="landing-resource-type">{{ t('home.resources.models.secai.type') }}</p>
              <h3>SecAI</h3>
              <p class="landing-resource-copy">{{ t('home.resources.models.secai.description') }}</p>
              <div class="landing-resource-rule"><span></span></div>
            </article>
          </div>

          <p class="landing-resources-footnote">{{ t('home.resources.footnote') }}</p>
        </div>
      </section>

      <div class="landing-story-shell">
        <section class="landing-story-section" aria-labelledby="home-access-title" data-reveal>
          <div class="landing-section-heading">
            <div>
              <p class="landing-section-kicker">{{ t('home.story.access.eyebrow') }}</p>
              <h2 id="home-access-title">{{ t('home.story.access.title') }}</h2>
            </div>
            <p class="landing-section-description">{{ t('home.story.access.description') }}</p>
          </div>

          <div class="landing-flow" role="group" :aria-label="t('home.story.access.flowLabel')">
            <div class="landing-flow-step">
              <span class="landing-flow-index">01</span>
              <Icon name="key" size="md" class="landing-flow-icon" />
              <div>
                <strong>{{ t('home.story.access.key') }}</strong>
                <span>{{ t('home.story.access.keyDescription') }}</span>
              </div>
            </div>
            <Icon name="arrowRight" size="sm" class="landing-flow-arrow" aria-hidden="true" />
            <div class="landing-flow-step">
              <span class="landing-flow-index">02</span>
              <Icon name="server" size="md" class="landing-flow-icon" />
              <div>
                <strong>{{ t('home.story.access.group') }}</strong>
                <span>{{ t('home.story.access.groupDescription') }}</span>
              </div>
            </div>
            <Icon name="arrowRight" size="sm" class="landing-flow-arrow" aria-hidden="true" />
            <div class="landing-flow-step">
              <span class="landing-flow-index">03</span>
              <Icon name="sparkles" size="md" class="landing-flow-icon" />
              <div>
                <strong>{{ t('home.story.access.models') }}</strong>
                <span>{{ t('home.story.access.modelsDescription') }}</span>
              </div>
            </div>
          </div>
        </section>

        <section class="landing-story-section" aria-labelledby="home-management-title" data-reveal>
          <div class="landing-section-heading">
            <div>
              <p class="landing-section-kicker">{{ t('home.story.management.eyebrow') }}</p>
              <h2 id="home-management-title">{{ t('home.story.management.title') }}</h2>
            </div>
            <p class="landing-section-description">{{ t('home.story.management.description') }}</p>
          </div>

          <div class="landing-capability-list">
            <article class="landing-capability-row">
              <span class="landing-capability-index">01</span>
              <div>
                <h3>{{ t('home.story.management.items.keys.title') }}</h3>
                <p>{{ t('home.story.management.items.keys.description') }}</p>
              </div>
              <Icon name="key" size="md" class="landing-capability-icon" aria-hidden="true" />
            </article>
            <article class="landing-capability-row">
              <span class="landing-capability-index">02</span>
              <div>
                <h3>{{ t('home.story.management.items.groups.title') }}</h3>
                <p>{{ t('home.story.management.items.groups.description') }}</p>
              </div>
              <Icon name="server" size="md" class="landing-capability-icon" aria-hidden="true" />
            </article>
            <article class="landing-capability-row">
              <span class="landing-capability-index">03</span>
              <div>
                <h3>{{ t('home.story.management.items.status.title') }}</h3>
                <p>{{ t('home.story.management.items.status.description') }}</p>
              </div>
              <Icon name="checkCircle" size="md" class="landing-capability-icon" aria-hidden="true" />
            </article>
          </div>
        </section>

        <section class="landing-story-section" aria-labelledby="home-visibility-title" data-reveal>
          <div class="landing-section-heading">
            <div>
              <p class="landing-section-kicker">{{ t('home.story.visibility.eyebrow') }}</p>
              <h2 id="home-visibility-title">{{ t('home.story.visibility.title') }}</h2>
            </div>
            <p class="landing-section-description">{{ t('home.story.visibility.description') }}</p>
          </div>

          <div class="landing-visibility-layout">
            <p class="landing-visibility-note">{{ t('home.story.visibility.note') }}</p>
            <dl class="landing-visibility-list">
              <div>
                <dt>{{ t('home.story.visibility.items.trend.label') }}</dt>
                <dd>{{ t('home.story.visibility.items.trend.description') }}</dd>
              </div>
              <div>
                <dt>{{ t('home.story.visibility.items.records.label') }}</dt>
                <dd>{{ t('home.story.visibility.items.records.description') }}</dd>
              </div>
              <div>
                <dt>{{ t('home.story.visibility.items.keys.label') }}</dt>
                <dd>{{ t('home.story.visibility.items.keys.description') }}</dd>
              </div>
            </dl>
          </div>
        </section>

        <section class="landing-closing" aria-labelledby="home-closing-title" data-reveal>
          <div>
            <p class="landing-section-kicker">{{ t('home.story.closing.eyebrow') }}</p>
            <h2 id="home-closing-title">{{ t('home.story.closing.title') }}</h2>
            <p>{{ t('home.story.closing.description') }}</p>
          </div>
          <button
            v-if="!isAuthenticated"
            type="button"
            data-testid="home-final-cta"
            class="landing-closing-action"
            @click="openLogin"
          >
            {{ t('home.story.closing.action') }}
            <Icon name="arrowRight" size="sm" aria-hidden="true" />
          </button>
          <RouterLink
            v-else
            to="/data-dashboard"
            data-testid="home-final-cta"
            class="landing-closing-action"
          >
            {{ t('home.story.closing.dashboardAction') }}
            <Icon name="arrowRight" size="sm" aria-hidden="true" />
          </RouterLink>
        </section>
      </div>
    </main>

    <div
      class="home-auth-control-anchor pointer-events-none fixed inset-x-0 z-40 flex justify-center"
      :style="{ bottom: 'calc(env(safe-area-inset-bottom, 0px) + 2rem)' }"
    >
      <div
        class="home-auth-control-shell pointer-events-auto w-[7rem] max-w-[36rem] rounded-[1.2rem] border border-white/55 bg-white/65 p-1 shadow-[0_16px_45px_-18px_rgba(15,23,42,0.18),0_4px_12px_-5px_rgba(15,23,42,0.08)] backdrop-blur-2xl dark:border-white/[0.05] dark:bg-[#181b1f]/70 dark:shadow-[0_18px_48px_-18px_rgba(0,0,0,0.5)] sm:p-1.5"
        :class="{ 'home-auth-control-expanded': showAuthenticatedNav }"
        :style="{ '--user-nav-count': authenticatedNavRef?.itemCount ?? 3 }"
      >
        <div class="home-auth-control-stage relative h-12 sm:h-14">
          <UserBottomNav
            ref="authenticatedNavRef"
            class="home-auth-control-content home-auth-control-user"
            :class="{ 'home-auth-control-content-active': showAuthenticatedNav }"
            active-path="/data-dashboard"
            embedded
          />
          <button
            ref="loginButtonRef"
            type="button"
            data-testid="home-login-button"
            :disabled="authTransition !== 'idle' || !showGuestLoginButton"
            :aria-busy="authTransition !== 'idle'"
            class="home-auth-control-content home-auth-control-guest relative flex h-12 sm:h-14 min-h-[3rem] min-w-0 flex-1 flex-col items-center justify-center gap-0.5 rounded-xl bg-gray-900 px-1 py-1.5 text-[9px] font-medium leading-tight text-white shadow-[0_7px_18px_-7px_rgba(17,24,39,0.58)] ring-1 ring-black/5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-400/70 dark:bg-gray-100 dark:text-gray-900 dark:shadow-[0_7px_18px_-7px_rgba(0,0,0,0.65)] dark:ring-white/30 sm:min-w-[6rem] sm:gap-1 sm:px-3 sm:py-2 sm:text-[10px]"
            :class="{ 'home-auth-control-content-active': showGuestLoginButton }"
            @click="openLogin"
          >
            <span
              v-if="authTransition !== 'idle'"
              class="h-4 w-4 animate-spin rounded-full border-2 border-current border-r-transparent"
              aria-hidden="true"
            />
            <Icon v-else name="login" size="md" :stroke-width="1.8" />
            <span>{{ authTransition !== 'idle' ? t('auth.signingIn') : t('home.loginSystem') }}</span>
          </button>
        </div>
      </div>
    </div>

    <BaseDialog
      :show="loginOpen"
      :title="authMode === 'register' ? t('auth.createAccount') : t('auth.welcomeBack')"
      width="normal"
      content-class="home-auth-dialog-panel"
      :content-style="authDialogStyle"
      :overlay-style="authOverlayStyle"
      :close-on-click-outside="authTransition === 'idle'"
      :close-on-escape="authTransition === 'idle'"
      :show-close-button="authTransition === 'idle'"
      @close="closeLogin"
    >
      <LoginView
        v-if="authMode === 'login'"
        embedded
        :redirect-to="redirectTarget"
        @auth-pending="beginAuthTransition"
        @auth-failure="restoreAuthDialog"
        @auth-step-up="restoreAuthDialog"
        @success="completeLogin"
        @register="openRegister"
      />
      <RegisterView
        v-else
        embedded
        :redirect-to="redirectTarget"
        @success="completeLogin"
        @login="openLogin"
      />
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { markLoginEntry } from '@/utils/dashboard-entry'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import type { CSSProperties } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import UserBottomNav from '@/components/layout/UserBottomNav.vue'
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import TokenConsumptionHero from '@/components/user/dashboard/TokenConsumptionHero.vue'
import type { TrendDataPoint } from '@/types'
import { getPublicActivityTrend } from '@/api/public'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import { toggleThemeWithReveal } from '@/utils/theme-transition'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => t('home.heroDescription'))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isDark = ref(document.documentElement.classList.contains('dark'))
const homePage = ref<HTMLElement | null>(null)
const chartRef = ref<InstanceType<typeof TokenConsumptionHero> | null>(null)
const authenticatedNavRef = ref<InstanceType<typeof UserBottomNav> | null>(null)
const loginButtonRef = ref<HTMLButtonElement | null>(null)
const authTransition = ref<'idle' | 'pending' | 'nav-morph'>(route.query.auth_transition === '1' ? 'pending' : 'idle')
const authDialogStyle = ref<CSSProperties>({})
const authOverlayStyle = ref<CSSProperties>({})
const showAuthenticatedNav = computed(() => isAuthenticated.value && authTransition.value !== 'pending')
const showGuestLoginButton = computed(() => !isAuthenticated.value || authTransition.value === 'pending')
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))
const showModelPlazaEntry = computed(() => modelPlazaEnabled.value && (
  isAuthenticated.value || appStore.cachedPublicSettings?.model_plaza_require_auth !== true
))
const authMode = computed(() => route.query.auth === 'register' ? 'register' : 'login')
const loginOpen = computed(() => route.path === '/login' || route.query.auth === 'login' || route.query.auth === 'register')
const shouldRunOAuthTransition = computed(() => route.query.auth_transition === '1')
const redirectTarget = computed(() => normalizeRedirect(route.query.redirect))
const trendData = ref<TrendDataPoint[]>([])
const trendLoading = ref(true)
const trendError = ref(false)
let trendRequest: Promise<void> | null = null
let oauthTransitionStarted = false
let revealObserver: IntersectionObserver | null = null
let authCollapsePromise: Promise<void> | null = null

const AUTH_COLLAPSE_MS = 460
const AUTH_CONTROL_MORPH_MS = 380

function waitForAnimation(duration: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, duration))
}

function prefersReducedMotion(): boolean {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function beginAuthTransition(): Promise<void> {
  if (authTransition.value !== 'idle') return authCollapsePromise ?? Promise.resolve()

  authTransition.value = 'pending'
  authCollapsePromise = (async () => {
    const reducedMotion = prefersReducedMotion()
    await nextTick()

    const panel = document.querySelector<HTMLElement>('.home-auth-dialog-panel')
    const button = loginButtonRef.value
    if (!panel || !button) {
      authDialogStyle.value = { opacity: 0, pointerEvents: 'none' }
      authOverlayStyle.value = { backgroundColor: 'transparent', backdropFilter: 'none' }
      return
    }

    const panelRect = panel.getBoundingClientRect()
    const buttonRect = button.getBoundingClientRect()
    const duration = reducedMotion ? 1 : AUTH_COLLAPSE_MS
    const dx = buttonRect.left - panelRect.left
    const dy = buttonRect.top - panelRect.top
    const sx = buttonRect.width / panelRect.width
    const sy = buttonRect.height / panelRect.height

    authDialogStyle.value = {
      transformOrigin: 'top left',
      transform: 'none',
      opacity: 1,
      pointerEvents: 'auto',
      transition: 'none',
      willChange: 'transform, opacity'
    }
    authOverlayStyle.value = { transition: 'none' }
    await nextTick()
    // Force the initial state to paint before setting the FLIP destination.
    panel.getBoundingClientRect()
    await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()))

    authDialogStyle.value = {
      transformOrigin: 'top left',
      transform: `translate3d(${dx}px, ${dy}px, 0) scale(${sx}, ${sy})`,
      opacity: 0,
      pointerEvents: 'none',
      transition: `transform ${duration}ms cubic-bezier(0.22, 1, 0.36, 1), opacity ${duration}ms ease`,
      willChange: 'transform, opacity'
    }
    authOverlayStyle.value = {
      backgroundColor: 'transparent',
      backdropFilter: 'none',
      transition: `background-color ${duration}ms ease, backdrop-filter ${duration}ms ease`
    }
    await waitForAnimation(duration)
  })()

  return authCollapsePromise
}

async function restoreAuthDialog(): Promise<void> {
  authTransition.value = 'idle'
  authCollapsePromise = null
  const reducedMotion = prefersReducedMotion()
  const duration = reducedMotion ? 1 : 180
  authDialogStyle.value = {
    transformOrigin: 'top left',
    transform: 'none',
    opacity: 1,
    pointerEvents: 'auto',
    transition: `transform ${duration}ms ease, opacity ${duration}ms ease`
  }
  authOverlayStyle.value = {
    backgroundColor: 'rgba(0, 0, 0, 0.6)',
    backdropFilter: 'blur(4px)',
    transition: `background-color ${duration}ms ease, backdrop-filter ${duration}ms ease`
  }
  await waitForAnimation(duration)
  authDialogStyle.value = {}
  authOverlayStyle.value = {}
}

function normalizeRedirect(value: unknown): string {
  if (typeof value !== 'string' || !value.startsWith('/') || value.startsWith('//') || value.includes('://')) {
    return '/data-dashboard'
  }
  return value
}

function loadTrend(): Promise<void> {
  trendLoading.value = true
  trendError.value = false

  trendRequest = getPublicActivityTrend()
    .then((response) => {
      trendData.value = response.trend.map((point) => ({
        date: point.date,
        requests: 0,
        input_tokens: 0,
        output_tokens: 0,
        cache_creation_tokens: 0,
        cache_read_tokens: 0,
        total_tokens: Math.max(0, Math.min(100, point.activity_percent)),
        cost: 0,
        actual_cost: 0
      }))
    })
    .catch(() => {
      trendError.value = true
    })
    .finally(() => {
      trendLoading.value = false
    })

  return trendRequest
}

function toggleTheme(event: MouseEvent) {
  toggleThemeWithReveal(event, isDark.value, (nextDark) => {
    isDark.value = nextDark
  })
}

function cleanLoginQuery() {
  const query = { ...route.query }
  delete query.auth
  delete query.auth_transition
  delete query.redirect
  return query
}

async function openLogin() {
  if (route.path === '/login' && route.query.auth !== 'register') return
  await router.replace({
    path: '/home',
    query: { ...route.query, auth: 'login', redirect: redirectTarget.value }
  })
}

async function openRegister() {
  await router.replace({
    path: '/home',
    query: { ...route.query, auth: 'register', redirect: redirectTarget.value }
  })
}

async function closeLogin() {
  if (authTransition.value !== 'idle') return
  await router.replace({ path: '/home', query: cleanLoginQuery() })
}

async function completeLogin() {
  if (authTransition.value === 'idle') void beginAuthTransition()
  await authCollapsePromise

  await nextTick()
  await chartRef.value?.playRetractAnimation()

  const reducedMotion = prefersReducedMotion()
  authTransition.value = 'nav-morph'
  await nextTick()
  await waitForAnimation(reducedMotion ? 1 : AUTH_CONTROL_MORPH_MS)
  markLoginEntry()
  await router.replace(redirectTarget.value)
}

async function completeOAuthTransition() {
  if (oauthTransitionStarted) return
  oauthTransitionStarted = true

  const target = redirectTarget.value
  let shouldAnimate = false
  try {
    shouldAnimate = sessionStorage.getItem('sub2api:home-login-transition') === target
    sessionStorage.removeItem('sub2api:home-login-transition')
  } catch {
    shouldAnimate = false
  }

  if (shouldAnimate) {
    authTransition.value = 'pending'
    await nextTick()
    await chartRef.value?.playRetractAnimation()
    authTransition.value = 'nav-morph'
    await nextTick()
    await waitForAnimation(prefersReducedMotion() ? 1 : AUTH_CONTROL_MORPH_MS)
  } else {
    authTransition.value = 'idle'
  }
  markLoginEntry()
  await router.replace(target)
}

watch(
  () => [shouldRunOAuthTransition.value, isAuthenticated.value] as const,
  ([transition, authenticated]) => {
    if (transition && authenticated) void completeOAuthTransition()
  },
  { flush: 'post' }
)

onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }

  authStore.checkAuth()
  if (shouldRunOAuthTransition.value && !isAuthenticated.value) {
    authTransition.value = 'idle'
  }
  void loadTrend()
  if (!appStore.publicSettingsLoaded) void appStore.fetchPublicSettings()
  if (shouldRunOAuthTransition.value && isAuthenticated.value) {
    void completeOAuthTransition()
  }

  const page = homePage.value
  if (page) {
    const revealSections = page.querySelectorAll<HTMLElement>('[data-reveal]')
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (!reducedMotion && 'IntersectionObserver' in window) {
      page.classList.add('landing-motion-ready')
      revealObserver = new IntersectionObserver((entries, observer) => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue
          entry.target.classList.add('is-visible')
          observer.unobserve(entry.target)
        }
      }, { threshold: 0.08, rootMargin: '0px 0px -5% 0px' })
      revealSections.forEach((section) => revealObserver?.observe(section))
    } else {
      revealSections.forEach((section) => section.classList.add('is-visible'))
    }
  }
})

onUnmounted(() => revealObserver?.disconnect())
</script>

<style scoped>
.home-auth-control-shell {
  transition: width 380ms cubic-bezier(0.22, 1, 0.36, 1);
}

.home-auth-control-shell.home-auth-control-expanded {
  width: calc(100% - 2rem);
}

@media (min-width: 640px) {
  .home-auth-control-shell.home-auth-control-expanded {
    /* Same item widths, gaps, padding and border as UserBottomNav. */
    width: calc(var(--user-nav-count) * 6.5rem + 0.25rem + 2px);
  }
}

.home-auth-control-content {
  position: absolute;
  inset: 0;
  width: 100%;
  transform-origin: center;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transform: scale(0.98);
  transition:
    opacity 260ms ease,
    transform 380ms cubic-bezier(0.22, 1, 0.36, 1),
    visibility 0s linear 380ms;
}

.home-auth-control-user {
  z-index: 2;
}

.home-auth-control-guest {
  z-index: 1;
}

.home-auth-control-content-active {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
  transform: scale(1);
  transition:
    opacity 260ms ease,
    transform 380ms cubic-bezier(0.22, 1, 0.36, 1),
    visibility 0s;
}

@media (prefers-reduced-motion: reduce) {
  .home-auth-control-shell,
  .home-auth-control-content,
  .home-auth-control-content-active {
    transition-duration: 1ms;
    transition-delay: 0ms;
  }
}

.landing-page {
  isolation: isolate;
  overflow-x: clip;
}

.landing-hero {
  flex: 0 0 auto;
}

.landing-resources {
  display: flex;
  min-height: 100svh;
  align-items: center;
  padding: clamp(96px, 12vh, 144px) 0 clamp(150px, 18vh, 200px);
  border-top: 1px solid #e1e6eb;
}

.landing-resources-inner {
  width: min(1240px, calc(100% - 2 * clamp(24px, 6vw, 96px)));
  margin: 0 auto;
}

.landing-resources-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 48px;
  margin-bottom: clamp(36px, 5vw, 62px);
}

.landing-resources-heading h2 {
  max-width: 720px;
  margin: 18px 0 0;
  color: #17212d;
  font-size: clamp(34px, 4.6vw, 66px);
  font-weight: 500;
  letter-spacing: -0.075em;
  line-height: 1.06;
}

.dark .landing-resources-heading h2 {
  color: #f3f5f7;
}

.landing-resources-description {
  max-width: 490px;
  margin: 16px 0 0;
  color: #7b8795;
  font-size: 13px;
  line-height: 1.8;
}

.dark .landing-resources-description {
  color: #9aa6b3;
}

.landing-resource-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  border-top: 1px solid #cfd8e1;
  border-bottom: 1px solid #cfd8e1;
}

.dark .landing-resource-grid {
  border-color: rgb(255 255 255 / 0.16);
}

.landing-resource-card {
  position: relative;
  display: flex;
  min-height: 310px;
  flex-direction: column;
  padding: 25px clamp(18px, 2.5vw, 34px) 24px;
  overflow: hidden;
}

.landing-resource-card + .landing-resource-card {
  border-left: 1px solid #dbe2e8;
}

.dark .landing-resource-card + .landing-resource-card {
  border-color: rgb(255 255 255 / 0.13);
}

.landing-resource-card-top {
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  min-height: 70px;
}

.landing-resource-logo {
  width: 32px;
  height: 32px;
}

.landing-resource-logo-openai {
  color: #70a78e;
}

.landing-resource-logo-deepseek {
  color: #6dabc3;
}

.landing-resource-logo-secai {
  color: #7d97b3;
}

.landing-resource-type {
  margin: 0 0 8px;
  color: #708095;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: 9px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.dark .landing-resource-type {
  color: #8d9caf;
}

.landing-resource-card h3 {
  margin: 0;
  color: #1e2935;
  font-size: clamp(29px, 3vw, 40px);
  font-weight: 500;
  letter-spacing: -0.07em;
}

.dark .landing-resource-card h3 {
  color: #eef1f4;
}

.landing-resource-copy {
  max-width: 280px;
  min-height: 48px;
  margin: 11px 0 22px;
  color: #7c8997;
  font-size: 12px;
  line-height: 1.7;
}

.dark .landing-resource-copy {
  color: #9aa6b2;
}

.landing-resource-rule {
  height: 1px;
  margin-top: auto;
  overflow: hidden;
  background: #e1e6eb;
}

.landing-resource-rule span {
  display: block;
  width: 34%;
  height: 100%;
  background: #7d97b3;
}

.landing-resource-card-ds .landing-resource-rule span {
  width: 48%;
}

.landing-resource-card-secai .landing-resource-rule span {
  width: 62%;
}

.dark .landing-resource-rule {
  background: rgb(255 255 255 / 0.12);
}

.dark .landing-resource-rule span {
  background: #718aa6;
}

.landing-resources-footnote {
  margin: 17px 0 0;
  color: #97a2ae;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: 9px;
  letter-spacing: 0.035em;
  line-height: 1.7;
}

.landing-story-shell {
  width: min(1240px, calc(100% - 2 * clamp(24px, 6vw, 96px)));
  margin: 0 auto;
}

.landing-story-section {
  padding: clamp(72px, 9vw, 126px) 0;
  border-top: 1px solid #e1e6eb;
}

.dark .landing-story-section,
.dark .landing-capability-row,
.dark .landing-visibility-list > div {
  border-color: rgb(255 255 255 / 0.12);
}

.landing-section-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 48px;
  margin-bottom: 52px;
}

.landing-section-kicker {
  margin: 0;
  color: #6682a4;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: 10px;
  letter-spacing: 0.14em;
}

.dark .landing-section-kicker {
  color: #93adc9;
}

.landing-section-heading h2,
.landing-closing h2 {
  max-width: 720px;
  margin: 20px 0 0;
  color: #17212d;
  font-size: clamp(32px, 4.2vw, 62px);
  font-weight: 500;
  letter-spacing: -0.075em;
  line-height: 1.08;
}

.dark .landing-section-heading h2,
.dark .landing-closing h2 {
  color: #f3f5f7;
}

.landing-section-description {
  max-width: 300px;
  margin: 0 0 4px;
  color: #8390a0;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: 10px;
  letter-spacing: 0.015em;
  line-height: 1.8;
}

.landing-flow {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: clamp(14px, 2.6vw, 34px);
  padding: 30px 0;
  border-top: 1px solid #cfd8e1;
  border-bottom: 1px solid #cfd8e1;
}

.dark .landing-flow {
  border-color: rgb(255 255 255 / 0.16);
}

.landing-flow-step {
  display: grid;
  grid-template-columns: 30px 24px minmax(0, 1fr);
  align-items: center;
  gap: 13px;
  min-width: 0;
}

.landing-flow-index,
.landing-capability-index {
  color: #99a5b3;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: 10px;
  letter-spacing: 0.08em;
}

.landing-flow-icon,
.landing-capability-icon {
  color: #6682a4;
}

.dark .landing-flow-icon,
.dark .landing-capability-icon {
  color: #9bb6d3;
}

.landing-flow-step strong,
.landing-flow-step span:last-child {
  display: block;
}

.landing-flow-step strong {
  color: #25303b;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.025em;
}

.landing-flow-step div span {
  margin-top: 5px;
  color: #8995a3;
  font-size: 11px;
  line-height: 1.5;
}

.dark .landing-flow-step strong {
  color: #e7ebef;
}

.landing-flow-arrow {
  color: #9ba8b7;
}

.landing-capability-list {
  border-top: 1px solid #cfd8e1;
}

.dark .landing-capability-list {
  border-color: rgb(255 255 255 / 0.16);
}

.landing-capability-row {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr) 24px;
  align-items: center;
  gap: 20px;
  min-height: 106px;
  border-bottom: 1px solid #e1e6eb;
}

.landing-capability-row h3 {
  margin: 0;
  color: #202a35;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: -0.035em;
}

.landing-capability-row p {
  margin: 7px 0 0;
  color: #7f8b99;
  font-size: 12px;
  line-height: 1.7;
}

.dark .landing-capability-row h3 {
  color: #e9edf1;
}

.landing-capability-row .landing-capability-index {
  align-self: start;
  padding-top: 3px;
}

.landing-visibility-layout {
  display: grid;
  grid-template-columns: minmax(220px, 0.78fr) minmax(0, 1.22fr);
  gap: clamp(44px, 9vw, 132px);
}

.landing-visibility-note {
  max-width: 380px;
  margin: 0;
  color: #556473;
  font-size: 15px;
  letter-spacing: -0.025em;
  line-height: 1.8;
}

.dark .landing-visibility-note {
  color: #aab5c1;
}

.landing-visibility-list {
  margin: 0;
  border-top: 1px solid #cfd8e1;
}

.dark .landing-visibility-list {
  border-color: rgb(255 255 255 / 0.16);
}

.landing-visibility-list > div {
  display: grid;
  grid-template-columns: minmax(96px, 0.42fr) minmax(0, 1fr);
  gap: 20px;
  padding: 18px 0;
  border-bottom: 1px solid #e1e6eb;
}

.landing-visibility-list dt {
  color: #526277;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.landing-visibility-list dd {
  margin: 0;
  color: #758292;
  font-size: 12px;
  line-height: 1.7;
}

.dark .landing-visibility-list dt {
  color: #c2ccd5;
}

.dark .landing-visibility-list dd {
  color: #8d9aa8;
}

.landing-closing {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 32px;
  padding: clamp(52px, 7vw, 88px) 0 12px;
  border-top: 1px solid #cfd8e1;
}

.dark .landing-closing {
  border-color: rgb(255 255 255 / 0.16);
}

.landing-closing h2 {
  margin-top: 16px;
  font-size: clamp(30px, 4vw, 52px);
}

.landing-closing > div > p:last-child {
  max-width: 460px;
  margin: 14px 0 0;
  color: #7d8997;
  font-size: 13px;
  line-height: 1.75;
}

.landing-closing-action {
  display: inline-flex;
  min-height: 44px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  gap: 20px;
  padding: 0 17px;
  border: 1px solid #17212d;
  color: #fff;
  background: #17212d;
  font-size: 12px;
  font-weight: 500;
  text-decoration: none;
  transition: color 180ms ease, background-color 180ms ease, border-color 180ms ease;
}

.landing-closing-action:hover {
  color: #17212d;
  background: transparent;
}

.dark .landing-closing-action {
  border-color: #e8edf2;
  color: #111820;
  background: #e8edf2;
}

.dark .landing-closing-action:hover {
  color: #e8edf2;
  background: transparent;
}

.landing-motion-ready [data-reveal] {
  opacity: 0;
  transform: translateY(14px);
  transition: opacity 620ms ease, transform 620ms cubic-bezier(0.2, 0.7, 0.2, 1);
}

.landing-motion-ready [data-reveal].is-visible {
  opacity: 1;
  transform: translateY(0);
}

.landing-header-action {
  display: inline-flex;
  height: 2.5rem;
  width: 2.5rem;
  align-items: center;
  justify-content: center;
  border-radius: 0.85rem;
  color: rgb(107 114 128);
  transition: background-color 160ms ease, color 160ms ease;
}

.landing-header-action:hover {
  background: rgb(15 23 42 / 0.055);
  color: rgb(31 41 55);
}

:global(.dark) .landing-header-action {
  color: rgb(156 163 175);
}

:global(.dark) .landing-header-action:hover {
  background: rgb(255 255 255 / 0.07);
  color: rgb(243 244 246);
}

.landing-chart-stage :deep(section) {
  width: 100%;
}

@media (max-width: 680px) {
  .landing-hero {
    min-height: calc(100svh - 4.5rem);
  }

  .landing-section-heading,
  .landing-closing,
  .landing-resources-heading {
    align-items: flex-start;
    flex-direction: column;
    gap: 18px;
  }

  .landing-resources {
    min-height: auto;
    padding: 88px 0 150px;
  }

  .landing-resources-heading {
    margin-bottom: 30px;
  }

  .landing-resource-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .landing-resource-card {
    min-height: 210px;
    padding: 8px 0;
  }

  .landing-resource-card + .landing-resource-card {
    border-top: 1px solid #dbe2e8;
    border-left: 0;
  }

  .dark .landing-resource-card + .landing-resource-card {
    border-color: rgb(255 255 255 / 0.13);
  }

  .landing-resource-card-top {
    min-height: 25px;
  }

  .landing-resource-copy {
    min-height: 0;
    margin: 3px 0 7px;
  }

  .landing-section-heading {
    margin-bottom: 34px;
  }

  .landing-section-description {
    max-width: 360px;
  }

  .landing-flow {
    grid-template-columns: minmax(0, 1fr);
    gap: 0;
    padding: 8px 0;
  }

  .landing-flow-step {
    min-height: 74px;
  }

  .landing-flow-arrow {
    display: none;
  }

  .landing-flow-step + .landing-flow-step {
    border-top: 1px solid #e1e6eb;
  }

  .dark .landing-flow-step + .landing-flow-step {
    border-color: rgb(255 255 255 / 0.12);
  }

  .landing-visibility-layout {
    grid-template-columns: minmax(0, 1fr);
    gap: 28px;
  }

  .landing-closing-action {
    min-width: 200px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .landing-header-action {
    transition: none;
  }

  .landing-closing-action,
  .landing-motion-ready [data-reveal] {
    transition: none;
    opacity: 1;
    transform: none;
  }
}

.landing-page > header,
.landing-hero > div:not(.landing-chart-stage),
.landing-resources,
.landing-story-shell {
  transition: opacity 600ms ease, transform 600ms cubic-bezier(0.22, 1, 0.36, 1);
}

.home-auth-success-fade > header,
.home-auth-success-fade .landing-hero > div:not(.landing-chart-stage),
.home-auth-success-fade .landing-resources,
.home-auth-success-fade .landing-story-shell {
  opacity: 0;
  transform: translateY(8px);
  pointer-events: none;
}

@media (prefers-reduced-motion: reduce) {
  .home-auth-success-fade > header,
  .home-auth-success-fade .landing-hero > div:not(.landing-chart-stage),
  .home-auth-success-fade .landing-resources,
  .home-auth-success-fade .landing-story-shell {
    transition-duration: 1ms;
  }

  .landing-page > header,
  .landing-hero > div:not(.landing-chart-stage),
  .landing-resources,
  .landing-story-shell {
    transition-duration: 1ms;
  }
}
</style>
