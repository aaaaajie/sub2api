<template>
  <BaseDialog
    :show="true"
    :title="t('admin.userMail.title')"
    width="extra-wide"
    :close-on-escape="!busy"
    @close="close"
  >
    <div class="mb-5 flex items-center gap-1 border-b border-gray-200 dark:border-dark-600">
      <button
        v-for="tab in tabs"
        :key="tab"
        type="button"
        :aria-pressed="mode === tab"
        :disabled="!!busy"
        class="-mb-px border-b-2 px-3 pb-2.5 pt-1 text-sm font-medium transition-colors"
        :class="mode === tab ? 'border-gray-900 text-gray-900 dark:border-white dark:text-white' : 'border-transparent text-gray-500 dark:text-gray-400'"
        @click="switchMode(tab)"
      >
        {{ t(`admin.userMail.${tab}`) }}
      </button>
    </div>

    <div v-if="error" class="mb-4 rounded-lg bg-red-50 px-3 py-2.5 text-sm text-red-700 dark:bg-red-950/30 dark:text-red-300" role="alert">
      {{ error }}
    </div>

    <div v-if="mode === 'compose'" class="grid gap-5 lg:grid-cols-2">
      <fieldset :disabled="!!busy" class="min-w-0 space-y-4">
        <div class="flex min-w-0 items-center gap-2.5 rounded-lg border border-gray-200 px-3 py-2.5 dark:border-dark-600">
          <Icon name="users" size="sm" class="shrink-0 text-gray-500" />
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.userMail.to') }}</span>
          <span class="min-w-0 truncate text-sm font-medium text-gray-900 dark:text-white" :title="audienceLabel">{{ audienceLabel }}</span>
        </div>

        <div>
          <p class="input-label">{{ t('admin.userMail.template') }}</p>
          <div class="grid grid-cols-3 gap-2">
            <button
              v-for="template in templates"
              :key="template.id"
              type="button"
              :aria-pressed="draft.template === template.id"
              class="flex min-h-12 items-center justify-center gap-1.5 rounded-lg border px-2 py-2 text-center text-xs font-medium transition-colors sm:text-sm"
              :class="draft.template === template.id
                ? 'border-gray-900 bg-gray-900 text-white dark:border-gray-100 dark:bg-gray-100 dark:text-gray-900'
                : 'border-gray-200 bg-white text-gray-600 hover:border-gray-400 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-300 dark:hover:border-gray-500'"
              @click="draft.template = template.id"
            >
              <Icon :name="template.icon" size="sm" class="hidden shrink-0 sm:block" />
              {{ t(`admin.userMail.templates.${template.id}`) }}
            </button>
          </div>
        </div>

        <label class="block">
          <span class="input-label">{{ t('admin.userMail.subject') }}</span>
          <input v-model="draft.subject" type="text" class="input" maxlength="160" :placeholder="t('admin.userMail.subjectPlaceholder')" />
        </label>
        <label class="block">
          <span class="input-label">{{ t('admin.userMail.body') }}</span>
          <textarea v-model="draft.body" class="input min-h-40 resize-y leading-relaxed" rows="7" maxlength="12000" :placeholder="t('admin.userMail.bodyPlaceholder')" />
        </label>
        <details class="rounded-lg border border-gray-200 px-3 py-2.5 dark:border-dark-600">
          <summary class="cursor-pointer text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('admin.userMail.linkOptional') }}</summary>
          <div class="mt-3 grid gap-3 sm:grid-cols-[1fr_1.4fr]">
            <label>
              <span class="input-label">{{ t('admin.userMail.linkText') }}</span>
              <input v-model="draft.action_label" type="text" class="input" maxlength="40" />
            </label>
            <label>
              <span class="input-label">{{ t('admin.userMail.linkURL') }}</span>
              <input v-model="draft.action_url" type="url" class="input" placeholder="https://" maxlength="2048" />
            </label>
          </div>
        </details>
      </fieldset>

      <section class="min-w-0 overflow-hidden rounded-xl border border-gray-200 dark:border-dark-600" :aria-label="t('admin.userMail.preview')">
        <div class="flex min-h-12 flex-wrap items-center justify-between gap-2 border-b border-gray-200 px-3.5 py-2.5 dark:border-dark-600">
          <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.userMail.preview') }}</h4>
          <div v-if="preview" class="text-xs tabular-nums text-gray-600 dark:text-gray-300" aria-live="polite">
            {{ t('admin.userMail.recipientCount', { count: preview.user_ids.length }) }}
            <span v-if="preview.skipped">· {{ t('admin.userMail.skippedCount', { count: preview.skipped }) }}</span>
          </div>
        </div>
        <iframe
          v-if="preview"
          :srcdoc="preview.html"
          :title="t('admin.userMail.preview')"
          sandbox=""
          class="block h-[440px] w-full bg-[#eef0f2]"
          tabindex="-1"
        />
        <div v-else class="grid min-h-[230px] place-content-center px-6 py-10 text-center lg:min-h-[440px]">
          <span class="mx-auto mb-3 grid h-11 w-11 place-items-center rounded-xl border border-gray-200 bg-gray-50 text-gray-500 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-400">
            <Icon name="mail" size="md" />
          </span>
          <p class="max-w-64 text-sm leading-relaxed text-gray-500 dark:text-gray-400">{{ t('admin.userMail.previewHint') }}</p>
        </div>
      </section>
    </div>

    <div v-else class="grid min-w-0 gap-4 lg:grid-cols-[minmax(240px,0.8fr)_minmax(0,1.2fr)]">
      <section class="min-w-0 overflow-hidden rounded-xl border border-gray-200 dark:border-dark-600">
        <p v-if="historyLoading" class="p-4 text-sm text-gray-500">{{ t('admin.userMail.loading') }}</p>
        <p v-else-if="!jobs.length" class="p-4 text-sm text-gray-500">{{ t('admin.userMail.emptyHistory') }}</p>
        <ul v-else class="divide-y divide-gray-100 dark:divide-dark-600">
          <li v-for="job in jobs" :key="job.id" class="min-w-0">
            <button
              type="button"
              class="block w-full min-w-0 px-4 py-3.5 text-left transition-colors"
              :class="selectedID === job.id ? 'bg-gray-100 dark:bg-dark-700' : 'hover:bg-gray-50 dark:hover:bg-dark-800'"
              @click="openJob(job.id)"
            >
              <span class="flex items-center gap-2">
                <span class="min-w-0 truncate text-sm font-semibold text-gray-900 dark:text-white">{{ job.subject }}</span>
                <span v-if="job.is_test" class="shrink-0 rounded bg-gray-200 px-1.5 py-0.5 text-[10px] font-medium text-gray-700 dark:bg-dark-600 dark:text-gray-300">{{ t('admin.userMail.testTag') }}</span>
              </span>
              <span class="mt-1 block text-xs tabular-nums text-gray-500 dark:text-gray-400">{{ formatDateTime(job.created_at) }} · {{ t('admin.userMail.recipientCount', { count: total(job) }) }}</span>
            </button>
          </li>
        </ul>
        <div v-if="historyPages > 1" class="flex items-center justify-between border-t border-gray-200 px-3 py-2 text-xs dark:border-dark-600">
          <button type="button" class="btn btn-secondary btn-sm" :disabled="historyPage === 1 || historyLoading" @click="loadHistory(historyPage - 1)">{{ t('admin.userMail.previous') }}</button>
          <span class="tabular-nums text-gray-500">{{ historyPage }} / {{ historyPages }}</span>
          <button type="button" class="btn btn-secondary btn-sm" :disabled="historyPage >= historyPages || historyLoading" @click="loadHistory(historyPage + 1)">{{ t('admin.userMail.next') }}</button>
        </div>
      </section>

      <section v-if="detail" class="min-w-0 overflow-hidden rounded-xl border border-gray-200 dark:border-dark-600" aria-live="polite">
        <div class="border-b border-gray-200 px-4 py-4 dark:border-dark-600">
          <h4 class="break-words text-base font-semibold text-gray-900 dark:text-white">{{ detail.subject }}</h4>
          <div class="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-xs text-gray-600 dark:text-gray-300">
            <span v-for="key in summaryKeys" :key="key" class="tabular-nums">{{ t(`admin.userMail.status.${key}`) }} {{ detail[key] }}</span>
          </div>
          <progress class="mt-3 h-1.5 w-full accent-gray-800 dark:accent-gray-200" :value="total(detail) - detail.pending" :max="total(detail)" :aria-label="t('admin.userMail.progress')" />
        </div>
        <ul class="max-h-80 divide-y divide-gray-100 overflow-y-auto dark:divide-dark-600">
          <li v-for="recipient in detail.recipients" :key="recipient.user_id" class="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1 px-4 py-2.5 text-xs">
            <span class="min-w-0 flex-1 truncate font-medium text-gray-900 dark:text-gray-100" :title="recipient.email">{{ recipient.email }}</span>
            <span class="rounded px-1.5 py-0.5 font-medium" :class="statusClass(recipient.status)">
              {{ t(`admin.userMail.status.${recipient.status}`) }}
            </span>
            <button v-if="recipient.status === 'failed'" type="button" class="btn btn-secondary btn-sm" :disabled="retrying !== null" @click="retryRecipient(recipient.user_id)">{{ t('admin.userMail.retry') }}</button>
            <p v-if="recipient.reason" class="w-full text-gray-500 dark:text-gray-400">{{ t(`admin.userMail.reasons.${recipient.reason}`) }}</p>
          </li>
        </ul>
      </section>
      <p v-else class="px-1 py-4 text-sm text-gray-500">{{ t(selectedID ? 'admin.userMail.loading' : 'admin.userMail.chooseJob') }}</p>
    </div>

    <template #footer>
      <div v-if="mode === 'compose'" class="flex w-full flex-wrap justify-end gap-2">
        <button type="button" class="btn btn-secondary" :disabled="!!busy || !canPreview" @click="renderPreview">
          {{ t(busy === 'preview' ? 'admin.userMail.loading' : 'admin.userMail.previewAction') }}
        </button>
        <button v-if="preview" type="button" class="btn btn-secondary" :disabled="!!busy || testQueued" @click="submit(true)">{{ t('admin.userMail.sendTest') }}</button>
        <button v-if="preview" type="button" class="btn btn-primary" :disabled="!!busy || !preview.user_ids.length" @click="submit(false)">
          {{ t('admin.userMail.sendTo', { count: preview.user_ids.length }) }}
        </button>
      </div>
      <button v-else type="button" class="btn btn-secondary" @click="close">{{ t('common.close') }}</button>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'
import { getUserMail, listUserMail, previewUserMail, queueUserMail, retryUserMailRecipient } from '@/api/admin/userMail'
import type { UserMailAudience, UserMailDraft, UserMailJob, UserMailJobDetail, UserMailPreview, UserMailStatus } from '@/api/admin/userMail'

const props = defineProps<{ audience: UserMailAudience; audienceLabel: string; startInHistory?: boolean }>()
const emit = defineEmits<{ close: [] }>()
const { t, locale } = useI18n()
const tabs = ['compose', 'history'] as const
const templates = [{ id: 'maintenance', icon: 'clock' }, { id: 'service', icon: 'sparkles' }, { id: 'account', icon: 'user' }] as const
const summaryKeys = ['pending', 'accepted', 'failed', 'skipped', 'uncertain'] as const
const mode = ref<(typeof tabs)[number]>('compose')
const draft = reactive<UserMailDraft>({ template: 'service', subject: '', body: '', action_label: '', action_url: '', locale: locale.value })
const preview = ref<UserMailPreview>()
const jobID = ref('')
const testID = ref('')
const testQueued = ref(false)
const busy = ref<'preview' | 'send' | 'test' | null>(null)
const error = ref('')
const canPreview = computed(() => !!draft.subject.trim() && !!draft.body.trim())
const jobs = ref<UserMailJob[]>([])
const historyLoading = ref(false)
const historyPage = ref(1)
const historyPages = ref(1)
const selectedID = ref('')
const detail = ref<UserMailJobDetail>()
const retrying = ref<number | null>(null)
let pollTimer: ReturnType<typeof setTimeout> | undefined
let detailVersion = 0

watch(draft, () => {
  preview.value = undefined
  jobID.value = ''
  testID.value = ''
  testQueued.value = false
  error.value = ''
})

function message(cause: unknown) {
  const reason = (cause as { reason?: string } | undefined)?.reason
  const translated = reason ? t(`admin.userMail.errors.${reason}`) : ''
  return translated && translated !== `admin.userMail.errors.${reason}` ? translated : t('admin.userMail.error')
}

const total = (job: UserMailJob) => summaryKeys.reduce((sum, key) => sum + job[key], 0)

function statusClass(status: UserMailStatus) {
  if (status === 'accepted') return 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300'
  if (status === 'failed') return 'bg-red-50 text-red-700 dark:bg-red-950/40 dark:text-red-300'
  if (status === 'uncertain') return 'bg-amber-50 text-amber-800 dark:bg-amber-950/40 dark:text-amber-200'
  return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'
}

function clearPoll() {
  clearTimeout(pollTimer)
  detailVersion++
}

function close() {
  if (!busy.value) emit('close')
}

async function switchMode(tab: (typeof tabs)[number]) {
  if (mode.value === tab) return
  mode.value = tab
  error.value = ''
  if (tab === 'history') {
    await loadHistory(1)
    if (selectedID.value && mode.value === tab) void openJob(selectedID.value)
  } else {
    clearPoll()
  }
}

async function renderPreview() {
  if (busy.value || !canPreview.value) return
  busy.value = 'preview'
  error.value = ''
  try {
    preview.value = await previewUserMail({ ...draft }, props.audience)
  } catch (cause) {
    error.value = message(cause)
  } finally {
    busy.value = null
  }
}

async function submit(isTest: boolean) {
  if (busy.value || !preview.value || (isTest ? testQueued.value : !preview.value.user_ids.length)) return
  busy.value = isTest ? 'test' : 'send'
  error.value = ''
  try {
    const requestID = isTest ? testID : jobID
    requestID.value ||= crypto.randomUUID()
    const job = await queueUserMail(requestID.value, { ...draft }, isTest ? [] : preview.value.user_ids, isTest)
    if (isTest) testQueued.value = true
    else {
      draft.subject = ''
      draft.body = ''
      draft.action_label = ''
      draft.action_url = ''
    }
    mode.value = 'history'
    await loadHistory(1)
    await openJob(job.id)
  } catch (cause) {
    error.value = message(cause)
  } finally {
    busy.value = null
  }
}

async function loadHistory(page: number) {
  historyLoading.value = true
  try {
    const result = await listUserMail(page)
    jobs.value = result.items
    historyPage.value = result.page
    historyPages.value = result.pages
  } catch (cause) {
    error.value = message(cause)
  } finally {
    historyLoading.value = false
  }
}

async function openJob(id: string) {
  clearPoll()
  if (selectedID.value !== id) detail.value = undefined
  selectedID.value = id
  const version = detailVersion
  try {
    const result = await getUserMail(id)
    if (version !== detailVersion || mode.value !== 'history') return
    detail.value = result
    if (result.pending) pollTimer = setTimeout(() => { void openJob(id) }, 2500)
  } catch (cause) {
    if (version === detailVersion) error.value = message(cause)
  }
}

async function retryRecipient(userId: number) {
  if (!detail.value || retrying.value !== null) return
  const id = detail.value.id
  retrying.value = userId
  error.value = ''
  try {
    await retryUserMailRecipient(id, userId)
    if (selectedID.value === id && mode.value === 'history') await openJob(id)
  } catch (cause) {
    error.value = message(cause)
  } finally {
    retrying.value = null
  }
}

if (props.startInHistory) void switchMode('history')
onBeforeUnmount(clearPoll)
</script>
