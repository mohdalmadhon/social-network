<script setup>
defineProps({
    group: {
        type: Object,
        required: true
    }
})

const emit = defineEmits(['join-group', 'view-group'])
</script>

<template>
    <article class="group-card">
        <p class="group-meta">Community</p>

        <h3 class="group-title">
            {{ group.title }}
        </h3>

        <p class="group-description">
            {{ group.description }}
        </p>

        <div class="group-footer">
            <span class="member-count">
                {{ group.memberCount }} members
            </span>

            <div class="actions">
                <button class="view-button" @click="emit('view-group', group.id)">
                    View
                </button>

                <button v-if="!group.isMember && !group.requested" class="join-button"
                    @click="emit('join-group', group.id)">
                    Join
                </button>

                <button v-else-if="group.requested" class="requested-button" disabled>
                    Requested
                </button>
            </div>
        </div>
    </article>
</template>

<style scoped>
.group-card {
    display: flex;
    flex-direction: column;
    min-height: 220px;
    padding: var(--space-5);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-surface);
    box-shadow: var(--shadow-raised);
    transition:
        transform 0.15s ease,
        border-color 0.15s ease,
        background 0.15s ease;
}

.group-card:hover {
    transform: translateY(-3px);
    border-color: var(--color-violet);
    background: var(--color-surface-raised);
}

.group-meta {
    margin: 0 0 var(--space-2);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.625rem;
    font-weight: 600;
    letter-spacing: 1.5px;
    text-transform: uppercase;
}

.group-title {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: 1.125rem;
    font-weight: 600;
    line-height: 1.3;
}

.group-description {
    margin: var(--space-2) 0 var(--space-4);
    color: var(--color-text-muted);
    font-size: 0.8125rem;
    line-height: 1.6;
}

.group-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    margin-top: auto;
    padding-top: var(--space-4);
    border-top: 1px solid var(--color-border);
}

.member-count {
    color: var(--color-text-faint);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
}

.actions {
    display: flex;
    gap: var(--space-2);
}

.view-button,
.join-button,
.requested-button {
    min-height: 36px;
    padding: 0 var(--space-3);
    border-radius: var(--radius-small);
    font-size: 0.75rem;
    font-weight: 600;
}

.view-button {
    border: 1px solid var(--color-border);
    background: var(--color-surface-raised);
    color: var(--color-text-soft);
    cursor: pointer;
}

.view-button:hover {
    border-color: var(--color-violet);
    color: var(--color-text);
}

.join-button {
    border: 0;
    background: var(--gradient-action);
    color: #fff;
    cursor: pointer;
}

.join-button:hover {
    opacity: 0.9;
}

.requested-button {
    border: 1px solid var(--color-border);
    background: var(--color-input);
    color: var(--color-text-faint);
    cursor: default;
}

@media (max-width: 28rem) {
    .group-footer {
        align-items: flex-start;
        flex-direction: column;
    }

    .actions {
        width: 100%;
    }

    .view-button,
    .join-button,
    .requested-button {
        flex: 1;
    }
}
</style>