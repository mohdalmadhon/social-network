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

export async function inviteUserToGroup(groupId, userId) {
  const response = await fetch(`/api/groups/${groupId}/invitations`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      userId,
    }),
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Could not send invitation");
  }

  return result;
}

export async function getInviteUsers(groupId) {
  const response = await fetch(`/api/groups/${groupId}/invite-users`, {
    method: "GET",
    credentials: "include",
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Could not load invite users");
  }

  return result;
}

export async function undoGroupInvitation(groupId, invitationId) {
  const response = await fetch(`/api/groups/${groupId}/invitations/${invitationId}`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!checkSessionResponse(response)) {
    router.replace("/login");
    return;
  }

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Could not undo invitation");
  }

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

export function getGroupPosts(groupId) {
  return groupContentRequest(`/api/groups/${groupId}/posts`);
}

export function createGroupPost(groupId, formData) {
  return groupContentRequest(`/api/groups/${groupId}/posts`, {
    method: "POST",
    body: formData,
  });
}

export function getGroupPostComments(groupId, postId) {
  return groupContentRequest(`/api/groups/${groupId}/posts/${postId}/comments`);
}

export function createGroupPostComment(groupId, postId, formData) {
  return groupContentRequest(`/api/groups/${groupId}/posts/${postId}/comments`, {
    method: "POST",
    body: formData,
  });
}

export function deleteGroupPost(groupId, postId) {
  return groupContentRequest(`/api/groups/${groupId}/posts/${postId}`, {
    method: "DELETE",
  });
}

export function deleteGroupPostComment(groupId, postId, commentId) {
  return groupContentRequest(`/api/groups/${groupId}/posts/${postId}/comments/${commentId}`, {
    method: "DELETE",
  });
}
