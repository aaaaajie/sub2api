<template>
  <span v-if="value == null">-</span>
  <span v-else-if="!schedule?.periods.length || aligned">{{ price(multiplier) }}</span>
  <span v-else class="inline-flex flex-col gap-0.5">
    <span>{{ price(1) }} <span class="font-sans text-[10px] opacity-70">{{ t('modelPlaza.table.officialOtherHours') }}</span></span>
    <span v-for="p in schedule.periods" :key="p.start_time" class="text-[11px]">
      {{ price(p.multiplier) }}
      <span class="font-sans text-[10px] opacity-70">
        {{ schedule.weekdays_only ? t('modelPlaza.table.timePricingWeekdays') : '' }}
        {{ clock(p.start_time) }}–{{ clock(p.end_time) }} ({{ schedule.timezone }})
      </span>
    </span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { PlazaTimePricing, PlazaTimePricingPeriod } from '@/api/modelPlaza'
import { formatScaled } from '@/utils/pricing'
import { clock, officialTimeMultiplier } from './officialPricing'

const props = defineProps<{
  value: number | null | undefined
  schedule?: PlazaTimePricing
  channelSchedule?: PlazaTimePricing
  period: PlazaTimePricingPeriod | null
}>()
const { t } = useI18n()
const aligned = computed(() => officialTimeMultiplier(props.schedule, props.channelSchedule, props.period) != null)
const multiplier = computed(() => officialTimeMultiplier(props.schedule, props.channelSchedule, props.period) ?? 1)
const price = (rate: number) => formatScaled(props.value == null ? null : props.value * rate, 1_000_000, 2)
</script>
