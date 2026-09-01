<script setup>
import { ref } from 'vue';

const activeTab = ref('information');

const props = defineProps({
    about: {
        type: Object,
        required: true
    }
});
console.log(props.about)
const tabs = [
    {
        id: 'information',
        label: 'Information'
    },
    {
        id: 'links',
        label: 'Social Links'
    }
];
</script>

<template>
    <section class="about-section">
        <div class="section-heading">
            <p class="eyebrow">PROFILE</p>
            <h2>About</h2>
        </div>

        <div class="about-card">
            <nav class="about-navigation">
                <button
                    v-for="tab in tabs"
                    :key="tab.id"
                    type="button"
                    :class="{ active: activeTab === tab.id }"
                    @click="activeTab = tab.id"
                >
                    {{ tab.label }}
                </button>
            </nav>

            <main class="about-content">
                <div
                    v-if="activeTab === 'information'"
                    class="content-section"
                >
                    <div class="content-heading">
                        <p class="eyebrow">PERSONAL</p>
                        <h3>Information</h3>
                    </div>

                    <div class="info-list">
                        <div
                            v-if="props.about.work"
                            class="info-item"
                        >
                            <span>Work</span>
                            <p>{{ props.about.work }}</p>
                        </div>

                        <div
                            v-if="props.about.education"
                            class="info-item"
                        >
                            <span>Education</span>
                            <p>{{ props.about.education }}</p>
                        </div>

                        <div
                            v-if="props.about.hobbies"
                            class="info-item"
                        >
                            <span>Hobbies</span>
                            <p>{{ props.about.hobbies }}</p>
                        </div>

                        <div
                            v-if="props.about.intrests"
                            class="info-item"
                        >
                            <span>Interests</span>
                            <p>{{ props.about.intrests }}</p>
                        </div>

                        <div
                            v-if="props.about.travel"
                            class="info-item"
                        >
                            <span>Travel</span>
                            <p>{{ props.about.travel }}</p>
                        </div>

                        <p
                            v-if="
                                !props.about.work &&
                                !props.about.education &&
                                !props.about.hobbies &&
                                !props.about.intrests &&
                                !props.about.travel
                            "
                            class="empty"
                        >
                            No information provided.
                        </p>
                    </div>
                </div>

                <div
                    v-if="activeTab === 'links'"
                    class="content-section"
                >
                    <div class="content-heading">
                        <p class="eyebrow">SOCIAL</p>
                        <h3>Links</h3>
                    </div>

                    <div class="links-list">
                        <a
                            v-if="props.about.website"
                            :href="props.about.website"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="link-item"
                        >
                            <span>Website</span>
                            <p>{{ props.about.website }}</p>
                        </a>

                        <a
                            v-if="props.about.linkedin"
                            :href="props.about.linkedin"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="link-item"
                        >
                            <span>LinkedIn</span>
                            <p>{{ props.about.linkedin }}</p>
                        </a>

                        <a
                            v-if="props.about.twitter"
                            :href="props.about.twitter"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="link-item"
                        >
                            <span>Twitter / X</span>
                            <p>{{ props.about.twitter }}</p>
                        </a>

                        <a
                            v-if="props.about.instgram"
                            :href="props.about.instgram"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="link-item"
                        >
                            <span>Instagram</span>
                            <p>{{ props.about.instgram }}</p>
                        </a>

                        <p
                            v-if="
                                !props.about.website &&
                                !props.about.linkedin &&
                                !props.about.twitter &&
                                !props.about.instgram
                            "
                            class="empty"
                        >
                            No social links provided.
                        </p>
                    </div>
                </div>
            </main>
        </div>
    </section>
</template>

<style scoped>
.about-section {
    width: 100%;
    scroll-margin-top: 100px;
}

.section-heading {
    margin-bottom: clamp(14px, 2.5vw, 20px);
}

.eyebrow {
    margin: 0 0 5px;
    color: #a855f7;
    font-size: clamp(10px, 1.2vw, 11px);
    font-weight: 600;
    letter-spacing: 2px;
}

h2 {
    margin: 0;
    color: #fff;
    font-size: clamp(18px, 4vw, 20px);
    font-weight: 700;
}

.about-card {
    display: grid;
    grid-template-columns: minmax(150px, 200px) 1fr;
    min-height: 360px;
    border: 1px solid #232332;
    border-radius: 16px;
    background: #12121c;
    overflow: hidden;
}

.about-navigation {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: clamp(12px, 2vw, 18px);
    border-right: 1px solid #232332;
    background: #171724;
}

.about-navigation button {
    width: 100%;
    padding: 11px 10px;
    border: 0;
    border-radius: 8px;
    background: transparent;
    color: #8b8b9e;
    font-size: clamp(12px, 1.4vw, 13px);
    font-weight: 500;
    text-align: left;
    cursor: pointer;
    transition: 0.15s ease;
}

.about-navigation button:hover {
    color: #fff;
}

.about-navigation button.active {
    background: #1c1c2a;
    color: #a855f7;
    font-weight: 600;
}

.about-content {
    min-width: 0;
    padding: clamp(18px, 3vw, 28px);
}

.content-section {
    width: 100%;
}

.content-heading {
    margin-bottom: clamp(18px, 3vw, 25px);
}

.info-list {
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.info-item {
    display: grid;
    grid-template-columns: minmax(90px, 130px) 1fr;
    gap: 20px;
    padding-bottom: 14px;
    border-bottom: 1px solid #232332;
}

.info-item span,
.link-item span {
    color: #8b8b9e;
    font-size: clamp(10px, 1.2vw, 11px);
    font-weight: 600;
    letter-spacing: 0.6px;
    text-transform: uppercase;
}

.info-item p {
    margin: 0;
    color: #c4c4d4;
    font-size: clamp(12px, 1.8vw, 13px);
    line-height: 1.5;
    word-break: break-word;
}

.links-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.link-item {
    display: grid;
    grid-template-columns: minmax(90px, 130px) 1fr;
    gap: 20px;
    padding: 14px;
    border: 1px solid #232332;
    border-radius: 12px;
    background: #171724;
    color: #c4c4d4;
    text-decoration: none;
    transition: 0.15s ease;
}

.link-item:hover {
    border-color: #a855f7;
}

.link-item p {
    margin: 0;
    color: #c4c4d4;
    font-size: clamp(11px, 1.7vw, 12px);
    word-break: break-all;
}

.empty {
    margin: 0;
    color: #6b6b7d;
    font-size: 12px;
}

@media (max-width: 600px) {
    .about-card {
        grid-template-columns: 1fr;
    }

    .about-navigation {
        flex-direction: row;
        border-right: none;
        border-bottom: 1px solid #232332;
        overflow-x: auto;
    }

    .about-navigation button {
        width: auto;
        min-width: max-content;
        text-align: center;
    }
}

@media (max-width: 500px) {
    .info-item,
    .link-item {
        grid-template-columns: 1fr;
        gap: 6px;
    }
}
</style>