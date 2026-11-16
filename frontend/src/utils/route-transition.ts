import { nextTick } from 'vue'
import type { Router } from 'vue-router'

const pages = new Set(['/data-dashboard', '/keys', '/profile', '/purchase', '/orders', '/usage', '/dashboard', '/redeem', '/affiliate', '/subscriptions'])
type Transition = { finished: Promise<void>; skipTransition: () => void }
type TransitionDocument = Document & {
  startViewTransition?: (update: () => Promise<void>) => Transition
}

// Capture the old page before committing navigation, then wait for Vue's new DOM.
export function installRouteTransitions(router: Router): () => void {
  let active: Transition | undefined
  let finishUpdate: (() => void) | undefined
  let sequence = 0
  const root = document.documentElement
  const finish = () => { finishUpdate?.(); finishUpdate = undefined }
  const removeAfter = router.afterEach(async () => { await nextTick(); finish() })
  const removeError = router.onError(() => { finish(); active?.skipTransition() })
  const removeBefore = router.beforeResolve((to, from) => {
    finish()
    active?.skipTransition()
    const doc = document as TransitionDocument
    if (document.body.classList.contains('driver-active') || !doc.startViewTransition || to.path === from.path || !pages.has(to.path) || !pages.has(from.path)
      || window.matchMedia('(prefers-reduced-motion: reduce)').matches
      || root.classList.contains('theme-transition-to-dark') || root.classList.contains('theme-transition-to-light')) return

    const current = ++sequence
    root.classList.add('route-transition')
    return new Promise<void>((proceed) => {
      try {
        active = doc.startViewTransition!(() => new Promise<void>((resolve) => {
          finishUpdate = resolve
          proceed()
        }))
        const cleanup = () => {
          proceed()
          if (current !== sequence) return
          root.classList.remove('route-transition')
          active = undefined
        }
        void active.finished.then(cleanup, cleanup)
      } catch {
        root.classList.remove('route-transition')
        proceed()
      }
    })
  })
  return () => {
    ++sequence
    finish()
    active?.skipTransition()
    root.classList.remove('route-transition')
    removeBefore(); removeAfter(); removeError()
  }
}
