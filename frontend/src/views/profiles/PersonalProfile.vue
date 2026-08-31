<script setup>
import { onMounted, ref } from 'vue';

import { getUserData } from '@/api/users/personalProfile';
import { profileData } from '@/data/usersData';

import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';

import ProfileAbout from '@/components/personalProfile/ProfileAbout.vue';
import ProfileFollowers from '@/components/personalProfile/ProfileFollowers.vue';
import ProfileFollowing from '@/components/personalProfile/ProfileFollowing.vue';
import ProfileFriends from '@/components/personalProfile/ProfileFriends.vue';
import ProfileGroups from '@/components/personalProfile/ProfileGroups.vue';
import ProfileHeader from '@/components/personalProfile/ProfileHeader.vue';
import ProfilePosts from '@/components/personalProfile/ProfilePosts.vue';
import ProfileTabs from '@/components/personalProfile/ProfileTabs.vue';
import AboutTab from '@/components/profile/AboutTab.vue';
import FollowersTab from '@/components/profile/FollowersTab.vue';

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
    <div class="facebook-layout">
        <TopNavigation />

        <div class="page-layout">
            <SideNavigation />

            <main class="profile-page">

                <div v-if="loading">
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