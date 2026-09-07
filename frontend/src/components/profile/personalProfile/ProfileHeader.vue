<script setup>
import { ref, watch } from 'vue';
import { requestFollow } from '@/api/users/profiles';
import { useRoute } from 'vue-router';
import { addNotification } from '@/data/notifications';

const route = useRoute();

const props = defineProps({
    addEdit: {
        type: Boolean,
        default: false
    },
    firstName: {
        type: String,
        default: ''
    },
    lastName: {
        type: String,
        default: ''
    },
    username: {
        type: String,
        default: ''
    },
    bio: {
        type: String,
        default: ''
    },
    avatarPath: {
        type: String,
        default: ''
    },
    numOfPosts: {
        type: Number,
        default: 0
    },
    numOfFollowing: {
        type: Number,
        default: 0
    },
    numOfFollowers: {
        type: Number,
        default: 0
    },
    isFollowing: {
        type: Number,
        default: -1
    }
});

const emit = defineEmits([
    'follow',
    'unfollow',
    'cancel-request'
]);

const followingStatus = ref(props.isFollowing);

watch(
    () => props.isFollowing,
    (newStatus) => {
        followingStatus.value = newStatus;
    }
);

async function handleFollow() {
    const id = route.query.id;

    if (!id) {
        addNotification('Could not follow user', 'error');
        return;
    }

    try {
        const result = await requestFollow(id, 'POST');

        if (!result.status) {
            addNotification('Could not follow user', 'error');
            return;
        }

        followingStatus.value = result.followStatus;

        if (result.followStatus === 0) {
            emit('follow');
        } else if (result.followStatus === 1) {
            emit('follow');
        }
    } catch (err) {
        console.error(err);
        addNotification('Could not follow user', 'error');
    }
}

async function handleRemoveFollow() {
    const id = route.query.id;
    const oldStatus = followingStatus.value;

    if (!id) {
        addNotification('Could not unfollow user', 'error');
        return;
    }

    try {
        const result = await requestFollow(id, 'DELETE');

        if (!result.status) {
            addNotification('Could not unfollow user', 'error');
            return;
        }

        followingStatus.value = result.followStatus;

        if (oldStatus === 0) {
            emit('cancel-request');
        } else if (oldStatus === 1) {
            emit('unfollow');
        }
    } catch (err) {
        console.error(err);
        addNotification('Could not unfollow user', 'error');
    }
}
</script>

<template>
    <section class="profile-header">
        <div class="profile-information">
            <div class="avatar">
                <img
                    v-if="avatarPath"
                    :src="avatarPath"
                    alt="Profile avatar"
                >
            </div>

            <div class="profile-details">
                <div class="name-row">
                    <div>
                        <h1>
                            {{ firstName }} {{ lastName }}
                        </h1>

                        <p class="username">
                            {{ username }}
                        </p>
                    </div>

                    <a
                        v-if="addEdit"
                        href="/me/edit"
                        class="edit-button"
                    >
                        Edit profile
                    </a>

                    <button
                        v-else-if="followingStatus === -1"
                        class="relationship-button follow"
                        type="button"
                        @click="handleFollow"
                    >
                        Follow
                    </button>

                    <button
                        v-else-if="followingStatus === 0"
                        class="relationship-button requested"
                        type="button"
                        @click="handleRemoveFollow"
                    >
                        Requested
                    </button>

                    <button
                        v-else-if="followingStatus === 1"
                        class="relationship-button following"
                        type="button"
                        @click="handleRemoveFollow"
                    >
                        Following
                    </button>
                </div>

                <p
                    v-if="bio"
                    class="about"
                >
                    {{ bio }}
                </p>

                <div class="profile-stats">
                    <span>
                        <strong>{{ numOfPosts }}</strong>
                        Posts
                    </span>

                    <span>
                        <strong>{{ numOfFollowing }}</strong>
                        Following
                    </span>

                    <span>
                        <strong>{{ numOfFollowers }}</strong>
                        Followers
                    </span>
                </div>
            </div>
        </div>
    </section>
</template>

<style scoped>
.profile-header {
    overflow: hidden;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.profile-information {
    display: flex;
    gap: var(--space-6);
    padding: var(--space-6);
}

.avatar {
    flex-shrink: 0;
    width: 9.375rem;
    height: 9.375rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 4px solid var(--color-surface);
    border-radius: 50%;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 4rem;
    box-shadow: var(--shadow-raised);
    overflow: hidden;
}

.avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 50%;
}

.profile-details {
    width: 100%;
    padding-top: 0.5rem;
}

.name-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-5);
}

h1 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 2.25rem;
    line-height: 1;
}

.username {
    margin: var(--space-2) 0 0;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
}

.edit-button,
.relationship-button {
    flex-shrink: 0;
    padding: var(--space-3) var(--space-4);
    border: none;
    border-radius: 1.5625rem;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.75rem;
    font-weight: 700;
    text-decoration: none;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.edit-button:hover,
.relationship-button:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

.relationship-button.requested {
    background: var(--color-surface-raised);
    border: 1px solid var(--color-border);
    color: var(--color-text-soft);
}

.relationship-button.following {
    background: var(--color-surface-raised);
    border: 1px solid var(--color-violet);
    color: var(--color-text);
}

.relationship-button:active {
    transform: translateY(0);
}

.about {
    max-width: 40.625rem;
    margin: var(--space-4) 0;
    color: var(--color-text-soft);
    font-size: 0.875rem;
    line-height: 1.6;
}

.profile-stats {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-5);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
}

.profile-stats strong {
    color: var(--color-text);
    font-size: 0.75rem;
}

@media (max-width: 40.625rem) {
    .profile-information {
        display: block;
        padding: var(--space-5);
    }

    .avatar {
        width: 7.1875rem;
        height: 7.1875rem;
        font-size: 3rem;
    }

    .name-row {
        align-items: center;
    }

    .profile-stats {
        gap: var(--space-3);
    }
}
</style>
