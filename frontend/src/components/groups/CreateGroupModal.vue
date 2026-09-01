<script setup>
import { ref } from 'vue'

defineProps({
    show: {
        type: Boolean,
        default: false
    }
})

const emit = defineEmits(['close', 'create'])

const title = ref('')
const description = ref('')

function submitGroup() {
    if (!title.value.trim() || !description.value.trim()) {
        console.log("Error: the title and description could not be empty!")
        return
    }

    emit('create', {
        title: title.value,
        description: description.value
    })

    title.value = ''
    description.value = ''
}
</script>

<template>
    <div v-if="show" class="modal-overlay" @click.self="emit('close')">
        <div class="modal">
            <h2>Create Group</h2>

            <form @submit.prevent="submitGroup">
                <label>
                    Group Title
                </label>

                <input v-model="title" type="text" placeholder="Enter group title" />

                <label>
                    Description
                </label>

                <textarea v-model="description" placeholder="Enter group description"></textarea>

                <div class="actions">
                    <button type="button" @click="emit('close')">
                        Cancel
                    </button>

                    <button type="submit">
                        Create
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped>
.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);

    display: flex;
    justify-content: center;
    align-items: center;
}

.modal {
    width: 400px;
    padding: 20px;
    background: white;
    border-radius: 10px;
}

form {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

input,
textarea {
    padding: 10px;
}

textarea {
    min-height: 100px;
    resize: vertical;
}

.actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}

button {
    padding: 8px 14px;
    cursor: pointer;
}
</style>