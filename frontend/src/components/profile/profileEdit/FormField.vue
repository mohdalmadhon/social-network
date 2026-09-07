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
.form-field {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
}

label {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
}

input,
textarea {
    width: 100%;
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-medium);
    background: var(--color-input);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.8125rem;
    box-sizing: border-box;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

input::placeholder,
textarea::placeholder {
    color: var(--color-text-faint);
}

textarea {
    min-height: 6.875rem;
    resize: vertical;
}

input:focus,
textarea:focus {
    outline: none;
    border-color: var(--color-violet);
    box-shadow: var(--focus-ring);
}
</style>
