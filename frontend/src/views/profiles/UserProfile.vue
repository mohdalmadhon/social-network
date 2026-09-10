<script setup>
import { onMounted, ref } from 'vue';

import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import ProfileHeader from '@/components/personalProfile/ProfileHeader.vue';
import ProfileTabs from '@/components/personalProfile/ProfileTabs.vue';
import PrivateProfileIcon from '@/components/ProfileEdit/PrivateProfileIcon.vue';

import { getProfileData } from '@/api/users/profiles';
import { profileData } from '@/data/usersData';

import { useRoute } from 'vue-router';
import { addNotification } from '@/data/notifications';
import AboutTab from '@/components/Profile/AboutTab.vue';
import FollowersTab from '@/components/Profile/FollowersTab.vue';

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
                    <ProfileHeader
                        :first-name="profileData.userInfo.firstName"
                        :last-name="profileData.userInfo.lastName"
                        :username="profileData.userInfo.userName"
                        :bio="profileData.about.bio"
                        :avatar-path="`/uploads/${profileData.userInfo.avatar}`"
                        :num-of-posts="profileData.numOfPosts"
                        :num-of-following="profileData.numOfFollowing"
                        :num-of-followers="profileData.numOfFollowers"
                        :add-edit="false"
                        :is-following="profileData.isFollowing"
                        @unfollow="handleUnfollow"
                        @follow="handleFollow"
                    />

                    <template v-if="!showPrivateProfile">
                        <section class="profile-content">
                            <ProfileTabs
                                v-if="profileData.show"
                                @change-tab="activeTab = $event"
                            />

                            <AboutTab
                                v-if="profileData.show && activeTab === 'about'"
                                :about="profileData.about"
                            />

                            <FollowersTab
                                v-if="profileData.show && activeTab === 'followers'"
                                :followers="profileData.followers"
                            />

                            <FollowersTab
                                v-if="profileData.show && activeTab === 'following'"
                                :followers="profileData.following"
                            />
                            
                            <PrivateProfileIcon
                                v-else-if="!profileData.show"
                            />
                        </section>
                    </template>

                    <section
                        v-else
                        class="profile-content"
                    >
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
    background: var(--color-background);
}

.page-layout {
    display: flex;
    padding-top: 64px;
}

.profile-page {
    width: 100%;
    max-width: var(--content-max-width);
    margin: 0 auto;
    padding: var(--space-6) var(--space-6) var(--space-7);
}

.profile-content {
    width: 100%;
}

@media (max-width: 800px) {
    .page-layout {
        display: block;
    }

    .profile-page {
        padding: var(--space-5) var(--space-4) var(--space-6);
    }
}
</style>