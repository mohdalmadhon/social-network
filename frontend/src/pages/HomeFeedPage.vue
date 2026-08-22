<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import PostComposer from '@/components/posts/PostComposer.vue'
</script>

<template>
  <AuthenticatedLayout active-page="home">
    <div class="home-feed">
      <h1 class="visually-hidden">Home feed</h1>
      <PostComposer />
    </div>
  </AuthenticatedLayout>
</template>

<style scoped>
.home-feed {
  width: 100%;
  max-width: 48rem;
}
</style>
