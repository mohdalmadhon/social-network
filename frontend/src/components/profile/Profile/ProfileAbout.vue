<script setup>
import { ref } from 'vue';
import AboutInfo from './AboutInfo.vue';
import AboutLinks from './AboutLinks.vue';
import { userAbout } from '@/data/usersData';

const activeTab = ref('information');

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
                <AboutInfo v-if="activeTab === 'information'" :about="userAbout" />
                <AboutLinks v-if="activeTab === 'links'" :about="userAbout" />
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
</style>
