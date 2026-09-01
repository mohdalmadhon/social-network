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
    border: 1px solid #232332;
    border-radius: 16px;
    background: #12121c;
}

.cover {
    position: relative;
    height: 150px;
    z-index: -1;
    overflow: hidden;
    background: #171724;
    border-bottom: 1px solid #232332;
}

.cover::before,
.cover::after {
    position: absolute;
    content: "";
    border-radius: 50%;
}

.cover::before {
    width: 220px;
    height: 220px;
    right: -60px;
    top: -110px;
    background: rgba(168, 85, 247, 0.14);
}

.cover::after {
    width: 140px;
    height: 140px;
    left: 60px;
    bottom: -90px;
    background: rgba(16, 185, 129, 0.08);
}

.cover-grid {
    position: absolute;
    inset: 0;
    opacity: 0.15;
    background-image:
        linear-gradient(rgba(168, 85, 247, 0.6) 1px, transparent 1px),
        linear-gradient(90deg, rgba(168, 85, 247, 0.6) 1px, transparent 1px);
    background-size: 25px 25px;
}

.profile-information {
    display: flex;
    gap: 28px;
    padding: 0 35px 30px;
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
    border: 4px solid #12121c;
    border-radius: 50%;
    background: #1d1d2b;
    color: white;
    font-size: 65px;
    font-weight: 700;
}

.profile-details {
    width: 100%;
    padding-top: 22px;
}

.name-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 20px;
}

h1 {
    margin: 0;
    color: #fff;
    font-size: 22px;
    font-weight: 700;
    line-height: 1.2;
}

.username {
    margin: 4px 0 0;
    color: #8b8b9e;
    font-size: 13px;
}

.edit-button,
.relationship-button {
    flex-shrink: 0;
    padding: 8px 16px;
    border: none;
    border-radius: 8px;
    background: #f5f5f7;
    color: #111;
    font-size: 13px;
    font-weight: 600;
    text-decoration: none;
    cursor: pointer;
    transition: opacity 0.15s;
}

.edit-button:hover,
.relationship-button:hover {
    opacity: 0.85;
}

.relationship-button.follow {
    background: #a855f7;
    color: white;
}

.relationship-button.requested {
    background: #1c1c2a;
    color: #8b8b9e;
    border: 1px solid #2a2a3a;
}

.relationship-button.following {
    background: transparent;
    color: #a855f7;
    border: 1px solid #a855f7;
}

.relationship-button:active {
    opacity: 0.7;
}

.about {
    max-width: 650px;
    margin: 14px 0;
    color: #c4c4d4;
    font-size: 13px;
    line-height: 1.6;
}

.profile-stats {
    display: flex;
    flex-wrap: wrap;
    gap: 22px;
    color: #8b8b9e;
    font-size: 12px;
}

.profile-stats strong {
    color: #fff;
    font-size: 14px;
}

@media (max-width: 650px) {
    .cover {
        height: 150px;
    }

    .profile-information {
        display: block;
        padding: 0 20px 25px;
    }

    .avatar {
        width: 115px;
        height: 115px;
        margin-top: -58px;
        font-size: 48px;
    }

    .profile-details {
        padding-top: 20px;
    }

    h1 {
        font-size: 19px;
    }

    .name-row {
        align-items: center;
    }

    .profile-stats {
        gap: 12px;
    }
}
</style>
