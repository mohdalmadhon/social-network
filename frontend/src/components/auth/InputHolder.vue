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
    height: 38px;

    padding: 0 13px;

    border: 1px solid transparent;

    border-radius: 7px;

    outline: none;

    background: #222741;

    color: #f4f3ff;

    font-family: inherit;

    font-size: 12px;
    font-weight: 400;

    transition: 0.2s ease;
}

.flip-card__input::placeholder {
    color: #59627f;
}

.flip-card__input:focus {
    border-color: #835cff;

    box-shadow: 0 0 0 1px rgba(131, 92, 255, 0.15);
}
</style>
