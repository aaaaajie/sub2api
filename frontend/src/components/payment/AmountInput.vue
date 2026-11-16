<template>
  <div class="recharge-amounts">
    <div>
      <div class="amount-heading"><span>{{ t('payment.quickAmounts') }}</span><span class="amount-currency">{{ normalizePaymentCurrency(currency) }}</span></div>
      <div class="amount-grid">
        <button v-for="amt in filteredAmounts" :key="amt" type="button" :class="['amount-option', { 'has-bonus': quoteFor(amt).percent > 0 }]" :aria-pressed="modelValue === amt"
          :data-testid="`quick-amount-${amt}`" @click="selectAmount(amt)">
          <span class="amount-value"><small>{{ currencySymbol(currency) }}</small>{{ amt }}</span>
          <span v-if="quoteFor(amt).percent > 0" class="amount-bonus" data-testid="quick-amount-bonus-badge">{{ badgeText(amt) }}</span>
          <span v-if="showSecondLine" class="amount-preview" data-testid="quick-amount-credited">{{ secondLine(amt) }}</span>
        </button>
      </div>
    </div>
    <div class="custom-amount">
      <label for="recharge-custom-amount">{{ t('payment.customAmount') }}</label>
      <div class="custom-amount-input">
        <span aria-hidden="true">{{ currencySymbol(currency) }}</span>
        <input id="recharge-custom-amount" type="text" inputmode="decimal" :value="customText" :placeholder="placeholderText" @input="handleInput" />
        <small>{{ normalizePaymentCurrency(currency) }}</small>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { RechargeBonusTier } from '@/types/payment'
import { formatRechargeBonusNumber, quoteRechargeBonus, type RechargeBonusMode } from '@/utils/rechargeBonus'
import { formatPaymentAmount, currencySymbol, normalizePaymentCurrency } from './currency'

const props = withDefaults(defineProps<{
  amounts?: number[]
  modelValue: number | null
  min?: number
  max?: number
  /** 充值优惠阶梯（按 min_amount 升序）；为空时不显示价签与第二行 */
  bonusTiers?: RechargeBonusTier[]
  /** 阶梯模式：bonus 赠金 / discount 折扣 */
  bonusMode?: RechargeBonusMode
  /** 充值倍率（1 支付币种 = multiplier USD），用于计算到账金额 */
  multiplier?: number
  /** 支付币种（折扣模式第二行实付金额的币种与精度） */
  currency?: string
}>(), {
  amounts: () => [10, 20, 50, 100, 200, 500, 1000, 2000, 5000],
  min: 0,
  max: 0,
  bonusTiers: () => [],
  bonusMode: 'bonus',
  multiplier: 1,
  currency: undefined,
})

const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

const { t } = useI18n()

const customText = ref('')

// 0 = no limit
const filteredAmounts = computed(() =>
  props.amounts.filter((a) => (props.min <= 0 || a >= props.min) && (props.max <= 0 || a <= props.max))
)

const showSecondLine = computed(() => props.bonusTiers.length > 0)

function currencyDigits(): number {
  if (!props.currency) return 2
  try {
    return new Intl.NumberFormat(undefined, { style: 'currency', currency: props.currency }).resolvedOptions().maximumFractionDigits ?? 2
  } catch {
    return 2
  }
}

function quoteFor(amt: number) {
  return quoteRechargeBonus(props.bonusTiers, amt, {
    multiplier: props.multiplier,
    mode: props.bonusMode,
    currencyDigits: currencyDigits(),
  })
}

// 价签文案：赠金「+20%」，折扣「20% OFF」
function badgeText(amt: number): string {
  const percent = formatRechargeBonusNumber(quoteFor(amt).percent)
  return t(props.bonusMode === 'discount' ? 'payment.checkout.badgeDiscount' : 'payment.checkout.badgeBonus', { percent })
}

function secondLine(amt: number): string {
  const quote = quoteFor(amt)
  if (props.bonusMode === 'discount') {
    return t('payment.rechargeBonus.payShort', { amount: formatPaymentAmount(quote.payBase, props.currency) })
  }
  return t('payment.rechargeBonus.creditedShort', { amount: '$' + quote.credited.toFixed(2) })
}

const placeholderText = computed(() => {
  if (props.min > 0 && props.max > 0) return `${props.min} - ${props.max}`
  if (props.min > 0) return `≥ ${props.min}`
  if (props.max > 0) return `≤ ${props.max}`
  return t('payment.enterAmount')
})

const AMOUNT_PATTERN = /^\d*(\.\d{0,2})?$/

function selectAmount(amt: number) {
  customText.value = String(amt)
  emit('update:modelValue', amt)
}

function handleInput(e: Event) {
  const input = e.target as HTMLInputElement
  const val = input.value
  if (!AMOUNT_PATTERN.test(val)) {
    input.value = customText.value
    return
  }
  customText.value = val
  if (val === '') {
    emit('update:modelValue', null)
    return
  }
  const num = parseFloat(val)
  if (!isNaN(num) && num > 0) {
    emit('update:modelValue', num)
  } else {
    emit('update:modelValue', null)
  }
}

watch(() => props.modelValue, (v) => {
  if (v !== null && String(v) !== customText.value) {
    customText.value = String(v)
  }
}, { immediate: true })
</script>

<style scoped>
.recharge-amounts { display: flex; flex-direction: column; gap: 22px; }
.amount-heading { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; font-size: 12px; font-weight: 500; @apply text-gray-700 dark:text-gray-300; }
.amount-heading .amount-currency { font-size: 10px; font-weight: 400; @apply text-gray-400 dark:text-gray-500; }
.amount-grid { padding-top: 6px; display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 18px 10px; }
.amount-option { position: relative; display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 68px; gap: 6px; padding: 14px 6px; border: 1px solid; border-radius: 13px; transition: background-color 160ms, border-color 160ms; @apply border-gray-200/80 text-gray-700 dark:border-white/10 dark:text-gray-300; }
.amount-option:hover { @apply border-gray-400 dark:border-white/25; }
.amount-option[aria-pressed='true'] { @apply border-gray-900 bg-gray-900 text-white dark:border-gray-200 dark:bg-gray-200 dark:text-gray-900; }
.amount-value { display: flex; align-items: baseline; gap: 3px; font-size: 18px; font-weight: 500; letter-spacing: -.025em; font-variant-numeric: tabular-nums; }
.amount-value small { font-size: 11px; opacity: .6; }
.amount-option.has-bonus { @apply border-[#42816d]/30 dark:border-[#9acdb8]/30; }
.amount-option.has-bonus[aria-pressed='true'] { @apply border-gray-900 dark:border-gray-200; }
.amount-bonus { position: absolute; top: -9px; right: 8px; border-radius: 6px; padding: 3px 7px; font-size: 10px; font-weight: 500; color: #42816d; background: #e5efe9; }
:global(.dark) .amount-bonus { color: #9acdb8; background: #1c3029; }
.amount-option[aria-pressed='true'] .amount-bonus { color: inherit; background: rgb(128 128 128 / 15%); }
.amount-preview { font-size: 10px; opacity: .65; line-height: 1.4; }
.custom-amount label { display: block; margin-bottom: 10px; font-size: 12px; @apply text-gray-500 dark:text-gray-400; }
.custom-amount-input { display: flex; align-items: center; gap: 10px; padding: 0 14px; border: 1px solid; border-radius: 12px; @apply border-gray-200 dark:border-white/10; }
.custom-amount-input:focus-within { @apply border-gray-500 dark:border-gray-400; }
.custom-amount-input > span, .custom-amount-input small { flex-shrink: 0; font-size: 11px; @apply text-gray-400 dark:text-gray-500; }
.custom-amount-input input { width: 100%; min-width: 0; border: 0; background: transparent; padding: 12px 0; outline: none; box-shadow: none; font-size: 13px; font-variant-numeric: tabular-nums; @apply text-gray-900 dark:text-gray-100; }
.amount-option:focus-visible { outline: 2px solid currentColor; outline-offset: 3px; }
@media (prefers-reduced-motion: reduce) { .amount-option { transition: none; } }
</style>
