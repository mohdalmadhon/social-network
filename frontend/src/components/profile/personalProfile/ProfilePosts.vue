<template>
    <section class="profile-section">
        <div class="section-heading">
            <div>
                <p class="eyebrow">ACTIVITY</p>
                <h2>Posts</h2>
            </div>

            <button class="primary-button">Create post</button>
        </div>

        <article class="post-card">
            <div class="post-header">
                <div class="small-avatar">A</div>

                <div>
                    <strong>Alex Anderson</strong>
                    <span>@alexanderson · 2h</span>
                </div>
            </div>

            <p class="post-text">
                Just finished working on something I've been thinking
                about for a while. There is something really satisfying
                about finally seeing an idea become real.
            </p>

            <div class="post-actions">
                <button>♡ 24</button>
                <button>○ 8 comments</button>
                <button>↗ Share</button>
            </div>
        </article>

        <article class="post-card">
            <div class="post-header">
                <div class="small-avatar">A</div>

                <div>
                    <strong>Alex Anderson</strong>
                    <span>@alexanderson · 1d</span>
                </div>
            </div>

            <p class="post-text">
                Sometimes keeping things simple is the hardest part.
            </p>

            <div class="post-actions">
                <button>♡ 41</button>
                <button>○ 12 comments</button>
                <button>↗ Share</button>
            </div>
        </article>
    </section>
</template>

<style scoped>
.profile-section {
    scroll-margin-top: 100px;
    margin-bottom: 35px;
}

.section-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 18px;
}

.eyebrow {
    margin: 0 0 5px;
    color: #a855f7;
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    color: #fff;
    font-size: 20px;
    font-weight: 700;
}

.primary-button {
    padding: 9px 16px;
    border: none;
    border-radius: 8px;
    background: #a855f7;
    color: white;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
}

.post-card {
    margin-bottom: 16px;
    padding: 20px;
    border: 1px solid #232332;
    border-radius: 16px;
    background: #12121c;
}

.post-header {
    display: flex;
    align-items: center;
    gap: 12px;
}

.small-avatar {
    width: 42px;
    height: 42px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: #a855f7;
    color: white;
    font-weight: 700;
    font-size: 16px;
}

.post-header strong,
.post-header span {
    display: block;
}

.post-header strong {
    color: #fff;
    font-size: 13px;
    font-weight: 600;
}

.post-header span {
    margin-top: 3px;
    color: #6b6b7d;
    font-size: 11px;
}

.post-text {
    margin: 18px 0;
    color: #c4c4d4;
    font-size: 13px;
    line-height: 1.7;
}

.post-actions {
    display: flex;
    gap: 16px;
    padding-top: 14px;
    border-top: 1px solid #232332;
}

.post-actions button {
    border: 0;
    background: transparent;
    color: #8b8b9e;
    font-size: 12px;
    cursor: pointer;
}

.post-actions button:hover {
    color: #a855f7;
}
</style>