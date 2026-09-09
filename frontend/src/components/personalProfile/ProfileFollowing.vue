<template>
    <section class="profile-section">
        <div class="section-heading">
            <p class="eyebrow">NETWORK</p>
            <h2>Following</h2>
        </div>

        <div class="list">
            <article class="list-item">
                <div class="list-avatar">T</div>
                <div>
                    <strong>Tech Daily</strong>
                    <span>@techdaily</span>
                </div>
                <button>Following</button>
            </article>

            <article class="list-item">
                <div class="list-avatar">C</div>
                <div>
                    <strong>Creative Corner</strong>
                    <span>@creativecorner</span>
                </div>
                <button>Following</button>
            </article>

            <article class="list-item">
                <div class="list-avatar">D</div>
                <div>
                    <strong>Design Weekly</strong>
                    <span>@designweekly</span>
                </div>
                <button>Following</button>
            </article>
        </div>
    </section>
</template>

<style scoped>
.profile-section {
    scroll-margin-top: 6.25rem;
}

.section-heading {
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-magenta);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.8125rem;
    color: var(--color-text);
}

.list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

.list-item {
    display: flex;
    align-items: center;
    gap: var(--space-4);
    padding: var(--space-4) var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
}

.list-avatar {
    width: 3rem;
    height: 3rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: var(--gradient-sunset);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 1.25rem;
}

.list-item div:nth-child(2) {
    flex: 1;
}

.list-item strong,
.list-item span {
    display: block;
}

.list-item strong {
    color: var(--color-text);
    font-size: 0.75rem;
}

.list-item span {
    margin-top: 0.1875rem;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.5rem;
}

.list-item button {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-violet);
    border-radius: 1.5625rem;
    background: var(--color-surface-violet);
    color: var(--color-violet-soft);
    font-family: var(--font-body);
    font-size: 0.5rem;
    font-weight: 700;
    cursor: pointer;
    transition: filter 0.15s ease;
}

.list-item button:hover {
    filter: brightness(1.15);
}
</style>
