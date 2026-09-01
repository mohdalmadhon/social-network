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

.about-card {
    padding: 4px 20px;
    border: 1px solid #232332;
    border-radius: 16px;
    background: #12121c;
}

.about-row {
    display: grid;
    grid-template-columns: 180px 1fr;
    gap: 20px;
    padding: 18px 0;
    border-bottom: 1px solid #232332;
}

.label {
    color: #8b8b9e;
    font-size: 11px;
    font-weight: 600;
}

.about-row strong {
    color: #c4c4d4;
    font-size: 13px;
    font-weight: 600;
}

.about-description {
    padding: 18px 0;
}

.about-description p {
    margin: 10px 0 0;
    color: #c4c4d4;
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