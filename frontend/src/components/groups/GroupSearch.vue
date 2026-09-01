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
    <input type="text" placeholder="Search groups..." :value="modelValue"
        @input="emit('update:modelValue', $event.target.value)" />
</template>

<style scoped>
input {
    width: 100%;
    padding: 10px;
    box-sizing: border-box;
}
</style>