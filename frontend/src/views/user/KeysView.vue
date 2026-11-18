<template>
  <AppLayout>
    <TablePageLayout class="keys-page-layout">
      <template #filters>
        <div class="flex flex-col gap-3">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div class="api-key-filter-shell flex h-12 w-full min-w-0 flex-none items-center gap-1 rounded-xl border border-gray-200 bg-white px-1.5 transition-colors focus-within:border-primary-500 focus-within:ring-2 focus-within:ring-primary-500/20 dark:border-dark-700 dark:bg-dark-950 sm:w-80">
              <Select
                :model-value="filterGroupId"
                class="key-filter-select"
                dropdown-class="keys-filter-dropdown"
                :options="groupFilterOptions"
                :aria-label="t('keys.group')"
                :title="selectedGroupFilterLabel"
                @update:model-value="onGroupFilterChange"
              >
                <template #selected="{ option }">
                  <span
                    v-if="option?.kind === 'group'"
                    class="flex h-7 w-7 items-center justify-center rounded-lg"
                    :class="platformBadgeLightClass((option as unknown as GroupOption).platform)"
                  >
                    <PlatformIcon :platform="(option as unknown as GroupOption).platform" size="sm" />
                  </span>
                  <span
                    v-else
                    class="flex h-7 w-7 items-center justify-center rounded-lg bg-gray-100 text-gray-500 dark:bg-dark-950 dark:text-gray-300"
                  >
                    <Icon :name="option?.kind === 'none' ? 'xCircle' : 'grid'" size="sm" />
                  </span>
                </template>
                <template #option="{ option, selected }">
                  <GroupOptionItem
                    v-if="option.kind === 'group'"
                    :name="(option as unknown as GroupOption).label"
                    :platform="(option as unknown as GroupOption).platform"
                    :subscription-type="(option as unknown as GroupOption).subscriptionType"
                    :rate-multiplier="(option as unknown as GroupOption).rate"
                    :user-rate-multiplier="(option as unknown as GroupOption).userRate"
                    :peak-rate-enabled="(option as unknown as GroupOption).peakRateEnabled"
                    :peak-start="(option as unknown as GroupOption).peakStart"
                    :peak-end="(option as unknown as GroupOption).peakEnd"
                    :peak-rate-multiplier="(option as unknown as GroupOption).peakRateMultiplier"
                    :description="(option as unknown as GroupOption).description"
                    :selected="selected"
                  />
                  <div v-else class="flex w-full items-center justify-between gap-3">
                    <span class="flex min-w-0 items-center gap-2.5">
                      <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-500 dark:bg-dark-950 dark:text-gray-300">
                        <Icon :name="option.kind === 'none' ? 'xCircle' : 'grid'" size="sm" />
                      </span>
                      <span class="truncate">{{ option.label }}</span>
                    </span>
                    <Icon v-if="selected" name="check" size="sm" class="shrink-0 text-primary-500" :stroke-width="2" />
                  </div>
                </template>
              </Select>

              <Select
                :model-value="filterStatus"
                class="key-filter-select"
                dropdown-class="keys-filter-dropdown"
                :options="statusFilterOptions"
                :aria-label="t('common.status')"
                :title="selectedStatusFilterLabel"
                @update:model-value="onStatusFilterChange"
              >
                <template #selected="{ option }">
                  <span
                    v-if="option"
                    :class="['flex h-7 w-7 items-center justify-center rounded-lg', option.badgeClass]"
                  >
                    <Icon :name="(option as StatusFilterOption).icon" size="xs" />
                  </span>
                </template>
                <template #option="{ option, selected }">
                  <div class="flex w-full items-center justify-between gap-3">
                    <span :class="['inline-flex min-w-0 items-center gap-1.5 rounded-full px-2 py-1 text-xs font-medium', option.badgeClass]">
                      <Icon :name="(option as StatusFilterOption).icon" size="xs" />
                      <span class="truncate">{{ option.label }}</span>
                    </span>
                    <Icon v-if="selected" name="check" size="sm" class="shrink-0 text-primary-500" :stroke-width="2" />
                  </div>
                </template>
              </Select>

              <span class="mx-1 h-6 w-px shrink-0 bg-gray-200 dark:bg-dark-600" aria-hidden="true" />

              <SearchInput
                v-model="filterSearch"
                :placeholder="t('keys.searchPlaceholder')"
                class="key-filter-search min-w-0 flex-1"
                @search="onFilterChange"
              />
            </div>

            <div
              v-if="selectedIds.length"
              class="flex h-11 shrink-0 items-center gap-1 rounded-xl border border-gray-200 bg-white px-2 dark:border-dark-700 dark:bg-dark-950"
            >
              <span class="whitespace-nowrap px-1 text-xs font-medium text-gray-500 dark:text-gray-400">
                {{ t('keys.bulkEdit.selectedCount', { count: selectedIds.length }) }}
              </span>
              <button
                class="inline-flex h-8 items-center gap-1.5 rounded-lg px-2 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-dark-800"
                :disabled="loading"
                data-test="bulk-edit-keys"
                :title="t('keys.bulkEdit.title')"
                :aria-label="t('keys.bulkEdit.title')"
                @click="showBulkEditModal = true"
              >
                <Icon name="edit" size="sm" />
                <span class="hidden sm:inline">{{ t('keys.bulkEdit.title') }}</span>
              </button>
              <button
                class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dark-800"
                :title="t('keys.bulkEdit.clearSelection')"
                :aria-label="t('keys.bulkEdit.clearSelection')"
                @click="selectedIds = []"
              >
                <Icon name="x" size="sm" />
              </button>
            </div>

            <div class="ml-auto flex shrink-0 items-center gap-2">
              <button
                v-if="apiKeys.length"
                @click="togglePageSelection"
                class="btn btn-secondary keys-toolbar-button px-2 md:px-3"
                :title="allPageSelected ? t('keys.bulkEdit.clearSelection') : t('common.selectAll')"
              >
                <Icon :name="allPageSelected ? 'checkCircle' : 'check'" size="sm" class="md:mr-1.5" />
                <span class="hidden md:inline">{{ allPageSelected ? t('keys.bulkEdit.clearSelection') : t('common.selectAll') }}</span>
              </button>
              <button @click="showCreateModal = true" class="btn btn-primary" data-tour="keys-create-btn">
                <Icon name="plus" size="md" class="mr-2" />
                {{ t('keys.createKey') }}
              </button>
            </div>
          </div>
          <EndpointPopover
            v-if="publicSettings?.api_base_url || (publicSettings?.custom_endpoints?.length ?? 0) > 0"
            :api-base-url="publicSettings?.api_base_url || ''"
            :custom-endpoints="publicSettings?.custom_endpoints || []"
          />
        </div>
      </template>

      <template #table>
        <div v-if="showInitialSkeleton" class="key-card-grid" aria-label="Loading API keys">
          <div
            v-for="placeholder in 3"
            :key="placeholder"
            class="key-skeleton-card rounded-xl border border-gray-200/90 bg-white p-4 dark:border-dark-800 dark:bg-dark-950 sm:p-[1.125rem]"
            :style="{ animationDelay: `${(placeholder - 1) * 140}ms` }"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="key-skeleton-block h-7 w-7 rounded-lg bg-gray-200 dark:bg-dark-800" />
                <div class="key-skeleton-block h-3 w-24 rounded bg-gray-200 dark:bg-dark-700" />
              </div>
              <div class="key-skeleton-block h-5 w-14 rounded-full bg-gray-200 dark:bg-dark-800" />
            </div>
            <div class="key-skeleton-block mt-2.5 flex h-8 items-center rounded-lg bg-gray-100 px-2.5 dark:bg-dark-800/70">
              <div class="h-2.5 w-32 rounded bg-gray-200 dark:bg-dark-700" />
            </div>
            <div class="mt-2.5 flex items-center justify-between gap-2">
              <div class="key-skeleton-block h-5 w-20 rounded-md bg-gray-100 dark:bg-dark-800" />
              <div class="key-skeleton-block h-2.5 w-20 rounded bg-gray-100 dark:bg-dark-800" />
            </div>
            <div class="mt-3 flex items-center justify-end gap-4">
              <div class="flex items-center gap-1.5">
                <div class="key-skeleton-block h-2 w-7 rounded bg-gray-100 dark:bg-dark-800" />
                <div class="key-skeleton-block h-2.5 w-10 rounded bg-gray-200 dark:bg-dark-700" />
              </div>
              <div class="flex items-center gap-1.5">
                <div class="key-skeleton-block h-2 w-10 rounded bg-gray-100 dark:bg-dark-800" />
                <div class="key-skeleton-block h-2.5 w-10 rounded bg-gray-200 dark:bg-dark-700" />
              </div>
            </div>
            <div class="mt-3 flex items-center justify-between gap-2 border-t border-gray-100 pt-2 dark:border-dark-800">
              <div class="key-skeleton-block h-8 w-24 rounded-lg bg-gray-200 dark:bg-dark-800" />
              <div class="flex items-center gap-1.5">
                <div class="key-skeleton-block h-7 w-7 rounded-lg bg-gray-100 dark:bg-dark-800" />
                <div class="key-skeleton-block h-7 w-7 rounded-lg bg-gray-100 dark:bg-dark-800" />
                <div class="key-skeleton-block h-7 w-7 rounded-lg bg-gray-100 dark:bg-dark-800" />
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="loading && !apiKeys.length" class="key-card-grid min-h-[14rem]" aria-hidden="true" />

        <EmptyState
          v-else-if="!loading && !apiKeys.length"
          :title="t('keys.noKeysYet')"
          :description="t('keys.createFirstKey')"
          :action-text="t('keys.createKey')"
          @action="showCreateModal = true"
        />

        <div v-else class="key-card-grid">
          <article
            v-for="row in apiKeys"
            :key="row.id"
            role="group"
            tabindex="0"
            :aria-label="`${t('keys.bulkEdit.selectKey', { name: row.name })}${selectedIds.includes(row.id) ? ` ${t('common.selectedCount', { count: 1 })}` : ''}`"
            @click="handleKeyCardClick(row, $event)"
            @keydown="handleKeyCardKeydown(row, $event)"
            :class="[
              'flex min-w-0 cursor-pointer flex-col rounded-xl border p-4 shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500/50 sm:p-[1.125rem]',
              selectedIds.includes(row.id)
                ? 'border-primary-500/70 bg-primary-50/40 ring-1 ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:ring-dark-600/30'
                : 'border-gray-200/90 bg-white hover:border-gray-300 dark:border-dark-800 dark:bg-dark-950 dark:hover:border-dark-700'
            ]"
          >
            <div class="flex min-w-0 items-center justify-between gap-2">
              <div class="flex min-w-0 items-center gap-2">
                <span
                  :class="[
                    'flex h-7 w-7 shrink-0 items-center justify-center rounded-lg border',
                    row.group?.platform ? platformBadgeLightClass(groupDisplayPlatform(row.group)) : 'border-gray-200 bg-gray-100 text-gray-500 dark:border-dark-700 dark:bg-dark-800 dark:text-gray-300'
                  ]"
                >
                  <PlatformIcon v-if="row.group?.platform" :platform="groupDisplayPlatform(row.group) || undefined" size="md" />
                  <Icon v-else name="key" size="sm" />
                </span>
                <h2 class="truncate text-sm font-semibold text-gray-900 dark:text-white">{{ row.name }}</h2>
                <Icon
                  v-if="row.ip_whitelist?.length > 0 || row.ip_blacklist?.length > 0"
                  name="shield"
                  size="xs"
                  class="shrink-0 text-gray-500 dark:text-gray-400"
                  :title="t('keys.ipRestrictionEnabled')"
                />
              </div>
              <span
                :class="[
                  'inline-flex shrink-0 items-center gap-1.5 rounded-full px-2 py-0.5 text-[10px] font-medium',
                  keyStatusClass(row.status)
                ]"
              >
                {{ t('keys.status.' + row.status) }}
                <button
                  v-if="row.status === 'active' || row.status === 'inactive'"
                  type="button"
                  role="switch"
                  :aria-checked="row.status === 'active'"
                  :aria-label="row.status === 'active' ? t('keys.disable') : t('keys.enable')"
                  :title="row.status === 'active' ? t('keys.disable') : t('keys.enable')"
                  @click.stop="toggleKeyStatus(row)"
                  :class="[
                    'relative inline-flex h-3.5 w-6 shrink-0 items-center rounded-full transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-1 dark:focus-visible:ring-offset-dark-900',
                    row.status === 'active' ? 'bg-emerald-500/70' : 'bg-gray-400/70 dark:bg-dark-500'
                  ]"
                >
                  <span
                    :class="[
                      'h-2.5 w-2.5 rounded-full bg-white shadow-sm transition-transform',
                      row.status === 'active' ? 'translate-x-3' : 'translate-x-0.5'
                    ]"
                  />
                </button>
              </span>
            </div>

            <div class="mt-2.5 flex min-w-0 items-center gap-2 rounded-lg bg-gray-50 px-2.5 py-1.5 dark:bg-dark-800/70">
              <div class="min-w-0 flex-1">
                <code class="block truncate font-mono text-[11px] text-gray-700 dark:text-gray-200" :title="t('keys.apiKey')">
                  {{ maskApiKey(row.key) }}
                </code>
              </div>
              <button
                type="button"
                @click="copyToClipboard(row.key, row.id)"
                :aria-label="copiedKeyId === row.id ? t('keys.copied') : t('keys.copyToClipboard')"
                :title="copiedKeyId === row.id ? t('keys.copied') : t('keys.copyToClipboard')"
                :class="[
                  'flex h-6 w-6 shrink-0 items-center justify-center rounded-md transition-colors hover:bg-white dark:hover:bg-dark-700',
                  copiedKeyId === row.id
                    ? 'text-emerald-600 dark:text-emerald-400'
                    : 'text-gray-500 dark:text-gray-400'
                ]"
              >
                <Icon v-if="copiedKeyId === row.id" name="check" size="sm" :stroke-width="2" />
                <Icon v-else name="clipboard" size="sm" />
              </button>
            </div>

            <div class="mt-2.5 flex min-w-0 items-center justify-between gap-2">
              <button
                :ref="(el) => setGroupButtonRef(row.id, el)"
                type="button"
                @click="openGroupSelector(row)"
                class="group/dropdown flex min-w-0 items-center gap-1.5 rounded-lg text-left"
                :title="t('keys.clickToChangeGroup')"
              >
                <GroupBadge
                  v-if="row.group"
                  :name="row.group.name"
                  :platform="groupDisplayPlatform(row.group) || undefined"
                  :subscription-type="row.group.subscription_type"
                  :rate-multiplier="row.group.rate_multiplier"
                  :user-rate-multiplier="userGroupRates[row.group.id]"
                  :peak-rate-enabled="row.group.peak_rate_enabled"
                  :peak-start="row.group.peak_start"
                  :peak-end="row.group.peak_end"
                  :peak-rate-multiplier="row.group.peak_rate_multiplier"
                  :show-rate="false"
                />
                <span v-else class="text-xs text-gray-400 dark:text-dark-500">{{ t('keys.noGroup') }}</span>
                <Icon
                  name="chevronDown"
                  size="xs"
                  class="shrink-0 text-gray-400 transition-colors group-hover/dropdown:text-gray-700 dark:group-hover/dropdown:text-gray-200"
                />
              </button>
              <span class="truncate text-[9px] tabular-nums text-gray-400 dark:text-gray-500">{{ formatDateTime(row.created_at) }}</span>
            </div>

            <div class="mt-2.5 flex flex-wrap items-center justify-end gap-x-2 gap-y-1 text-[9px] text-gray-400 dark:text-gray-500">
              <div class="flex items-center gap-1">
                <span>{{ t('keys.today') }}</span>
                <span
                  v-if="usageLoading && !usageStats[row.id]"
                  data-test="usage-pending"
                  class="min-w-[2.75rem] tabular-nums text-gray-400 dark:text-gray-500"
                  :aria-label="t('common.loading')"
                >$—</span>
                <span v-else class="min-w-[2.75rem] tabular-nums text-gray-500 dark:text-gray-400"><span>$</span>{{ (usageStats[row.id]?.today_actual_cost ?? 0).toFixed(4) }}</span>
              </div>
              <div class="flex items-center gap-1">
                <span>{{ t('keys.total') }}</span>
                <span
                  v-if="usageLoading && !usageStats[row.id]"
                  class="min-w-[2.75rem] tabular-nums text-gray-400 dark:text-gray-500"
                  :aria-label="t('common.loading')"
                >$—</span>
                <span v-else class="min-w-[2.75rem] tabular-nums text-gray-500 dark:text-gray-400"><span>$</span>{{ (usageStats[row.id]?.total_actual_cost ?? 0).toFixed(4) }}</span>
              </div>
            </div>

            <div v-if="row.quota > 0" class="mt-2 rounded-lg border border-gray-100 px-2.5 py-2 dark:border-dark-800">
              <div class="flex items-center justify-between gap-2 text-[10px]">
                <span class="text-gray-500 dark:text-gray-400">{{ t('keys.quota') }}</span>
                <span
                  :class="[
                    'font-medium tabular-nums',
                    row.quota_used >= row.quota ? 'text-red-500' :
                    row.quota_used >= row.quota * 0.8 ? 'text-amber-600 dark:text-amber-400' :
                    'text-gray-700 dark:text-gray-300'
                  ]"
                >
                  <span>$</span>{{ row.quota_used?.toFixed(2) || '0.00' }} / <span>$</span>{{ row.quota?.toFixed(2) }}
                </span>
              </div>
              <div class="mt-2 h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                <div
                  :class="[
                    'h-full rounded-full',
                    row.quota_used >= row.quota ? 'bg-red-500' :
                    row.quota_used >= row.quota * 0.8 ? 'bg-amber-500' :
                    'bg-emerald-600 dark:bg-emerald-400'
                  ]"
                  :style="{ width: Math.min((row.quota_used / row.quota) * 100, 100) + '%' }"
                />
              </div>
            </div>

            <div
              v-if="rateLimitRows(row).length"
              class="mt-2 space-y-2 rounded-lg border border-gray-100 px-2.5 py-2 dark:border-dark-800"
            >
              <div v-for="limit in rateLimitRows(row)" :key="limit.label">
                <div class="flex items-center justify-between gap-2 text-[10px]">
                  <span class="text-gray-500 dark:text-gray-400">{{ limit.label }}</span>
                  <span class="font-medium tabular-nums text-gray-700 dark:text-gray-300">
                    <span>$</span>{{ limit.usage.toFixed(2) }} / <span>$</span>{{ limit.limit.toFixed(2) }}
                  </span>
                </div>
                <div class="mt-1.5 h-1 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                  <div
                    class="h-full rounded-full bg-emerald-600 dark:bg-emerald-400"
                    :style="{ width: Math.min((limit.usage / limit.limit) * 100, 100) + '%' }"
                  />
                </div>
                <div v-if="limit.resetAt && formatResetTime(limit.resetAt)" class="mt-1 text-[9px] text-gray-400 dark:text-gray-500">
                  ⟳ {{ formatResetTime(limit.resetAt) }}
                </div>
              </div>
              <button
                v-if="row.usage_5h > 0 || row.usage_1d > 0 || row.usage_7d > 0"
                type="button"
                @click.stop="confirmResetRateLimit(row)"
                class="inline-flex items-center gap-1 rounded px-1 py-0.5 text-[10px] text-gray-500 transition-colors hover:text-gray-800 dark:hover:text-gray-200"
                :title="t('keys.resetRateLimitUsage')"
              >
                <Icon name="refresh" size="xs" />
                {{ t('keys.resetUsage') }}
              </button>
            </div>

            <div
              v-if="row.expires_at || row.current_concurrency > 0"
              class="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-[10px] text-gray-500 dark:text-gray-400"
            >
              <span v-if="row.current_concurrency > 0" class="inline-flex items-center gap-1">
                <Icon name="users" size="xs" />
                {{ t('keys.currentConcurrency') }}: {{ row.current_concurrency }}
              </span>
              <span v-if="row.expires_at" class="inline-flex items-center gap-1">
                <Icon name="clock" size="xs" />
                <span :class="new Date(row.expires_at) < new Date() ? 'text-red-500' : ''">
                  {{ formatDateTime(row.expires_at) }}
                </span>
              </span>
            </div>

            <div class="mt-2 flex items-center justify-between gap-2 border-t border-gray-100 pt-2 dark:border-dark-800">
              <button
                type="button"
                @click="openUseKeyModal(row)"
                class="inline-flex items-center gap-1.5 rounded-lg bg-gray-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
              >
                <Icon name="terminal" size="sm" />
                {{ t('keys.useKey') }}
              </button>
              <div class="flex items-center gap-0.5">
                <button
                  v-if="!publicSettings?.hide_ccs_import_button"
                  type="button"
                  @click="importToCcswitch(row)"
                  data-tour="key-import-ccs"
                  class="inline-flex items-center gap-1 rounded-lg px-2 py-1.5 text-[10px] font-medium text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dark-800"
                  :title="t('keys.importToCcSwitch')"
                  :aria-label="t('keys.importToCcSwitch')"
                >
                  <Icon name="upload" size="sm" />
                  <span class="hidden sm:inline">CCS</span>
                </button>
                <button
                  type="button"
                  @click="editKey(row)"
                  class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dark-800"
                  :title="t('common.edit')"
                  :aria-label="t('common.edit')"
                >
                  <Icon name="edit" size="sm" />
                </button>
                <button
                  type="button"
                  @click="confirmDelete(row)"
                  class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:text-gray-400 dark:hover:bg-red-900/20 dark:hover:text-red-400"
                  :title="t('common.delete')"
                  :aria-label="t('common.delete')"
                >
                  <Icon name="trash" size="sm" />
                </button>
              </div>
            </div>
          </article>
        </div>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > pagination.page_size"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <!-- Create/Edit Modal -->
    <BaseDialog
      :show="showCreateModal || showEditModal"
      :title="showEditModal ? t('keys.editKey') : t('keys.createKey')"
      width="normal"
      @close="closeModals"
    >
      <form id="key-form" @submit.prevent="handleSubmit" class="key-form space-y-5">
        <div>
          <label class="input-label">{{ t('keys.nameLabel') }}</label>
          <input
            v-model="formData.name"
            type="text"
            required
            class="input"
            :placeholder="t('keys.namePlaceholder')"
            data-tour="key-form-name"
          />
        </div>

        <fieldset v-if="!showEditModal" data-tour="key-form-provider">
          <legend class="input-label">{{ t('keys.providerLabel') }}</legend>
          <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
            <label
              v-for="provider in createProviderOptions"
              :key="provider.value"
              class="relative min-w-0"
              :class="provider.count === 0 ? 'cursor-not-allowed' : 'cursor-pointer'"
            >
              <input
                type="radio"
                name="key-provider"
                :value="provider.value"
                :checked="createProvider === provider.value"
                :disabled="provider.count === 0"
                class="peer sr-only"
                @change="selectCreateProvider(provider.value)"
              />
              <span
                class="flex h-full flex-col items-center gap-2 rounded-xl border border-gray-200 bg-white px-2 py-3 text-center transition-colors peer-checked:border-primary-500 peer-checked:bg-primary-50/60 peer-checked:ring-1 peer-checked:ring-primary-500 peer-focus-visible:outline peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-primary-500 peer-disabled:opacity-40 dark:border-dark-800 dark:bg-dark-900/70 dark:peer-checked:border-gray-500 dark:peer-checked:bg-dark-800 dark:peer-checked:ring-1 dark:peer-checked:ring-white/10"
                :class="provider.count > 0 && 'hover:border-primary-300 dark:hover:border-dark-600'"
              >
                <span class="flex h-8 items-center justify-center gap-1.5" aria-hidden="true">
                  <span
                    v-for="platform in KEY_GROUP_PROVIDER_ICONS[provider.value]"
                    :key="platform"
                    class="flex h-8 w-8 items-center justify-center rounded-lg"
                    :class="platformBadgeLightClass(platform)"
                  >
                    <PlatformIcon :platform="platform" size="lg" />
                  </span>
                </span>
                <span class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ provider.label }}</span>
              </span>
              <span
                v-if="createProvider === provider.value"
                class="absolute right-1.5 top-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-gray-900 text-white dark:bg-gray-100 dark:text-dark-950"
                aria-hidden="true"
              >
                <Icon name="check" size="xs" :stroke-width="3" />
              </span>
            </label>
          </div>
          <p class="mt-2 text-xs leading-5 text-gray-500 dark:text-gray-400" aria-live="polite">
            {{ groups.length === 0 ? t('common.noGroupsAvailable') : t(`keys.providerHints.${createProvider}`) }}
          </p>
        </fieldset>

        <div>
          <label class="input-label" for="key-form-group">{{ t('keys.groupLabel') }}</label>
          <Select
            :key="showEditModal ? 'edit' : createProvider"
            id="key-form-group"
            :aria-label="t('keys.groupLabel')"
            v-model="formData.group_id"
            :options="formGroupOptions"
            :placeholder="t('keys.selectGroup')"
            :empty-text="t('common.noGroupsAvailable')"
            :searchable="true"
            :search-placeholder="t('keys.searchGroup')"
            dropdown-class="key-form-dropdown"
            data-tour="key-form-group"
          >
            <template #selected="{ option }">
              <GroupBadge
                v-if="option"
                :name="(option as unknown as GroupOption).label"
                :platform="(option as unknown as GroupOption).platform"
                :subscription-type="(option as unknown as GroupOption).subscriptionType"
                :rate-multiplier="(option as unknown as GroupOption).rate"
                :user-rate-multiplier="(option as unknown as GroupOption).userRate"
                :peak-rate-enabled="(option as unknown as GroupOption).peakRateEnabled"
                :peak-start="(option as unknown as GroupOption).peakStart"
                :peak-end="(option as unknown as GroupOption).peakEnd"
                :peak-rate-multiplier="(option as unknown as GroupOption).peakRateMultiplier"
              />
              <span v-else class="text-gray-400">{{ t('keys.selectGroup') }}</span>
            </template>
            <template #option="{ option, selected }">
              <GroupOptionItem
                :name="(option as unknown as GroupOption).label"
                :platform="(option as unknown as GroupOption).platform"
                :subscription-type="(option as unknown as GroupOption).subscriptionType"
                :rate-multiplier="(option as unknown as GroupOption).rate"
                :user-rate-multiplier="(option as unknown as GroupOption).userRate"
                :peak-rate-enabled="(option as unknown as GroupOption).peakRateEnabled"
                :peak-start="(option as unknown as GroupOption).peakStart"
                :peak-end="(option as unknown as GroupOption).peakEnd"
                :peak-rate-multiplier="(option as unknown as GroupOption).peakRateMultiplier"
                :description="(option as unknown as GroupOption).description"
                :selected="selected"
              />
            </template>
          </Select>
        </div>

        <!-- Custom Key Section (only for create) -->
        <div v-if="!showEditModal" class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.customKeyLabel') }}</label>
            <button
              type="button"
              @click="formData.use_custom_key = !formData.use_custom_key"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.use_custom_key ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.use_custom_key ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>
          <div v-if="formData.use_custom_key">
            <input
              v-model="formData.custom_key"
              type="text"
              class="input font-mono"
              :placeholder="t('keys.customKeyPlaceholder')"
              :class="{ 'border-red-500 dark:border-red-500': customKeyError }"
            />
            <p v-if="customKeyError" class="mt-1 text-sm text-red-500">{{ customKeyError }}</p>
            <p v-else class="input-hint">{{ t('keys.customKeyHint') }}</p>
          </div>
        </div>

        <div v-if="showEditModal">
          <label class="input-label">{{ t('keys.statusLabel') }}</label>
          <Select
            v-model="formData.status"
            :options="statusOptions"
            :placeholder="t('keys.selectStatus')"
            dropdown-class="key-form-dropdown"
          />
        </div>

        <!-- IP Restriction Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.ipRestriction') }}</label>
            <button
              type="button"
              @click="formData.enable_ip_restriction = !formData.enable_ip_restriction"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_ip_restriction ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_ip_restriction ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>

          <div v-if="formData.enable_ip_restriction" class="space-y-4 pt-2">
            <div>
              <label class="input-label">{{ t('keys.ipWhitelist') }}</label>
              <textarea
                v-model="formData.ip_whitelist"
                rows="3"
                class="input font-mono text-sm"
                :placeholder="t('keys.ipWhitelistPlaceholder')"
              />
              <p class="input-hint">{{ t('keys.ipWhitelistHint') }}</p>
            </div>

            <div>
              <label class="input-label">{{ t('keys.ipBlacklist') }}</label>
              <textarea
                v-model="formData.ip_blacklist"
                rows="3"
                class="input font-mono text-sm"
                :placeholder="t('keys.ipBlacklistPlaceholder')"
              />
              <p class="input-hint">{{ t('keys.ipBlacklistHint') }}</p>
            </div>
          </div>
        </div>

        <!-- Quota Limit Section -->
        <div class="space-y-3">
          <label class="input-label">{{ t('keys.quotaLimit') }}</label>
          <!-- Switch commented out - always show input, 0 = unlimited
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.quotaLimit') }}</label>
            <button
              type="button"
              @click="formData.enable_quota = !formData.enable_quota"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_quota ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_quota ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>
          -->

          <div class="space-y-4">
            <div>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.quota"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="t('keys.quotaAmountPlaceholder')"
                />
              </div>
              <p class="input-hint">{{ t('keys.quotaAmountHint') }}</p>
            </div>

            <!-- Quota used display (only in edit mode) -->
            <div v-if="showEditModal && selectedKey && selectedKey.quota > 0">
              <label class="input-label">{{ t('keys.quotaUsed') }}</label>
              <div class="flex items-center gap-2">
                <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-900">
                  <span class="font-medium text-gray-900 dark:text-white">
                    ${{ selectedKey.quota_used?.toFixed(4) || '0.0000' }}
                  </span>
                  <span class="mx-2 text-gray-400">/</span>
                  <span class="text-gray-500 dark:text-gray-400">
                    ${{ selectedKey.quota?.toFixed(2) || '0.00' }}
                  </span>
                </div>
                <button
                  type="button"
                  @click="confirmResetQuota"
                  class="btn btn-secondary text-sm"
                  :title="t('keys.resetQuotaUsed')"
                >
                  {{ t('keys.reset') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Rate Limit Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.rateLimitSection') }}</label>
            <button
              type="button"
              @click="formData.enable_rate_limit = !formData.enable_rate_limit"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_rate_limit ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_rate_limit ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>

          <div v-if="formData.enable_rate_limit" class="space-y-4 pt-2">
            <p class="input-hint -mt-2">{{ t('keys.rateLimitHint') }}</p>
            <!-- 5-Hour Limit -->
            <div>
              <label class="input-label">{{ t('keys.rateLimit5h') }}</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.rate_limit_5h"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="'0'"
                />
              </div>
              <!-- Usage info (edit mode only) -->
              <div v-if="showEditModal && selectedKey && selectedKey.rate_limit_5h > 0" class="mt-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-900 text-sm">
                    <span :class="[
                      'font-medium',
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h ? 'text-red-500' :
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h * 0.8 ? 'text-yellow-500' :
                      'text-gray-900 dark:text-white'
                    ]">
                      ${{ selectedKey.usage_5h?.toFixed(4) || '0.0000' }}
                    </span>
                    <span class="mx-2 text-gray-400">/</span>
                    <span class="text-gray-500 dark:text-gray-400">
                      ${{ selectedKey.rate_limit_5h?.toFixed(2) || '0.00' }}
                    </span>
                  </div>
                </div>
                <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-800">
                  <div
                    :class="[
                      'h-full rounded-full transition-all',
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h ? 'bg-red-500' :
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h * 0.8 ? 'bg-yellow-500' :
                      'bg-green-500'
                    ]"
                    :style="{ width: Math.min((selectedKey.usage_5h / selectedKey.rate_limit_5h) * 100, 100) + '%' }"
                  />
                </div>
              </div>
            </div>

            <!-- Daily Limit -->
            <div>
              <label class="input-label">{{ t('keys.rateLimit1d') }}</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.rate_limit_1d"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="'0'"
                />
              </div>
              <!-- Usage info (edit mode only) -->
              <div v-if="showEditModal && selectedKey && selectedKey.rate_limit_1d > 0" class="mt-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-900 text-sm">
                    <span :class="[
                      'font-medium',
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d ? 'text-red-500' :
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d * 0.8 ? 'text-yellow-500' :
                      'text-gray-900 dark:text-white'
                    ]">
                      ${{ selectedKey.usage_1d?.toFixed(4) || '0.0000' }}
                    </span>
                    <span class="mx-2 text-gray-400">/</span>
                    <span class="text-gray-500 dark:text-gray-400">
                      ${{ selectedKey.rate_limit_1d?.toFixed(2) || '0.00' }}
                    </span>
                  </div>
                </div>
                <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-800">
                  <div
                    :class="[
                      'h-full rounded-full transition-all',
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d ? 'bg-red-500' :
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d * 0.8 ? 'bg-yellow-500' :
                      'bg-green-500'
                    ]"
                    :style="{ width: Math.min((selectedKey.usage_1d / selectedKey.rate_limit_1d) * 100, 100) + '%' }"
                  />
                </div>
              </div>
            </div>

            <!-- 7-Day Limit -->
            <div>
              <label class="input-label">{{ t('keys.rateLimit7d') }}</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.rate_limit_7d"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="'0'"
                />
              </div>
              <!-- Usage info (edit mode only) -->
              <div v-if="showEditModal && selectedKey && selectedKey.rate_limit_7d > 0" class="mt-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-900 text-sm">
                    <span :class="[
                      'font-medium',
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d ? 'text-red-500' :
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d * 0.8 ? 'text-yellow-500' :
                      'text-gray-900 dark:text-white'
                    ]">
                      ${{ selectedKey.usage_7d?.toFixed(4) || '0.0000' }}
                    </span>
                    <span class="mx-2 text-gray-400">/</span>
                    <span class="text-gray-500 dark:text-gray-400">
                      ${{ selectedKey.rate_limit_7d?.toFixed(2) || '0.00' }}
                    </span>
                  </div>
                </div>
                <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-800">
                  <div
                    :class="[
                      'h-full rounded-full transition-all',
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d ? 'bg-red-500' :
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d * 0.8 ? 'bg-yellow-500' :
                      'bg-green-500'
                    ]"
                    :style="{ width: Math.min((selectedKey.usage_7d / selectedKey.rate_limit_7d) * 100, 100) + '%' }"
                  />
                </div>
              </div>
            </div>

            <!-- Reset Rate Limit button (edit mode only) -->
            <div v-if="showEditModal && selectedKey && (selectedKey.rate_limit_5h > 0 || selectedKey.rate_limit_1d > 0 || selectedKey.rate_limit_7d > 0)">
              <button
                type="button"
                @click="confirmResetRateLimit(selectedKey)"
                class="btn btn-secondary text-sm"
              >
                {{ t('keys.resetRateLimitUsage') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Expiration Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.expiration') }}</label>
            <button
              type="button"
              @click="formData.enable_expiration = !formData.enable_expiration"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_expiration ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_expiration ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>

          <div v-if="formData.enable_expiration" class="space-y-4 pt-2">
            <!-- Quick select buttons (for both create and edit mode) -->
            <div class="flex flex-wrap gap-2">
              <button
                v-for="days in ['7', '30', '90']"
                :key="days"
                type="button"
                @click="setExpirationDays(parseInt(days))"
                :class="[
                  'rounded-lg px-3 py-1.5 text-sm transition-colors',
                  formData.expiration_preset === days
                    ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-900 dark:text-gray-400 dark:hover:bg-dark-800'
                ]"
              >
                {{ showEditModal ? t('keys.extendDays', { days }) : t('keys.expiresInDays', { days }) }}
              </button>
              <button
                type="button"
                @click="formData.expiration_preset = 'custom'"
                :class="[
                  'rounded-lg px-3 py-1.5 text-sm transition-colors',
                  formData.expiration_preset === 'custom'
                    ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-900 dark:text-gray-400 dark:hover:bg-dark-800'
                ]"
              >
                {{ t('keys.customDate') }}
              </button>
            </div>

            <!-- Date picker (always show for precise adjustment) -->
            <div>
              <label class="input-label">{{ t('keys.expirationDate') }}</label>
              <input
                v-model="formData.expiration_date"
                type="datetime-local"
                class="input"
              />
              <p class="input-hint">{{ t('keys.expirationDateHint') }}</p>
            </div>

            <!-- Current expiration display (only in edit mode) -->
            <div v-if="showEditModal && selectedKey?.expires_at" class="text-sm">
              <span class="text-gray-500 dark:text-gray-400">{{ t('keys.currentExpiration') }}: </span>
              <span class="font-medium text-gray-900 dark:text-white">
                {{ formatDateTime(selectedKey.expires_at) }}
              </span>
            </div>
          </div>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button @click="closeModals" type="button" class="btn btn-secondary">
            {{ t('common.cancel') }}
          </button>
          <button
            form="key-form"
            type="submit"
            :disabled="submitting"
            class="btn btn-primary"
            data-tour="key-form-submit"
          >
            <svg
              v-if="submitting"
              class="-ml-1 mr-2 h-4 w-4 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{
              submitting
                ? t('keys.saving')
                : showEditModal
                  ? t('common.update')
                  : t('common.create')
            }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <BulkEditKeysModal
      :show="showBulkEditModal"
      :selected-keys="selectedApiKeys"
      :groups="groups"
      @close="showBulkEditModal = false"
      @updated="handleBulkUpdated"
    />

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      :show="showDeleteDialog"
      :title="t('keys.deleteKey')"
      :message="t('keys.deleteConfirmMessage', { name: selectedKey?.name })"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="handleDelete"
      @cancel="showDeleteDialog = false"
    />

    <!-- Reset Quota Confirmation Dialog -->
    <ConfirmDialog
      :show="showResetQuotaDialog"
      :title="t('keys.resetQuotaTitle')"
      :message="t('keys.resetQuotaConfirmMessage', { name: selectedKey?.name, used: selectedKey?.quota_used?.toFixed(4) })"
      :confirm-text="t('keys.reset')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="resetQuotaUsed"
      @cancel="showResetQuotaDialog = false"
    />

    <!-- Reset Rate Limit Confirmation Dialog -->
    <ConfirmDialog
      :show="showResetRateLimitDialog"
      :title="t('keys.resetRateLimitTitle')"
      :message="t('keys.resetRateLimitConfirmMessage', { name: selectedKey?.name })"
      :confirm-text="t('keys.reset')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="resetRateLimitUsage"
      @cancel="showResetRateLimitDialog = false"
    />

    <!-- Use Key Modal -->
    <UseKeyModal
      :show="showUseKeyModal"
      :api-key="selectedKey?.key || ''"
      :base-url="publicSettings?.api_base_url || ''"
      :platform="selectedKey?.group?.platform || null"
      :claude-code-only="selectedKey?.group?.claude_code_only || false"
      :allow-messages-dispatch="selectedKey?.group?.allow_messages_dispatch || false"
      @close="closeUseKeyModal"
    />

    <!-- CC Switch Client Selection Dialog -->
    <BaseDialog
      :show="showCcsClientSelect"
      :title="t('keys.ccsClientSelect.title')"
      width="narrow"
      @close="closeCcsClientSelect"
    >
      <div class="space-y-4">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ t('keys.ccsClientSelect.description') }}
        </p>
        <div class="grid grid-cols-2 gap-3" data-tour="ccs-clients">
          <button
            v-for="option in ccsClientOptions"
            :key="option.clientType"
            type="button"
            :disabled="option.disabled"
            :title="option.disabled ? t('keys.ccsClientSelect.messagesDispatchRequired') : undefined"
            @click="handleCcsClientSelect(option.clientType)"
            :class="[
              'flex min-h-28 flex-col items-center justify-center gap-2 rounded-lg border-2 p-3 text-center transition-colors',
              selectedCcsClient === option.clientType
                ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                : 'border-gray-200 hover:border-primary-500 hover:bg-primary-50 dark:border-dark-600 dark:hover:bg-primary-900/20',
              option.disabled ? 'cursor-not-allowed opacity-50' : ''
            ]"
          >
            <PlatformIcon :platform="option.icon" size="lg" class="!h-7 !w-7 text-gray-700 dark:text-gray-200" aria-hidden="true" />
            <span class="font-medium text-gray-900 dark:text-white">{{ t(option.label) }}</span>
          </button>
        </div>
        <div v-if="!ccsModelsLoading && ccsModelsUnavailable" class="flex items-center justify-between gap-3 text-xs text-amber-600 dark:text-amber-400">
          <span>{{ ccsModelsError }}</span>
          <button type="button" class="shrink-0 underline" @click="retryLoadCcsModels">
            {{ t('keys.ccsClientSelect.retryModels') }}
          </button>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-between">
          <button type="button" @click="closeCcsClientSelect" class="btn btn-secondary">
            {{ t('common.cancel') }}
          </button>
          <button
            v-if="selectedCcsClient"
            type="button"
            class="btn btn-primary"
            :disabled="!canConfirmCcsImport"
            @click="confirmCcsClientSelect"
            data-tour="ccs-confirm"
            data-tour-manual
          >
            {{ t('keys.ccsClientSelect.import') }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Group Selector Dropdown (Teleported to body to avoid overflow clipping) -->
    <Teleport to="body">
      <div
        v-if="groupSelectorKeyId !== null && dropdownPosition"
        ref="dropdownRef"
        class="keys-group-dropdown animate-in fade-in slide-in-from-top-2 fixed z-[100000020] w-max max-w-[calc(100vw-16px)] overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl shadow-black/20 ring-1 ring-black/5 duration-200 sm:min-w-[380px] dark:border-white/10 dark:bg-[#0d0f11] dark:shadow-black/50 dark:ring-white/[0.04]"
        style="pointer-events: auto !important;"
        :style="{
          top: dropdownPosition.top !== undefined ? dropdownPosition.top + 'px' : undefined,
          bottom: dropdownPosition.bottom !== undefined ? dropdownPosition.bottom + 'px' : undefined,
          left: dropdownPosition.left + 'px'
        }"
      >
        <!-- Search box -->
        <div class="border-b border-gray-100 p-2 dark:border-white/[0.06]">
          <div class="relative">
            <svg class="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              v-model="groupSearchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-200 bg-gray-50 py-1.5 pl-8 pr-3 text-sm text-gray-900 placeholder-gray-400 outline-none focus:border-primary-300 focus:ring-1 focus:ring-primary-300 dark:border-white/[0.08] dark:bg-white/[0.04] dark:text-white dark:placeholder-gray-500 dark:focus:border-white/20 dark:focus:ring-white/10"
              :placeholder="t('keys.searchGroup')"
              @click.stop
            />
          </div>
        </div>
        <!-- Group list -->
        <div class="max-h-80 overflow-y-auto p-1.5">
          <button
            v-for="option in filteredGroupOptions"
            :key="option.value ?? 'null'"
            @click="changeGroup(selectedKeyForGroup!, option.value)"
            :class="[
              'flex w-full items-center justify-between rounded-lg px-3 py-2.5 text-sm transition-colors',
              selectedKeyForGroup?.group_id === option.value ||
              (!selectedKeyForGroup?.group_id && option.value === null)
                ? 'bg-primary-50 dark:bg-white/[0.08]'
                : 'hover:bg-gray-100 dark:hover:bg-white/[0.05]'
            ]"
            :title="option.description || undefined"
          >
            <GroupOptionItem
              :name="option.label"
              :platform="option.platform"
              :subscription-type="option.subscriptionType"
              :rate-multiplier="option.rate"
              :user-rate-multiplier="option.userRate"
              :peak-rate-enabled="option.peakRateEnabled"
              :peak-start="option.peakStart"
              :peak-end="option.peakEnd"
              :peak-rate-multiplier="option.peakRateMultiplier"
              :description="option.description"
              :selected="
                selectedKeyForGroup?.group_id === option.value ||
                (!selectedKeyForGroup?.group_id && option.value === null)
              "
            />
          </button>
          <!-- Empty state when search has no results -->
          <div v-if="filteredGroupOptions.length === 0" class="py-4 text-center text-sm text-gray-400 dark:text-gray-500">
            {{ t('keys.noGroupFound') }}
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { fetchModelCatalog, ModelCatalogError } from '@/api/modelCatalog'
	import { ref, computed, watch, onMounted, onUnmounted, type ComponentPublicInstance } from 'vue'
	import { useI18n } from 'vue-i18n'
	import { useAppStore } from '@/stores/app'
	import { useOnboardingStore } from '@/stores/onboarding'
	import { useClipboard } from '@/composables/useClipboard'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'

const { t } = useI18n()
import { keysAPI, authAPI, usageAPI, userGroupsAPI } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BulkEditKeysModal from '@/components/keys/BulkEditKeysModal.vue'
	import Pagination from '@/components/common/Pagination.vue'
	import BaseDialog from '@/components/common/BaseDialog.vue'
	import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
	import EmptyState from '@/components/common/EmptyState.vue'
	import Select from '@/components/common/Select.vue'
	import SearchInput from '@/components/common/SearchInput.vue'
	import Icon from '@/components/icons/Icon.vue'
	import UseKeyModal from '@/components/keys/UseKeyModal.vue'
	import EndpointPopover from '@/components/keys/EndpointPopover.vue'
	import GroupBadge from '@/components/common/GroupBadge.vue'
	import GroupOptionItem from '@/components/common/GroupOptionItem.vue'
	import type { ApiKey, Group, PublicSettings, SubscriptionType, GroupPlatform, UpdateApiKeyRequest } from '@/types'
import type { BatchApiKeyUsageStats } from '@/api/usage'
import { formatDateTime } from '@/utils/format'
import { maskApiKey } from '@/utils/maskApiKey'
import PlatformIcon from '@/components/common/PlatformIcon.vue'
import { groupDisplayPlatform } from '@/utils/groupPresentation'
import { platformBadgeLightClass } from '@/utils/platformColors'
import { KEY_GROUP_PROVIDERS, KEY_GROUP_PROVIDER_ICONS, getKeyGroupProvider, type KeyGroupProvider } from '@/utils/keyGroupProviders'
import {
  CC_SWITCH_USAGE_SCRIPT,
  buildCcSwitchImportDeeplink,
  type CcSwitchClientType
} from '@/utils/ccswitchImport'

// Helper to format date for datetime-local input
const formatDateTimeLocal = (isoDate: string): string => {
  const date = new Date(isoDate)
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

interface GroupOption {
  value: number
  label: string
  description: string | null
  rate: number
  userRate: number | null
  peakRateEnabled: boolean
  peakStart: string
  peakEnd: string
  peakRateMultiplier: number
  subscriptionType: SubscriptionType
  platform: GroupPlatform
}

type CcsClientIcon = 'openai' | 'anthropic' | 'deepseek' | 'gemini'

interface CcsClientOption {
  clientType: CcSwitchClientType
  label: string
  icon: CcsClientIcon
  disabled?: boolean
}

const appStore = useAppStore()
const onboardingStore = useOnboardingStore()
const { copyToClipboard: clipboardCopy } = useClipboard()

const apiKeys = ref<ApiKey[]>([])
const selectedIds = ref<number[]>([])
const showBulkEditModal = ref(false)
const selectedApiKeys = computed(() => apiKeys.value.filter((key) => selectedIds.value.includes(key.id)))

const handleSelectionChange = (ids: Array<string | number>) => {
  const visibleIds = new Set(apiKeys.value.map((key) => key.id))
  selectedIds.value = [...new Set(ids.map(Number))].filter((id) => visibleIds.has(id))
}

const allPageSelected = computed(() =>
  apiKeys.value.length > 0 && apiKeys.value.every((key) => selectedIds.value.includes(key.id))
)

const togglePageSelection = () => {
  selectedIds.value = allPageSelected.value ? [] : apiKeys.value.map((key) => key.id)
}

const toggleKeyCardSelection = (keyId: number) => {
  handleSelectionChange(selectedIds.value.includes(keyId)
    ? selectedIds.value.filter((id) => id !== keyId)
    : [...selectedIds.value, keyId])
}

const handleKeyCardClick = (row: ApiKey, event: MouseEvent) => {
  const target = event.target
  if (target instanceof Element && target.closest('button, a, input, select, textarea, [role="switch"]')) return
  toggleKeyCardSelection(row.id)
}

const handleKeyCardKeydown = (row: ApiKey, event: KeyboardEvent) => {
  if (event.target !== event.currentTarget || (event.key !== 'Enter' && event.key !== ' ')) return
  event.preventDefault()
  toggleKeyCardSelection(row.id)
}

const handleBulkUpdated = (succeededIds: number[]) => {
  const succeeded = new Set(succeededIds)
  selectedIds.value = selectedIds.value.filter((id) => !succeeded.has(id))
  loadApiKeys()
}

const groups = ref<Group[]>([])
const loading = ref(false)
const showInitialSkeleton = ref(false)
const submitting = ref(false)
const now = ref(new Date())
let resetTimer: ReturnType<typeof setInterval> | null = null
let initialSkeletonTimer: ReturnType<typeof setTimeout> | null = null
let initialSkeletonVisibleAt: number | null = null
const INITIAL_SKELETON_DELAY_MS = 180
const INITIAL_SKELETON_MIN_VISIBLE_MS = 280
const usageStats = ref<Record<string, BatchApiKeyUsageStats>>({})
const usageLoading = ref(false)
const userGroupRates = ref<Record<number, number>>({})

const pagination = ref({
  page: 1,
  page_size: getPersistedPageSize(),
  total: 0,
  pages: 0
})
const sortState = ref({
  sort_by: 'created_at',
  sort_order: 'desc' as 'asc' | 'desc'
})

// Filter state
const filterSearch = ref('')
const filterStatus = ref('')
const filterGroupId = ref<string | number>('')

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)
const showResetQuotaDialog = ref(false)
const showResetRateLimitDialog = ref(false)
const showUseKeyModal = ref(false)
const showCcsClientSelect = ref(false)
const pendingCcsRow = ref<ApiKey | null>(null)
const selectedCcsClient = ref<CcSwitchClientType | null>(null)
const selectedCcsModel = ref('')
const ccsModelOptions = ref<string[]>([])
const ccsModelsLoading = ref(false)
const ccsModelsUnavailable = ref(false)
const ccsModelsError = ref('')
let ccsModelsRequest = 0
let ccsModelsController: AbortController | undefined
const requiresCcsModel = computed(() =>
  selectedCcsClient.value === 'codex' ||
  selectedCcsClient.value === 'claude-desktop' ||
  selectedCcsClient.value === 'dsh'
)
const canConfirmCcsImport = computed(() =>
  !!selectedCcsClient.value &&
  (!requiresCcsModel.value || selectedCcsModel.value.trim().length > 0) &&
  !ccsModelsLoading.value
)
const ccsClientOptions = computed<CcsClientOption[]>(() => {
  const platform = pendingCcsRow.value?.group?.platform || 'anthropic'
  const messagesDispatchEnabled = pendingCcsRow.value?.group?.allow_messages_dispatch === true
  const option = (
    clientType: CcSwitchClientType,
    label: string,
    icon: CcsClientIcon,
    disabled = false
  ): CcsClientOption => ({ clientType, label, icon, disabled })

  if (platform === 'gemini') {
    return [option('gemini', 'keys.ccsClientSelect.geminiCli', 'gemini')]
  }

  if (platform === 'antigravity') {
    return [
      option('claude', 'keys.ccsClientSelect.claudeCode', 'anthropic'),
      option('gemini', 'keys.ccsClientSelect.geminiCli', 'gemini')
    ]
  }

  const claudeOptions = [
    option('claude', 'keys.ccsClientSelect.claudeCode', 'anthropic'),
    option(
      'claude-desktop',
      'keys.ccsClientSelect.claudeDesktop',
      'anthropic'
    )
  ]

  if (platform === 'openai') {
    return [
      option('codex', 'keys.ccsClientSelect.codex', 'openai'),
      ...claudeOptions.map((item) => ({ ...item, disabled: !messagesDispatchEnabled })),
      option('dsh', 'keys.ccsClientSelect.deepSeekHarness', 'deepseek')
    ]
  }

  if (platform !== 'anthropic') {
    return [
      option('codex', 'keys.ccsClientSelect.codex', 'openai'),
      ...claudeOptions,
      option('dsh', 'keys.ccsClientSelect.deepSeekHarness', 'deepseek')
    ]
  }

  return [
    ...claudeOptions,
    option('dsh', 'keys.ccsClientSelect.deepSeekHarness', 'deepseek')
  ]
})
const selectedKey = ref<ApiKey | null>(null)
const copiedKeyId = ref<number | null>(null)
const groupSelectorKeyId = ref<number | null>(null)
const publicSettings = ref<PublicSettings | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)
const dropdownPosition = ref<{ top?: number; bottom?: number; left: number } | null>(null)
const groupButtonRefs = ref<Map<number, HTMLElement>>(new Map())
let abortController: AbortController | null = null

// Get the currently selected key for group change
const selectedKeyForGroup = computed(() => {
  if (groupSelectorKeyId.value === null) return null
  return apiKeys.value.find((k) => k.id === groupSelectorKeyId.value) || null
})

const setGroupButtonRef = (keyId: number, el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    groupButtonRefs.value.set(keyId, el)
  } else {
    groupButtonRefs.value.delete(keyId)
  }
}

const formData = ref({
  name: '',
  group_id: null as number | null,
  status: 'active' as 'active' | 'inactive',
  use_custom_key: false,
  custom_key: '',
  enable_ip_restriction: false,
  ip_whitelist: '',
  ip_blacklist: '',
  // Quota settings (empty = unlimited)
  enable_quota: false,
  quota: null as number | null,
  // Rate limit settings
  enable_rate_limit: false,
  rate_limit_5h: null as number | null,
  rate_limit_1d: null as number | null,
  rate_limit_7d: null as number | null,
  enable_expiration: false,
  expiration_preset: '30' as '7' | '30' | '90' | 'custom',
  expiration_date: ''
})

// 自定义Key验证
const customKeyError = computed(() => {
  if (!formData.value.use_custom_key || !formData.value.custom_key) {
    return ''
  }
  const key = formData.value.custom_key
  if (key.length < 16) {
    return t('keys.customKeyTooShort')
  }
  // 检查字符：只允许字母、数字、下划线、连字符
  if (!/^[a-zA-Z0-9_-]+$/.test(key)) {
    return t('keys.customKeyInvalidChars')
  }
  return ''
})

const statusOptions = computed(() => [
  { value: 'active', label: t('common.active') },
  { value: 'inactive', label: t('common.inactive') }
])

const shouldSubmitEditStatus = (key: ApiKey, status: 'active' | 'inactive') => {
  if (key.status === 'quota_exhausted' || key.status === 'expired') {
    return status === 'active'
  }
  return true
}

// Convert groups to Select options format with rate multiplier and subscription type
const groupOptions = computed(() =>
  groups.value.map((group) => ({
    value: group.id,
    label: group.name,
    description: group.description,
    rate: group.rate_multiplier,
    userRate: userGroupRates.value[group.id] ?? null,
    peakRateEnabled: group.peak_rate_enabled,
    peakStart: group.peak_start,
    peakEnd: group.peak_end,
    peakRateMultiplier: group.peak_rate_multiplier,
    subscriptionType: group.subscription_type,
    platform: groupDisplayPlatform(group) as GroupPlatform
  }))
)

// Filter dropdown options
const groupFilterOptions = computed(() => [
  { value: '', label: t('keys.allGroups'), kind: 'all' as const },
  { value: 0, label: t('keys.noGroup'), kind: 'none' as const },
  ...groupOptions.value.map((option) => ({ ...option, kind: 'group' as const }))
])

interface StatusFilterOption {
  value: string
  label: string
  icon: 'filter' | 'checkCircle' | 'ban' | 'exclamationTriangle' | 'clock'
  badgeClass: string
  [key: string]: unknown
}

const statusFilterOptions = computed<StatusFilterOption[]>(() => [
  {
    value: '',
    label: t('keys.allStatus'),
    icon: 'filter',
    badgeClass: 'bg-gray-100 text-gray-600 dark:bg-dark-950 dark:text-gray-300'
  },
  {
    value: 'active',
    label: t('keys.status.active'),
    icon: 'checkCircle',
    badgeClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
  },
  {
    value: 'inactive',
    label: t('keys.status.inactive'),
    icon: 'ban',
    badgeClass: 'bg-gray-100 text-gray-600 dark:bg-dark-950 dark:text-gray-300'
  },
  {
    value: 'quota_exhausted',
    label: t('keys.status.quota_exhausted'),
    icon: 'exclamationTriangle',
    badgeClass: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  },
  {
    value: 'expired',
    label: t('keys.status.expired'),
    icon: 'clock',
    badgeClass: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
  }
])

const selectedGroupFilterLabel = computed(() =>
  groupFilterOptions.value.find((option) => option.value === filterGroupId.value)?.label ?? t('keys.allGroups')
)

const selectedStatusFilterLabel = computed(() =>
  statusFilterOptions.value.find((option) => option.value === filterStatus.value)?.label ?? t('keys.allStatus')
)

const onFilterChange = () => {
  selectedIds.value = []
  pagination.value.page = 1
  loadApiKeys()
}

const onGroupFilterChange = (value: string | number | boolean | null) => {
  filterGroupId.value = value as string | number
  onFilterChange()
}

const onStatusFilterChange = (value: string | number | boolean | null) => {
  filterStatus.value = value as string
  onFilterChange()
}

const createProvider = ref<KeyGroupProvider>('anthropic')
const createProviderOptions = computed(() => KEY_GROUP_PROVIDERS.map((value) => ({
  value,
  label: t(`keys.providers.${value}`),
  count: groups.value.filter((group) => getKeyGroupProvider(groupDisplayPlatform(group) as GroupPlatform) === value).length
})))

const formGroupOptions = computed(() => showEditModal.value
  ? groupOptions.value
  : groupOptions.value.filter((group) => getKeyGroupProvider(group.platform) === createProvider.value)
)

const selectCreateProvider = (provider: KeyGroupProvider) => {
  if (createProvider.value === provider) return
  createProvider.value = provider
  formData.value.group_id = null
}

// Also handles groups arriving after the create dialog has already opened.
watch([showCreateModal, createProviderOptions], ([isOpen, providers], [wasOpen]) => {
  if (!isOpen) return
  if (!wasOpen || !providers.some((provider) => provider.value === createProvider.value && provider.count > 0)) {
    selectCreateProvider(providers.find((provider) => provider.count > 0)?.value ?? 'anthropic')
  }
  if (!formGroupOptions.value.some((group) => group.value === formData.value.group_id)) {
    formData.value.group_id = null
  }
})

// Group dropdown search
const groupSearchQuery = ref('')
const filteredGroupOptions = computed(() => {
  const query = groupSearchQuery.value.trim().toLowerCase()
  if (!query) return groupOptions.value
  return groupOptions.value.filter((opt) => {
    return opt.label.toLowerCase().includes(query) ||
      (opt.description && opt.description.toLowerCase().includes(query))
  })
})

const copyToClipboard = async (text: string, keyId: number) => {
  const success = await clipboardCopy(text, t('keys.copied'))
  if (success) {
    copiedKeyId.value = keyId
    setTimeout(() => {
      copiedKeyId.value = null
    }, 800)
  }
}

const isAbortError = (error: unknown) => {
  if (!error || typeof error !== 'object') return false
  const { name, code } = error as { name?: string; code?: string }
  return name === 'AbortError' || code === 'ERR_CANCELED'
}

const scheduleInitialSkeleton = () => {
  if (initialSkeletonTimer) clearTimeout(initialSkeletonTimer)
  initialSkeletonTimer = null
  initialSkeletonVisibleAt = null
  showInitialSkeleton.value = false

  if (apiKeys.value.length > 0) return

  initialSkeletonTimer = setTimeout(() => {
    showInitialSkeleton.value = true
    initialSkeletonVisibleAt = Date.now()
    initialSkeletonTimer = null
  }, INITIAL_SKELETON_DELAY_MS)
}

const finishInitialSkeleton = async (controller: AbortController) => {
  if (initialSkeletonTimer) clearTimeout(initialSkeletonTimer)
  initialSkeletonTimer = null

  if (initialSkeletonVisibleAt !== null) {
    const visibleFor = Date.now() - initialSkeletonVisibleAt
    const remaining = INITIAL_SKELETON_MIN_VISIBLE_MS - visibleFor
    if (remaining > 0) {
      await new Promise<void>((resolve) => setTimeout(resolve, remaining))
    }
  }

  if (abortController !== controller) return
  showInitialSkeleton.value = false
  initialSkeletonVisibleAt = null
}

const loadApiKeys = async () => {
  abortController?.abort()
  const controller = new AbortController()
  abortController = controller
  const { signal } = controller
  loading.value = true
  scheduleInitialSkeleton()
  usageLoading.value = false
  try {
    // Build filters
    const filters: {
      search?: string
      status?: string
      group_id?: number | string
      sort_by?: string
      sort_order?: 'asc' | 'desc'
    } = {}
    if (filterSearch.value) filters.search = filterSearch.value
    if (filterStatus.value) filters.status = filterStatus.value
    if (filterGroupId.value !== '') filters.group_id = filterGroupId.value
    filters.sort_by = sortState.value.sort_by
    filters.sort_order = sortState.value.sort_order

    const response = await keysAPI.list(pagination.value.page, pagination.value.page_size, filters, {
      signal
    })
    if (signal.aborted) return
    apiKeys.value = response.items
    handleSelectionChange(selectedIds.value)
    pagination.value.total = response.total
    pagination.value.pages = response.pages

    // 用量统计要聚合 30 天的 usage_logs，可能比列表慢很多：不阻塞表格渲染，单独加载
    usageLoading.value = response.items.length > 0
    if (response.items.length > 0) {
      void loadUsageStats(response.items.map((k) => k.id), signal)
    }
  } catch (error) {
    if (isAbortError(error)) {
      return
    }
    appStore.showError(t('keys.failedToLoad'))
  } finally {
    if (abortController === controller) {
      await finishInitialSkeleton(controller)
    }
    if (abortController === controller) {
      loading.value = false
    }
  }
}

const loadUsageStats = async (keyIds: number[], signal: AbortSignal) => {
  try {
    const usageResponse = await usageAPI.getDashboardApiKeysUsage(keyIds, { signal })
    if (signal.aborted) return
    usageStats.value = usageResponse.stats
  } catch (e) {
    if (!isAbortError(e)) {
      console.error('Failed to load usage stats:', e)
    }
  } finally {
    if (!signal.aborted) {
      usageLoading.value = false
    }
  }
}

const loadGroups = async () => {
  try {
    groups.value = await userGroupsAPI.getAvailable()
  } catch (error) {
    console.error('Failed to load groups:', error)
  }
}

const loadUserGroupRates = async () => {
  try {
    userGroupRates.value = await userGroupsAPI.getUserGroupRates()
  } catch (error) {
    console.error('Failed to load user group rates:', error)
  }
}

const loadPublicSettings = async () => {
  try {
    publicSettings.value = await authAPI.getPublicSettings()
  } catch (error) {
    console.error('Failed to load public settings:', error)
  }
}

const openUseKeyModal = (key: ApiKey) => {
  selectedKey.value = key
  showUseKeyModal.value = true
}

const closeUseKeyModal = () => {
  showUseKeyModal.value = false
  selectedKey.value = null
}

const handlePageChange = (page: number) => {
  selectedIds.value = []
  pagination.value.page = page
  loadApiKeys()
}

const handlePageSizeChange = (pageSize: number) => {
  selectedIds.value = []
  pagination.value.page_size = pageSize
  pagination.value.page = 1
  loadApiKeys()
}

const keyStatusClass = (status: string) => {
  switch (status) {
    case 'active':
      return 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/25 dark:text-emerald-300'
    case 'quota_exhausted':
      return 'bg-amber-50 text-amber-700 dark:bg-amber-900/25 dark:text-amber-300'
    case 'expired':
      return 'bg-red-50 text-red-700 dark:bg-red-900/25 dark:text-red-300'
    default:
      return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'
  }
}

const rateLimitRows = (key: ApiKey) => [
  { label: '5h', usage: key.usage_5h || 0, limit: key.rate_limit_5h || 0, resetAt: key.reset_5h_at },
  { label: '1d', usage: key.usage_1d || 0, limit: key.rate_limit_1d || 0, resetAt: key.reset_1d_at },
  { label: '7d', usage: key.usage_7d || 0, limit: key.rate_limit_7d || 0, resetAt: key.reset_7d_at }
].filter((window) => window.limit > 0)

const editKey = (key: ApiKey) => {
  selectedKey.value = key
  const hasIPRestriction = (key.ip_whitelist?.length > 0) || (key.ip_blacklist?.length > 0)
  const hasExpiration = !!key.expires_at
  formData.value = {
    name: key.name,
    group_id: key.group_id,
    status: key.status === 'quota_exhausted' || key.status === 'expired' ? 'inactive' : key.status,
    use_custom_key: false,
    custom_key: '',
    enable_ip_restriction: hasIPRestriction,
    ip_whitelist: (key.ip_whitelist || []).join('\n'),
    ip_blacklist: (key.ip_blacklist || []).join('\n'),
    enable_quota: key.quota > 0,
    quota: key.quota > 0 ? key.quota : null,
    enable_rate_limit: (key.rate_limit_5h > 0) || (key.rate_limit_1d > 0) || (key.rate_limit_7d > 0),
    rate_limit_5h: key.rate_limit_5h || null,
    rate_limit_1d: key.rate_limit_1d || null,
    rate_limit_7d: key.rate_limit_7d || null,
    enable_expiration: hasExpiration,
    expiration_preset: 'custom',
    expiration_date: key.expires_at ? formatDateTimeLocal(key.expires_at) : ''
  }
  showEditModal.value = true
}

const toggleKeyStatus = async (key: ApiKey) => {
  const newStatus = key.status === 'active' ? 'inactive' : 'active'
  try {
    await keysAPI.toggleStatus(key.id, newStatus)
    appStore.showSuccess(
      newStatus === 'active' ? t('keys.keyEnabledSuccess') : t('keys.keyDisabledSuccess')
    )
    loadApiKeys()
  } catch (error) {
    appStore.showError(t('keys.failedToUpdateStatus'))
  }
}

const openGroupSelector = (key: ApiKey) => {
  if (groupSelectorKeyId.value === key.id) {
    groupSelectorKeyId.value = null
    dropdownPosition.value = null
  } else {
    const buttonEl = groupButtonRefs.value.get(key.id)
    if (buttonEl) {
      const rect = buttonEl.getBoundingClientRect()
      const dropdownEstHeight = 400 // estimated max dropdown height
      const dropdownEstWidth = Math.min(380, window.innerWidth - 16)
      const spaceBelow = window.innerHeight - rect.bottom
      const spaceAbove = rect.top
      // 夹取 left，避免窄屏下浮层超出视口右缘
      const left = Math.max(8, Math.min(rect.left, window.innerWidth - dropdownEstWidth - 8))

      if (spaceBelow < dropdownEstHeight && spaceAbove > spaceBelow) {
        // Not enough space below, pop upward
        dropdownPosition.value = {
          bottom: window.innerHeight - rect.top + 4,
          left
        }
      } else {
        // Default: pop downward
        dropdownPosition.value = {
          top: rect.bottom + 4,
          left
        }
      }
    }
    groupSelectorKeyId.value = key.id
    groupSearchQuery.value = ''
  }
}

const changeGroup = async (key: ApiKey, newGroupId: number | null) => {
  groupSelectorKeyId.value = null
  dropdownPosition.value = null
  if (key.group_id === newGroupId) return

  try {
    await keysAPI.update(key.id, { group_id: newGroupId })
    appStore.showSuccess(t('keys.groupChangedSuccess'))
    loadApiKeys()
  } catch (error) {
    appStore.showError(t('keys.failedToChangeGroup'))
  }
}

const closeGroupSelector = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  // Check if click is inside the dropdown or the trigger button
  if (!target.closest('.group\\/dropdown') && !dropdownRef.value?.contains(target)) {
    groupSelectorKeyId.value = null
    dropdownPosition.value = null
  }
}

const confirmDelete = (key: ApiKey) => {
  selectedKey.value = key
  showDeleteDialog.value = true
}

const handleSubmit = async () => {
  // Validate group_id is required
  if (formData.value.group_id === null) {
    appStore.showError(t('keys.groupRequired'))
    return
  }

  // Validate custom key if enabled
  if (!showEditModal.value && formData.value.use_custom_key) {
    if (!formData.value.custom_key) {
      appStore.showError(t('keys.customKeyRequired'))
      return
    }
    if (customKeyError.value) {
      appStore.showError(customKeyError.value)
      return
    }
  }

  // Parse IP lists only if IP restriction is enabled
  const parseIPList = (text: string): string[] =>
    text.split('\n').map(ip => ip.trim()).filter(ip => ip.length > 0)
  const ipWhitelist = formData.value.enable_ip_restriction ? parseIPList(formData.value.ip_whitelist) : []
  const ipBlacklist = formData.value.enable_ip_restriction ? parseIPList(formData.value.ip_blacklist) : []

  // Calculate quota value (null/empty/0 = unlimited, stored as 0)
  const quota = formData.value.quota && formData.value.quota > 0 ? formData.value.quota : 0

  // Calculate expiration
  let expiresInDays: number | undefined
  let expiresAt: string | null | undefined
  if (formData.value.enable_expiration && formData.value.expiration_date) {
    if (!showEditModal.value) {
      // Create mode: calculate days from date
      const expDate = new Date(formData.value.expiration_date)
      const now = new Date()
      const diffDays = Math.ceil((expDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
      expiresInDays = diffDays > 0 ? diffDays : 1
    } else {
      // Edit mode: use custom date directly
      expiresAt = new Date(formData.value.expiration_date).toISOString()
    }
  } else if (showEditModal.value) {
    // Edit mode: if expiration disabled or date cleared, send empty string to clear
    expiresAt = ''
  }

  // Calculate rate limit values (send 0 when toggle is off)
  const rateLimitData = formData.value.enable_rate_limit ? {
    rate_limit_5h: formData.value.rate_limit_5h && formData.value.rate_limit_5h > 0 ? formData.value.rate_limit_5h : 0,
    rate_limit_1d: formData.value.rate_limit_1d && formData.value.rate_limit_1d > 0 ? formData.value.rate_limit_1d : 0,
    rate_limit_7d: formData.value.rate_limit_7d && formData.value.rate_limit_7d > 0 ? formData.value.rate_limit_7d : 0,
  } : { rate_limit_5h: 0, rate_limit_1d: 0, rate_limit_7d: 0 }

  submitting.value = true
  try {
    if (showEditModal.value && selectedKey.value) {
      const updates: UpdateApiKeyRequest = {
        name: formData.value.name,
        group_id: formData.value.group_id,
        ip_whitelist: ipWhitelist,
        ip_blacklist: ipBlacklist,
        quota: quota,
        expires_at: expiresAt,
        rate_limit_5h: rateLimitData.rate_limit_5h,
        rate_limit_1d: rateLimitData.rate_limit_1d,
        rate_limit_7d: rateLimitData.rate_limit_7d,
      }
      if (shouldSubmitEditStatus(selectedKey.value, formData.value.status)) {
        updates.status = formData.value.status
      }
      await keysAPI.update(selectedKey.value.id, updates)
      appStore.showSuccess(t('keys.keyUpdatedSuccess'))
    } else {
      const customKey = formData.value.use_custom_key ? formData.value.custom_key : undefined
      await keysAPI.create(
        formData.value.name,
        formData.value.group_id,
        customKey,
        ipWhitelist,
        ipBlacklist,
        quota,
        expiresInDays,
        rateLimitData
      )
      appStore.showSuccess(t('keys.keyCreatedSuccess'))
    }
    const advanceTour = onboardingStore.isCurrentStep('[data-tour="key-form-submit"]')
    closeModals()
    await loadApiKeys()
    if (advanceTour) void onboardingStore.nextStep()
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToSave')
    appStore.showError(errorMsg)
    // Don't advance tour on error
  } finally {
    submitting.value = false
  }
}

/**
 * 处理删除 API Key 的操作
 * 优化：错误处理改进，优先显示后端返回的具体错误消息（如权限不足等），
 * 若后端未返回消息则显示默认的国际化文本
 */
const handleDelete = async () => {
  if (!selectedKey.value) return

  try {
    await keysAPI.delete(selectedKey.value.id)
    appStore.showSuccess(t('keys.keyDeletedSuccess'))
    showDeleteDialog.value = false
    loadApiKeys()
  } catch (error: any) {
    // 优先使用后端返回的错误消息，提供更具体的错误信息给用户
    const errorMsg = error?.message || t('keys.failedToDelete')
    appStore.showError(errorMsg)
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  selectedKey.value = null
  formData.value = {
    name: '',
    group_id: null,
    status: 'active',
    use_custom_key: false,
    custom_key: '',
    enable_ip_restriction: false,
    ip_whitelist: '',
    ip_blacklist: '',
    enable_quota: false,
    quota: null,
    enable_rate_limit: false,
    rate_limit_5h: null,
    rate_limit_1d: null,
    rate_limit_7d: null,
    enable_expiration: false,
    expiration_preset: '30',
    expiration_date: ''
  }
}

// Show reset quota confirmation dialog
const confirmResetQuota = () => {
  showResetQuotaDialog.value = true
}

// Set expiration date based on quick select days
const setExpirationDays = (days: number) => {
  formData.value.expiration_preset = days.toString() as '7' | '30' | '90'
  const expDate = new Date()
  expDate.setDate(expDate.getDate() + days)
  formData.value.expiration_date = formatDateTimeLocal(expDate.toISOString())
}

// Reset quota used for an API key
const resetQuotaUsed = async () => {
  const key = selectedKey.value
  if (!key) return
  showResetQuotaDialog.value = false
  try {
    const updatedKey = await keysAPI.update(key.id, { reset_quota: true })
    appStore.showSuccess(t('keys.quotaResetSuccess'))
    key.quota_used = updatedKey.quota_used
    if (key.status !== updatedKey.status) {
      key.status = updatedKey.status
      if (selectedKey.value?.id === key.id) {
        formData.value.status = updatedKey.status === 'active' ? 'active' : 'inactive'
      }
    }
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToResetQuota')
    appStore.showError(errorMsg)
  }
}

// Show the reset rate limit confirmation dialog, optionally selecting a key card.
const confirmResetRateLimit = (row?: ApiKey) => {
  if (row) selectedKey.value = row
  showResetRateLimitDialog.value = true
}

// Reset rate limit usage for an API key
const resetRateLimitUsage = async () => {
  if (!selectedKey.value) return
  showResetRateLimitDialog.value = false
  try {
    await keysAPI.update(selectedKey.value.id, { reset_rate_limit_usage: true })
    appStore.showSuccess(t('keys.rateLimitResetSuccess'))
    // Refresh key data
    await loadApiKeys()
    // Update the editing key with fresh data
    const refreshedKey = apiKeys.value.find(k => k.id === selectedKey.value!.id)
    if (refreshedKey) {
      selectedKey.value = refreshedKey
    }
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToResetRateLimit')
    appStore.showError(errorMsg)
  }
}

const importToCcswitch = (row: ApiKey) => {
  const platform = row.group?.platform || 'anthropic'

  ++ccsModelsRequest
  ccsModelsController?.abort()
  ccsModelsError.value = ''
  ccsModelsLoading.value = false
  pendingCcsRow.value = row
  selectedCcsClient.value = null
  selectedCcsModel.value = ''
  ccsModelOptions.value = []
  ccsModelsUnavailable.value = false
  showCcsClientSelect.value = true
  if (platform !== 'antigravity' && platform !== 'gemini') {
    void loadCcsModels(row)
  }
}

const loadCcsModels = async (row: ApiKey) => {
  const request = ++ccsModelsRequest
  ccsModelsController?.abort()
  ccsModelsController = new AbortController()
  ccsModelsLoading.value = true
  ccsModelsUnavailable.value = false
  ccsModelsError.value = ''
  ccsModelOptions.value = []
  selectedCcsModel.value = ''
  try {
    const models = await fetchModelCatalog(publicSettings.value?.api_base_url || window.location.origin, row.key, ccsModelsController.signal)
    if (request !== ccsModelsRequest || pendingCcsRow.value?.id !== row.id) return
    ccsModelOptions.value = models
    ccsModelsUnavailable.value = models.length === 0
    if (!models.length) ccsModelsError.value = t('keys.ccsClientSelect.modelsEmpty')
    const suggestedModel = row.group?.default_model?.trim()
    selectedCcsModel.value = suggestedModel && models.includes(suggestedModel) ? suggestedModel : models[0] || ''
  } catch (error) {
    if (request !== ccsModelsRequest || pendingCcsRow.value?.id !== row.id) return
    ccsModelsUnavailable.value = true
    if (error instanceof ModelCatalogError && error.kind === 'http') {
      const detail = [error.code, error.detail].filter(Boolean).join(': ')
      ccsModelsError.value = `${t('keys.ccsClientSelect.modelsHttpError', { status: error.status })}${detail ? ` ${detail}` : ''}`
    } else {
      ccsModelsError.value = t(error instanceof ModelCatalogError && error.kind === 'network'
        ? 'keys.ccsClientSelect.modelsNetworkError' : 'keys.ccsClientSelect.modelsInvalid')
    }
  } finally {
    if (request === ccsModelsRequest && pendingCcsRow.value?.id === row.id) ccsModelsLoading.value = false
  }
}

const executeCcsImport = (
  row: ApiKey,
  clientType: CcSwitchClientType,
  model = ''
) => {
  const baseUrl = publicSettings.value?.api_base_url || window.location.origin
  const platform = row.group?.platform || 'anthropic'

  const usageScript = CC_SWITCH_USAGE_SCRIPT
  const providerName = (publicSettings.value?.site_name || 'sub2api').trim() || 'sub2api'
  const deeplink = buildCcSwitchImportDeeplink({
    baseUrl,
    platform,
    clientType,
    providerName,
    apiKey: row.key,
    usageScript,
    model: model.trim() || undefined
  })

  try {
    window.open(deeplink, '_self')

    // Check if the protocol handler worked by detecting if we're still focused
    setTimeout(() => {
      if (document.hasFocus()) {
        // Still focused means the protocol handler likely failed
        appStore.showError(t('keys.ccSwitchNotInstalled'))
      }
    }, 100)
  } catch (error) {
    appStore.showError(t('keys.ccSwitchNotInstalled'))
  }
}

const handleCcsClientSelect = (clientType: CcSwitchClientType) => {
  const option = ccsClientOptions.value.find((item) => item.clientType === clientType)
  if (!option || option.disabled) return

  selectedCcsClient.value = clientType
}

const retryLoadCcsModels = () => {
  if (pendingCcsRow.value) {
    void loadCcsModels(pendingCcsRow.value)
  }
}

const confirmCcsClientSelect = () => {
  if (!pendingCcsRow.value || !selectedCcsClient.value || !canConfirmCcsImport.value) return
  executeCcsImport(pendingCcsRow.value, selectedCcsClient.value, selectedCcsModel.value)
  const advanceTour = onboardingStore.isCurrentStep('[data-tour="ccs-confirm"]')
  closeCcsClientSelect()
  if (advanceTour) void onboardingStore.nextStep()
}

const closeCcsClientSelect = () => {
  ++ccsModelsRequest
  ccsModelsController?.abort()
  ccsModelsError.value = ''
  showCcsClientSelect.value = false
  pendingCcsRow.value = null
  selectedCcsClient.value = null
  selectedCcsModel.value = ''
  ccsModelOptions.value = []
  ccsModelsLoading.value = false
  ccsModelsUnavailable.value = false
}

function formatResetTime(resetAt: string | null): string {
  if (!resetAt) return ''
  const diff = new Date(resetAt).getTime() - now.value.getTime()
  if (diff <= 0) return t('keys.resetNow')
  const days = Math.floor(diff / 86400000)
  const hours = Math.floor((diff % 86400000) / 3600000)
  const mins = Math.floor((diff % 3600000) / 60000)
  if (days > 0) return `${days}d ${hours}h`
  if (hours > 0) return `${hours}h ${mins}m`
  return `${mins}m`
}

onMounted(() => {
  loadApiKeys()
  loadGroups()
  loadUserGroupRates()
  loadPublicSettings()
  document.addEventListener('click', closeGroupSelector)
  resetTimer = setInterval(() => { now.value = new Date() }, 60000)
})

onUnmounted(() => {
  ++ccsModelsRequest
  ccsModelsController?.abort()
  document.removeEventListener('click', closeGroupSelector)
  if (resetTimer) clearInterval(resetTimer)
  if (initialSkeletonTimer) clearTimeout(initialSkeletonTimer)
  abortController?.abort()
  abortController = null
})
</script>

<style scoped>
.keys-page-layout {
  width: 100%;
  max-width: 72rem;
  margin-inline: auto;
}

.keys-page-layout :deep(.table-scroll-container) {
  border: 0;
  border-radius: 0;
  background: transparent;
  padding: 0;
  box-shadow: none;
}

.key-form :deep(.input) {
  @apply dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100 dark:placeholder:text-dark-400;
  @apply dark:focus:border-dark-500 dark:focus:ring-white/10;
}

.key-form :deep(.select-trigger) {
  @apply dark:border-dark-700 dark:bg-dark-900 dark:text-gray-100 dark:hover:border-dark-500;
  @apply dark:focus:border-dark-500 dark:focus:ring-white/10;
}

.key-form :deep(.select-trigger-open) {
  @apply dark:border-dark-500 dark:ring-white/10;
}

.key-card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(100%, 20rem), 20rem));
  align-items: start;
  gap: 1rem;
}

.key-skeleton-block {
  animation: key-skeleton-breathe 2.8s ease-in-out infinite;
}

@keyframes key-skeleton-breathe {
  0%,
  100% {
    opacity: 0.52;
  }
  50% {
    opacity: 0.88;
  }
}

@media (prefers-reduced-motion: reduce) {
  .key-skeleton-block {
    animation: none;
    opacity: 0.68;
  }
}

.api-key-filter-shell:focus-within {
  outline: none;
}

:deep(.key-filter-select) {
  width: 2.5rem !important;
  flex: 0 0 2.5rem;
}

:deep(.key-filter-select .select-trigger) {
  width: 2.5rem;
  height: 2.5rem;
  justify-content: center;
  gap: 0;
  padding: 0;
  border-color: transparent;
  background: transparent;
  box-shadow: none;
}

:deep(.key-filter-select .select-trigger:hover) {
  filter: brightness(0.96);
}

:deep(.key-filter-select .select-trigger:focus),
:deep(.key-filter-select .select-trigger-open) {
  border-color: transparent;
  box-shadow: none;
}

:deep(.key-filter-select .select-value) {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  overflow: visible;
}

:deep(.key-filter-select .select-icon) {
  display: none;
}

:deep(.key-filter-search) {
  width: auto !important;
  min-width: 0;
  flex: 1 1 10rem;
}

:deep(.key-filter-search > div) {
  width: 100%;
}

:deep(.key-filter-search .input) {
  height: 2.75rem;
  border: 0;
  background: transparent;
  padding-top: 0;
  padding-bottom: 0;
  box-shadow: none;
}

:deep(.key-filter-search .input:focus) {
  border-color: transparent;
  box-shadow: none;
}

:deep(.keys-toolbar-button) {
  @apply dark:border-dark-700 dark:bg-dark-950 dark:hover:border-dark-600 dark:hover:bg-dark-900;
}

:global(.select-dropdown-portal.keys-filter-dropdown) {
  @apply rounded-2xl border-white/10 bg-[#0d0f11] p-1.5 text-gray-200 shadow-2xl shadow-black/50;
  backdrop-filter: blur(16px);
}

:global(.select-dropdown-portal.keys-filter-dropdown .select-options) {
  @apply space-y-1 p-0;
}

:global(.select-dropdown-portal.keys-filter-dropdown .select-option) {
  @apply rounded-xl px-2.5 py-2 text-gray-300;
  background: transparent;
}

:global(.select-dropdown-portal.keys-filter-dropdown .select-option:hover),
:global(.select-dropdown-portal.keys-filter-dropdown .select-option-focused) {
  @apply bg-white/[0.06] text-white;
}

:global(.select-dropdown-portal.keys-filter-dropdown .select-option-selected) {
  @apply bg-white/[0.08] text-white;
}

:global(.select-dropdown-portal.keys-filter-dropdown .select-option-selected:hover),
:global(.select-dropdown-portal.keys-filter-dropdown .select-option-selected.select-option-focused) {
  @apply bg-white/[0.1] text-white;
}

:global(.select-dropdown-portal.key-form-dropdown) {
  @apply rounded-2xl border-white/10 bg-[#090a0c] p-1.5 text-gray-200 shadow-2xl shadow-black/50;
  backdrop-filter: blur(16px);
}

:global(.select-dropdown-portal.key-form-dropdown .select-option) {
  @apply rounded-xl text-gray-300;
  background: transparent;
}

:global(.select-dropdown-portal.key-form-dropdown .select-option:hover),
:global(.select-dropdown-portal.key-form-dropdown .select-option-focused) {
  @apply bg-white/[0.06] text-white;
}

:global(.select-dropdown-portal.key-form-dropdown .select-option-selected) {
  @apply bg-white/[0.08] text-white;
}

:global(.select-dropdown-portal.key-form-dropdown .select-option-selected:hover),
:global(.select-dropdown-portal.key-form-dropdown .select-option-selected.select-option-focused) {
  @apply bg-white/[0.1] text-white;
}

</style>
