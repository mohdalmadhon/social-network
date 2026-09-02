<script setup>
defineProps({
    modelValue: {
        type: String,
        default: ''
    }
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
    <div class="group-search">
        <input class="search-input" type="text" placeholder="Search groups..." :value="modelValue"
            @input="emit('update:modelValue', $event.target.value)" />
    </div>
</template>

<style scoped>
.group-search {
    position: relative;
    width: 100%;
}

.search-input {
    width: 100%;
    min-height: var(--touch-target);
    padding: 0 var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 0.875rem;
    transition:
        border-color 0.15s ease,
        background 0.15s ease;
}

.search-input::placeholder {
    color: var(--color-text-faint);
}

.search-input:hover {
    border-color: #434c70;
}

.search-input:focus {
    border-color: var(--color-violet);
    outline: none;
    box-shadow: var(--focus-ring);
}
</style>