<script setup>
import { onMounted, ref } from 'vue';

import { getUserData } from '@/api/users/personalProfile';
import { profileData } from '@/data/usersData';

import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import ProfileHeader from '@/components/profile/personalProfile/ProfileHeader.vue';
import ProfileTabs from '@/components/profile/personalProfile/ProfileTabs.vue';
import AboutTab from '@/components/profile/profile/AboutTab.vue';
import FollowersTab from '@/components/profile/profile/FollowersTab.vue';


const activeTab = ref('personal');
const loading = ref(true);

async function getData() {
    try {
        await getUserData();
        console.log(profileData)
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
    <div class="page-shell">
        <TopNavigation />

        <div class="page-body">
            <SideNavigation active-page="profile" />

            <main class="page-content">

                <div v-if="loading" class="loading-state">
                    Loading profile...
                </div>

                <template v-else>
                    <ProfileHeader :first-name="profileData.userInfo.firstName" :last-name="profileData.userInfo.lastName"
                        :username="profileData.userInfo.userName" :add-edit="true" :bio="profileData.userInfo.about" :avatar-path="`/uploads/${profileData.userInfo.avatar}`"
                        :num-of-posts="profileData.numOfPosts" :num-of-following="profileData.numOfFollowing"
                        :num-of-followers="profileData.numOfFollowers" />

                    <ProfileTabs @change-tab="activeTab = $event" />
                    <AboutTab v-if="activeTab === 'about'" :about="profileData.about" />

                    <FollowersTab v-if="activeTab === 'followers'" :followers="profileData.followers" />
                    <FollowersTab v-if="activeTab === 'following'" :followers="profileData.following" />
                </template>

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
    max-width: 68.75rem;
    margin: 0 auto;
    padding: var(--space-6) var(--space-5) calc(4.25rem + var(--space-6));
}

.loading-state {
    padding: var(--space-6);
    color: var(--color-text-muted);
    font-family: var(--font-meta);
    font-size: 0.8125rem;
    text-align: center;
}

@media (min-width: 64rem) {
    .page-content {
        padding-bottom: var(--space-7);
    }
}
</style>
