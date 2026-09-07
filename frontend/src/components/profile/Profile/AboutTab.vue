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
    scroll-margin-top: 6.25rem;
}

.section-heading {
    margin-bottom: var(--space-4);
}

.eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.15em;
}

h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.8125rem;
    color: var(--color-text);
}

.about-card {
    display: grid;
    grid-template-columns: minmax(9.375rem, 12.5rem) 1fr;
    min-height: 22.5rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-large);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    overflow: hidden;
}

.about-navigation {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: var(--space-4);
    border-right: 1px solid var(--color-border);
    background: var(--color-sidebar);
}

.about-navigation button {
    width: 100%;
    padding: var(--space-3) var(--space-3);
    border: none;
    border-radius: var(--radius-medium);
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-body);
    font-size: 0.75rem;
    font-weight: 600;
    text-align: left;
    cursor: pointer;
    transition: background 0.15s ease, color 0.15s ease;
}

.about-navigation button:hover {
    background: var(--color-surface-raised);
    color: var(--color-text-soft);
}

.about-navigation button.active {
    background: var(--gradient-action);
    color: var(--color-text);
}

.about-content {
    min-width: 0;
    padding: var(--space-6);
}

.content-section {
    width: 100%;
}

.content-heading {
    margin-bottom: var(--space-5);
}

.info-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
}

.info-item {
    display: grid;
    grid-template-columns: minmax(5.625rem, 8.125rem) 1fr;
    gap: var(--space-5);
    padding-bottom: var(--space-4);
    border-bottom: 1px solid var(--color-border);
}

.info-item span,
.link-item span {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
}

.info-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 0.875rem;
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
    grid-template-columns: minmax(5.625rem, 8.125rem) 1fr;
    gap: var(--space-5);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface-raised);
    color: var(--color-text-soft);
    text-decoration: none;
    transition: border-color 0.15s ease, transform 0.15s ease;
}

.link-item:hover {
    transform: translateY(-1px);
    border-color: var(--color-violet);
}

.link-item p {
    margin: 0;
    color: var(--color-text-soft);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    word-break: break-all;
}

.empty {
    margin: 0;
    color: var(--color-text-faint);
    font-family: var(--font-meta);
    font-size: 0.625rem;
}

@media (max-width: 37.5rem) {
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

@media (max-width: 31.25rem) {
    .info-item,
    .link-item {
        grid-template-columns: 1fr;
        gap: var(--space-2);
    }
}
</style>
