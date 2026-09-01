<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import PageIntroduction from '@/components/layout/PageIntroduction.vue'
</script>

<template>
  <AuthenticatedLayout active-page="notifications">
    <PageIntroduction
      label="Orbit page"
      title="Notifications"
      description="The Notifications screen is ready for notification lists, filters, and actions."
    />
  </AuthenticatedLayout>
</template>
