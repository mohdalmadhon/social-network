<script setup>
import { router } from '@/router/router.js'
import { useGroups } from '@/composables/useGroups.js'
import GroupsHeader from '@/components/groups/GroupsHeader.vue'
import GroupsTabs from '@/components/groups/GroupsTabs.vue'
import GroupSearch from '@/components/groups/GroupSearch.vue'
import GroupsList from '@/components/groups/GroupsList.vue'
import CreateGroupModal from '@/components/groups/CreateGroupModal.vue'

const {
  activeTab,
  searchInputValue,
  modalStatus,
  filteredGroup,
  tabChanged,
  userSearchInput,
  toggleJoinRequest,
  openModal,
  closeModal,
  createGroup,
} = useGroups()

function viewGroup(id) {
  router.push(`/groups/${id}`)
}
</script>

<template>
  <main class="groups-page">
    <div class="groups-content">
      <GroupsHeader @create-group="openModal" />

      <GroupsTabs @change-tab="tabChanged" :active-tab="activeTab" />

      <GroupSearch @update:model-value="userSearchInput" :model-value="searchInputValue" />

      <GroupsList @toggle-join-request="toggleJoinRequest" @view-group="viewGroup" :groups="filteredGroup" />
    </div>

    <CreateGroupModal @close="closeModal" @create="createGroup" :show="modalStatus" />
  </main>
</template>

<style scoped>
.groups-page {
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding: 24px 32px;
  box-sizing: border-box;
}

.groups-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

@media (max-width: 48rem) {
  .groups-page {
    padding: 20px 16px;
  }
}
</style>