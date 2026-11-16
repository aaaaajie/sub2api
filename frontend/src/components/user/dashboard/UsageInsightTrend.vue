<template>
  <div class="insight-chart" role="img" :aria-label="label">
    <Bar v-if="kind === 'cost' || kind === 'requests'" :data="barData" :options="barOptions" />
    <Line v-else :data="lineData" :options="lineOptions" />
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePreferredReducedMotion } from '@vueuse/core'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip, Legend, Filler, type ScriptableContext, type ChartOptions } from 'chart.js'
import { Line, Bar } from 'vue-chartjs'
import type { TrendDataPoint } from '@/types'
import { formatTokensK } from '@/utils/format'
import { chartPalette } from './chartPalette'
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip, Legend, Filler)
const props = defineProps<{ points: TrendDataPoint[]; kind: 'tokens' | 'cost' | 'requests' | 'cache'; dark: boolean; label: string; details?: boolean }>()
const { t } = useI18n()
const reduced = usePreferredReducedMotion()
const colors = computed(() => chartPalette(props.dark))
const ink = computed(() => props.kind === 'cache' ? colors.value.cache : props.kind === 'cost' ? colors.value.cost : props.kind === 'requests' ? colors.value.requests : colors.value.total)
function areaFill(context: ScriptableContext<'line'>) {
  const { ctx, chartArea } = context.chart
  if (!chartArea) return `${ink.value}15`
  const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom)
  gradient.addColorStop(0, `${ink.value}${props.dark ? '38' : '25'}`)
  gradient.addColorStop(1, `${ink.value}00`)
  return gradient
}
const labels = computed(() => props.points.map(p => p.date.includes(' ') ? p.date.slice(5, 16) : p.date.slice(5)))
const hit = (p: TrendDataPoint) => { const base = p.input_tokens + p.cache_read_tokens + p.cache_creation_tokens; return base ? p.cache_read_tokens / base * 100 : 0 }
const lineData = computed(() => ({ labels: labels.value, datasets: [
  { label: props.kind === 'cache' ? t('dashboard.tokenConsumption.cacheHitRate') : t('dashboard.tokenConsumption.totalTokens'), data: props.points.map(p => props.kind === 'cache' ? hit(p) : p.total_tokens), borderColor: ink.value, backgroundColor: areaFill, fill: true, borderWidth: 2, tension: .3, pointRadius: 2, pointHoverRadius: 5 },
  ...(props.details && props.kind === 'tokens' ? (['input_tokens', 'output_tokens', 'cache_creation_tokens', 'cache_read_tokens'] as const).map((key, index) => ({ label: [t('dashboard.input'), t('dashboard.output'), t('dashboard.tokenConsumption.cacheCreation'), t('dashboard.tokenConsumption.cacheRead')][index]!, data: props.points.map(p => p[key]), borderColor: colors.value.parts[index]!, backgroundColor: 'transparent', fill: false, borderWidth: 1.8, tension: .3, pointRadius: 0, pointHoverRadius: 4 })) : [])
] }))
const barData = computed(() => ({ labels: labels.value, datasets: [{ label: props.label, data: props.points.map(p => props.kind === 'cost' ? p.actual_cost : p.requests), backgroundColor: ink.value, borderRadius: 5, maxBarThickness: 18 }] }))
const base = computed(() => ({ responsive: true, maintainAspectRatio: false, animation: reduced.value === 'reduce' ? false as const : { duration: 500 }, interaction: { mode: 'index' as const, intersect: false }, plugins: { legend: { display: props.details && props.kind === 'tokens', position: 'bottom' as const, labels: { usePointStyle: true, boxWidth: 6, boxHeight: 6, padding: 15, color: props.dark ? '#a6b3bb' : '#65737d', font: { size: 10 } } }, tooltip: { displayColors: true, padding: 12 } }, scales: { x: { grid: { display: false }, border: { display: false }, ticks: { color: props.dark ? '#89969e' : '#74828b', font: { size: 10 }, maxTicksLimit: 9, maxRotation: 0 } }, y: { beginAtZero: true, ...(props.kind === 'cache' ? { max: 100 } : {}), border: { display: false }, grid: { color: props.dark ? '#ffffff0b' : '#52606c12' }, ticks: { maxTicksLimit: 5, color: props.dark ? '#89969e' : '#74828b', font: { size: 10 }, callback: (value: string | number) => props.kind === 'cache' ? `${value}%` : props.kind === 'cost' ? `$${Number(value).toFixed(2)}` : formatTokensK(Number(value)) } } } }))
const lineOptions = computed<ChartOptions<'line'>>(() => base.value)
const barOptions = computed<ChartOptions<'bar'>>(() => base.value)
</script>
<style scoped>
.insight-chart { height: 240px; min-width: 0; }
@media(max-width: 600px) { .insight-chart { height: 210px; } }
</style>
