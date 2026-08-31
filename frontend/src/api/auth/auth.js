import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";

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

export async function loggingSession(userLogger) {
    const resp = await fetch("/api/session", {
        method: "POST",
        credentials: 'include',
        body: JSON.stringify(userLogger)
    });

    const result = await resp.json()

    if (!resp.ok) {
        throw new Error(result.message || `Logging failed: ${resp.status}`)
    }

    return result
}

export async function logout() {
    const resp = await fetch("/api/session", {
        method: "DELETE",
        credentials: 'include'
    });
    
    if(!checkSessionResponse(resp)) {
        router.push("/login");
    }

    if (!resp.ok && resp.status != 401) {
        throw new Error("could not logout")
        
    }

    router.push("/login")
}