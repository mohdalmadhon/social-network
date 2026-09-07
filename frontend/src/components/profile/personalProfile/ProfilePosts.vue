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
    scroll-margin-top: 6.25rem;
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
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.8125rem;
    color: var(--color-text);
}

.primary-button {
    padding: var(--space-3) var(--space-4);
    border: none;
    border-radius: 1.5625rem;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.75rem;
    font-weight: 700;
    cursor: pointer;
    transition: transform 0.15s ease, filter 0.15s ease;
}

.primary-button:hover {
    filter: brightness(1.08);
    transform: translateY(-1px);
}

.post-card {
    margin-bottom: var(--space-4);
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.post-header {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

.small-avatar {
    width: 2.625rem;
    height: 2.625rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: var(--gradient-action);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 1.125rem;
}

.post-header strong,
.post-header span {
    display: block;
    color: var(--color-text);
}

.post-header strong {
    font-size: 0.75rem;
}

.post-header span {
    margin-top: 0.1875rem;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.5625rem;
}

.post-text {
    margin: var(--space-5) 0;
    color: var(--color-text-soft);
    font-size: 0.875rem;
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
    font-size: 0.5625rem;
    cursor: pointer;
}

.post-actions button:hover {
    color: var(--color-violet);
}
</style>
