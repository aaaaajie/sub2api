<template>
  <div class="order-list" :aria-busy="loading">
    <div v-if="loading && !orders.length" class="order-skeleton" aria-live="polite">
      <span class="sr-only">{{ t('common.loading') }}</span>
      <div v-for="n in 4" :key="n"><span /><span /><span /></div>
    </div>
    <div v-else-if="!orders.length" class="order-empty">
      <Icon name="creditCard" size="lg" />
      <p>{{ t('payment.orders.empty') }}</p>
    </div>
    <div v-else :class="{ 'is-updating': loading }">
      <details v-for="order in orders" :key="order.id" class="order-entry">
        <summary>
          <span class="order-mark"><Icon :name="order.order_type === 'subscription' ? 'calendar' : 'creditCard'" size="md" /></span>
          <span class="order-identity"><strong>{{ t(order.order_type === 'subscription' ? 'payment.tabSubscribe' : 'payment.tabTopUp') }}</strong><time :datetime="order.created_at">{{ formatDate(order.created_at) }}</time></span>
          <span class="order-payment"><strong>{{ formatPaymentAmount(order.pay_amount, order.currency) }} <small>{{ normalizePaymentCurrency(order.currency) }}</small></strong><span>{{ t('payment.methods.' + order.payment_type, order.payment_type) }}</span></span>
          <OrderStatusBadge :status="order.status" class="order-status" />
          <Icon name="chevronDown" size="sm" class="order-chevron" />
        </summary>
        <div class="order-details">
          <dl>
            <div><dt>{{ t('payment.orders.orderNo') }}</dt><dd>{{ order.out_trade_no || '#' + order.id }}</dd></div>
            <div><dt>{{ t('payment.orders.payAmount') }}</dt><dd>{{ formatPaymentAmount(order.pay_amount, order.currency) }} {{ normalizePaymentCurrency(order.currency) }}</dd></div>
            <template v-if="order.order_type === 'balance'">
              <div><dt>{{ t('payment.orders.creditedAmount') }}</dt><dd>{{ formatPaymentAmount(order.amount, 'USD') }} USD</dd></div>
              <div v-if="(order.bonus_amount ?? 0) > 0"><dt>{{ t('payment.orders.bonusAmount') }}</dt><dd class="order-benefit">{{ t('payment.orders.bonusIncluded', { amount: formatPaymentAmount((order.bonus_amount ?? 0), 'USD') }) }}</dd></div>
            </template>
            <div v-if="order.fee_rate > 0"><dt>{{ t('payment.orders.fee') }}</dt><dd>{{ order.fee_rate }}%</dd></div>
          </dl>
          <slot name="actions" :row="order" />
        </div>
      </details>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { PaymentOrder } from '@/types/payment'
import Icon from '@/components/icons/Icon.vue'
import OrderStatusBadge from './OrderStatusBadge.vue'
import { formatPaymentAmount, normalizePaymentCurrency } from './currency'
defineProps<{ orders: PaymentOrder[]; loading: boolean }>()
const { t } = useI18n()
function formatDate(value: string) { return new Date(value).toLocaleString() }
</script>

<style scoped>
.order-list { overflow: hidden; border: 1px solid; border-radius: 20px; @apply border-gray-200/70 bg-white/60 dark:border-white/[0.07] dark:bg-dark-950; }
.order-entry + .order-entry { border-top: 1px solid; @apply border-gray-200/60 dark:border-white/[0.06]; }
summary { display: grid; grid-template-columns: 38px minmax(0, 1fr) minmax(120px, .6fr) 96px 16px; align-items: center; gap: 18px; padding: 24px; cursor: pointer; list-style: none; transition: background .2s; }
summary::-webkit-details-marker { display: none; }
summary:hover { @apply bg-gray-50 dark:bg-white/[0.025]; }
summary:focus-visible { outline: 2px solid currentColor; outline-offset: -3px; }
.order-mark { display: grid; place-items: center; width: 38px; height: 38px; border-radius: 12px; @apply bg-gray-100 text-gray-500 dark:bg-white/[0.05] dark:text-gray-400; }
.order-identity, .order-payment { display: flex; flex-direction: column; gap: 6px; min-width: 0; }
.order-identity strong { font-size: 13px; font-weight: 500; @apply text-gray-900 dark:text-gray-100; }
.order-identity time, .order-payment > span { font-size: 11px; @apply text-gray-500 dark:text-gray-400; }
.order-payment strong { font-size: 17px; font-weight: 500; letter-spacing: -.025em; font-variant-numeric: tabular-nums; @apply text-gray-900 dark:text-gray-100; }
.order-payment small { font-size: 9px; font-weight: 400; margin-left: 3px; @apply text-gray-400; }
.order-status { justify-self: start; white-space: nowrap; font-size: 10px; }
.order-chevron { transition: transform .2s; @apply text-gray-400; }
[open] .order-chevron { transform: rotate(180deg); }
.order-details { padding: 0 24px 24px 80px; }
dl { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px 28px; padding: 20px 0; border-top: 1px dashed; @apply border-gray-200 dark:border-white/10; }
dt { font-size: 10px; margin-bottom: 5px; @apply text-gray-500 dark:text-gray-400; }
dd { font-size: 12px; font-variant-numeric: tabular-nums; overflow-wrap: anywhere; @apply text-gray-700 dark:text-gray-300; }
dd.order-benefit { @apply text-[#42816d] dark:text-[#9acdb8]; }
.order-empty { display: flex; flex-direction: column; align-items: center; gap: 14px; padding: 80px 24px; font-size: 12px; @apply text-gray-400 dark:text-gray-500; }
.is-updating { opacity: .55; transition: opacity .2s; }
.order-skeleton > div { display: flex; align-items: center; gap: 24px; height: 100px; padding: 24px; }
.order-skeleton span:not(.sr-only) { height: 12px; width: 25%; border-radius: 6px; @apply bg-gray-100 dark:bg-white/[0.05]; }
.order-skeleton span:first-child { width: 38px; height: 38px; }
@media (max-width: 640px) {
  summary { grid-template-columns: 32px minmax(0, 1fr) auto; gap: 12px; padding: 18px; }
  .order-mark { width: 32px; height: 32px; }
  .order-payment { text-align: right; }
  summary .order-status { grid-column: 2; justify-self: start; }
  .order-chevron { display: none; }
  .order-details { padding: 0 18px 20px; }
  dl { grid-template-columns: 1fr; }
}
@media (prefers-reduced-motion: reduce) { summary, .order-chevron, .is-updating { transition: none; } }
</style>
