import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import UserMailDialog from '../UserMailDialog.vue'

const { previewUserMail, queueUserMail } = vi.hoisted(() => ({
  previewUserMail: vi.fn(),
  queueUserMail: vi.fn(),
}))

vi.mock('@/api/admin/userMail', () => ({
  previewUserMail,
  queueUserMail,
  getUserMail: vi.fn(),
  listUserMail: vi.fn(),
  retryUserMailRecipient: vi.fn(),
}))

vi.mock('vue-i18n', async (importOriginal) => ({
  ...await importOriginal<typeof import('vue-i18n')>(),
  useI18n: () => ({ t: (key: string) => key, locale: { value: 'zh' } }),
}))

function mountDialog(props: { selectedUserIds: number[]; filters: { api_key_group_id?: number }; singleEmail?: string }) {
  return mount(UserMailDialog, {
    props,
    global: {
      stubs: { BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' }, Icon: true },
    },
  })
}

async function setDraft(dialog: ReturnType<typeof mountDialog>) {
  await dialog.get('input[placeholder="admin.userMail.subjectPlaceholder"]').setValue('维护通知')
  await dialog.get('textarea').setValue('维护时间为 02:00。')
}

async function preview(dialog: ReturnType<typeof mountDialog>) {
  const button = dialog.findAll('button').find(button => button.text() === 'admin.userMail.previewAction')
  await button!.trigger('click')
  await flushPromises()
}

describe('user mail recipient scopes', () => {
  beforeEach(() => {
    previewUserMail.mockReset()
    queueUserMail.mockReset()
    previewUserMail.mockResolvedValue({ html: '<p>维护时间</p>', user_ids: [7], skipped: 0 })
  })

  it('clears the preview when changing scopes and sends an empty filter for all site users', async () => {
    const dialog = mountDialog({ selectedUserIds: [7], filters: { api_key_group_id: 2 } })
    const scope = dialog.get('select')
    expect(scope.element.value).toBe('selected')
    await setDraft(dialog)
    await preview(dialog)
    expect(previewUserMail).toHaveBeenLastCalledWith(expect.any(Object), { user_ids: [7] })
    expect(dialog.find('iframe').exists()).toBe(true)

    await scope.setValue('filtered')
    expect(dialog.find('iframe').exists()).toBe(false)
    await preview(dialog)
    expect(previewUserMail).toHaveBeenLastCalledWith(expect.any(Object), { filters: { api_key_group_id: 2 } })

    await scope.setValue('all')
    expect(dialog.find('iframe').exists()).toBe(false)
    expect(dialog.findAll('button').map(button => button.text())).not.toContain('admin.userMail.sendTo')
    await preview(dialog)
    expect(previewUserMail).toHaveBeenLastCalledWith(expect.any(Object), { filters: {} })
    expect(queueUserMail).not.toHaveBeenCalled()
    dialog.unmount()
  })

  it('defaults to all users when the list has no filters or selections', () => {
    const dialog = mountDialog({ selectedUserIds: [], filters: {} })
    expect(dialog.get('select').element.value).toBe('all')
    dialog.unmount()
  })

  it('keeps single-user mail locked to that recipient', async () => {
    const dialog = mountDialog({ selectedUserIds: [7], filters: { api_key_group_id: 2 }, singleEmail: 'user@example.test' })
    expect(dialog.find('select').exists()).toBe(false)
    expect(dialog.text()).toContain('user@example.test')
    await setDraft(dialog)
    await preview(dialog)
    expect(previewUserMail).toHaveBeenLastCalledWith(expect.any(Object), { user_ids: [7] })
    dialog.unmount()
  })
})
