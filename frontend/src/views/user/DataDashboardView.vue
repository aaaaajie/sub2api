<template>
  <AppLayout>
    <section class="-mx-4 w-[calc(100%+2rem)] pb-8 md:-mx-6 md:w-[calc(100%+3rem)] lg:-mx-8 lg:w-[calc(100%+4rem)]" aria-labelledby="data-dashboard-title">
      <h2 id="data-dashboard-title" class="sr-only">{{ t('nav.dataDashboard') }}</h2>
      <div class="dashboard-first-screen">
      <TokenConsumptionHero
        :trend-data="trendData"
        :range="selectedRange"
        :granularity="granularity"
        :loading="loading"
        :error="loadError"
        @retry="loadTrend"
        @update:range="changeRange"
      />
      </div>
      <DataDashboardInsights />
    </section>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TokenConsumptionHero from '@/components/user/dashboard/TokenConsumptionHero.vue'
import DataDashboardInsights from '@/components/user/dashboard/DataDashboardInsights.vue'
import { usageAPI } from '@/api/usage'
import type { TrendDataPoint } from '@/types'
import { formatDateLocalInput } from '@/utils/format'

const { t } = useI18n()
const trendData = ref<TrendDataPoint[]>([])
const loading = ref(true)
const loadError = ref(false)
type DashboardRange = 'recent' | '7d' | '30d'
const selectedRange = ref<DashboardRange>('recent')
const granularity = ref<'day' | 'hour'>('hour')
let loadSequence = 0

function getDateRange(days: 7 | 30) {
  const end = new Date()
  const localEnd = new Date(end.getFullYear(), end.getMonth(), end.getDate())
  const start = new Date(localEnd)
  start.setDate(start.getDate() - (days - 1))

  return {
    startDate: formatDateLocalInput(start),
    endDate: formatDateLocalInput(localEnd)
  }
}

function createZeroTrendPoint(date: Date, pointGranularity: 'day' | 'hour'): TrendDataPoint {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const datePart = `${year}-${month}-${day}`
  const dateValue = pointGranularity === 'hour'
    ? `${datePart} ${String(date.getHours()).padStart(2, '0')}:00`
    : datePart

  return {
    date: dateValue,
    requests: 0,
    input_tokens: 0,
    output_tokens: 0,
    cache_creation_tokens: 0,
    cache_read_tokens: 0,
    total_tokens: 0,
    cost: 0,
    actual_cost: 0
  }
}

function createEmptyTrend(range: DashboardRange, pointGranularity: 'day' | 'hour'): TrendDataPoint[] {
  if (pointGranularity === 'hour') {
    const currentHour = new Date()
    currentHour.setMinutes(0, 0, 0)
    return Array.from({ length: 24 }, (_, index) => {
      const date = new Date(currentHour.getTime() - (23 - index) * 60 * 60 * 1000)
      return createZeroTrendPoint(date, pointGranularity)
    })
  }

  const days = range === '30d' ? 30 : 7
  const today = new Date()
  const end = new Date(today.getFullYear(), today.getMonth(), today.getDate())
  const start = new Date(end)
  start.setDate(start.getDate() - (days - 1))

  return Array.from({ length: days }, (_, index) => {
    const date = new Date(start)
    date.setDate(start.getDate() + index)
    return createZeroTrendPoint(date, pointGranularity)
  })
}

async function loadTrend(range: DashboardRange = selectedRange.value) {
  const sequence = ++loadSequence
  loading.value = true
  loadError.value = false

  try {
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
    let points: TrendDataPoint[]
    if (range === 'recent') {
      granularity.value = 'hour'
      const response = await usageAPI.getDashboardRecentHours(timezone)
      points = response.trend ?? []
    } else {
      const days = range === '7d' ? 7 : 30
      const { startDate, endDate } = getDateRange(days)
      granularity.value = 'day'
      const response = await usageAPI.getDashboardTrend({
        start_date: startDate,
        end_date: endDate,
        granularity: 'day',
        timezone
      })
      points = response.trend ?? []
    }
    if (sequence !== loadSequence) return
    const sortedPoints = points.sort((left, right) => left.date.localeCompare(right.date))
    trendData.value = sortedPoints.length > 0
      ? sortedPoints
      : createEmptyTrend(range, granularity.value)
  } catch {
    if (sequence === loadSequence) loadError.value = true
  } finally {
    if (sequence === loadSequence) loading.value = false
  }
}

function changeRange(range: DashboardRange) {
  if (range === selectedRange.value) return
  selectedRange.value = range
  void loadTrend(range)
}

onMounted(loadTrend)
</script>

<style scoped>
.dashboard-first-screen {
  min-height: calc(100svh - 5rem);
  padding-bottom: 5rem;
}
@media (min-width: 768px) {
  .dashboard-first-screen { min-height: calc(100svh - 5.5rem); }
}
@media (min-width: 1024px) {
  .dashboard-first-screen { min-height: calc(100svh - 6rem); }
}
</style>
