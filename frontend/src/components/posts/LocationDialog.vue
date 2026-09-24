```vue
<script setup>
import {
  getLocationData,
  getSuggestedLocations
} from '@/api/common/location'
import { onBeforeUnmount, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits([
  'update:modelValue',
  'close'
])

const search = ref('')
const locations = ref([])
const selectedLocation = ref(null)
const isLoading = ref(false)
const error = ref('')
let searchTimeout = null

function handleLocationSearch() {
  clearTimeout(searchTimeout)

  selectedLocation.value = null
  error.value = ''

  const query = search.value.trim()

  if (!query) {
    locations.value = []
    return
  }

  searchTimeout = setTimeout(async () => {
    isLoading.value = true

    try {
      const result = await getSuggestedLocations(query)
      locations.value = getLocationData(result)
    } catch (error) {
      locations.value = []
      error.value =
        error?.message || 'Could not search locations.'
    } finally {
      isLoading.value = false
    }
  }, 500)
}

function selectLocation(location) {
  selectedLocation.value = location

  search.value =
    location.city ||
    location.state ||
    location.country ||
    ''

  if (location.country && location.city) {
    search.value = `${location.city}, ${location.country}`
  }

  locations.value = []
  error.value = ''
}

function addLocation() {
  if (!selectedLocation.value) {
    return
  }

  const location = selectedLocation.value

  const locationValue = [
    location.country,
    location.city || location.state
  ]
    .filter(Boolean)
    .join(', ')

  emit(
    'update:modelValue',
    `${locationValue}:${location.lat}:${location.lon}`
  )

  closeDialog()
}

function removeLocation() {
  emit('update:modelValue', '')
  closeDialog()
}

function closeDialog() {
  clearTimeout(searchTimeout)

  search.value = ''
  locations.value = []
  selectedLocation.value = null
  error.value = ''
  isLoading.value = false

  emit('close')
}

function clearSearch() {
  search.value = ''
  locations.value = []
  selectedLocation.value = null
  error.value = ''
}

function handleKeydown(event) {
  if (event.key === 'Escape') {
    closeDialog()
  }
}

watch(
  () => props.modelValue,
  (value) => {
    if (value) {
      search.value = value.split(':')[0].trim()
    }
  },
  { immediate: true }
)

onBeforeUnmount(() => {
  clearTimeout(searchTimeout)
})
</script>

<template>
  <Teleport to="body">
    <div
      class="location-dialog"
      role="dialog"
      aria-modal="true"
      aria-labelledby="location-dialog-title"
      @keydown="handleKeydown"
    >
      <button
        class="location-dialog__backdrop"
        type="button"
        aria-label="Close location dialog"
        @click="closeDialog"
      />

      <div class="location-dialog__panel">
        <div class="location-dialog__header">
          <div>
            <h2 id="location-dialog-title">
              Add location
            </h2>

            <p>
              Add a place to your post
            </p>
          </div>

          <button
            class="location-dialog__close"
            type="button"
            aria-label="Close location dialog"
            @click="closeDialog"
          >
            <svg
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path d="m7 7 10 10M17 7 7 17" />
            </svg>
          </button>
        </div>

        <div class="location-dialog__body">
          <div class="location-input">
            <svg
              class="location-input__icon"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path d="M12 21s7-6.1 7-12a7 7 0 1 0-14 0c0 5.9 7 12 7 12Z" />
              <circle
                cx="12"
                cy="9"
                r="2.25"
              />
            </svg>

            <input
              v-model="search"
              type="text"
              placeholder="Search city or country"
              autocomplete="off"
              autofocus
              @input="handleLocationSearch"
              @keydown.escape="closeDialog"
            />

            <button
              v-if="search"
              class="location-input__clear"
              type="button"
              aria-label="Clear search"
              @click="clearSearch"
            >
              <svg
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path d="m7 7 10 10M17 7 7 17" />
              </svg>
            </button>
          </div>

          <div
            v-if="isLoading"
            class="location-state"
          >
            Searching locations...
          </div>

          <div
            v-else-if="error"
            class="location-state location-state--error"
          >
            {{ error }}
          </div>

          <div
            v-else-if="locations.length"
            class="location-suggestions"
          >
            <div class="location-suggestions__header">
              Suggested locations
            </div>

            <button
              v-for="location in locations"
              :key="location.link"
              class="location-suggestion"
              type="button"
              @click="selectLocation(location)"
            >
              <span class="location-suggestion__icon">
                <svg
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path d="M12 21s7-6.1 7-12a7 7 0 1 0-14 0c0 5.9 7 12 7 12Z" />
                  <circle
                    cx="12"
                    cy="9"
                    r="2.25"
                  />
                </svg>
              </span>

              <span class="location-suggestion__content">
                <strong>
                  {{ location.city || location.state || location.country }}
                </strong>

                <small>
                  {{ location.country }}
                </small>
              </span>

              <svg
                class="location-suggestion__arrow"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path d="m9 18 6-6-6-6" />
              </svg>
            </button>
          </div>

          <div
            v-else-if="search.trim()"
            class="location-state"
          >
            No locations found.
          </div>

          <div
            v-if="selectedLocation"
            class="location-selected"
          >
            <div class="location-selected__icon">
              <svg
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path d="M12 21s7-6.1 7-12a7 7 0 1 0-14 0c0 5.9 7 12 7 12Z" />
                <circle
                  cx="12"
                  cy="9"
                  r="2.25"
                />
              </svg>
            </div>

            <div class="location-selected__content">
              <span>Selected location</span>

              <strong>
                {{
                  selectedLocation.city ||
                  selectedLocation.state ||
                  selectedLocation.country
                }}
              </strong>

              <small>
                {{ selectedLocation.country }}
              </small>
            </div>
          </div>
        </div>

        <div class="location-dialog__footer">
          <button
            class="location-dialog__cancel"
            type="button"
            @click="closeDialog"
          >
            Cancel
          </button>

          <button
            class="location-dialog__add"
            type="button"
            :disabled="!selectedLocation"
            @click="addLocation"
          >
            Add location
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.location-dialog {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.location-dialog__backdrop {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  padding: 0;
  border: 0;
  background: rgb(0 0 0 / 58%);
  cursor: default;
}

.location-dialog__panel {
  position: relative;
  z-index: 1;
  width: min(100%, 34rem);
  max-height: min(90vh, 42rem);
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: 1rem;
  background: var(--color-surface-raised);
  box-shadow: 0 1.5rem 4rem rgb(0 0 0 / 35%);
}

.location-dialog__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.25rem;
  border-bottom: 1px solid var(--color-border);
}

.location-dialog__header h2 {
  margin: 0;
  color: var(--color-text);
  font-size: 1.1rem;
}

.location-dialog__header p {
  margin: 0.3rem 0 0;
  color: var(--color-text-faint);
  font-size: 0.75rem;
}

.location-dialog__close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  flex-shrink: 0;
  padding: 0;
  border: 0;
  border-radius: 0.65rem;
  background: transparent;
  color: var(--color-text-faint);
  cursor: pointer;
}

.location-dialog__close:hover {
  background: var(--color-input);
  color: var(--color-text);
}

.location-dialog__close svg {
  width: 1rem;
  height: 1rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
}

.location-dialog__body {
  position: relative;
  padding: 1.25rem;
  overflow-y: auto;
}

.location-input {
  position: relative;
}

.location-input__icon {
  position: absolute;
  top: 50%;
  left: 0.9rem;
  width: 1.1rem;
  height: 1.1rem;
  transform: translateY(-50%);
  pointer-events: none;
  fill: none;
  stroke: var(--color-text-faint);
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.location-input input {
  width: 100%;
  min-height: 3.1rem;
  box-sizing: border-box;
  padding: 0 2.75rem 0 2.7rem;
  border: 1px solid var(--color-border);
  border-radius: 0.75rem;
  outline: none;
  background: var(--color-input);
  color: var(--color-text);
  font: inherit;
  font-size: 0.875rem;
}

.location-input input:focus {
  border-color: var(--color-violet-soft);
  box-shadow:
    0 0 0 3px color-mix(
      in srgb,
      var(--color-violet) 12%,
      transparent
    );
}

.location-input__clear {
  position: absolute;
  top: 50%;
  right: 0.6rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.8rem;
  height: 1.8rem;
  padding: 0;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  color: var(--color-text-faint);
  transform: translateY(-50%);
  cursor: pointer;
}

.location-input__clear:hover {
  background: var(--color-border);
  color: var(--color-text);
}

.location-input__clear svg {
  width: 0.9rem;
  height: 0.9rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
}

.location-suggestions {
  margin-top: 0.75rem;
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: 0.8rem;
  background: var(--color-input);
}

.location-suggestions__header {
  padding: 0.7rem 0.9rem;
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text-faint);
  font-size: 0.68rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.location-suggestion {
  display: flex;
  align-items: center;
  width: 100%;
  min-height: 4rem;
  gap: 0.75rem;
  padding: 0.7rem 0.9rem;
  border: 0;
  border-bottom: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text);
  cursor: pointer;
  text-align: left;
}

.location-suggestion:last-child {
  border-bottom: 0;
}

.location-suggestion:hover {
  background: var(--color-surface-raised);
}

.location-suggestion__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  flex-shrink: 0;
  border-radius: 0.65rem;
  background: color-mix(
    in srgb,
    var(--color-violet) 10%,
    var(--color-input)
  );
  color: var(--color-violet-soft);
}

.location-suggestion__icon svg {
  width: 1rem;
  height: 1rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.location-suggestion__content {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 0.15rem;
}

.location-suggestion__content strong {
  overflow: hidden;
  color: var(--color-text);
  font-size: 0.85rem;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.location-suggestion__content small {
  color: var(--color-text-faint);
  font-size: 0.72rem;
}

.location-suggestion__arrow {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  fill: none;
  stroke: var(--color-text-faint);
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.location-selected {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 0.9rem;
  padding: 0.8rem;
  border: 1px solid color-mix(
    in srgb,
    var(--color-violet) 30%,
    var(--color-border)
  );
  border-radius: 0.8rem;
  background: color-mix(
    in srgb,
    var(--color-violet) 7%,
    var(--color-input)
  );
}

.location-selected__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.3rem;
  height: 2.3rem;
  flex-shrink: 0;
  border-radius: 0.65rem;
  background: color-mix(
    in srgb,
    var(--color-violet) 14%,
    var(--color-input)
  );
  color: var(--color-violet-soft);
}

.location-selected__icon svg {
  width: 1rem;
  height: 1rem;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.location-selected__content {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 0.15rem;
}

.location-selected__content span {
  color: var(--color-text-faint);
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
}

.location-selected__content strong {
  color: var(--color-text);
  font-size: 0.85rem;
}

.location-selected__content small {
  color: var(--color-text-faint);
  font-size: 0.72rem;
}

.location-state {
  padding: 1.5rem 0.5rem;
  color: var(--color-text-faint);
  font-size: 0.8rem;
  text-align: center;
}

.location-state--error {
  color: var(--color-coral);
}

.location-dialog__footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  border-top: 1px solid var(--color-border);
}

.location-dialog__cancel,
.location-dialog__add {
  min-height: 2.75rem;
  padding: 0 1.1rem;
  border-radius: 0.7rem;
  font: inherit;
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
}

.location-dialog__cancel {
  border: 1px solid var(--color-border);
  background: var(--color-input);
  color: var(--color-text-muted);
}

.location-dialog__cancel:hover {
  background: var(--color-border);
  color: var(--color-text);
}

.location-dialog__add {
  border: 1px solid transparent;
  background: var(--gradient-action);
  color: white;
}

.location-dialog__add:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-soft);
}

.location-dialog__add:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

@media (max-width: 36rem) {
  .location-dialog {
    align-items: flex-end;
    padding: 0;
  }

  .location-dialog__panel {
    width: 100%;
    max-height: 90vh;
    border-radius: 1rem 1rem 0 0;
  }

  .location-dialog__footer {
    padding-bottom: calc(
      1rem + env(safe-area-inset-bottom)
    );
  }
}
</style>
```
