import { router } from "@/router";
import { userData } from "@/stores/userData";

export async function getUserData() {
    try {
        const resp = await fetch("/api/me", {
            method: "GET"
        });

        if (!resp.ok) {
            router.replace("/login");
            return;
        }

        const result = await resp.json();
        if (!result.status) {
            router.replace("/login");
            return;
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
    } catch (err) {
        ;
        console.error(err);
        router.replace("/login");
        return;
    }
}

export async function updateUserData(userData = {}) {
    try {
        const resp = await fetch("/api/user", {
            method: "UPDATE",
            body: JSON.stringify(userData)
        });

        if(!resp.ok) {
            router.replace("/login");
            return;
        }

        const result = resp.json();
        if(!result.status) {
            router.replace("/login");
            return;
        } 
    } catch(err) {
        console.error(err)
        router.replace("/login")
    }
}