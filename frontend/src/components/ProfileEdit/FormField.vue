<script setup>
defineProps({
  label: String,
  id: String,
  modelValue: [String, Number],
  type: { type: String, default: 'text' },
  placeholder: String,
})
defineEmits(['update:modelValue'])
</script>

<template>
  <div class="form-field">
    <label :for="id">{{ label }}</label>
    <textarea v-if="type === 'textarea'" :id="id" :placeholder="placeholder" :value="modelValue" maxlength="1000" @input="$emit('update:modelValue', $event.target.value)" />
    <input v-else :id="id" :type="type" :placeholder="placeholder" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)" />
  </div>
</template>

<style scoped>
.form-field { display: grid; min-width: 0; gap: var(--space-2); }
label { color: var(--color-text-muted); font-size: .8125rem; font-weight: 600; }
input, textarea { width: 100%; min-width: 0; padding: var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-small); outline: none; background: var(--color-input); color: var(--color-text); }
input::placeholder, textarea::placeholder { color: var(--color-text-faint); }
textarea { min-height: 8rem; line-height: 1.55; }
input:focus, textarea:focus { border-color: var(--color-violet); box-shadow: var(--focus-ring); }
.form-field:focus-within label { color: var(--color-text-soft); }
</style>
