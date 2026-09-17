import { checkSessionResponse } from "@/helpers/auth/auth";
import { router } from "@/router/router";
import { disconnectRealtime } from "@/services/realtime";

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
        disconnectRealtime()
        router.push("/login");
        return
    }

    if (!resp.ok && resp.status != 401) {
        throw new Error("could not logout")
        
    }

    disconnectRealtime()
    router.push("/login")
}

export async function authorizeSession() {
    const resp = await fetch("/api/session", {
        method: "GET",
        credentials: "include"
    });

    const result = await resp.json();

    return result;
}
