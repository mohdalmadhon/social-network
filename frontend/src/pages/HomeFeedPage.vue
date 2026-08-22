<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import PostCard from '@/components/posts/PostCard.vue'
import PostComposer from '@/components/posts/PostComposer.vue'

const samplePosts = [
  {
    id: 1,
    author: 'Kiko Tanaka',
    avatarColor: '#45d9d0',
    time: '2h ago',
    privacy: 'Followers only',
    content: 'Summit at 6:12 am — worth every step of the climb. ⛰️',
    likes: 48,
    comments: 12,
    hasMedia: true,
    mediaDescription: 'A colorful sunrise over a layered mountain range',
  },
  {
    id: 2,
    author: 'Mara Voss',
    avatarColor: '#ff8b5c',
    time: '5h ago',
    privacy: 'Public',
    content: 'Migrating the social-network database tonight: SQLite + golang-migrate, wish me luck. ☕',
    likes: 17,
    comments: 3,
    hasMedia: false,
    mediaDescription: '',
  },
]
</script>

<template>
  <AuthenticatedLayout active-page="home">
    <div class="home-feed">
      <h1 class="visually-hidden">Home feed</h1>
      <PostComposer />
      <PostCard v-for="post in samplePosts" :key="post.id" :post="post" />
    </div>
  </AuthenticatedLayout>
</template>

<style scoped>
.home-feed {
  display: grid;
  width: 100%;
  max-width: 48rem;
  gap: var(--space-4);
}

@media (min-width: 48rem) {
  .home-feed {
    gap: var(--space-5);
  }
}
</style>
