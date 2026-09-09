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
    scroll-margin-top: 100px;
}

.section-heading {
    margin-bottom: 20px;
}

.eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: 29px;
}

.list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.list-item {
    display: flex;
    align-items: center;
    gap: 15px;
    padding: 15px 20px;
    border: 2px solid var(--main-color);
    border-radius: 6px;
    background: var(--bg-color);
    box-shadow: 4px 4px var(--main-color);
}

.list-avatar {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--input-focus);
    color: white;
    font-family: "Liter", serif;
    font-size: 20px;
}

.list-item div:nth-child(2) {
    flex: 1;
}

.list-item strong,
.list-item span {
    display: block;
}

.list-item strong {
    font-size: 12px;
}

.list-item span {
    margin-top: 3px;
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 8px;
}

.list-item button {
    padding: 8px 12px;
    border: 2px solid var(--main-color);
    border-radius: 4px;
    background: var(--page-background);
    font-family: "JetBrains Mono", monospace;
    font-size: 8px;
    font-weight: 600;
}
</style>