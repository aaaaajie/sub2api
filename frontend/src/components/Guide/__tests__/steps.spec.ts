import { describe, expect, it } from 'vitest'
import { getUserSteps } from '../steps'
const t = (key: string) => key

describe('user quick start', () => {
  it('walks through creation, group choice, import confirmation and documentation', () => {
    const steps = getUserSteps(t, { showImport: true, showDocs: true })
    expect(steps.map(step => step.element)).toEqual([
      undefined, '[data-tour="user-my-keys"]', '[data-tour="keys-create-btn"]',
      '[data-tour="key-form-name"]', '[data-tour="key-form-provider"]',
      '[data-tour="key-form-group"]', '[data-tour="key-form-submit"]',
      '[data-tour="key-import-ccs"]', '[data-tour="ccs-clients"]',
      '[data-tour="ccs-confirm"]', '[data-tour="header-docs"]'
    ])
    expect(steps[6].popover?.showButtons).toEqual(['close'])
    expect(steps[9].popover?.showButtons).toEqual(['close'])
  })
  it('omits entries disabled or not configured by the site', () => {
    const steps = getUserSteps(t, { showImport: false, showDocs: false })
    expect(steps).toHaveLength(7)
    expect(steps.at(-1)?.element).toBe('[data-tour="key-form-submit"]')
  })
})
