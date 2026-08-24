<script setup>
import { userData } from '@/stores/userData'

function formatDate(date) {
    if (!date) return ''

    return new Date(date).toLocaleDateString('en-GB', {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
    })
}

function formatJoinedDate(date) {
    if (!date) return ''

    return new Date(date).toLocaleDateString('en-GB', {
        month: 'long',
        year: 'numeric'
    })
}
</script>

<template>
    <div class="about-card">
        <div class="about-header">
            <h2>About</h2>
        </div>

        <ul class="about-list">
            <li>
                <svg
                    viewBox="0 0 24 24"
                    width="16"
                    height="16"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <rect x="3" y="5" width="18" height="14" rx="2" />
                    <path d="M3 7l9 6 9-6" />
                </svg>

                <span class="about-value">
                    {{ userData.email }}
                </span>
            </li>

            <li>
                <svg
                    viewBox="0 0 24 24"
                    width="16"
                    height="16"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <rect x="4" y="4" width="16" height="16" rx="2" />
                    <path d="M4 10h16M9 4v6M15 4v6" />
                </svg>

                <span class="about-value">
                    Born {{ formatDate(userData.dob) }}
                </span>
            </li>

            <li>
                <svg
                    viewBox="0 0 24 24"
                    width="16"
                    height="16"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <path d="M9 11l3 3 8-8" />
                    <path d="M20 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2h11" />
                </svg>

                <span class="about-value">
                    Joined {{ formatJoinedDate(userData.createdAt) }}
                </span>
            </li>
        </ul>

        <div class="divider"></div>

        <p class="bio-text">
            {{ userData.about || 'No bio added yet.' }}
        </p>

        <p class="joke-text">
            Password never leaves the auth layer — bcrypt hashed.
        </p>
    </div>
</template>

<style scoped>
.about-card {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    background: #12121c;
    border: 1px solid #232332;
    border-radius: 16px;
    padding: 20px;
}

.about-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 18px;
}

.about-header h2 {
    color: #fff;
    font-size: 15px;
    font-weight: 700;
    margin: 0;
}

.visibility-badge {
    flex-shrink: 0;
    background: #1c1c2a;
    color: #8b8b9e;
    font-size: 11px;
    padding: 5px 10px;
    border-radius: 12px;
}

.about-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.about-list li {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    min-width: 0;
    color: #c4c4d4;
    font-size: 13px;
    line-height: 1.4;
}

.about-list svg {
    color: #6b6b7d;
    flex-shrink: 0;
    margin-top: 1px;
}

.about-value {
    min-width: 0;
    overflow-wrap: anywhere;
    word-break: break-word;
}

.divider {
    height: 1px;
    background: #232332;
    margin: 20px 0;
}

.bio-text {
    color: #c4c4d4;
    font-size: 13px;
    line-height: 1.6;
    margin: 0;
    overflow-wrap: anywhere;
}

.joke-text {
    color: #4b4b5d;
    font-size: 11px;
    line-height: 1.5;
    margin: 20px 0 0;
}

@media (max-width: 700px) {
    .about-card {
        padding: 16px;
    }
}
</style>