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
@import '../../styles/global.css';
@import '../../styles/variables.css';
.form-field {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
}

label {
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: clamp(8px, 1.2vw, 9px);
    font-weight: 600;
    letter-spacing: 1px;
}

input,
textarea {
    width: 100%;
    padding: clamp(10px, 1.8vw, 12px) clamp(10px, 2vw, 14px);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    color: var(--color-text);
    font-family: var(--font-display);
    font-size: clamp(12px, 1.8vw, 13px);
    box-sizing: border-box;
    transition: border-color 0.15s, box-shadow 0.15s;
}

input::placeholder,
textarea::placeholder {
    color: var(--color-text-faint);
}

textarea {
    min-height: 110px;
    resize: vertical;
}

input:focus,
textarea:focus {
    outline: none;
    border-color: var(--color-cyan);
    box-shadow: 0 0 0 3px rgb(34 229 229 / 24%);
}

textarea:focus {
    border-color: var(--color-mint);
    box-shadow: 0 0 0 3px rgb(46 235 181 / 24%);
}

input[type="password"]:focus {
    border-color: var(--color-magenta);
    box-shadow: 0 0 0 3px rgb(255 63 216 / 24%);
}

input[type="email"]:focus {
    border-color: var(--color-blue);
    box-shadow: 0 0 0 3px rgb(56 182 255 / 24%);
}

.form-field:focus-within label {
    color: var(--color-text-soft);
}
</style>