<script setup>
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';

import { onMounted, ref } from 'vue';
import { getUserData } from '@/api/users/personalProfile';
import { profileData } from '@/data/usersData';
import EditProfileTabs from '@/components/profile/ProfileEdit/EditProfileTabs.vue';
import EditPersonalInfo from '@/components/profile/ProfileEdit/EditPersonalInfo.vue';
import EditAdditionalInfo from '@/components/profile/ProfileEdit/EditAdditionalInfo.vue';

const activeTab = ref('personal');
const loading = ref(true);

async function getData() {
    try {
        await getUserData();
    } catch (err) {
        addNotification('could not get user data', err)
        console.error(err);
    } finally {
        loading.value = false;
    }
}

onMounted(getData);

</script>

<template>
    <div class="page-shell">
        <TopNavigation />

        <div class="page-body">
            <SideNavigation active-page="profile" />

            <main class="page-content">
                <div class="page-heading">
                    <p class="eyebrow">SETTINGS</p>
                    <h1>Edit Profile</h1>
                </div>

                <EditProfileTabs v-model:activeTab="activeTab" />

                <section class="profile-content">
                    <EditPersonalInfo v-if="!loading && activeTab === 'personal'" :first-name="profileData.userInfo.firstName"
                        :last-name="profileData.userInfo.lastName" :username="profileData.userInfo.userName" :email="profileData.userInfo.email"
                        :bio="profileData.userInfo.about" :avatar_path="`/uploads/${profileData.userInfo.avatar}`" />

                    <EditAdditionalInfo v-else-if="activeTab === 'additional'" />
                </section>
            </main>
        </div>
    </div>
</template>

<style scoped>
.page-shell {
    min-height: 100vh;
    background: var(--color-background);
    color: var(--color-text);
    font-family: var(--font-body);
}

.page-body {
    display: flex;
    align-items: flex-start;
}

.page-content {
    flex: 1;
    min-width: 0;
    width: 100%;
    max-width: 46rem;
    margin: 0 auto;
    padding: var(--space-6) var(--space-5) calc(4.25rem + var(--space-6));
}

.page-heading .eyebrow {
    margin: 0 0 var(--space-1);
    color: var(--color-violet);
    font-family: var(--font-meta);
    font-size: 0.6875rem;
    font-weight: 600;
    letter-spacing: 0.15em;
}

.page-heading h1 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: clamp(1.75rem, 4vw, 2.25rem);
    color: var(--color-text);
}

.profile-content {
    width: 100%;
    margin-top: var(--space-5);
}

@media (min-width: 64rem) {
    .page-content {
        padding-bottom: var(--space-7);
    }
}
</style>
