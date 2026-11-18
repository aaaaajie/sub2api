<template>
<section class="w-full" :aria-label="presentation === 'landing' ? t('home.relativeTrendAccessible') : t('nav.dataDashboard')">
    <header v-if="presentation === 'dashboard'" class="flex flex-col gap-1 px-4 pb-1 md:px-6 lg:px-8">
      <div class="flex min-h-[3.75rem] items-center justify-between gap-4">
        <Transition name="token-summary" appear>
          <div
            v-if="summaryVisible || error || (!loading && !trendData.length)"
            class="flex min-w-0 items-baseline gap-2.5 tabular-nums"
          >
            <p v-if="summaryVisible" class="flex min-w-0 items-baseline gap-2">
              <span class="token-summary-value truncate text-3xl font-semibold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
                <span
                  v-for="candidate in summaryWidthCandidates"
                  :key="candidate"
                  class="token-summary-value-reserve"
                  aria-hidden="true"
                >{{ candidate }}</span>
                <span class="token-summary-reveal token-summary-reveal--value token-summary-value-current">
                  {{ formatNumber(Math.round(displayedTokens)) }}
                </span>
              </span>
              <span class="token-summary-reveal token-summary-reveal--unit text-xs font-medium text-gray-400 dark:text-gray-500">
                {{ t('dashboard.tokenConsumption.tokens') }}
              </span>
            </p>
            <Select
              class="token-period-select token-summary-reveal token-summary-reveal--period shrink-0"
              :model-value="range"
              :options="rangeOptions"
              :aria-label="t('dashboard.tokenConsumption.rangeLabel')"
              :dropdown-class="isDarkMode
                ? 'token-range-dropdown token-range-dropdown-dark'
                : 'token-range-dropdown token-range-dropdown-light'"
              @update:model-value="handleRangeChange"
            />
          </div>
        </Transition>

        <Transition name="token-detail-toggle" appear>
          <TrendDetailToggle
            v-if="presentation === 'dashboard'"
            v-model="detailsEnabled"
            class="ml-auto"
          />
        </Transition>
      </div>

      <p v-if="range === 'recent'" class="min-h-[15px] text-[10px] text-gray-400 dark:text-gray-500">
        <span v-if="summaryVisible && !loading && !error && trendData.some(point => point.requests > 0)" class="token-summary-reveal token-summary-reveal--recent-hours inline-block">
          {{ t('dashboard.recentActiveHours', { count: trendData.length }) }}
          <span class="ml-2">{{ trendData[0]?.date.slice(0, 10) }} — {{ trendData[trendData.length - 1]?.date.slice(0, 10) }}</span>
        </span>
      </p>
    </header>

    <div
      class="relative"
      :class="presentation === 'landing'
        ? 'mt-0 h-[clamp(14rem,42svh,34rem)] sm:h-[clamp(20rem,52svh,43rem)]'
        : 'mt-2 h-[clamp(18rem,48svh,28rem)] sm:mt-3 sm:h-[clamp(24rem,55svh,38rem)]'"
    >
      <div
        v-if="loading"
        class="flex h-full items-center justify-center"
        role="status"
        :aria-label="t('dashboard.tokenConsumption.loading')"
      >
        <span class="text-xs text-gray-400 dark:text-gray-500">
          {{ t('dashboard.tokenConsumption.loading') }}
        </span>
      </div>

      <div v-else-if="error" class="flex h-full flex-col items-center justify-center gap-3 text-center" role="alert">
        <p class="text-sm text-gray-500 dark:text-gray-400">
          {{ t('dashboard.tokenConsumption.error') }}
        </p>
        <button
          type="button"
          class="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 dark:border-white/10 dark:text-gray-300 dark:hover:bg-white/[0.05]"
          @click="emit('retry')"
        >
          {{ t('dashboard.tokenConsumption.retry') }}
        </button>
      </div>

      <div
        v-else-if="!trendData.length"
        class="relative h-full w-full"
        role="img"
        :aria-label="presentation === 'landing'
          ? t('home.relativeTrendAccessible')
          : `${t('dashboard.noUsageRecords')}, 0 ${t('dashboard.tokenConsumption.tokens')}`"
      >
        <p class="absolute inset-x-0 top-1/2 -translate-y-1/2 text-center text-sm text-gray-400 dark:text-gray-500">
          {{ t('dashboard.noUsageRecords') }}
        </p>
        <div class="absolute bottom-[30px] left-10 right-3 border-t border-gray-200/80 dark:border-white/10" aria-hidden="true" />
        <span class="absolute bottom-[23px] left-3 -translate-y-1/2 text-[10px] leading-none text-gray-400 dark:text-gray-500" aria-hidden="true">
          0
        </span>
      </div>

      <div v-else class="relative h-full w-full">
        <div
          v-if="presentation === 'landing'"
          class="pointer-events-none absolute right-4 top-2 z-10 rounded-full border border-gray-200/80 bg-white/70 px-3 py-1 text-[10px] font-medium tracking-wide text-gray-500 backdrop-blur-md dark:border-white/[0.08] dark:bg-[#141719]/70 dark:text-gray-400 sm:right-8"
        >
          {{ t('home.relativeActivity') }}
        </div>
        <div
          class="token-detail-legend pointer-events-none absolute inset-x-4 top-1 z-10 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-[10px] font-medium text-gray-500 dark:text-gray-400 md:inset-x-6 lg:inset-x-8"
          :class="detailsEnabled ? 'opacity-100' : 'opacity-0'"
          :aria-label="t('dashboard.tokenConsumption.detailLegend')"
          :aria-hidden="!detailsEnabled"
          role="list"
        >
          <span
            v-for="item in legendItems"
            :key="item.label"
            class="inline-flex items-center gap-1.5 whitespace-nowrap"
            role="listitem"
          >
            <span
              class="w-3 shrink-0 border-t-2"
              :class="item.dashed ? 'border-dashed' : 'border-solid'"
              :style="{ borderColor: item.color }"
              aria-hidden="true"
            />
            {{ item.label }}
          </span>
        </div>

        <div
          class="h-full w-full"
          role="img"
          :aria-label="presentation === 'landing'
            ? t('home.relativeTrendAccessible')
            : t('dashboard.tokenConsumption.chartAccessible', {
              period: rangeOptions.find((option) => option.value === range)?.label ?? '',
              total: formatNumber(totalTokens)
            })"
        >
          <Line
            ref="chartRef"
            :data="chartData"
            :options="chartOptions"
            :plugins="[tokenLineRevealPlugin]"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  LineElement,
  LinearScale,
  PointElement,
  Tooltip,
  type Chart,
  type ChartDataset,
  type ChartOptions,
  type Plugin,
  type ScriptableContext
} from 'chart.js'
import { Line } from 'vue-chartjs'
import type { ChartComponentRef } from 'vue-chartjs'
import { useI18n } from 'vue-i18n'
import Select from '@/components/common/Select.vue'
import TrendDetailToggle from './TrendDetailToggle.vue'
import { useAuthStore } from '@/stores/auth'
import type { TrendDataPoint } from '@/types'
import { formatNumber } from '@/utils/format'

ChartJS.register(CategoryScale, LineElement, LinearScale, PointElement, Tooltip)

const props = withDefaults(defineProps<{
  trendData: TrendDataPoint[]
  range: 'recent' | '7d' | '30d'
  presentation?: 'dashboard' | 'landing'
  granularity?: 'day' | 'hour'
  loading?: boolean
  error?: boolean
  autoReveal?: boolean
  initiallyRevealed?: boolean
}>(), {
  presentation: 'dashboard',
  granularity: 'day',
  loading: false,
  error: false,
  autoReveal: true,
  initiallyRevealed: false
})

const emit = defineEmits<{
  (event: 'retry'): void
  (event: 'update:range', value: 'recent' | '7d' | '30d'): void
}>()

const { t, locale } = useI18n()
const authStore = useAuthStore()
const chartRef = ref<ChartComponentRef | null>(null)
const isDarkMode = ref(false)
const prefersReducedMotion = ref(false)
const summaryVisible = ref(false)
const displayedTokens = ref(0)
const detailsEnabled = ref(false)
const detailsRendered = ref(false)
const animationDuration = 1300
const countUpDuration = 800
const detailFadeDuration = 240
const detailsStorageKey = computed(() => authStore.user
  ? `sub2api:data-dashboard-detail:${authStore.user.id}`
  : '')
let syncingDetailsPreference = false
let detailAnimationFrame = 0
let detailAnimationSequence = 0

watch(detailsStorageKey, (key) => {
  if (detailAnimationFrame) cancelAnimationFrame(detailAnimationFrame)
  detailAnimationFrame = 0
  detailAnimationSequence += 1
  syncingDetailsPreference = true
  let enabled = false
  if (key && props.presentation === 'dashboard') {
    try {
      enabled = localStorage.getItem(key) === 'true'
    } catch {
      enabled = false
    }
  }
  detailsEnabled.value = enabled
  detailsRendered.value = enabled
  syncingDetailsPreference = false

  void nextTick(() => {
    const chart = chartRef.value?.chart as RevealableChart | undefined
    if (!chart) return
    chart.$detailOpacity = enabled ? 1 : 0
    chart.$detailRevealProgress = 1
    chart.$detailRevealActive = false
    chart.draw()
  })
}, { immediate: true, flush: 'sync' })

watch(detailsEnabled, (enabled) => {
  if (syncingDetailsPreference) return
  const key = detailsStorageKey.value
  if (key) {
    try {
      localStorage.setItem(key, String(enabled))
    } catch {
      // Keep the local toggle usable when browser storage is unavailable.
    }
  }
  void animateDetailVisibility(enabled)
}, { flush: 'sync' })

const rangeOptions = computed(() => [
  { value: 'recent', label: t('dashboard.recentUsage') },
  { value: '7d', label: t('dates.last7Days') },
  { value: '30d', label: t('dates.last30Days') }
])

function handleRangeChange(value: string | number | boolean | null) {
  if (value === 'recent' || value === '7d' || value === '30d') {
    emit('update:range', value)
  }
}

type RevealableChart = Chart<'line'> & {
  $tokenRevealProgress?: number
  $tokenRevealClipActive?: boolean
  $tokenEndpointPulse?: number
  $detailOpacity?: number
  $detailRevealProgress?: number
  $detailRevealActive?: boolean
  $detailDatasetDrawActive?: boolean
}

const colors = computed(() => isDarkMode.value
  ? {
      line: '#e8f7f8',
      glow: 'rgba(185, 243, 243, 0.42)',
      glowFade: 'rgba(185, 243, 243, 0)',
      lineGlowBlur: 16,
      cursorGlowBlur: 18,
      pointBorder: '#090b0d',
      valueLabel: '#eef8f8',
      dateLabel: '#737d86',
      tooltip: '#111417',
      tooltipBorder: 'rgba(255, 255, 255, 0.1)',
      tooltipText: '#f3f4f6'
    }
  : {
      line: '#171b20',
      glow: 'rgba(23, 27, 32, 0.22)',
      glowFade: 'rgba(23, 27, 32, 0)',
      lineGlowBlur: 15,
      cursorGlowBlur: 18,
      pointBorder: '#ffffff',
      valueLabel: '#171b20',
      dateLabel: '#748597',
      tooltip: '#ffffff',
      tooltipBorder: 'rgba(15, 23, 42, 0.12)',
      tooltipText: '#111827'
    })

const detailColors = {
  input: '#3b82f6',
  output: '#10b981',
  cacheCreation: '#f59e0b',
  cacheRead: '#06b6d4',
  cacheHitRate: '#8b5cf6'
}

const aggregateCacheHitRate = computed(() => {
  const inputTokens = props.trendData.reduce((sum, point) => sum + point.input_tokens, 0)
  const cacheCreationTokens = props.trendData.reduce((sum, point) => sum + point.cache_creation_tokens, 0)
  const cacheReadTokens = props.trendData.reduce((sum, point) => sum + point.cache_read_tokens, 0)
  const promptTokens = inputTokens + cacheCreationTokens + cacheReadTokens
  return promptTokens > 0 ? (cacheReadTokens / promptTokens) * 100 : 0
})

const formattedCacheHitRate = computed(() => new Intl.NumberFormat(locale.value, {
  minimumFractionDigits: 1,
  maximumFractionDigits: 1
}).format(aggregateCacheHitRate.value))

const legendItems = computed(() => [
  { label: t('dashboard.tokenConsumption.totalTokens'), color: colors.value.line, dashed: false },
  { label: t('dashboard.input'), color: detailColors.input, dashed: false },
  { label: t('dashboard.output'), color: detailColors.output, dashed: false },
  { label: t('dashboard.tokenConsumption.cacheCreation'), color: detailColors.cacheCreation, dashed: false },
  { label: t('dashboard.tokenConsumption.cacheRead'), color: detailColors.cacheRead, dashed: false },
  {
    label: `${t('dashboard.tokenConsumption.cacheHitRate')} ${formattedCacheHitRate.value}%`,
    color: detailColors.cacheHitRate,
    dashed: true
  }
])

const totalTokens = computed(() => props.trendData.reduce(
  (total, point) => total + (Number.isFinite(point.total_tokens) ? point.total_tokens : 0),
  0
))

const summaryWidthCandidates = computed(() => {
  const total = Math.floor(totalTokens.value)
  const values = [total]
  const compactBoundaries = [1_000, 10_000, 1_000_000, 100_000_000, 1_000_000_000, 1_000_000_000_000]

  for (const boundary of compactBoundaries) {
    for (const factor of [0.9, 0.99, 0.999]) {
      const candidate = Math.floor(boundary * factor)
      if (candidate <= total) values.push(candidate)
    }
  }

  return [...new Set(values.map((value) => formatNumber(value)))]
})

const chartData = computed(() => {
  const datasets: ChartDataset<'line', number[]>[] = [{
    label: t('dashboard.tokenConsumption.totalTokens'),
    data: props.trendData.map((point) => point.total_tokens),
    borderColor: colors.value.line,
    backgroundColor: colors.value.line,
    borderWidth: 2.2,
    borderCapStyle: 'round' as const,
    borderJoinStyle: 'round' as const,
    fill: false,
    tension: 0.36,
    cubicInterpolationMode: 'monotone' as const,
    pointRadius: (context: ScriptableContext<'line'>) => Number(context.raw) > 0 ? 3.5 : 0,
    pointHoverRadius: (context: ScriptableContext<'line'>) => Number(context.raw) > 0 ? 6 : 0,
    pointBackgroundColor: colors.value.line,
    pointBorderColor: colors.value.pointBorder,
    pointBorderWidth: 2,
    pointHitRadius: (context: ScriptableContext<'line'>) => Number(context.raw) > 0 ? 12 : 0,
    clip: { top: 8, right: 12, bottom: 8, left: 12 }
  }]

  if (detailsRendered.value) {
    const detailDatasets: Array<{
      label: string
      data: number[]
      color: string
      dashed?: boolean
      yAxisID?: string
    }> = [
      {
        label: t('dashboard.input'),
        data: props.trendData.map((point) => point.input_tokens),
        color: detailColors.input
      },
      {
        label: t('dashboard.output'),
        data: props.trendData.map((point) => point.output_tokens),
        color: detailColors.output
      },
      {
        label: t('dashboard.tokenConsumption.cacheCreation'),
        data: props.trendData.map((point) => point.cache_creation_tokens),
        color: detailColors.cacheCreation
      },
      {
        label: t('dashboard.tokenConsumption.cacheRead'),
        data: props.trendData.map((point) => point.cache_read_tokens),
        color: detailColors.cacheRead
      },
      {
        label: t('dashboard.tokenConsumption.cacheHitRate'),
        data: props.trendData.map((point) => {
          const totalPromptTokens = point.input_tokens + point.cache_read_tokens + point.cache_creation_tokens
          return totalPromptTokens > 0 ? (point.cache_read_tokens / totalPromptTokens) * 100 : 0
        }),
        color: detailColors.cacheHitRate,
        dashed: true,
        yAxisID: 'yPercent'
      }
    ]

    datasets.push(...detailDatasets.map((dataset) => ({
      label: dataset.label,
      data: dataset.data,
      borderColor: dataset.color,
      backgroundColor: dataset.color,
      borderWidth: dataset.dashed ? 1.7 : 1.8,
      borderDash: dataset.dashed ? [5, 5] : [],
      borderCapStyle: 'round' as const,
      borderJoinStyle: 'round' as const,
      fill: false,
      tension: 0.36,
      cubicInterpolationMode: 'monotone' as const,
      pointRadius: (context: ScriptableContext<'line'>) => dataset.yAxisID === 'yPercent' && Number(context.raw) > 0 ? 2.3 : 0,
      pointHoverRadius: (context: ScriptableContext<'line'>) => dataset.yAxisID === 'yPercent' && Number(context.raw) > 0 ? 4.5 : 0,
      pointBackgroundColor: dataset.color,
      pointBorderColor: colors.value.pointBorder,
      pointBorderWidth: 1.4,
      pointHitRadius: (context: ScriptableContext<'line'>) => dataset.yAxisID === 'yPercent' && Number(context.raw) > 0 ? 8 : 0,
      clip: { top: 8, right: 12, bottom: 8, left: 12 },
      ...('yAxisID' in dataset ? { yAxisID: dataset.yAxisID } : {})
    })))
  }

  return {
    labels: props.trendData.map((point) => point.date),
    datasets
  }
})

const chartOptions = computed<ChartOptions<'line'>>(() => {
  const allZero = props.trendData.every((point) => point.total_tokens === 0
    && point.input_tokens === 0
    && point.output_tokens === 0
    && point.cache_creation_tokens === 0
    && point.cache_read_tokens === 0)

  return {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    interaction: detailsEnabled.value
      ? { mode: 'index', intersect: false }
      : { mode: 'nearest', intersect: true },
    layout: { padding: { top: 52, right: 12, bottom: 30, left: 12 } },
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: colors.value.tooltip,
        titleColor: colors.value.tooltipText,
        bodyColor: colors.value.tooltipText,
        borderColor: colors.value.tooltipBorder,
        borderWidth: 1,
        padding: 12,
        displayColors: detailsEnabled.value,
        filter: (item) => detailsEnabled.value
          ? (props.trendData[item.dataIndex]?.total_tokens ?? 0) > 0
          : typeof item.parsed.y === 'number' && item.parsed.y > 0,
        callbacks: {
          title: (items) => {
            const point = props.trendData[items[0]?.dataIndex ?? -1]
            return point ? formatFullDate(point.date, props.granularity) : ''
          },
          label: (context) => props.presentation === 'landing'
            ? `${t('home.relativeIntensity')}: ${formatActivityPercent(Number(context.parsed.y))}`
            : context.dataset.yAxisID === 'yPercent'
              ? `${context.dataset.label}: ${Number(context.parsed.y).toFixed(1)}%`
              : `${context.dataset.label}: ${formatNumber(Number(context.parsed.y))} ${t('dashboard.tokenConsumption.tokens')}`
        }
      }
    },
    scales: {
      x: {
        offset: false,
        grid: { display: false },
        border: { display: false },
        display: false,
        ticks: { display: false }
      },
      y: {
        beginAtZero: true,
        min: 0,
        ...(props.presentation === 'landing' ? { max: 100 } : allZero ? { max: 1 } : {}),
        display: false,
        grid: { display: false },
        border: { display: false },
        ticks: { display: false }
      },
      yPercent: {
        type: 'linear',
        position: 'right',
        min: 0,
        max: 100,
        display: false,
        grid: { display: false, drawOnChartArea: false },
        border: { display: false },
        ticks: { display: false }
      }
    }
  }
})

const tokenLineRevealPlugin: Plugin<'line'> = {
  id: 'tokenLineReveal',
  beforeInit(chart) {
    const revealChart = chart as RevealableChart
    revealChart.$tokenRevealProgress = prefersReducedMotion.value || props.initiallyRevealed ? 1 : 0
    revealChart.$tokenEndpointPulse = 0.15
    revealChart.$detailOpacity = detailsRendered.value ? 1 : 0
    revealChart.$detailRevealProgress = 1
    revealChart.$detailRevealActive = false
  },
  beforeDatasetsDraw(chart) {
    const revealChart = chart as RevealableChart
    const progress = revealChart.$tokenRevealProgress
    if (!chart.chartArea) return

    const { ctx, chartArea } = chart
    ctx.save()
    if (progress !== undefined && progress < 1) {
      const edgePadding = 12
      const clipLeft = chartArea.left - edgePadding
      const revealRight = chartArea.left + chartArea.width * Math.max(progress, 0)
      const clipRight = Math.min(chartArea.right + edgePadding, revealRight + edgePadding)
      ctx.beginPath()
      ctx.rect(clipLeft, chartArea.top, clipRight - clipLeft, chartArea.height)
      ctx.clip()
    }
    ctx.shadowColor = colors.value.glow
    ctx.shadowBlur = colors.value.lineGlowBlur
    revealChart.$tokenRevealClipActive = true
  },
  beforeDatasetDraw(chart, args) {
    if (args.index === 0) return

    const revealChart = chart as RevealableChart
    const { ctx, chartArea } = chart
    ctx.save()
    ctx.globalAlpha *= revealChart.$detailOpacity ?? 1

    const progress = revealChart.$detailRevealProgress ?? 1
    if (revealChart.$detailRevealActive && progress < 1 && chartArea) {
      const edgePadding = 12
      const clipLeft = chartArea.left - edgePadding
      const revealRight = chartArea.left + chartArea.width * Math.max(progress, 0)
      const clipRight = Math.min(chartArea.right + edgePadding, revealRight + edgePadding)
      ctx.beginPath()
      ctx.rect(clipLeft, chartArea.top, clipRight - clipLeft, chartArea.height)
      ctx.clip()
    }
    revealChart.$detailDatasetDrawActive = true
  },
  afterDatasetDraw(chart, args) {
    if (args.index === 0) return
    const revealChart = chart as RevealableChart
    if (!revealChart.$detailDatasetDrawActive) return
    chart.ctx.restore()
    revealChart.$detailDatasetDrawActive = false
  },
  afterDatasetsDraw(chart) {
    const revealChart = chart as RevealableChart
    const { ctx } = chart
    if (revealChart.$tokenRevealClipActive) {
      ctx.restore()
      revealChart.$tokenRevealClipActive = false
    }

    const progress = revealChart.$tokenRevealProgress ?? 1
    const points = chart.getDatasetMeta(0).data
    if (!points.length) return

    const pixelPoints = points.map((point) => {
      const pointProps = point.getProps(['x', 'y'], true)
      return { x: Number(pointProps.x), y: Number(pointProps.y) }
    })
    const first = pixelPoints[0]
    const last = pixelPoints[pixelPoints.length - 1]
    const cursorX = first.x + (last.x - first.x) * Math.min(1, Math.max(0, progress))

    ctx.save()
    const compactDates = chart.width < 620
    const valueFontSize = compactDates ? 9 : 10
    const dateFontSize = compactDates ? 8 : 9
    ctx.textAlign = 'center'
    ctx.textBaseline = 'alphabetic'

    if (!prefersReducedMotion.value) {
      const burstDuration = 250
      const pointSpan = last.x - first.x || 1
      const drawBurst = (
        point: { x: number; y: number },
        index: number,
        color: string,
        highlight: string,
        glowFade: string,
        revealProgress: number
      ) => {
        const usage = props.trendData[index]
        if (!usage || usage.total_tokens <= 0) return
        const pointProgress = (point.x - first.x) / pointSpan
        const age = (revealProgress - pointProgress) * animationDuration
        if (age < 0 || age >= burstDuration) return

        const phase = age / burstDuration
        const edgeSpace = Math.min(point.x, chart.width - point.x)
        const burstRadius = Math.max(6, Math.min(18, edgeSpace - 2))
        const alpha = Math.pow(1 - phase, 1.6)

        if (!isDarkMode.value) {
          const flare = Math.sin(Math.PI * phase)
          const haloRadius = 6 + phase * 10
          const nodeRadius = 3.4 + flare * 1.8

          ctx.save()
          ctx.globalCompositeOperation = 'source-over'
          ctx.globalAlpha = Math.pow(1 - phase, 0.8)

          const halo = ctx.createRadialGradient(point.x, point.y, 0, point.x, point.y, haloRadius)
          halo.addColorStop(0, 'rgba(23, 27, 32, 0.2)')
          halo.addColorStop(0.38, 'rgba(23, 27, 32, 0.1)')
          halo.addColorStop(1, 'rgba(23, 27, 32, 0)')
          ctx.fillStyle = halo
          ctx.beginPath()
          ctx.arc(point.x, point.y, haloRadius, 0, Math.PI * 2)
          ctx.fill()

          ctx.globalAlpha = 0.72 * Math.pow(1 - phase, 0.85)
          ctx.strokeStyle = color
          ctx.lineWidth = 1.25 * (1 - phase * 0.45)
          ctx.beginPath()
          ctx.arc(point.x, point.y, 4 + phase * 10, 0, Math.PI * 2)
          ctx.stroke()

          ctx.globalAlpha = 1
          ctx.beginPath()
          ctx.fillStyle = '#ffffff'
          ctx.arc(point.x, point.y, nodeRadius + 1.2, 0, Math.PI * 2)
          ctx.fill()
          ctx.fillStyle = color
          ctx.beginPath()
          ctx.arc(point.x, point.y, nodeRadius, 0, Math.PI * 2)
          ctx.fill()
          ctx.restore()
          return
        }

        ctx.save()
        ctx.globalCompositeOperation = 'lighter'
        ctx.globalAlpha = alpha
        ctx.shadowColor = color
        ctx.shadowBlur = 14 * (1 - phase)

        const flash = ctx.createRadialGradient(point.x, point.y, 0, point.x, point.y, burstRadius)
        flash.addColorStop(0, highlight)
        flash.addColorStop(0.18, color)
        flash.addColorStop(1, glowFade)
        ctx.fillStyle = flash
        ctx.beginPath()
        ctx.arc(point.x, point.y, burstRadius, 0, Math.PI * 2)
        ctx.fill()

        const ringRadius = 3 + burstRadius * phase
        ctx.globalAlpha = alpha * 0.85
        ctx.strokeStyle = highlight
        ctx.lineWidth = 1.6 * (1 - phase)
        ctx.beginPath()
        ctx.arc(point.x, point.y, ringRadius, 0, Math.PI * 2)
        ctx.stroke()

        ctx.strokeStyle = color
        ctx.lineWidth = 1.2 * (1 - phase)
        for (let ray = 0; ray < 8; ray += 1) {
          const angle = (Math.PI * 2 * ray) / 8
          const innerRadius = 4 + phase * 3
          const outerRadius = Math.min(burstRadius, innerRadius + 3 + (1 - phase) * 3)
          ctx.beginPath()
          ctx.moveTo(point.x + Math.cos(angle) * innerRadius, point.y + Math.sin(angle) * innerRadius)
          ctx.lineTo(point.x + Math.cos(angle) * outerRadius, point.y + Math.sin(angle) * outerRadius)
          ctx.stroke()
        }
        ctx.restore()
      }

      const drawCacheHitCross = (point: { x: number; y: number }, index: number, revealProgress: number) => {
        const usage = props.trendData[index]
        if (!usage || usage.total_tokens <= 0 || usage.cache_read_tokens <= 0) return
        const pointProgress = (point.x - first.x) / pointSpan
        const age = (revealProgress - pointProgress) * animationDuration
        if (age < 0 || age >= burstDuration) return

        const phase = age / burstDuration
        const flare = Math.sin(Math.PI * phase)
        const alpha = Math.pow(1 - phase, 1.6)
        const armLength = 4 + flare * 8

        ctx.save()
        ctx.globalCompositeOperation = 'lighter'
        ctx.globalAlpha = alpha
        ctx.lineCap = 'round'
        ctx.shadowColor = detailColors.cacheHitRate
        ctx.shadowBlur = 5 + flare * 12

        ctx.beginPath()
        ctx.moveTo(point.x - armLength, point.y)
        ctx.lineTo(point.x + armLength, point.y)
        ctx.moveTo(point.x, point.y - armLength)
        ctx.lineTo(point.x, point.y + armLength)
        ctx.strokeStyle = detailColors.cacheHitRate
        ctx.lineWidth = 1.2 + flare * 0.8
        ctx.stroke()

        ctx.shadowBlur = 0
        ctx.globalAlpha = alpha * 0.9
        ctx.beginPath()
        ctx.moveTo(point.x - armLength * 0.72, point.y)
        ctx.lineTo(point.x + armLength * 0.72, point.y)
        ctx.moveTo(point.x, point.y - armLength * 0.72)
        ctx.lineTo(point.x, point.y + armLength * 0.72)
        ctx.strokeStyle = '#f5f3ff'
        ctx.lineWidth = 0.8
        ctx.stroke()

        ctx.globalAlpha = alpha
        ctx.fillStyle = '#ffffff'
        ctx.beginPath()
        ctx.arc(point.x, point.y, 1.4 + flare * 0.8, 0, Math.PI * 2)
        ctx.fill()
        ctx.restore()
      }

      if (progress < 1) {
        pixelPoints.forEach((point, index) => {
          drawBurst(point, index, colors.value.line, colors.value.valueLabel, colors.value.glowFade, progress)
        })
      }

      if (detailsEnabled.value) {
        const detailsProgress = revealChart.$detailRevealActive
          ? revealChart.$detailRevealProgress ?? progress
          : progress
        const cacheHitPoints = chart.getDatasetMeta(5).data
          .map((point) => {
            const pointProps = point.getProps(['x', 'y'], true)
            return { x: Number(pointProps.x), y: Number(pointProps.y) }
          })
        cacheHitPoints.forEach((point, index) => {
          if (detailsProgress >= 1) return
          drawCacheHitCross(point, index, detailsProgress)
        })
      }
    }

    const maxDateLabels = Math.max(2, Math.floor(chart.width / 112))
    const dateStride = Math.max(1, Math.ceil(pixelPoints.length / maxDateLabels))
    const valueLabelStride = props.presentation === 'landing' && chart.width < 560
      ? Math.max(1, Math.ceil(pixelPoints.length / 8))
      : 1
    const valueLabelTracks: Array<Array<{ left: number; right: number }>> = []

    pixelPoints.forEach((point, index) => {
      if (progress < 1 && point.x > cursorX + 1) return

      const value = props.trendData[index]?.total_tokens ?? 0
      const showValueLabel = index === pixelPoints.length - 1 || index % valueLabelStride === 0
      if (value > 0 && showValueLabel) {
        const valueText = props.presentation === 'landing'
          ? formatActivityPercent(value)
          : formatNumber(value)
        ctx.font = `600 ${valueFontSize}px system-ui, -apple-system, sans-serif`
        const halfValueWidth = ctx.measureText(valueText).width / 2
        const valueLabelX = Math.min(chart.width - halfValueWidth - 8, Math.max(halfValueWidth + 8, point.x))
        const interval = { left: valueLabelX - halfValueWidth - 5, right: valueLabelX + halfValueWidth + 5 }
        let track = valueLabelTracks.findIndex((items) => items.every(
          (item) => interval.right <= item.left || interval.left >= item.right
        ))
        if (track < 0) track = Math.min(valueLabelTracks.length, 3)
        valueLabelTracks[track] ??= []
        valueLabelTracks[track].push(interval)

        ctx.fillStyle = colors.value.valueLabel
        ctx.textBaseline = 'bottom'
        ctx.fillText(valueText, valueLabelX, point.y - 9 - track * (valueFontSize + 4))
      }

      const showDateLabel = index === 0 || index === pixelPoints.length - 1 || index % dateStride === 0
      if (showDateLabel) {
        const dateText = formatPointDate(props.trendData[index]?.date ?? '', props.granularity)
        ctx.font = `400 ${dateFontSize}px system-ui, -apple-system, sans-serif`
        const halfDateWidth = ctx.measureText(dateText).width / 2
        const dateLabelX = Math.min(chart.width - halfDateWidth - 8, Math.max(halfDateWidth + 8, point.x))
        ctx.fillStyle = colors.value.dateLabel
        ctx.textBaseline = 'top'
        ctx.fillText(dateText, dateLabelX, point.y + 8)
      }
    })

    if (progress >= 1) {
      const lastValue = props.trendData[props.trendData.length - 1]?.total_tokens ?? 0
      if (lastValue > 0) {
        const endpointPulse = revealChart.$tokenEndpointPulse ?? 0.15
        const endpointRadius = 12 + endpointPulse * 8
        const endGlow = ctx.createRadialGradient(last.x, last.y, 0, last.x, last.y, endpointRadius)
        endGlow.addColorStop(0, colors.value.glow)
        endGlow.addColorStop(1, colors.value.glowFade)
        ctx.globalAlpha = 0.42 + endpointPulse * 0.26
        ctx.fillStyle = endGlow
        ctx.beginPath()
        ctx.arc(last.x, last.y, endpointRadius, 0, Math.PI * 2)
        ctx.fill()
      }
      ctx.restore()
      return
    }

    if (prefersReducedMotion.value || pixelPoints.length < 2) {
      ctx.restore()
      return
    }

    let segmentIndex = pixelPoints.findIndex((point) => point.x >= cursorX)
    if (segmentIndex < 1) segmentIndex = 1
    const previous = pixelPoints[segmentIndex - 1]
    const current = pixelPoints[segmentIndex]
    const segmentWidth = current.x - previous.x || 1
    const segmentProgress = Math.min(1, Math.max(0, (cursorX - previous.x) / segmentWidth))
    const cursorY = previous.y + (current.y - previous.y) * segmentProgress

    const cursorGlowRadius = 22
    const glow = ctx.createRadialGradient(cursorX, cursorY, 0, cursorX, cursorY, cursorGlowRadius)
    glow.addColorStop(0, colors.value.glow)
    glow.addColorStop(1, colors.value.glowFade)
    ctx.fillStyle = glow
    ctx.beginPath()
    ctx.arc(cursorX, cursorY, cursorGlowRadius, 0, Math.PI * 2)
    ctx.fill()

    ctx.shadowColor = colors.value.glow
    ctx.shadowBlur = colors.value.cursorGlowBlur
    ctx.fillStyle = colors.value.line
    ctx.beginPath()
    ctx.arc(cursorX, cursorY, 2.7, 0, Math.PI * 2)
    ctx.fill()
    ctx.restore()
  }
}

let animationFrame = 0
let animationStartFrame = 0
let countUpFrame = 0
let revealSequence = 0
let themeObserver: MutationObserver | null = null
let motionQuery: MediaQueryList | null = null

function runRevealAnimation() {
  const chart = chartRef.value?.chart as RevealableChart | undefined
  if (!chart) return

  if (animationFrame) cancelAnimationFrame(animationFrame)
  if (prefersReducedMotion.value) {
    chart.$tokenRevealProgress = 1
    chart.$tokenEndpointPulse = 0.15
    chart.draw()
    displayedTokens.value = totalTokens.value
    summaryVisible.value = true
    return
  }

  chart.$tokenRevealProgress = 0
  chart.$tokenEndpointPulse = 0
  chart.draw()
  const startedAt = performance.now()
  const drawFrame = (now: number) => {
    const progress = Math.min(1, (now - startedAt) / animationDuration)
    chart.$tokenRevealProgress = progress
    chart.draw()
    if (progress < 1) animationFrame = requestAnimationFrame(drawFrame)
    else {
      animationFrame = 0
      displayedTokens.value = 0
      summaryVisible.value = true
      runCountUpAnimation()
      startEndpointBreathing(chart)
    }
  }
  animationFrame = requestAnimationFrame(drawFrame)
}

function runCountUpAnimation() {
  if (countUpFrame) cancelAnimationFrame(countUpFrame)

  const target = totalTokens.value
  if (prefersReducedMotion.value || target <= 0) {
    displayedTokens.value = target
    countUpFrame = 0
    return
  }

  const startedAt = performance.now()
  const countFrame = (now: number) => {
    const progress = Math.min(1, (now - startedAt) / countUpDuration)
    const easedProgress = 1 - Math.pow(1 - progress, 3)
    displayedTokens.value = target * easedProgress

    if (progress < 1) countUpFrame = requestAnimationFrame(countFrame)
    else {
      displayedTokens.value = target
      countUpFrame = 0
    }
  }

  countUpFrame = requestAnimationFrame(countFrame)
}

function startEndpointBreathing(chart: RevealableChart) {
  if (prefersReducedMotion.value) {
    chart.$tokenEndpointPulse = 0.15
    chart.draw()
    return
  }

  const duration = 1200
  const startedAt = performance.now()
  const pulseFrame = (now: number) => {
    const progress = Math.min(1, (now - startedAt) / duration)
    chart.$tokenEndpointPulse = 0.5 - 0.5 * Math.cos(progress * Math.PI * 4)
    chart.draw()
    if (progress < 1) animationFrame = requestAnimationFrame(pulseFrame)
    else {
      chart.$tokenEndpointPulse = 0.15
      chart.draw()
      animationFrame = 0
    }
  }
  animationFrame = requestAnimationFrame(pulseFrame)
}

async function animateDetailVisibility(enabled: boolean) {
  const sequence = ++detailAnimationSequence
  if (detailAnimationFrame) cancelAnimationFrame(detailAnimationFrame)
  detailAnimationFrame = 0

  if (enabled) {
    detailsRendered.value = true
    await nextTick()
    if (sequence !== detailAnimationSequence) return

    const chart = chartRef.value?.chart as RevealableChart | undefined
    if (!chart) return

    chart.$detailOpacity = 1
    chart.$detailRevealProgress = 0
    chart.$detailRevealActive = true
    if (prefersReducedMotion.value || props.loading || props.error || !props.trendData.length) {
      chart.$detailOpacity = 1
      chart.$detailRevealProgress = 1
      chart.$detailRevealActive = false
      chart.draw()
      return
    }

    chart.draw()
    const startedAt = performance.now()
    const drawDetailsFrame = (now: number) => {
      if (sequence !== detailAnimationSequence) return
      const progress = Math.min(1, (now - startedAt) / animationDuration)
      chart.$detailRevealProgress = progress
      chart.draw()
      if (progress < 1) detailAnimationFrame = requestAnimationFrame(drawDetailsFrame)
      else {
        detailAnimationFrame = 0
        chart.$detailRevealActive = false
        chart.draw()
      }
    }
    detailAnimationFrame = requestAnimationFrame(drawDetailsFrame)
    return
  }

  if (!detailsRendered.value) return
  const chart = chartRef.value?.chart as RevealableChart | undefined
  if (!chart) {
    detailsRendered.value = false
    return
  }

  if (prefersReducedMotion.value || props.loading || props.error || !props.trendData.length) {
    chart.$detailOpacity = 0
    chart.$detailRevealActive = false
    chart.$detailRevealProgress = 1
    detailsRendered.value = false
    chart.draw()
    return
  }

  const startOpacity = chart.$detailOpacity ?? 1
  if (startOpacity <= 0) {
    chart.$detailOpacity = 0
    chart.$detailRevealActive = false
    chart.$detailRevealProgress = 1
    detailsRendered.value = false
    chart.draw()
    return
  }

  const startedAt = performance.now()
  const fadeDetailsFrame = (now: number) => {
    if (sequence !== detailAnimationSequence) return
    const progress = Math.min(1, (now - startedAt) / detailFadeDuration)
    chart.$detailOpacity = startOpacity * (1 - progress)
    chart.draw()
    if (progress < 1) detailAnimationFrame = requestAnimationFrame(fadeDetailsFrame)
    else {
      detailAnimationFrame = 0
      chart.$detailOpacity = 0
      chart.$detailRevealProgress = 1
      chart.$detailRevealActive = false
      detailsRendered.value = false
      chart.draw()
    }
  }
  detailAnimationFrame = requestAnimationFrame(fadeDetailsFrame)
}

async function scheduleRevealAnimation() {
  const sequence = ++revealSequence
  summaryVisible.value = false
  displayedTokens.value = 0
  if (countUpFrame) cancelAnimationFrame(countUpFrame)
  countUpFrame = 0
  if (animationFrame) cancelAnimationFrame(animationFrame)
  if (animationStartFrame) cancelAnimationFrame(animationStartFrame)
  animationFrame = 0
  animationStartFrame = 0
  if (props.loading || props.error || !props.trendData.length) return

  await nextTick()
  if (sequence !== revealSequence || props.loading || props.error) return
  if (!props.autoReveal) {
    const chart = chartRef.value?.chart as RevealableChart | undefined
    if (chart) {
      chart.$tokenRevealProgress = 1
      chart.$tokenEndpointPulse = 0.15
      chart.draw()
    }
    return
  }
  animationStartFrame = requestAnimationFrame(() => {
    animationStartFrame = 0
    runRevealAnimation()
  })
}

function playRetractAnimation(): Promise<void> {
  const chart = chartRef.value?.chart as RevealableChart | undefined
  if (!chart) return Promise.resolve()

  if (animationFrame) cancelAnimationFrame(animationFrame)
  if (animationStartFrame) cancelAnimationFrame(animationStartFrame)
  if (countUpFrame) cancelAnimationFrame(countUpFrame)
  animationFrame = 0
  animationStartFrame = 0
  countUpFrame = 0
  summaryVisible.value = false
  displayedTokens.value = 0

  if (prefersReducedMotion.value) {
    chart.$tokenRevealProgress = 0
    chart.$tokenEndpointPulse = 0
    chart.draw()
    return Promise.resolve()
  }

  return new Promise((resolve) => {
    const sequence = ++revealSequence
    chart.$tokenRevealProgress = 1
    chart.$tokenEndpointPulse = 0
    chart.draw()
    const startedAt = performance.now()
    const retractFrame = (now: number) => {
      if (sequence !== revealSequence) {
        resolve()
        return
      }

      const progress = Math.max(0, 1 - (now - startedAt) / 1200)
      chart.$tokenRevealProgress = progress
      chart.$tokenEndpointPulse = 0
      chart.draw()
      if (progress > 0) {
        animationFrame = requestAnimationFrame(retractFrame)
      } else {
        animationFrame = 0
        resolve()
      }
    }
    animationFrame = requestAnimationFrame(retractFrame)
  })
}

watch(
  () => [props.trendData, props.range, props.granularity, props.loading, props.error] as const,
  () => { void scheduleRevealAnimation() },
  { deep: true, flush: 'post' }
)

function formatFullDate(value: string, granularity: 'day' | 'hour'): string {
  const date = parseDate(value)
  if (!date) return value
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    ...(granularity === 'hour' ? { hour: '2-digit', minute: '2-digit' } : {})
  }
  return new Intl.DateTimeFormat(locale.value, options).format(date)
}

function formatPointDate(value: string, granularity: 'day' | 'hour'): string {
  const date = parseDate(value)
  if (!date) return value
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  if (granularity === 'hour') {
    const hour = String(date.getHours()).padStart(2, '0')
    const minute = String(date.getMinutes()).padStart(2, '0')
    return `${hour}:${minute}`
  }
  return `${month}.${day}`
}

function formatActivityPercent(value: number): string {
  const digits = value >= 10 ? 0 : value >= 1 ? 1 : 2
  return `${new Intl.NumberFormat(locale.value, { maximumFractionDigits: digits }).format(value)}%`
}

function parseDate(value: string): Date | null {
  if (/[zZ]$|[+-]\d{2}:?\d{2}$/.test(value)) {
    const timestamp = new Date(value)
    return Number.isNaN(timestamp.getTime()) ? null : timestamp
  }
  const match = /^(\d{4})-(\d{2})-(\d{2})(?:[ T](\d{2}):(\d{2}))?/.exec(value)
  if (!match) return null
  return new Date(
    Number(match[1]),
    Number(match[2]) - 1,
    Number(match[3]),
    Number(match[4] ?? 0),
    Number(match[5] ?? 0)
  )
}

function updateTheme() {
  isDarkMode.value = document.documentElement.classList.contains('dark')
}

function updateMotionPreference(event?: MediaQueryListEvent) {
  prefersReducedMotion.value = event?.matches ?? motionQuery?.matches ?? false
  if (prefersReducedMotion.value) {
    if (animationFrame) cancelAnimationFrame(animationFrame)
    if (animationStartFrame) cancelAnimationFrame(animationStartFrame)
    animationFrame = 0
    animationStartFrame = 0
    if (countUpFrame) cancelAnimationFrame(countUpFrame)
    countUpFrame = 0
    if (detailAnimationFrame) cancelAnimationFrame(detailAnimationFrame)
    detailAnimationFrame = 0
    detailAnimationSequence += 1
    const chart = chartRef.value?.chart as RevealableChart | undefined
    detailsRendered.value = detailsEnabled.value
    if (chart) {
      chart.$tokenRevealProgress = 1
      chart.$tokenEndpointPulse = 0.15
      chart.$detailOpacity = detailsEnabled.value ? 1 : 0
      chart.$detailRevealProgress = 1
      chart.$detailRevealActive = false
      chart.draw()
      displayedTokens.value = totalTokens.value
      summaryVisible.value = true
    }
  }
}

onMounted(() => {
  updateTheme()
  themeObserver = new MutationObserver(updateTheme)
  themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })

  motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
  updateMotionPreference()
  motionQuery.addEventListener('change', updateMotionPreference)
  if (props.autoReveal && !props.initiallyRevealed) {
    void scheduleRevealAnimation()
  }
})

onUnmounted(() => {
  if (animationFrame) cancelAnimationFrame(animationFrame)
  if (animationStartFrame) cancelAnimationFrame(animationStartFrame)
  if (countUpFrame) cancelAnimationFrame(countUpFrame)
  if (detailAnimationFrame) cancelAnimationFrame(detailAnimationFrame)
  detailAnimationSequence += 1
  themeObserver?.disconnect()
  motionQuery?.removeEventListener('change', updateMotionPreference)
})

defineExpose({ playRetractAnimation })
</script>

<style scoped>
.token-detail-legend {
  transition: opacity 240ms ease;
}

.token-detail-toggle-enter-active,
.token-detail-toggle-appear-active {
  transition: opacity 480ms ease, transform 480ms cubic-bezier(0.2, 0.7, 0.2, 1);
  transition-delay: 260ms;
}

.token-detail-toggle-enter-from,
.token-detail-toggle-appear-from {
  opacity: 0;
  transform: translateY(-7px);
}

.token-summary-enter-active,
.token-summary-leave-active {
  transition: opacity 240ms ease, transform 240ms cubic-bezier(0.22, 1, 0.36, 1);
}

.token-summary-enter-from,
.token-summary-leave-to {
  opacity: 0;
  transform: translateY(6px);
}

.token-summary-value {
  display: inline-grid;
  grid-template-areas: 'value';
  vertical-align: baseline;
}

.token-summary-value > span {
  grid-area: value;
  white-space: nowrap;
}

.token-summary-value-reserve {
  visibility: hidden;
  pointer-events: none;
}

.token-summary-value-current {
  justify-self: start;
}

.token-summary-reveal {
  animation: token-summary-rise 800ms cubic-bezier(0.22, 1, 0.36, 1) both;
}

.token-summary-reveal--unit {
  animation-delay: 90ms;
}

.token-summary-reveal--period {
  animation-delay: 180ms;
}

.token-summary-reveal--recent-hours {
  animation-delay: 260ms;
}

@keyframes token-summary-rise {
  from { opacity: 0; transform: translateY(28px); }
  to { opacity: 1; transform: translateY(0); }
}

:deep(.token-period-select .select-trigger) {
  width: auto;
  min-height: 30px;
  gap: 0.35rem;
  border-color: transparent;
  border-radius: 9999px;
  padding: 0.2rem 0.55rem;
  background: transparent;
  color: rgb(107 114 128);
  font-size: 0.75rem;
  line-height: 1rem;
  transition: background-color 160ms ease, color 160ms ease;
}

:deep(.token-period-select .select-trigger:hover),
:deep(.token-period-select .select-trigger-open) {
  border-color: transparent;
  background: rgb(15 23 42 / 0.06);
  color: rgb(31 41 55);
  box-shadow: none;
}

:deep(.dark .token-period-select .select-trigger) {
  color: rgb(156 163 175);
}

:deep(.dark .token-period-select .select-trigger:hover),
:deep(.dark .token-period-select .select-trigger-open) {
  background: rgb(255 255 255 / 0.08);
  color: rgb(229 231 235);
}

:global(.select-dropdown-portal.token-range-dropdown) { min-width: 144px; }

:global(.select-dropdown-portal.token-range-dropdown-light) {
  background-color: #f9fafb !important;
  border-color: #e5e7eb !important;
  color: #374151 !important;
}

:global(.select-dropdown-portal.token-range-dropdown-light .select-option) { color: #374151 !important; }

:global(.select-dropdown-portal.token-range-dropdown-light .select-option-focused),
:global(.select-dropdown-portal.token-range-dropdown-light .select-option:hover) {
  background-color: #f3f4f6 !important;
  color: #111827 !important;
}

:global(.select-dropdown-portal.token-range-dropdown-light .select-option-selected) {
  background-color: #e5e7eb !important;
  color: #111827 !important;
}

:global(.select-dropdown-portal.token-range-dropdown-dark) {
  background-color: #090b0d !important;
  border-color: #252a2e !important;
  color: #d1d5db !important;
}

:global(.select-dropdown-portal.token-range-dropdown-dark .select-option) { color: #d1d5db !important; }

:global(.select-dropdown-portal.token-range-dropdown-dark .select-option-focused),
:global(.select-dropdown-portal.token-range-dropdown-dark .select-option:hover) {
  background-color: #1c2024 !important;
  color: #f9fafb !important;
}

:global(.select-dropdown-portal.token-range-dropdown-dark .select-option-selected) {
  background-color: #20262a !important;
  color: #f9fafb !important;
}

@media (prefers-reduced-motion: reduce) {
  .token-detail-legend { transition: none; }
  .token-detail-toggle-enter-active,
  .token-detail-toggle-appear-active {
    transition-duration: 1ms;
    transition-delay: 0ms;
  }
  .token-summary-enter-active,
  .token-summary-leave-active { transition: none; }
  .token-summary-reveal { animation: none; }
}
</style>
