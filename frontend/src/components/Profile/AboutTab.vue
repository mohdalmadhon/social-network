<script setup>
import { ref } from 'vue';

const activeTab = ref('information');

const props = defineProps({
    about: {
        type: Object,
        required: true
    }
});

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
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 9px);
    letter-spacing: 2px;
    text-transform: uppercase;
}

h2 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: clamp(20px, 4vw, 29px);
}

.about-card {
    display: grid;
    grid-template-columns: minmax(150px, 200px) 1fr;
    min-height: 360px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    overflow: hidden;
}

.about-navigation {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: clamp(12px, 2vw, 18px);
    border-right: 1px solid var(--color-border);
    background: var(--color-surface-violet);
}

.about-navigation button {
    width: 100%;
    padding: var(--space-3) var(--space-2);
    border: 1px solid transparent;
    border-radius: var(--radius-small);
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: clamp(9px, 1.4vw, 11px);
    font-weight: 600;
    text-align: left;
    cursor: pointer;
    transition: 0.15s ease;
}

.about-navigation button:hover {
    border-color: var(--color-border);
    background: var(--color-surface-raised);
    color: var(--color-text);
}

.about-navigation button:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.about-navigation button.active {
    border-color: transparent;
    background: var(--gradient-cyber);
    color: var(--color-text);
    box-shadow: var(--shadow-raised);
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
    gap: var(--space-3);
}

.info-item {
    display: grid;
    grid-template-columns: minmax(90px, 130px) 1fr;
    gap: var(--space-5);
    padding-bottom: var(--space-3);
    border-bottom: 1px solid var(--color-border);
}

.info-item span,
.link-item span {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 10px);
    font-weight: 600;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.info-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-display);
    font-size: clamp(12px, 1.8vw, 14px);
    line-height: 1.5;
    word-break: break-word;
}

.links-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

.link-item {
    display: grid;
    grid-template-columns: minmax(90px, 130px) 1fr;
    gap: var(--space-5);
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-left: 3px solid var(--color-amber);
    border-radius: var(--radius-small);
    background: var(--color-surface);
    color: var(--color-text);
    text-decoration: none;
    transition: 0.15s ease;
}

.link-item:nth-of-type(2) {
    border-left-color: var(--color-blue);
}

.link-item:nth-of-type(3) {
    border-left-color: var(--color-cyan);
}

.link-item:nth-of-type(4) {
    border-left-color: var(--color-magenta);
}

.link-item:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-raised);
    border-color: var(--color-violet);
}

.link-item:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
}

.link-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-display);
    font-size: clamp(11px, 1.7vw, 13px);
    word-break: break-all;
}

.empty {
    margin: 0;
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 10px;
}

@media (max-width: 600px) {
    .about-card {
        grid-template-columns: 1fr;
    }

    .about-navigation {
        flex-direction: row;
        border-right: none;
        border-bottom: 1px solid var(--color-border);
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
        gap: var(--space-1);
    }
}
</style>