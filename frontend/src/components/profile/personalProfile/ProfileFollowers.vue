<template>
    <section class="profile-section">
        <div class="section-heading">
            <p class="eyebrow">YOUR COMMUNITY</p>
            <h2>Followers</h2>
        </div>

        <div class="list">
            <article class="list-item">
                <div class="list-avatar">L</div>
                <div>
                    <strong>Lucas Martin</strong>
                    <span>@lucasmartin</span>
                </div>
                <button>Follow back</button>
            </article>

            <article class="list-item">
                <div class="list-avatar">N</div>
                <div>
                    <strong>Nora Smith</strong>
                    <span>@norasmith</span>
                </div>
                <button>Follow back</button>
            </article>

            <article class="list-item">
                <div class="list-avatar">R</div>
                <div>
                    <strong>Ryan Cooper</strong>
                    <span>@ryancooper</span>
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

.list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.list-item {
    display: flex;
    align-items: center;
    gap: 15px;
    padding: 15px 16px;
    border: 1px solid #232332;
    border-radius: 12px;
    background: #171724;
}

.list-avatar {
    width: 44px;
    height: 44px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: #a855f7;
    color: white;
    font-weight: 700;
    font-size: 16px;
}

.list-item div:nth-child(2) {
    flex: 1;
}

.list-item strong,
.list-item span {
    display: block;
}

.list-item strong {
    color: #fff;
    font-size: 13px;
    font-weight: 600;
}

.list-item span {
    margin-top: 3px;
    color: #6b6b7d;
    font-size: 11px;
}

.list-item button {
    flex-shrink: 0;
    padding: 7px 12px;
    border: 1px solid #2a2a3a;
    border-radius: 8px;
    background: #1c1c2a;
    color: #8b8b9e;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
}
</style>