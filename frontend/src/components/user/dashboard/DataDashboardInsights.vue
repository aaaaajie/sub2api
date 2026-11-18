<template>
  <section ref="section" class="analytics" :class="{ revealed: visible }" :aria-label="t('nav.dataDashboard')">
    <div class="analytics-toolbar">
      <DateRangePicker v-model:start-date="startDate" v-model:end-date="endDate" @change="refreshRange" />
      <div class="toolbar-group">
        <span>{{ t('dashboard.dataDashboard.accountWide') }}</span>
        <div class="segmented"><button v-for="unit in (['day', 'hour'] as const)" :key="unit" :class="{ selected: granularity === unit }" :aria-pressed="granularity === unit" @click="setGranularity(unit)">{{ t(`dashboard.${unit}`) }}</button></div>
        <button class="refresh-button" :disabled="loading || overviewLoading || logsLoading" @click="refreshAll"><Icon name="refresh" size="sm" />{{ t('common.refresh') }}</button>
      </div>
    </div>

    <div class="summary-strip" :aria-busy="loading">
      <div v-for="metric in metrics" :key="metric.label"><span>{{ metric.label }}</span><strong>{{ loading || statsError ? '—' : metric.value }}</strong><small>{{ loading || statsError ? rangeLabel : metric.hint || rangeLabel }}</small></div>
      <button v-if="statsError" type="button" class="retry" @click="load">{{ t('dashboard.dataDashboard.retry') }}</button>
    </div>

    <div class="two-columns">
      <UsageInsightPanel :title="t('dashboard.dataDashboard.composition')" :loading="loading" :error="snapshotError" @retry="load">
        <div class="composition-layout">
          <div class="composition-ring"><Doughnut v-if="visible" :data="donutData" :options="donutOptions" /><div class="ring-label"><strong>{{ compact(tokenTotal) }}</strong><span>{{ t('dashboard.dataDashboard.total') }}</span></div></div>
          <div class="composition-key"><div v-for="(part, i) in parts" :key="part.label"><i :style="{ background: palette[i] }" /><span>{{ part.label }}</span><strong :title="number(part.value)">{{ compact(part.value) }}</strong><small>{{ percent(part.value, tokenTotal) }}</small></div></div>
        </div>
      </UsageInsightPanel>
      <UsageDistribution :title="t('dashboard.modelDistribution')" :rows="modelRows" :loading="loading" :error="snapshotError" @retry="load" />
    </div>

    <UsageInsightPanel :title="t('dashboard.dataDashboard.requestTrend')" :loading="loading" :error="snapshotError" :empty="!trend.length" @retry="load">
      <template #actions><div class="chart-actions"><div class="segmented"><button v-for="item in trendModes" :key="item.value" type="button" :class="{ selected: trendMode === item.value }" :aria-pressed="trendMode === item.value" @click="trendMode = item.value">{{ item.label }}</button></div><TrendDetailToggle v-if="trendMode === 'tokens'" v-model="details" /></div></template>
      <UsageInsightTrend :points="trend" :kind="trendMode" :dark="dark" :details="details" :label="t('dashboard.dataDashboard.requestTrend')" />
    </UsageInsightPanel>

    <UsageInsightPanel :title="t('dashboard.dataDashboard.modelDetails')" :loading="loading" :error="snapshotError" :empty="!models.length" @retry="load">
      <div class="table-scroll"><table><thead><tr><th>{{ t('dashboard.model') }}</th><th>{{ t('dashboard.dataDashboard.total') }}</th><th v-for="part in parts" :key="part.label">{{ part.label }}</th><th>{{ t('dashboard.requests') }}</th><th>{{ t('dashboard.dataDashboard.spend') }}</th><th>{{ t('usage.standardCost') }}</th></tr></thead><tbody><tr v-for="item in sortedModels" :key="item.model"><td :title="item.model">{{ item.model }}</td><td>{{ number(item.total_tokens) }}</td><td>{{ number(item.input_tokens) }}</td><td>{{ number(item.output_tokens) }}</td><td>{{ number(item.cache_creation_tokens) }}</td><td>{{ number(item.cache_read_tokens) }}</td><td>{{ number(item.requests) }}</td><td>{{ money(item.actual_cost) }}</td><td class="muted">{{ money(item.cost) }}</td></tr></tbody></table></div>
    </UsageInsightPanel>

    <div class="two-columns">
      <UsageDistribution :title="t('dashboard.groupDistribution')" :rows="groupRows" :loading="loading" :error="snapshotError" @retry="load" />
      <UsageDistribution :title="t('usage.endpointDistribution')" :rows="endpointRows" :loading="loading" :error="statsError" @retry="load" />
    </div>

    <UsageInsightPanel :title="t('dashboard.dataDashboard.cacheUsage')" :loading="loading" :error="snapshotError" :empty="!trend.length" @retry="load">
      <div class="cache-layout"><div class="cache-stats"><span>{{ t('dashboard.tokenConsumption.cacheHitRate') }}<strong>{{ cacheRate }}</strong></span><span>{{ t('dashboard.tokenConsumption.cacheRead') }}<b>{{ compact(sum('cache_read_tokens')) }}</b></span><span>{{ t('dashboard.tokenConsumption.cacheCreation') }}<b>{{ compact(sum('cache_creation_tokens')) }}</b></span></div><UsageInsightTrend :points="trend" kind="cache" :dark="dark" :label="t('dashboard.tokenConsumption.cacheHitRate')" /></div>
    </UsageInsightPanel>

    <UsageInsightPanel :title="t('dashboard.dataDashboard.accountStatus')" :loading="overviewLoading" :error="overviewError" @retry="loadOverview">
      <div v-if="overview" class="account-status">
        <RouterLink to="/keys"><span>{{ t('dashboard.apiKeys') }}</span><strong>{{ overview.active_api_keys }} / {{ overview.total_api_keys }}</strong><small>{{ t('dashboard.dataDashboard.activeOfTotal') }}</small></RouterLink>
        <div><span>RPM</span><strong>{{ number(overview.rpm) }}</strong><small>{{ t('dashboard.dataDashboard.lastFiveMinutes') }}</small></div>
        <div><span>TPM</span><strong>{{ compact(overview.tpm) }}</strong><small>{{ t('dashboard.dataDashboard.lastFiveMinutes') }}</small></div>
      </div>
    </UsageInsightPanel>

    <div class="two-columns">
      <UsageInsightPanel v-for="period in (['total', 'today'] as const)" :key="period" :title="t(`dashboard.dataDashboard.${period}Overview`)" :loading="overviewLoading" :error="overviewError" @retry="loadOverview">
        <template #actions><span class="muted">{{ t('dashboard.dataDashboard.accountWide') }}</span></template>
        <dl v-if="overview" class="overview-list"><div><dt>{{ t('dashboard.dataDashboard.total') }}</dt><dd>{{ number(period === 'today' ? overview.today_tokens : overview.total_tokens) }}</dd></div><div><dt>{{ t('dashboard.requests') }}</dt><dd>{{ number(period === 'today' ? overview.today_requests : overview.total_requests) }}</dd></div><div><dt>{{ t('dashboard.dataDashboard.spend') }}</dt><dd>{{ money(period === 'today' ? overview.today_actual_cost : overview.total_actual_cost) }}</dd></div><div><dt>{{ t('usage.standardCost') }}</dt><dd>{{ money(period === 'today' ? overview.today_cost : overview.total_cost) }}</dd></div></dl>
      </UsageInsightPanel>
    </div>

    <UsageInsightPanel :title="t('dashboard.dataDashboard.recentLogs')" :loading="logsLoading" :error="logsError" :empty="!logs.length" @retry="loadLogs">
      <template #actions><RouterLink to="/usage" class="all-records">{{ t('dashboard.viewAllUsage') }}<Icon name="arrowRight" size="xs" /></RouterLink></template>
      <div class="table-scroll"><table><thead><tr><th>{{ t('dashboard.dataDashboard.time') }}</th><th>{{ t('dashboard.model') }}</th><th>{{ t('keys.apiKey') }} / {{ t('dashboard.group') }}</th><th>{{ t('dashboard.input') }} / {{ t('dashboard.output') }}</th><th>{{ t('dashboard.tokenConsumption.cacheRead') }} / {{ t('dashboard.tokenConsumption.cacheCreation') }}</th><th>{{ t('usage.firstToken') }} / {{ t('dashboard.dataDashboard.duration') }}</th><th>{{ t('dashboard.dataDashboard.spend') }} / {{ t('usage.standardCost') }}</th></tr></thead>
        <tbody><tr v-for="log in logs" :key="log.id"><td>{{ formatDateTimeToMinute(log.created_at, locale) }}</td><td :title="log.model">{{ log.model }}</td><td><span>{{ log.api_key?.name || '—' }}</span><small class="cell-secondary">{{ log.group?.name || '—' }}</small></td><td>{{ number(log.input_tokens) }}<small class="cell-secondary">{{ number(log.output_tokens) }}</small></td><td>{{ number(log.cache_read_tokens) }}<small class="cell-secondary">{{ number(log.cache_creation_tokens) }}</small></td><td>{{ log.first_token_ms == null ? '—' : duration(log.first_token_ms) }}<small class="cell-secondary">{{ log.duration_ms == null ? '—' : duration(log.duration_ms) }}</small></td><td>{{ money(log.actual_cost) }}<small class="cell-secondary">{{ money(log.total_cost) }}</small></td></tr></tbody></table></div>
    </UsageInsightPanel>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { usePreferredReducedMotion } from '@vueuse/core'
import { Chart as ChartJS, ArcElement, Tooltip, type ChartOptions } from 'chart.js'
import { Doughnut } from 'vue-chartjs'
import { usageAPI, type UserDashboardStats } from '@/api/usage'
import type { GroupStat, ModelStat, TrendDataPoint, UsageLog, UsageStatsResponse } from '@/types'
import { formatDateLocalInput, formatDateTimeToMinute, formatTokensK } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import UsageInsightPanel from './UsageInsightPanel.vue'
import TrendDetailToggle from './TrendDetailToggle.vue'
import UsageInsightTrend from './UsageInsightTrend.vue'
import UsageDistribution from './UsageDistribution.vue'
import { chartPalette } from './chartPalette'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
ChartJS.register(ArcElement, Tooltip)
const startDate = ref(formatDateLocalInput(new Date(Date.now() - 6 * 86400000)))
const endDate = ref(formatDateLocalInput(new Date()))
const granularity = ref<'day' | 'hour'>('day')
const { t, locale } = useI18n()
const dark = ref(document.documentElement.classList.contains('dark'))
const reduced = usePreferredReducedMotion()
const section = ref<HTMLElement | null>(null)
const visible = ref(false)
const models = ref<ModelStat[]>([])
const groups = ref<GroupStat[]>([])
const trend = ref<TrendDataPoint[]>([])
const stats = ref<UsageStatsResponse | null>(null)
const overview = ref<UserDashboardStats | null>(null)
const logs = ref<UsageLog[]>([])
const loading = ref(true), snapshotError = ref(false), statsError = ref(false)
const logsLoading = ref(true), logsError = ref(false), overviewLoading = ref(true), overviewError = ref(false)
const trendMode = ref<'tokens' | 'requests' | 'cost'>('tokens')
const details = ref(true)
let sequence = 0, logSequence = 0
let observer: IntersectionObserver | undefined, themeObserver: MutationObserver | undefined
const rangeLabel = computed(() => `${startDate.value} — ${endDate.value}`)
const trendModes = computed(() => [{ value: 'tokens' as const, label: 'Token' }, { value: 'requests' as const, label: t('dashboard.requests') }, { value: 'cost' as const, label: t('dashboard.dataDashboard.spend') }])
const palette = computed(() => chartPalette(dark.value).parts)
function sum(key: 'input_tokens' | 'output_tokens' | 'cache_creation_tokens' | 'cache_read_tokens') { return models.value.reduce((acc, item) => acc + item[key], 0) }
const parts = computed(() => [{ label: t('dashboard.input'), value: sum('input_tokens') }, { label: t('dashboard.output'), value: sum('output_tokens') }, { label: t('dashboard.tokenConsumption.cacheCreation'), value: sum('cache_creation_tokens') }, { label: t('dashboard.tokenConsumption.cacheRead'), value: sum('cache_read_tokens') }])
const tokenTotal = computed(() => parts.value.reduce((acc, item) => acc + item.value, 0))
const sortedModels = computed(() => [...models.value].sort((a, b) => b.total_tokens - a.total_tokens))
const cacheRate = computed(() => percent(sum('cache_read_tokens'), sum('input_tokens') + sum('cache_creation_tokens') + sum('cache_read_tokens')))
const metrics = computed(() => [
  { label: t('dashboard.dataDashboard.total'), value: number(stats.value?.total_tokens ?? 0) },
  { label: t('dashboard.requests'), value: number(stats.value?.total_requests ?? 0) },
  { label: t('dashboard.dataDashboard.spend'), value: money(stats.value?.total_actual_cost ?? 0), hint: `${t('usage.standardCost')} ${money(stats.value?.total_cost ?? 0)}` },
  { label: t('dashboard.tokenConsumption.cacheHitRate'), value: stats.value ? percent(stats.value.total_cache_read_tokens, stats.value.total_input_tokens + stats.value.total_cache_creation_tokens + stats.value.total_cache_read_tokens) : '—' },
  { label: t('dashboard.avgResponse'), value: duration(stats.value?.average_duration_ms ?? 0) }
])
const donutData = computed(() => ({ labels: tokenTotal.value ? parts.value.map(p => p.label) : [''], datasets: [{ data: tokenTotal.value ? parts.value.map(p => p.value) : [1], backgroundColor: tokenTotal.value ? palette.value : [dark.value ? '#293238' : '#e8edef'], borderWidth: 0, spacing: tokenTotal.value ? 2 : 0, hoverOffset: 3 }] }))
const donutOptions = computed<ChartOptions<'doughnut'>>(() => ({ responsive: true, maintainAspectRatio: false, cutout: '73%', layout: { padding: 5 }, animation: reduced.value === 'reduce' ? false : { duration: 600 }, plugins: { legend: { display: false }, tooltip: { enabled: tokenTotal.value > 0, callbacks: { label: ctx => `${number(ctx.parsed)} · ${percent(ctx.parsed, tokenTotal.value)}` } } } }))
const number = (n: number) => n.toLocaleString(locale.value)
const compact = formatTokensK
const money = (n: number) => `$${n.toLocaleString(locale.value, { minimumFractionDigits: 4, maximumFractionDigits: 4 })}`
const duration = (n: number) => n >= 1000 ? `${(n / 1000).toFixed(2)} s` : `${Math.round(n)} ms`
const percent = (n: number, total: number) => `${total ? (n / total * 100).toFixed(1) : '0'}%`
const modelRows = computed(() => models.value.map(row => ({ ...row, id: row.model, name: row.model })))
const groupRows = computed(() => groups.value.map(row => ({ ...row, id: String(row.group_id), name: row.group_name || '—' })))
const endpointRows = computed(() => (stats.value?.endpoints ?? []).map(row => ({ ...row, id: row.endpoint, name: row.endpoint || '—' })))
function params() {
  return { start_date: startDate.value, end_date: endDate.value, timezone: Intl.DateTimeFormat().resolvedOptions().timeZone }
}
function refreshRange() { void load(); void loadLogs() }
function setGranularity(value: 'day' | 'hour') { if (granularity.value === value) return; granularity.value = value; void load() }
function refreshAll() { refreshRange(); void loadOverview() }

async function load() {
  const id = ++sequence
  loading.value = true
  snapshotError.value = statsError.value = false
  const query = params()
  const results = await Promise.allSettled([
    usageAPI.getDashboardSnapshotV2({ ...query, granularity: granularity.value, include_trend: true, include_model_stats: true, include_group_stats: true }),
    usageAPI.getStats(query)
  ])
  if (id !== sequence) return
  const [snapshot, summary] = results
  if (snapshot.status === 'fulfilled') {
    models.value = snapshot.value.models ?? []
    groups.value = snapshot.value.groups ?? []
    trend.value = (snapshot.value.trend ?? []).slice().sort((a, b) => a.date.localeCompare(b.date))
  } else { models.value = []; groups.value = []; trend.value = []; snapshotError.value = true }
  if (summary.status === 'fulfilled') stats.value = summary.value
  else { stats.value = null; statsError.value = true }
  loading.value = false
}
async function loadLogs() {
  const id = ++logSequence
  logsLoading.value = true; logsError.value = false
  try { const result = await usageAPI.query({ ...params(), page: 1, page_size: 8, sort_by: 'created_at', sort_order: 'desc' }); if (id === logSequence) logs.value = result.items }
  catch { if (id === logSequence) { logsError.value = true; logs.value = [] } }
  finally { if (id === logSequence) logsLoading.value = false }
}
async function loadOverview() {
  overviewLoading.value = true; overviewError.value = false
  try { overview.value = await usageAPI.getDashboardStats() }
  catch { overviewError.value = true }
  finally { overviewLoading.value = false }
}
onMounted(() => {
  refreshAll()
  themeObserver = new MutationObserver(() => { dark.value = document.documentElement.classList.contains('dark') })
  themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  observer = new IntersectionObserver(entries => { if (entries.some(entry => entry.isIntersecting)) { visible.value = true; observer?.disconnect() } }, { threshold: 0 })
  if (section.value) observer.observe(section.value)
})
onUnmounted(() => { sequence++; logSequence++; observer?.disconnect(); themeObserver?.disconnect() })
</script>

<style scoped>
.refresh-button { display: flex; align-items: center; gap: 6px; padding: 7px 10px; border-radius: 7px; border: 1px solid color-mix(in srgb, var(--muted) 25%, transparent); }
.refresh-button:disabled { opacity: .45; }
.account-status { display: grid; grid-template-columns: repeat(3, minmax(0,1fr)); gap: 20px; }
.account-status > * { display: grid; gap: 8px; }
.account-status span, .account-status small, .cell-secondary { color: var(--muted); font-size: 10px; }
.account-status strong { font-size: 23px; font-weight: 550; font-variant-numeric: tabular-nums; }
.cell-secondary { display: block; margin-top: 4px; }
.analytics-toolbar { flex-wrap: wrap; }

.analytics { --accent: rgb(var(--color-primary-900)); --muted: rgb(var(--color-primary-500)); max-width: 1340px; margin: auto; padding: 24px 32px 40px; display: grid; gap: 20px; color: rgb(var(--color-primary-900)); opacity: 0; transform: translateY(12px); transition: opacity .6s, transform .6s; font-size: 12px; }
.analytics.revealed { opacity: 1; transform: none; }
.analytics-toolbar, .toolbar-group, .chart-actions { display: flex; align-items: center; gap: 12px; }
.analytics-toolbar { justify-content: space-between; flex-wrap: wrap; padding: 0 4px 4px; }
.toolbar-group > span, .toolbar-group label { color: var(--muted); font-size: 11px; }
.segmented { display: inline-flex; padding: 3px; gap: 2px; border: 1px solid rgb(var(--color-primary-200) / .6); border-radius: 8px; background: rgb(var(--color-primary-50)); }
.segmented button { padding: 6px 12px; border-radius: 5px; font-size: 11px; white-space: nowrap; transition: background .2s; }
.segmented .selected { background: var(--accent); color: rgb(var(--color-primary-50)); box-shadow: 0 2px 4px rgb(0 0 0 / .08); }
.summary-strip { position: relative; display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); padding: 23px 12px; border-radius: 13px; border: 1px solid rgb(var(--color-primary-200) / .4); background: rgb(255 255 255 / .75); box-shadow: 0 2px 12px #14283c06; }
.summary-strip > div { padding: 0 20px; border-right: 1px solid rgb(var(--color-primary-200) / .6); min-width: 0; }
.summary-strip > div:last-of-type { border: 0; }
.summary-strip span, .summary-strip small { display: block; font-size: 10px; color: var(--muted); }
.summary-strip strong { display: block; margin: 7px 0; font-size: clamp(15px, 1.65vw, 23px); font-weight: 600; letter-spacing: -.04em; overflow-wrap: anywhere; font-variant-numeric: tabular-nums; }
.summary-strip small { color: var(--muted); }
.retry { position: absolute; right: 12px; bottom: 5px; color: var(--accent); text-decoration: underline; }
.two-columns { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 20px; }
.composition-layout { min-height: 210px; display: flex; align-items: center; gap: 25px; }
.composition-ring { width: 180px; height: 180px; flex-shrink: 0; position: relative; }
.ring-label { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; pointer-events: none; }
.ring-label strong { font-size: 22px; letter-spacing: -.04em; }
.ring-label span { color: var(--muted); font-size: 10px; margin-top: 3px; }
.composition-key { flex: 1; min-width: 0; display: grid; gap: 19px; }
.composition-key > div { display: flex; align-items: center; gap: 7px; font-size: 10px; }
.composition-key i { width: 7px; height: 7px; border-radius: 2px; flex-shrink: 0; }
.composition-key span { flex: 1; }
.composition-key small { color: var(--muted); width: 34px; text-align: right; }
.muted { font-size: 10px; color: var(--muted); }
.chart-actions { flex-wrap: wrap; }
input[type=checkbox] { accent-color: var(--accent); }
.table-scroll { overflow: auto; max-height: 370px; }
table { width: 100%; border-collapse: collapse; white-space: nowrap; font-size: 11px; }
th { color: var(--muted); font-weight: 500; background: rgb(var(--color-primary-100) / .45); }
th, td { padding: 11px 12px; text-align: right; border-bottom: 1px solid rgb(var(--color-primary-200) / .35); font-variant-numeric: tabular-nums; }
th:first-child, td:first-child { text-align: left; }
td:first-child, td:nth-child(2) { max-width: 240px; overflow: hidden; text-overflow: ellipsis; }
tbody tr:last-child td { border-bottom: 0; }
tbody tr:hover { background: rgb(var(--color-primary-100) / .35); }
.cache-layout { display: grid; grid-template-columns: 180px minmax(0, 1fr); gap: 30px; }
.cache-stats { display: grid; align-content: center; gap: 22px; }
.cache-stats span { font-size: 11px; color: var(--muted); }
.cache-stats strong, .cache-stats b { display: block; color: rgb(var(--color-primary-900)); margin-top: 5px; font-variant-numeric: tabular-nums; }
.cache-stats strong { font-size: 27px; }
.cache-stats b { font-size: 17px; }
.overview-list > div { display: flex; align-items: center; justify-content: space-between; padding: 13px 0; border-bottom: 1px solid rgb(var(--color-primary-200) / .35); }
.overview-list > div:last-child { border: 0; }
.overview-list dt { color: var(--muted); }
.overview-list dd { font-weight: 600; font-size: 15px; font-variant-numeric: tabular-nums; }
.all-records { display: flex; align-items: center; gap: 5px; font-size: 11px; color: var(--muted); }
button:focus-visible, select:focus-visible, a:focus-visible { outline: 2px solid var(--accent); outline-offset: 3px; }
.dark .analytics { --accent: rgb(var(--color-primary-100)); --muted: rgb(var(--color-primary-400)); color: rgb(var(--color-primary-100)); }
.dark .summary-strip { background: rgb(var(--color-primary-950) / .55); border-color: rgb(var(--color-primary-700) / .3); }
.dark .segmented { background: rgb(var(--color-primary-950)); border-color: rgb(var(--color-primary-700) / .4); }
.dark .segmented .selected { color: rgb(var(--color-primary-950)); }
.dark .summary-strip > div, .dark th, .dark td, .dark .overview-list > div { border-color: rgb(var(--color-primary-700) / .3); }
.dark th { background: rgb(var(--color-primary-800) / .25); }
.dark tbody tr:hover { background: rgb(var(--color-primary-800) / .3); }
.dark .cache-stats strong, .dark .cache-stats b { color: rgb(var(--color-primary-100)); }
@media(max-width: 1100px) { .composition-layout { gap: 12px; } .composition-ring { width: 145px; height: 145px; } .summary-strip > div { padding: 0 12px; } }
@media(max-width: 850px) { .two-columns { grid-template-columns: 1fr; } .composition-layout { gap: 28px; } .composition-ring { width: 180px; height: 180px; } }
@media(max-width: 600px) { .analytics { padding: 16px; gap: 14px; } .summary-strip { grid-template-columns: repeat(2,minmax(0,1fr)); gap: 20px 0; } .summary-strip > div { border: 0; } .summary-strip strong { font-size: 20px; } .composition-ring { width: 135px; height: 135px; } .composition-layout { gap: 12px; } .composition-key > div { flex-wrap: wrap; } .composition-key small { display: none; } .cache-layout { grid-template-columns: 1fr; gap: 20px; } .cache-stats { grid-template-columns: repeat(3,minmax(0,1fr)); gap: 10px; } .cache-stats strong { font-size: 20px; } .cache-stats b { font-size: 15px; } }
@media(prefers-reduced-motion: reduce) { .analytics { transition: none; transform: none; } }
</style>
