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

.count {
    color: #8b8b9e;
    font-size: 12px;
}

.people-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
}

.person-card {
    padding: 20px 15px;
    border: 1px solid #232332;
    border-radius: 12px;
    background: #171724;
    text-align: center;
}

.person-avatar {
    width: 60px;
    height: 60px;
    margin: 0 auto 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: #a855f7;
    color: white;
    font-weight: 700;
    font-size: 24px;
}

.person-card h3 {
    margin: 0;
    color: #fff;
    font-size: 13px;
    font-weight: 600;
}

.person-card p {
    margin: 5px 0 15px;
    color: #6b6b7d;
    font-size: 11px;
}

.person-card button {
    width: 100%;
    padding: 7px 0;
    border: 1px solid #2a2a3a;
    border-radius: 8px;
    background: #1c1c2a;
    color: #8b8b9e;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
}

@media (max-width: 850px) {
    .people-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
</style>