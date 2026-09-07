<script setup>
import { onMounted, ref } from 'vue';

import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';

import { getProfileData } from '@/api/users/profiles';
import { profileData } from '@/data/usersData';

import { useRoute } from 'vue-router';
import { addNotification } from '@/data/notifications';
import ProfileHeader from '@/components/profile/personalProfile/ProfileHeader.vue';
import ProfileTabs from '@/components/profile/personalProfile/ProfileTabs.vue';
import AboutTab from '@/components/profile/profile/AboutTab.vue';
import FollowersTab from '@/components/profile/profile/FollowersTab.vue';

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
    <div class="page-shell">
        <TopNavigation />

        <div class="page-body">
            <SideNavigation active-page="profile" />

            <main class="page-content">
                <div v-if="loading" class="loading-state">
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

.profile-content {
    width: 100%;
}

@media (min-width: 64rem) {
    .page-content {
        padding-bottom: var(--space-7);
    }
}
</style>
