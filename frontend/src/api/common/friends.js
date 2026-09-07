export async function getFriends(searchValue = "") {
    const resp = await fetch(`/api/friends?search=${searchValue}`, {
        method: "GET",
        credentials: 'include'
    });

    if (!resp.ok) {
        throw new Error(`could not fetch data`)
    }

    const result = await resp.json();
    if (!result.status) {
        throw new Error(`could not fetch data`)
    }
    return result.data;
}

export async function addPostGroup(data = {}) {
    const resp = await fetch("/api/post/groups", {
        method: "POST",
        credentials: 'include',
        body: JSON.stringify(data)
    });

    if (!resp.ok) {
        throw new Error('could no connect to server')
    }

    const result = await resp.json();
    return result;
}

export async function getPostGroups() {
    const resp = await fetch("/api/post/groups", {
        method: "GET",
        credentials: 'include'
    });

    if (!resp.ok) {
        throw new Error('could no connect to server')
    }

    const result = await resp.json();
    return result.data;
}

export async function updatePostGroup() {
    
}

export async function deletePostGroup(groupID) {
    const resp = await fetch("/api/post/groups", {
        method: "DELETE",
        credentials: 'include',
        body: JSON.stringify({
            GroupID: groupID
        })
    });
    
    if (!resp.ok) {
        throw new Error("could not connect to server")
    }

    const result = await resp.json();
    return result;
}