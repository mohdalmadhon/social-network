<script setup>
import { ref, onMounted } from 'vue';

import HomePosts from '@/components/home/HomePosts.vue';
import HomeSearch from '@/components/home/HomeSearch.vue';
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';

const posts = ref([]);
const loading = ref(false);
const hasMore = ref(true);
const offset = ref(0);

async function loadPosts() {
    if (loading.value || !hasMore.value) {
        return;
    }

    loading.value = true;

    try {
        const response = await fetch(`/api/posts?offset=${offset.value}`, {
            credentials: 'include'
        });
        const data = await response.json();
        console.log(data)
        if (!response.ok || !data.status) {
            return;
        }

        const newPosts = data.posts || [];

        posts.value.push(...newPosts);

        offset.value += newPosts.length;

        if (newPosts.length < 13) {
            hasMore.value = false;
        }
    } catch (error) {
        console.error(error);
    } finally {
        loading.value = false;
    }
}

onMounted(() => {
    loadPosts();
});
</script>

<template>
    <div class="home-page">
        <TopNavigation />

        <div class="page-layout">
            <SideNavigation />

            <main class="main-content">
                <div class="content-container">
                    <HomeSearch />

                    <section class="posts">
                        <HomePosts
                            v-for="post in posts"
                            :key="post.id"
                            v-bind="post"
                        />
                    </section>
                </div>
            </main>
        </div>
    </div>
</template>

<style scoped>
.home-page {
    min-height: 100vh;
    padding-top: 64px;
}

.page-layout {
    display: flex;
    align-items: flex-start;
    min-height: calc(100vh - 64px);
}

.main-content {
    flex: 1;
    min-width: 0;
}

.content-container {
    width: 100%;
    max-width: 760px;
    margin: 0 auto;
    padding: 30px 25px 60px;
}

.posts {
    display: flex;
    flex-direction: column;
    gap: 28px;
    margin-top: 30px;
}

@media (max-width: 800px) {
    .page-layout {
        display: block;
    }

    .content-container {
        padding: 20px 15px 50px;
    }
}

@media (max-width: 650px) {
    .content-container {
        padding-left: 10px;
        padding-right: 10px;
    }
}
</style>