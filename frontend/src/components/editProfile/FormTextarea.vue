<script setup>
defineProps({
  id: { type: String, required: true },
  label: { type: String, required: true },
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  rows: { type: [String, Number], default: 5 },
  maxlength: { type: [String, Number], default: undefined },
  hint: { type: String, default: '' },
  error: { type: String, default: '' }
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

    <textarea :id="id" class="field__textarea" :rows="rows" :maxlength="maxlength" :placeholder="placeholder"
      :value="modelValue" @input="onInput"></textarea>

    <div class="field__bottom">
      <span v-if="hint" class="field__hint">
        {{ hint }}
      </span>

      <span v-if="maxlength" class="field__counter">
        {{ modelValue.length }} / {{ maxlength }}
      </span>
    </div>

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

.field__textarea {
  width: 100%;
  min-height: 8rem;
  padding: var(--space-3);
  background: var(--color-input);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-small);
  color: var(--color-text);
  font: inherit;
  line-height: 1.55;
  resize: vertical;
  transition:
    border-color 0.18s ease,
    background 0.18s ease,
    box-shadow 0.18s ease;
}

.field__textarea::placeholder {
  color: var(--color-text-faint);
}

.field__textarea:hover {
  border-color: var(--color-text-faint);
  background: var(--color-surface-raised);
}

.field__textarea:focus {
  outline: none;
  border-color: var(--color-violet);
  background: var(--color-surface-raised);
  box-shadow: var(--focus-ring);
}

.field__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
}

.field__hint,
.field__counter {
  color: var(--color-text-faint);
  font-family: var(--font-meta);
  font-size: 0.68rem;
}

.field__counter {
  color: var(--color-text-muted);
}

.field__error {
  color: var(--color-danger, #ef4444);
  font-size: 0.75rem;
  line-height: 1.4;
}
</style>
