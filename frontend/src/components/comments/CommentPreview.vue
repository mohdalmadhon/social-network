<script setup>
defineProps({
  comment: {
    type: Object,
    required: true,
  },
})

function imageUrl(imagePath) {
  if (!imagePath) return ''
  return imagePath.startsWith('/') ? imagePath : `/uploads/${imagePath}`
}

function initials(author) {
  return author.slice(0, 2).toUpperCase()
}
</script>

<template>
  <div class="comment-preview">
    <div class="comment-preview__avatar" :style="{ background: comment.avatarColor }" aria-hidden="true">
      <img v-if="comment.avatarPath" :src="imageUrl(comment.avatarPath)" alt="" />
      <span v-else>{{ initials(comment.author) }}</span>
    </div>
    <div class="comment-preview__body">
      <strong>{{ comment.author }}</strong>
      <p>{{ comment.content }}</p>
    </div>
  </div>
</template>

<style scoped>
.comment-preview {
  display: grid;
  grid-template-columns: 2.25rem minmax(0, 1fr);
  align-items: start;
  gap: var(--space-3);
  margin-top: var(--space-3);
  padding: var(--space-3);
  border-radius: var(--radius-small);
  background: var(--color-input);
}

.comment-preview__avatar {
  display: grid;
  width: 2.25rem;
  height: 2.25rem;
  aspect-ratio: 1;
  place-items: center;
  overflow: hidden;
  border-radius: 50%;
  background: var(--color-input);
  color: var(--color-text);
  font-size: 0.8125rem;
  font-weight: 700;
}

.comment-preview__avatar img {
  width: 100%;
  height: 100%;
  border-radius: inherit;
  object-fit: cover;
  object-position: center;
}

.comment-preview__body {
  min-width: 0;
}

.comment-preview strong,
.comment-preview p {
  margin: 0;
}

.comment-preview strong {
  color: var(--color-text);
  font-size: 0.875rem;
  font-weight: 500;
}

.comment-preview p {
  color: var(--color-text-soft);
  font-size: 0.875rem;
  line-height: 1.45;
  overflow-wrap: anywhere;
}

</style>
