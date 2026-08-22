<template>
    <div class="profile-tabs">
        <span class="tab">Posts</span>
        <span class="tab active">Followers</span>
        <span class="tab">Following</span>
        <span class="tab">About</span>
    </div>
</template>

<style scoped>
.profile-tabs {
    display: flex;
    gap: 28px;
    border-bottom: 1px solid #232332;
    padding: 0 4px;
    margin-top: 20px;
}
.tab {
    color: #8b8b9e;
    font-size: 14px;
    font-weight: 500;
    padding: 12px 0;
    cursor: pointer;
}
.tab.active {
    color: #a855f7;
    font-weight: 600;
    border-bottom: 2px solid #a855f7;
}
</style>