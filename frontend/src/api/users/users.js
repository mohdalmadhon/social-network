export async function getUsers() {
  const response = await fetch("/api/users", {
    credentials: "include",
  });

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Could not load users");
  }

  return result;
}
