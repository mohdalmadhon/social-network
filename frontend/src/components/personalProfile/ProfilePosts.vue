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
    margin-bottom: var(--space-6);
}

.section-heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-4);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-orange);
    font-family: var(--font-meta);
    font-size: 9px;
    font-weight: 600;
    letter-spacing: 2px;
    text-transform: uppercase;
}

h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 29px;
}

.primary-button {
    padding: var(--space-3) var(--space-4);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: var(--gradient-sunset);
    color: var(--color-text);
    font-family: var(--font-meta);
    font-size: 10px;
    font-weight: 600;
    cursor: pointer;
    box-shadow: var(--shadow-raised);
    transition: transform 0.1s;
}

.primary-button:hover {
    transform: translateY(-1px);
}

.primary-button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.post-card {
    margin-bottom: var(--space-5);
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.post-header {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

.small-avatar {
    width: 42px;
    height: 42px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--color-border);
    border-radius: 50%;
    background: var(--gradient-cyber);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 18px;
}

.post-header strong,
.post-header span {
    display: block;
}

.post-header strong {
    color: var(--color-text);
    font-size: 12px;
}

.post-header span {
    margin-top: var(--space-1);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 9px;
}

.post-text {
    margin: var(--space-5) 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 14px;
    line-height: 1.7;
}

.post-actions {
    display: flex;
    gap: var(--space-3);
    padding-top: var(--space-4);
    border-top: 1px solid var(--color-border);
}

.post-actions button {
    border: 0;
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 9px;
    cursor: pointer;
    transition: color 0.15s;
}

.post-actions button:nth-child(1):hover {
    color: var(--color-coral);
}

.post-actions button:nth-child(2):hover {
    color: var(--color-cyan);
}

.post-actions button:nth-child(3):hover {
    color: var(--color-lime);
}

.post-actions button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}
</style>