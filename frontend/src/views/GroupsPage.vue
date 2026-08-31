<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import PageIntroduction from '@/components/layout/PageIntroduction.vue'
</script>

<template>
  <AuthenticatedLayout active-page="groups">
    <PageIntroduction
      label="Orbit page"
      title="Groups"
      description="The Groups screen is ready for groups, invitations, requests, and events."
    />
  </AuthenticatedLayout>
</template>
