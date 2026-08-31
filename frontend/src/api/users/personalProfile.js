import { profileData } from "@/data/usersData";
import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";

export async function getUserData() {

    const resp = await fetch("/api/user", {
        method: "GET",
        credentials: 'include'
    });

    const result = await resp.json()

    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return;
    }
    if (!resp.ok) {
        throw new Error(result.message || `Registration failed: ${resp.status}`)
    }

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
    profileData.followers = result.data.Followers;
    profileData.following = result.data.Following;
}
