<script setup>
import { onMounted, ref } from 'vue'
import { getUserData } from '@/api/users/personalProfile'
import { getPosts } from '@/api/posts/posts.js'
import { profileData } from '@/data/usersData'
import { profilePosts as selectProfilePosts } from '@/helpers/profilePosts.js'
import AuthenticatedLayout from '@/components/layout/AuthenticatedLayout.vue'
import ProfileHeader from '@/components/personalProfile/ProfileHeader.vue'
import ProfileTabs from '@/components/personalProfile/ProfileTabs.vue'
import AboutTab from '@/components/Profile/AboutTab.vue'
import FollowersTab from '@/components/Profile/FollowersTab.vue'
import PostCard from '@/components/posts/PostCard.vue'

const activeTab = ref('posts')
const loading = ref(true)
const error = ref('')
const posts = ref([])

onMounted(async () => {
  try {
    await getUserData()
    const result = await getPosts()
    posts.value = selectProfilePosts(result?.posts, profileData.userInfo)
  } catch (err) {
    console.error(err)
    error.value = err.message || 'Could not load your profile.'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <AuthenticatedLayout active-page="profile">
    <main class="profile-page">
      <p v-if="loading" class="profile-state orbit-surface">Loading profile...</p>
      <div v-else-if="error" class="profile-state profile-state--error orbit-surface" role="alert">
        <p>{{ error }}</p>
        <button type="button" @click="$router.go(0)">Try again</button>
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
          :is-private="profileData.userInfo.isPrivate === 1"
          add-edit
          @select-tab="activeTab = $event"
        />

        <ProfileTabs v-model="activeTab" type="personal" />

        <section v-if="activeTab === 'posts'" class="profile-posts" aria-labelledby="profile-posts-heading">
          <header><p class="orbit-meta">Activity</p><h2 id="profile-posts-heading">Your posts</h2></header>
          <PostCard v-for="post in posts" :key="post.id" :post="post" />
          <p v-if="!posts.length" class="profile-empty orbit-surface">No posts yet.</p>
        </section>

        <AboutTab v-else-if="activeTab === 'about'" :about="profileData.about" :profile="profileData.userInfo" own-profile />
        <FollowersTab v-else-if="activeTab === 'followers'" type="followers" :target-id="profileData.userInfo.id" :follower-list="profileData.followers" />
        <FollowersTab v-else-if="activeTab === 'following'" type="following" :target-id="profileData.userInfo.id" :follower-list="profileData.following" />
        <FollowersTab v-else-if="activeTab === 'friends'" type="friends" :target-id="profileData.userInfo.id" :follower-list="profileData.friends" />
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
