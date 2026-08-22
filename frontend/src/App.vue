<script setup>
    import Login from './pages/Login.vue';

    async function checkToken() {
        try {
            const resp = await fetch("/api/session", {
                method: "POST",
            });

            if(!resp.ok) {
                // to login
            }

            const result = await resp.json();
            if(!result.status) {
                // to login
            }
        } catch (err) {

        }
    }
</script>

<template> 
    <RouterView/>
</template>

<style>
</style>