<script setup>
<<<<<<< HEAD
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import EditProfileTabs from '@/components/ProfileEdit/EditProfileTabs.vue';
import EditPersonalInfo from '@/components/ProfileEdit/EditPersonalInfo.vue';
import EditAdditionalInfo from '@/components/ProfileEdit/EditAdditionalInfo.vue';

import { onMounted, ref } from 'vue';
import { getUserData } from '@/api/users/personalProfile';
import { addNotification } from '@/data/notifications';

const activeTab = ref('personal');
const loading = ref(true);
const user = ref(null);

async function getData() {
    try {
        user.value = await getUserData();
    } catch (err) {
        addNotification('could not get user data', 'error')
        console.error(err);
    } finally {
        loading.value = false;
    }
}

onMounted(getData);

</script>

<template>
    <div class="facebook-layout">
        <TopNavigation />

        <div class="page-layout">
            <SideNavigation />

            <main class="profile-page">
                <div class="page-heading">
                    <p class="eyebrow">SETTINGS</p>
                    <h1>Edit Profile</h1>
                </div>

                <EditProfileTabs v-model:activeTab="activeTab" />

                <section class="profile-content">
                    <template v-if="!loading && user">
                        <EditPersonalInfo v-if="activeTab === 'personal'" :first-name="user.firstName"
                            :last-name="user.lastName" :username="user.username" :email="user.email"
                            :bio="user.Profile.About?.bio" :avatar_path="`/uploads/${user.Profile.avatar}`"
                            :is-private="user.isPrivate" />

                        <EditAdditionalInfo v-else-if="activeTab === 'additional'" :about="user.Profile.About" />
                    </template>
                </section>
            </main>
=======
import { onMounted, ref } from 'vue'
import { getUserData } from '@/api/users/personalProfile'
import { profileData } from '@/data/usersData'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import EditProfileTabs from '@/components/ProfileEdit/EditProfileTabs.vue'
import EditPersonalInfo from '@/components/ProfileEdit/EditPersonalInfo.vue'
import EditAdditionalInfo from '@/components/ProfileEdit/EditAdditionalInfo.vue'

const activeTab = ref('personal')
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    await getUserData()
  } catch (err) {
    console.error(err)
    error.value = err.message || 'Could not load your profile settings.'
  } finally {
    loading.value = false
  }
})

function onAvatarChange(avatarPath) {
  if (avatarPath) profileData.userInfo.avatar = avatarPath
}
</script>

<template>
  <AuthenticatedLayout active-page="profile">
    <main class="edit-profile-page">
      <header class="edit-profile-heading">
        <div>
          <p class="orbit-meta">Profile settings</p>
          <h1>Edit profile</h1>
          <p>Keep your identity, visibility, and profile details up to date.</p>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
        </div>
        <RouterLink to="/me">Back to profile</RouterLink>
      </header>

      <p v-if="loading" class="edit-state orbit-surface">Loading profile settings...</p>
      <p v-else-if="error" class="edit-state edit-state--error orbit-surface" role="alert">{{ error }}</p>

      <template v-else>
        <EditProfileTabs v-model:active-tab="activeTab" />
        <EditPersonalInfo
          v-if="activeTab === 'personal'"
          :first-name="profileData.userInfo.firstName"
          :last-name="profileData.userInfo.lastName"
          :username="profileData.userInfo.userName"
          :email="profileData.userInfo.email"
          :bio="profileData.about.bio"
          :avatar-path="profileData.userInfo.avatar ? `/uploads/${profileData.userInfo.avatar}` : ''"
          :is-private="profileData.userInfo.isPrivate === 1"
          @avatar-change="onAvatarChange"
        />
        <EditAdditionalInfo v-else />
      </template>
    </main>
  </AuthenticatedLayout>
</template>

<style scoped>
<<<<<<< HEAD
.facebook-layout {
    min-height: 100vh;
}

.page-layout {
    display: flex;
    padding-top: 64px;
}

.profile-page {
    width: 100%;
    max-width: 1100px;
    margin: 0 auto;
    padding: 25px 30px 60px;
}

.page-heading .eyebrow {
    margin: 0 0 5px;
    color: var(--input-focus);
    font-family: "JetBrains Mono", monospace;
    font-size: 9px;
    letter-spacing: 2px;
}

.page-heading h1 {
    margin: 0;
    font-family: "Liter", serif;
    font-size: 36px;
}

.profile-content {
    width: 100%;
}

@media (max-width: 800px) {
    .page-layout {
        display: block;
    }

    .profile-page {
        padding: 20px 15px 50px;
    }
}
</style>
=======
.edit-profile-page { display: grid; width: 100%; max-width: 58rem; margin: 0 auto; gap: var(--space-5); }
.edit-profile-heading { display: flex; align-items: end; justify-content: space-between; gap: var(--space-4); }
.edit-profile-heading .orbit-meta { margin: 0; color: var(--color-violet-soft); }
.edit-profile-heading h1 { margin: var(--space-2) 0; font-family: var(--font-display); font-size: 2.25rem; line-height: 1.15; letter-spacing: 0; }
.edit-profile-heading p:last-child { margin: 0; color: var(--color-text-muted); }
.edit-profile-heading a { display: inline-flex; min-height: var(--touch-target); flex: 0 0 auto; align-items: center; padding: 0 var(--space-4); border: 1px solid var(--color-border); border-radius: var(--radius-small); color: var(--color-text-soft); text-decoration: none; }
.edit-profile-heading a:hover { border-color: var(--color-violet); background: var(--color-input); }
.edit-state { margin: 0; padding: var(--space-6); color: var(--color-text-muted); text-align: center; }
.edit-state--error { color: var(--color-coral); }
@media (max-width: 600px) {
  .edit-profile-heading { align-items: stretch; flex-direction: column; }
  .edit-profile-heading h1 { font-size: 1.875rem; }
  .edit-profile-heading a { justify-content: center; }
}
</style>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
