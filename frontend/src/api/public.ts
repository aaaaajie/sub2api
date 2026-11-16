import { apiClient } from './client'

export interface PublicActivityTrendPoint {
  date: string
  activity_percent: number
}

export interface PublicActivityTrendResponse {
  trend: PublicActivityTrendPoint[]
}

export async function getPublicActivityTrend(): Promise<PublicActivityTrendResponse> {
  const { data } = await apiClient.get<PublicActivityTrendResponse>('/public/recent-usage-activity')
  return data
}
