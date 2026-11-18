<template>
  <section class="insight-panel" :aria-busy="loading">
    <header><h3>{{ title }}</h3><slot name="actions" /></header>
    <div v-if="error" class="panel-state" role="alert"><p>{{ t('dashboard.dataDashboard.error') }}</p><button type="button" @click="$emit('retry')">{{ t('dashboard.dataDashboard.retry') }}</button></div>
    <div v-else-if="loading" class="panel-state" role="status">{{ t('common.loading') }}</div>
    <div v-else-if="empty" class="panel-state">{{ t('dashboard.dataDashboard.empty') }}</div>
    <slot v-else />
  </section>
</template>
<script setup lang="ts">
import { useI18n } from 'vue-i18n'
defineProps<{ title: string; loading?: boolean; error?: boolean; empty?: boolean }>()
defineEmits<{ (event: 'retry'): void }>()
const { t } = useI18n()
</script>
<style scoped>
.insight-panel { min-width: 0; padding: 20px 22px; border: 1px solid rgb(var(--color-primary-200) / .42); border-radius: 14px; background: rgb(255 255 255 / .75); box-shadow: 0 2px 12px rgb(20 40 60 / .025); }
header { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; margin-bottom: 22px; }
h3 { font-size: 13px; font-weight: 600; }
.panel-state { min-height: 180px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; color: rgb(var(--color-primary-500)); font-size: 12px; }
.panel-state button { text-decoration: underline; text-underline-offset: 4px; }
.dark .insight-panel { background: rgb(var(--color-primary-950) / .55); border-color: rgb(var(--color-primary-700) / .3); box-shadow: none; }
@media(max-width: 600px) { .insight-panel { padding: 16px; } }
</style>
