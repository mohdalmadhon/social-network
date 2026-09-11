export async function getNotifications(offset = 0) {
    const response = await fetch(`/api/notifications?offset=${offset}`, {
        credentials: 'include'
    });

    return await response.json();
}

export async function acceptFollowRequest(userID) {
    const response = await fetch(
        `/api/follow/accept?targetid=${userID}`,
        {
            method: 'POST',
            credentials: 'include'
        }
    );

    return await response.json();
}