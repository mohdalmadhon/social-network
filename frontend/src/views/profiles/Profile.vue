<script setup>
import { onMounted, ref } from 'vue';

import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import ProfileHeader from '@/components/personalProfile/ProfileHeader.vue';
import ProfileTabs from '@/components/personalProfile/ProfileTabs.vue';
import PrivateProfileIcon from '@/components/ProfileEdit/PrivateProfileIcon.vue';
import AboutTab from '@/components/profile/AboutTab.vue';

import { getProfileData } from '@/api/users/profiles';
import { profileData } from '@/data/usersData';

import { useRoute } from 'vue-router';
import FollowersTab from '@/components/profile/FollowersTab.vue';
import { addNotification } from '@/data/notifications';
import { getFriends } from '@/api/common/friends';

const route = useRoute();

const activeTab = ref('about');
const loading = ref(true);
const showPrivateProfile = ref(false);

async function getData() {
    const id = route.query.id;
    const count = 10;
    try {
        await getProfileData(id, count);
        showPrivateProfile.value = !profileData.show;
        const result = await getFriends("", id)
        profileData.friends = result.data;
    } catch (err) {
        addNotification('could not get user data', 'error')
        console.error(err);
    } finally {
        loading.value = false;
    }
}

function handleUnfollow() {
    if (profileData.userInfo.isPrivate === 1) {
        showPrivateProfile.value = true;
    }
}

function handleFollow() {
    showPrivateProfile.value = false;
}

let id = route.query.id;
if (!id) {
    id = ""
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
                    <ProfileHeader :first-name="profileData.userInfo.firstName"
                        :last-name="profileData.userInfo.lastName" :username="profileData.userInfo.userName"
                        :bio="profileData.about.bio" :avatar-path="`/uploads/${profileData.userInfo.avatar}`"
                        :num-of-posts="profileData.numOfPosts" :num-of-following="profileData.numOfFollowing"
                        :num-of-followers="profileData.numOfFollowers" :add-edit="false"
                        :is-following="profileData.isFollowing" @unfollow="handleUnfollow" @follow="handleFollow" />

                    <template v-if="!showPrivateProfile">
                        <section class="profile-content">
                            <ProfileTabs v-if="profileData.show" type="user" @change-tab="activeTab = $event" />

                            <AboutTab v-if="profileData.show && activeTab === 'about'" :about="profileData.about" />

                            <FollowersTab :target-id="id" v-if="profileData.show && activeTab === 'followers'"
                                :follower-list="profileData.followers" />

                            <FollowersTab :target-id="id" v-if="profileData.show && activeTab === 'following'"
                                :follower-list="profileData.following" />

                            <FollowersTab :target-id="id" v-if="profileData.show && activeTab === 'friends'"
                                :follower-list="profileData.friends" />

                            <PrivateProfileIcon v-else-if="!profileData.show" />
                        </section>
                    </template>

                    <section v-else class="profile-content">
                        <PrivateProfileIcon />
                    </section>
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
