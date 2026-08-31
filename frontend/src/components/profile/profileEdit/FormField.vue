<script setup>
defineProps({
    label: String,
    id: String,
    modelValue: String,
    type: {
        type: String,
        default: 'text'
    },
    placeholder: String
})

defineEmits(['update:modelValue'])
</script>

<template>
    <div class="form-field">
        <label :for="id">{{ label }}</label>

        <textarea
            v-if="type === 'textarea'"
            :id="id"
            :placeholder="placeholder"
            :value="modelValue"
            @input="$emit('update:modelValue', $event.target.value)"
        ></textarea>

        <input
            v-else
            :id="id"
            :type="type"
            :placeholder="placeholder"
            :value="modelValue"
            @input="$emit('update:modelValue', $event.target.value)"
        >
    </div>
</template>

<style scoped>
@import '../../../styles/global.css';

.form-field {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
}

label {
    color: var(--color-text-soft);
    font-size: 0.82rem;
    font-weight: 600;
}

input,
textarea {
    width: 100%;
    padding: 0 var(--space-3);
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    color: var(--color-text);
    font: inherit;
    box-sizing: border-box;
    transition:
        border-color 0.18s ease,
        background 0.18s ease,
        box-shadow 0.18s ease;
}

input {
    min-height: var(--touch-target);
}

textarea {
    min-height: 8rem;
    padding-block: var(--space-3);
    line-height: 1.55;
    resize: vertical;
}

input::placeholder,
textarea::placeholder {
    color: var(--color-text-faint);
}

input:hover,
textarea:hover {
    border-color: var(--color-text-faint);
    background: var(--color-surface-raised);
}

input:focus,
textarea:focus {
    outline: none;
    border-color: var(--color-violet);
    background: var(--color-surface-raised);
    box-shadow: var(--focus-ring);
}
</style>
