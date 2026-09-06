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
</script>

<template>
  <div class="comment-preview">
    <div class="comment-preview__avatar" :style="{ background: comment.avatarColor }" aria-hidden="true">
      {{ comment.author.charAt(0) }}
    </div>
    <div class="comment-preview__body">
      <strong>{{ comment.author }}</strong>
      <p>{{ comment.content }}</p>
      <img v-if="comment.imagePath" :src="imageUrl(comment.imagePath)" alt="Image attached to this comment" />
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
  place-items: center;
  border-radius: 50%;
  color: var(--color-text);
  font-size: 0.8125rem;
  font-weight: 700;
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

.comment-preview img {
  display: block;
  width: min(100%, 20rem);
  max-height: 16rem;
  margin-top: var(--space-2);
  border-radius: var(--radius-small);
  object-fit: contain;
}
</style>
