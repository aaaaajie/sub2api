import { describe, expect, it } from 'vitest'
import { groupDisplayPlatform } from '../groupPresentation'
import { getKeyGroupProvider } from '../keyGroupProviders'

describe('group presentation', () => {
 it('categorizes a branded OpenCode group as domestic without changing its protocol', () => {
  const group = { platform: 'opencode_go' as const, display_platform: 'deepseek' as const }
  expect(groupDisplayPlatform(group)).toBe('deepseek')
  expect(getKeyGroupProvider(groupDisplayPlatform(group) || group.platform)).toBe('domestic')
  expect(group.platform).toBe('opencode_go')
 })
 it('follows the actual platform when cleared', () => {
  expect(groupDisplayPlatform({platform:'opencode_go',display_platform:''})).toBe('opencode_go')
  expect(groupDisplayPlatform(null)).toBe('')
 })
})
