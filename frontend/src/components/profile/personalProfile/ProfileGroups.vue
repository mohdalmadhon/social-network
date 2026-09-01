<template>
    <section class="profile-section">
        <div class="section-heading">
            <div>
                <p class="eyebrow">COMMUNITIES</p>
                <h2>Groups</h2>
            </div>
        </div>

        <div class="groups-grid">
            <article class="group-card">
                <div class="group-icon">G</div>
                <div>
                    <h3>Developers Hub</h3>
                    <p>2.4k members</p>
                </div>
                <button>View group</button>
            </article>

            <article class="group-card">
                <div class="group-icon">D</div>
                <div>
                    <h3>Design & Creative</h3>
                    <p>1.8k members</p>
                </div>
                <button>View group</button>
            </article>

            <article class="group-card">
                <div class="group-icon">P</div>
                <div>
                    <h3>Photography</h3>
                    <p>936 members</p>
                </div>
                <button>View group</button>
            </article>
        </div>
    </section>
</template>

<style scoped>
.profile-section {
    scroll-margin-top: 100px;
}

.section-heading {
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

.groups-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
}

.group-card {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 15px;
    padding: 18px;
    border: 1px solid #232332;
    border-radius: 12px;
    background: #171724;
}

.group-icon {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 10px;
    background: #a855f7;
    color: white;
    font-weight: 700;
    font-size: 20px;
}

.group-card h3 {
    margin: 3px 0 5px;
    color: #fff;
    font-size: 13px;
    font-weight: 600;
}

.group-card p {
    margin: 0;
    color: #6b6b7d;
    font-size: 11px;
}

.group-card button {
    grid-column: 1 / -1;
    padding: 8px 0;
    border: 1px solid #2a2a3a;
    border-radius: 8px;
    background: #1c1c2a;
    color: #8b8b9e;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
}

@media (max-width: 800px) {
    .groups-grid {
        grid-template-columns: 1fr;
    }
}
</style>