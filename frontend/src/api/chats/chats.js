export async function getPrivateChatsLists(offset) {
    const resp = await fetch(`/api/groups?offset=${offset}&private=1`, {
        method: "GET",
        credentials: 'include'
    });

    const result = await resp.json();
    if (!resp.ok) {
        throw new Error(result.message || "coud not get data")
    }

    return result;
}

export async function searchChats(offset, searchValue) {
    const resp = await fetch(`/api/groups/search?offset=${offset}&private=1&search=${searchValue}`, {
        method: "GET",
        credentials: 'include'
    });

    const result = await resp.json();
    if (!resp.ok) {
        throw new Error(result.message || "coud not get data")
    }

    console.log(result)
    return result;
}

export async function sendMessage(data) {
    const resp = await fetch(`/api/chats`, {
        method: "POST",
        credentials: 'include',
        body: JSON.stringify(data)
    });
    
    const result = await resp.json();
    if (!resp.ok) {
        throw new Error(result.message || "coud not get data")
    }

    console.log(result)
    return result;
}