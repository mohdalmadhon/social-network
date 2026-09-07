export async function addPost(post = {}, image = null) {
    const formData = new FormData();

    formData.append('content', post.content || '');
    formData.append('allowComments', String(post.allowComments));
    formData.append('privatePost', String(post.privatePost));
    formData.append('groupID', String(post.groupID));
    formData.append('location', post.location || '');
    formData.append(
        'taggedPeople',
        JSON.stringify(post.taggedPeople || [])
    );

    if (image) {
        formData.append('image', image);
    }

    const resp = await fetch('/api/post', {
        method: 'POST',
        credentials: 'include',
        body: formData
    });

    if (!resp.ok) {
        throw new Error('could not send post');
    }

    return await resp.json();
}