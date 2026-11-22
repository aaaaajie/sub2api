import type { PlazaTimePricing, PlazaTimePricingPeriod } from '@/api/modelPlaza'

export const clock = (value: string) => value.replace(/^(\d{2}:\d{2}):00$/, '$1')

const windowKey = (period: PlazaTimePricingPeriod) =>
  `${clock(period.start_time)}-${clock(period.end_time)}`

// Compare prices only when both schedules cover the same periods.
export function officialTimeMultiplier(
  official: PlazaTimePricing | undefined,
  channel: PlazaTimePricing | undefined,
  period: PlazaTimePricingPeriod | null
): number | undefined {
  if (!official?.periods.length && !channel?.periods.length) return 1
  if (!official || !channel || official.timezone !== channel.timezone ||
      !!official.weekdays_only !== !!channel.weekdays_only ||
      official.periods.length !== channel.periods.length ||
      !official.periods.every(p => channel.periods.some(q => windowKey(p) === windowKey(q)))) {
    return undefined
  }
  return period ? official.periods.find(p => windowKey(p) === windowKey(period))?.multiplier : 1
}
