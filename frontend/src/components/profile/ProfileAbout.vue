<template>
    <div class="about-card">
        <div class="about-header">
            <h2>About Noa</h2>
            <span class="visibility-badge">everyone</span>
        </div>
        <ul class="about-list">
            <li>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="5" width="18" height="14" rx="2"/>
                    <path d="M3 7l9 6 9-6"/>
                </svg>
                noa@orbit.app
            </li>
            <li>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="4" y="4" width="16" height="16" rx="2"/>
                    <path d="M4 10h16M9 4v6M15 4v6"/>
                </svg>
                Born 14 March 1999
            </li>
            <li>
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M9 11l3 3 8-8"/>
                    <path d="M20 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2h11"/>
                </svg>
                Joined February 2026
            </li>
        </ul>
        <div class="divider"></div>
        <p class="bio-text">
            Trail runner, pixel-art hobbyist, espresso before noon. Building a social network in Go for fun.
        </p>
        <p class="joke-text">Password never leaves the auth layer — bcrypt hashed.</p>
    </div>
</template>

<style scoped>
.about-card {
    background: #12121c;
    border: 1px solid #232332;
    border-radius: 16px;
    padding: 20px;
}
.about-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
}
.about-header h2 {
    color: #fff;
    font-size: 15px;
    font-weight: 700;
    margin: 0;
}
.visibility-badge {
    background: #1c1c2a;
    color: #8b8b9e;
    font-size: 11px;
    padding: 4px 10px;
    border-radius: 12px;
}
.about-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.about-list li {
    display: flex;
    align-items: center;
    gap: 10px;
    color: #c4c4d4;
    font-size: 13px;
}
.about-list svg {
    color: #6b6b7d;
    flex-shrink: 0;
}
.divider {
    height: 1px;
    background: #232332;
    margin: 18px 0;
}
.bio-text {
    color: #c4c4d4;
    font-size: 13px;
    line-height: 1.6;
    margin: 0;
}
.joke-text {
    color: #4b4b5d;
    font-size: 11px;
    margin: 20px 0 0;
}
</style>