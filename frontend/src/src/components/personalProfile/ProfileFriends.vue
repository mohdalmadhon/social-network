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
    margin-bottom: 20px;
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: 29px;
}

.count {
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
}

.people-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 18px;
}

.person-card {
    padding: 20px 15px;
    border: 2px solid var(--main-color);
    border-radius: 7px;
    background: var(--bg-color);
    text-align: center;
    box-shadow: 4px 4px var(--main-color);
}

.person-avatar {
    width: 75px;
    height: 75px;
    margin: 0 auto 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--input-focus);
    color: white;
    font-family: "Liter", serif;
    font-size: 30px;
}

.person-card h3 {
    margin: 0;
    font-size: 13px;
}

.person-card p {
    margin: 5px 0 15px;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 8px;
}

.person-card button {
    width: 100%;
    padding: 8px;
    border: 2px solid var(--main-color);
    border-radius: 4px;
    background: var(--page-background);
    color: var(--main-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 8px;
    font-weight: 600;
}

@media (max-width: 850px) {
    .people-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
</style>