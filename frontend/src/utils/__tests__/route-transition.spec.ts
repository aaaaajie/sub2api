import { afterEach, describe, expect, it, vi } from 'vitest'
import { createMemoryHistory, createRouter } from 'vue-router'
import { installRouteTransitions } from '../route-transition'

async function setup(reduced = false) {
  vi.stubGlobal('matchMedia', () => ({ matches: reduced }))
  const start = vi.fn((update: () => Promise<void>) => {
    const finished = Promise.resolve().then(update)
    return { finished, skipTransition: vi.fn() }
  })
  Object.defineProperty(document, 'startViewTransition', { configurable: true, value: start })
  const router = createRouter({ history: createMemoryHistory(), routes: ['/home', '/keys', '/profile', '/purchase'].map(path => ({ path, component: { template: '<div />' } })) })
  await router.push('/keys')
  return { router, start, cleanup: installRouteTransitions(router) }
}
afterEach(() => { vi.unstubAllGlobals(); Reflect.deleteProperty(document, 'startViewTransition'); document.documentElement.classList.remove('route-transition') })

describe('route transitions', () => {
  it('commits navigation and completes the captured update without blocking', async () => {
    const { router, start, cleanup } = await setup()
    try {
      await router.push('/profile')
      await Promise.resolve()
      expect(router.currentRoute.value.path).toBe('/profile')
      expect(start).toHaveBeenCalledTimes(1)
      await router.push('/purchase')
      expect(router.currentRoute.value.path).toBe('/purchase')
    } finally { cleanup() }
    expect(document.documentElement.classList.contains('route-transition')).toBe(false)
  })
  it('skips reduced-motion, query-only and public-page navigation', async () => {
    const { router, start, cleanup } = await setup(true)
    try {
      await router.push('/profile')
      expect(start).not.toHaveBeenCalled()
      vi.stubGlobal('matchMedia', () => ({ matches: false }))
      await router.push('/profile?tab=security')
      await router.push('/home')
      expect(start).not.toHaveBeenCalled()
    } finally { cleanup() }
  })
  it('navigates normally without browser transition support', async () => {
    const { router, start, cleanup } = await setup()
    Reflect.deleteProperty(document, 'startViewTransition')
    try {
      await router.push('/profile')
      expect(router.currentRoute.value.path).toBe('/profile')
      expect(start).not.toHaveBeenCalled()
    } finally { cleanup() }
  })
})
