<template>
    <section class="profile-section">
        <div class="section-heading">
            <p class="eyebrow">PROFILE</p>
            <h2>About</h2>
        </div>

        <div class="about-card">
            <div class="about-row">
                <span class="label">Username</span>
                <strong>@alexanderson</strong>
            </div>

            <div class="about-row">
                <span class="label">Email</span>
                <strong>alex@example.com</strong>
            </div>

            <div class="about-row">
                <span class="label">Date of birth</span>
                <strong>March 18, 2002</strong>
            </div>

            <div class="about-row">
                <span class="label">Joined</span>
                <strong>August 2026</strong>
            </div>

            <div class="about-description">
                <span class="label">About me</span>

                <p>
                    Building things, learning new stuff and enjoying
                    the little moments. Always looking for something
                    interesting to create.
                </p>
            </div>
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

.about-card {
    padding: 10px 25px;
    border: 2px solid var(--main-color);
    border-radius: 7px;
    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);
}

.about-row {
    display: grid;
    grid-template-columns: 180px 1fr;
    gap: 20px;
    padding: 20px 0;
    border-bottom: 2px solid var(--page-background);
}

.label {
    color: var(--font-color-sub);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    font-weight: 600;
}

.about-row strong {
    font-size: 12px;
}

.about-description {
    padding: 20px 0;
}

.about-description p {
    margin: 12px 0 0;
    color: var(--font-color);
    font-size: 13px;
    line-height: 1.7;
}

@media (max-width: 550px) {
    .about-row {
        grid-template-columns: 1fr;
        gap: 7px;
    }
}
</style>