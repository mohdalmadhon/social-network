import { ref, computed, onMounted } from "vue";
import {
  getGroups,
  createGroupApi,
  groupJoinRequest,
  undoJoinGroup,
  deleteGroupApi,
} from "@/api/groups/Groups";

export function useGroups() {
  const activeTab = ref("discover");
  const searchInputValue = ref("");
  const modalStatus = ref(false);

  const AllGroupsdata = ref([]);
  onMounted(async () => {
    try {
      const result = await getGroups();

      AllGroupsdata.value = result.groups;
    } catch (error) {
      console.error(error);
    }
  });

  const filteredGroup = computed(() => {
    let groups = AllGroupsdata.value;

    if (activeTab.value === "discover") {
      groups = groups.filter((group) => !group.isMember);
    }

    if (activeTab.value === "my-groups") {
      groups = groups.filter((group) => group.isMember);
    }

    if (searchInputValue.value.trim()) {
      const search = searchInputValue.value.trim().toLowerCase();

      groups = groups.filter((group) => {
        return group.title.toLowerCase().trim().includes(search);
      });
    }

    return groups;
  });

  function tabChanged(tab) {
    activeTab.value = tab;
  }

  function userSearchInput(input) {
    searchInputValue.value = input;
  }

  async function toggleJoinRequest(id) {
    const groupToModify = AllGroupsdata.value.find((group) => group.id === id);

    if (!groupToModify) {
      return;
    }

    if (!groupToModify.isRequested) {
      const result = await groupJoinRequest(id);
      if (result?.status) {
        groupToModify.isRequested = true;
      }
    } else {
      const result = await undoJoinGroup(id);
      if (result?.status) {
        groupToModify.isRequested = false;
      }
    }
  }

  function openModal() {
    modalStatus.value = true;
  }

  function closeModal() {
    modalStatus.value = false;
  }

  async function createGroup(data) {
    try {
      const result = await createGroupApi(data);

      const group = result.group;

      AllGroupsdata.value.push(group);

      closeModal();
    } catch (error) {
      console.error(error);
    }
  }

  async function deleteGroup(groupID) {
    try {
      const result = await deleteGroupApi(groupID);

      if (result?.status) {
        AllGroupsdata.value = AllGroupsdata.value.filter(
          (group) => group.id !== groupID,
        );
      }
    } catch (error) {
      console.error(error);
    }
  }

  return {
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
    deleteGroup,
  };
}
