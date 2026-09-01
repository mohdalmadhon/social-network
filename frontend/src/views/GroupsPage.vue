<script setup>
import { ref } from 'vue'
import GroupsHeader from '@/components/groups/GroupsHeader.vue'
import GroupsTabs from '@/components/groups/GroupsTabs.vue'
import GroupSearch from '@/components/groups/GroupSearch.vue'
import GroupsList from '@/components/groups/GroupsList.vue'
import CreateGroupModal from '@/components/groups/CreateGroupModal.vue'

const activeTab = ref('discover')
const searchInputValue = ref('')
const modalStatus = ref(false)

function createGroupModal() {
  console.log("create group clicked!")

  modalStatus.value = true
}

function tabChanged(tab) {
  console.log(tab, "tab changed!")

  activeTab.value = tab
}

function userSearchInput(input) {
  console.log(input)
}

function joinGroup(id) {
  console.log("user wants join group that has ID:", id)
}

function viewGroup(id) {
  console.log("user wants view group that has ID:", id)
}

function closeModal() {
  console.log("modal closed")

  modalStatus.value = false
}

function createGroup(data) {
  console.log("Group created")
  console.log("the title is:", data.title)
  console.log("the description is:", data.description)
}

const Groupsdata = [
  {
    id: 1,
    title: "friends",
    description: "it our privite group.",
    memberCount: 7,
    isMember: false,
    isRequested: false,
  },
  {
    id: 2,
    title: "work",
    description: "it our work group.",
    memberCount: 77,
    isMember: false,
    requested: false,
  }
]
</script>

<template>
  <main class="groups-page">
    <GroupsHeader @create-group="createGroupModal" />

    <GroupsTabs @change-tab="tabChanged" :active-tab="activeTab" />

    <GroupSearch @update:model-value="userSearchInput" :model-value="searchInputValue" />

    <GroupsList @join-group="joinGroup" @view-group="viewGroup" :groups="Groupsdata" />

    <CreateGroupModal @close="closeModal" @create="createGroup" :show="modalStatus" />
  </main>
</template>

<style scoped></style>