import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import OfficialTokenPrice from '../OfficialTokenPrice.vue'
import type { PlazaTimePricing } from '@/api/modelPlaza'
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (key: string) => key }) }))
const schedule: PlazaTimePricing = {
  timezone: 'Asia/Shanghai', weekdays_only: true,
  periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: 2 },
    { start_time: '14:00:00', end_time: '18:00:00', multiplier: 2 }]
}
describe('official peak prices', () => {
  it('uses official multiplier, not channel markup, and normalizes clock precision', () => {
    const channelSchedule = { ...schedule, periods: schedule.periods.map(p => ({ ...p, start_time: p.start_time.slice(0,5), end_time: p.end_time.slice(0,5), multiplier: 3 })) }
    const wrapper = mount(OfficialTokenPrice, { props: { value: 0.15e-6, schedule, channelSchedule, period: channelSchedule.periods[0] } })
    expect(wrapper.text()).toBe('$0.30')
  })
  it('keeps base and non time-priced model references unchanged', () => {
    expect(mount(OfficialTokenPrice, { props: { value: 0.15e-6, schedule, channelSchedule: schedule, period: null } }).text()).toBe('$0.15')
    expect(mount(OfficialTokenPrice, { props: { value: 3e-6, channelSchedule: schedule, period: schedule.periods[0] } }).text()).toBe('$3.00')
  })
  it.each([undefined, { ...schedule, timezone: 'UTC' }, { ...schedule, weekdays_only: false }])('labels independent official windows for unmatched schedules', (channelSchedule) => {
    const wrapper = mount(OfficialTokenPrice, { props: { value: 0.003e-6, schedule, channelSchedule, period: null } })
    expect(wrapper.text()).toContain('$0.003')
    expect(wrapper.text()).toContain('$0.006')
    expect(wrapper.text()).toContain('09:00–12:00')
    expect(wrapper.text()).toContain('Asia/Shanghai')
  })
})
