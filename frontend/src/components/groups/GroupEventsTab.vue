<script setup>
defineProps({
    events: {
        type: Array,
        default: () => []
    }
});
</script>

<template>
    <section class="tab-panel">
        <div v-if="!events.length" class="empty-state">
            No events scheduled yet.
        </div>

        <div v-else class="events-list">
            <article
                v-for="event in events"
                :key="event.id"
                class="event-card"
            >
                <div class="event-date">
                    {{ event.date }}
                </div>

                <div class="event-info">
                    <h3 class="event-title">
                        {{ event.title }}
                    </h3>

                    <p v-if="event.description" class="event-description">
                        {{ event.description }}
                    </p>
                </div>
            </article>
        </div>
    </section>
</template>

<style scoped>
.tab-panel {
    width: 100%;
}

.empty-state {
    padding: 40px 20px;
    text-align: center;

    border: 2px solid var(--main-color);
    border-radius: 8px;

    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);

    color: var(--font-color-sub);

    font-family: "JetBrains Mono", monospace;
    font-size: 12px;
}

.events-list {
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.event-card {
    display: flex;
    align-items: center;
    gap: 16px;

    padding: 16px 18px;

    border: 2px solid var(--main-color);
    border-radius: 8px;

    background: var(--bg-color);
    box-shadow: 5px 5px var(--main-color);
}

.event-date {
    flex-shrink: 0;

    padding: 10px 14px;

    border: 2px solid var(--main-color);
    border-radius: 6px;

    background: var(--input-focus);
    color: #fff;

    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
    font-weight: 700;
    text-align: center;
    white-space: nowrap;
}

.event-info {
    flex: 1;
    min-width: 0;
}

.event-title {
    margin: 0 0 4px;

    color: var(--font-color);

    font-family: "Liter", serif;
    font-size: 16px;
    font-weight: 700;
}

.event-description {
    margin: 0;

    color: var(--font-color-sub);

    font-size: 13px;
    line-height: 1.4;
}
</style>
