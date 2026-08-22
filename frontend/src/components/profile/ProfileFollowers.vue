<template>
    <div class="followers-card">
        <div class="followers-header">
            <h2>Followers · 342</h2>
            <input class="search-input" type="text" placeholder="Search followers..." />
        </div>
        <div class="followers-grid">
            <div class="follower" v-for="f in followers" :key="f.name">
                <div class="follower-avatar" :style="{ background: f.color }">{{ f.initial }}</div>
                <span class="follower-name">{{ f.name }}</span>
                <button :class="['follow-btn', { active: f.followingBack }]">
                    {{ f.followingBack ? '✓ Following back' : '+ Follow back' }}
                </button>
            </div>
            <div class="more-card">
                <span class="more-num">+ 336 more</span>
                <span class="more-label">shown to followers only when private</span>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            followers: [
                { name: 'Kiko Tanaka', initial: 'K', color: '#134e3a', followingBack: true },
                { name: 'Mara Voss', initial: 'M', color: '#4a3a1e', followingBack: true },
                { name: 'Leo Marchetti', initial: 'L', color: '#1e3a4a', followingBack: false },
                { name: 'Aya Benali', initial: 'A', color: '#3a2e5c', followingBack: true },
                { name: 'Rin Kato', initial: 'R', color: '#5c2e3a', followingBack: false },
                { name: 'Dario Fontana', initial: 'D', color: '#134e3a', followingBack: true },
            ],
        };
    },
};
</script>

<style scoped>
.followers-card {
    background: #12121c;
    border: 1px solid #232332;
    border-radius: 16px;
    padding: 20px;
}
.followers-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
}
.followers-header h2 {
    color: #fff;
    font-size: 15px;
    font-weight: 700;
    margin: 0;
}
.search-input {
    background: #1a1a26;
    border: 1px solid #2a2a3a;
    border-radius: 8px;
    color: #c4c4d4;
    font-size: 13px;
    padding: 8px 14px;
    width: 200px;
    outline: none;
}
.followers-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
}
.follower {
    background: #171724;
    border: 1px solid #232332;
    border-radius: 12px;
    padding: 16px;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
}
.follower-avatar {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-weight: 700;
    font-size: 16px;
}
.follower-name {
    color: #fff;
    font-size: 13px;
    font-weight: 600;
}
.follow-btn {
    width: 100%;
    background: #1c1c2a;
    color: #8b8b9e;
    border: 1px solid #2a2a3a;
    border-radius: 8px;
    font-size: 12px;
    padding: 7px 0;
    cursor: pointer;
}
.follow-btn.active {
    background: transparent;
    color: #a855f7;
    border-color: #a855f7;
}
.more-card {
    background: #171724;
    border: 1px solid #232332;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    padding: 16px;
}
.more-num {
    color: #fff;
    font-size: 14px;
    font-weight: 700;
}
.more-label {
    color: #6b6b7d;
    font-size: 11px;
    text-align: center;
}
</style>