import { profileData } from "@/data/usersData";
import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";

export async function getProfileData(id, count) {
    const resp = await fetch(`/api/profile?id=${id}`, {
        method: "GET",
        credentials: 'include',
    });

    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return;
    }

    const result = await resp.json();

    if (!resp.ok) {
        if (result.message == "cannot view your own profile this way") {
            router.push("/me")
            return;
        }
        throw new Error("network error, could not connet to server")
    }

    profileData.show = result.showProfile;
    profileData.isFollowing = result.followStatus;
    if (profileData.show) {
        profileData.userInfo.firstName = result.data.UserInfo.FirstName
        profileData.userInfo.lastName = result.data.UserInfo.LastName
        profileData.userInfo.userName = result.data.UserInfo.UserName
        profileData.userInfo.email = result.data.UserInfo.Email
        profileData.userInfo.dob = result.data.UserInfo.DOB
        profileData.userInfo.avatar = result.data.UserInfo.Avatar
        profileData.userInfo.isPrivate = result.data.UserInfo.IsPrivate

        profileData.numOfFollowers = result.data.NumOfFollowers;
        profileData.NumOfPosts = result.data.NumOfPosts;
        profileData.NumOfFollowing = result.data.NumOfFollowing;

        profileData.about.bio = result.data.About.Bio
        profileData.about.work = result.data.About.Work
        profileData.about.education = result.data.About.Education
        profileData.about.travel = result.data.About.Travel
        profileData.about.intrests = result.data.About.Intrests
        profileData.about.hobbies = result.data.About.Hobbies
        profileData.about.website = result.data.About.Website
        profileData.about.linkedin = result.data.About.Linkedin
        profileData.about.instgram = result.data.About.Instgram
        profileData.about.twitter = result.data.About.Twitter
    } else {
        profileData.userInfo.firstName = result.data.UserInfo.FirstName
        profileData.userInfo.lastName = result.data.UserInfo.LastName
        profileData.userInfo.avatar = result.data.UserInfo.Avatar
        profileData.userInfo.isPrivate = result.data.UserInfo.IsPrivate
        profileData.about.bio = result.data.About.Bio
        profileData.numOfFollowers = result.data.NumOfFollowers;
        profileData.NumOfPosts = result.data.NumOfPosts;
        profileData.NumOfFollowing = result.data.NumOfFollowing;
    }


    // get followers
    try {
        const result = await getFollowers(id, count)
        if (!result.status) {
            throw new Error("could not get user data")
        }
        profileData.followers = result.data;
    } catch (err) {
        throw new Error("network error, could not connet to server")
    }
    try {
        const result = await getFollowing(id, count)
        if (!result.status) {
            throw new Error("could not get user data")
        }
        profileData.following = result.data;
    } catch (err) {
        throw new Error("network error, could not connet to server")
    }
    console.log(profileData)
}

export async function requestFollow(id, method) {
    const resp = await fetch(`/api/profile/follow?targetid=${id}`, {
        method: method,
        credentials: 'include'
    });

    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return;
    }

    if (!resp.ok) {
        throw new Error("network error, could not connet to server")
    }

    const result = await resp.json();
    return result
}

export async function getFollowers(id, count) {
    const resp = await fetch(`/api/profile/follow?targetid=${id}&count=${count}`, {
        method: "GET",
        credentials: 'include'
    });

    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return;
    }

    if (!resp.ok) {
        throw new Error('could not connect to network')
    }

    const result = await resp.json();
    return result;
}

export async function getFollowing(id, count) {
    const resp = await fetch(`/api/profile/following?targetid=${id}&count=${count}`, {
        method: "GET",
        credentials: 'include'
    });

    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return;
    }

    if (!resp.ok) {
        throw new Error('could not connect to network')
    }

    const result = await resp.json();
    return result;
}