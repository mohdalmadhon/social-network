import { ref, computed } from "vue";

export function useGroups() {
  const activeTab = ref("discover");
  const searchInputValue = ref("");
  const modalStatus = ref(false);

  const AllGroupsdata = ref([
    {
      id: 1,
      title: "friends",
      description: "it our privite group.",
      memberCount: 7,
      isMember: false,
      isRequested: false,
    },
    {
      id: 20,
      title: "work",
      description: "it our work group.",
      memberCount: 77,
      isMember: false,
      isRequested: false,
    },
    {
      id: 2,
      title: "reboot",
      description: "it our reboot group.",
      memberCount: 11,
      isMember: true,
      isRequested: false,
    },
  ]);

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

  function toggleJoinRequest(id) {
    const groupToModify = AllGroupsdata.value.find((group) => group.id === id);

    if (groupToModify) {
      groupToModify.isRequested = !groupToModify.isRequested;
    }
  }

  function openModal() {
    modalStatus.value = true;
  }

  function closeModal() {
    modalStatus.value = false;
  }

  function createGroup(data) {
    AllGroupsdata.value.push({
      id: 99,
      title: data.title,
      description: data.description,
      memberCount: 1,
      isMember: true,
      isRequested: false,
    });
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
  };
}
