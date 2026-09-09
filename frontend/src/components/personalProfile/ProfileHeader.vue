<script setup>
import { ref } from 'vue';
import { requestFollow } from '@/api/users/profiles';
import { useRoute } from 'vue-router';
import { addNotification } from '@/data/notifications';

const route = useRoute();

const props = defineProps({
    addEdit: Boolean,
    firstName: String,
    lastName: String,
    username: String,
    bio: String,
    avatarPath: String,
    numOfPosts: Number,
    numOfFollowing: Number,
    numOfFollowers: Number,
    isFollowing: Number
});

const emit = defineEmits([
    'follow',
    'unfollow',
    'cancel-request'
]);

const followingStatus = ref(props.isFollowing);

async function handleFollow() {
    const id = route.query.id;

    try {
        const result = await requestFollow(id, "POST");

        if (result.status) {
            followingStatus.value = result.followStatus;
            emit('follow');
        } else {
            addNotification("could not follow user", 'error')
        }
    } catch (err) {
        addNotification("could not follow user", 'error')
        console.error(err);
    }
}

async function handleRemoveFollow() {
    const id = route.query.id;
    const oldStatus = followingStatus.value;

    try {
        const result = await requestFollow(id, "DELETE");

        if (result.status) {
            followingStatus.value = result.followStatus;
            addNotification("could not unfollow user", 'error')
            if (oldStatus === 0) {
                emit('cancel-request');
            } else if (oldStatus === 1) {
                emit('unfollow');
            }
        }
    } catch (err) {
        addNotification("could not unfollow user", 'error')
        console.error(err);
    }
}
</script>

<template>
    <section class="profile-header">
        <div class="cover">
            <div class="cover-grid"></div>
        </div>

        <div class="profile-information">
            <div class="avatar">
                <img
                    
                    v-if="props.avatarPath"
                    :src="props.avatarPath"
                    alt="Profile avatar"
                >
            </div>

            <div class="profile-details">
                <div class="name-row">
                    <div>
                        <h1>
                            {{ props.firstName }} {{ props.lastName }}
                        </h1>

                        <p class="username">
                            {{ props.username || '' }}
                        </p>
                    </div>

                    <a
                        v-if="props.addEdit"
                        href="/me/edit"
                        class="edit-button"
                    >
                        Edit profile
                    </a>

                    <button
                        v-else-if="followingStatus === -1"
                        class="relationship-button follow"
                        @click="handleFollow"
                    >
                        Follow
                    </button>

                    <button
                        v-else-if="followingStatus === 0"
                        class="relationship-button requested"
                        @click="handleRemoveFollow"
                    >
                        Requested
                    </button>

                    <button
                        v-else-if="followingStatus === 1"
                        class="relationship-button following"
                        @click="handleRemoveFollow"
                    >
                        Following
                    </button>
                </div>

                <p class="about">
                    {{ props.bio }}
                </p>

                <div class="profile-stats">
                    <span>
                        <strong>{{ props.numOfPosts }}</strong> Posts
                    </span>

                    <span>
                        <strong>{{ props.numOfFollowing }}</strong> Following
                    </span>

                    <span>
                        <strong>{{ props.numOfFollowers }}</strong> Followers
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

.cover {
    position: relative;
    height: 150px;
    z-index: -1;
    overflow: hidden;
    background: var(--gradient-aurora);
    border-bottom: 1px solid var(--color-border);
}

.cover::before,
.cover::after {
    position: absolute;
    content: "";
    border: 1px solid var(--color-border);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.cover::before {
    width: 150px;
    height: 150px;
    right: 100px;
    top: -70px;
    transform: rotate(25deg);
}

.cover::after {
    width: 80px;
    height: 80px;
    left: 120px;
    bottom: -40px;
    transform: rotate(45deg);
}

.cover-grid {
    position: absolute;
    inset: 0;
    opacity: 0.2;
    background-image:
        linear-gradient(rgb(11 13 31 / 60%) 1px, transparent 1px),
        linear-gradient(90deg, rgb(11 13 31 / 60%) 1px, transparent 1px);
    background-size: 25px 25px;
}

.profile-information {
    display: flex;
    gap: var(--space-6);
    padding: 0 var(--space-7) var(--space-6);
}

.avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 50%;
}

.avatar {
    flex-shrink: 0;
    width: 150px;
    height: 150px;
    margin-top: -75px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 4px solid var(--color-surface);
    outline: 2px solid var(--color-cyan);
    border-radius: 50%;
    background: var(--gradient-cyber);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 65px;
    box-shadow: var(--shadow-raised);
}

.profile-details {
    width: 100%;
    padding-top: var(--space-5);
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
    font-size: 36px;
    line-height: 1;
}

.username {
    margin: var(--space-2) 0 0;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 10px;
}

.edit-button,
.relationship-button {
    flex-shrink: 0;
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: 10px;
    font-weight: 600;
    text-decoration: none;
    cursor: pointer;
    transition: transform 0.1s, box-shadow 0.1s, background 0.15s;
}

.edit-button:hover,
.relationship-button:hover {
    transform: translateY(-1px);
    box-shadow: var(--shadow-raised);
}

.edit-button:focus-visible,
.relationship-button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.relationship-button.follow {
    background: var(--gradient-cyber);
    border-color: transparent;
    color: var(--color-text);
}

.relationship-button.requested {
    background: var(--color-surface-amber);
    border-color: var(--color-amber);
    color: var(--color-amber-soft);
}

.relationship-button.following {
    background: var(--color-surface-violet);
    border-color: var(--color-violet);
    color: var(--color-violet-soft);
}

.relationship-button:active {
    transform: translateY(1px);
    box-shadow: none;
}

.about {
    max-width: 650px;
    margin: var(--space-5) 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 14px;
    line-height: 1.6;
}

.profile-stats {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-5);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 10px;
}

.profile-stats span:nth-child(1) strong {
    color: var(--color-cyan);
}

.profile-stats span:nth-child(2) strong {
    color: var(--color-violet);
}

.profile-stats span:nth-child(3) strong {
    color: var(--color-coral);
}

.profile-stats strong {
    font-size: 12px;
}

@media (max-width: 650px) {
    .cover {
        height: 150px;
    }

    .profile-information {
        display: block;
        padding: 0 var(--space-5) var(--space-5);
    }

    .avatar {
        width: 115px;
        height: 115px;
        margin-top: -58px;
        font-size: 48px;
    }

    .profile-details {
        padding-top: var(--space-5);
    }

    h1 {
        font-size: 29px;
    }

    .name-row {
        align-items: center;
    }

    .profile-stats {
        gap: var(--space-3);
    }
}
</style>