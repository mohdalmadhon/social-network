<script>
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import ProfileHeader from '@/components/profile/ProfileHeader.vue';
import ProfileTabs from '@/components/profile/ProfileTabs.vue';
import ProfileAbout from '@/components/profile/ProfileAbout.vue';
import ProfileFollowers from '@/components/profile/ProfileFollowers.vue';

export default {
    components: {
        SideNavigation,
        TopNavigation,
        ProfileHeader,
        ProfileTabs,
        ProfileAbout,
        ProfileFollowers,
    },
};
</script>

<template>
    <header>
        <TopNavigation/>
    </header>
    <div class="app-body">
        <SideNavigation/>
        <main>
            <ProfileHeader/>
            <ProfileTabs/>
            <div class="profile-content">
                <ProfileAbout/>
                <ProfileFollowers/>
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