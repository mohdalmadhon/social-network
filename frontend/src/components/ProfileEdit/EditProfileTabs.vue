<script setup>
defineProps({ activeTab: String })
defineEmits(['update:activeTab'])
const tabs = [{ id: 'personal', label: 'Profile' }, { id: 'additional', label: 'Additional info' }]
</script>

<template>
  <nav class="edit-tabs" aria-label="Profile setting sections">
    <button v-for="tab in tabs" :key="tab.id" type="button" :class="{ active: activeTab === tab.id }" @click="$emit('update:activeTab', tab.id)">
      {{ tab.label }}
    </button>
  </nav>
</template>

<style scoped>
.edit-tabs { display: inline-flex; width: fit-content; padding: var(--space-1); border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-surface); }
.edit-tabs button { min-height: var(--touch-target); padding: 0 var(--space-4); border: 0; border-radius: var(--radius-small); background: transparent; color: var(--color-text-muted); cursor: pointer; font-weight: 600; }
.edit-tabs button.active { background: var(--color-input); color: var(--color-mint); }
@media (max-width: 480px) { .edit-tabs { display: grid; width: 100%; grid-template-columns: 1fr 1fr; } }
</style>
