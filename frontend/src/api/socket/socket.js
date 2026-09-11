let ws = null;

export function connectToWS() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        return ws;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

    ws = new WebSocket(`${protocol}//${window.location.host}/api/ws`);

    ws.onopen = () => {
        console.log("websocket connected");
    };

    ws.onmessage = (event) => {
        const payload = JSON.parse(event.data);

        switch (payload.type) {
            case "notification":
                console.log("Notification:", payload.data);
                break;

            case "message":
                console.log("Message:", payload.data);
                break;
        }
    };

    ws.onclose = () => {
        console.log("websocket disconnected");
        ws = null;
    };

    ws.onerror = (error) => {
        console.error("websocket error:", error);
    };

    return ws;
}

export function sendWS(payload) {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        return;
    }

    ws.send(JSON.stringify(payload));
}