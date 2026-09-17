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
<<<<<<< HEAD
.edit-tabs {
    position: sticky;
    top: 64px;
    z-index: 50;
    display: flex;
    margin: clamp(16px, 3vw, 25px) 0;
    overflow-x: auto;
    border: 2px solid var(--main-color);
    border-radius: 6px;
    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);
    -webkit-overflow-scrolling: touch;
}

.tab-button {
    flex: 1 1 0;
    min-width: max-content;
    padding: clamp(11px, 2vw, 15px) clamp(12px, 2.5vw, 18px);
    border: none;
    border-right: 2px solid var(--main-color);
    background: transparent;
    color: var(--font-color-sub);
    text-align: center;
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(9px, 1.4vw, 10px);
    font-weight: 600;
    white-space: nowrap;
}

.tab-button:last-child {
    border-right: 0;
}

.tab-button:hover {
    background: var(--page-background);
    color: var(--main-color);
}

.tab-button.active {
    background: var(--input-focus);
    color: white;
}

@media (max-width: 800px) {
    .edit-tabs {
        top: 0;
    }
}
</style>
=======
.edit-tabs { display: inline-flex; width: fit-content; padding: var(--space-1); border: 1px solid var(--color-border); border-radius: var(--radius-small); background: var(--color-surface); }
.edit-tabs button { min-height: var(--touch-target); padding: 0 var(--space-4); border: 0; border-radius: var(--radius-small); background: transparent; color: var(--color-text-muted); cursor: pointer; font-weight: 600; }
.edit-tabs button.active { background: var(--color-input); color: var(--color-mint); }
@media (max-width: 480px) { .edit-tabs { display: grid; width: 100%; grid-template-columns: 1fr 1fr; } }
</style>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
