<script setup>
defineProps(['name', 'maxLength', 'minLength', 'type', 'id', 'placeHolder', 'required']);
</script>

<template>
    <input :name="name" :maxlength="maxLength" :minlength="minLength" :type="type || 'text'" :id="id"
        :placeholder="placeHolder" :required="required" class="flip-card__input">
</template>

<style>
.flip-card__input {
    width: 100%;
    min-height: 44px;

    padding: 0 13px;

    border: 1px solid transparent;

    border-radius: 7px;

    outline: none;

    background: var(--color-input);

    color: var(--color-text);

    font-family: inherit;

    font-size: 0.875rem;
    font-weight: 400;

    transition: 0.2s ease;
}

.flip-card__input::placeholder {
    color: var(--color-text-faint);
}

.flip-card__input:focus {
    border-color: var(--color-violet);

    box-shadow: var(--focus-ring);
}
</style>
