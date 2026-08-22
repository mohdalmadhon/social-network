<script setup>
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import PageIntroduction from '@/components/layout/PageIntroduction.vue'
</script>

<template>
  <AuthenticatedLayout active-page="profile">
    <PageIntroduction
      label="Orbit page"
      title="Profile"
      description="The Profile screen is ready for profile and follow features."
    />
  </AuthenticatedLayout>
</template>
