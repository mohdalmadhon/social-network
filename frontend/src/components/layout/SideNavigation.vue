<script setup>
import { logout } from '@/api/auth/auth';
import { addNotification } from '@/data/notifications';


async function logoutHandler() {
    try {
        await logout();
    } catch (err) {
        addNotification(err, 'error');
    }
}
</script>

<template>
    <aside class="side-navigation">
        <nav>
            <a href="/">
                <span class="icon">⌂</span>
                <span class="label">Home</span>
            </a>

            <a href="/profile" class="active">
                <span class="icon">◉</span>
                <span class="label">Profile</span>
            </a>

            <a href="/friends">
                <span class="icon">♧</span>
                <span class="label">Friends</span>
            </a>

            <a href="/groups">
                <span class="icon">▦</span>
                <span class="label">Groups</span>
            </a>

            <a href="/following">
                <span class="icon">→</span>
                <span class="label">Following</span>
            </a>

            <a href="/followers">
                <span class="icon">←</span>
                <span class="label">Followers</span>
            </a>
        </nav>

        <div class="side-bottom">
            <a href="/settings">
                <span class="icon">⚙</span>
                <span class="label">Settings</span>
            </a>

            <a @click="logoutHandler">
                <span class="icon">↪</span>
                <span class="label">Log out</span>
            </a>
        </div>
    </aside>
</template>

<style scoped>
.side-navigation {
    position: sticky;
    top: 64px;
    flex-shrink: 0;
    width: clamp(72px, 16vw, 220px);
    height: calc(100vh - 64px);
    padding: clamp(16px, 2.5vw, 25px) clamp(10px, 1.5vw, 18px);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    background: var(--bg-color);
    border-right: 2px solid var(--main-color);
    box-sizing: border-box;
}

.side-navigation nav,
.side-bottom {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.side-navigation a {
    min-height: 44px;
    display: flex;
    align-items: center;
    gap: 13px;
    padding: 0 clamp(8px, 1.2vw, 13px);
    border: 2px solid transparent;
    border-radius: 5px;
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: clamp(10px, 1.4vw, 11px);
    font-weight: 600;
    overflow: hidden;
}

.side-navigation a:hover {
    border-color: var(--main-color);
    background: var(--page-background);
}

.side-navigation a.active {
    border: 2px solid var(--main-color);
    background: var(--input-focus);
    color: white;
    box-shadow: 3px 3px var(--main-color);
}

.icon {
    flex-shrink: 0;
    width: 20px;
    text-align: center;
    font-size: 16px;
}

@media (max-width: 1000px) {
    .side-navigation .label {
        display: none;
    }

    .side-navigation a {
        justify-content: center;
        padding: 0;
    }
}

@media (max-width: 800px) {
    .side-navigation {
        position: static;
        width: 100%;
        height: auto;
        flex-shrink: initial;
        padding: 8px 15px;
        border-right: 0;
        border-bottom: 2px solid var(--main-color);
        flex-direction: row;
        overflow-x: auto;
        -webkit-overflow-scrolling: touch;
    }

    .side-navigation nav,
    .side-bottom {
        flex-direction: row;
    }

    .side-navigation a {
        white-space: nowrap;
        flex-shrink: 0;
    }

    .side-navigation .label {
        display: inline;
    }

    .side-bottom {
        display: none;
    }
}
</style>