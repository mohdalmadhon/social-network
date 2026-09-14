<script setup>
import ChatsSideBar from '@/components/chats/ChatsSideBar.vue';
import ChatWindow from '@/components/chats/ChatWindow.vue';
import SideNavigation from '@/components/layout/SideNavigation.vue';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import { activePage } from '@/data/chatState';
import { ref } from 'vue';

activePage.value = 'chat:';

const activeChat = ref(null);

function handleSelectChat(chat) {
    console.log(chat)
    activeChat.value = chat;
    activePage.value = 'chat:' + chat.UserID;

    console.log(activeChat.value);
}
</script>

<template>
    <div class="app-shell">
        <TopNavigation />

        <div class="body-layout">
            <SideNavigation />

            <main class="chats-page">
                <ChatsSideBar @select-chat="handleSelectChat" />

                <ChatWindow
                    :chat="activeChat"
                    :userID="activeChat?.UserID"
                    :groupID="activeChat?.GroupID"
                />
            </main>
        </div>
    </div>
</template>

<style scoped>
.app-shell {
    min-height: 100vh;
}

.body-layout {
    display: flex;
    padding-top: 64px;
}

.chats-page {
    flex: 1;
    display: flex;
    gap: 20px;
    padding: 20px clamp(16px, 3vw, 30px);
    box-sizing: border-box;
    min-width: 0;
}

@media (max-width: 800px) {
    .chats-page {
        flex-direction: column;
    }
}
</style>