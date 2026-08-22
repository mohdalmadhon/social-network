<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import PageIntroduction from '@/components/layout/PageIntroduction.vue'
</script>

<template>
  <AuthenticatedLayout active-page="chats">
    <PageIntroduction
      label="Orbit page"
      title="Chats"
      description="The Chats screen is ready for private and group conversations."
    />
  </AuthenticatedLayout>
</template>
