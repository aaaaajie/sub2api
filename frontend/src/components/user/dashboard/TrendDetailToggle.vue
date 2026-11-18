<template>
  <div class="trend-view-switch" role="group" :aria-label="t('dashboard.tokenConsumption.viewMode')">
    <span class="trend-view-switch__selection" :class="{ 'is-detailed': modelValue }" aria-hidden="true" />
    <button type="button" :aria-pressed="!modelValue" @click="$emit('update:modelValue', false)">
      {{ t('dashboard.tokenConsumption.overview') }}
    </button>
    <button type="button" :aria-pressed="modelValue" @click="$emit('update:modelValue', true)">
      {{ t('dashboard.tokenConsumption.breakdown') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
defineProps<{ modelValue: boolean }>()
defineEmits<{ 'update:modelValue': [value: boolean] }>()
const { t } = useI18n()
</script>

<style scoped>
.trend-view-switch {
  position: relative;
  display: inline-grid;
  grid-template-columns: repeat(2, 1fr);
  flex-shrink: 0;
  padding: 3px;
  border: 1px solid rgb(0 0 0 / 6%);
  border-radius: 999px;
  background: rgb(0 0 0 / 3%);
  isolation: isolate;
}
.trend-view-switch__selection {
  position: absolute;
  inset: 3px auto 3px 3px;
  width: calc((100% - 6px) / 2);
  border-radius: inherit;
  background: #fff;
  box-shadow: 0 1px 3px rgb(0 0 0 / 7%);
  transition: transform 220ms cubic-bezier(.22, 1, .36, 1);
  z-index: -1;
}
.trend-view-switch__selection.is-detailed { transform: translateX(100%); }
button {
  min-width: 54px;
  padding: 6px 12px;
  border-radius: inherit;
  color: rgb(107 114 128);
  font-size: 11px;
  font-weight: 500;
  line-height: 16px;
  transition: color 160ms ease;
}
button[aria-pressed='true'] { color: rgb(17 24 39); }
button:hover { color: rgb(17 24 39); }
button:focus-visible { outline: 2px solid currentColor; outline-offset: 3px; }
:global(.dark) .trend-view-switch { background: rgb(255 255 255 / 3%); border-color: rgb(255 255 255 / 7%); }
:global(.dark) .trend-view-switch__selection { background: rgb(255 255 255 / 9%); box-shadow: inset 0 1px 0 rgb(255 255 255 / 5%); }
:global(.dark) button { color: rgb(156 163 175); }
:global(.dark) button[aria-pressed='true'], :global(.dark) button:hover { color: rgb(243 244 246); }
@media (prefers-reduced-motion: reduce) {
  .trend-view-switch__selection, button { transition: none; }
}
</style>
