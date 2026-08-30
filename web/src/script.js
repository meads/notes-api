import { fetchClient } from "./auth_client";

const getErrorData = async (fetchResponse) => {
    let errorData = null;
    const contentType = fetchResponse.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        errorData = await fetchResponse.json();
        errorData = errorData.error;
    } else {
        errorData = await fetchResponse.text();
    }
    
    return errorData
}

export async function login(username, password) {
    const response = await fetch('http://localhost:8080/login/', {
        mode:'cors',
        credentials: 'include',
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ username, password }),
    })

    if (response.ok) {
        const data = await response.json()

        sessionStorage.setItem("sessionId", data.sessionId)
        sessionStorage.setItem("userId", data.userId)

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}

export async function register(username, password) {
    const response = await fetch('http://localhost:8080/register/', {
        mode:'cors',
        credentials: 'include',
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ username, password }),
    })

    if (response.ok) {
        const data = await response.json()

        sessionStorage.setItem("sessionId", data.sessionId)
        sessionStorage.setItem("userId", data.userId)

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}

export async function logout() {
    let sessionId = sessionStorage.getItem("sessionId")
    const response = await fetch(`http://localhost:8080/logout/${sessionId}`, {
        mode:'cors',
        method: 'POST',
    })

    if (response.ok) { 
        const data = await response.text()

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}

export async function createNote(title, content) {
    const userIdString = sessionStorage.getItem("userId");
    const userId = parseInt(userIdString, 10);
    const response = await fetchClient(`http://localhost:8080/users/${userIdString}/notes/`, {
        mode:'cors',
        credentials: 'include',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title, content, userId }),
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}

export async function updateNote(id, title, content) {
    const userIdString = sessionStorage.getItem("userId");
    const userId = parseInt(userIdString, 10);
    const response = await fetchClient(`http://localhost:8080/users/${userIdString}/notes/`, {
        mode:'cors',
        credentials: 'include',
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id, title, content, userId }),
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}

export async function deleteNote(id) {
    const userIdString = sessionStorage.getItem("userId");
    const response = await fetchClient(`http://localhost:8080/users/${userIdString}/notes/${id}`, {
        mode:'cors',
        credentials: 'include',
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}

export async function getNotes() {
    const userIdString = sessionStorage.getItem("userId");
    const response = await fetchClient(`http://localhost:8080/users/${userIdString}/notes/`, {
        mode:'cors',
        credentials: 'include',
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    })

    if (response.ok) {
        const data = await response.json()

        return { success: true, data: data };
    }
    
    const errorData = await getErrorData(response)
    return { success: false, data: errorData }
}
