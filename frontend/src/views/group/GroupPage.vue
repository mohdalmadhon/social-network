<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()

const groupId = route.params.groupId
</script>

<template>
    <main class="group-page">
        <section class="group-container">
            <header class="group-header">
                <h1>Group</h1>
                <p>Group ID: {{ groupId }}</p>
            </header>

            <section class="group-content">
                <p>Group information will be displayed here.</p>
            </section>
        </section>
    </main>
</template>

<style scoped>
.group-page {
    width: 100%;
    min-height: 100vh;
    padding: var(--space-6);
}

.group-container {
    width: 100%;
    max-width: 900px;
    margin: 0 auto;
}

.group-header {
    margin-bottom: var(--space-6);
}

.group-header h1 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-heading);
}

.group-header p {
    margin-top: var(--space-2);
    color: var(--color-text-muted);
}

.group-content {
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    color: var(--color-text);
}
</style>