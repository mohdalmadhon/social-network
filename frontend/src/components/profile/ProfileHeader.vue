<script setup>
import { computed } from 'vue'
import { userData } from '@/stores/userData'

const initial = computed(() => {
    if (userData.value.firstName) {
        return userData.value.firstName.charAt(0).toUpperCase()
    }

    if (userData.value.username) {
        return userData.value.username.charAt(0).toUpperCase()
    }

    return '?'
})

const fullName = computed(() => {
    const name = `${userData.value.firstName || ''} ${userData.value.lastName || ''}`.trim()

    return name || userData.value.username || ''
})

function formatJoinedDate(date) {
    if (!date) return ''

    return new Date(date).toLocaleDateString('en-GB', {
        month: 'long',
        year: 'numeric'
    })
}
</script>

<template>
    <div class="profile-header">
        <div class="profile-info">
            <div class="avatar-wrap">
                <div class="avatar">
                    <img v-if="userData.avatar_path" :src="userData.avatar_path" alt="Avatar">

                    <span v-else>{{ initial }}</span>
                </div>

                <div class="avatar-camera">
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 7h4l2-2h6l2 2h4v13H3z" />
                        <circle cx="12" cy="13" r="4" />
                    </svg>
                </div>
            </div>

            <div class="identity">
                <h1>{{ fullName }}</h1>

                <p class="handle">
                    @{{ userData.username || '' }} · joined {{ formatJoinedDate(userData.createdAt) }}
                </p>

                <p class="bio">
                    {{ userData.about || 'No bio added yet.' }}
                </p>
            </div>
        </div>

        <div class="actions-row">
            <div class="actions-left">
                <RouterLink to="/profile/edit">
                    <button class="edit-btn">Edit profile</button>
                </RouterLink>
            </div>

            <div class="stats">
                <div class="stat">
                    <span class="stat-num">{{ userData.numOfPosts }}</span>
                    <span class="stat-label">posts</span>
                </div>

                <div class="stat">
                    <span class="stat-num">{{ userData.numOfFollowers }}</span>
                    <span class="stat-label">followers</span>
                </div>

                <div class="stat">
                    <span class="stat-num">{{ userData.numOfFollowing }}</span>
                    <span class="stat-label">following</span>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.profile-header {
    background: #12121c;
    border: 1px solid #232332;
    border-radius: 16px;
    overflow: hidden;
}

.avatar {
    width: 96px;
    height: 96px;
    border-radius: 50%;
    background: #1d1d2b;
    border: 4px solid #12121c;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 32px;
    font-weight: 700;
    color: #fff;
    overflow: hidden;
}

.avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.profile-info {
    display: flex;
    align-items: flex-end;
    gap: 20px;
    padding: 24px 32px 0;
}

.avatar-wrap {
    position: relative;
}

.avatar {
    width: 96px;
    height: 96px;
    border-radius: 50%;
    background: #1d1d2b;
    border: 4px solid #12121c;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 32px;
    font-weight: 700;
    color: #fff;
}

.avatar-camera {
    position: absolute;
    bottom: 4px;
    right: 4px;
    width: 26px;
    height: 26px;
    border-radius: 50%;
    background: #1d1d2b;
    border: 2px solid #12121c;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #9ca3af;
}

.identity {
    padding-bottom: 12px;
}

.identity h1 {
    font-size: 22px;
    font-weight: 700;
    color: #fff;
    margin: 0;
}

.handle {
    color: #8b8b9e;
    font-size: 13px;
    margin: 2px 0 0;
}

.bio {
    color: #c4c4d4;
    font-size: 13px;
    margin: 6px 0 0;
}

.actions-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 32px 24px;
}

.actions-left {
    display: flex;
    align-items: center;
    gap: 14px;
}

.public-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    background: rgba(16, 185, 129, 0.1);
    border: 1px solid rgba(16, 185, 129, 0.3);
    color: #34d399;
    font-size: 13px;
    font-weight: 500;
    padding: 6px 10px;
    border-radius: 20px;
}

.public-badge .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #34d399;
}

.toggle {
    width: 32px;
    height: 18px;
    border-radius: 20px;
    background: #10b981;
    display: flex;
    align-items: center;
    padding: 2px;
    margin-left: 4px;
}

.toggle-knob {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: #fff;
    margin-left: auto;
}

.edit-btn {
    background: #f5f5f7;
    color: #111;
    border: none;
    font-size: 13px;
    font-weight: 600;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
}

.hint {
    color: #6b6b7d;
    font-size: 12px;
}

.stats {
    display: flex;
    gap: 32px;
}

.stat {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.stat-num {
    color: #fff;
    font-size: 18px;
    font-weight: 700;
}

.stat-label {
    color: #8b8b9e;
    font-size: 12px;
    margin-top: 2px;
}
</style>