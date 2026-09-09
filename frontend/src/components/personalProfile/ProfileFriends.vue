<template>
    <section class="profile-section">
        <div class="section-heading">
            <div>
                <p class="eyebrow">PEOPLE</p>
                <h2>Friends</h2>
            </div>

            <span class="count">342 friends</span>
        </div>

        <div class="people-grid">
            <article class="person-card">
                <div class="person-avatar">J</div>
                <h3>James Wilson</h3>
                <p>@jameswilson</p>
                <button>View profile</button>
            </article>

            <article class="person-card">
                <div class="person-avatar">S</div>
                <h3>Sarah Miller</h3>
                <p>@sarahmiller</p>
                <button>View profile</button>
            </article>

            <article class="person-card">
                <div class="person-avatar">M</div>
                <h3>Michael Brown</h3>
                <p>@michaelbrown</p>
                <button>View profile</button>
            </article>

            <article class="person-card">
                <div class="person-avatar">E</div>
                <h3>Emma Davis</h3>
                <p>@emmadavis</p>
                <button>View profile</button>
            </article>
        </div>
    </section>
</template>

<style scoped>
.profile-section {
    scroll-margin-top: 100px;
}

.section-heading {
    display: flex;
    justify-content: space-between;
    align-items: end;
    margin-bottom: var(--space-5);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-lime);
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

.count {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 9px;
}

.people-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--space-4);
}

.person-card {
    padding: var(--space-5) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    text-align: center;
    box-shadow: var(--shadow-raised);
    transition: transform 0.1s, border-color 0.1s;
}

.person-card:hover {
    transform: translateY(-2px);
    border-color: var(--color-cyan);
}

.person-card:nth-child(4n+1) .person-avatar {
    background: var(--gradient-cyber);
}

.person-card:nth-child(4n+2) .person-avatar {
    background: var(--gradient-sunset);
}

.person-card:nth-child(4n+3) .person-avatar {
    background: var(--gradient-aurora);
}

.person-card:nth-child(4n+4) .person-avatar {
    background: var(--gradient-action);
}

.person-avatar {
    width: 75px;
    height: 75px;
    margin: 0 auto var(--space-3);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--color-border);
    border-radius: 50%;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 30px;
}

.person-card h3 {
    margin: 0;
    color: var(--color-text);
    font-size: 13px;
}

.person-card p {
    margin: var(--space-1) 0 var(--space-4);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 8px;
}

.person-card button {
    width: 100%;
    padding: var(--space-2);
    border: 1px solid var(--color-cyan);
    border-radius: var(--radius-small);
    background: transparent;
    color: var(--color-cyan);
    font-family: var(--font-meta);
    font-size: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
}

.person-card button:hover {
    background: var(--color-cyan);
    color: var(--color-background);
}

.person-card button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

@media (max-width: 850px) {
    .people-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
</style>