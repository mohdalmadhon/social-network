<script setup>
import { getUserData } from '@/api/users'
import ProfileAvatarCard from '@/components/editProfile/ProfileAvatarCard.vue'
import ProfileDetailsForm from '@/components/editProfile/ProfileDetailsForm.vue'
import SideNavigation from '@/components/layout/SideNavigation.vue'
import TopNavigation from '@/components/layout/TopNavigation.vue'

getUserData();
</script>

<template>
  <div class="app-shell">
    <TopNavigation />

    <div class="app-body">
      <SideNavigation active-page="profile" />

      <main class="edit-profile">
        <header class="edit-profile__header">
          <div>
            <p class="edit-profile__eyebrow">Profile</p>
            <h1 class="edit-profile__title">Edit profile</h1>
            <p class="edit-profile__subtitle">Update your photo and personal details.</p>
          </div>

          <div class="edit-profile__header-actions">
            <button class="btn btn--ghost" type="button">Cancel</button>
          </div>
        </header>

        <div class="edit-profile__layout">
          <ProfileAvatarCard />
          <ProfileDetailsForm />
        </div>

        <div class="edit-profile__footer-actions">
          <button class="btn btn--ghost btn--block" type="button">Cancel</button>
          <button class="btn btn--primary btn--block" type="button">Save changes</button>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
@import '../styles/global.css';
@import '../styles/variables.css';

.app-shell {
  min-height: 100vh;
  background: var(--color-background);
}

.app-body {
  display: flex;
  align-items: flex-start;
}

.edit-profile {
  flex: 1;
  min-width: 0;
  padding: var(--space-4) var(--space-3) calc(var(--space-7) + 4.25rem);
}

.edit-profile__header {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-3);
  padding-inline: var(--space-1);
}

.edit-profile__eyebrow {
  color: var(--color-violet);
  font-family: var(--font-meta);
  font-size: 0.75rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.edit-profile__title {
  margin-top: var(--space-1);
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: 1.85rem;
  font-weight: 700;
}

.edit-profile__subtitle {
  margin-top: var(--space-1);
  color: var(--color-text-muted);
  font-size: 0.95rem;
}

.edit-profile__header-actions {
  display: none;
  gap: var(--space-3);
}

.edit-profile__layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-5);
  margin-top: var(--space-5);
}

.edit-profile__footer-actions {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin-top: var(--space-5);
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: var(--touch-target);
  padding-inline: var(--space-5);
  border: 1px solid transparent;
  border-radius: var(--radius-small);
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  transition: filter 0.15s ease, border-color 0.15s ease, color 0.15s ease;
}

.btn--block {
  width: 100%;
}

.btn--ghost {
  background: var(--color-input);
  border-color: var(--color-border);
  color: var(--color-text-soft);
}

.btn--ghost:hover,
.btn--ghost:focus-visible {
  border-color: var(--color-violet);
  color: var(--color-text);
  outline: none;
}

.btn--primary {
  background: var(--gradient-action);
  color: white;
  box-shadow: var(--shadow-raised);
}

.btn--primary:hover,
.btn--primary:focus-visible {
  filter: brightness(1.08);
  outline: none;
}

@media (min-width: 48rem) {
  .edit-profile {
    padding-inline: var(--space-5);
  }
}

@media (min-width: 64rem) {
  .edit-profile {
    padding: var(--space-6) var(--space-6) var(--space-7);
  }

  .edit-profile__header-actions {
    display: flex;
  }

  .edit-profile__footer-actions {
    display: none;
  }

  .edit-profile__layout {
    grid-template-columns: 20rem 1fr;
    align-items: start;
  }
}

@media (min-width: 90rem) {
  .edit-profile {
    padding-inline: var(--space-7);
  }
}
</style>