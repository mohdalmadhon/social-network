<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

import SideNavigation from '@/components/layout/SideNavigation.vue'
import TopNavigation from '@/components/layout/TopNavigation.vue'
import ProfileHeader from '@/components/profile/ProfileHeader.vue'
import ProfileTabs from '@/components/profile/ProfileTabs.vue'
import ProfileAbout from '@/components/profile/ProfileAbout.vue'
import ProfileFollowers from '@/components/profile/ProfileFollowers.vue'
import { userData } from '@/stores/userData'

const router = useRouter()

async function getData() {
    try {
        const resp = await fetch('/api/me', {
            method: 'GET',
            credentials: 'include'
        })

        const result = await resp.json()


        if (!resp.ok || !result.status) {
            router.replace('/login')
            return
        }

        userData.value = {
            ...userData.value,
            ...result.data.user,
        }

        userData.value.numOfFollowers = result.data.profile.Followers;
        userData.value.numOfFollowing = result.data.profile.Following;
        userData.value.numOfPosts = result.data.profile.Posts;
        userData.value.about = result.data.profile.About;
        if (result.data.profile.Avatar_Path) {
            userData.value.avatar_path = 'http://localhost:4031' +
                result.data.profile.Avatar_Path
                    .replaceAll('\\', '/')
                    .replace('..', '')
        }
        console.log(userData)
    } catch (err) {
        console.error(err)
        router.replace('/login')
    }
}

onMounted(() => {
    getData()
})
</script>

<template>
    <header>
        <TopNavigation />
    </header>

    <div class="app-body">
        <SideNavigation />

        <main>
            <ProfileHeader />
            <ProfileTabs />

            <div class="profile-content">
                <ProfileAbout />
                <ProfileFollowers />
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
    display: grid;
    grid-template-columns: 320px 1fr;
    gap: 24px;
    margin-top: 24px;
}
</style>