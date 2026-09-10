<script setup>
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import EditProfileTabs from '@/components/ProfileEdit/EditProfileTabs.vue';
import EditPersonalInfo from '@/components/ProfileEdit/EditPersonalInfo.vue';
import EditAdditionalInfo from '@/components/ProfileEdit/EditAdditionalInfo.vue';
import '@/styles/variables.css';
import '@/styles/global.css';
import { onMounted, ref } from 'vue';
import { getUserData } from '@/api/users/personalProfile';
import { profileData } from '@/data/usersData';

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

function onAvatarChange(avatar) {
    profileData.userInfo.avatar = avatar;
}

onMounted(getData);

</script>

<template>
    <div class="facebook-layout">
        <TopNavigation />

        <div class="page-layout">
            <SideNavigation />

            <main class="profile-page">
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
@import '../../styles/global.css';
@import '../../styles/variables.css';
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
    max-width: 1100px;
    margin: 0 auto;
    padding: 25px 30px 60px;
}

.page-heading .eyebrow {
    margin: 0 0 5px;
    color: var(--color-cyan);
    font-family: var(--font-meta);
    font-size: 9px;
    letter-spacing: 2px;
}

.page-heading h1 {
    margin: 0;
    color: var(--color-text);
    font-family: var(--font-body);
    font-size: 36px;
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