<template>
    <div class="profile-header">
        <div class="banner">
            <div class="banner-glow"></div>
            <span class="status-dot"></span>
        </div>
        <div class="profile-info">
            <div class="avatar-wrap">
                <div class="avatar">N</div>
                <div class="avatar-camera">
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 7h4l2-2h6l2 2h4v13H3z"/>
                        <circle cx="12" cy="13" r="4"/>
                    </svg>
                </div>
            </div>
            <div class="identity">
                <h1>Noa Ferreira</h1>
                <p class="handle">@noa.png · joined February 2026</p>
                <p class="bio">Trail runner, pixel-art hobbyist, espresso before noon.</p>
            </div>
        </div>
        <div class="actions-row">
            <div class="actions-left">
                <span class="public-badge">
                    <span class="dot"></span>
                    Public profile
                    <span class="toggle"><span class="toggle-knob"></span></span>
                </span>
                <button class="edit-btn">Edit profile</button>
                <span class="hint">Switch to private — only followers see this page</span>
            </div>
            <div class="stats">
                <div class="stat">
                    <span class="stat-num">128</span>
                    <span class="stat-label">posts</span>
                </div>
                <div class="stat">
                    <span class="stat-num">342</span>
                    <span class="stat-label">followers</span>
                </div>
                <div class="stat">
                    <span class="stat-num">198</span>
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
.banner {
    height: 160px;
    background: linear-gradient(120deg, #6d28d9, #a855f7 45%, #ec4899 100%);
    position: relative;
}
.banner-glow {
    position: absolute;
    top: -60px;
    right: 40px;
    width: 220px;
    height: 220px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.08);
}
.status-dot {
    position: absolute;
    bottom: 22px;
    right: 165px;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #22c55e;
    box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.25);
}
.profile-info {
    display: flex;
    align-items: flex-end;
    gap: 20px;
    padding: 0 32px;
    margin-top: -50px;
}
.avatar-wrap {
    position: relative;
}
.avatar {
    width: 96px;
    height: 96px;
    border-radius: 50%;
    background: linear-gradient(135deg, #d946ef, #ec4899);
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