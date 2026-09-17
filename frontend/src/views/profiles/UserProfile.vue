<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getProfileData } from '@/api/users/profiles'
import { getPosts } from '@/api/posts/posts.js'
import { profileData } from '@/data/usersData'
import { profilePosts as selectProfilePosts } from '@/helpers/profilePosts.js'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import ProfileHeader from '@/components/personalProfile/ProfileHeader.vue'
import ProfileTabs from '@/components/personalProfile/ProfileTabs.vue'
import PrivateProfileIcon from '@/components/ProfileEdit/PrivateProfileIcon.vue'
import AboutTab from '@/components/Profile/AboutTab.vue'
import FollowersTab from '@/components/Profile/FollowersTab.vue'
import PostCard from '@/components/posts/PostCard.vue'

<<<<<<< HEAD
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import ProfileHeader from '@/components/personalProfile/ProfileHeader.vue';
import ProfileTabs from '@/components/personalProfile/ProfileTabs.vue';
import PrivateProfileIcon from '@/components/ProfileEdit/PrivateProfileIcon.vue';
import AboutTab from '@/components/profile/AboutTab.vue';

import { getProfileData } from '@/api/users/profiles';

import { useRoute } from 'vue-router';
import FollowersTab from '@/components/profile/FollowersTab.vue';
import { addNotification } from '@/data/notifications';
import { getFriends } from '@/api/common/friends';

const route = useRoute();

const activeTab = ref('about');
const loading = ref(true);
const showPrivateProfile = ref(false);
const user = ref(null);

async function getData() {
    const id = route.query.id;
    const count = 10;
    try {
        user.value = await getProfileData(id, count);
        showPrivateProfile.value = !user.value.show;
        const result = await getFriends("", id)
        user.value.Profile.friends = result.data;
    } catch (err) {
        addNotification('could not get user data', 'error')
        console.error(err);
    } finally {
        loading.value = false;
=======
const route = useRoute()
const activeTab = ref('posts')
const loading = ref(true)
const error = ref('')
const posts = ref([])

async function loadProfile(idValue) {
  const id = Number(idValue)
  activeTab.value = 'posts'
  posts.value = []
  error.value = ''
  if (!Number.isInteger(id) || id <= 0) {
    loading.value = false
    error.value = 'Select a valid member to view their profile.'
    return
  }

  loading.value = true
  try {
    const profileResult = await getProfileData(id, 10)
    if (!profileResult) return
    if (profileResult.showProfile) {
      const postResult = await getPosts()
      posts.value = selectProfilePosts(postResult?.posts, { ...profileData.userInfo, id })
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
    }
  } catch (err) {
    console.error(err)
    error.value = err.message || 'Could not load this profile.'
  } finally {
    loading.value = false
  }
}

<<<<<<< HEAD
function handleUnfollow() {
    if (user.value.isPrivate === 1) {
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
=======
function relationshipChanged(status) {
  profileData.isFollowing = status
  if (profileData.userInfo.isPrivate === 1 && status !== 1) {
    profileData.show = false
    posts.value = []
  }
}

watch(() => route.query.id, loadProfile, { immediate: true })
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
</script>

<template>
  <AuthenticatedLayout active-page="profile">
    <main class="profile-page">
      <p v-if="loading" class="profile-state orbit-surface">Loading profile...</p>
      <div v-else-if="error" class="profile-state profile-state--error orbit-surface" role="alert">
        <p>{{ error }}</p>
        <button type="button" @click="loadProfile(route.query.id)">Try again</button>
      </div>

      <template v-else>
        <ProfileHeader
          :first-name="profileData.userInfo.firstName"
          :last-name="profileData.userInfo.lastName"
          :username="profileData.userInfo.userName"
          :bio="profileData.about.bio"
          :avatar-path="profileData.userInfo.avatar ? `/uploads/${profileData.userInfo.avatar}` : ''"
          :num-of-posts="profileData.numOfPosts"
          :num-of-following="profileData.numOfFollowing"
          :num-of-followers="profileData.numOfFollowers"
          :is-following="profileData.isFollowing"
          :is-private="profileData.userInfo.isPrivate === 1"
          @relationship-change="relationshipChanged"
          @select-tab="activeTab = $event"
        />

<<<<<<< HEAD
            <main class="profile-page">
                <div v-if="loading">
                    Loading profile...
                </div>

                <template v-else-if="user">
                    <ProfileHeader :first-name="user.firstName"
                        :last-name="user.lastName" :username="user.username"
                        :bio="user.Profile.About?.bio" :avatar-path="`/uploads/${user.Profile.avatar}`"
                        :num-of-posts="user.Profile.numOfPosts" :num-of-following="user.Profile.numOfFollowing"
                        :num-of-followers="user.Profile.numOfFollowers" :add-edit="false"
                        :is-following="user.isFollowing" @unfollow="handleUnfollow" @follow="handleFollow" />

                    <template v-if="!showPrivateProfile">
                        <section class="profile-content">
                            <ProfileTabs v-if="user.show" type="user" @change-tab="activeTab = $event" />

                            <AboutTab v-if="user.show && activeTab === 'about'" :about="user.Profile.About" />

                            <FollowersTab :target-id="id" v-if="user.show && activeTab === 'followers'"
                                :follower-list="user.Profile.followers" />

                            <FollowersTab :target-id="id" v-if="user.show && activeTab === 'following'"
                                :follower-list="user.Profile.following" />

                            <FollowersTab :target-id="id" v-if="user.show && activeTab === 'friends'"
                                :follower-list="user.Profile.friends" />

                            <PrivateProfileIcon v-else-if="!user.show" />
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
=======
        <template v-if="profileData.show">
          <ProfileTabs v-model="activeTab" type="user" />
          <section v-if="activeTab === 'posts'" class="profile-posts" aria-labelledby="member-posts-heading">
            <header><p class="orbit-meta">Activity</p><h2 id="member-posts-heading">Posts</h2></header>
            <PostCard v-for="post in posts" :key="post.id" :post="post" />
            <p v-if="!posts.length" class="profile-empty orbit-surface">No posts available.</p>
          </section>
          <AboutTab v-else-if="activeTab === 'about'" :about="profileData.about" :profile="profileData.userInfo" />
          <FollowersTab v-else-if="activeTab === 'followers'" type="followers" :target-id="profileData.userInfo.id" :follower-list="profileData.followers" />
          <FollowersTab v-else-if="activeTab === 'following'" type="following" :target-id="profileData.userInfo.id" :follower-list="profileData.following" />
        </template>
        <PrivateProfileIcon v-else />
      </template>
    </main>
  </AuthenticatedLayout>
</template>

<style scoped>
.profile-page { display: grid; width: 100%; max-width: 64rem; margin: 0 auto; gap: var(--space-5); }
.profile-state { margin: 0; padding: var(--space-6); color: var(--color-text-muted); text-align: center; }
.profile-state p { margin: 0; }
.profile-state--error { color: var(--color-coral); }
.profile-state button { min-height: var(--touch-target); margin-top: var(--space-3); padding: 0 var(--space-4); border: 0; border-radius: var(--radius-small); background: var(--gradient-action); color: white; cursor: pointer; font-weight: 700; }
.profile-posts { display: grid; max-width: 48rem; margin-inline: auto; gap: var(--space-4); }
.profile-posts header { margin-bottom: var(--space-1); }
.profile-posts .orbit-meta { margin: 0; color: var(--color-violet-soft); }
.profile-posts h2 { margin: var(--space-1) 0 0; font-family: var(--font-display); font-size: 1.5rem; letter-spacing: 0; }
.profile-empty { margin: 0; padding: var(--space-6); color: var(--color-text-muted); text-align: center; }
</style>
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
