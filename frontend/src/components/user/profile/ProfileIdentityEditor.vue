<template>
  <div class="profile-identity-editor">
    <ProfileAvatarCard v-if="kind === 'avatar'" :user="user" embedded @saved="$emit('saved')" />
    <ProfileEditForm v-else :initial-username="user?.username || ''" embedded @saved="$emit('saved')" />
    <p v-for="hint in sourceHints.filter(h => h.key === kind)" :key="hint.key" class="mt-4 text-xs text-gray-500 dark:text-gray-400">{{ hint.text }}</p>
  </div>
</template>
<script setup lang="ts">
import type { User } from '@/types'
import ProfileAvatarCard from './ProfileAvatarCard.vue'
import ProfileEditForm from './ProfileEditForm.vue'
import { useProfileIdentity } from './useProfileIdentity'
const props = defineProps<{ kind: 'avatar' | 'username'; user: User | null; oidcProviderName: string }>()
defineEmits<{ saved: [] }>()
const { sourceHints } = useProfileIdentity(() => props.user, () => props.oidcProviderName)
</script>
<style scoped>
.profile-identity-editor :deep(.input), .profile-identity-editor :deep(.btn-secondary) { @apply bg-transparent dark:bg-transparent border-gray-200 dark:border-white/10; }
</style>
