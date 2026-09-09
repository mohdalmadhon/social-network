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
</style>