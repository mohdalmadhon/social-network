
export async function registerUser(userData) {
    const resp = await fetch("/api/user", {
        method: "POST",
        body: userData
    })

    const result = await resp.json()

    if (!resp.ok) {
        throw new Error(result.message || `Registration failed: ${resp.status}`)
    }

    return result
}