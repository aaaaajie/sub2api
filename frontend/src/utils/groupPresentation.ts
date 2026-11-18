import type { GroupPlatform } from '@/types'

// Presentation only. Protocol selection and routing must continue using platform.
export function groupDisplayPlatform(group?: { platform: GroupPlatform; display_platform?: GroupPlatform | '' } | null): GroupPlatform | '' {
  return group?.display_platform || group?.platform || ''
}
