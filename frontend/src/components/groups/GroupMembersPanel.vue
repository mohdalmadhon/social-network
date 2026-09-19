<script setup>
import { ref, watch, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import Groupssearch from './Groupssearch.vue';

const props = defineProps({
    show: {
        type: Boolean,
        default: false
    },
    groupID: {
        type: [Number, String],
        required: true
    },
    members: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits(['close']);

const router = useRouter();

const search = ref('');
const memberList = ref([...props.members]);

const offset = ref(props.members.length);
const limit = 20;

const loading = ref(false);
const hasMore = ref(true);

let debounceTimer = null;
let scrollTimer = null;
let requestNumber = 0;

function getMemberId(member) {
    return member.ID ?? member.id;
}

function getMemberName(member) {
    return `${member.firstName || ''} ${member.lastName || ''}`.trim();
}

function getMemberAvatar(member) {
    const avatar = member.avatar ?? member.Avatar;

    if (!avatar) {
        return '';
    }

    if (avatar.startsWith('/')) {
        return avatar;
    }

    return `/uploads/${avatar}`;
}

function initials(name) {
    if (!name) {
        return '?';
    }

    return name
        .trim()
        .split(/\s+/)
        .slice(0, 2)
        .map(word => word[0]?.toUpperCase())
        .join('');
}

async function getMembers(reset = false) {
    if (loading.value) {
        return;
    }

    if (!reset && !hasMore.value) {
        return;
    }

    loading.value = true;

    const currentRequest = ++requestNumber;

    try {
        const currentOffset = reset ? 0 : offset.value;

        const params = new URLSearchParams({
            groupID: String(props.groupID),
            offset: String(currentOffset),
            search: search.value.trim()
        });

        const response = await fetch(`/api/group/search?${params.toString()}`);

        if (!response.ok) {
            throw new Error('Failed to get members');
        }

        const result = await response.json();

        if (currentRequest !== requestNumber) {
            return;
        }

        if (!result.status) {
            throw new Error(result.message || 'Failed to get members');
        }

        const newMembers = result.data || [];

        if (reset) {
            memberList.value = newMembers;
            offset.value = newMembers.length;
        } else {
            memberList.value.push(...newMembers);
            offset.value += newMembers.length;
        }

        hasMore.value = newMembers.length === limit;
    } catch (error) {
        console.error(error);
    } finally {
        if (currentRequest === requestNumber) {
            loading.value = false;
        }
    }
}

function handleSearch(value) {
    search.value = value;
}

function handleScroll(event) {
    if (scrollTimer) {
        return;
    }

    scrollTimer = setTimeout(() => {
        scrollTimer = null;

        const element = event.currentTarget;

        const distanceFromBottom =
            element.scrollHeight -
            element.scrollTop -
            element.clientHeight;

        if (distanceFromBottom <= 80) {
            getMembers();
        }
    }, 200);
}

function openProfile(member) {
    const id = getMemberId(member);

    if (!id) {
        return;
    }

    emit('close');

    router.push(`/user?id=${id}`);
}

watch(search, () => {
    clearTimeout(debounceTimer);

    debounceTimer = setTimeout(() => {
        offset.value = 0;
        hasMore.value = true;
        getMembers(true);
    }, 300);
});

watch(
    () => props.members,
    value => {
        if (!search.value.trim()) {
            memberList.value = [...value];
            offset.value = value.length;
            hasMore.value = value.length >= limit;
        }
    }
);

watch(
    () => props.show,
    visible => {
        if (visible && memberList.value.length === 0) {
            getMembers(true);
        }
    }
);

onBeforeUnmount(() => {
    clearTimeout(debounceTimer);
    clearTimeout(scrollTimer);
});
</script>

<template>
    <div v-if="show" class="members-panel">
        <div class="members-panel-header">
            <h3 class="members-title">
                {{ memberList.length }}
                {{ memberList.length === 1 ? 'Member' : 'Members' }}
            </h3>

            <button
                class="close-button"
                type="button"
                @click="$emit('close')"
            >
                ×
            </button>
        </div>

        <Groupssearch
            :model-value="search"
            placeholder="Search members..."
            @update:model-value="handleSearch"
        />

        <div
            class="members-list"
            @scroll="handleScroll"
        >
            <div
                v-for="member in memberList"
                :key="getMemberId(member)"
                class="member-row"
                @click="openProfile(member)"
            >
                <div class="member-avatar">
                    <img
                        v-if="getMemberAvatar(member)"
                        :src="getMemberAvatar(member)"
                        :alt="getMemberName(member)"
                        class="member-avatar-img"
                    >

                    <span
                        v-else
                        class="member-avatar-fallback"
                    >
                        {{ initials(getMemberName(member)) }}
                    </span>
                </div>

                <div class="member-info">
                    <span class="member-name">
                        {{ getMemberName(member) }}
                    </span>
                </div>
            </div>

            <div
                v-if="loading"
                class="loading-members"
            >
                Loading...
            </div>

            <p
                v-if="!loading && !memberList.length"
                class="no-members"
            >
                No members found.
            </p>
        </div>
    </div>
</template>

<style scoped>
.members-panel {
    position: absolute;
    z-index: 20;
    top: calc(100% + 10px);
    right: 24px;
    width: 320px;
    max-width: calc(100vw - 48px);
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 18px;
    border: 2px solid var(--main-color);
    border-radius: 8px;
    background: var(--bg-color);
    box-shadow: 6px 6px var(--main-color);
}

.members-panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.members-title {
    margin: 0;
    color: var(--font-color);
    font-family: "Liter", serif;
    font-size: 16px;
    font-weight: 700;
}

.close-button {
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--bg-color);
    color: var(--main-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 14px;
    line-height: 1;
    cursor: pointer;
}

.close-button:hover {
    background: var(--main-color);
    color: var(--bg-color);
}

.members-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
    max-height: 320px;
    overflow-y: auto;
}

.member-row {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 5px;
    border-radius: 6px;
    cursor: pointer;
}

.member-row:hover {
    background: var(--input-focus);
}

.member-avatar {
    flex-shrink: 0;
    width: 38px;
    height: 38px;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--input-focus);
}

.member-avatar-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.member-avatar-fallback {
    color: #fff;
    font-family: "JetBrains Mono", monospace;
    font-size: 12px;
    font-weight: 600;
}

.member-info {
    flex: 1;
    min-width: 0;
}

.member-name {
    color: var(--font-color);
    font-size: 13px;
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.loading-members {
    padding: 8px;
    text-align: center;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
}

.no-members {
    margin: 10px 0;
    text-align: center;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 12px;
}

@media (max-width: 650px) {
    .members-panel {
        right: 12px;
        left: 12px;
        width: auto;
    }
}
</style>