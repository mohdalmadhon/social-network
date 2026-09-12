import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";

export async function getGroups() {
  const response = await fetch("/api/groups", {
    method: "GET",
    credentials: "include",
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  if (!response.ok) {
    throw new Error("Error: Could not get groups");
  }

  const result = await response.json();
  return result;
}

export async function createGroupApi(groupData) {
  const response = await fetch("/api/groups", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(groupData),
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  if (!response.ok) {
    throw new Error("Error: Could not create a new group");
  }

  const result = await response.json();
  return result;
}

export async function getGroup(groupID) {
  const response = await fetch(`/api/groups/${groupID}`, {
    method: "GET",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  if (!response.ok) {
    throw new Error("Error: Could not get a group");
  }

  const result = await response.json();
  return result;
}

export async function deleteGroupApi(groupID) {
  const response = await fetch(`/api/groups/${groupID}`, {
    method: "DELETE",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  if (!response.ok) {
    throw new Error("Error: Could not delete a group");
  }

  const result = await response.json();
  return result;
}

export async function groupJoinRequest(groupID) {
  const response = await fetch(`/api/groups/${groupID}/join-requests`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  if (!response.ok) {
    throw new Error("Error: Could not create a group join request");
  }

  const result = await response.json();
  return result;
}

export async function undoJoinGroup(groupID) {
  const response = await fetch(`/api/groups/${groupID}/join-requests`, {
    method: "DELETE",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  if (!response.ok) {
    throw new Error("Error: Could not undo a group join request");
  }

  const result = await response.json();
  return result;
}

async function groupContentRequest(url, options = {}) {
  const response = await fetch(url, {
    credentials: "include",
    ...options,
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  const result = await response.json();
  if (!response.ok) {
    throw new Error(result.message || "Could not load group content");
  }

  return result;
}

export function getGroupPosts(groupID) {
  return groupContentRequest(`/api/groups/${groupID}/posts`);
}

export function createGroupPost(groupID, formData) {
  return groupContentRequest(`/api/groups/${groupID}/posts`, {
    method: "POST",
    body: formData,
  });
}

export function getGroupPostComments(groupID, postID) {
  return groupContentRequest(`/api/groups/${groupID}/posts/${postID}/comments`);
}

export function createGroupPostComment(groupID, postID, formData) {
  return groupContentRequest(`/api/groups/${groupID}/posts/${postID}/comments`, {
    method: "POST",
    body: formData,
  });
}

export function deleteGroupPost(groupID, postID) {
  return groupContentRequest(`/api/groups/${groupID}/posts/${postID}`, {
    method: "DELETE",
  });
}

export function deleteGroupPostComment(groupID, postID, commentID) {
  return groupContentRequest(`/api/groups/${groupID}/posts/${postID}/comments/${commentID}`, {
    method: "DELETE",
  });
}
