import { Users, Profile, About } from "@/models/users";
import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";

export async function getUserData() {
    const resp = await fetch("/api/user", {
        method: "GET",
        credentials: "include"
    });

<<<<<<< HEAD
    const result = await resp.json();

=======
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
    if (!checkSessionResponse(resp)) {
        router.replace("/login");
        return null;
    }
<<<<<<< HEAD

=======
    const result = await resp.json().catch(() => null)
    if (!result) {
        throw new Error('Could not load your profile')
    }
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
    if (!resp.ok) {
        throw new Error(result.message || `Failed to get user data: ${resp.status}`);
    }
<<<<<<< HEAD

    const data = result.data;
    const user = data.UserInfo;

    const about = new About(
        data.About.Bio,
        data.About.Work,
        data.About.Education,
        data.About.Travel,
        data.About.interests,
        data.About.Hobbies,
        data.About.Website,
        data.About.Linkedin,
        data.About.instagram,
        data.About.Twitter
    );

    const profile = new Profile(
        data.NumOfFollowing,
        data.NumOfFollowers,
        data.NumOfPosts,
        user.avatar,
        about,
        data.Following,
        data.Followers,
        data.Friends
    );

    
    const userData = new Users(
        user.ID,
        user.firstName,
        user.lastName,
        user.username,
        user.email,
        user.DOB,
        user.isPrivate,
        profile
    );
    console.log(userData)
    return userData;

}
=======
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
>>>>>>> 10a41907a21c0baac510b189cc5b0eead2b57f53
