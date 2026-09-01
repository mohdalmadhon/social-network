<script setup>
import { onMounted, ref } from 'vue'

import TopNavigation from '@/components/layout/TopNavigation.vue';
import SideNavigation from '@/components/layout/SideNavigation.vue';
import ProfileHeader from '@/components/profile/personalProfile/ProfileHeader.vue';
import { profileData } from '@/data/usersData';
import { getUserData } from '@/api/users/personalProfile';
import { addNotification } from '@/data/notifications';
import ProfileTabs from '@/components/profile/personalProfile/ProfileTabs.vue';
import ProfilePosts from '@/components/profile/personalProfile/ProfilePosts.vue';
import ProfileFriends from '@/components/profile/personalProfile/ProfileFriends.vue';
import ProfileGroups from '@/components/profile/personalProfile/ProfileGroups.vue';
import ProfileFollowing from '@/components/profile/personalProfile/ProfileFollowing.vue';
import AboutTab from '@/components/profile/Profile/AboutTab.vue';
import FollowersTab from '@/components/profile/Profile/FollowersTab.vue';

function handleChangeTab(tab) {
    activeTab.value = tab
}

const activeTab = ref('about');
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

onMounted(getData)
</script>

<template>
    <header>
        <TopNavigation />
    </header>

    <div class="app-body">
        <SideNavigation active-page="personalPage" />

        <div v-if="loading">
            Loading profile...
        </div>
        <main v-else>
            <ProfileHeader add-edit :first-name="profileData.userInfo.firstName"
                :last-name="profileData.userInfo.lastName" :username="profileData.userInfo.username"
                :bio="profileData.about.bio" :avatar-path="`/uploads/${profileData.userInfo.avatar}`"
                :num-of-posts="profileData.numOfPosts" :num-of-following="profileData.numOfFollowing"
                :num-of-followers="profileData.numOfFollowers" />

            <ProfileTabs @change-tab="handleChangeTab" />

            <div class="profile-content">
                <ProfilePosts v-if="activeTab === 'posts'" />
                <ProfileFriends v-else-if="activeTab === 'friends'" />
                <ProfileGroups v-else-if="activeTab === 'groups'" />
                <FollowersTab :followers="profileData.following" v-else-if="activeTab === 'following'" />
                <FollowersTab :followers="profileData.following" v-else-if="activeTab === 'followers'" />
                <AboutTab :about="profileData.about" v-else-if="activeTab === 'about'" />
            </div>
        </main>
    </div>
</template>

<style scoped>
.app-body {
    display: flex;
    align-items: flex-start;
}

main {
    flex: 1;
    min-width: 0;
    padding: 24px 32px;
    max-width: 1240px;
    margin: 0 auto;
    box-sizing: border-box;
}

.profile-content {
    margin-top: 24px;
}
</style>