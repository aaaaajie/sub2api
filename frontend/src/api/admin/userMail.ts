import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'
import type { UserListFilters } from './users'

export type UserMailTemplate = 'maintenance' | 'service' | 'account'
export type UserMailStatus = 'pending' | 'sending' | 'accepted' | 'failed' | 'skipped' | 'uncertain'

export type UserMailAudience =
  | { user_ids: number[]; filters?: never }
  | { filters: UserListFilters; user_ids?: never }

export interface UserMailDraft {
  template: UserMailTemplate
  subject: string
  body: string
  action_label: string
  action_url: string
  locale: string
}

export interface UserMailPreview {
  html: string
  user_ids: number[]
  skipped: number
}

export interface UserMailJob {
  id: string
  subject: string
  is_test: boolean
  created_at: string
  pending: number
  accepted: number
  failed: number
  skipped: number
  uncertain: number
}

export interface UserMailRecipient {
  user_id: number
  email: string
  status: UserMailStatus
  reason?: string
}

export interface UserMailJobDetail extends UserMailJob {
  recipients: UserMailRecipient[]
}

const base = '/admin/user-mail'

export async function previewUserMail(draft: UserMailDraft, audience: UserMailAudience) {
  const { data } = await apiClient.post<UserMailPreview>(`${base}/preview`, { draft, ...audience })
  return data
}

export async function queueUserMail(id: string, draft: UserMailDraft, userIds: number[], isTest: boolean) {
  const { data } = await apiClient.post<UserMailJob>(base, { id, draft, user_ids: userIds, is_test: isTest })
  return data
}

export async function listUserMail(page: number) {
  const { data } = await apiClient.get<PaginatedResponse<UserMailJob>>(base, { params: { page, page_size: 10 } })
  return data
}

export async function getUserMail(id: string) {
  const { data } = await apiClient.get<UserMailJobDetail>(`${base}/${id}`)
  return data
}

export async function retryUserMailRecipient(jobId: string, userId: number) {
  const { data } = await apiClient.post<UserMailJobDetail>(`${base}/${jobId}/recipients/${userId}/retry`)
  return data
}
