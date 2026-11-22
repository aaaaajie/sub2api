import { syncThemeColor } from './theme-meta'

type ViewTransitionHandle = {
  finished: Promise<void>
}

type ViewTransitionDocument = Document & {
  startViewTransition?: (update: () => void) => ViewTransitionHandle
}

let transitionSequence = 0

export function toggleThemeWithReveal(
  event: MouseEvent,
  currentDark: boolean,
  updateLocalState: (isDark: boolean) => void
) {
  const root = document.documentElement
  const nextDark = !currentDark
  let applied = false
  const applyTheme = () => {
    if (applied) return
    applied = true
    root.classList.toggle('dark', nextDark)
    localStorage.setItem('theme', nextDark ? 'dark' : 'light')
    syncThemeColor(nextDark)
    updateLocalState(nextDark)
  }

  const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  const transitionDocument = document as ViewTransitionDocument
  const startViewTransition = transitionDocument.startViewTransition
  if (reducedMotion || !startViewTransition) {
    applyTheme()
    return
  }

  const trigger = event.currentTarget instanceof HTMLElement ? event.currentTarget : null
  const rect = trigger?.getBoundingClientRect()
  root.style.setProperty('--theme-wipe-x', `${rect ? rect.left + rect.width / 2 : window.innerWidth / 2}px`)
  root.style.setProperty('--theme-wipe-y', `${rect ? rect.top + rect.height / 2 : window.innerHeight / 2}px`)

  const sequence = ++transitionSequence
  const transitionClass = nextDark ? 'theme-transition-to-dark' : 'theme-transition-to-light'
  root.classList.remove('theme-transition-to-dark', 'theme-transition-to-light')
  root.classList.add(transitionClass)

  const cleanup = () => {
    if (sequence !== transitionSequence) return
    root.classList.remove('theme-transition-to-dark', 'theme-transition-to-light')
    root.style.removeProperty('--theme-wipe-x')
    root.style.removeProperty('--theme-wipe-y')
  }

  try {
    const transition = startViewTransition.call(transitionDocument, applyTheme)
    void transition.finished.then(cleanup, cleanup)
  } catch {
    cleanup()
    applyTheme()
  }
}
