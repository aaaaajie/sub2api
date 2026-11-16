<template>
  <UsageInsightPanel :title="title" :loading="loading" :error="error" :empty="!rows.length" @retry="$emit('retry')">
    <template #actions>
      <div class="metrics">
        <button v-for="option in options" :key="option.key" :aria-pressed="metric === option.key" :class="{ selected: metric === option.key }" @click="metric = option.key">{{ option.label }}</button>
      </div>
    </template>
    <div class="distribution">
      <div v-for="row in sorted" :key="row.id" class="distribution-row" :style="{ '--series-color': `var(--chart-series-${distributionColorIndex(row.id)})` }">
        <div class="row-caption"><span :title="row.name">{{ row.name }}</span><strong>{{ format(row[metric]) }} <small>{{ total ? (row[metric] / total * 100).toFixed(1) : '0' }}%</small></strong></div>
        <progress :value="row[metric]" :max="total || 1" :aria-label="row.name" />
        <div class="row-detail"><span>{{ row.requests.toLocaleString() }} {{ t('dashboard.requests') }}</span><span>{{ formatTokensK(row.total_tokens) }} Token</span><span>${{ row.actual_cost.toFixed(4) }}</span></div>
      </div>
    </div>
  </UsageInsightPanel>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatTokensK } from '@/utils/format'
import UsageInsightPanel from './UsageInsightPanel.vue'
import { distributionColorIndex } from './chartPalette'
const props = defineProps<{ title: string; loading: boolean; error: boolean; rows: { id: string; name: string; total_tokens: number; requests: number; actual_cost: number }[] }>()
defineEmits<{ (event: 'retry'): void }>()
const { t } = useI18n()
const metric = ref<'total_tokens' | 'requests' | 'actual_cost'>('total_tokens')
const options = computed(() => [{ key: 'total_tokens' as const, label: 'Token' }, { key: 'requests' as const, label: t('dashboard.requests') }, { key: 'actual_cost' as const, label: t('dashboard.dataDashboard.spend') }])
const sorted = computed(() => [...props.rows].sort((a,b) => b[metric.value] - a[metric.value]))
const total = computed(() => props.rows.reduce((n, row) => n + row[metric.value], 0))
const format = (n: number) => metric.value === 'actual_cost' ? `$${n.toFixed(4)}` : metric.value === 'total_tokens' ? formatTokensK(n) : n.toLocaleString()
</script>
<style scoped>
.dark .distribution { --chart-series-0: #60a5fa; --chart-series-1: #a78bfa; --chart-series-2: #f472b6; --chart-series-3: #fbbf24; --chart-series-4: #22d3ee; --chart-series-5: #fb7185; }
.metrics { display: flex; gap: 3px; }
.metrics button { font-size: 10px; border-radius: 5px; padding: 5px 7px; color: var(--muted); }
.metrics button.selected { color: var(--accent); background: color-mix(in srgb, var(--accent) 12%, transparent); }
.distribution { --chart-series-0: #2563eb; --chart-series-1: #7c3aed; --chart-series-2: #db2777; --chart-series-3: #b45309; --chart-series-4: #0891b2; --chart-series-5: #e11d48; max-height: 300px; min-height: 200px; overflow-y: auto; display: grid; align-content: start; gap: 18px; }
.row-caption { display: flex; align-items: baseline; justify-content: space-between; gap: 14px; font-size: 12px; }
.row-caption > span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.row-caption strong { white-space: nowrap; font-size: 11px; font-variant-numeric: tabular-nums; }
small { color: var(--muted); font-size: 10px; margin-left: 8px; font-weight: 400; }
progress { appearance: none; width: 100%; height: 5px; border: 0; border-radius: 3px; overflow: hidden; background: color-mix(in srgb, var(--accent) 10%, transparent); margin: 8px 0 5px; }
progress::-webkit-progress-bar { background: color-mix(in srgb, var(--accent) 10%, transparent); }
progress::-webkit-progress-value { background: var(--series-color); }
progress::-moz-progress-bar { background: var(--series-color); }
.row-detail { display: flex; gap: 16px; color: var(--muted); font-size: 10px; font-variant-numeric: tabular-nums; }
</style>
