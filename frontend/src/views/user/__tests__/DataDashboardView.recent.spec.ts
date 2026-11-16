import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { enableAutoUnmount, flushPromises, shallowMount } from '@vue/test-utils'
import DataDashboardView from '../DataDashboardView.vue'
import TokenConsumptionHero from '@/components/user/dashboard/TokenConsumptionHero.vue'
const api = vi.hoisted(() => ({ getDashboardRecentHours: vi.fn(), getDashboardTrend: vi.fn() }))
vi.mock('@/api/usage', () => ({ usageAPI: api }))
vi.mock('vue-i18n', async (importOriginal) => ({ ...await importOriginal<typeof import('vue-i18n')>(), useI18n: () => ({ t: (key: string) => key }) }))
enableAutoUnmount(afterEach)
beforeEach(() => { vi.clearAllMocks(); api.getDashboardRecentHours.mockResolvedValue({trend:[]}); api.getDashboardTrend.mockResolvedValue({trend:[]}) })
async function open() {
 const wrapper = shallowMount(DataDashboardView, {global:{stubs:{AppLayout:{template:'<div><slot /></div>'}}}})
 await flushPromises()
 return wrapper
}
describe('recent active hours', () => {
 it('uses historical active hours without a date cutoff, and retains zero-token requests', async () => {
  api.getDashboardRecentHours.mockResolvedValue({trend:[{date:'2000-01-01 10:00',requests:1,total_tokens:0},{date:'1999-12-01 08:00',requests:5,total_tokens:100}]})
  const w=await open(); const hero=w.getComponent(TokenConsumptionHero)
  expect(api.getDashboardTrend).not.toHaveBeenCalled()
  expect(api.getDashboardRecentHours).toHaveBeenCalledWith(Intl.DateTimeFormat().resolvedOptions().timeZone)
  expect(hero.props('trendData').map(p=>p.date)).toEqual(['1999-12-01 08:00','2000-01-01 10:00'])
  hero.vm.$emit('update:range','7d'); await flushPromises()
  expect(api.getDashboardTrend).toHaveBeenCalledWith(expect.objectContaining({granularity:'day'}))
 })
 it('keeps the empty zero baseline only for successful empty results', async () => {
  const w=await open(); const hero=w.getComponent(TokenConsumptionHero)
  expect(hero.props('trendData')).toHaveLength(24)
  expect(hero.props('trendData').every(p=>p.requests===0)).toBe(true)
 })
 it('shows failure and supports retry without turning errors into zero records', async () => {
  api.getDashboardRecentHours.mockRejectedValueOnce(new Error('offline'))
  const w=await open(); const hero=w.getComponent(TokenConsumptionHero)
  expect(hero.props('error')).toBe(true)
  expect(hero.props('trendData')).toEqual([])
  hero.vm.$emit('retry'); await flushPromises()
  expect(hero.props('error')).toBe(false)
  expect(hero.props('trendData')).toHaveLength(24)
 })
})
