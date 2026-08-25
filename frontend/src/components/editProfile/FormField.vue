<script setup>
defineProps({
  id: { type: String, required: true },
  label: { type: String, required: true },
  modelValue: { type: String, default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  autocomplete: { type: String, default: 'off' },
  maxlength: { type: [String, Number], default: undefined },
  error: { type: String, default: '' },
  prefix: { type: String, default: '' }
})

const emit = defineEmits(['update:modelValue', 'input'])

function onInput(event) {
  emit('update:modelValue', event.target.value)
  emit('input', event)
}
</script>

<template>
  <div class="field">
    <label class="field__label" :for="id">
      {{ label }}
    </label>

    <div v-if="prefix" class="field__prefixed" :class="{ 'field__prefixed--error': error }">
      <span class="field__prefix">{{ prefix }}</span>

      <input :id="id" class="field__input field__input--prefixed" :type="type" :value="modelValue"
        :placeholder="placeholder" :autocomplete="autocomplete" :maxlength="maxlength" @input="onInput" />
    </div>

    <input v-else :id="id" class="field__input" :class="{ 'field__input--error': error }" :type="type"
      :value="modelValue" :placeholder="placeholder" :autocomplete="autocomplete" :maxlength="maxlength"
      @input="onInput" />

    <span v-if="error" class="field__error">
      {{ error }}
    </span>
  </div>
</template>

<style scoped>
@import '../../styles/variables.css';

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.field__label {
  color: var(--color-text-soft);
  font-size: 0.82rem;
  font-weight: 600;
}

.field__input {
  width: 100%;
  min-height: var(--touch-target);
  padding: 0 var(--space-3);
  background: var(--color-input);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  color: var(--color-text);
  font: inherit;
  transition:
    border-color 0.18s ease,
    background 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease;
}

.field__input::placeholder {
  color: var(--color-text-faint);
}

.field__input:hover {
  border-color: var(--color-text-faint);
  background: var(--color-surface-raised);
}

.field__input:focus {
  outline: none;
  border-color: var(--color-violet);
  background: var(--color-surface-raised);
  box-shadow: var(--focus-ring);
}

.field__input:focus-visible {
  outline: none;
}

.field__prefixed {
  display: flex;
  align-items: center;
  min-height: var(--touch-target);
  background: var(--color-input);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  transition:
    border-color 0.18s ease,
    background 0.18s ease,
    box-shadow 0.18s ease;
}

.field__prefixed:hover {
  border-color: var(--color-text-faint);
  background: var(--color-surface-raised);
}

.field__prefixed:focus-within {
  border-color: var(--color-violet);
  background: var(--color-surface-raised);
  box-shadow: var(--focus-ring);
}

.field__prefix {
  padding-left: var(--space-3);
  color: var(--color-violet);
  font-family: var(--font-meta);
  font-size: 0.9rem;
  font-weight: 600;
}

.field__input--prefixed {
  min-height: auto;
  padding-left: var(--space-2);
  border: 0;
  background: transparent;
  box-shadow: none;
}

.field__input--prefixed:hover,
.field__input--prefixed:focus {
  border: 0;
  background: transparent;
  box-shadow: none;
}

.field__error {
  color: var(--color-danger, #ef4444);
  font-size: 0.75rem;
  line-height: 1.4;
}

.field__input--error {
  border-color: var(--color-danger, #ef4444);
}

.field__input--error:hover,
.field__input--error:focus {
  border-color: var(--color-danger, #ef4444);
}

.field__input--error:focus {
  box-shadow: 0 0 0 3px rgb(239 68 68 / 12%);
}

.field__prefixed--error {
  border-color: var(--color-danger, #ef4444);
}

.field__prefixed--error:hover,
.field__prefixed--error:focus-within {
  border-color: var(--color-danger, #ef4444);
}

.field__prefixed--error:focus-within {
  box-shadow: 0 0 0 3px rgb(239 68 68 / 12%);
}
</style>
