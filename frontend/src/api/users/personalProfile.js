import { profileData } from "@/data/usersData";
import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";

export async function getUserData() {
    const resp = await fetch("/api/user", {
        method: "GET",
        credentials: 'include'
    });

    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return;
    }
    const result = await resp.json().catch(() => null)
    if (!result) {
        throw new Error('Could not load your profile')
    }
    if (!resp.ok) {
        throw new Error(result.message || `Registration failed: ${resp.status}`)
    }
    profileData.userInfo.id = result.data.UserInfo.ID
    profileData.userInfo.firstName = result.data.UserInfo.FirstName
    profileData.userInfo.lastName = result.data.UserInfo.LastName
    profileData.userInfo.userName = result.data.UserInfo.UserName
    profileData.userInfo.email = result.data.UserInfo.Email
    profileData.userInfo.dob = result.data.UserInfo.DOB
    profileData.userInfo.avatar = result.data.UserInfo.Avatar
    profileData.userInfo.about = result.data.About.Bio
    profileData.userInfo.isPrivate = Number(result.data.IsPrivate) === 1 ? 1 : 0

    profileData.numOfFollowers = result.data.NumOfFollowers;
    profileData.numOfPosts = result.data.NumOfPosts;
    profileData.numOfFollowing = result.data.NumOfFollowing;

    profileData.about.bio = result.data.About.Bio
    profileData.about.work = result.data.About.Work
    profileData.about.education = result.data.About.Education
    profileData.about.travel = result.data.About.Travel
    profileData.about.intrests = result.data.About.interests
    profileData.about.hobbies = result.data.About.Hobbies
    profileData.about.website = result.data.About.Website
    profileData.about.linkedin = result.data.About.Linkedin
    profileData.about.instgram = result.data.About.instagram
    profileData.about.twitter = result.data.About.Twitter
    profileData.followers = result.data.Followers;
    profileData.following = result.data.Following;
    profileData.friends = result.data.Friends;
    return result
}
